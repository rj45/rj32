`default_nettype none

module sram (
  input        cs,
  input        we,
  input        oe,
  input [16:0] addr,
  input [15:0] data_i,
  output [15:0] data_o
);

  reg [15:0] mem [0:131071];

  // initial $readmemh("./finch.hex",mem);

  assign data_o = (!cs && !oe) ? mem[addr] : 16'bz;

  always_latch @ (we, cs, addr, data_i) begin
    if (!cs && !we)
      mem[addr] = data_i;
  end

  always @(we, oe) begin
    if (!we && !oe) begin
      $error("we and oe conflict!");
    end
  end

endmodule

module top_verilator (
  input clk_cpu,
  input clk_vga,

  input run,
  input step,

  output led1,
  output led2,
  output led3,

  output error,
  output halt,

  output [7:0] r,
  output [7:0] g,
  output [7:0] b,
  output vga_hs,
  output vga_vs,
  output vga_de
);

  wire v_we;
  wire [16:0] v_adr;
  reg [15:0] v_datr;
  wire [15:0] v_datr_in;
  wire [15:0] v_datw_in;
  reg [15:0] v_datw;
  wire [15:0] v_dat;
  wire v_oe_pin;
  wire v_oe_sram;

  wire flash_ss;
	wire flash_miso;
	wire flash_mosi;


  assign led1 = 0;
  assign led2 = 0;
  assign led3 = 0;

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

		.flash_ss(flash_ss),
    .flash_miso(flash_miso),
    .flash_mosi(flash_mosi),

    .halt(halt),
    .error(error)
  );

  // spiflash flashrom (
  //   .csb(flash_ss),
  //   .clk(clk_vga),
  //   .io0(flash_mosi), // MOSI
  //   .io1(flash_miso) // MISO
  // );

  sram testmem (
    .we(!v_we),
    .oe(!v_oe_sram),
    .cs(0),
    .addr(v_adr),
    .data_i(v_datw),
    .data_o(v_datr_in)
  );

  always @(posedge clk_vga) begin
    v_datr <= v_datr_in;
    v_datw <= v_datw_in;
  end

endmodule