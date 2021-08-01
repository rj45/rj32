module top (
  input step,
  input run,

  input clk_cpu,
  input clk_vga,

  input [15:0] D_in,

  output [3:0] r,
  output [3:0] g,
  output [3:0] b,
  output vga_hs,
  output vga_vs,
  output vga_de,

  output [15:0] D_out,
  output [13:0] A_data,
  output w_en,

  output halt,
  output error
);
  reg [15:0] D_prog;
  wire [7:0] A_prog;
  wire stall;
  wire skip;

  wire nclock;

  wire [10:0] tileA;
  reg [15:0] tileD;

  wire [9:0] mapA;
  reg [7:0] mapD;

  reg [7:0] fbr_D;
  wire fbwen;
  wire [8:0] fbw_A;
  wire [7:0] fbw_D;
  wire [8:0] fbr_A;
  reg [3:0] fbc_D;

  wire [25:0] db;

  reg ack;
  wire req;

  always @(posedge clk_cpu)
  begin
    ack <= req;
  end


  rj32 cpu(
    //inputs
    .clock(clk_cpu),
    .step(step),
    .run_slow(run),
    .D_prog(D_prog),
    .D_in(D_in),
    .run_fast(1'b0),
    .run_faster(1'b0),
    .erun(1'b0),
    .ack(ack),

    // outputs
    .halt(halt),
    .error(error),
    .skip(skip),
    .stall(stall),
    .A_prog(A_prog),
    .clock_m(nclock),
    .A_data(A_data),
    .D_out(D_out),
    .w_en(w_en),
    .db(db),
    .req(req)
  );

  prog_bram progmem (
    .clk_i(clk_cpu),
    .addr_i(A_prog),
    .data_o(D_prog)
  );

`ifdef RES_720x400
  localparam [11:0] res_H = 720;
  localparam [11:0] fp_H = 36;
  localparam [11:0] sync_H = 72;
  localparam [11:0] bp_H = 108;
  localparam neg_H = 1'b1;
  localparam [11:0] res_V = 400;
  localparam [11:0] fp_V = 1;
  localparam [11:0] sync_V = 3;
  localparam [11:0] bp_V = 42;
  localparam neg_V = 1'b0;
`elsif RES_720x480
  localparam [11:0] res_H = 720;
  localparam [11:0] fp_H = 16;
  localparam [11:0] sync_H = 62;
  localparam [11:0] bp_H = 60;
  localparam neg_H = 1'b1;
  localparam [11:0] res_V = 480;
  localparam [11:0] fp_V = 9;
  localparam [11:0] sync_V = 6;
  localparam [11:0] bp_V = 30;
  localparam neg_V = 1'b1;
`elsif RES_640x480 // 640x480
  localparam [11:0] res_H = 640;
  localparam [11:0] fp_H = 16;
  localparam [11:0] sync_H = 96;
  localparam [11:0] bp_H = 48;
  localparam neg_H = 1'b1;
  localparam [11:0] res_V = 480;
  localparam [11:0] fp_V = 10;
  localparam [11:0] sync_V = 2;
  localparam [11:0] bp_V = 33;
  localparam neg_V = 1'b1;
`else // 640x400
  localparam [11:0] res_H = 640;
  localparam [11:0] fp_H = 16;
  localparam [11:0] sync_H = 96;
  localparam [11:0] bp_H = 48;
  localparam neg_H = 1'b1;
  localparam [11:0] res_V = 400;
  localparam [11:0] fp_V = 12;
  localparam [11:0] sync_V = 2;
  localparam [11:0] bp_V = 35;
  localparam neg_V = 1'b0;
`endif

  tile_bram tilemem (
    .clk_i(clk_vga),
    .addr_i(tileA),
    .data_o(tileD)
  );

  map_bram mapmem (
    .clk_i(clk_vga),
    .addr_i(mapA),
    .data_o(mapD)
  );

  framebuf_bram fbmem (
    .clk(clk_vga),
    .rd_addr(fbr_A),
    .wr_addr(fbw_A),
    .data_in(fbw_D),
    .wr_en(fbwen),
    .data_out(fbr_D)
  );

  fbcolour_bram colourmem (
    .clk_i(clk_vga),
    .addr_i(fbr_A),
    .data_o(fbc_D)
  );

  vga_blinkenlights vgafrontpanel(
    .clock(clk_vga),
    .db(db),

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

    .TD(tileD),
    .MD(mapD),
    .fbr_D(fbr_D),
    .fbc_D(fbc_D),

    .R(r),
    .G(g),
    .B(b),
    .hs(vga_hs),
    .vs(vga_vs),
    .de(vga_de),
    .TA(tileA),
    .MA(mapA),
    .fbr_A(fbr_A),
    .fbw_A(fbw_A),
    .fbw_D(fbw_D),
    .wen(fbwen)
  );

endmodule