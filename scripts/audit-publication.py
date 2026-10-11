#!/usr/bin/env python3
"""Audit public candidate files without executing Git or printing secrets.

Uses ripgrep's ignore engine, and reads index filenames to catch private files
already tracked. Historical/staged blob contents are outside this audit.
"""
import argparse
import hashlib
import json
import os
import re
import shutil
import struct
import subprocess
import tempfile
import zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
PRIVATE_PARTS = {'.git', '.aws', '.codex', '.agents', '.pnpm-store', 'node_modules', 'tmp', 'marketplace-profiles', 'coverage', 'dist', 'playwright-report', 'test-results', '__pycache__', 'backups'}
PATTERNS = {
    'private_key': re.compile(rb'-----BEGIN (?:RSA |EC |OPENSSH |DSA )?PRIVATE KEY-----'),
    'provider_key': re.compile(rb'\bsk-(?:proj-)?[A-Za-z0-9_-]{30,}\b'),
    'github_token': re.compile(rb'\b(?:gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{30,})\b'),
    'aws_access_key': re.compile(rb'\b(?:AKIA|ASIA)[A-Z0-9]{16}\b'),
}

def private_path(name):
    p = Path(name)
    if set(p.parts) & PRIVATE_PARTS: return True
    if name.startswith(('checkoutguard_BE/evidence/', 'checkoutguard_BE/bin/', '.entire/tmp/', '.entire/metadata/', '.entire/logs/', '.entire/redactors/local/')): return True
    if name == '.entire/settings.local.json': return True
    if p.name in ('Cookies', 'Cookies-journal', 'Login Data', 'Login Data-journal') or re.match(r'(?:storageState|auth-state).*\.json$', p.name): return True
    if p.name.endswith('.env.example'): return False
    return p.name == '.env' or p.name.startswith('.env.') or bool(re.search(r'\.(env|local|pem|key|log|tsbuildinfo|pyc|test|prof|exe|db|sqlite|sqlite3|dump)$', p.name)) or '.env.' in p.name

def visible_files(root):
    rg = shutil.which('rg')
    if not rg: raise ValueError('ripgrep is required; Git is never invoked')
    result = subprocess.run([rg, '--files', '--hidden', '--no-require-git', '--no-ignore-parent', '--no-ignore-global', '--no-ignore-dot', '-g', '!.git/**', '-0'], cwd=root, capture_output=True)
    if result.returncode not in (0, 1): raise ValueError('file enumeration failed')
    return {os.fsdecode(p).replace(os.sep, '/') for p in result.stdout.split(b'\0') if p}

