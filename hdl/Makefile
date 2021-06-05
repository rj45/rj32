.PHONY: all
all: lint icesugar.bin

.PHONY: icesugar
icesugar: icesupload

icesugar.json: rj32.v digits.v mem.v icesugar_top.v top.v test.hex vga_blinkenlights.v
	yosys -p 'read_verilog -nooverwrite icesugar_top.v top.v rj32.v vga_blinkenlights.v digits.v mem.v; synth_ice40 -dsp -top icesugar_top -json icesugar.json'

icesugar.asc: icesugar.json icesugar.pcf
	nextpnr-ice40 --up5k --package sg48 --freq 12.5 --json icesugar.json --pcf icesugar.pcf --asc icesugar.asc --opt-timing --placer heap --pcf-allow-unconstrained

icesugar.bin: icesugar.asc
	icepack icesugar.asc icesugar.bin

.PHONY: lint
lint:
	verilator --lint-only -Wall -Wno-PINMISSING -Wno-DECLFILENAME -Wno-MODDUP --top-module top rj32.v digits.v mem.v top.v sb_spram256ka.v

.PHONY: icesupload
icesupload: icesugar.bin
	icesprog icesugar.bin

.PHONY: clean
clean:
	rm -rf *.rpt *.asc *.json *.bin
