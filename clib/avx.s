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
	mov	qword ptr [rsp + 40], rdx       # 8-byte Spill
	mov	qword ptr [rsp + 48], rsi       # 8-byte Spill
	cmp	dword ptr [rbp + 16], 0
	jle	.LBB0_45
# %bb.1:
	mov	r15, r8
	mov	r11, rdi
	mov	eax, dword ptr [rbp + 16]
	mov	edi, eax
	cmp	eax, 8
	jb	.LBB0_2
# %bb.3:
	lea	rax, [r11 + 2*rdi]
	cmp	rax, r9
	jbe	.LBB0_6
# %bb.4:
	lea	rax, [r9 + 2*rdi]
	cmp	rax, r11
	jbe	.LBB0_6
.LBB0_2:
	xor	eax, eax
.LBB0_18:
	mov	rdx, rax
	not	rdx
	add	rdx, rdi
	mov	rsi, rdi
	and	rsi, 3
	je	.LBB0_20
	.p2align	4, 0x90
.LBB0_19:                               # =>This Inner Loop Header: Depth=1
	movzx	ebx, word ptr [r11 + 2*rax]
	mov	word ptr [r9 + 2*rax], bx
	add	rax, 1
	add	rsi, -1
	jne	.LBB0_19
.LBB0_20:
	cmp	rdx, 3
	jb	.LBB0_22
	.p2align	4, 0x90
.LBB0_21:                               # =>This Inner Loop Header: Depth=1
	movzx	edx, word ptr [r11 + 2*rax]
	mov	word ptr [r9 + 2*rax], dx
	movzx	edx, word ptr [r11 + 2*rax + 2]
	mov	word ptr [r9 + 2*rax + 2], dx
	movzx	edx, word ptr [r11 + 2*rax + 4]
	mov	word ptr [r9 + 2*rax + 4], dx
	movzx	edx, word ptr [r11 + 2*rax + 6]
	mov	word ptr [r9 + 2*rax + 6], dx
	add	rax, 4
	cmp	rdi, rax
	jne	.LBB0_21
.LBB0_22:
	cmp	dword ptr [rbp + 16], 0
	jle	.LBB0_45
# %bb.23:
	test	ecx, ecx
	jle	.LBB0_45
# %bb.24:
	mov	eax, ecx
	mov	qword ptr [rsp + 32], rax       # 8-byte Spill
	lea	rax, [r9 + 2*rdi]
	mov	qword ptr [rsp + 8], rax        # 8-byte Spill
	lea	rax, [r15 + 2*rdi]
	mov	qword ptr [rsp + 24], rax       # 8-byte Spill
	mov	eax, edi
	and	eax, -64
	mov	esi, edi
	and	esi, -8
	mov	r8, rdi
	neg	r8
	lea	rcx, [r15 + 96]
	mov	qword ptr [rsp], rcx            # 8-byte Spill
	lea	rcx, [r15 + 2]
	mov	qword ptr [rsp + 16], rcx       # 8-byte Spill
	xor	r12d, r12d
	movsxd	r11, dword ptr [rbp + 16]
	jmp	.LBB0_26
	.p2align	4, 0x90
.LBB0_25:                               #   in Loop: Header=BB0_26 Depth=1
	add	r12, 1
	cmp	r12, qword ptr [rsp + 32]       # 8-byte Folded Reload
	je	.LBB0_45
.LBB0_26:                               # =>This Loop Header: Depth=1
                                        #     Child Loop BB0_34 Depth 2
                                        #     Child Loop BB0_38 Depth 2
                                        #     Child Loop BB0_44 Depth 2
	mov	rcx, qword ptr [rsp + 48]       # 8-byte Reload
	movsx	r13, word ptr [rcx + 2*r12]
	imul	r13, r11
	mov	rcx, qword ptr [rsp + 40]       # 8-byte Reload
	movsx	r10d, byte ptr [rcx + r12]
	cmp	dword ptr [rbp + 16], 8
	jb	.LBB0_27
# %bb.28:                               #   in Loop: Header=BB0_26 Depth=1
	lea	r14, [r15 + 2*r13]
	mov	rcx, qword ptr [rsp + 24]       # 8-byte Reload
	lea	rdx, [rcx + 2*r13]
	cmp	rdx, r9
	jbe	.LBB0_31
# %bb.29:                               #   in Loop: Header=BB0_26 Depth=1
	cmp	r14, qword ptr [rsp + 8]        # 8-byte Folded Reload
	jae	.LBB0_31
