#!/usr/bin/env python3
"""Calculate a dated runway under explicitly supplied constant cash-flow assumptions.

Writes financial inputs/results only inside ignored tmp. This is a local
calculation, not an export from cfo.ai. No Git, account or network operations.
"""
import argparse
import calendar
import json
import re
from datetime import date
from decimal import Decimal, InvalidOperation, ROUND_HALF_UP
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
CATEGORIES = ('operating', 'founder', 'tax', 'capex', 'debt_service', 'refunds')


def amount(value, name, missing):
    if value is None:
        missing.append(name)
        return None
    if isinstance(value, bool) or not isinstance(value, (int, float, str)):
        raise ValueError('Invalid amount: ' + name)
    try:
        parsed = Decimal(str(value))
    except InvalidOperation:
        raise ValueError('Invalid amount: ' + name) from None
    if not parsed.is_finite() or parsed < 0:
        raise ValueError('Amount must be finite and nonnegative: ' + name)
    return parsed


def calculate(data, months=24):
    if not isinstance(months, int) or isinstance(months, bool) or not 1 <= months <= 120:
        raise ValueError('Forecast horizon must be 1–120 months')
    missing = []
    basis = data.get('cash_basis')
    if basis is None:
        missing.append('cash_basis')
    elif basis not in ('actual', 'simulation'):
        raise ValueError('cash_basis must be actual or simulation')
    currency = data.get('currency')
    if currency is None:
        missing.append('currency')
    elif not isinstance(currency, str) or len(currency) != 3 or not currency.isalpha() or currency != currency.upper():
        raise ValueError('Use a three-letter currency code; no automatic currency conversion is performed')
    dated = data.get('cash_as_of')
    if dated is None:
        missing.append('cash_as_of')
        as_of = None
    else:
        if not isinstance(dated, str) or not re.fullmatch(r'\d{4}-\d{2}-\d{2}', dated):
            raise ValueError('cash_as_of must be an ISO date YYYY-MM-DD')
        try:
            as_of = date.fromisoformat(dated)
        except (ValueError, TypeError):
            raise ValueError('cash_as_of must be an ISO date YYYY-MM-DD') from None
    cash = amount(data.get('unrestricted_cash'), 'unrestricted_cash', missing)
    buffer = amount(data.get('minimum_cash_buffer'), 'minimum_cash_buffer', missing)
    receipts = amount(data.get('monthly_cash_receipts'), 'monthly_cash_receipts', missing)
    payments = data.get('monthly_cash_payments')
    if payments is None:
        payments = {}
    if not isinstance(payments, dict) or set(payments) - set(CATEGORIES):
        raise ValueError('Use the documented payment categories without overlapping totals')
    values = [amount(payments.get(k), 'monthly_cash_payments.' + k, missing) for k in CATEGORIES]
    if data.get('no_future_financing_assumed') is not True:
        missing.append('no_future_financing_assumed: explicit true required for this flat model')
    if missing:
        return {'status': 'needs_input', 'missing_inputs': missing, 'runway_months': None,
                'source': 'Local calculator; not a cfo.ai export'}
    total_payments = sum(values, Decimal(0))
    burn = total_payments - receipts
    available = max(Decimal(0), cash - buffer)
    if cash <= buffer:
        runway, state = Decimal(0), 'at_or_below_minimum_cash_buffer'
    elif burn > 0:
        runway, state = available / burn, 'depleting_under_flat_assumptions'
    else:
        runway, state = None, 'not_depleting_under_flat_assumptions'

    def shown(value):
        return float(value.quantize(Decimal('.01'), rounding=ROUND_HALF_UP))

    schedule, opening = [], cash
    first = as_of.isoformat() if cash <= buffer else None
    for i in range(1, months + 1):
        # Anniversary periods retain the original day; Jan 31 becomes Feb 28
        # then Mar 31. The model does not infer exact within-period payment dates.
        index = as_of.year * 12 + as_of.month - 1 + i
        year, month_zero = divmod(index, 12)
        end = date(year, month_zero + 1, min(as_of.day, calendar.monthrange(year, month_zero + 1)[1]))
        closing = opening + receipts - total_payments
        schedule.append({'period': i, 'period_end': end.isoformat(), 'opening_cash': shown(opening),
                         'cash_receipts': shown(receipts), 'cash_payments': shown(total_payments),
                         'closing_cash': shown(closing), 'at_or_below_buffer': closing <= buffer})
        if first is None and closing <= buffer:
            first = end.isoformat()
        opening = closing
    return {'status': 'calculated', 'source': 'Local calculator; not a cfo.ai export',
            'cash_basis': basis, 'cash_as_of': as_of.isoformat(), 'currency': currency,
            'unrestricted_cash': shown(cash), 'minimum_cash_buffer': shown(buffer),
            'monthly_cash_receipts': shown(receipts), 'monthly_cash_payments': shown(total_payments),
            'monthly_net_cash_burn': shown(burn), 'runway_months': shown(runway) if runway is not None else None,
            'cash_status': state, 'first_period_end_at_or_below_buffer': first,
            'forecast_periods': months, 'cash_schedule': schedule,
            'assumptions': ['Constant monthly receipts/payments and no future financing.',
                            'Periods run from the supplied balance date to monthly anniversary dates.',
                            'A threshold period-end is not the exact date cash runs out within that period.',
                            'Forecast cash flows are assumptions; no invoice or balance has been independently verified.',
                            'No currency conversion; all inputs must use the same currency.'],
            'git_commands_executed': False}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--input', type=Path, required=True)
    parser.add_argument('--output', type=Path, default=Path('tmp/runway-result.json'))
    parser.add_argument('--months', type=int, default=24)
    args = parser.parse_args()
    output = args.output.resolve()
    if not output.is_relative_to(ROOT / 'tmp'):
        raise ValueError('Financial results must stay inside ignored tmp until the user selects what to publish')
    report = calculate(json.loads(args.input.read_text()), args.months)
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(json.dumps(report, indent=2) + '\n')
    output.chmod(0o600)
    # Keep cash amounts out of automatic tool output. Open the private report
    # when the user is ready to inspect their own inputs and calculated result.
    print(json.dumps({'status': report['status'], 'missing_inputs': report.get('missing_inputs', []),
                      'output': str(output.relative_to(ROOT)), 'git_commands_executed': False}, indent=2))


if __name__ == '__main__':
    main()
