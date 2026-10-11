#!/usr/bin/env python3
"""Read-only checks of the original cfo.ai XLSX against independent arithmetic.

No workbook rewriting, native Excel recalculation, Git or network calls.
The 2026-10 planning seed and all customer forecasts are assumptions.
"""
import argparse
import hashlib
import importlib.util
import json
import math
import re
import xml.etree.ElementTree as ET
import zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
NS = {'s': 'http://schemas.openxmlformats.org/spreadsheetml/2006/main',
      'p': 'http://schemas.openxmlformats.org/package/2006/relationships'}
spec = importlib.util.spec_from_file_location('publication_audit', ROOT / 'scripts/audit-publication.py')
audit = importlib.util.module_from_spec(spec)
spec.loader.exec_module(audit)


def workbook(path):
    sheets, error_count, privacy_count, external_count, formulas = {}, 0, 0, 0, 0
    secrets = audit.known_secrets(ROOT)
    identities = re.compile(rb'[\w.+-]+@[\w.-]+\.[A-Za-z]{2,}|https?://app\.(?:cfo\.ai|neatlogs\.com)/|/Users/[^/]+/')
    with zipfile.ZipFile(path) as archive:
        if sum(x.file_size for x in archive.infolist()) > 20 << 20:
            raise ValueError('Workbook too large to review')
        for name in archive.namelist():
            if 'vbaProject' in name or 'externalLink' in name:
                external_count += 1
            data = archive.read(name)
            if identities.search(data) or any(p.search(data) for p in audit.PATTERNS.values()) or any(s in data for s in secrets):
                privacy_count += 1
            if name.endswith('.rels'):
                external_count += sum(r.get('TargetMode') == 'External' for r in ET.fromstring(data))
        shared = []
        if 'xl/sharedStrings.xml' in archive.namelist():
            shared = [''.join(t.text or '' for t in si.iterfind('.//s:t', NS))
                      for si in ET.fromstring(archive.read('xl/sharedStrings.xml'))]
        rels = {r.get('Id'): r.get('Target') for r in ET.fromstring(archive.read('xl/_rels/workbook.xml.rels'))}
        for sheet in ET.fromstring(archive.read('xl/workbook.xml')).findall('s:sheets/s:sheet', NS):
            target = rels[sheet.get('{http://schemas.openxmlformats.org/officeDocument/2006/relationships}id')]
            member = target.lstrip('/') if target.startswith('/') else 'xl/' + target
            cells = {}
            for cell in ET.fromstring(archive.read(member)).iterfind('.//s:c', NS):
                kind = cell.get('t')
                value = cell.find('s:v', NS)
                if kind == 'inlineStr':
                    parsed = ''.join(t.text or '' for t in cell.iterfind('.//s:t', NS))
                elif value is None or value.text is None:
                    continue
                elif kind == 's':
                    parsed = shared[int(value.text)]
                elif kind in ('str', 'e'):
                    parsed = value.text
                else:
                    parsed = float(value.text)
                error_count += kind == 'e'
                formulas += cell.find('s:f', NS) is not None
                cells[cell.get('r')] = parsed
            sheets[sheet.get('name')] = cells
    if privacy_count or external_count or error_count:
        raise ValueError('Workbook privacy, external-content or error-cell check failed; values withheld')
    return sheets, {'sheets': len(sheets), 'cached_formula_cells': formulas,
                    'excel_error_cells': error_count, 'matched_private_or_secret_members': privacy_count,
                    'external_relationships_or_macro_members': external_count}


def column(number):
    result = ''
    while number:
        number, rem = divmod(number - 1, 26)
        result = chr(65 + rem) + result
    return result


