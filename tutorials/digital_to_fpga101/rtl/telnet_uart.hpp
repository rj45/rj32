/*
Telnet UART module for simulation with verilator

Original code at: github.com/rj45/rj32

Copyright (c) 2021 rj45 (github.com/rj45)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

#include <stdio.h>
#include <memory.h>
#include <string.h>
#include <arpa/inet.h>
#include <sys/socket.h>

// make sure this is the same as the values in uart.v
const int clk_freq = 29625000;
const int baud_rate = 115200;

const int port = 2023;
const int baud_tick = clk_freq / baud_rate;

// states for the rx state machine
enum rx_states {
  rx_idle,
  rx_bits,
  rx_stop
};

// states for the tx state machine
enum tx_states {
  tx_idle,
  tx_start,
  tx_bits,
  tx_stop
};

struct telnet_uart_rx {
  char buf[1];
  rx_states state;
  int ticks;
  int bit;
};

struct telnet_uart_tx {
  char buf[1];
  tx_states state;
  int ticks;
  int bit;
};

struct telnet_uart{
  int srv;
  int sock;

  struct telnet_uart_rx rx;
  struct telnet_uart_tx tx;
};

void telnet_uart_init(struct telnet_uart *telnet, Vfpga101 *top) {
  memset(telnet, 0, sizeof(telnet_uart));
  telnet->rx.state = rx_idle;
  telnet->tx.state = tx_idle;

  // open the socket handle
  telnet->srv = socket(AF_INET, SOCK_STREAM, 0);
  if (telnet->srv < 0) {
    printf("socket failed to open: %d\n", errno);
    exit(1);
  }

  // this prevents problems restarting the program waiting
  // for the socket to timeout by telling the OS we will
  // reuse the socket
  int opt = 1;
  if (setsockopt(telnet->srv, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)) != 0) {
    printf("failed setting socket reuse: %d\n", errno);
    exit(1);
  }

  // set up the socket address struct
  struct sockaddr_in addr;
  memset(&addr, 0, sizeof(struct sockaddr_in));
  addr.sin_family = AF_INET;
  addr.sin_port = htons(port);
  addr.sin_addr.s_addr = htonl(INADDR_ANY);

  // bind to 0.0.0.0:`port`
  // if this fails, it's probably something already
  // using port `port`
  if (bind(telnet->srv, (const sockaddr*)&addr, sizeof(addr)) != 0) {
    printf("bind failed: %d\n", errno);
    exit(1);
  }

  // start listening for connections from telnet
  // if this fails, probably something is listening on `port`
  if (listen(telnet->srv, 1) != 0) {
    printf("listen failed: %d\n", errno);
    exit(1);
  }

  printf("\nTo connect to uart do:\n");
  printf("telnet localhost %d\nWaiting... ", port);

  // wait for an incomming connection from telnet
  telnet->sock = accept(telnet->srv, 0, 0);
  if (telnet->sock < 0) {
    printf("\naccept failed: %d\n", errno);
    exit(1);
  }
  printf("Connected!\n\n");

  // rx starts high
  top->rx = 1;
}

void telnet_uart_cleanup(struct telnet_uart *telnet) {
  close(telnet->sock);
  close(telnet->srv);
}


bool telnet_uart_step_tx(int sock, struct telnet_uart_tx *tx, Vfpga101 *top) {
  // handle the tx state machine
  switch (tx->state) {
  case tx_idle:
    {
      // try a non-blocking receive of a byte
      int bytes = recv(sock, tx->buf, 1, MSG_DONTWAIT);
      if (bytes == 1) {
        if (tx->buf[0] == 4 || tx->buf[0] == 6) {
          // exit on a control code (ctrl-D or ctrl-C or whatever)
          return true;
        }

        tx->state = tx_start;
        printf("%c", tx->buf[0]);
      } else if (bytes < 0 && errno != EAGAIN && errno != EWOULDBLOCK) {
        return true;
      }
    }
    break;

  // start transmitting start bit
  case tx_start:
    top->rx = 0;
    tx->ticks = baud_tick;
    tx->bit = -1;
    tx->state = tx_bits;
    break;

  // the bits are being transmitted
  case tx_bits:
    tx->ticks--;
    if (tx->ticks <= 0) {
      tx->ticks = baud_tick;

      // send the next bit
      top->rx = tx->buf[0] & 1;
      tx->buf[0] >>= 1;

      tx->bit++;
      if (tx->bit >= 8) {
        // done now transmit stop bit, stop bit is high
        top->rx = 1;
        tx->state = tx_stop;
      }
    }
    break;

  case tx_stop:
    tx->ticks--;
    if (tx->ticks <= 0) {
      // stop bit is done, return to idle
      tx->state = tx_idle;
    }
  }

  return false;
}

bool telnet_uart_step_rx(int sock, struct telnet_uart_rx *rx, Vfpga101 *top) {
  // handle the rx state machine
  switch (rx->state) {
  case rx_idle:
    if (top->tx == 0) {
      // seeing a start bit
      rx->state = rx_bits;
      rx->buf[0] = 0;
      rx->bit = 0;

      // first bit is at 1.5 baud ticks
      rx->ticks = baud_tick + (baud_tick / 2);
    }
    break;
  case rx_bits:
    rx->ticks--;
    if (rx->ticks <= 0) {
      // time to sample the next bit
      rx->ticks = baud_tick;
      rx->buf[0] >>= 1;

      // check bit and shift it in as the topmost bit
      if (top->tx == 1) {
        // fancy way to set the top bit
        rx->buf[0] = (rx->buf[0] & 0x7f) | 0x80;
      } else {
        // clear the top bit
        rx->buf[0] &= 0x7f;
      }

      rx->bit++;
      if (rx->bit >= 8) {
        // done, receive the stop bit
        rx->state = rx_stop;

        // wait till after the stop bit
        rx->ticks = baud_tick + (baud_tick / 2);

        printf("%c", rx->buf[0]);

        // send to the telnet side
        if (send(sock, rx->buf, 1, 0) != 1) {
          // socket probably disconnected
          printf("\nFailed to send\n");
          return true;
        }
      }
    }
    break;

  case rx_stop:
    rx->ticks--;
    if (rx->ticks <= 0) {
      // stop bit fully received, return to idle
      rx->state = rx_idle;
    }
    break;
  }

  return false;
}


bool telnet_uart_step(struct telnet_uart *telnet, Vfpga101 *top) {
  if(telnet_uart_step_tx(telnet->sock, &telnet->tx, top)) {
    return true;
  }

  if(telnet_uart_step_rx(telnet->sock, &telnet->rx, top)) {
    return true;
  }

  return false;
}
