
.PHONY: fib
fib: mc
	customasm -f logisim16 programs/fib.asm -o dig/test.hex

.PHONY: mc
mc: dig/microcode.hex

dig/microcode.hex: microcode/microcode.asm
	customasm -f logisim16 microcode/microcode.asm -o dig/microcode.hex

.PHONY: addtest
addtest: mc
	customasm -f logisim16 programs/tests/add.asm -o dig/test.hex

.PHONY: jumptest
jumptest: mc
	customasm -f logisim16 programs/tests/jump.asm -o dig/test.hex

.PHONY: calltest
calltest: mc
	customasm -f logisim16 programs/tests/call.asm -o dig/test.hex
