# HDL

This is the HDL version of the processor. In Digital, key circuits are exported to verilog and saved in this folder. Then a verilog top level file provides the peripherals that Digital doesn't support like block ram and hdmi out.

## Boards

- icesugar from MuseLabs
  - up5k and 3 pmods
- icezero from trendz
  - hx4k, 512 KB 10 ns SRAM, 4 pmods, raspberry pi zero hat
- BlackIce MX from tindie / folknology
  - hx4k, 512 KB SDRAM, 6 pmods, arm usb programmer
  - not fully implemented, sdram controller needed

Each board has a .pcf file and a top level adaptor verilog.

## PMODs

- icebreaker DVI/HDMI 12bit PMOD
  - pmod2 + pmod3 on icesugar
  - P3 + P4 on icezero
- digilent 410-077 debounced button PMOD
  - pmod 1 bottom half on icesugar
  - P1 bottom half on icezero

## Building

Run either `make icesugar` or `make icezero` to build and program either board.

Both boards need an up to date install of yosys and nextpnr.

The icesugar board needs [icesprog](https://github.com/wuxx/icesugar/tree/master/tools/src).

The icezero board needs a raspberry pi set up with ssh, and you can use [icezprog](https://github.com/cliffordwolf/icotools/tree/master/examples/icezero) on the pi. I just have a script waiting in a loop for the file to upload, and it will program it and restart it.

Make sure to remove the `-Os` from the makefile for icezprog, since optimizing the binary messes up the timing causing errors.

## Running programs

The program is stored in `test.hex`. Build for Digital, then just remove the header for the hdl version.
