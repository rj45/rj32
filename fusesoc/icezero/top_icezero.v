`default_nettype none

module top_icezero (
  input clk_100m,
  input [7:0] p1,

  output led1,
  output led2,
  output led3,

  output [7:0] p4,
  output [7:0] p3//,

  // output sram_ce,
  // output sram_we,
  // output sram_oe,
  // output sram_lb,
  // output sram_ub,
  // output [17:0] sram_a,
  // inout  [15:0] sram_d,

  // output flash_ss,
  // output flash_sck,
  // output flash_mosi_dq0,
  // input  flash_miso_dq1
);
  wire clk_cpu;
  wire clk_vga;

  wire [7:0] r;
  wire [7:0] g;
  wire [7:0] b;
  wire vga_hs;
  wire vga_vs;
  wire vga_de;

  reg [15:0] dm_dat_i;
  wire [15:0] dm_dat_o;
  wire [15:0] dm_adr;
  wire dm_we;
  wire dm_req;
  reg  dm_ack;

  always @(posedge clk_cpu)
  begin
    dm_ack <= dm_req;
  end

  wire v_we;
  wire [17:0] v_adr;
  reg [15:0] v_datr;
  reg [15:0] v_datw;
  wire v_oe_pin;
  wire v_oe_sram;
  wire [15:0] v_datr_in;
  wire [15:0] v_datw_in;

  wire run;
  wire step;

  assign run = p1[4];
  assign step = p1[6];

  assign led1 = 0;
  assign led2 = 0;
  assign led3 = 0;

// Are we using the 24bpp HDMI DVI PMOD?
`ifdef hdmi24bpp
  // Pinout Legend
  // Pin   DBus LDat HDat
  // ---------------------
  // p3[0] D11  G3   R7
  // p3[1] D9   G1   R5
  // p3[2] D7   B7   R3
  // p3[3] D5   B5   R1
  // p3[4] D10  G2   R6
  // p3[5] D8   G0   R4
  // p3[6] D6   B6   R2
  // p3[7] D4   B4   R0
  //
  // p4[0] D3   B3   G7
  // p4[1] D1   B1   G5
  // p4[2] CK   --   --
  // p4[3] HS   --   --
  // p4[4] D2   B2   G6
  // p4[5] D0   B0   G4
  // p4[6] DE   --   --
  // p4[7] VS   --   --
  SB_IO #(
    .PIN_TYPE(6'b01_0000)  // PIN_OUTPUT_DDR
  ) dvi_ddr_iob [15:0](
    .PACKAGE_PIN ({p3[0], p3[1],  p3[2],  p3[3],
                  p3[4],  p3[5],  p3[6],  p3[7],
                  p4[0],  p4[1],  p4[2],  p4[3],
                  p4[4],  p4[5],  p4[6],  p4[7]}),
    .D_OUT_0     ({r[7],   r[5],   r[3],   r[1],
                  r[6],   r[4],   r[2],   r[0],
                  g[7],   g[5],   1'b1,   vga_hs,
                  g[6],   g[4],   vga_de, vga_vs}),
    .D_OUT_1     ({g[3],   g[1],   b[7],   b[5],
                  g[2],   g[0],   b[6],   b[4],
                  b[3],   b[1],   1'b0,   vga_hs,
                  b[2],   b[0],   vga_de, vga_vs}),
    .OUTPUT_CLK  (clk_vga)
  );
