Smaller / Weaker NETs provided by Alan (Al) Cooper aka Scally, for use as Levels for Human Play for Picochess or any Chess GUI
 
These were all built on a Raspberry Pi 4 booted from a SD Drive with a M.2 backup drive using the Zahak Engine, Cutechess, fengen & zahak-trainer 
 
Run 1 was built on small input data and was originally graded around 2000 elo but was found to be defective (NET not included)
 
Run 2 was built on a smaller NET with more input data and is graded around 2100 elo (scally147.nn included)
 
Run 3 was built on an even smaller NET with even more input data and is also graded around 2100 elo (scally1523.nn included)
This is probably the best NET to use at this time
 
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
 
