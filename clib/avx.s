	.text
	.intel_syntax noprefix
	.file	"avx.c"
	.globl	update_hidden                   # -- Begin function update_hidden
	.p2align	4, 0x90
	.type	update_hidden,@function
update_hidden:                          # @update_hidden
# %bb.0:
	push	rbp
	mov	rbp, rsp
	push	r15
	push	r14
	push	r13
	push	r12
	push	rbx
	and	rsp, -8
	sub	rsp, 56
	mov	r12d, dword ptr [rbp + 16]
	test	r12d, r12d
	jle	.LBB0_33
# %bb.1:
	mov	r13d, r12d
	cmp	r12d, 32
	jb	.LBB0_2
# %bb.3:
	lea	rax, [rdi + 4*r13]
	cmp	rax, r9
	jbe	.LBB0_6
# %bb.4:
	lea	rax, [r9 + 4*r13]
	cmp	rax, rdi
	jbe	.LBB0_6
.LBB0_2:
	xor	r14d, r14d
.LBB0_12:
	mov	r10, r14
	not	r10
	add	r10, r13
	mov	rax, r13
	and	rax, 3
	je	.LBB0_14
	.p2align	4, 0x90
.LBB0_13:                               # =>This Inner Loop Header: Depth=1
	mov	ebx, dword ptr [rdi + 4*r14]
	mov	dword ptr [r9 + 4*r14], ebx
	add	r14, 1
	add	rax, -1
	jne	.LBB0_13
.LBB0_14:
	cmp	r10, 3
	jb	.LBB0_16
	.p2align	4, 0x90
.LBB0_15:                               # =>This Inner Loop Header: Depth=1
	mov	eax, dword ptr [rdi + 4*r14]
	mov	dword ptr [r9 + 4*r14], eax
	mov	eax, dword ptr [rdi + 4*r14 + 4]
	mov	dword ptr [r9 + 4*r14 + 4], eax
	mov	eax, dword ptr [rdi + 4*r14 + 8]
	mov	dword ptr [r9 + 4*r14 + 8], eax
	mov	eax, dword ptr [rdi + 4*r14 + 12]
	mov	dword ptr [r9 + 4*r14 + 12], eax
	add	r14, 4
	cmp	r13, r14
	jne	.LBB0_15
	jmp	.LBB0_16
.LBB0_6:
	mov	r14d, r13d
	and	r14d, -32
	lea	rax, [r14 - 32]
	mov	r10, rax
	shr	r10, 5
	add	r10, 1
	test	rax, rax
	je	.LBB0_34
# %bb.7:
	mov	r11, r10
	and	r11, -2
	neg	r11
	xor	ebx, ebx
	.p2align	4, 0x90
.LBB0_8:                                # =>This Inner Loop Header: Depth=1
	vmovups	ymm0, ymmword ptr [rdi + 4*rbx]
	vmovups	ymm1, ymmword ptr [rdi + 4*rbx + 32]
	vmovups	ymm2, ymmword ptr [rdi + 4*rbx + 64]
	vmovups	ymm3, ymmword ptr [rdi + 4*rbx + 96]
	vmovups	ymmword ptr [r9 + 4*rbx], ymm0
	vmovups	ymmword ptr [r9 + 4*rbx + 32], ymm1
	vmovups	ymmword ptr [r9 + 4*rbx + 64], ymm2
	vmovups	ymmword ptr [r9 + 4*rbx + 96], ymm3
	vmovdqu	ymm0, ymmword ptr [rdi + 4*rbx + 128]
	vmovdqu	ymm1, ymmword ptr [rdi + 4*rbx + 160]
	vmovdqu	ymm2, ymmword ptr [rdi + 4*rbx + 192]
	vmovdqu	ymm3, ymmword ptr [rdi + 4*rbx + 224]
	vmovdqu	ymmword ptr [r9 + 4*rbx + 128], ymm0
	vmovdqu	ymmword ptr [r9 + 4*rbx + 160], ymm1
	vmovdqu	ymmword ptr [r9 + 4*rbx + 192], ymm2
	vmovdqu	ymmword ptr [r9 + 4*rbx + 224], ymm3
	add	rbx, 64
	add	r11, 2
	jne	.LBB0_8
# %bb.9:
	test	r10b, 1
	je	.LBB0_11
