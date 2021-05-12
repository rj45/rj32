.PHONY: all
all: upload

test.json: rj32.v digits.v mem.v top.v debouncer.v test.hex vga_blinkenlights.v
	yosys -p 'read_verilog -nooverwrite rj32.v digits.v mem.v top.v debouncer.v vga_blinkenlights.v; synth_ice40 -dsp -top top -json test.json'

test.asc: test.json icesugar.pcf
	nextpnr-ice40 --up5k --package sg48 --freq 12.5 --json test.json --pcf icesugar.pcf --asc test.asc --opt-timing --placer heap --pcf-allow-unconstrained

testout.bin: test.asc
	icepack test.asc testout.bin

.PHONY: lint
lint:
	verilator --lint-only -Wall rj32.v digits.v mem.v top.v debouncer.v sb_spram256ka.v

.PHONY: upload
upload: testout.bin
	icesprog testout.bin
