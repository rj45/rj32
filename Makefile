
.PHONY: fib
fib: mc
	customasm -f logisim16 programs/fib.asm -o dig/test.hex
	tail -n +2 dig/test.hex > hdl/test.hex

.PHONY: sieve
sieve: mc
	customasm -f logisim16 programs/sieve.asm -o dig/test.hex
	tail -n +2 dig/test.hex > hdl/test.hex

.PHONY: hello
hello: mc
	customasm -f logisim16 programs/hello.asm -o dig/test.hex
	tail -n +2 dig/test.hex > hdl/test.hex

.PHONY: mc
mc: dig/microcode.hex

.PHONY: displaymc
displaymc: dig/displaymc.hex

dig/microcode.hex: microcode/microcode.asm
	customasm -f intelhex microcode/microcode.asm -o dig/microcode.hex

dig/displaymc.hex: microcode/display.asm
	customasm -f logisim16 microcode/display.asm -o dig/displaymc.hex


.PHONY: addtest
addtest: mc
	customasm -f logisim16 programs/tests/add.asm -o dig/test.hex
	tail -n +2 dig/test.hex > hdl/test.hex

.PHONY: jumptest
jumptest: mc
	customasm -f logisim16 programs/tests/jump.asm -o dig/test.hex
	tail -n +2 dig/test.hex > hdl/test.hex

.PHONY: calltest
calltest: mc
	customasm -f logisim16 programs/tests/call.asm -o dig/test.hex
	tail -n +2 dig/test.hex > hdl/test.hex

.PHONY: loadstoretest
loadstoretest: mc
	customasm -f logisim16 programs/tests/loadstore.asm -o dig/test.hex
	tail -n +2 dig/test.hex > hdl/test.hex

.PHONY: immtest
immtest: mc
	customasm -f logisim16 programs/tests/imm.asm -o dig/test.hex
	tail -n +2 dig/test.hex > hdl/test.hex


.PHONY: updatesuite
updatesuite:
	ruby scripts/updatetestsuite.rb > dig/testsuite2.dig
	diff -u dig/testsuite.dig dig/testsuite2.dig || :
	mv dig/testsuite2.dig dig/testsuite.dig

.PHONY: testemu
testemu:
	cd emu && go build emu.go && cd ..
	find programs/tests/*.asm | \
		xargs -I testname sh -c 'echo testname && \
		customasm -f logisim16 testname -qp | \
		emu/emu -novdp -run - -trace -maxcycles 100'
	rm emu/emu
	@echo "All passed!"


###########
# Fusesoc #
###########

.PHONY: sim
sim:
	fusesoc run --target sim rj45:rj32:soc

.PHONY: icezero
icezero:
	fusesoc run --target icezero rj45:rj32:soc

.PHONY: upload
upload: icezero
	scp build/rj45_rj32_soc_1.0.0/icezero-icestorm/rj45_rj32_soc_1.0.0.bin pi@raspberrypi.local:~/icezero.bin

.PHONY: clean
clean:
	rm -rf build fusesoc/build