.LBB0_27:                               #   in Loop: Header=BB0_26 Depth=1
	xor	ebx, ebx
.LBB0_40:                               #   in Loop: Header=BB0_26 Depth=1
	mov	rdx, rbx
	test	dil, 1
	je	.LBB0_42
# %bb.41:                               #   in Loop: Header=BB0_26 Depth=1
	lea	rdx, [rbx + r13]
	movzx	edx, word ptr [r15 + 2*rdx]
	imul	dx, r10w
	add	word ptr [r9 + 2*rbx], dx
	mov	rdx, rbx
	or	rdx, 1
.LBB0_42:                               #   in Loop: Header=BB0_26 Depth=1
	not	rbx
	cmp	rbx, r8
	je	.LBB0_25
# %bb.43:                               #   in Loop: Header=BB0_26 Depth=1
	mov	rcx, qword ptr [rsp + 16]       # 8-byte Reload
	lea	rbx, [rcx + 2*r13]
	.p2align	4, 0x90
.LBB0_44:                               #   Parent Loop BB0_26 Depth=1
                                        # =>  This Inner Loop Header: Depth=2
	movzx	ecx, word ptr [rbx + 2*rdx - 2]
	imul	cx, r10w
	add	word ptr [r9 + 2*rdx], cx
	movzx	ecx, word ptr [rbx + 2*rdx]
	imul	cx, r10w
	add	word ptr [r9 + 2*rdx + 2], cx
	add	rdx, 2
	cmp	rdi, rdx
	jne	.LBB0_44
	jmp	.LBB0_25
	.p2align	4, 0x90
.LBB0_31:                               #   in Loop: Header=BB0_26 Depth=1
	cmp	dword ptr [rbp + 16], 64
	jae	.LBB0_33
# %bb.32:                               #   in Loop: Header=BB0_26 Depth=1
	xor	edx, edx
	jmp	.LBB0_37
.LBB0_33:                               #   in Loop: Header=BB0_26 Depth=1
	mov	rbx, r15
	vmovd	xmm0, r10d
	vpshuflw	xmm0, xmm0, 0                   # xmm0 = xmm0[0,0,0,0,4,5,6,7]
	vpshufd	xmm0, xmm0, 0                   # xmm0 = xmm0[0,0,0,0]
	vinsertf128	ymm0, ymm0, xmm0, 1
	mov	rcx, qword ptr [rsp]            # 8-byte Reload
	lea	r15, [rcx + 2*r13]
	xor	edx, edx
	.p2align	4, 0x90
.LBB0_34:                               #   Parent Loop BB0_26 Depth=1
                                        # =>  This Inner Loop Header: Depth=2
	vpmullw	xmm8, xmm0, xmmword ptr [r15 + 2*rdx - 96]
	vextractf128	xmm2, ymm0, 1
	vpmullw	xmm3, xmm2, xmmword ptr [r15 + 2*rdx - 80]
	vpmullw	xmm9, xmm0, xmmword ptr [r15 + 2*rdx - 64]
	vpmullw	xmm5, xmm2, xmmword ptr [r15 + 2*rdx - 48]
	vpmullw	xmm10, xmm0, xmmword ptr [r15 + 2*rdx - 32]
	vpmullw	xmm7, xmm2, xmmword ptr [r15 + 2*rdx - 16]
	vpmullw	xmm11, xmm0, xmmword ptr [r15 + 2*rdx]
	vpmullw	xmm2, xmm2, xmmword ptr [r15 + 2*rdx + 16]
	vpaddw	xmm3, xmm3, xmmword ptr [r9 + 2*rdx + 16]
	vpaddw	xmm4, xmm8, xmmword ptr [r9 + 2*rdx]
	vpaddw	xmm5, xmm5, xmmword ptr [r9 + 2*rdx + 48]
	vpaddw	xmm6, xmm9, xmmword ptr [r9 + 2*rdx + 32]
	vpaddw	xmm7, xmm7, xmmword ptr [r9 + 2*rdx + 80]
	vpaddw	xmm1, xmm10, xmmword ptr [r9 + 2*rdx + 64]
	vpaddw	xmm8, xmm2, xmmword ptr [r9 + 2*rdx + 112]
	vpaddw	xmm2, xmm11, xmmword ptr [r9 + 2*rdx + 96]
	vmovdqu	xmmword ptr [r9 + 2*rdx], xmm4
	vmovdqu	xmmword ptr [r9 + 2*rdx + 16], xmm3
	vmovdqu	xmmword ptr [r9 + 2*rdx + 32], xmm6
	vmovdqu	xmmword ptr [r9 + 2*rdx + 48], xmm5
	vmovdqu	xmmword ptr [r9 + 2*rdx + 64], xmm1
	vmovdqu	xmmword ptr [r9 + 2*rdx + 80], xmm7
	vmovdqu	xmmword ptr [r9 + 2*rdx + 96], xmm2
	vmovdqu	xmmword ptr [r9 + 2*rdx + 112], xmm8
	add	rdx, 64
	cmp	rax, rdx
	jne	.LBB0_34
