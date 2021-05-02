.PHONY: all
all: upload

test.json: vga_display.v vgatop.v bram.v
	yosys -p 'synth_ice40 -dsp -top vgatop -json test.json' vga_display.v vgatop.v bram.v

test.asc: test.json icesugar.pcf
	nextpnr-ice40 --up5k --freq 22 --package sg48 --json test.json --pcf icesugar.pcf --asc test.asc --opt-timing --placer heap --pcf-allow-unconstrained

testout.bin: test.asc
	icepack test.asc testout.bin

.PHONY: lint
lint:
	verilator --lint-only -Wall vga_display.v vgatop.v

.PHONY: upload
upload: testout.bin
	icesprog testout.bin
