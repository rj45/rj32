.PHONY: all
all: upload

test.json: vga_display.v vgatop.v bram.v
	yosys -p 'synth_ice40 -dsp -abc2 -top vgatop -json test.json' -DRES_720x480 vga_display.v vgatop.v bram.v

test.asc: test.json icesugar.pcf
	nextpnr-ice40 --up5k --freq 22 --package sg48 --json test.json --pcf icesugar.pcf --asc test.asc --opt-timing --placer heap --pcf-allow-unconstrained

hx8k.asc: test.json blackicemx.pcf
	nextpnr-ice40 --hx8k --freq 22 --package tq144:4k --json test.json --pcf blackicemx.pcf --asc test.asc -r --opt-timing --placer sa --pcf-allow-unconstrained

testout.bin: test.asc
	icepack test.asc testout.bin

.PHONY: lint
lint:
	verilator --lint-only -Wall vga_display.v vgatop.v

.PHONY: upload
upload: testout.bin
	icesprog testout.bin