# %bb.35:                               #   in Loop: Header=BB0_26 Depth=1
	cmp	rax, rdi
	mov	r15, rbx
	je	.LBB0_25
# %bb.36:                               #   in Loop: Header=BB0_26 Depth=1
	mov	rdx, rax
	mov	rbx, rax
	test	dil, 56
	je	.LBB0_40
.LBB0_37:                               #   in Loop: Header=BB0_26 Depth=1
	vmovd	xmm0, r10d
	vpshuflw	xmm0, xmm0, 0                   # xmm0 = xmm0[0,0,0,0,4,5,6,7]
	vpshufd	xmm0, xmm0, 0                   # xmm0 = xmm0[0,0,0,0]
	.p2align	4, 0x90
.LBB0_38:                               #   Parent Loop BB0_26 Depth=1
                                        # =>  This Inner Loop Header: Depth=2
	vpmullw	xmm1, xmm0, xmmword ptr [r14 + 2*rdx]
	vpaddw	xmm1, xmm1, xmmword ptr [r9 + 2*rdx]
	vmovdqu	xmmword ptr [r9 + 2*rdx], xmm1
	add	rdx, 8
	cmp	rsi, rdx
	jne	.LBB0_38
# %bb.39:                               #   in Loop: Header=BB0_26 Depth=1
	mov	rbx, rsi
	cmp	rsi, rdi
	je	.LBB0_25
	jmp	.LBB0_40
.LBB0_45:
	lea	rsp, [rbp - 40]
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	vzeroupper
	ret
.LBB0_6:
	cmp	dword ptr [rbp + 16], 64
	jae	.LBB0_8
# %bb.7:
	xor	eax, eax
	jmp	.LBB0_15
.LBB0_8:
	mov	eax, edi
	and	eax, -64
	lea	rdx, [rax - 64]
	mov	r10, rdx
	shr	r10, 6
	add	r10, 1
	test	rdx, rdx
	je	.LBB0_46
# %bb.9:
	mov	rdx, r10
	and	rdx, -2
	neg	rdx
	xor	ebx, ebx
	.p2align	4, 0x90
.LBB0_10:                               # =>This Inner Loop Header: Depth=1
	vmovups	ymm0, ymmword ptr [r11 + 2*rbx]
	vmovups	ymm1, ymmword ptr [r11 + 2*rbx + 32]
	vmovups	ymm2, ymmword ptr [r11 + 2*rbx + 64]
	vmovups	ymm3, ymmword ptr [r11 + 2*rbx + 96]
	vmovups	ymmword ptr [r9 + 2*rbx], ymm0
	vmovups	ymmword ptr [r9 + 2*rbx + 32], ymm1
	vmovups	ymmword ptr [r9 + 2*rbx + 64], ymm2
	vmovups	ymmword ptr [r9 + 2*rbx + 96], ymm3
	vmovdqu	ymm0, ymmword ptr [r11 + 2*rbx + 128]
	vmovdqu	ymm1, ymmword ptr [r11 + 2*rbx + 160]
	vmovdqu	ymm2, ymmword ptr [r11 + 2*rbx + 192]
	vmovdqu	ymm3, ymmword ptr [r11 + 2*rbx + 224]
	vmovdqu	ymmword ptr [r9 + 2*rbx + 128], ymm0
	vmovdqu	ymmword ptr [r9 + 2*rbx + 160], ymm1
	vmovdqu	ymmword ptr [r9 + 2*rbx + 192], ymm2
	vmovdqu	ymmword ptr [r9 + 2*rbx + 224], ymm3
	sub	rbx, -128
	add	rdx, 2
	jne	.LBB0_10
# %bb.11:
	test	r10b, 1
	je	.LBB0_13
