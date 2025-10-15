#!/usr/bin/env python3
"""
Copy canonical SVGs into client icon folders and normalize filenames to kebab-case.
Writes a filename-map.txt in each client folder mapping original -> normalized.
Supports --dry-run for safe preview.
"""
import argparse
import os
import re
import shutil
from pathlib import Path

RE_FILENAME = re.compile(r"[^a-z0-9]+")

def normalize(name: str) -> str:
    name = name.lower()
    name = RE_FILENAME.sub('-', name)
    name = re.sub(r'-{2,}', '-', name)
    name = name.strip('-')
    if not name:
        name = 'icon'
    return name


def main():
    p = argparse.ArgumentParser()
    p.add_argument('--canonical', default='assets/canonical-svgs', help='folder with canonical SVGs')
    p.add_argument('--clients', nargs='+', default=[
        'SocialinkAndroid/app/src/main/assets/icons',
        'SocialinkiOS/Assets/Icons',
        'VignetteAndroid/app/src/main/assets/icons',
        'VignetteiOS/Assets/Icons',
    ], help='client icon folders to populate')
    p.add_argument('--dry-run', action='store_true')
    args = p.parse_args()

    canon = Path(args.canonical)
    if not canon.exists():
        print('Canonical folder not found:', canon)
        return

    svgs = [p for p in sorted(canon.iterdir()) if p.suffix.lower() == '.svg']
    if not svgs:
        print('No SVGs found in', canon)
        return

    for client in args.clients:
        client_dir = Path(client)
        if args.dry_run:
            print('\n[DRY RUN] Populate client:', client_dir)
        else:
            client_dir.mkdir(parents=True, exist_ok=True)

        mapping = []
        seen = set()
        for src in svgs:
            base = src.stem
            norm = normalize(base)
            # ensure uniqueness
            dest_name = norm + '.svg'
            i = 1
            while dest_name in seen or (not args.dry_run and (client_dir / dest_name).exists() and (client_dir / dest_name).read_bytes() != src.read_bytes()):
                dest_name = f"{norm}-{i}.svg"
                i += 1
            seen.add(dest_name)
            dest = client_dir / dest_name
            mapping.append((src.name, dest_name))
            if args.dry_run:
                print(f'  {src.name} -> {dest}')
            else:
                shutil.copy2(src, dest)

        # write filename-map.txt
        map_path = client_dir / 'filename-map.txt'
        if args.dry_run:
            print(f'  Would write mapping to {map_path} ({len(mapping)} entries)')
        else:
            with open(map_path, 'w') as f:
                for a, b in mapping:
                    f.write(f"{a} -> {b}\n")
            print('Wrote', map_path)

if __name__ == '__main__':
    main()
