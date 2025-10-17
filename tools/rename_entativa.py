#!/usr/bin/env python3
"""
Rename occurrences of 'socialink' -> 'entativa' across filesystem names and text file contents.
- Preserves case for simple variations (Socialink, socialink, SOCIALINK)
- Supports --dry-run (report planned actions) and --apply (perform them)
- Skips binary files (images, compiled files)
- Avoids descending into .git and vendor directories

Use with caution and review the --dry-run output first.
"""
import argparse
import os
import re
from pathlib import Path

REPLACEMENTS = [
    (re.compile(r'socialink', re.IGNORECASE), 'entativa'),
]

SKIP_DIRS = {'.git', 'node_modules', 'venv', '__pycache__', 'build'}
BINARY_EXTS = {'.png', '.jpg', '.jpeg', '.gif', '.ico', '.pdf', '.zip', '.tar', '.gz', '.jar', '.class'}


def preserve_case_replace(match, new):
    src = match.group(0)
    if src.isupper():
        return new.upper()
    if src[0].isupper():
        return new.capitalize()
    return new


def plan_and_apply(root: Path, dry_run: bool):
    file_content_changes = []
    file_rename_changes = []

    for p in root.rglob('*'):
        # skip top-level unwanted dirs
        if any(part in SKIP_DIRS for part in p.parts):
            continue
        # plan renames first (files and directories)
        name = p.name
        new_name = name
        for pattern, replacement in REPLACEMENTS:
            new_name = pattern.sub(lambda m: preserve_case_replace(m, replacement), new_name)
        if new_name != name:
            file_rename_changes.append((p, p.with_name(new_name)))

        # plan content replacements for text files
        if p.is_file():
            if p.suffix.lower() in BINARY_EXTS:
                continue
            try:
                data = p.read_text(encoding='utf-8')
            except Exception:
                # likely binary or unreadable as text
                continue
            new_data = data
            for pattern, replacement in REPLACEMENTS:
                new_data = pattern.sub(lambda m: preserve_case_replace(m, replacement), new_data)
            if new_data != data:
                file_content_changes.append((p, len(data), len(new_data)))

    # Report actions
    print(f'Planned {len(file_rename_changes)} renames and {len(file_content_changes)} content edits under {root}')
    for src, dst in file_rename_changes:
        print('RENAME:', src, '->', dst)
    for f, old_len, new_len in file_content_changes:
        print('EDIT CONTENT:', f, f'({old_len} -> {new_len} bytes)')

    if dry_run:
        return 0

    # Apply renames: directories should be renamed after children -> sort by depth descending
    for src, dst in sorted(file_rename_changes, key=lambda x: -len(x[0].parts)):
        try:
            src.rename(dst)
            print('Renamed', src, '->', dst)
        except Exception as e:
            print('Failed rename', src, '->', dst, e)

    # Apply content changes
    for p, old_len, new_len in file_content_changes:
        try:
            data = p.read_text(encoding='utf-8')
            for pattern, replacement in REPLACEMENTS:
                data = pattern.sub(lambda m: preserve_case_replace(m, replacement), data)
            p.write_text(data, encoding='utf-8')
            print('Edited', p)
        except Exception as e:
            print('Failed edit', p, e)

    return 0


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--root', default='.', help='project root to apply rename')
    group = parser.add_mutually_exclusive_group()
    group.add_argument('--dry-run', action='store_true')
    group.add_argument('--apply', action='store_true')
    args = parser.parse_args()

    root = Path(args.root)
    return plan_and_apply(root, args.dry_run)

if __name__ == '__main__':
    raise SystemExit(main())
