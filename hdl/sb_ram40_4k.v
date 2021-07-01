/* verilator lint_off UNUSED */
/* verilator lint_off UNDRIVEN */
/* verilator lint_off DECLFILENAME */
(* blackbox *)
module SB_RAM40_4K (
	output [15:0] RDATA,
	input         RCLK,
	input         RCLKE ,
	input         RE,
	input  [10:0] RADDR,
	input         WCLK,
	input         WCLKE ,
	input         WE,
	input  [10:0] WADDR,
	input  [15:0] MASK,
	input  [15:0] WDATA
);

	parameter WRITE_MODE = 0;
	parameter READ_MODE = 0;

endmodule
/* verilator lint_on UNUSED */
/* verilator lint_on UNDRIVEN */