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
  - If adjacent tiles on a grid share a palette (for example if a 2x2 grid of tiles (16x16px) all share the same palette number) then:
    - `-gridw` is the width of the grid in tiles (defaults to 1)
    - `-gridh` is the height of the grid in tiles (defaults to 1)
  - Optionally the image can be reclustered after an initial global dither pass to try to incorporate dither pixels into the palettes of tiles in a better way with `-recluster`
    - This does not always produce a better image
- Can produce a low tile count test image with `-gentest`
- Can produce image for the clustering with `-clusterfile`
- Can produce image for the final result with `-outfile`
- Can produce image for the tile set with `-tilesimg`
- Can emit JSON of the result, for easy creation of custom binary packing scripts with `-json`
  - Format is:
    - `palettes [][]{R:int,G:int,B:int}`: list of palettes
    - `tilepals []int`: the palette id for each tile id
    - `tiles    [][]int`: the list of colors in each tile
    - `tilemap  [][]int`: the list of tiles in each row of the full image broken down into 8x8px tiles
- Can emit hex files in the specific format rj32 requires (`-pal`, `-tiles`, and `-map`)

## License

Copyright (C) 2021 rj45 and contributors

This project is licensed under the MIT License - see the LICENSE file for details.
