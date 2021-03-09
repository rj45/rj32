
.PHONY: fib
fib:
	customasm -f logisim16 programs/fib.asm -o dig/test.hex

.PHONY: addtest
addtest:
	customasm -f logisim16 programs/tests/add.asm -o dig/test.hex
