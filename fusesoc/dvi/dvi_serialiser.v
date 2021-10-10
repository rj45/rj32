//-----------------------------------------------------------------
//                      DVI / HDMI Framebuffer
//                              V0.1
//                     github.com/ultraembedded
//                          Copyright 2020
//
//                 Email: admin@ultra-embedded.com
//
//                       License: MIT
//-----------------------------------------------------------------
// Copyright (c) 2020 github.com/ultraembedded
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//-----------------------------------------------------------------
`default_nettype none
module dvi_serialiser
(
      input       clk_i
    , input       rst_i
    , input       strobe_i
    , input [9:0] data_i
    , output      serial_p_o
    , output      serial_n_o
);

  reg [1:0] data_p, data_n;

  always @(posedge clk_i) data_p[0] <=  data_i[0];
  always @(posedge clk_i) data_n[0] <= ~data_i[0];
  always @(posedge clk_i) data_p[1] <=  data_i[1];
  always @(posedge clk_i) data_n[1] <= ~data_i[1];

  ODDRX1F
  ddr_p_instance
  (
    .D0(data_p[0]),
    .D1(data_p[1]),
    .Q(serial_p_o),
    .SCLK(clk_i),
    .RST(rst_i)
  );
  ODDRX1F
  ddr_n_instance
  (
    .D0(data_n[0]),
    .D1(data_n[1]),
    .Q(serial_n_o),
    .SCLK(clk_i),
    .RST(rst_i)
  );


endmodule