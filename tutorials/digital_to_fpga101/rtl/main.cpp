// Copyright (C)2021 rj45
// MIT Licensed

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


