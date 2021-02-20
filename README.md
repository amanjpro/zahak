# Zahak

![Build Status](https://github.com/amanjpro/zahak/workflows/Go/badge.svg)

A UCI comppatible chess AI written in Go. Still work in proress.

# Implemented Features:

- UCI Support
- Bitboards
- Alpha-Beta search
- Quiescence Search
- Iterative Deepining
- PV Search and PV
- Zero Windows
- Delta Pruning
- Null-Move Pruning
- Transposition Table
- Static Exchange Evaluation
- Mult-Cut Pruning
- Reverse Futility Pruning
- Razoring

# Building

To build the project, simply run `make build`, testing with `make test`, and running with `make run`.
Other features exist, for example you can run `perft` with `./zahak -perft` or profile it with `./zahak -profile`.
You can also run it in perfttree mode with `./zahak -preft-tree`.
