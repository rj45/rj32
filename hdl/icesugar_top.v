module icesugar_top (
  input clock,
  input step,
  input run,
  output nred,
  output ngreen,
  output nblue,

  output P2_1 ,
  output P2_2 ,
  output P2_3 ,
  output P2_4 ,
  output P2_7 ,
  output P2_8 ,
  output P2_9 ,
  output P2_10,
  output P3_1 ,
  output P3_2 ,
  output P3_3 ,
  output P3_4 ,
  output P3_7 ,
  output P3_8 ,
  output P3_9 ,
  output P3_10
);
  wire clk_12m;
  wire clk_25m;

  wire [3:0] r;
  wire [3:0] g;
  wire [3:0] b;
  wire vga_hs;
  wire vga_vs;
  wire vga_de;

  wire [15:0] D_in;
  wire [15:0] D_out;
  wire [13:0] A_data;
  wire w_en;

  assign nred = 1;
  assign ngreen = 1;
  assign nblue = 1;

  assign {P2_1,   P2_2,   P2_3,   P2_4,   P2_7,   P2_8,   P2_9,   P2_10} =
         {r[3],   r[1],   g[3],   g[1],   r[2],   r[0],   g[2],   g[0]};
  assign {P3_1,   P3_2,   P3_3,   P3_4,   P3_7,   P3_8,   P3_9,   P3_10} =
         {b[3],   clk_25m,b[0],   vga_hs, b[2],   b[1],   vga_de, vga_vs};

  // up5k SPRAM
  SB_SPRAM256KA dataram (
    .ADDRESS(A_data),
    .DATAIN(D_out),
    .MASKWREN({w_en, w_en, w_en, w_en}),
    .WREN(w_en),
    .CHIPSELECT(1'b1),
    .CLOCK(clk_12m),
    .STANDBY(1'b0),
    .SLEEP(1'b0),
    .POWEROFF(1'b1),
    .DATAOUT(D_in)
  );

  top topmod(
    .clk_cpu(clk_12m),
    .clk_vga(clk_25m),
    .step(step),
    .run(run),
    .r(r),
    .g(g),
    .b(b),
    .vga_hs(vga_hs),
    .vga_vs(vga_vs),
    .vga_de(vga_de),
    .A_data(A_data),
    .D_out(D_out),
    .D_in(D_in),
    .w_en(w_en)
  );

`ifdef RES_720x400
  // 12 MHz -> 35.25 MHz (goal: 35.5 MHz)
  localparam divf = 7'b0101110; // 46
  localparam divq = 3'b100; // 4
`elsif RES_720x480
  // 12 MHz -> 27.00 MHz (goal: 27.00 MHz)
  localparam divf = 7'b1000111; // 71
  localparam divq = 3'b101; // 5
`else // 640x480 or 640x400
  // 12 MHz -> 25.125 MHz (goal: 25.175 MHz)
  localparam divf = 7'b1000010; // 66
  localparam divq = 3'b101; // 5
`endif

  SB_PLL40_2F_PAD #(
		.FEEDBACK_PATH("SIMPLE"),
		.DIVR(4'b0000),		// DIVR =  0
		.DIVF(divf),	// DIVF = ?
		.DIVQ(divq),		// DIVQ =  ?
		.FILTER_RANGE(3'b001),	// FILTER_RANGE = 1
    .PLLOUT_SELECT_PORTB("GENCLK_HALF"),
	) uut (
		.RESETB(1'b1),
		.BYPASS(1'b0),
		.PACKAGEPIN(clock),

		.PLLOUTGLOBALA(clk_25m),
    .PLLOUTGLOBALB(clk_12m)
  );

endmodule