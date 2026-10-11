#!/usr/bin/env python3
"""Export only verified synthetic demo evidence; never raw runs or prompts."""
import hashlib
import json
import shutil
from datetime import datetime
from decimal import Decimal
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
CASES = [
    ('pricing-defect', 'b87b1ae0-83a0-4355-b3ac-850d3fb6a4be', 9800, 11050, 'findings', 21.30),
    ('corrected-retest', '6b111e4c-30ff-44e2-a6cc-0c533eb3af71', 9800, 9800, 'passed', 28.92),
    ('payment-timeout', 'b74d8ccb-9559-4a61-9c52-7eb3ae713815', 12400, 12400, 'passed', 16.66),
    ('priced-retest', '8b2d9a3c-7d16-41fa-8f99-5cba3ea097c1', 9800, 9800, 'passed', 34.98),
]

def main():
    destination = ROOT / 'submission/evidence'
    destination.mkdir(parents=True, exist_ok=True)
    summaries, files = [], []
    for label, run_id, expected, displayed, status, _ in CASES:
        run = json.loads((ROOT / f'tmp/model-reliability/{run_id}.json').read_text())
        trace = json.loads((ROOT / f'tmp/model-reliability/{run_id[:8]}-trace.json').read_text())[0]
        assert run['environment'] == 'sandbox' and run['status'] == status
        assert run['total_verified'] and run['order_count'] == 1
        assert run['expected_total_minor'] == expected and run['displayed_total_minor'] == displayed
        assert trace['trace_id'] == run['trace_id'] and trace['verified'] and trace['local_usage_matches']
        assert sum(s['action']['tool'] == 'click' and s['action'].get('selector') == '#place' for s in run['steps']) == 1
        assert all(s['action']['tool'] != 'model_error' for s in run['steps'])
        for step in run['steps']:
            file = ROOT / 'checkoutguard_BE' / step['evidence'].lstrip('/')
            assert file.read_bytes().startswith(b'\x89PNG\r\n\x1a\n')
        image = ROOT / f'tmp/model-reliability/{run_id[:8]}-final.png'
        final = ROOT / 'checkoutguard_BE' / run['steps'][-1]['evidence'].lstrip('/')
        assert image.read_bytes() == final.read_bytes(), 'Preview must match original evidence exactly'
        target = destination / f'{label}.png'
        shutil.copyfile(image, target)
        target.chmod(0o644)
        amount = None
        if run['provider'] == 'publicai' and run['model'] == 'aisingapore/Qwen-SEA-LION-v4-32B-IT':
            amount = float((Decimal(run['usage']['input_tokens']) * Decimal('0.25') + Decimal(run['usage']['output_tokens']) * Decimal('0.50')) / Decimal(1_000_000))
        summaries.append({
            'case': label, 'run_id': run_id, 'parent_run_id': run.get('parent_id'),
            'scenario': run['scenario'], 'basket': run['variant'], 'status': status,
            'synthetic': True, 'total_verified': True, 'expected_total_idr': expected,
            'displayed_total_idr': displayed, 'persisted_orders': 1, 'submit_actions': 1,
            'duration_seconds': round((datetime.fromisoformat(run['finished_at']) - datetime.fromisoformat(run['started_at'])).total_seconds(), 2),
            'findings': [{'kind': f['kind'], 'verified': f['verified']} for f in run['findings']],
            'provider': run['provider'], 'requested_model': run['requested_model'], 'returned_model_label': run['model'],
            'reported_usage': run['usage'], 'posthoc_known_usage_estimate_usd': amount,
            'posthoc_estimate_basis': 'Public AI catalog, checked 2026-10-11; not a verified invoice',
            'screenshot': target.name, 'verified_png_count': len(run['steps']), 'trace_id': run['trace_id'],
        })
        safe_trace = {key: trace[key] for key in ('trace_id', 'status', 'span_count', 'root_count', 'hierarchy_valid', 'duplicates', 'llm_count', 'input_tokens', 'output_tokens', 'total_tokens', 'local_usage_matches', 'verified')}
        safe_trace['spans'] = [{k: s[k] for k in ('name', 'kind', 'input_present', 'output_present')} for s in trace['spans']]
        trace_file = destination / f'{label}-trace.json'
        trace_file.write_text(json.dumps(safe_trace, indent=2) + '\n')
        files.extend([target, trace_file])
    assert summaries[1]['parent_run_id'] == summaries[0]['run_id']
    runs_file = destination / 'runs.json'
    runs_file.write_text(json.dumps({'generated_from_verified_local_records': True, 'raw_prompts_and_credentials_included': False, 'runs': summaries}, indent=2) + '\n')
    files.append(runs_file)
    overview = ROOT / 'tmp/ui/dashboard-overview.png'
    if overview.is_file():
        assert overview.read_bytes().startswith(b'\x89PNG\r\n\x1a\n')
        preview = destination / 'dashboard.png'
        shutil.copyfile(overview, preview)
        preview.chmod(0o644)
        files.append(preview)
    manifest = [{'file': f.name, 'bytes': f.stat().st_size, 'sha256': hashlib.sha256(f.read_bytes()).hexdigest()} for f in files]
    (destination / 'manifest.json').write_text(json.dumps(manifest, indent=2) + '\n')
    print(json.dumps({'exported_cases': len(summaries), 'public_files': len(files)+1, 'verified_pngs': sum(r['verified_png_count'] for r in summaries), 'destination': 'submission/evidence', 'raw_prompts_and_credentials_included': False}))

if __name__ == '__main__': main()
