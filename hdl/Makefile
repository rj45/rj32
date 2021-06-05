.PHONY: all
all: lint icesugar.bin

.PHONY: lint
lint:
	verilator --lint-only -Wall -Wno-PINMISSING -Wno-DECLFILENAME -Wno-MODDUP --top-module top rj32.v digits.v mem.v top.v sb_spram256ka.v

.PHONY: clean
clean:
	rm -rf *.rpt *.asc *.json *.bin

###############
## IceSugar
###############

.PHONY: icesugar
icesugar: icesupload

icesugar.json: rj32.v digits.v mem.v icesugar_top.v top.v test.hex vga_blinkenlights.v
	yosys -p 'read_verilog -nooverwrite icesugar_top.v top.v rj32.v vga_blinkenlights.v digits.v mem.v; synth_ice40 -dsp -top icesugar_top -json icesugar.json'

icesugar.asc: icesugar.json icesugar.pcf
	nextpnr-ice40 --up5k --package sg48 --freq 12.5 --json icesugar.json --pcf icesugar.pcf --asc icesugar.asc --opt-timing --placer heap --pcf-allow-unconstrained

icesugar.bin: icesugar.asc
	icepack icesugar.asc icesugar.bin

.PHONY: icesupload
icesupload: icesugar.bin
	icesprog icesugar.bin

###############
## IceZero
###############

.PHONY: icezero
icezero: icezupload

icezero.json: rj32.v digits.v mem.v icezero_top.v top.v test.hex vga_blinkenlights.v
	yosys -p 'read_verilog -nooverwrite icezero_top.v top.v rj32.v vga_blinkenlights.v digits.v mem.v; synth_ice40 -top icezero_top -json icezero.json'

icezero.asc: icezero.json icezero.pcf
	nextpnr-ice40  --hx8k --package tq144:4k --freq 12.5 --json icezero.json --pcf icezero.pcf --asc icezero.asc --opt-timing --placer heap --pcf-allow-unconstrained

icezero.bin: icezero.asc
	icepack icezero.asc icezero.bin
	icetime -d hx8k -c 12.5 -p icezero.pcf -t icezero.asc

.PHONY: icezupload
icezupload: icezero.bin
	scp icezero.bin pi@raspberrypi.local:~/


