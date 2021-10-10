`default_nettype none

module top_colorlight (
  input clk_25m,

  output led_o,

  output [3:0] hdmi_dp,
  output [3:0] hdmi_dn//,


  // output flash_ss,
  // output flash_sck,
  // output flash_mosi_dq0,
  // input  flash_miso_dq1
);
  wire clk_cpu;
  wire clk_vga;
  wire clk_5x_vga;

  wire [7:0] p1;
  wire [7:0] p4;
  wire [7:0] p3;

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

  reg	[25:0]	counter;
  always @(posedge clk_5x_vga)
    counter <= counter + 1'b1;
  assign led_o = ~vga_hs; //counter[25];

  wire vga_blank;
  assign vga_blank = ~vga_de;

  // raw DVI output via the HDMI connector
  dvi dvi_inst(
    .clk_i(clk_vga),
    .rst_i(~locked),
    .clk_x5_i(clk_5x_vga),
    .vga_red_i(r),
    .vga_green_i(g),
    .vga_blue_i(b),
    .vga_blank_i(vga_blank),
    .vga_hsync_i(vga_hs),
    .vga_vsync_i(vga_vs),
    .dvi_dp_o(hdmi_dp),
    .dvi_dn_o(hdmi_dn)
  );


  assign dm_dat_i = 16'h0000;


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


  localparam vga_width = 640;
  localparam vga_height = 400;
  localparam vga_refresh = 60;
  localparam vga_xadjust = 0; // adjust -3..+3 to fix sync issues
  localparam vga_yadjust = 0;

  function integer F_find_next_f(input integer f);
    if(25000000>=f)
      F_find_next_f=25000000;
    else if(27000000>=f)
      F_find_next_f=27000000;
    else if(40000000>=f)
      F_find_next_f=40000000;
    else if(50000000>=f)
      F_find_next_f=50000000;
    else if(54000000>=f)
      F_find_next_f=54000000;
    else if(60000000>=f)
      F_find_next_f=60000000;
    else if(65000000>=f)
      F_find_next_f=65000000;
    else if(75000000>=f)
      F_find_next_f=75000000;
    else if(80000000>=f)
      F_find_next_f=80000000;  // overclock
    else if(100000000>=f)
      F_find_next_f=100000000; // overclock
    else if(108000000>=f)
      F_find_next_f=108000000; // overclock
    else if(120000000>=f)
      F_find_next_f=120000000; // overclock
  endfunction

  // localparam xminblank         = vga_width/64; // initial estimate
  // localparam yminblank         = vga_height/64; // for minimal blank space
  // localparam min_pixel_f       = vga_refresh*(vga_width+xminblank)*(vga_height+yminblank);
  // localparam pixel_f           = F_find_next_f(min_pixel_f);
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

  localparam pixel_f           = 27000000;


  wire locked;

  // clock generator
  wire [3:0] clocks;
  assign clk_5x_vga = clocks[0];
  assign clk_vga = clocks[1];
  assign clk_cpu = clk_vga;
  ecp5pll
  #(
    .in_hz(25000000),
    .out0_hz(pixel_f*5),
    .out0_tol_hz(1000000*5),
    .out1_hz(pixel_f),
    .out1_tol_hz(1000000)
  )
  ecp5pll_inst
  (
    .clk_i(clk_25m),
    .clk_o(clocks),
    .locked(locked)
  );


endmodule