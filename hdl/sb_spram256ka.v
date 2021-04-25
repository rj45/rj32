/* verilator lint_off UNUSED */
/* verilator lint_off UNDRIVEN */
/* verilator lint_off DECLFILENAME */
(* blackbox *)
module SB_SPRAM256KA(
	input [13:0] ADDRESS,
	input [15:0] DATAIN,
	input [3:0] MASKWREN,
	input WREN,
	input CHIPSELECT,
	input CLOCK,
	input STANDBY,
	input SLEEP,
	input POWEROFF,
	output [15:0] DATAOUT
);
endmodule
/* verilator lint_on UNUSED */
/* verilator lint_on UNDRIVEN */