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
  output P2_10
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
  wire [15:0] L;
  wire [15:0] R;
  wire [1:0] rd;
  wire [1:0] rs;
  wire jump;
  wire halt;
  wire stall;
  wire error;
  wire skip;

  wire	[30:0]		db1;
  wire	[30:0]		db2;

  wire step_debounced;
  wire run_debounced;
  wire nclock;
  wire red;
  wire green;
  wire blue;

  assign nred = ~red;
  assign ngreen = ~green;
  assign nblue = ~blue;

  debouncer step_deb(
    .i_clk(clock),
    .i_in(~step),
    .o_debounced(step_debounced),
    .o_debug(db1)
  );

  debouncer run_deb(
    .i_clk(clock),
    .i_in(~run),
    .o_debounced(run_debounced),
    .o_debug(db2)
  );

  rj32 cpu(
    //inputs
    .clock(clock),
    .step(step_debounced),
    .run(run_debounced),
    .D_prog(D_prog),
    .D_in(D_in),

    // outputs
    .R0(R0),
    .R1(R1),
    .R2(R2),
    .R3(R3),
    .PC(PC),
    // .halt(halt),
    .halt(green),

    .error(error),
    // .skip(skip),
    .skip(red),
    .rd_valid(rd_valid),
    .rs_valid(rs_valid),
    .op(op),
    .cond(cond),
    .stall(stall),
    .L(L),
    .R(R),
    .rd(rd),
    .rs(rs),
    // .jump(jump),
    .jump(blue),
    .A_prog(A_prog),
    .clock_m(nclock),
    .A_data(A_data),
    .D_out(D_out),
    .w_en(w_en),
  );

  SB_SPRAM256KA dataram (
    .ADDRESS(A_data),
    .DATAIN(D_out),
    .MASKWREN({w_en, w_en, w_en, w_en}),
    .WREN(w_en),
    .CHIPSELECT(1'b1),
    .CLOCK(~clock),
    .STANDBY(1'b0),
    .SLEEP(1'b0),
    .POWEROFF(1'b1),
    .DATAOUT(D_in)
  );

  bram progmem (
    .clk_i(clock),
    .addr_i(A_prog),
    .data_o(D_prog),
  );

  reg [31:0] cnt;
  always @(posedge clock) cnt <= cnt+1;

  display r3out (
    .CLK(clock),
    // .byt({3'b0, op[0:4]}),
    // .byt(PC[7:0]),
    .byt(R3[7:0]),
    .P2_1(P2_1),
    .P2_2(P2_2),
    .P2_3(P2_3),
    .P2_4(P2_4),
    .P2_7(P2_7),
    .P2_8(P2_8),
    .P2_9(P2_9),
    .P2_10(P2_10)
  );

endmodule