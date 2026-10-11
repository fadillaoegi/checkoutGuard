#!/usr/bin/env python3
"""Exercise hosted formatting without credentials or network calls."""
import ast
import json
import re
import unittest
from pathlib import Path

source = ast.parse(Path(__file__).with_name('read-traces.py').read_text())
helpers = ast.Module(body=[node for node in source.body if isinstance(node, ast.FunctionDef)
                          and node.name in ('decoded_object', 'model_calls_from_output', 'model_diagnostics')], type_ignores=[])
scope = {'json': json, 're': re}
exec(compile(helpers, 'read-traces.py', 'exec'), scope)


class DiagnosticsReadbackTest(unittest.TestCase):
    def test_structured_and_hosted_text_preserve_types(self):
        calls = [{'step_number': 1, 'diagnostics': {'span_id': '1234567890abcdef',
                  'requested_model': 'requested', 'timeout_ms': 45000,
                  'returned_model_known': False, 'error_kind': 'timeout',
                  'http_response_received': False}}]
        encoded = 'checkoutguard.model_calls.v1=' + json.dumps(calls)
        for value in ({'model_calls': calls}, {'model_calls': encoded},
                      'Findings: none\nModel Calls: ' + encoded + '\nStatus: passed'):
            self.assertEqual(scope['model_calls_from_output'](value), calls)

    def test_malformed_and_unrelated_output_is_rejected(self):
        for value in ('Model Calls: invalid', 'Unrelated: [{"span_id":"123"}]',
                      'Model Calls: checkoutguard.model_calls.v1=[invalid]', None):
            self.assertEqual(scope['model_calls_from_output'](value), [])

    def test_only_correlated_diagnostics_are_exported(self):
        span = {'span_id': '123', 'span_metadata': {}}
        calls = {'123': {'span_id': '123', 'requested_model': 'requested',
                        'timeout_ms': 45000, 'raw_prompt': 'private prompt'}}
        result = scope['model_diagnostics'](span, calls)
        self.assertEqual(result['span_id'], '123')
        self.assertNotIn('raw_prompt', result)
        self.assertIsNone(scope['model_diagnostics']({'span_id': 'other'}, calls))


if __name__ == '__main__':
    unittest.main()
