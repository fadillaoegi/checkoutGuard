#!/usr/bin/env python3
"""Read finalized application traces through the SDK's authenticated v3 route.

Only safe structural metadata is printed/saved, never credentials or raw payloads.
The diagnostic-read marker is an HTTP read header; application spans are unchanged.
"""
import argparse
import json
import os
import re
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path


def load_env(path):
    values = {}
    for line in Path(path).read_text().splitlines():
        line = line.strip().removeprefix('export ')
        if not line or line.startswith('#') or '=' not in line:
            continue
        key, value = line.split('=', 1)
        values[key.strip()] = value.strip().strip('\"\'')
    return values


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *args, **kwargs):
        return None


parser = argparse.ArgumentParser()
parser.add_argument('trace_ids', nargs='+')
parser.add_argument('--output', default='tmp/trace-readback.json')
parser.add_argument('--run-file', action='append', help='Local run JSON or evaluation report for usage comparison; repeatable')
parser.add_argument('--expected-status', choices=['success', 'error'], default='success', help='Expected finalized workflow status; telemetry verification does not mean checkout passed')
args = parser.parse_args()
workspace = Path(__file__).resolve().parent.parent
local_config = workspace / '.env.neatlogs.local'
default_config = local_config if local_config.is_file() else workspace.parent / '.neatlogs-api.env'
config = load_env(os.environ.get('NEATLOGS_ENV_FILE', str(default_config)))
endpoint = os.environ.get('NEATLOGS_ENDPOINT', config.get('NEATLOGS_ENDPOINT', 'https://ingest.neatlogs.com')).rstrip('/')
key = os.environ.get('NEATLOGS_API_KEY', config.get('NEATLOGS_API_KEY', ''))
origin = urllib.parse.urlparse(endpoint)
if origin.scheme != 'https' or origin.hostname not in ('ingest.neatlogs.com', 'eu.ingest.neatlogs.com'):
    raise SystemExit('Unsupported telemetry read origin')
opener = urllib.request.build_opener(NoRedirect)
records = []
local_runs = []
for filename in args.run_file or ['tmp/live-evaluation.json', 'tmp/latest-run.json']:
    path = Path(filename)
    if path.is_file():
        report = json.loads(path.read_text())
        local_runs.extend(report.get('runs', []) if 'runs' in report else [report])
local_usage = {run['trace_id']: run['usage'] for run in local_runs if run.get('trace_id')}
for trace_id in args.trace_ids:
    if not re.fullmatch('[a-f0-9]{32}', trace_id):
        raise SystemExit('Invalid trace ID')
    request = urllib.request.Request(endpoint + '/api/traces/v3/' + trace_id, headers={'x-api-key': key, 'x-neatlogs-doctor': 'v1'})
    try:
        with opener.open(request, timeout=15) as response:
            data = json.load(response)
        spans = data.get('spans', [])
        ids = {span.get('span_id') for span in spans}
        roots = [span for span in spans if not span.get('parent_span_id')]
        llm = [span for span in spans if str(span.get('node_type', span.get('span_type', ''))).lower() == 'llm']
        summaries = []
        for span in spans:
            io = span.get('data', {}) or {}
            summaries.append({'name': span.get('node_name', span.get('span_name')), 'kind': span.get('node_type', span.get('span_type')), 'input_present': io.get('input_value') is not None, 'output_present': io.get('output_value') is not None})
        records.append({'trace_id': trace_id, 'returned_id': data.get('_id'), 'status': data.get('status'), 'span_count': len(spans), 'reported_span_count': data.get('spanCount'), 'root_count': len(roots), 'hierarchy_valid': all(not span.get('parent_span_id') or span['parent_span_id'] in ids for span in spans), 'duplicates': len(spans) - len(ids), 'llm_count': len(llm), 'input_tokens': data.get('promptTokens'), 'output_tokens': data.get('completionTokens'), 'total_tokens': data.get('totalTokensUsed'), 'spans': summaries})
    except urllib.error.HTTPError as error:
        records.append({'trace_id': trace_id, 'http_status': error.code, 'readback_verified': False})
    except (OSError, ValueError):
        records.append({'trace_id': trace_id, 'readback_verified': False, 'error': 'read request failed'})
for record in records:
    record['expected_status'] = args.expected_status
    usage = local_usage.get(record['trace_id'])
    record['local_usage_matches'] = (record.get('input_tokens') == usage['input_tokens'] and record.get('output_tokens') == usage['output_tokens']) if usage else None
    tokens = [record.get(name) for name in ('input_tokens', 'output_tokens', 'total_tokens')]
    record['verified'] = (record.get('returned_id') == record['trace_id'] and record.get('status') == args.expected_status and record.get('root_count') == 1 and record.get('duplicates') == 0 and record.get('hierarchy_valid') is True and record.get('span_count') == record.get('reported_span_count') and record.get('llm_count', 0) > 0 and all(span.get('input_present') and span.get('output_present') for span in record.get('spans', [])) and all(isinstance(value, (int, float)) and value >= 0 for value in tokens) and tokens[0] + tokens[1] == tokens[2] and record['local_usage_matches'] is not False)
path = Path(args.output)
path.parent.mkdir(parents=True, exist_ok=True)
path.write_text(json.dumps(records, indent=2))
print(json.dumps(records, indent=2))

if not all(record['verified'] for record in records):
    raise SystemExit(1)
