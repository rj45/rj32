`default_nettype none

module prog_bram (
    input             clk_i,
    input       [7:0] addr_i,
    output reg [15:0] data_o
);

reg [15:0] mem [0:255];

initial $readmemh("./data/test.hex",mem);

always @(negedge clk_i) begin
    data_o <= mem[addr_i];
end

endmodule

module palette_bram (
    input             clk_i,
    input      [7:0]  addr_i,
    output reg [15:0] data_o
);

reg [15:0] mem [0:255];

initial $readmemh("./data/finchpal.hex",mem);

always @(posedge clk_i) begin
    data_o <= mem[addr_i];
end

endmodule

module linebuffer_bram (
    input             clk_i,
    input             we_i,
    input      [8:0]  addrr_i,
    input      [8:0]  addrw_i,
    input      [31:0] dataw_i,
    output reg [31:0] datar_o
);

reg [31:0] mem [0:511];

always @(posedge clk_i) begin
    datar_o <= mem[addrr_i];
    if (we_i) begin
      mem[addrw_i] <= dataw_i;
    end
end

endmodule


module top (
  input step,
  input run,

  input clk_cpu,
  input clk_vga,

  input [15:0] dm_dat_i,
  output [15:0] dm_dat_o,
  output [15:0] dm_adr,
  output dm_we,
  output dm_req,
  input  dm_ack,

  output [7:0] r,
  output [7:0] g,
  output [7:0] b,
  output vga_hs,
  output vga_vs,
  output vga_de,

  output v_we,
  output v_oe_pin,
  output v_oe_sram,
  output [17:0] v_adr,
  input  [15:0] v_dat_i,
  output [15:0] v_dat_o,

  output flash_ss,
  output flash_mosi,
  input flash_miso,

  output halt,
  output error
);
  wire stall;
  wire skip;

  wire nclock;

  reg [15:0] prg_dat;
  wire [7:0] prg_adr;

  wire [7:0] pal_adr;
  reg [15:0] pal_dat;

  wire [8:0] lbr_adr;
  reg [31:0] lbr_dat;

  wire [8:0] lbw_adr;
  wire [31:0] lbw_dat;
  wire lbw_we;

  prog_bram progmem (
    .clk_i(clk_cpu),
    .addr_i(prg_adr),
    .data_o(prg_dat)
  );

  palette_bram palette (
    .clk_i(clk_vga),
    .addr_i(pal_adr),
    .data_o(pal_dat)
  );

  linebuffer_bram linebuffer (
    .clk_i(clk_vga),
    .addrr_i(lbr_adr),
    .datar_o(lbr_dat),
    .addrw_i(lbw_adr),
    .dataw_i(lbw_dat),
    .we_i(lbw_we)
  );

  wire [25:0] db;

  rj32 cpu(
    //inputs
    .clock(clk_cpu),
    .step(step),
    .run_slow(run),
    .D_prog(prg_dat),
    .D_in(dm_dat_i),
    .run_fast(1'b0),
    .run_faster(1'b0),
    .erun(1'b0),
    .ack(dm_ack),

    // outputs
    .halt(halt),
    .error(error),
    .skip(skip),
    .stall(stall),
    .A_prog(prg_adr),
    .clock_m(nclock),
    .A_data(dm_adr),
    .D_out(dm_dat_o),
    .w_en(dm_we),
    .db(db),
    .req(dm_req)
  );

`ifdef RES_720x400
  localparam [11:0] res_H = 720;
  localparam [11:0] fp_H = 36;
  localparam [11:0] sync_H = 72;
  localparam [11:0] bp_H = 108;
  localparam neg_H = 1'b1;
  localparam [10:0] res_V = 400;
  localparam [10:0] fp_V = 1;
  localparam [10:0] sync_V = 3;
  localparam [10:0] bp_V = 42;
  localparam neg_V = 1'b0;
