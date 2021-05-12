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

module font_bram (
    input             clk_i,
    input      [10:0] addr_i,
    output reg [15:0] data_o
);

reg [15:0] mem [0:2047];

initial $readmemh("./font16seg.hex",mem);

always @(posedge clk_i) begin
    data_o <= mem[addr_i];
end

endmodule