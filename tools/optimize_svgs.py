#!/usr/bin/env python3
"""
Optimize SVG files using scour (if available). The script writes to a temp file then replaces the original.
Supports --dry-run to preview actions.
"""
import argparse
import subprocess
import tempfile
import shutil
from pathlib import Path

def optimize_file(path: Path, dry_run: bool):
    out = Path(tempfile.mkstemp(suffix='.svg')[1])
    cmd = ['scour', '--enable-viewboxing', '--remove-metadata', '--shorten-ids', str(path), str(out)]
    if dry_run:
        print('DRY:', ' '.join(cmd))
        return True
    try:
        subprocess.run(cmd, check=True)
        shutil.move(str(out), str(path))
        print('Optimized', path)
        return True
    except FileNotFoundError:
        print('scour not found; install with `pip install scour` or use your OS package manager')
        return False
    except subprocess.CalledProcessError as e:
        print('scour failed for', path, '->', e)
        try:
            out.unlink()
        except Exception:
            pass
        return False


def main():
    p = argparse.ArgumentParser()
    p.add_argument('targets', nargs='+', help='directories or files to optimize')
    p.add_argument('--dry-run', action='store_true')
    args = p.parse_args()

    paths = []
    for t in args.targets:
        pth = Path(t)
        if pth.is_dir():
            paths.extend(sorted(pth.rglob('*.svg')))
        elif pth.is_file():
            paths.append(pth)
        else:
            print('Target not found:', t)

    failed = 0
    for f in paths:
        ok = optimize_file(f, args.dry_run)
        if not ok:
            failed += 1
    if failed:
        print('Completed with', failed, 'failures')
    else:
        print('All optimized')

if __name__ == '__main__':
    main()
