`default_nettype none

module bram(
    input wire clk,
    input wire en_wr,
    input wire [8:0] addr_rd,
    input wire [8:0] addr_wr,
    input wire [7:0] data_wr,
    output reg [7:0] data_rd
);
  reg [7:0] memory [0:511];

  always @(posedge clk) begin
    data_rd <= memory[addr_rd];
  end

  always @(posedge clk) begin
    if (en_wr) begin
      memory[addr_wr] <= data_wr;
    end
  end
endmodule