`elsif RES_720x480
  localparam [11:0] res_H = 720;
  localparam [11:0] fp_H = 16;
  localparam [11:0] sync_H = 62;
  localparam [11:0] bp_H = 60;
  localparam neg_H = 1'b1;
  localparam [10:0] res_V = 480;
  localparam [10:0] fp_V = 9;
  localparam [10:0] sync_V = 6;
  localparam [10:0] bp_V = 30;
  localparam neg_V = 1'b1;
`elsif RES_1280x720
  localparam [11:0] res_H = 1280;
  localparam [11:0] fp_H = 48;
  localparam [11:0] sync_H = 32;
  localparam [11:0] bp_H = 80;
  localparam neg_H = 1'b1;
  localparam [11:0] res_V = 720;
  localparam [11:0] fp_V = 7;
  localparam [11:0] sync_V = 8;
  localparam [11:0] bp_V = 6;
  localparam neg_V = 1'b1;
`elsif RES_1024x600
  localparam [11:0] res_H = 1024;
  localparam [11:0] fp_H = 40;
  localparam [11:0] sync_H = 104;
  localparam [11:0] bp_H = 144;
  localparam neg_H = 1'b1;
  localparam [10:0] res_V = 600;
  localparam [10:0] fp_V = 3;
  localparam [10:0] sync_V = 10;
  localparam [10:0] bp_V = 11;
  localparam neg_V = 1'b0;
`elsif RES_640x480 // 640x480
  localparam [11:0] res_H = 640;
  localparam [11:0] fp_H = 16;
  localparam [11:0] sync_H = 96;
  localparam [11:0] bp_H = 48;
  localparam neg_H = 1'b1;
  localparam [10:0] res_V = 480;
  localparam [10:0] fp_V = 10;
  localparam [10:0] sync_V = 2;
  localparam [10:0] bp_V = 33;
  localparam neg_V = 1'b1;
`elsif RES_640x360
  localparam [11:0] res_H = 640;
  localparam [11:0] fp_H = 24;
  localparam [11:0] sync_H = 56;
  localparam [11:0] bp_H = 80;
  localparam neg_H = 1'b1;
  localparam [10:0] res_V = 360;
  localparam [10:0] fp_V = 3;
  localparam [10:0] sync_V = 5;
  localparam [10:0] bp_V = 13;
  localparam neg_V = 1'b1;
`else // 640x400
  localparam [11:0] res_H = 640;
  localparam [11:0] fp_H = 16;
  localparam [11:0] sync_H = 96;
  localparam [11:0] bp_H = 48;
  localparam neg_H = 1'b1;
  localparam [10:0] res_V = 400;
  localparam [10:0] fp_V = 12;
  localparam [10:0] sync_V = 2;
  localparam [10:0] bp_V = 35;
  localparam neg_V = 1'b0;
`endif

  vdp vgavdp(
    .clk(clk_vga),
    .rst(0),

    .hfp(res_H),
    .hsn(res_H + fp_H),
    .hbp(res_H + fp_H + sync_H),
    .httl(res_H + fp_H + sync_H + bp_H),
    .hneg(neg_H),

    .vfp(res_V),
    .vsn(res_V + fp_V),
    .vbp(res_V + fp_V + sync_V),
    .vttl(res_V + fp_V + sync_V + bp_V),
    .vneg(neg_V),

    .R(r),
    .G(g),
    .B(b),
    .hs(vga_hs),
    .vs(vga_vs),
    .de(vga_de)

    //     .p_adr(p_adr),
    // .p_dat(p_dat),

    // .lbr_adr(lbr_adr),
    // .lbr_dat(lbr_dat),
    // .lbw_adr(lbw_adr),
    // .lbw_dat(lbw_dat),
    // .lbw_we(lbw_we),

    // .v_adr(v_adr),
    // .v_dat(v_dat_i),

    // .oe_pin(v_oe_pin),
    // .we_sram(v_we),
    // .oe_sram(v_oe_sram),


    // .ss(flash_ss),
    // .mosi(flash_mosi),
    // .miso(flash_miso),
    // .fl_d(v_dat_o)

  );

endmodule