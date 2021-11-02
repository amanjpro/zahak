// https://github.com/dgryski/go-prefetch
TEXT ·_prefetch(SB),4,$0-8
        MOVQ  e+0(FP), AX
        PREFETCHNTA (AX)
        RET
