# Zahak

![Build Status](https://github.com/amanjpro/zahak/workflows/Go/badge.svg) [![Join the chat at https://gitter.im/Zahak-Chess-Engine/zahak](https://badges.gitter.im/Zahak-Chess-Engine/zahak.svg)](https://gitter.im/Zahak-Chess-Engine/zahak?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

<img src="zahak_logo.svg" width="300"/>


A UCI compatible chess AI written in Go. Still work in progress. Zahak is also
compliant with [open-benchmark](https://github.com/AndyGrant/OpenBench).

# The name

Zahak (or Zahhak or Azhi Dahak) is an evil figure in Iranian/Kurdish/Perisan
mythology, evident in ancient Iranian folklore as Azhi Dahāka, the name by
which he also appears in the texts of the Avesta.  Legend has it, that he had two
giant snakes on his shoulders and he had to feed them two human brains on
daily basis, you can read more about him
[here](https://en.wikipedia.org/wiki/Zahhak)

# Play Zahak online

Zahak is new to LiChess, you can play him and be impressed with him. His
LiChess handle is [zahak_engine](https://lichess.org/@/zahak_engine). He is
currently running on an old RaspberryPi device, so do not expect a truly
amazing performance. But, hopefully he will be online 24/7.

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

Zahak is participating in [ZaTour tournament series](https://zatour.amanj.me) for open source and original chess engines
that are rated lower than 2800 as per CCRL.

He is also participating in the tournaments arranged by [Chess Engine
Diaries](https://chessengines.blogspot.com/), and several others.

# Rating

Zahak is in the [CCRL Blitz ratings](https://ccrl.chessdom.com/ccrl/404/) as
well as the [4040
rating](https://ccrl.chessdom.com/ccrl/4040/cgi/compare_engines.cgi?family=Zahak&print=Rating+list&print=Results+table&print=LOS+table&print=Ponder+hit+table&print=Eval+difference+table&print=Comopp+gamenum+table&print=Overlap+table&print=Score+with+common+opponents).
He is also listed in [Rebel's Gambit rating (i.e.
GRL)](http://rebel13.nl/grl-40-2.html), in time control 40 moves in 2 minutes.
__Please Note__ Up until including Zahak 5.0, GRL tests are done with the engine's
default hash size, which means it doesn't really represent the actual 
strength/weakness of the engines in the list.  As of version 4.0, Zahak is also
listed on [CEGT 40/4 Rating List](http://www.cegt.net/blitz.htm). Also, starting
from version 5.0, Zahak is listed on [Fast GM rating list 60+06 too](http://www.fastgm.de/60-0.60.html).

The history of Zahak's rating is summerized here:

| **version** |    **CCRL Blitz Rating**     |   **GRL 40/2 Rating**   |   **CEGT 40/4 Rating**    |   **CCRL 40/40 Rating**   |  **Fast GM 60+06** | **Fast GM 10m+6s** |
|-------------|------------------------------|-------------------------|---------------------------|---------------------------|--------------------|--------------------|
| 6.x         | N/A                          | 2785                    | 2664                      | 2800 (unstable rating)    | N/A                | 2720               |
| 5.0         | 2730                         | 2686                    | 2553                      | 2676                      | 2505               | N/A                |
| 4.0         | 2570                         | 2522                    | 2417                      | 2568 (unstable rating)    | N/A                | N/A                |
| 3.0         | 2407                         | 2378                    | N/A                       | N/A                       | N/A                | N/A                |
| 2.0.0       | 2105 (unstable rating)       | N/A                     | N/A                       | N/A                       | N/A                | N/A                |
| 1.0.0       | 2011                         | N/A                     | N/A                       | N/A                       | N/A                | N/A                |
| 0.3.0       | 1922                         | N/A                     | N/A                       | N/A                       | N/A                | N/A                |
| 0.2.1       | 1824                         | N/A                     | N/A                       | N/A                       | N/A                | N/A                |

# Implemented Features:

## Core Features

- UCI Support
- (Magic) Bitboards
- Multi-stage move generation
- Transposition Table
- Pawnhash
- PolyGlot opening book
- Compliant with OpenBench

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

### Selectivity
- Late Move Pruning
- Null-Move Pruning
- Delta Pruning
- Reverse Futility Pruning
- Futility Pruning
- Late Move Reduction
- Razoring
- Check Extensions
- Internal Iterative Deepening (for PV nodes only)
- Internal Iterative Reduction (for non-PV nodes only)
- SEE pruning both in QS and normal search
- ProbCut
- History Leaf Pruning
- Singular Extension
- Multi-Cut

## Evaluation

- NNUE
- Draw pattern recognition: RvN, RvB, RNvR, RBvR

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

# Opening Books

Currently only PolyGlot is supported. Then engine doesn't come with any books,
but you can attach your favourite one easily by passing the path to `-book`
command: `zahak -book PATH_TO_BOOK`.

A bunch of free books are available [here](https://github.com/michaeldv/donna_opening_books)

# Building

To build the project, simply run `make build`, testing with `make test`, and running with `make run`.
Other features exist, for example you can run `perft` with `./zahak -perft` or profile it with `./zahak -profile`.
You can also run it in perfttree mode with `./zahak -preft-tree`.

# Acknowledgement

Zahak wouldn't have been possible without:
- [VICE videos](https://www.youtube.com/playlist?list=PLZ1QII7yudbc-Ky058TEaOstZHVbT-2hg)
- [Chess Programming Wiki](https://www.chessprogramming.org/)
- Niels Abildskov the author of [Loki](https://github.com/BimmerBass/Loki), who helped me with Texel Tuning
- [Nasrin Zaza](https://www.linkedin.com/in/nasrin-zaza/) for the amazing logo
- OpenSource engines like: Weiss, Ethereal, CounterGo, Cheng and Berserk (in no specific order)
- OpenBench community on Discord
- No4b for helping me with some evaluation terms
