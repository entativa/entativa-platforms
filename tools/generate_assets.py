#!/usr/bin/env python3
"""
Generate Android PNG drawables and iOS PDF assets from SVG sources using cairosvg.
Supports --dry-run to preview.
"""
import argparse
import os
from pathlib import Path

try:
    from cairosvg import svg2png, svg2pdf
except Exception:
    svg2png = None
    svg2pdf = None

DENSITIES = {
    'mdpi': 1.0,
    'hdpi': 1.5,
    'xhdpi': 2.0,
    'xxhdpi': 3.0,
    'xxxhdpi': 4.0,
}
BASE_PX = 48


def ensure_dir(p: Path):
    p.mkdir(parents=True, exist_ok=True)


def gen_android(src: Path, res_base: Path, dry_run: bool):
    for density, scale in DENSITIES.items():
        out_dir = res_base / f'drawable-{density}'
        if dry_run:
            print('DRY: make dir', out_dir)
        else:
            ensure_dir(out_dir)
        size = int(round(BASE_PX * scale))
        out = out_dir / (src.stem + '.png')
        if dry_run:
            print(f'DRY: svg2png {src} -> {out} @ {size}x{size}')
        else:
            if svg2png is None:
                raise RuntimeError('cairosvg not installed; pip install cairosvg')
            svg2png(url=str(src), write_to=str(out), output_width=size, output_height=size)
            print('Wrote', out)


def gen_ios(src: Path, xcassets_dir: Path, dry_run: bool):
    imageset = xcassets_dir / (src.stem + '.imageset')
    if dry_run:
        print('DRY: make dir', imageset)
    else:
        ensure_dir(imageset)
    pdf_path = imageset / (src.stem + '.pdf')
    contents = imageset / 'Contents.json'
    if dry_run:
        print(f'DRY: svg2pdf {src} -> {pdf_path}')
        print(f'DRY: write {contents}')
    else:
        if svg2pdf is None:
            raise RuntimeError('cairosvg not installed; pip install cairosvg')
        svg2pdf(url=str(src), write_to=str(pdf_path))
        import json
        contents_json = {
            'images': [{ 'idiom': 'universal', 'filename': f'{src.stem}.pdf', 'scale': '1x' }],
            'info': { 'version': 1, 'author': 'xcode' }
        }
        with open(contents, 'w') as f:
            json.dump(contents_json, f, indent=2)
        print('Wrote', pdf_path)


def main():
    p = argparse.ArgumentParser()
    p.add_argument('--sources', nargs='+', default=['assets/canonical-svgs'], help='source folders with SVGs')
    p.add_argument('--android-projects', nargs='*', default=[
        'SocialinkAndroid',
        'VignetteAndroid',
    ], help='android project roots')
    p.add_argument('--ios-projects', nargs='*', default=[
        'SocialinkiOS',
        'VignetteiOS',
    ], help='ios project roots')
    p.add_argument('--dry-run', action='store_true')
    args = p.parse_args()

    svgs = []
    for s in args.sources:
        pth = Path(s)
        if not pth.exists():
            continue
        svgs.extend(sorted(pth.rglob('*.svg')))

    if not svgs:
        print('No SVG sources found in', args.sources)
        return

    # Android: generate pngs into app/src/main/res/drawable-*
    for proj in args.android_projects:
        res_base = Path(proj) / 'app' / 'src' / 'main' / 'res'
        for svg in svgs:
            try:
                gen_android(svg, res_base, args.dry_run)
            except Exception as e:
                print('Android generation failed for', svg, '->', e)

    # iOS: generate xcassets
    for proj in args.ios_projects:
        xcassets = Path(proj) / 'Assets' / 'Icons.xcassets'
        for svg in svgs:
            try:
                gen_ios(svg, xcassets, args.dry_run)
            except Exception as e:
                print('iOS generation failed for', svg, '->', e)

    print('Done')

if __name__ == '__main__':
    main()
