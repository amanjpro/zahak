# Training for Zahak

So you are thinking of trainig networks for Zahak? Great! You know what? If you
produce a net that beats the latest strongest at the time, your net will be
promoted to the master net! Not only that, your name will also be listed under
the authors list! The only condition here that needs to be respected is the
originality of the data, all the evaluation scores **must** come from Zahak's
evaluation/search itself. How do we know if a net is stronger than another? We
run a match between the two nets, in both STC and LTC to find out!

## Network Architecture

Zahak's network architecture is a simple `768*128*1`. The features represent
`(piece_type, piece_color, square number)`, notice that `6 * 2 * 64 = 768`.

## Training

### Data generation

To start training, you need to have a lot of chess positions. Zahak's first net
was trained over 58 million chess positions (FEN). It is important that your
positions cover a wide range of different positions, or else your network will
be overfitting for certain types of positions, while performing very badly in
others.

There are many ways to generate data, the way I did it for Zahak was first to
use `cutechess` to generate many self-play games searching until the depth of 9
per move, with this command:

```
cutechess-cli -tournament gauntlet -concurrency 15 \
  -pgnout zahak_games/PGN_NAME.pgn "fi" \
  -engine conf=zahak_latest tc="inf" depth=9 \
  -engine conf=zahak_latest tc="inf" depth=9 \
  -ratinginterval 1 \
  -recover \
  -event SELF_PLAY_GAMES \
  -draw movenumber=40 movecount=10 score=20 \
  -resign movecount=5 score=1000 \
  -resultformat per-color \
  -openings order=random policy=round file=SOME_BOOK_HERE format="epd" \
  -each proto=uci option.Hash=32 option.Threads=1 \
  -rounds 100000000
```

If you are using the above command (or a version of it), it is at most
important to use a book with millions of different openning positions, the
easiest way to generate a wide range of them is by using Zahak itself:

```
bin/zahak -gen-epds
```

The above command basically spits perfectly randomly generated starting
positions that you can use it as an EPD book. All you need is to capture the
output and store it in file, In Windows you can do this: `bin\zahak -gen-epds >
FILE` while in Linux/Mac you can do this: `bin/zahak -gen-epds | FILE`.

Now I have millions of games! What next?

### Games to Chess Positions 

