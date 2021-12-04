; func(c byte)
putc:
  UART_OUTPUT_COUNT = 0x007 ; read-only
  UART_DATA_OUT     = 0x005 ; write-only

  .wait:
    in t0, [UART_OUTPUT_COUNT]
    cmp t0, 0xFF
    br.u.g .wait

  out [UART_DATA_OUT], a0
