module prog_bram (
    input             clk_i,
    input       [7:0] addr_i,
    output reg [15:0] data_o
);

reg [15:0] mem [0:255];

initial $readmemh("./test.hex",mem);

always @(negedge clk_i) begin
    data_o <= mem[addr_i];
end

endmodule

module tile_bram (
    input             clk_i,
    input      [10:0] addr_i,
    output reg [15:0] data_o
);

reg [15:0] mem [0:2047];

initial $readmemh("./tiles.hex",mem);

always @(posedge clk_i) begin
    data_o <= mem[addr_i];
end

endmodule

module map_bram (
    input             clk_i,
    input      [9:0] addr_i,
    output reg [7:0] data_o
);

reg [7:0] mem [0:511];

initial $readmemh("./tilemap.hex",mem);

always @(posedge clk_i) begin
    data_o <= mem[addr_i];
end

endmodule

module fbcolour_bram (
    input             clk_i,
    input      [8:0] addr_i,
    output reg [3:0] data_o
);

reg [3:0] mem [0:511];

initial $readmemh("./fbcolour.hex",mem);

always @(posedge clk_i) begin
    data_o <= mem[addr_i];
end

endmodule

module framebuf_bram(
    input wire clk,
    input wire wr_en,
    input wire [8:0] rd_addr,
    input wire [8:0] wr_addr,
    input wire [7:0] data_in,
    output reg [7:0] data_out
);
   reg [7:0] memory [0:511];
   integer i;

   initial begin
      for(i = 0; i <= 511; i=i+1) begin
         memory[i] = 8'b0;
      end
   end

   always @(posedge clk) begin
      if (wr_en) begin
         memory[wr_addr] <= data_in;
      end
      data_out <= memory[rd_addr];
   end
endmodule
