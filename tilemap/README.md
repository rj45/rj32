# Tilemap

A program for producing tile-based graphics for homebrew VGA / Video Display Processors.

## Features

- Many dithering matrices supported (`-dither`):
  - `none` - turn off dithering
  - `floydsteinberg` (default)
  - `jarvisjudiceninke`
  - `stucki`
  - `atkinson`
  - `burkes`
  - `sierra`
  - `tworowsierra`
  - `sierralite`
- Color space reduction (`-colorconv`):
  - `24`: 8:8:8 bits RGB (no conversion)
  - `12`: 4:4:4 bits RGB (default)
  - `8`:  3:3:2 bits RGB
- Color reduction by k-means clustering:
  - Finds `-palettes` number of palettes of `-perpalette` colors each
  - Clusters the tiles into one of these palettes
  - Can ensure that the `-transparent` color is always the 0th entry in each palette that has it
  - If adjacent tiles on a grid share a palette (for example if a 2x2 grid of tiles (16x16px) all share the same palette number) then:
    - `-gridw` is the width of the grid in tiles (defaults to 1)
    - `-gridh` is the height of the grid in tiles (defaults to 1)
  - Optionally the image can be reclustered after an initial global dither pass to try to incorporate dither pixels into the palettes of tiles in a better way with `-recluster`
    - This does not always produce a better image
- Can produce a low tile count test image with `-gentest`
  - `-mapw` and `-maph` adjust the size in tiles
- Can produce image for the clustering with `-clusterfile`
- Can produce image for the final result with `-outfile`
- Can produce image for the tile set with `-tilesimg`
- Can emit JSON of the result, for easy creation of custom binary packing scripts with `-json`
  - Format is:
    - `palettes [][]{R:int,G:int,B:int}`: list of palettes
    - `tilepals [][]int`: the palette id for each tilemap tile
    - `tiles    [][]int`: the list of colors in each tile
    - `tilemap  [][]int`: the list of tiles in each row of the full image broken down into 8x8px tiles
- Can emit logisim (also Digitial) hex files in the specific format rj32 requires
  - `-pal` is in 16-bit 4:4:4:4 ARGB format
  - `-tiles` are 4 bpp, packed into 16-bit words
  - `-map` is a list of width x height 16 bit words:
    - lower 12 bits are the tile ID (max 4096 tiles)
    - upper 4 bits are the palette ID (max 16 palettes)
- Can split a large tilemap into multiple sheets (or blocks)
  - Can also be used to pad a tilemap to specific dimensions
  - `-splitmap` turns this on
  - `-mapw` is the width in tiles of each sheet
    - if the tile map is wider it is truncated
    - if the tile map is narrower, it is padded with zero tiles
  - `-maph` is the height in tiles of each sheet
    - tile map is split into N sheets of `-maph` height except the last sheet will have just the remaining height
  - a file is produced for each of `-json`, `-tiles`, `-map` and `-tilesimg`
    - will have a suffix `_0`, `_1`, etc for each sheet

## Changelog

### Aug 13, 2021

- Tiles will now be reused with different palettes
  - The JSON format has been adjusted to specify the palette for each location instead of each tile id
  - Reduces tile count slightly, but there still isn't any smart tile reduction going on, only duplicate elimination
- `-gentest` now takes width and height in tiles from `-mapw` and `-maph`
- Now has support for splitting a tilemap into multiple sheets with `-splitmap`
- Many internal refactorings to simplify some things and remove some dead code
- The hex files for rj32 have changed to be simpler and were documented above

### Aug 20, 2021

- Huge quality improvement:
  - Grids/tiles palette assignment now goes through two searches for the palette with the least error
  - After each search a new palette is quantized
  - Pretty much eliminates weird highly noticable glitches
  - There are still some tiles with poorly picked colors
- Changed default number of palettes to 32
- Hex file output now uses 5 bits per tile for palette id instead of 4 bits
- Added `-transparent` to specify the hex value of a color
  - ensures if any palette contains this color, it's always color zero
- Some code cleanup (removed `intset` that was mostly unused, renamed lab to luv)
- Hex palette generator will now emit 8:8:8 RGB, 4:4:4 and 3:3:2 depending on `colorconv`

## License

Copyright (C) 2021 rj45 and contributors

This project is licensed under the MIT License - see the LICENSE file for details.
