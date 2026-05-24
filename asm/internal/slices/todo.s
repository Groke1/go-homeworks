#include "textflag.h"

TEXT ·Sum(SB), NOSPLIT, $0
    MOVQ data+0(FP), AX
    MOVQ len+8(FP), BX
    XORQ CX, CX // index
    XORQ DX, DX // sum
    loop:
    MOVLQSX (AX)(CX*4), R10
    ADDQ R10, DX
    INCQ CX
    CMPQ CX, BX
    JL loop

    MOVQ DX, ret+24(FP)
    RET
