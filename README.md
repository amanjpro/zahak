# Zahak

![Build Status](https://github.com/amanjpro/zahak/workflows/Go/badge.svg)

A UCI compatible chess AI written in Go. Still work in progress.

# Implemented Features:

- UCI Support
- Bitboards
- Alpha-Beta search
- Quiescence Search
- Iterative Deepening
- PV Search and PV
- Aspiration Window with PVS
- Zero Windows
- ~Delta Pruning~ Disabled, somehow it makes the search slower
- Null-Move Pruning
- Transposition Table
- Static Exchange Evaluation
- Multi-Cut Pruning
- Reverse Futility Pruning
- Extended Futility Pruning
- Late Move Reduction
- Razoring
- Killer Moves Heuristics
- Move History Heuristics
- Check Extensions

# Building

To build the project, simply run `make build`, testing with `make test`, and running with `make run`.
Other features exist, for example you can run `perft` with `./zahak -perft` or profile it with `./zahak -profile`.
You can also run it in perfttree mode with `./zahak -preft-tree`.
