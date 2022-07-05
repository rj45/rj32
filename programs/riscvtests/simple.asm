; *****************************************************************************
; simple.S
; -----------------------------------------------------------------------------

; This is the most basic self checking test. If your simulator does not
; pass thiss then there is little chance that it will pass any of the
; more complicated self checking tests.

fence
li gp, 1
li a7, 93
li a0, 0
ecall
