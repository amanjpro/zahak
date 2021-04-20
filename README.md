# Zahak

![Build Status](https://github.com/amanjpro/zahak/workflows/Go/badge.svg) [![Join the chat at https://gitter.im/Zahak-Chess-Engine/zahak](https://badges.gitter.im/Zahak-Chess-Engine/zahak.svg)](https://gitter.im/Zahak-Chess-Engine/zahak?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

A UCI compatible chess AI written in Go. Still work in progress.

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

# Rating

Zahak is participating in some tournaments arranged by [Chess Engine
Diaries](https://chessengines.blogspot.com/), he recently advanced to the E7 in
the 44th edition of JCER competition.

He is also listed in the [CCRL ratings](https://ccrl.chessdom.com/ccrl/404/),
his current rating is around 1922.

# Implemented Features:

- UCI Support
- (Magic) Bitboards
- Multi-stage move generation
- Transposition Table
- Alpha-Beta search
- Quiescence Search
- Iterative Deepening
- PV Search and PV
- Zero Windows
- Aspiration Window with PVS
- Static Exchange Evaluation
- Late Move Pruning
- Null-Move Pruning
- Delta Pruning
- Reverse Futility Pruning
- Late Move Reduction
- Razoring
- Killer Moves Heuristics
- Move History Heuristics
- Check Extensions
- Internal Iterative Deepening
- PolyGlot opening book

# Command line options

```
Usage of bin/zahak:
  -book string
    Path to openning book in PolyGlot (bin) format
  -perft
    Provide this to run perft tests
  -perft-tree
    Run the engine in prefttree mode
  -profile
    Run the engine in profiling mode
  -slow
    Run all perft tests, even the very slow tests
  -test-positions
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

Zahak wouldn't have been possible without [VICE videos](https://www.youtube.com/playlist?list=PLZ1QII7yudbc-Ky058TEaOstZHVbT-2hg)
and [Chess Programming Wiki](https://www.chessprogramming.org/)