.LBB0_10:
	vmovdqu	ymm0, ymmword ptr [rdi + 4*rbx]
	vmovdqu	ymm1, ymmword ptr [rdi + 4*rbx + 32]
	vmovdqu	ymm2, ymmword ptr [rdi + 4*rbx + 64]
	vmovdqu	ymm3, ymmword ptr [rdi + 4*rbx + 96]
	vmovdqu	ymmword ptr [r9 + 4*rbx], ymm0
	vmovdqu	ymmword ptr [r9 + 4*rbx + 32], ymm1
	vmovdqu	ymmword ptr [r9 + 4*rbx + 64], ymm2
	vmovdqu	ymmword ptr [r9 + 4*rbx + 96], ymm3
.LBB0_11:
	cmp	r14, r13
	jne	.LBB0_12
.LBB0_16:
	test	r12d, r12d
	jle	.LBB0_33
# %bb.17:
	test	ecx, ecx
	jle	.LBB0_33
# %bb.18:
	mov	eax, ecx
	mov	qword ptr [rsp + 48], rax       # 8-byte Spill
	lea	rax, [r9 + 4*r13]
	mov	qword ptr [rsp + 16], rax       # 8-byte Spill
	lea	rax, [r8 + 4*r13]
	mov	qword ptr [rsp + 32], rax       # 8-byte Spill
	mov	r14d, r13d
	and	r14d, -32
	mov	rax, r13
	neg	rax
	mov	qword ptr [rsp + 40], rax       # 8-byte Spill
	lea	rax, [r8 + 96]
	mov	qword ptr [rsp + 8], rax        # 8-byte Spill
	lea	rax, [r8 + 4]
	mov	qword ptr [rsp + 24], rax       # 8-byte Spill
	xor	edi, edi
	movsxd	r15, r12d
	jmp	.LBB0_20
	.p2align	4, 0x90
.LBB0_19:                               #   in Loop: Header=BB0_20 Depth=1
	add	rdi, 1
	cmp	rdi, qword ptr [rsp + 48]       # 8-byte Folded Reload
	je	.LBB0_33
.LBB0_20:                               # =>This Loop Header: Depth=1
                                        #     Child Loop BB0_26 Depth 2
                                        #     Child Loop BB0_32 Depth 2
	movsx	r11, word ptr [rsi + 2*rdi]
	imul	r11, r15
	movsx	r10d, byte ptr [rdx + rdi]
	cmp	r12d, 32
	jb	.LBB0_21
# %bb.22:                               #   in Loop: Header=BB0_20 Depth=1
	mov	rax, qword ptr [rsp + 32]       # 8-byte Reload
	lea	rax, [rax + 4*r11]
	cmp	rax, r9
	jbe	.LBB0_25
# %bb.23:                               #   in Loop: Header=BB0_20 Depth=1
	lea	rax, [r8 + 4*r11]
	cmp	rax, qword ptr [rsp + 16]       # 8-byte Folded Reload
	jae	.LBB0_25
.LBB0_21:                               #   in Loop: Header=BB0_20 Depth=1
	xor	ebx, ebx
.LBB0_28:                               #   in Loop: Header=BB0_20 Depth=1
	mov	rax, rbx
	test	r13b, 1
	je	.LBB0_30
# %bb.29:                               #   in Loop: Header=BB0_20 Depth=1
	lea	rax, [rbx + r11]
	mov	eax, dword ptr [r8 + 4*rax]
	imul	eax, r10d
	add	dword ptr [r9 + 4*rbx], eax
	mov	rax, rbx
	or	rax, 1
.LBB0_30:                               #   in Loop: Header=BB0_20 Depth=1
	not	rbx
	cmp	rbx, qword ptr [rsp + 40]       # 8-byte Folded Reload
	je	.LBB0_19
# %bb.31:                               #   in Loop: Header=BB0_20 Depth=1
	mov	rcx, qword ptr [rsp + 24]       # 8-byte Reload
	lea	rbx, [rcx + 4*r11]
	.p2align	4, 0x90
.LBB0_32:                               #   Parent Loop BB0_20 Depth=1
                                        # =>  This Inner Loop Header: Depth=2
	mov	ecx, dword ptr [rbx + 4*rax - 4]
	imul	ecx, r10d
	add	dword ptr [r9 + 4*rax], ecx
	mov	ecx, dword ptr [rbx + 4*rax]
	imul	ecx, r10d
	add	dword ptr [r9 + 4*rax + 4], ecx
	add	rax, 2
	cmp	r13, rax
	jne	.LBB0_32
	jmp	.LBB0_19
	.p2align	4, 0x90