def review(path):
    original = json.loads((ROOT / 'submission/business/cfoai-original-manifest.json').read_text())
    if hashlib.sha256(path.read_bytes()).hexdigest() != original['export_sha256']:
        raise ValueError('Workbook differs from the supplied original export; review provenance before replacing it')
    w, checks = workbook(path)
    a = w['CheckoutGuard — Assumptions · C']
    plans = w['CheckoutGuard — Assumptions · W']
    observations = w['CheckoutGuard — Business plan ·']
    costs = w['CheckoutGuard — Assumptions (2)']
    native = w['Model inputs']
    comparisons = 0

    def equal(actual, expected, label):
        nonlocal comparisons
        if not isinstance(actual, (int, float)) or not math.isclose(actual, expected, rel_tol=1e-10, abs_tol=1e-7):
            raise ValueError('Independent calculation differs at ' + label)
        comparisons += 1

    runs = json.loads((ROOT / 'submission/evidence/runs.json').read_text())['runs']
    linked = []
    estimates = []
    for row in range(5, 8):
        matches = [r for r in runs if r['reported_usage']['input_tokens'] == observations[f'B{row}']]
        if len(matches) != 1:
            raise ValueError('Runtime observation cannot be uniquely linked to curated evidence')
        run = matches[0]
        equal(observations[f'C{row}'], run['reported_usage']['output_tokens'], 'runtime output tokens')
        equal(observations[f'D{row}'], run['duration_seconds'], 'runtime duration')
        equal(observations[f'F{row}'], run['verified_png_count'], 'runtime screenshots')
        estimate = (observations[f'B{row}'] * a['B5'] + observations[f'C{row}'] * a['B6']) / 1_000_000
        equal(observations[f'G{row}'], estimate, 'catalog token estimate')
        estimates.append(estimate)
        linked.append({'workbook_row': row, 'case': run['case'], 'run_id': run['run_id']})
    mean_token = sum(estimates) / len(estimates)

    def run_cost(provider=1, extra=None, browser=1, support=1):
        extra = a['B8'] if extra is None else extra
        return (mean_token * a['B7'] * provider * (1 + extra * a['B9'])
                + a['B10'] / 3600 * a['B11'] * browser * (1 + extra)
                + a['B12'] + a['B13'] + a['B14'] + a['B15'] / 60 * a['B16'] * support)

    cost = run_cost()
    fixed = sum(a[f'B{r}'] for r in range(20, 26))
    equal(costs['B8'], mean_token, 'known token estimate')
    equal(costs['B12'], cost, 'all-in variable cost')
    equal(costs['B13'], fixed, 'fixed monthly cost')
    settings = [(plans[f'B{r}'], plans[f'C{r}'], plans[f'E{r}'], plans[f'F{r}']) for r in range(5, 8)]
    sensitivities = w['CheckoutGuard — Assumptions · S']
    fee_rate, fee_flat = sensitivities['B5'], sensitivities['B6']

    def economics(counts, unit_cost=cost, utilization=a['B17'], support=1):
        revenue = sum(n * price for n, (price, _, _, _) in zip(counts, settings))
        completed = sum(n * quota * utilization for n, (_, quota, _, _) in zip(counts, settings))
        support_cost = sum(n * minutes / 60 * a['B16'] * support
                           for n, (_, _, minutes, _) in zip(counts, settings))
        fees = sum(n * (price * fee_rate + fee_flat) for n, (price, _, _, _) in zip(counts, settings))
        return revenue, revenue - completed * unit_cost - support_cost - fees

    pricing = w['CheckoutGuard — Business pl (2)']
    for row in range(5, 8):
        price, quota, minutes, _ = settings[row - 5]
        contribution = price - quota * a['B17'] * cost - minutes / 60 * a['B16'] - price * fee_rate - fee_flat
        full = contribution - quota * (1 - a['B17']) * cost
        equal(pricing[f'I{row}'], contribution, 'subscription contribution')
        equal(pricing[f'J{row}'], contribution / price, 'subscription margin')
        equal(pricing[f'K{row}'], full, 'full-quota contribution')
        equal(pricing[f'L{row}'], full / price, 'full-quota margin')

    cases = {
        'base': ([plans[f'G{r}'] for r in range(5, 8)], a['B18'], a['B19'], [6, 84, 19, 68, 31]),
        'conservative': ([native[f'AB{r}'] for r in range(293, 296)], native['AB382'], native['AB287'], [10, 88, 23, 69, 32]),
        'growth': ([native[f'AB{r}'] for r in range(297, 300)], native['AB383'], native['AB288'], [14, 92, 27, 70, 33]),
    }
    results = {}
    base_final = None
    for name, (conversion, growth, churn, rows) in cases.items():
        counts, cumulative, peak, first, snapshots = [0., 0., 0.], 0., 0., None, {}
        for month in range(25):
            if month:
                counts = [prior * (1 - churn) + prospects * (1 + growth) ** (month - 1) * rate
                          for prior, (_, _, _, prospects), rate in zip(counts, settings, conversion)]
            revenue, contribution = economics(counts)
            operating = contribution - fixed if month else 0.
            cumulative += operating
            peak = max(peak, -cumulative)
            date = f'{2026 + (9 + month) // 12}-{(9 + month) % 12 + 1:02}'
            if month and operating >= 0 and first is None:
                first = date
            for r, value in zip(rows, [sum(counts), revenue, contribution, operating, cumulative]):
                cell = column(month + 3) + str(r)
                equal(native[cell], value, 'Model inputs!' + cell)
            if month in (12, 24):
                snapshots[f'm{month}'] = {'revenue_usd': revenue, 'subscriptions_expected': sum(counts),
                                         'operating_result_usd': operating, 'cumulative_operating_result_usd': cumulative}
        results[name] = {'periods_checked': 25, 'first_nonnegative_operating_month': first,
                         'peak_cumulative_operating_deficit_usd': peak, **snapshots}
        if name == 'base':
            base_final = counts
            break_even = math.ceil(fixed / (contribution / sum(counts)))
            equal(w['CheckoutGuard — Forecast · Chec']['Z46'], break_even, 'break-even subscriptions')
            break_even_revenue = break_even * revenue / sum(counts)
            equal(w['CheckoutGuard — Forecast · Chec']['Z47'], break_even_revenue, 'break-even revenue')

    sensitivity_results = {}
    changes = {'provider_10x': (6, run_cost(provider=sensitivities['B7']), a['B17'], 1),
               'extra_attempts_100_percent': (7, run_cost(extra=sensitivities['B8']), a['B17'], 1),
               'browser_3x': (8, run_cost(browser=sensitivities['B11']), a['B17'], 1),
               'support_2x': (9, run_cost(support=sensitivities['B10']), a['B17'], sensitivities['B10']),
               'full_quota': (10, cost, a['B17'] * sensitivities['B9'], 1)}
    for name, (row, unit, use, support) in changes.items():
        _, contribution = economics(base_final, unit, use, support)
        value = contribution - fixed
        equal(w['CheckoutGuard — Forecast · One-'][f'Y{row}'], value, 'October 2028 sensitivity')
        sensitivity_results[name] = value
    return {'status': 'pass', 'xlsx_sha256': hashlib.sha256(path.read_bytes()).hexdigest(),
            'checks': {**checks, 'independent_numeric_comparisons': comparisons,
                       'forecast_periods_checked': 75, 'one_way_sensitivities_checked': 5},
            'runtime_evidence_links': linked, 'mean_known_token_estimate_usd': mean_token,
            'variable_cost_per_completed_run_usd': cost, 'fixed_monthly_cost_assumption_usd': fixed,
            'base_break_even_subscriptions': break_even, 'base_break_even_monthly_revenue_usd': break_even_revenue,
            'scenarios': results, 'october_2028_sensitivity_operating_results_usd': sensitivity_results,
            'limits': ['Cached XLSX values compared to independently calculated model; original export unchanged.',
                       'No native Excel recalculation or interactive cfo.ai simulator testing performed by this review.',
                       'Original cfo.ai plan used supplied figures; EVIDENCE.md was not attached there.',
                       'Customer, price, cost and growth figures are planning assumptions, not validated traction or invoices.',
                       'Cash runway is uncalculated; cumulative operating result is not available cash.',
                       'Pattern checks are not a guarantee against every form of sensitive information.'],
            'git_commands_executed': False}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--xlsx', type=Path, default=ROOT / 'checkoutguard_FE/public/business/CheckoutGuard-cfoai-model.xlsx')
    parser.add_argument('--output', type=Path)
    args = parser.parse_args()
    report = review(args.xlsx)
    if args.output:
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(json.dumps(report, indent=2) + '\n')
    print(json.dumps({'status': report['status'], **report['checks']}, indent=2))


if __name__ == '__main__':
    main()
