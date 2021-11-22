# Nets for Skill Levels

A set of smaller / weaker Zahak NETs by Alan Cooper (Al) aka Scally, for use as
Levels for Human Play for Picochess or any Chess GUI

These were all built on a Raspberry Pi 4 booted from a SD Drive with a M.2
backup drive using the latest Zahak Engine (at the time of data generation) to
generate the EPDs, Cutechess and Zahak to produce the PGNs, fengen to produce
the FENs & zahak-trainer to produce the new NETs

When we first started this Project we believed that the best NETs were those
with the lowest Validation Cost, so I was literally creating hundreds and in
some cases thousands of NETs until the lowest Validation Cost was reached.
Earlier NETs were produced with a smaller NetInputSize of 768, so all work done
on these were scrapped after Amanj increased this size to 769 when adding tempo
into the network architecture.

It later become apparent that we were overtraining these NETs, so we settled
for producing no more than 50 and testing these for strength.

Amanj was producing the strongest NET available, trying NET sizes of 769 x 128
x 1, then 769 x 256 x 1 & now 769 x 512 x 1.  My aim was to try and produce
weaker NETs, so after trying 769 x 128 x 1, then 769 x 64 x 1, I settled for
769 x 32 x 1.

Whilst I was creating and testing my NETs, Amanj was fine tuning his code, so
as soon as I produced NETs that were tested at grade X, they became stronger
with his new code.

I have settled on the following NETs, these are all based on the grade for
Zahak v0.2.1 set at 1824 elo.

Summary of Zahak NET Rebuilds:

```
62,158,444 zahak4.epds --> 4.3G
run4 started at 211101-1234 with 10240 games at depth 9
Wrote 921639 FENs

Rebuilds after NET configuration changes:
Rerun zahak4.epds with different PARMs
Network Size = 769 x 32 x 1
Learning Rate (LR) 0.01 unless stated otherwise
-sigmoid-scale of 2/1024 (0.001953125) 4d
-sigmoid-scale of 3/1024 (0.0029296875) 4e
-sigmoid-scale of 1/1024 (0.0009765625) 4f

All grades based on Zahak v0.2.1 @1824 elo
epoch-6.nnue   4e:  2289.1 elo
epoch-49.nnue  4d:  2189.7 elo [skills_6.nn]
epoch-9.nnue   4d:  2017.5 elo [skills_5.nn]
epoch-5.nnue   4e:  2011.2 elo
epoch-5.5.nnue 4f:  1891 elo (LR of 0.001 restarted from epoch-5)
epoch-4.nnue   4d:  1884.6 elo
epoch-4.nnue   4e:  1837.4 elo (appeared weaker on subsequent tests so not used)
epoch-5.nnue   4d:  1833.8 elo [skills_4.nn]
(Zahak v0.2.1 ----- @1824 elo)
epoch-6.nnue   4f:  1807.6 elo (appeared weaker on subsequent tests so not used)
epoch-5.3.nnue 4f:  1781.8 elo (LR of 0.001 restarted from epoch-5)
epoch-3.7.nnue 4e:  1753.7 elo (LR of 0.001 restarted from epoch-3)
epoch-3.6.nnue 4e:  1695.2 elo (LR of 0.001 restarted from epoch-3)
epoch-3.5.nnue 4e:  1683   elo (LR of 0.001 restarted from epoch-3)
epoch-5.nnue   4f:  1596.1 elo [skills_3.nn]
epoch-3.nnue   4d:  1586   elo
epoch-3.nnue   4e:  1515.5 elo
epoch-4.nnue   4f:  1421.8 elo [skills_2.nn]
epoch-2.nnue   4d:  1212   elo (appeared weaker on subsequent tests so not used)
epoch-3.nnue   4f:   959   elo
epoch-2.nnue   4e:  1108.5 elo [skills_1.nn]
```
Scally
