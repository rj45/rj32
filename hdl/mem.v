module bram (
    input             clk_i,
    input       [8:0] addr_i,
    output reg [15:0] data_o
);

parameter MEMFILE = "./test.hex";

reg [15:0] mem [0:255];

initial $readmemh(MEMFILE,mem);

always @(negedge clk_i) begin
    data_o <= mem[addr_i];
end

endmodule