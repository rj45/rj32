
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


.PHONY: clean
clean:
	rm -rf build fusesoc/build
