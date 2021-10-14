/*
 *  icebreaker examples - Async uart mirror using pll
 *
 *  Copyright (C) 2018 Piotr Esden-Tempski <piotr@esden.net>
 *  Copyright (C) 2021 rj45 <github.com/rj45>
 *
 *  Permission to use, copy, modify, and/or distribute this software for any
 *  purpose with or without fee is hereby granted, provided that the above
 *  copyright notice and this permission notice appear in all copies.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 *  WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 *  MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 *  ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 *  WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 *  ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 *  OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

`default_nettype none

module uart (
	input  clk,

  input  tx_start,
  input  [7:0] tx_data,
  output tx_busy,
  output tx,

  input  rx,
  output rx_ready,
  output [7:0] rx_data
);


/* local parameters */
// make sure this is the same as the
// "Achieved output frequency" in the `icepll` output
localparam clk_freq = 29_625_000; // 29.625 MHz
localparam baud = 115200;


/* instantiate the rx module */
uart_rx #(clk_freq, baud) urx (
	.clk(clk),
	.rx(rx),
	.rx_ready(rx_ready),
	.rx_data(rx_data)
);

/* instantiate the tx module */
uart_tx #(clk_freq, baud) utx (
	.clk(clk),
	.tx_start(tx_start),
	.tx_data(tx_data),
	.tx(tx),
	.tx_busy(tx_busy)
);

endmodule