.LBB0_25:                               #   in Loop: Header=BB0_20 Depth=1
	mov	ecx, r12d
	mov	rax, r8
	vmovd	xmm0, r10d
	vpshufd	xmm0, xmm0, 0                   # xmm0 = xmm0[0,0,0,0]
	vinsertf128	ymm0, ymm0, xmm0, 1
	mov	rbx, qword ptr [rsp + 8]        # 8-byte Reload
	lea	r12, [rbx + 4*r11]
	xor	r8d, r8d
	.p2align	4, 0x90
.LBB0_26:                               #   Parent Loop BB0_20 Depth=1
                                        # =>  This Inner Loop Header: Depth=2
	vpmulld	xmm8, xmm0, xmmword ptr [r12 + 4*r8 - 96]
	vextractf128	xmm2, ymm0, 1
	vpmulld	xmm3, xmm2, xmmword ptr [r12 + 4*r8 - 80]
	vpmulld	xmm9, xmm0, xmmword ptr [r12 + 4*r8 - 64]
	vpmulld	xmm5, xmm2, xmmword ptr [r12 + 4*r8 - 48]
	vpmulld	xmm10, xmm0, xmmword ptr [r12 + 4*r8 - 32]
	vpmulld	xmm7, xmm2, xmmword ptr [r12 + 4*r8 - 16]
	vpmulld	xmm11, xmm0, xmmword ptr [r12 + 4*r8]
	vpmulld	xmm2, xmm2, xmmword ptr [r12 + 4*r8 + 16]
	vpaddd	xmm3, xmm3, xmmword ptr [r9 + 4*r8 + 16]
	vpaddd	xmm4, xmm8, xmmword ptr [r9 + 4*r8]
	vpaddd	xmm5, xmm5, xmmword ptr [r9 + 4*r8 + 48]
	vpaddd	xmm6, xmm9, xmmword ptr [r9 + 4*r8 + 32]
	vpaddd	xmm7, xmm7, xmmword ptr [r9 + 4*r8 + 80]
	vpaddd	xmm1, xmm10, xmmword ptr [r9 + 4*r8 + 64]
	vpaddd	xmm8, xmm2, xmmword ptr [r9 + 4*r8 + 112]
	vpaddd	xmm2, xmm11, xmmword ptr [r9 + 4*r8 + 96]
	vmovdqu	xmmword ptr [r9 + 4*r8], xmm4
	vmovdqu	xmmword ptr [r9 + 4*r8 + 16], xmm3
	vmovdqu	xmmword ptr [r9 + 4*r8 + 32], xmm6
	vmovdqu	xmmword ptr [r9 + 4*r8 + 48], xmm5
	vmovdqu	xmmword ptr [r9 + 4*r8 + 64], xmm1
	vmovdqu	xmmword ptr [r9 + 4*r8 + 80], xmm7
	vmovdqu	xmmword ptr [r9 + 4*r8 + 96], xmm2
	vmovdqu	xmmword ptr [r9 + 4*r8 + 112], xmm8
	add	r8, 32
	cmp	r14, r8
	jne	.LBB0_26
# %bb.27:                               #   in Loop: Header=BB0_20 Depth=1
	mov	rbx, r14
	cmp	r14, r13
	mov	r8, rax
	mov	r12d, ecx
	je	.LBB0_19
	jmp	.LBB0_28
.LBB0_33:
	lea	rsp, [rbp - 40]
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	vzeroupper
	ret
.LBB0_34:
	xor	ebx, ebx
	test	r10b, 1
	jne	.LBB0_10
	jmp	.LBB0_11
.Lfunc_end0:
	.size	update_hidden, .Lfunc_end0-update_hidden
                                        # -- End function
	.globl	quick_feed                      # -- Begin function quick_feed
	.p2align	4, 0x90
	.type	quick_feed,@function
quick_feed:                             # @quick_feed
# %bb.0:
	push	rbp
	mov	rbp, rsp
	and	rsp, -8
	test	ecx, ecx
	jle	.LBB1_1
# %bb.2:
	mov	r10d, ecx
	xor	r9d, r9d
	cmp	ecx, 32
	jae	.LBB1_4
# %bb.3:
	xor	ecx, ecx
	xor	esi, esi
	jmp	.LBB1_7
.LBB1_1:
	xor	esi, esi
	jmp	.LBB1_8
.LBB1_4:
	mov	ecx, r10d
	and	ecx, -32
	vpxor	xmm9, xmm9, xmm9
	xor	esi, esi
	vpxor	xmm8, xmm8, xmm8
	vpxor	xmm12, xmm12, xmm12
	vpxor	xmm10, xmm10, xmm10
	vpxor	xmm11, xmm11, xmm11
	.p2align	4, 0x90
