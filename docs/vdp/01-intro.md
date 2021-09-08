# rj32 Video Display Processor

A retro Video Display Processor (VDP) for hobby CPUs.

## Video Display Processor (VDP)

- Displays pixels from memory on a screen
- Designed to make use of small amounts of memory
- Usually no bitmap frame buffer
- 8x8 pixel square (tile/character) based
- Scrollable backgrounds (tile maps)
- Often has a text mode
- Movable object graphics (sprites)
- Also known as:
  - Video Display Unit (VDU)
  - Picture Processing Unit (PPU)
  - Video Interface Chip (VIC)

## Goals and Requirements

- Retro: late 80s / early 90s era technology
- Pixel art graphics
- Designed for demo scene or games
- Designed for a fairly slow CPU
- Low memory requirements
  - No frame buffer (except in simulation)
- Tile based
- Built in an inexpensive FPGA (ice40)
- DVI output over HDMI connector, or VGA
- Ensure 1080p is possible if FPGA is fast enough
  - Optional pixel doubling for pixel art
- One circuit that can do
  - Lots of sprites (movable graphics)
  - Scrollable backgrounds
  - Scrollable text

## Specifications

**NOTE:** This is _PRELIMINARY_ and subject to change!

Please see **github.com/rj45/rj32** for latest specs!

- Resolution:
  - Up to 1080p if FPGA is fast enough
    - Pixel doubling, tripling
  - 640x360 typical resolution (1/2 720p, 1/3 1080p)
  - Palette of 1024 total 24-bit colours
- Sprites:
  - Up to 512 total sprites, but 256 fitted standard
  - Depending on resolution, approx. 800 sprite pixels minimum per line
  - Sprite dimensions from 8x8 to 1024x1024
    - Width and height independently configurable
  - Tile based with 8x8 pixel tiles
  - Tiles are laid out in "sprite sheets"
    - 128x128 grid of tiles
    - Many sprite sheets available (limited by memory)
    - Each sheet can have up to 2048 unique tiles
    - Sprite sheets can overlap to save memory
  - Each sprite can have a configurable tile set
    - Tile sets can share the same sprite sheet
- Colours
  - A pixel in a tile can have one of 16 colours (15 if transparent)
  - A sprite sheet can specify one of 32 palettes for a tile
  - A sprite can specify one of two palette sets for a sprite
  - Total possible colours: 1024 24-bit colours
    - Current hardware can only display 12-bit colours
- Text modes:
  - Sprite sheets can be used as a text buffer
  - Possible to use sprites to configure 8x8 or 8x16 font
  - Text buffer is 128x128 characters in size
    - Supports smooth scrolling in x and y direction
- Most FPGA dev boards have only one RAM chip
  - Must be able to share that RAM chip with the CPU
  - Must support variable latency (SDRAM)

## Data flow

![Data Flow](./dataflow.svg "Data Flow")

- Each sprite
  - Has a rectangle of pixels on a sprite sheet
  - Sprite sheet contains palette and tile IDs
  - Tiles contain 4-bit pixels
  - Pixels are drawn on a line buffer with palette ID
  - Line buffer pixels are looked up in the palette
  - Palette colours are drawn to the screen

## Plan

- Instead of build videos
  - Progress report videos
  - Focus on telling the story of each module
  - Explain deeper how each module works
  - Unlocks me to make faster progress
  - Let me know what you think in the comments!
