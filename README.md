# Zahak

![Build Status](https://github.com/amanjpro/zahak/workflows/Go/badge.svg)

<img src="zahak_logo.svg" width="300"/>


A UCI compatible chess AI written in Go. Still work in progress. 

Courtesy to the OpenBench community, Zahak is now part of the [official
OpenBench instance](http://chess.grantnet.us/index/)

# The name

Zahak (or Zahhak or Azhi Dahak) is an evil figure in Iranian/Kurdish/Perisan
mythology, evident in ancient Iranian folklore as Azhi Dahāka, the name by
which he also appears in the texts of the Avesta.  Legend has it, that he had two
giant snakes on his shoulders and he had to feed them two human brains on
daily basis, you can read more about him
[here](https://en.wikipedia.org/wiki/Zahhak)

# Play Zahak online

Zahak is new to LiChess, you can play him and be impressed with him. His
LiChess handle is [zahak_engine](https://lichess.org/@/zahak_engine).

# Play Zahak on your Android Phone/Desktop

Zahak is a bare chess engine AI, that means it doesn't come with any GUI
interface.  That also means, it is easy to plug it into any chess GUI that
supports UCI protocol.

- [Arena Chess GUI](http://www.playwitharena.de/)
- [CuteChess](https://cutechess.com/)
- [Tarrasch](https://www.triplehappy.com/)
- [The Shredder GUI](https://www.shredderchess.com/)
- [Fritz / Chessbase series](https://en.chessbase.com/)
- [Scid vs PC (database)](http://scidvspc.sourceforge.net/)
- [Banksia GUI](https://banksiagui.com/)
- [DroidFish](https://play.google.com/store/apps/details?id=org.petero.droidfish) is a good choice on Android

# Tournaments

Zahak recently got invitation to the [TCEC tournament](https://tcec-chess.com/)!

Zahak is also participating in [ZaTour tournament
series](https://zatour.amanj.me) for open source and original chess engines.
And many other tournaments, like the ones arranged by the amazing Graham Banks.

# Rating

Zahak is in almost all major rating lists, here are the details:


- [CCRL rating list](https://ccrl.chessdom.com/ccrl)

| **version** | [**CCRL Blitz Rating**](https://ccrl.chessdom.com/ccrl/404/cgi/compare_engines.cgi?family=Zahak&print=Rating+list&print=Results+table&print=LOS+table&print=Ponder+hit+table&print=Eval+difference+table&print=Comopp+gamenum+table&print=Overlap+table&print=Score+with+common+opponents)    | [**CCRL Blitz (8CPU) Rating**](https://ccrl.chessdom.com/ccrl/404/cgi/compare_engines.cgi?family=Zahak&print=Rating+list&print=Results+table&print=LOS+table&print=Ponder+hit+table&print=Eval+difference+table&print=Comopp+gamenum+table&print=Overlap+table&print=Score+with+common+opponents)    |   [**CCRL 40/40 Rating**](https://ccrl.chessdom.com/ccrl/4040/cgi/compare_engines.cgi?family=Zahak&print=Rating+list&print=Results+table&print=LOS+table&print=Ponder+hit+table&print=Eval+difference+table&print=Comopp+gamenum+table&print=Overlap+table&print=Score+with+common+opponents) |    [**CCRL 40/4040 (4CPU) Rating**](https://ccrl.chessdom.com/ccrl/4040/cgi/compare_engines.cgi?family=Zahak&print=Rating+list&print=Results+table&print=LOS+table&print=Ponder+hit+table&print=Eval+difference+table&print=Comopp+gamenum+table&print=Overlap+table&print=Score+with+common+opponents) 
|-------------|------------------------------|------------------------------|---------------------------|--------------------|
| 9.x         | 3278                         | 3406                         | 3213                      | 3282               |
| 8.x         | 3133                         | N/A                          | 3098                      | N/A                |
| 7.x         | 2964 (32 bit: 2897)          | N/A                          | 2938                      | 3006               |
| 6.x         | 2833                         | N/A                          | 2800 (unstable rating)    | N/A                |
| 5.0         | 2730                         | N/A                          | 2676                      | N/A                |
| 4.0         | 2570                         | N/A                          | 2568 (unstable rating)    | N/A                |
| 3.0         | 2407                         | N/A                          | N/A                       | N/A                |
| 2.0.0       | 2105 (unstable rating)       | N/A                          | N/A                       | N/A                |
| 1.0.0       | 2011                         | N/A                          | N/A                       | N/A                |
| 0.3.0       | 1922                         | N/A                          | N/A                       | N/A                |
| 0.2.1       | 1824                         | N/A                          | N/A                       | N/A                |

- [CEGT rating list](http://www.cegt.net)

| **version** |   [**CEGT 40/4 Rating**](http://www.cegt.net/40_4_Ratinglist/40_4_AllVersion/rangliste.html)    |   [**CEGT 40/20 Rating**](http://www.cegt.net/40_40%20Rating%20List/40_40%20All%20Versions/rangliste.html)    |   [**CEGT 5"+3' Rating (Ponder On)**](http://www.cegt.net/5Plus3Rating/5Plus3AllVersion/rangliste.html)    |
|-------------|---------------------------|---------------------------|---------------------------|
| 9.x         | 3182                      | 3162                      | 3199                      |
| 8.x         | 3046                      | 3022                      | 3051                      |
| 7.x         | 2840                      | N/A                       | 2858                      |
| 6.x         | 2664                      | N/A                       | 2676                      |
| 5.0         | 2553                      | N/A                       | N/A                       |
| 4.0         | 2417                      | N/A                       | N/A                       |


- Other well-known rating lists

| **version** |   [**GRL 40/2 Rating**](http://rebel13.nl/history.html#57)   |     [**SPCC Rating**](https://www.sp-cc.de/index.htm)       |     [**BRUCE Rating**](https://e4e6.com/)      |  [**Fast GM 60+06**](http://fastgm.de/60-0.60.html) | [**Fast GM 10m+6s**](http://fastgm.de/10min.html) |
|-------------|-------------------------|---------------------------|---------------------------|--------------------|--------------------|
| 9.x         | 3244                    | 3273                      | 3283                      | N/A                | N/A                |
| 8.x         | 3140                    | N/A                       | 3169                      | N/A                | 3057               |
| 7.x         | 2929                    | N/A                       | 2981                      | 2749               | N/A                |
| 6.x         | 2785                    | N/A                       | 2841                      | 2584               | 2720               |
| 5.0         | 2686                    | N/A                       | 2683                      | 2505               | N/A                |
| 4.0         | 2522                    | N/A                       | N/A                       | N/A                | N/A                |
| 3.0         | 2378                    | N/A                       | N/A                       | N/A                | N/A                |

# Implemented Features:

## Core Features

- UCI Support
- (Magic) Bitboards
- Multi-stage move generation
- Transposition Table
- Pawnhash
- PolyGlot opening book
- Compliant with OpenBench
- Syzygy Support
- MultiPV
- Skill Levels, 1 to 7 (strongest)

## Search

### Basics
- Alpha-Beta search
- Quiescence Search
- Iterative Deepening
- PV Search and PV
- Search with Zero Windows
- Aspiration Window with PVS
- Pondering
- Multi-Threading (LazySMP)

### Move Ordering

- Hash move
- Promotions
- Static Exchange Evaluation followed by LVA-MVV for equal captures according to SEE
- Killer Moves Heuristics
- Countermove Heuristics
- Move History Heuristics
- Countermove History Heuristics
- FollowUp History Heuristics

### Selectivity
- Late Move Pruning
- Null-Move Pruning
- Delta Pruning
- Reverse Futility Pruning
- Futility Pruning
- Late Move Reduction
- Razoring
- Check Extensions
- Internal Iterative Deepening
- SEE pruning both in QS and normal search
- Threat Pruning
- Singular Extension
- Multi-Cut

## Evaluation

- NNUE

# Command line options

```
bash-3.2$ bin/zahak -help
Usage of bin/zahak:
  Commands:
   ./zahak         Runs Zahak in UCI mode
   ./zahak bench   Runs Zahak in OpenBench mode
   
  Options:
  
  -book string
        Path to openning book in PolyGlot (bin) format
  -perft
        Provide this to run perft tests
  -perft-tree
        Run the engine in prefttree mode
  -gen-epds
        Generate opening EPDs for self-play
  -profile
        Run the engine in profiling mode
  -slow
        Run all perft tests, even the very slow tests
  -test-positions string
        Path to EPD positions, used to test the strength of the engine
```

# Skill Levels

Anchored around Rustic Alpha 3, I found that, based on CCRL the ratings will
probably translate to the following:

- Skill Level 1: 1270
- Skill Level 2: 1440
- Skill Level 3: 1630
- Skill Level 4: 1856
- Skill Level 5: 2004
- Skill Level 6: 2074

# Opening Books

Currently only PolyGlot is supported. Then engine doesn't come with any books,
but you can attach your favourite one easily by passing the path to `-book`
command: `zahak -book PATH_TO_BOOK`.

A bunch of free books are available [here](https://github.com/michaeldv/donna_opening_books)

# Training Networks

Please refer to [this guide](https://github.com/amanjpro/zahak/blob/master/training.md).

# Building

To build the project, simply run `make build`, testing with `make test`, and running with `make run`.
Other features exist, for example you can run `perft` with `./zahak -perft` or profile it with `./zahak -profile`.
You can also run it in perfttree mode with `./zahak -preft-tree`.

# Contributors

Thanks to the following for their valuable contributions:

- Basti Dangca: for generating data for training Zahak's network
- Alan Cooper (Scally): for generating weaker networks for different skill levels.

# Acknowledgement

Zahak wouldn't have been possible without:
- [VICE videos](https://www.youtube.com/playlist?list=PLZ1QII7yudbc-Ky058TEaOstZHVbT-2hg)
- [Chess Programming Wiki](https://www.chessprogramming.org/)
- [The official OpenBench instance](http://chess.grantnet.us/index/), and the all the hardware donators (noopwn4ftw and others)
- Aryan Parekh the author of [Bit-Genie](https://github.com/Aryan1508/Bit-Genie), who helped me with NNUE
- Niels Abildskov the author of [Loki](https://github.com/BimmerBass/Loki), who helped me with Texel Tuning
- [Nasrin Zaza](https://www.linkedin.com/in/nasrin-zaza/) for the amazing logo
- OpenSource engines like: Weiss, Ethereal, CounterGo, Cheng and Berserk (in no specific order)
- OpenBench community on Discord
- No4b for helping me with some evaluation terms
