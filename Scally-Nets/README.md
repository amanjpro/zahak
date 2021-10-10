Smaller / Weaker NETs provided by Alan Cooper (Al) aka Scally, for use as Levels for Human Play for Picochess or any Chess GUI
  
These were all built on a Raspberry Pi 4 booted from a SD Drive with a M.2 backup drive using the latest Zahak Engine to generate the EPDs, Cutechess & Zahak to produce the PGNs, fengen to produce the FENs & zahak-trainer to produce the new NETs 
 
Run 1 was built on small input data and was originally graded around 2000 elo after a winning a small Cutechess tournament against Zahak v0.2.1, Zahak v0.3.0 & Zahak v1.0.0, however a wrong lower epoch defective NET was tested, later the correct NET was tested, graded round 2050 (Scally2050.nn from epoch-36.nnue included) 
 
Run 2 was built on a smaller NET and lower depth with more input data and was graded around 2100 elo after winning a Cutechess tournament against the same 3 engines as above (scally2100.nn from epoch-147.nnue included) 
 
Run 3 was built on an even smaller NET with even more input data and was graded around 2000 elo after winning a Cutechess tournament aginst Zahak v1.0.0 & Zahak v2.0.0 (scally2000.nn from epoch-1523.nnue included) 
 
A later Cutechess tournament between the 3 above NETs, confirmed the grading order

More NETs will be added soon. I hope to provide NETs from 1600 - 2100 elo 
 
Zahak Net Build Summary by Run Number:
======================================
 
30295959 zahak1.epds --> 2.1G  
run1 started at 211005-1332 with 5120 games at depth 9  
Wrote 432217 FENs  
Network Size = 768 x 128 x 1  
The Best NET with the lowest Validation Cost is : epoch-36.nnue  
  
51263866 zahak2.epds --> 3.6G  
run2 started at 211006-1555 with 10240 games at depth 7  
Wrote 696618 FENs  
Network Size = 768 x 64 x 1  
The Best NET with the lowest Validation Cost is : epoch-147.nnue  
  
52810187 zahak3.epds --> 3.7G  
run3 started at 211008-2035 with 20480 games at depth 7  
Wrote 1349435 FENs  
Network Size = 768 x 32 x 1  
The Best NET with the lowest Validation Cost is : epoch-1523.nnue  
 
 
Thanks to Amanj for being very patient in helping me understand the processes  
 
 
Scally 
 
