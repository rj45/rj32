module vgatop (
  input clock,
  output nred,
  output ngreen,
  // output nblue,

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

  wire red;
  wire green;
  wire blue;

  assign nred = ~red;
  assign ngreen = ~green;
  assign nblue = ~blue;

  wire [3:0] r;
  wire [3:0] g;
  wire [3:0] b;
  wire vga_hs;
  wire vga_vs;
  wire vga_de;

  wire [10:0] fontA;
  reg [15:0] fontD;


`ifdef RES_720x400
  localparam [11:0] res_H = 720;
  localparam [11:0] fp_H = 36;
  localparam [11:0] sync_H = 72;
  localparam [11:0] bp_H = 108;
  localparam neg_H = 1;
  localparam [11:0] res_V = 400;
  localparam [11:0] fp_V = 1;
  localparam [11:0] sync_V = 3;
  localparam [11:0] bp_V = 42;
  localparam neg_V = 0;
  // 12 MHz -> 35.25 MHz (goal: 35.5 MHz)
  localparam divf = 7'b0101110; // 46
  localparam divq = 3'b100; // 4
`elsif RES_720x480
  localparam [11:0] res_H = 720;
  localparam [11:0] fp_H = 16;
  localparam [11:0] sync_H = 62;
  localparam [11:0] bp_H = 60;
  localparam neg_H = 1;
  localparam [11:0] res_V = 480;
  localparam [11:0] fp_V = 9;
  localparam [11:0] sync_V = 6;
  localparam [11:0] bp_V = 30;
  localparam neg_V = 1;
  // 12 MHz -> 27.00 MHz (goal: 27.00 MHz)
  localparam divf = 7'b1000111; // 71
  localparam divq = 3'b101; // 5
`else // 640x480
  localparam [11:0] res_H = 640;
  localparam [11:0] fp_H = 16;
  localparam [11:0] sync_H = 96;
  localparam [11:0] bp_H = 48;
  localparam neg_H = 1;
  localparam [11:0] res_V = 480;
  localparam [11:0] fp_V = 10;
  localparam [11:0] sync_V = 2;
  localparam [11:0] bp_V = 33;
  localparam neg_V = 1;
  // 12 MHz -> 25.125 MHz (goal: 25.175 MHz)
  localparam divf = 7'b1000010; // 66
  localparam divq = 3'b101; // 5
`endif



  // 12 MHz -> 25.125 MHz (goal: 25.175 MHz)
  SB_PLL40_PAD #(
  // SB_PLL40_CORE #(
		.FEEDBACK_PATH("SIMPLE"),
		.DIVR(4'b0000),		// DIVR =  0
		.DIVF(divf),	// DIVF = ?
		.DIVQ(divq),		// DIVQ =  ?
		.FILTER_RANGE(3'b001)	// FILTER_RANGE = 1
	) uut (
		.RESETB(1'b1),
		.BYPASS(1'b0),
		.PACKAGEPIN(clock),
    //.REFERENCECLK(clock),

		.PLLOUTCORE(clk_25m),
  );

  bram progmem (
    .clk_i(clk_25m),
    .addr_i(fontA),
    .data_o(fontD),
  );

  vga_display myvd(
    .clock(clk_25m),

    .res_H(res_H),
    .fp_H(fp_H),
    .sync_H(sync_H),
    .bp_H(bp_H),
    .neg_H(neg_H),
    .res_V(res_V),
    .fp_V(fp_V),
    .sync_V(sync_V),
    .bp_V(bp_V),
    .neg_V(neg_V),

    .CD(fontD),

    .R(r),
    .G(g),
    .B(b),
    .hs(vga_hs),
    .vs(vga_vs),
    .de(vga_de),
    .CA(fontA)
  );

  reg [31:0] cnt;
  always @(posedge clk_25m) cnt <= cnt+1;

  assign red = cnt[21];
  assign green = ~vga_vs;
  // assign blue = r[3];



  assign {P2_1,   P2_2,   P2_3,   P2_4,   P2_7,   P2_8,   P2_9,   P2_10} =
         {r[3],   r[1],   g[3],   g[1],   r[2],   r[0],   g[2],   g[0]};
  assign {P3_1,   P3_2,   P3_3,   P3_4,   P3_7,   P3_8,   P3_9,   P3_10} =
         {b[3],   clk_25m,b[0],   vga_hs, b[2],   b[1],   vga_de, vga_vs};

endmodule