def index_paths(data, hash_length=20):
    if len(data) < 12 or data[:4] != b'DIRC': raise ValueError('invalid index header')
    version, count = struct.unpack_from('>II', data, 4)
    if version not in (2, 3, 4): raise ValueError('unsupported index version')
    offset, previous, names = 12, b'', set()
    for _ in range(count):
        start = offset
        flags_at = start + 40 + hash_length
        flags = struct.unpack_from('>H', data, flags_at)[0]
        offset = flags_at + 2 + (2 if flags & 0x4000 else 0)
        if version == 4:
            byte = data[offset]; offset += 1
            trim = byte & 0x7f
            while byte & 0x80:
                byte = data[offset]; offset += 1
                trim = ((trim + 1) << 7) + (byte & 0x7f)
            if trim > len(previous): raise ValueError('invalid index path prefix')
            end = data.index(b'\0', offset)
            name = previous[:len(previous)-trim] + data[offset:end]
            offset = end + 1
        else:
            end = data.index(b'\0', offset)
            name = data[offset:end]
            offset = start + ((end + 1 - start + 7) // 8) * 8
        previous = name
        value = os.fsdecode(name)
        if not value or Path(value).is_absolute() or any(part in ('.', '..', '.git') for part in Path(value).parts): raise ValueError('unsupported index path')
        names.add(value)
    # Split indices require joining another file; do not silently skip entries.
    end_of_extensions = len(data) - hash_length
    while offset + 8 <= end_of_extensions:
        tag, length = data[offset:offset+4], struct.unpack_from('>I', data, offset+4)[0]
        if tag == b'link': raise ValueError('split index is not supported by this audit')
        offset += 8 + length
    return names

def tracked_paths(root):
    directory = root / '.git'
    if directory.is_file():
        pointer = directory.read_text().strip()
        if not pointer.startswith('gitdir: '): raise ValueError('unsupported metadata pointer')
        directory = (root / pointer[8:]).resolve()
    index = directory / 'index'
    if not index.is_file(): return set()
    config = directory / 'config'
    hash_length = 32 if config.is_file() and re.search(r'objectformat\s*=\s*sha256', config.read_text(), re.I) else 20
    return index_paths(index.read_bytes(), hash_length)

def known_secrets(root):
    values = set()
    paths = [root.parent / '.apertus-api.env', root.parent / '.neatlogs-api.env']
    paths.extend(root.glob('.env*'))
    for path in paths:
        if not path.is_file() or path.name.endswith('.example'): continue
        for raw in path.read_text().splitlines():
            line = raw.strip().removeprefix('export ')
            if not line or line.startswith('#') or '=' not in line: continue
            name, value = line.split('=', 1)
            if not re.search(r'(?:API_KEY|TOKEN|SECRET|PRIVATE_KEY)$', name.strip()): continue
            value = value.strip().strip('\"\'')
            if len(value) >= 8: values.add(value.encode())
    for name, value in os.environ.items():
        if re.search(r'(?:API_KEY|TOKEN|SECRET|PRIVATE_KEY)$', name) and len(value) >= 8: values.add(value.encode())
    return values

def ignore_contract(root):
    public = ['README.md', '.env.example', 'docs/examples/model-candidate.env.example', 'docs/examples/publicai-pricing.env.example', '.entire/settings.json', '.entire/.gitignore', 'checkoutguard_FE/pnpm-lock.yaml', 'checkoutguard_BE/go.sum', 'submission/evidence/runs.json', 'submission/evidence/retest.png', '.github/workflows/ci.yml', 'checkoutguard_FE/public/business/CheckoutGuard-business-plan.pdf', 'checkoutguard_FE/public/business/CheckoutGuard-cfoai-model.xlsx', 'submission/CFO-RESULTS.md']
    private = ['.env', '.env.model-candidate.local', 'config/secret.env', 'config/secret.env.backup', '.pnpm-store/v11/index.db', '.codex/hooks.json', '.aws/credentials', '.agents/session.json', '.entire/settings.local.json', '.entire/metadata/session.json', 'checkoutguard_BE/tmp/marketplace-profiles/shopee/Cookies', 'checkoutguard_BE/evidence/run/0.png', 'checkoutguard_BE/bin/server', 'checkoutguard_FE/node_modules/library/index.js', 'checkoutguard_FE/dist/index.html', 'tmp/run.json', 'coverage/results.json', 'secret.key', 'client.pem', 'Cookies', 'Login Data', 'storageState.json', 'auth-state.json', 'backup.dump']
    with tempfile.TemporaryDirectory(prefix='checkoutguard-ignore-') as directory:
        fixture = Path(directory)
        for name in public + private:
            path = fixture / name
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_bytes(b'')
        shutil.copyfile(root / '.gitignore', fixture / '.gitignore')
        visible = visible_files(fixture)
        return [{'file': p, 'kind': 'public_file_ignored'} for p in public if p not in visible] + [{'file': p, 'kind': 'private_file_not_ignored'} for p in private if p in visible]

def secret_kinds(data, secrets):
    kinds = {kind for kind, pattern in PATTERNS.items() if pattern.search(data)}
    if any(value in data for value in secrets): kinds.add('known_secret_value')
    return kinds

def compressed_secret_kinds(path, secrets):
    # XLSX stores text and formulas in compressed XML; raw-byte scanning alone
    # cannot detect credentials pasted into a cell or document property.
    kinds = set()
    try:
        with zipfile.ZipFile(path) as archive:
            entries = archive.infolist()
            if len(entries) > 2000 or sum(item.file_size for item in entries) > 20 << 20:
                return {'compressed_file_too_large_for_audit'}
            for item in entries:
                kinds.update(secret_kinds(archive.read(item), secrets))
    except (OSError, ValueError, RuntimeError, zipfile.BadZipFile):
        return {'compressed_file_could_not_be_audited'}
    return kinds

def audit(root):
    root = root.resolve()
    contract = ignore_contract(root)
    issues = list(contract)
    candidates, tracked = visible_files(root), tracked_paths(root)
    secrets = known_secrets(root)
    checked, manifest = 0, []
    for name in sorted(candidates | tracked):
        path = root / name
        if private_path(name):
            issues.append({'file': name, 'kind': 'private_path_in_public_candidates'})
            continue
        if path.is_symlink() or not path.resolve().is_relative_to(root):
            issues.append({'file': name, 'kind': 'symlink_or_external_path'})
            continue
        if not path.is_file():
            issues.append({'file': name, 'kind': 'tracked_file_not_in_worktree'})
            continue
        if path.stat().st_size > 20 << 20:
            issues.append({'file': name, 'kind': 'file_too_large_for_audit'})
            continue
        data = path.read_bytes()
        checked += 1
        kinds = secret_kinds(data, secrets)
        if path.suffix.lower() in ('.xlsx', '.xlsm', '.zip'):
            kinds.update(compressed_secret_kinds(path, secrets))
        for kind in sorted(kinds): issues.append({'file': name, 'kind': kind})
        manifest.append({'file': name, 'bytes': len(data), 'sha256': hashlib.sha256(data).hexdigest()})
    return {'status': 'pass' if not issues else 'fail', 'files_checked': checked, 'tracked_paths_checked': len(tracked), 'ignore_contract': 'pass' if not contract else 'fail', 'issues': issues, 'git_commands_executed': False, 'staged_and_historical_blobs_checked': False, 'manifest': manifest}

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--output', default='tmp/publication-audit.json')
    parser.add_argument('--archive', help='Optional public source ZIP inside ignored tmp; created only after the audit passes')
    args = parser.parse_args()
    root = ROOT.resolve()
    output = (root / args.output).resolve()
    if not output.is_relative_to(root / 'tmp'): raise ValueError('audit reports must stay inside the ignored tmp directory')
    report = audit(root)
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(json.dumps(report, indent=2) + '\n')
    output.chmod(0o600)
    if args.archive:
        archive = (root / args.archive).resolve()
        if not archive.is_relative_to(root / 'tmp'): raise ValueError('archive must stay inside ignored tmp')
        if report['status'] != 'pass': raise ValueError('refusing to package files after a failed audit')
        archive.parent.mkdir(parents=True, exist_ok=True)
        with zipfile.ZipFile(archive, 'w', compression=zipfile.ZIP_DEFLATED) as package:
            for item in report['manifest']:
                data = (root / item['file']).read_bytes()
                if hashlib.sha256(data).hexdigest() != item['sha256']: raise ValueError('source changed after audit; run again')
                package.writestr(item['file'], data)
            package.writestr('PUBLICATION-MANIFEST.json', json.dumps(report, indent=2) + '\n')
        report['archive'] = str(archive.relative_to(root))
        output.write_text(json.dumps(report, indent=2) + '\n')
    print(json.dumps({k: v for k, v in report.items() if k != 'manifest'}, indent=2))
    raise SystemExit(0 if report['status'] == 'pass' else 1)

if __name__ == '__main__':
    try: main()
    except (OSError, ValueError, struct.error, IndexError):
        print(json.dumps({'status': 'fail', 'error': 'audit could not complete; inspect paths and repository metadata without printing credentials', 'git_commands_executed': False}))
        raise SystemExit(1)
