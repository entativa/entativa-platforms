ICON asset pipeline

This document explains how icon assets in this repository are organized and how to regenerate platform-specific assets from the canonical SVG sources.

Summary
- Canonical SVGs: stored in `assets/canonical-svgs/` (the authoritative source used to regenerate client and platform assets).
- Per-client copies: each client project keeps its own `assets/icons/` folder for working SVG sources. These are normalized to kebab-case filenames. A `filename-map.txt` in each folder records original → normalized name mappings.
- Platform artifacts:
  - Android: PNG drawables at multiple densities are placed in `app/src/main/res/drawable-<density>/`.
  - iOS: PDF vector assets are placed in `Assets/Icons.xcassets/<icon>.imageset/<icon>.pdf` along with a `Contents.json`.

Helper scripts (tools/)
We added small helper scripts under `tools/` to make the pipeline reproducible:

- `tools/normalize_icons.py`
  - Copies canonical SVGs from `assets/canonical-svgs/` into each client `assets/icons/` directory,
    normalizes filenames to kebab-case, and writes `filename-map.txt` per client.
  - Usage (dry-run):

```bash
python3 tools/normalize_icons.py --dry-run
```

- `tools/optimize_svgs.py`
  - Runs `scour` on target directories or files, writing to a temporary file then replacing the original.
  - Usage (dry-run):

```bash
python3 tools/optimize_svgs.py --dry-run assets/canonical-svgs
```

- `tools/generate_assets.py`
  - Generates Android PNG drawables and iOS PDF assets from SVG sources using CairoSVG.
  - Usage (dry-run):

```bash
python3 tools/generate_assets.py --dry-run
```

Recommended workflow
1. Update canonical SVGs in `assets/canonical-svgs/`.
2. Normalize and copy to clients:

```bash
python3 tools/normalize_icons.py
```

3. Optimize the canonical and client SVGs (optional):

```bash
python3 tools/optimize_svgs.py assets/canonical-svgs SocialinkAndroid/app/src/main/assets/icons
```

4. Generate platform assets (PNGs + PDFs):

```bash
python3 tools/generate_assets.py
```

Notes
- Keep canonical SVGs under `assets/canonical-svgs/` to avoid drift. Generated platform assets should be committed for CI and release, but treat `assets/canonical-svgs/` as the single source of truth.
- The scripts support `--dry-run` so you can preview changes before they write files.
- If you want Android VectorDrawable conversion instead of PNG rasterization, that can be added as an optional step.
