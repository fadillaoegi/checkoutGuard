#!/usr/bin/env python3
"""Regression checks for publication boundaries; no Git commands or real keys."""
import importlib.util
import struct
import tempfile
import unittest
import zipfile
from pathlib import Path

spec = importlib.util.spec_from_file_location('publication_audit', Path(__file__).with_name('audit-publication.py'))
audit = importlib.util.module_from_spec(spec)
spec.loader.exec_module(audit)

def index_fixture(names, version=2):
    data = b'DIRC' + struct.pack('>II', version, len(names))
    previous = b''
    for value in names:
        name = value.encode()
        entry = bytes(60) + struct.pack('>H', len(name))
        if version == 4:
            common = 0
            while common < min(len(previous), len(name)) and previous[common] == name[common]: common += 1
            trim = len(previous) - common
            if trim >= 128: raise ValueError('test fixture only uses short prefixes')
            entry += bytes([trim]) + name[common:] + b'\0'
        else:
            entry += name + b'\0'
            entry += bytes((-len(entry)) % 8)
        data += entry
        previous = name
    return data + bytes(20)

class PublicationAuditTests(unittest.TestCase):
    def test_compressed_workbook_secret_is_blocked_without_printing_it(self):
        with tempfile.TemporaryDirectory() as name:
            root = Path(name)
            (root / '.gitignore').write_bytes((audit.ROOT / '.gitignore').read_bytes())
            (root / '.env.model.local').write_text('MODEL_API_KEY=synthetic-test-value\n')
            with zipfile.ZipFile(root / 'plan.xlsx', 'w', compression=zipfile.ZIP_DEFLATED) as archive:
                archive.writestr('xl/worksheets/sheet1.xml', '<t>synthetic-test-value</t>')
            report = audit.audit(root)
            self.assertIn({'file': 'plan.xlsx', 'kind': 'known_secret_value'}, report['issues'])
            self.assertNotIn('synthetic-test-value', str(report))

    def test_invalid_workbook_fails_the_audit(self):
        with tempfile.TemporaryDirectory() as name:
            root = Path(name)
            (root / '.gitignore').write_bytes((audit.ROOT / '.gitignore').read_bytes())
            (root / 'plan.xlsx').write_bytes(b'not an archive')
            report = audit.audit(root)
            self.assertIn({'file': 'plan.xlsx', 'kind': 'compressed_file_could_not_be_audited'}, report['issues'])

    def test_ignored_private_files_and_visible_reviewer_evidence(self):
        self.assertEqual(audit.ignore_contract(audit.ROOT), [])

    def test_index_paths_include_ignored_tracked_files(self):
        names = ['.env.model.local', 'README.md', 'submission/evidence/runs.json']
        for version in (2, 3, 4):
            self.assertEqual(audit.index_paths(index_fixture(names, version)), set(names))

    def test_tracked_secret_file_is_blocked_even_when_ignored(self):
        with tempfile.TemporaryDirectory() as name:
            root = Path(name)
            (root / '.git').mkdir()
            (root / '.git/index').write_bytes(index_fixture(['.env.model.local']))
            (root / '.gitignore').write_bytes((audit.ROOT / '.gitignore').read_bytes())
            (root / '.env.model.local').write_text('MODEL_API_KEY=synthetic-test-value\n')
            (root / 'README.md').write_text('Public readme\n')
            report = audit.audit(root)
            self.assertEqual(report['status'], 'fail')
            self.assertIn({'file': '.env.model.local', 'kind': 'private_path_in_public_candidates'}, report['issues'])

    def test_known_secret_value_is_reported_without_printing_it(self):
        with tempfile.TemporaryDirectory() as name:
            root = Path(name)
            (root / '.gitignore').write_bytes((audit.ROOT / '.gitignore').read_bytes())
            (root / '.env.model.local').write_text('MODEL_API_KEY=synthetic-test-value\n')
            (root / 'README.md').write_text('accidental synthetic-test-value\n')
            report = audit.audit(root)
            self.assertIn({'file': 'README.md', 'kind': 'known_secret_value'}, report['issues'])
            self.assertNotIn('synthetic-test-value', str(report))

if __name__ == '__main__': unittest.main()
