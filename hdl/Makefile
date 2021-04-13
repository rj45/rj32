.PHONY: all
all: upload

test.json: veritest.v digits.v mem.v top.v debouncer.v test.hex
	yosys -p 'synth_ice40 -dsp -top top -json test.json' veritest.v digits.v mem.v top.v debouncer.v

test.asc: test.json icesugar.pcf
	nextpnr-ice40 --up5k --package sg48 --json test.json --pcf icesugar.pcf --asc test.asc --opt-timing --placer heap --pcf-allow-unconstrained

testout.bin: test.asc
	icepack test.asc testout.bin

.PHONY: upload
upload: testout.bin
	icesprog testout.bin
