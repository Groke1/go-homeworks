#include "textflag.h"

TEXT ·LowerBound(SB), NOSPLIT, $0
    MOVQ data+0(FP), AX
    MOVQ len+8(FP), BX
    MOVQ value+24(FP), CX
    MOVQ $0, R10 // left
    MOVQ BX, R11 // right

    loop:
    CMPQ R10, R11
    JGE end // if left >= right

    MOVQ R10, R12
    ADDQ R11, R12
    SHRQ $1, R12 // mid

    CMPQ (AX)(R12*8), CX
    JL change_left // array[mid] < value

    MOVQ R12, R11
    JMP loop

    change_left:
    LEAQ 1(R12), R10
    JMP loop

    end:
    MOVQ R10, ret+32(FP)
    RET
