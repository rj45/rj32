module pll(
	input  clock_in,
	output clock_out,
	output locked
	);

  assign locked = 1'b1;
  assign clock_out = clock_in;

endmodule
