module icezero_top (
  input clk_100m,
  input [7:0] p1,

  output led1,
  output led2,
  output led3,

  output [7:0] p4,
  output [7:0] p3
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

  wire run;
  wire step;

  assign run = p1[4];
  assign step = p1[6];

  assign led1 = 0;
  assign led2 = 0;
  assign led3 = 0;

  assign {p3[0],  p3[1],  p3[2],  p3[3],  p3[4],  p3[5],  p3[6],  p3[7]} =
         {r[3],   r[1],   g[3],   g[1],   r[2],   r[0],   g[2],   g[0]};
  assign {p4[0],  p4[1],  p4[2],  p4[3],  p4[4],  p4[5],  p4[6],  p4[7]} =
         {b[3],   clk_25m,b[0],   vga_hs, b[2],   b[1],   vga_de, vga_vs};

  // todo: figure out SRAM
  assign D_in = 16'h0000;

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

`ifdef RES_1280x720
  // 100 MHz -> 73.750 MHz (goal: 74.25 MHz)
  localparam divr = 4'b1001; // 9
  localparam divf = 7'b0111010; // 58
  localparam divq = 3'b011; // 3
  localparam rang = 3'b001; // 1
`elsif RES_800x600
  // 100 MHz -> 40.000 MHz (goal: 40 MHz)
  localparam divr = 4'b0100; // 4
  localparam divf = 7'b0011111; // 31
  localparam divq = 3'b100; // 4
  localparam rang = 3'b010; // 2
`elsif RES_720x400
  // 100 MHz -> 35.417 MHz (goal: 35.5 MHz)
  localparam divr = 4'b0010; // 2
  localparam divf = 7'b0010000; // 16
  localparam divq = 3'b100; // 4
  localparam rang = 3'b011; // 3
`elsif RES_720x480
  // 100 MHz -> 26.953 MHz (goal: 27.00 MHz)
  localparam divr = 4'b0111; // 7
  localparam divf = 7'b1000100; // 68
  localparam divq = 3'b101; // 5
  localparam rang = 3'b001; // 1
`else // 640x480 or 640x400
  // 100 MHz -> 25.312 MHz (goal: 25.175 MHz)
  localparam divr = 4'b1001; // 9
  localparam divf = 7'b1010000; // 80
  localparam divq = 3'b101; // 5
  localparam rang = 3'b001; // 1
`endif

  SB_PLL40_2F_PAD #(
		.FEEDBACK_PATH("SIMPLE"),
		.DIVR(divr),		// DIVR = ?
		.DIVF(divf),   	// DIVF = ?
		.DIVQ(divq),		// DIVQ = ?
		.FILTER_RANGE(rang),	// FILTER_RANGE = ?
    .PLLOUT_SELECT_PORTB("GENCLK_HALF"),
	) uut (
		.RESETB(1'b1),
		.BYPASS(1'b0),
		.PACKAGEPIN(clk_100m),

		.PLLOUTGLOBALA(clk_25m),
    .PLLOUTGLOBALB(clk_12m)
  );

endmodule