# Verilator Simulation

This code lets you simulate the verilog in verilator complete with VGA output.

This is mainly from the [Project F Verilator with SDL blog post](https://projectf.io/posts/verilog-sim-verilator-sdl/).

This runs the simulation slightly faster than Digital can. But it cannot run anywhere near as fast as an FPGA unfortunately.

## Running on linux

Install the build software required:

```sh
apt install build-essential verilator libsdl2-dev
```

On Mac:

```sh
brew install verilator sdl2
```

Then just run make.

### Running on Windows

The linked blog article may be able to help you, sorry I cannot.

### Keyboard Keys

- Space steps
- R runs
- Space will exit run mode back to step mode

## License

Most of the code is copyright:

(C)2021 Will Green, open source software released under the MIT License

With some modifications by me. Thank you Will Green for making this code available!
