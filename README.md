# rj32

A 16-bit RISC CPU with 32 instructions built with [Digital](https://github.com/hneemann/Digital).

## Description

This is a CPU built from scratch in a visual way using a digital circuit simulator called [Digital](https://github.com/hneemann/Digital). This is then able to be exported to verilog or VHDL and that can be converted into a digital circuit that can run on an [FPGA](https://en.wikipedia.org/wiki/Field-programmable_gate_array).

I have built a couple CPUs before, but this time I decided to record it as a youtube series. I start out being a bit of a youtube n00b, but the quality hopefully gets better as the series goes on.

[![Introduction Video - Building a CPU From Scratch](https://img.youtube.com/vi/FSVhlqE7EgA/maxresdefault.jpg | width=640)](https://www.youtube.com/watch?v=FSVhlqE7EgA&list=PLilenfQGj6CEG6iZ4TQJ10PI7pCWsy1AO&index=1)

### Graphics Hardware

There is also a video display circuit designed to work with DVI over HDMI with a PMOD.

The circuit is a retro Video Display Processor (VDP), sometimes also known as a Video Display Unit (VDU), or Picture Processing Unit (PPU). It's designed to work similarily to a late 80s graphics system like in the Commodore 64 / Amiga, or the various Nintendo or Sega consoles of that era.

It uses a 8x8 tile based system, with a text display of 16x32 characters, of which are made of tiles with 3 colors plus transparent. The resolution is 640x400 with a 12bpp colour depth.

Currently this is hard-coded to display the frontpanel blinkenlights, but the text is now in a framebuffer so that could allow CPU access to it soon.

Here is a playlist just of the videos showing how this part was built:

[![Building a GPU From Scratch](https://img.youtube.com/vi/nVaOJ6CwIic/maxresdefault.jpg | width=640)](https://www.youtube.com/watch?v=nVaOJ6CwIic&list=PLilenfQGj6CEbC7-IoXsmrmDfBjiUi6a1&index=1)

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

## Contributing

**Contributions are welcome!**

- Please follow the existing style and implement as much as possible in the Digital simulation rather than verilog.
- Fork and submit a PR, please update any documentation and tests and explain exactly what you changed (preferably with screenshots) in the PR description.
- I will be showing all contributions in youtube videos showing what you changed. I will give you credit.
- I reserve the right to reject any changes that take the processor in a different direction than I want it to go
- Changes I would like to make are listed in the issues. If I haven't begun working on something, feel free to take it up.
- If no issues looks interesting, feel free to submit an issue for something you'd like to do, and I can approve it.
- Feel free to submit PRs without an issue or approval if you don't mind me deciding not to accept it.

## Contributors

- rj45
- advice from BigEd, joanlluch, oldben, robfinch, robinsonb5, DiTBho, MichaelM on AnyCPU

## License

This project is licensed under the MIT License - see the LICENSE file for details.

NOTE: The youtube videos do not fall under this license. They are under the standard youtube copyright, and I ask you not re-publish them elsewhere. You can create your own youtube videos with the content in this repo, however.

## Acknowledgments

- [Ben Eater](https://eater.net/)
- Many awesome folks on [AnyCPU](http://anycpu.org/forum/)
- John Lluch's [CPU74](https://github.com/John-Lluch/CPU74/)
- [Dieter "ttlworks"](http://www.6502.org/users/dieter/)
- So many inspirations I can't hope to enumerate them all
