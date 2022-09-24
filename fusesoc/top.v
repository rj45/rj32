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
  localparam [11:0] h_res = 720;
  localparam [11:0] h_fp = 36;
  localparam [11:0] h_sync = 72;
  localparam [11:0] h_bp = 108;
  localparam h_neg = 1'b1;
  localparam [11:0] v_res = 400;
  localparam [11:0] v_fp = 1;
  localparam [11:0] v_sync = 3;
  localparam [11:0] v_bp = 42;
  localparam v_neg = 1'b0;
`elsif RES_720x480
  localparam [11:0] h_res = 720;
  localparam [11:0] h_fp = 16;
  localparam [11:0] h_sync = 62;
  localparam [11:0] h_bp = 60;
  localparam h_neg = 1'b1;
  localparam [11:0] v_res = 480;
  localparam [11:0] v_fp = 9;
  localparam [11:0] v_sync = 6;
  localparam [11:0] v_bp = 30;
  localparam v_neg = 1'b1;
`elsif RES_1280x720
  localparam [11:0] h_res = 1280;
  localparam [11:0] h_fp = 48;
  localparam [11:0] h_sync = 32;
  localparam [11:0] h_bp = 80;
  localparam h_neg = 1'b1;
  localparam [11:0] v_res = 720;
  localparam [11:0] v_fp = 7;
  localparam [11:0] v_sync = 8;
  localparam [11:0] v_bp = 6;
  localparam v_neg = 1'b1;
`elsif RES_1024x600
  localparam [11:0] h_res = 1024;
  localparam [11:0] h_fp = 40;
  localparam [11:0] h_sync = 104;
  localparam [11:0] h_bp = 144;
  localparam h_neg = 1'b1;
  localparam [11:0] v_res = 600;
  localparam [11:0] v_fp = 3;
  localparam [11:0] v_sync = 10;
  localparam [11:0] v_bp = 11;
  localparam v_neg = 1'b0;
`elsif RES_640x480 // 640x480
  localparam [11:0] h_res = 640;
  localparam [11:0] h_fp = 16;
  localparam [11:0] h_sync = 96;
  localparam [11:0] h_bp = 48;
  localparam h_neg = 1'b1;
  localparam [11:0] v_res = 480;
  localparam [11:0] v_fp = 10;
  localparam [11:0] v_sync = 2;
  localparam [11:0] v_bp = 33;
  localparam v_neg = 1'b1;
`elsif RES_640x360
  localparam [11:0] h_res = 640;
  localparam [11:0] h_fp = 24;
  localparam [11:0] h_sync = 56;
  localparam [11:0] h_bp = 80;
  localparam h_neg = 1'b1;
  localparam [11:0] v_res = 360;
  localparam [11:0] v_fp = 3;
  localparam [11:0] v_sync = 5;
  localparam [11:0] v_bp = 13;
  localparam v_neg = 1'b1;
`else // 640x400
  localparam [11:0] h_res = 640;
  localparam [11:0] h_fp = 16;
  localparam [11:0] h_sync = 96;
  localparam [11:0] h_bp = 48;
  localparam h_neg = 1'b1;
  localparam [11:0] v_res = 400;
  localparam [11:0] v_fp = 12;
  localparam [11:0] v_sync = 2;
  localparam [11:0] v_bp = 35;
  localparam v_neg = 1'b0;
`endif

  // localparam vga_width = 720;
  // localparam vga_height = 480;
  // localparam vga_refresh = 60;
  // localparam vga_xadjust = 0;
  // localparam vga_yadjust = 0;
  // localparam yminblank         = vga_height/64; // for minimal blank space
  // localparam pixel_f           = 27000000;
  // localparam yframe            = vga_height+yminblank;
  // localparam xframe            = pixel_f/(vga_refresh*yframe);
  // localparam xblank            = xframe-vga_width;
  // localparam yblank            = yframe-vga_height;
  // localparam hsync_front_porch = xblank/3;
  // localparam hsync_pulse_width = xblank/3;
  // localparam hsync_back_porch  = xblank-hsync_pulse_width-hsync_front_porch+vga_xadjust;
  // localparam vsync_front_porch = yblank/3;
  // localparam vsync_pulse_width = yblank/3;
  // localparam vsync_back_porch  = yblank-vsync_pulse_width-vsync_front_porch+vga_yadjust;

  // localparam [11:0] h_res = vga_width;
  // localparam [11:0] h_fp = hsync_front_porch;
  // localparam [11:0] h_sync = hsync_pulse_width;
  // localparam [11:0] h_bp = hsync_back_porch;
  // localparam h_neg = 1'b0;
  // localparam [11:0] v_res = vga_height;
  // localparam [11:0] v_fp = vsync_front_porch;
  // localparam [11:0] v_sync = vsync_pulse_width;
  // localparam [11:0] v_bp = vsync_back_porch;
  // localparam v_neg = 1'b0;

  vdp vgavdp(
    .clk(clk_vga),
    .rst(1'b0),

    .h_fp(h_res - 1),
    .h_sync(h_res + h_fp - 1),
    .h_bp(h_res + h_fp + h_sync - 1),
    .h_total(h_res + h_fp + h_sync + h_bp - 1),
    .h_neg(h_neg),

    .v_fp(v_res - 1),
    .v_sync(v_res + v_fp - 1),
    .v_bp(v_res + v_fp + v_sync - 1),
    .v_total(v_res + v_fp + v_sync + v_bp - 1),
    .v_neg(v_neg),

    .r(r),
    .g(g),
    .b(b),
    .hsync(vga_hs),
    .vsync(vga_vs),
    .en_disp(vga_de)

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