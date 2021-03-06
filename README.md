# Zahak

![Build Status](https://github.com/amanjpro/zahak/workflows/Go/badge.svg) [![Join the chat at https://gitter.im/Zahak-Chess-Engine/zahak](https://badges.gitter.im/Zahak-Chess-Engine/zahak.svg)](https://gitter.im/Zahak-Chess-Engine/zahak?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

A UCI compatible chess AI written in Go. Still work in progress.

# Implemented Features:

- UCI Support
- Bitboards
- Transposition Table
- Alpha-Beta search
- Quiescence Search
- Iterative Deepening
- PV Search and PV
- Zero Windows
- Aspiration Window with PVS
- Static Exchange Evaluation
- Multi-Cut Pruning
- Null-Move Pruning
- Delta Pruning
- Reverse Futility Pruning
- Extended Futility Pruning
- Late Move Reduction
- Razoring
- Killer Moves Heuristics
- Move History Heuristics
- Check Extensions
- Internal Iterative Deepening

# Building

To build the project, simply run `make build`, testing with `make test`, and running with `make run`.
Other features exist, for example you can run `perft` with `./zahak -perft` or profile it with `./zahak -profile`.
You can also run it in perfttree mode with `./zahak -preft-tree`.
