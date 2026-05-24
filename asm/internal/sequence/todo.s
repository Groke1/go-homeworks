#include "textflag.h"

TEXT ·Fibonacci(SB), NOSPLIT, $0
    MOVQ n+0(FP), AX
    MOVQ $-1, BX // n - 1
    MOVQ $1, CX // n
    XORQ DX, DX // ind
    loop:
    MOVQ BX, R10
    MOVQ CX, BX
    ADDQ R10, CX
    INCQ DX
    CMPQ DX, AX
    JLE loop

    MOVQ CX, ret+8(FP)
    RET
