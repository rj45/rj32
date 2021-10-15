/*
The MIT License (MIT)

Copyright (c) 2021 rj45 (github.com/rj45)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

#include <stdio.h>
#include <verilated.h>
#include "Vfpga101.h"
#include "telnet_uart.hpp"

#if VM_TRACE
#include <verilated_fst_c.h>
#endif

int main(int argc, char *argv[])
{
  #if VM_TRACE
  const std::unique_ptr<VerilatedContext> contextp{new VerilatedContext};
  #endif

  Verilated::commandArgs(argc, argv);


  // initialize Verilog module
  Vfpga101 *top = new Vfpga101;

  // initialize the telnet uart module
  struct telnet_uart telnet;
  telnet_uart_init(&telnet, top);

#if VM_TRACE
  // if tracing turn that on
  Verilated::traceEverOn(true);
  VerilatedFstC* tfp = new VerilatedFstC;
  top->trace(tfp, 4);
  tfp->open("run.fst");
#endif

  // the initial clock and eval
  top->clk = 0;
  top->eval();

#if VM_TRACE
  tfp->dump(contextp->time());
#endif


  bool done = false;

  // main loop
  while (!done)
  {
    // cycle the clock
    top->clk = 0;
    top->eval();

#if VM_TRACE
    contextp->timeInc(1);
    tfp->dump(contextp->time());
#endif

    top->clk = 1;
    top->eval();

#if VM_TRACE
    contextp->timeInc(1);
    tfp->dump(contextp->time());
#endif

    done = telnet_uart_step(&telnet, top);
  }

  // clean up and close out

#if VM_TRACE
  tfp->close();
#endif

  top->final(); // simulation done

  telnet_uart_cleanup(&telnet);

  return 0;
}