You now need to convert all the games that you generated above into FENs to be
able to feed it to the network trainer. Easiest way is to use the latest
version of [fengen](https://github.com/amanjpro/fengen/releases), which expects
games stored in PGN format, that has a the score in the very first section of a
move's comment. Something like the following:

```
[Event "self-play-1"]
[Site "?"]
[Date "2021.09.06"]
[Round "28"]
[White "zahak-linux-amd64-6.2"]
[Black "zahak-linux-amd64-6.2"]
[Result "1-0"]
[FEN "rnbqk1nr/1ppp1p1p/4p1pb/p7/P1P5/R6N/1P1PPPPP/1NBQKB1R w Kkq - 0 1"]
[GameDuration "00:00:01"]
[GameEndTime "2021-09-06T17:23:04.131 EDT"]
[GameStartTime "2021-09-06T17:23:02.218 EDT"]
[PlyCount "51"]
[SetUp "1"]
[TimeControl "inf"]

1. Rg3 {-0.17/9 0.027s} Nc6 {+0.29/9 0.022s} 2. Nc3 {-0.14/9 0.014s}
Nf6 {+0.54/9 0.022s} 3. d3 {+0.12/9 0.020s} Bxc1 {+0.39/9 0.006s}
4. Qxc1 {-0.45/9 0.012s} h6 {+0.35/9 0.028s} 5. e4 {-0.04/9 0.037s}
Nh5 {+0.07/9 0.022s} 6. Re3 {+0.09/9 0.052s} Nd4 {-0.20/9 0.042s}
7. Be2 {+0.18/9 0.028s} Nxe2 {-0.05/9 0.021s} 8. Nxe2 {-0.01/9 0.019s}
d6 {+0.03/9 0.022s} 9. O-O {-0.03/9 0.034s} e5 {+0.06/9 0.030s}
10. Qc3 {-0.13/9 0.065s} Qf6 {+0.05/9 0.064s} 11. f4 {-0.04/9 0.069s}
O-O {+0.01/9 0.029s} 12. fxe5 {-0.14/9 0.033s} Qxe5 {-0.28/9 0.012s}
13. d4 {+0.31/9 0.044s} Qe8 {+0.09/9 0.083s} 14. Nhf4 {+0.39/9 0.057s}
Nxf4 {-0.41/9 0.032s} 15. Nxf4 {+0.41/9 0.018s} c6 {-0.59/9 0.059s}
16. e5 {+0.76/9 0.035s} dxe5 {-0.59/9 0.009s} 17. Rxe5 {+0.65/9 0.033s}
Qd7 {-0.67/9 0.030s} 18. d5 {+0.70/9 0.070s} Qd8 {-0.68/9 0.075s}
19. Qd4 {+0.83/9 0.090s} cxd5 {-0.87/9 0.041s} 20. Nxd5 {+0.93/9 0.019s}
Be6 {-0.71/9 0.050s} 21. c5 {+0.78/9 0.12s} Rc8 {-0.90/9 0.042s}
22. Nf6+ {+0.69/9 0.065s} Kg7 {-1.03/9 0.055s} 23. Nd7 {+1.55/9 0.011s}
Bxd7 {-1.39/9 0.022s} 24. Re8+ {+1.50/9 0.011s} Kh7 {-1.63/9 0.034s}
25. Rxf7+ {+M1/2 0s} Rxf7 {-M2/1 0s} 26. Qh8# {0.00/1 0.006s, White mates} 1-0
```

Notice that the score is stored as the first section, and is separated with the
rest with `/`: `{-0.17/9 0.027s}`.

You can use `fengen` in with many threads, and it can recieve multiple files at
the same time:

```
./fengen -help
Usage of ./fengen:
  -input string
    	Comma separated set of paths to PGN files
  -limit int
    	Maximum allowed difference between Quiescence Search result and Static Evaluation, the bigger it is the more tactical positions are included
  -output string
    	Directory to write produced FENs in
  -threas int
    	Number of parallelism (default 8)
```

And example run will be: `./fengen -input self-play-1.pgn,self-play-2.pgn
-output training-data.ep`, to make the positions more tactical. negative values
are ignored, 0 means no tactical position in the training set, the higher the
limit, the more tactical the training data becomes. Zahak's default net is
trained with `limit` set to 0.  You can play with the value of `-limit`, to
make the positions more tactical. negative values are ignored, 0 means no
tactical position in the training set, the higher the limit, the more tactical
the training data becomes. Zahak's default net is trained with `limit` set to
0.

### Finally, training a net

Easiest way to train a net for Zahak is to use the latest version of
[zahak-trainer](https://github.com/amanjpro/zahak-trainer/releases), it expects
data in exactly the same format as what is generated by `fengen`:
`<FEN>;score:<SCORE>;eval:<EVAL>;qs:<QUIESCENCE SEARCH
SCORE>;outcome:[1.0|0.5|0.0]` Only the `fen`, `score` and `outcome` parts are
used, `eval` and `qs` parts do not affect the training procedure. Both `eval`
and `outcome` are reported in the eyes of white.

`zahak-trainer` is somewhat flexible, you can tweak a some of the internal
parametrs to give your net special characteristics:

```
$ ./zahak-trainer -help
Usage of ./zahak-trainer:
  -epochs int
      Number of epochs (default 100)
  -from-net string
      Path to a network, to be used as a starting point
  -hiddens string
      Number of hidden neurons, for multi-layer you can send comma separated numbers (default "128")
  -input-path string
      Path to input dataset (FENs)
  -inputs int
      Number of inputs (default 768)
  -lr float
      Learning Rate (default 0.009999999776482582)
  -network-id int
      A unique id for the network (default 285061292)
  -output-path string
      Final NNUE path directory
  -outputs int
      Number of outputs (default 1)
  -profile
      Profile the trainer
  -sigmoid-scale float
      Sigmoid scale (default 0.00244140625)
```

Unfortunately, due to the limitations of the probing code in Zahak, `hiddens`,
`inputs` `outputs` should be left as the default value, or your net won't be
compatible with the engine. You can continue your training session, from an
already semi-trained network by using `from-net` argument.

Here is how I trained the default net of Zahak 7.0: `./zahak-trainer -epochs
30000 -input-path training-dataset.epd -output-path epoch-nets/`

You can play with `-sigmoid-scale` to give your network different
characteristics, there is really no good way (that I know of) to pick the
perfect number, it is mostly trial and error. I used 3.5 / 1024, but you can as
well use other numbers. The only real test is a match with the best net, if you
beat it, your value is good. Unfortunately,this value also heavily depends on
the dataset, so best value for a dataset is not necessarily the best value for
another.

`lr` is the learning rate, the bigger it is the faster the training reaches a
local minima (a good net basically), but it might not be the best. THe smaller
it is, the slower the training session converges, but the higher the chance it
finds the perfect weights for the net.

How do you know a net is fully trained? You can never be 100% sure really, but
there is a BIG clue for, the trainer after each epoch emits the validation and
training cost, the goal is to make these two costs as low as possible. At
somepoint the training cost goes down, while the validation cost goes up. This
is usually an undesired state, because it means that we are over-training the
net. You always want to pick the epoch with the lowest validation cost:

```
Training and validation cost progression
==============================================================================
Epoch                Training Cost                Validation Cost
==============================================================================
1                      0.075133                      0.074623
2                      0.074407                      0.074167
3                      0.074063                      0.073939
```

The result net of each epoch is stored, so make sure you use the right epoch at
the end of the training session. One tip, always use a high number of epoch
when starting a training session, because it is much easier to cut it short,
than to continue from a half-trained net.