.LBB1_5:                                # =>This Inner Loop Header: Depth=1
	vpmaxsd	xmm5, xmm8, xmmword ptr [rdi + 4*rsi + 16]
	vpmaxsd	xmm6, xmm8, xmmword ptr [rdi + 4*rsi]
	vpmaxsd	xmm7, xmm8, xmmword ptr [rdi + 4*rsi + 48]
	vpmaxsd	xmm4, xmm8, xmmword ptr [rdi + 4*rsi + 32]
	vpmaxsd	xmm0, xmm8, xmmword ptr [rdi + 4*rsi + 80]
	vpmaxsd	xmm1, xmm8, xmmword ptr [rdi + 4*rsi + 64]
	vpmaxsd	xmm2, xmm8, xmmword ptr [rdi + 4*rsi + 112]
	vpmaxsd	xmm3, xmm8, xmmword ptr [rdi + 4*rsi + 96]
	vpmulld	xmm6, xmm6, xmmword ptr [rdx + 4*rsi]
	vpmulld	xmm5, xmm5, xmmword ptr [rdx + 4*rsi + 16]
	vpmulld	xmm4, xmm4, xmmword ptr [rdx + 4*rsi + 32]
	vpmulld	xmm7, xmm7, xmmword ptr [rdx + 4*rsi + 48]
	vpmulld	xmm1, xmm1, xmmword ptr [rdx + 4*rsi + 64]
	vpmulld	xmm0, xmm0, xmmword ptr [rdx + 4*rsi + 80]
	vpmulld	xmm13, xmm3, xmmword ptr [rdx + 4*rsi + 96]
	vpmulld	xmm2, xmm2, xmmword ptr [rdx + 4*rsi + 112]
	vextractf128	xmm3, ymm9, 1
	vpaddd	xmm3, xmm5, xmm3
	vpaddd	xmm5, xmm9, xmm6
	vinsertf128	ymm9, ymm5, xmm3, 1
	vextractf128	xmm3, ymm12, 1
	vpaddd	xmm3, xmm7, xmm3
	vpaddd	xmm4, xmm12, xmm4
	vinsertf128	ymm12, ymm4, xmm3, 1
	vextractf128	xmm3, ymm10, 1
	vpaddd	xmm0, xmm0, xmm3
	vpaddd	xmm1, xmm10, xmm1
	vinsertf128	ymm10, ymm1, xmm0, 1
	vextractf128	xmm0, ymm11, 1
	vpaddd	xmm0, xmm2, xmm0
	vpaddd	xmm1, xmm13, xmm11
	vinsertf128	ymm11, ymm1, xmm0, 1
	add	rsi, 32
	cmp	rcx, rsi
	jne	.LBB1_5
# %bb.6:
	vextractf128	xmm0, ymm9, 1
	vextractf128	xmm1, ymm12, 1
	vpaddd	xmm0, xmm1, xmm0
	vpaddd	xmm1, xmm12, xmm9
	vextractf128	xmm2, ymm10, 1
	vextractf128	xmm3, ymm11, 1
	vpaddd	xmm2, xmm2, xmm3
	vpaddd	xmm0, xmm0, xmm2
	vpaddd	xmm2, xmm10, xmm11
	vpaddd	xmm1, xmm1, xmm2
	vpaddd	xmm0, xmm1, xmm0
	vpshufd	xmm1, xmm0, 238                 # xmm1 = xmm0[2,3,2,3]
	vpaddd	xmm0, xmm0, xmm1
	vpshufd	xmm1, xmm0, 85                  # xmm1 = xmm0[1,1,1,1]
	vpaddd	xmm0, xmm0, xmm1
	vmovd	esi, xmm0
	cmp	rcx, r10
	je	.LBB1_8
	.p2align	4, 0x90
.LBB1_7:                                # =>This Inner Loop Header: Depth=1
	mov	eax, dword ptr [rdi + 4*rcx]
	test	eax, eax
	cmovs	eax, r9d
	imul	eax, dword ptr [rdx + 4*rcx]
	add	esi, eax
	add	rcx, 1
	cmp	r10, rcx
	jne	.LBB1_7
.LBB1_8:
	mov	dword ptr [r8], esi
	mov	rsp, rbp
	pop	rbp
	vzeroupper
	ret
.Lfunc_end1:
	.size	quick_feed, .Lfunc_end1-quick_feed
                                        # -- End function
	.ident	"Ubuntu clang version 12.0.0-3ubuntu1~21.04.2"
	.section	".note.GNU-stack","",@progbits
	.addrsig
