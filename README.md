# rj32

A 16-bit RISC CPU with 32 instructions built with [Digital](https://github.com/hneemann/Digital).

## Description

This is a CPU built from scratch in a visual way using a digital circuit simulator called [Digital](https://github.com/hneemann/Digital). This is then able to be exported to verilog or VHDL and that can be converted into a digital circuit that can run on an [FPGA](https://en.wikipedia.org/wiki/Field-programmable_gate_array).

I have built a couple CPUs before, but this time I decided to record it as a youtube series. I start out being a bit of a youtube n00b, but the quality hopefully gets better as the series goes on.

[![Introduction Video - Building a CPU From Scratch](https://img.youtube.com/vi/FSVhlqE7EgA/0.jpg)](https://www.youtube.com/watch?v=FSVhlqE7EgA&list=PLilenfQGj6CEG6iZ4TQJ10PI7pCWsy1AO&index=1)

## Building and Running

For the simulation, open `dig/frontpanel.dig` in [Digital](https://github.com/hneemann/Digital).

For the verilog version, see [the HDL documentation](./hdl/README.md).

## Design

A [minimal instruction set computer](https://en.wikipedia.org/wiki/Minimal_instruction_set_computer) with exactly 32 instructions. Well, 32 opcodes anyway, technically there's more instructions. Though that may change with the next instruction set update.

Currently it is a RISC instruction set with a Harvard memory architecture. In other words, data memory is accessed only through load/store instructions, and data memory is separate from program memory. Program memory can only be used to execute code.

There is currently two pipeline stages: fetch and execute. In the first stage the instruction is fetched from program memory. In the second stage the instruction is decoded, has its operands loaded from the register file, executed, then the result written.

Currently the CPU is designed to run on ice40 FPGAs using the open source toolchain.

### Instruction Set

[Instruction Set Documentation](./docs/instructions.md).

## License

This project is licensed under the MIT License - see the LICENSE file for details.

NOTE: The youtube videos do not fall under this license. They are under the standard youtube copyright, and I ask you not re-publish them elsewhere. You can create your own youtube videos with the content in this repo, however.

## Acknowledgments

* [Ben Eater](https://eater.net/)
* Many awesome folks on [AnyCPU](http://anycpu.org/forum/)
* John Lluch's [CPU74](https://github.com/John-Lluch/CPU74/)
* [Dieter "ttlworks"](http://www.6502.org/users/dieter/)
* So many inspirations I can't hope to enumerate them all
