# C-lib

A bunch of helper functions to help me generate Go-ASM code. The procedure is simple:

- Install c2goasm, and its friends:

```
go get -u github.com/minio/asm2plan9s
go get -u github.com/minio/c2goasm
go get -u github.com/klauspost/asmfmt/cmd/asmfmt
```
- Install yasm 1.2 from sources: https://yasm.tortall.net/releases/Release1.2.0.html
- Install Clang:

```
$ sudo apt-get install clang
```
- Compile the C code to assembly:
```
$ clang -S -mavx -ffast-math -masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -O3 -c clib/avx.c
```
- Convert the assembly code to GoASM:
```
$ c2goasm -a -f clib/avx.s engine/nnue_instructions_amd64.s
```

Enjoy!

