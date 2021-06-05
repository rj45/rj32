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
  - P4 + P3 on icezero
- digilent 410-077 debounced button PMOD
  - pmod 1 on icesugar
  - P1 on icezero
