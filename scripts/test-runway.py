#!/usr/bin/env python3
"""Check cash-flow edge cases with synthetic financial inputs only."""
import importlib.util
import unittest
from pathlib import Path

spec = importlib.util.spec_from_file_location('runway', Path(__file__).with_name('calculate-runway.py'))
runway = importlib.util.module_from_spec(spec)
spec.loader.exec_module(runway)


def inputs():
    return {'cash_basis': 'simulation', 'cash_as_of': '2026-01-31', 'currency': 'USD',
            'unrestricted_cash': 1000, 'minimum_cash_buffer': 100, 'monthly_cash_receipts': 50,
            'monthly_cash_payments': {'operating': 250, 'founder': 0, 'tax': 0,
                                      'capex': 0, 'debt_service': 0, 'refunds': 0},
            'no_future_financing_assumed': True}


class RunwayTests(unittest.TestCase):
    def test_cash_rollforward_and_buffer_crossing(self):
        result = runway.calculate(inputs())
        self.assertEqual(result['monthly_net_cash_burn'], 200)
        self.assertEqual(result['runway_months'], 4.5)
        self.assertEqual(result['first_period_end_at_or_below_buffer'], '2026-06-30')
        self.assertEqual(result['cash_schedule'][0]['closing_cash'], 800)

    def test_anniversary_periods_do_not_drift_after_february(self):
        result = runway.calculate(inputs())
        self.assertEqual([r['period_end'] for r in result['cash_schedule'][:3]],
                         ['2026-02-28', '2026-03-31', '2026-04-30'])

    def test_missing_inputs_are_not_zero(self):
        data = inputs(); data['unrestricted_cash'] = None
        result = runway.calculate(data)
        self.assertEqual(result['status'], 'needs_input')
        self.assertIsNone(result['runway_months'])

    def test_real_zero_cash_has_zero_runway(self):
        data = inputs(); data['unrestricted_cash'] = 0
        self.assertEqual(runway.calculate(data)['runway_months'], 0)

    def test_nonpositive_burn_is_not_a_fabricated_infinite_runway(self):
        for receipts in (250, 300):
            data = inputs(); data['monthly_cash_receipts'] = receipts
            result = runway.calculate(data)
            self.assertIsNone(result['runway_months'])
            self.assertIsNone(result['first_period_end_at_or_below_buffer'])

    def test_unknown_funding_assumption_requires_input(self):
        data = inputs(); data['no_future_financing_assumed'] = None
        self.assertEqual(runway.calculate(data)['status'], 'needs_input')

    def test_invalid_dates_and_amounts_are_rejected(self):
        for value in (float('nan'), float('inf'), -1, True):
            data = inputs(); data['unrestricted_cash'] = value
            with self.assertRaises(ValueError): runway.calculate(data)
        data = inputs(); data['cash_as_of'] = '2026-02-30'
        with self.assertRaises(ValueError): runway.calculate(data)
        data['cash_as_of'] = '20260131'
        with self.assertRaises(ValueError): runway.calculate(data)


if __name__ == '__main__': unittest.main()
