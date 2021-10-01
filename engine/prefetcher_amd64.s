TEXT ·_prefetch(SB), $0-8
       MOVQ       e+0(FP), AX
       PREFETCHNTA (AX)
       RET
