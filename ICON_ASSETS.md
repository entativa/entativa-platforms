ICON asset pipeline

This document explains how icon assets in this repository are organized and how to regenerate platform-specific assets from the canonical SVG sources.

Summary
- Canonical SVGs: originally stored at the repository root (now removed). The canonical sources should be kept in a directory such as `assets/canonical-svgs/` if you plan to update and regenerate assets.
- Per-client copies: each client project keeps its own `assets/icons/` folder for working SVG sources. These were generated from the canonical SVGs and normalized to kebab-case filenames. A `filename-map.txt` in each folder records original → normalized name mappings.
- Platform artifacts:
  - Android: PNG drawables at multiple densities are placed in `app/src/main/res/drawable-<density>/`.
  - iOS: PDF vector assets are placed in `Assets/Icons.xcassets/<icon>.imageset/<icon>.pdf` along with a `Contents.json`.

Regenerating assets (recommended workflow)
1. Update canonical SVGs in `assets/canonical-svgs/`.
2. Run the normalization script (tools/normalize_icons.py) to copy/normalize into each client `assets/icons/` and write `filename-map.txt` files.
3. Optimize the SVGs with scour: `scour input.svg output.svg --enable-viewboxing --remove-metadata` (or use the repository helper script).
4. Generate platform assets with the helper script (tools/generate_assets.py) which uses CairoSVG to produce Android PNGs and iOS PDF assets.

Notes
- Keep canonical SVGs under a single directory to avoid drift.
- If you prefer Android VectorDrawable XML, add an extra conversion step; currently we generate PNGs as fallbacks.
- The `filename-map.txt` files are authoritative for mapping names.

If you want, I can move the canonical SVGs into `assets/canonical-svgs/` before deleting the root copies.