`else
  // 12b Module - Facing module pins
  //      ----------------------------        ----------------------------
  //     | 0-R3 1-R1 2-G3 3-G1 GND 3V |      | 0-B3 1-ck 2-B0 3-HS GND 3V |
  //     | 4-R2 5-R0 6-G2 7-G0 GND 3V |      | 4-B2 5-B1 6-DE 7-VS GND 3V |
  //  ___|____________________________|______|____________________________|__
  // |       1 bit squared HDMI 12bpp color PMOD board                      |
  //  -----------------------------------------------------------------------
  //       pmod_*_*<0> = r[3]                    pmod_*_*<0> = b[3]
  //       pmod_*_*<1> = r[1]                    pmod_*_*<1> = ck
  //       pmod_*_*<2> = g[3]                    pmod_*_*<2> = b[0]
  //       pmod_*_*<3> = g[1]                    pmod_*_*<3> = hs
  //       pmod_*_*<4> = r[2]                    pmod_*_*<4> = b[2]
  //       pmod_*_*<5> = r[0]                    pmod_*_*<5> = b[1]
  //       pmod_*_*<6> = g[2]                    pmod_*_*<6> = de
  //       pmod_*_*<7> = g[0]                    pmod_*_*<7> = vs
  // wire [3:0] pr;
  // wire [3:0] pg;
  // wire [3:0] pb;

  // assign pr = r[7:4];
  // assign pg = g[7:4];
  // assign pb = b[7:4];

  // assign {p3[0],  p3[1],  p3[2],  p3[3],  p3[4],  p3[5],  p3[6],  p3[7]} =
  //        {pr[3],  pr[1],  pg[3],  pg[1],  pr[2],  pr[0],  pg[2],  pg[0]};
  // assign {p4[0],  p4[1],  p4[2],  p4[3],  p4[4],  p4[5],  p4[6],  p4[7]} =
  //        {pb[3],  clk_vga,pb[0],  vga_hs, pb[2],  pb[1],  vga_de, vga_vs};

  SB_IO #(
    .PIN_TYPE(6'b01_0000)  // PIN_OUTPUT_DDR
  ) dvi_clk_iob (
    .PACKAGE_PIN (p4[1]),
    .D_OUT_0     (1'b0),
    .D_OUT_1     (1'b1),
    .OUTPUT_CLK  (clk_vga)
  );

  SB_IO #(
    .PIN_TYPE(6'b01_0100)  // PIN_OUTPUT_REGISTERED
  ) dvi_data_iob [14:0] (
    .PACKAGE_PIN ({p3[0],  p3[1],  p3[2],  p3[3],  p3[4],  p3[5],  p3[6],  p3[7],
                   p4[0],          p4[2],  p4[3],  p4[4],  p4[5],  p4[6],  p4[7]}),
    .D_OUT_0     ({r[7],   r[5],   g[7],   g[5],   r[6],   r[4],   g[6],   g[4],
                   b[7],           b[4],   vga_hs, b[6],   b[5],   vga_de, vga_vs}),
    .OUTPUT_CLK  (clk_vga)
  );
`endif

  // todo: figure out SRAM
  assign dm_dat_i = 16'h0000;


  // assign sram_ub = 0;
  // assign sram_lb = 0;
  // assign sram_ce = 0;
  // assign sram_oe = !v_oe_sram;
  // assign sram_we = !v_we;

  // assign sram_a = v_adr;

  // SB_IO #(
	// 	.PIN_TYPE(6'b 1010_01),
	// 	.PULLUP(1'b 0)
  // ) databuf [15:0] (
	// 	.PACKAGE_PIN(sram_d),
	// 	.OUTPUT_ENABLE(v_oe_pin),
	// 	.D_OUT_0(v_datw),
	// 	.D_IN_0(v_datr_in)
	// );

  // always @(posedge clk_vga) begin
  //   v_datr <= v_datr_in;
  //   v_datw <= v_datw_in;
  // end

  // assign flash_sck = clk_vga;

  top topmod(
    .run(run),
    .step(step),

    .clk_cpu(clk_cpu),
    .clk_vga(clk_vga),

    .dm_dat_i(dm_dat_i),
    .dm_dat_o(dm_dat_o),
    .dm_adr(dm_adr),
    .dm_we(dm_we),
    .dm_req(dm_req),
    .dm_ack(dm_ack),

    .r(r),
    .g(g),
    .b(b),
    .vga_hs(vga_hs),
    .vga_vs(vga_vs),
    .vga_de(vga_de),

    .v_we(v_we),
    .v_adr(v_adr),
    .v_dat_o(v_datw_in),
    .v_dat_i(v_datr),
    .v_oe_pin(v_oe_pin),
    .v_oe_sram(v_oe_sram),

		// .flash_ss(flash_ss),
    // .flash_miso(flash_miso),
    // .flash_mosi(flash_mosi),

    // .halt(led1),
    // .error(led3)
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
`elsif RES_640x360
  // 100 MHz -> 25.781 MHz (goal: 25.75 MHz)
  localparam divr = 4'b0011; // 3
  localparam divf = 7'b0100000; // 32
  localparam divq = 3'b101; // 5
  localparam rang = 3'b010; // 2
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

		.PLLOUTGLOBALA(clk_vga),
    .PLLOUTGLOBALB(clk_cpu)
  );

endmodule