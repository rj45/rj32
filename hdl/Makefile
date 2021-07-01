.PHONY: all
all: lintsugar lintzero icesugar.bin icezero.bin

.PHONY: lintsugar
lintsugar:
	verilator --lint-only -Wall -Wno-PINMISSING -Wno-DECLFILENAME -Wno-MODDUP -Wno-UNUSED -Wno-WIDTH --top-module icesugar_top top.v rj32.v digits.v mem.v top.v icesugar_top.v sb_ram40_4k.v sb_spram256ka.v sb_pll40_2f_pad.v

.PHONY: lintzero
lintzero:
	verilator --lint-only -Wall -Wno-PINMISSING -Wno-DECLFILENAME -Wno-MODDUP -Wno-UNUSED -Wno-WIDTH --top-module icezero_top top.v rj32.v digits.v mem.v top.v icezero_top.v sb_ram40_4k.v sb_pll40_2f_pad.v


.PHONY: clean
clean:
	rm -rf *.rpt *.asc *.json *.bin *.log

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

# memory_dff  vgafrontpanel.vga_debugrender

icezero.json: rj32.v digits.v mem.v icezero_top.v top.v test.hex vga_blinkenlights.v
	yosys -ql icezero_yosys.log -p 'read_verilog -nooverwrite icezero_top.v top.v rj32.v vga_blinkenlights.v digits.v mem.v; synth_ice40 -top icezero_top; write_json icezero.json'

icezero.asc: icezero.json icezero.pcf
	nextpnr-ice40  --hx8k --package tq144:4k --freq 12.5 --json icezero.json --pcf icezero.pcf --asc icezero.asc --opt-timing --placer heap --pcf-allow-unconstrained

icezero.bin: icezero.asc
	icepack icezero.asc icezero.bin
	# icetime -d hx8k -c 12.5 -p icezero.pcf -t icezero.asc

.PHONY: icezupload
icezupload: icezero.bin
	scp icezero.bin pi@raspberrypi.local:~/


