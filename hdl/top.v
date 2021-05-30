module top (
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

  reg [15:0] D_prog;
  wire [15:0] D_in;

  wire [15:0] R0;
  wire [15:0] R1;
  wire [15:0] R2;
  wire [15:0] R3;
  wire [7:0] A_prog;
  wire [13:0] A_data;
  wire [15:0] D_out;
  wire w_en;
  wire [7:0] PC;
  wire [4:0] op;
  wire rd_valid;
  wire rs_valid;
  wire [2:0] cond;
  wire [15:0] L_b;
  wire [15:0] R_b;
  wire [15:0] result;
  wire [3:0] rd;
  wire [3:0] rs;
  wire jump;
  wire halt;
  wire stall;
  wire error;
  wire skip;
  wire immv;
  wire [2:0] aluop;
  wire [15:0] imm;

  wire	[30:0]		db1;
  wire	[30:0]		db2;

  wire step_debounced;
  wire run_debounced;
  wire nclock;
  wire red;
  wire green;
  wire blue;

  assign red = 0;
  assign green = 0;
  assign blue = 0;

  assign nred = ~red;
  assign ngreen = ~green;
  assign nblue = ~blue;

  wire [3:0] r;
  wire [3:0] g;
  wire [3:0] b;
  wire vga_hs;
  wire vga_vs;
  wire vga_de;

  wire [10:0] tileA;
  reg [15:0] tileD;

  wire [8:0] mapA;
  reg [7:0] mapD;

  debouncer step_deb(
    .i_clk(clk_12m),
    .i_in(~step),
    .o_debounced(step_debounced),
    .o_debug(db1)
  );

  debouncer run_deb(
    .i_clk(clk_12m),
    .i_in(~run),
    .o_debounced(run_debounced),
    .o_debug(db2)
  );

  rj32 cpu(
    //inputs
    .clock(clk_12m),
    .step(step_debounced),
    .run_slow(run_debounced),
    .D_prog(D_prog),
    .D_in(D_in),
    .run_fast(0),
    .run_faster(0),
    .erun(0),

    // outputs
    .R0(R0),
    .R1(R1),
    .R2(R2),
    .R3(R3),
    .PC(PC),
    .halt(halt),

    .error(error),
    .skip(skip),
    .rd_valid(rd_valid),
    .rs_valid(rs_valid),
    .op(op),
    .cond(cond),
    .stall(stall),
    .L(L_b),
    .R(R_b),
    .result(result),
    .rd(rd),
    .rs(rs),
    .jump(jump),
    .A_prog(A_prog),
    .clock_m(nclock),
    .A_data(A_data),
    .D_out(D_out),
    .w_en(w_en),
    .immv(immv),
    .imm(imm),
    .aluop(aluop),
  );

  wire n_clk_12m;

  assign n_clk_12m = ~clk_12m;

  SB_SPRAM256KA dataram (
    .ADDRESS(A_data),
    .DATAIN(D_out),
    .MASKWREN({w_en, w_en, w_en, w_en}),
    .WREN(w_en),
    .CHIPSELECT(1'b1),
    .CLOCK(n_clk_12m),
    .STANDBY(1'b0),
    .SLEEP(1'b0),
    .POWEROFF(1'b1),
    .DATAOUT(D_in)
  );

  prog_bram progmem (
    .clk_i(clk_12m),
    .addr_i(A_prog),
    .data_o(D_prog),
  );

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
  SB_PLL40_2F_PAD #(
  // SB_PLL40_CORE #(
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
    //.REFERENCECLK(clock),

		.PLLOUTGLOBALA(clk_25m),
    .PLLOUTGLOBALB(clk_12m),
  );

  tile_bram tilemem (
    .clk_i(clk_25m),
    .addr_i(tileA),
    .data_o(tileD),
  );

  map_bram mapmem (
    .clk_i(clk_25m),
    .addr_i(mapA),
    .data_o(mapD),
  );

  vga_blinkenlights vgafrontpanel(
    .clock(clk_25m),
    .clock_slow(clk_12m),

    // .res_H(res_H),
    // .fp_H(fp_H),
    // .sync_H(sync_H),
    // .bp_H(bp_H),
    // .neg_H(neg_H),
    // .res_V(res_V),
    // .fp_V(fp_V),
    // .sync_V(sync_V),
    // .bp_V(bp_V),
    // .neg_V(neg_V),

    .PC(PC),
    .R0(R0),
    .R1(R1),
    .R2(R2),
    .R3(R3),
    .L_b(L_b),
    .R_b(R_b),
    .result(result),
    .skip(skip),
    .jump(jump),
    .op(op),
    .cond(cond),
    .rd(rd),
    .rs(rs),
    .rdv(rd_valid),
    .rsv(rs_valid),
    .halt(halt),
    .error(error),
    .imm(imm),
    .immv(immv),
    .aluop(aluop),

    .TD(tileD),
    .MD(mapD),

    .R(r),
    .G(g),
    .B(b),
    .hs(vga_hs),
    .vs(vga_vs),
    .de(vga_de),
    .TA(tileA),
    .MA(mapA)
  );

  assign {P2_1,   P2_2,   P2_3,   P2_4,   P2_7,   P2_8,   P2_9,   P2_10} =
         {r[3],   r[1],   g[3],   g[1],   r[2],   r[0],   g[2],   g[0]};
  assign {P3_1,   P3_2,   P3_3,   P3_4,   P3_7,   P3_8,   P3_9,   P3_10} =
         {b[3],   clk_25m,b[0],   vga_hs, b[2],   b[1],   vga_de, vga_vs};


  // reg [31:0] cnt;
  // always @(posedge clock) cnt <= cnt+1;

  // display r3out (
  //   .CLK(clock),
  //   // .byt({3'b0, op[0:4]}),
  //   // .byt(PC[7:0]),
  //   .byt(R3[7:0]),
  //   .P2_1(P2_1),
  //   .P2_2(P2_2),
  //   .P2_3(P2_3),
  //   .P2_4(P2_4),
  //   .P2_7(P2_7),
  //   .P2_8(P2_8),
  //   .P2_9(P2_9),
  //   .P2_10(P2_10)
  // );

endmodule