.LBB0_12:
	vmovdqu	ymm0, ymmword ptr [r11 + 2*rbx]
	vmovdqu	ymm1, ymmword ptr [r11 + 2*rbx + 32]
	vmovdqu	ymm2, ymmword ptr [r11 + 2*rbx + 64]
	vmovdqu	ymm3, ymmword ptr [r11 + 2*rbx + 96]
	vmovdqu	ymmword ptr [r9 + 2*rbx], ymm0
	vmovdqu	ymmword ptr [r9 + 2*rbx + 32], ymm1
	vmovdqu	ymmword ptr [r9 + 2*rbx + 64], ymm2
	vmovdqu	ymmword ptr [r9 + 2*rbx + 96], ymm3
.LBB0_13:
	cmp	rax, rdi
	je	.LBB0_22
# %bb.14:
	test	dil, 56
	je	.LBB0_18
.LBB0_15:
	mov	rdx, rax
	mov	eax, edi
	and	eax, -8
	.p2align	4, 0x90
.LBB0_16:                               # =>This Inner Loop Header: Depth=1
	vmovdqu	xmm0, xmmword ptr [r11 + 2*rdx]
	vmovdqu	xmmword ptr [r9 + 2*rdx], xmm0
	add	rdx, 8
	cmp	rax, rdx
	jne	.LBB0_16
# %bb.17:
	cmp	rax, rdi
	je	.LBB0_22
	jmp	.LBB0_18
.LBB0_46:
	xor	ebx, ebx
	test	r10b, 1
	jne	.LBB0_12
	jmp	.LBB0_13
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
	vpxor	xmm8, xmm8, xmm8
	xor	esi, esi
	vpxor	xmm1, xmm1, xmm1
	vpxor	xmm2, xmm2, xmm2
	vpxor	xmm3, xmm3, xmm3
	vpxor	xmm4, xmm4, xmm4
	.p2align	4, 0x90
.LBB1_5:                                # =>This Inner Loop Header: Depth=1
	vpmaxsw	xmm5, xmm8, xmmword ptr [rdi + 2*rsi]
	vpmaxsw	xmm6, xmm8, xmmword ptr [rdi + 2*rsi + 16]
	vpmaxsw	xmm7, xmm8, xmmword ptr [rdi + 2*rsi + 32]
	vpmaxsw	xmm0, xmm8, xmmword ptr [rdi + 2*rsi + 48]
	vpmullw	xmm5, xmm5, xmmword ptr [rdx + 2*rsi]
	vpaddw	xmm1, xmm5, xmm1
	vpmullw	xmm5, xmm6, xmmword ptr [rdx + 2*rsi + 16]
	vpaddw	xmm2, xmm5, xmm2
	vpmullw	xmm5, xmm7, xmmword ptr [rdx + 2*rsi + 32]
	vpmullw	xmm0, xmm0, xmmword ptr [rdx + 2*rsi + 48]
	vpaddw	xmm3, xmm5, xmm3
	vpaddw	xmm4, xmm0, xmm4
	add	rsi, 32
	cmp	rcx, rsi
	jne	.LBB1_5
# %bb.6:
	vpaddw	xmm0, xmm2, xmm1
	vpaddw	xmm0, xmm3, xmm0
	vpaddw	xmm0, xmm4, xmm0
	vpshufd	xmm1, xmm0, 238                 # xmm1 = xmm0[2,3,2,3]
	vpaddw	xmm0, xmm0, xmm1
	vpshufd	xmm1, xmm0, 85                  # xmm1 = xmm0[1,1,1,1]
	vpaddw	xmm0, xmm0, xmm1
	vpsrld	xmm1, xmm0, 16
	vpaddw	xmm0, xmm0, xmm1
	vmovd	esi, xmm0
	cmp	rcx, r10
	je	.LBB1_8
	.p2align	4, 0x90
.LBB1_7:                                # =>This Inner Loop Header: Depth=1
	movzx	eax, word ptr [rdi + 2*rcx]
	test	ax, ax
	cmovs	eax, r9d
	imul	ax, word ptr [rdx + 2*rcx]
	add	esi, eax
	add	rcx, 1
	cmp	r10, rcx
	jne	.LBB1_7
.LBB1_8:
	mov	word ptr [r8], si
	mov	rsp, rbp
	pop	rbp
	ret
.Lfunc_end1:
	.size	quick_feed, .Lfunc_end1-quick_feed
                                        # -- End function
	.ident	"Ubuntu clang version 12.0.0-3ubuntu1~21.04.2"
	.section	".note.GNU-stack","",@progbits
	.addrsig
