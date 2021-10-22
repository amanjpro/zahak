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
	sub	rsp, 88
	mov	qword ptr [rsp + 80], rdx       # 8-byte Spill
	cmp	dword ptr [rbp + 48], 0
	jle	.LBB0_37
# %bb.1:
	mov	r13, rcx
	mov	r11, qword ptr [rbp + 40]
	mov	rax, qword ptr [rbp + 32]
	mov	ecx, dword ptr [rbp + 48]
	mov	ebx, ecx
	cmp	ecx, 16
	mov	qword ptr [rsp + 32], r9        # 8-byte Spill
	mov	qword ptr [rsp + 24], r8        # 8-byte Spill
	mov	qword ptr [rsp + 16], r13       # 8-byte Spill
	jae	.LBB0_3
# %bb.2:
	xor	r14d, r14d
.LBB0_17:
	mov	rcx, r14
	not	rcx
	test	bl, 1
	je	.LBB0_19
# %bb.18:
	vmovss	xmm0, dword ptr [rdi + 4*r14]   # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [rax + 4*r14], xmm0
	vmovss	xmm0, dword ptr [rsi + 4*r14]   # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [r11 + 4*r14], xmm0
	or	r14, 1
.LBB0_19:
	add	rcx, rbx
	je	.LBB0_20
	.p2align	4, 0x90
.LBB0_38:                               # =>This Inner Loop Header: Depth=1
	vmovss	xmm0, dword ptr [rdi + 4*r14]   # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [rax + 4*r14], xmm0
	vmovss	xmm0, dword ptr [rsi + 4*r14]   # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [r11 + 4*r14], xmm0
	vmovss	xmm0, dword ptr [rdi + 4*r14 + 4] # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [rax + 4*r14 + 4], xmm0
	vmovss	xmm0, dword ptr [rsi + 4*r14 + 4] # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [r11 + 4*r14 + 4], xmm0
	add	r14, 2
	cmp	rbx, r14
	jne	.LBB0_38
	jmp	.LBB0_20
.LBB0_3:
	lea	r14, [r11 + 4*rbx]
	cmp	r14, rax
	seta	byte ptr [rsp + 40]             # 1-byte Folded Spill
	lea	rcx, [rax + 4*rbx]
	cmp	rcx, r11
	seta	r15b
	lea	rdx, [rdi + 4*rbx]
	cmp	rdx, rax
	seta	r12b
	cmp	rcx, rdi
	seta	r10b
	lea	r9, [rsi + 4*rbx]
	cmp	r9, rax
	seta	r13b
	cmp	rcx, rsi
	seta	byte ptr [rsp + 14]             # 1-byte Folded Spill
	cmp	rdx, r11
	seta	r8b
	cmp	r14, rdi
	seta	byte ptr [rsp + 13]             # 1-byte Folded Spill
	cmp	r9, r11
	seta	dl
	cmp	r14, rsi
	seta	cl
	xor	r14d, r14d
	test	byte ptr [rsp + 40], r15b       # 1-byte Folded Reload
	jne	.LBB0_4
# %bb.5:
	and	r12b, r10b
	mov	r9, qword ptr [rsp + 32]        # 8-byte Reload
	jne	.LBB0_6
# %bb.7:
	and	r13b, byte ptr [rsp + 14]       # 1-byte Folded Reload
	jne	.LBB0_6
# %bb.8:
	and	r8b, byte ptr [rsp + 13]        # 1-byte Folded Reload
	mov	r8, qword ptr [rsp + 24]        # 8-byte Reload
	mov	r13, qword ptr [rsp + 16]       # 8-byte Reload
	jne	.LBB0_17
# %bb.9:
	and	dl, cl
	jne	.LBB0_17
# %bb.10:
	mov	r14d, ebx
	and	r14d, -16
	lea	rcx, [r14 - 16]
	mov	r10, rcx
	shr	r10, 4
	add	r10, 1
	test	rcx, rcx
	je	.LBB0_11
# %bb.12:
	mov	rcx, r10
	and	rcx, -2
	neg	rcx
	xor	edx, edx
	.p2align	4, 0x90
.LBB0_13:                               # =>This Inner Loop Header: Depth=1
	vmovups	ymm0, ymmword ptr [rdi + 4*rdx]
	vmovups	ymm1, ymmword ptr [rdi + 4*rdx + 32]
	vmovups	ymmword ptr [rax + 4*rdx], ymm0
	vmovups	ymmword ptr [rax + 4*rdx + 32], ymm1
	vmovups	ymm0, ymmword ptr [rsi + 4*rdx]
	vmovups	ymm1, ymmword ptr [rsi + 4*rdx + 32]
	vmovups	ymmword ptr [r11 + 4*rdx], ymm0
	vmovups	ymmword ptr [r11 + 4*rdx + 32], ymm1
	vmovups	ymm0, ymmword ptr [rdi + 4*rdx + 64]
	vmovups	ymm1, ymmword ptr [rdi + 4*rdx + 96]
	vmovups	ymmword ptr [rax + 4*rdx + 64], ymm0
	vmovups	ymmword ptr [rax + 4*rdx + 96], ymm1
	vmovups	ymm0, ymmword ptr [rsi + 4*rdx + 64]
	vmovups	ymm1, ymmword ptr [rsi + 4*rdx + 96]
	vmovups	ymmword ptr [r11 + 4*rdx + 64], ymm0
	vmovups	ymmword ptr [r11 + 4*rdx + 96], ymm1
	add	rdx, 32
	add	rcx, 2
	jne	.LBB0_13
# %bb.14:
	test	r10b, 1
	je	.LBB0_16
.LBB0_15:
	vmovups	ymm0, ymmword ptr [rdi + 4*rdx]
	vmovups	ymm1, ymmword ptr [rdi + 4*rdx + 32]
	vmovups	ymmword ptr [rax + 4*rdx], ymm0
	vmovups	ymmword ptr [rax + 4*rdx + 32], ymm1
	vmovups	ymm0, ymmword ptr [rsi + 4*rdx]
	vmovups	ymm1, ymmword ptr [rsi + 4*rdx + 32]
	vmovups	ymmword ptr [r11 + 4*rdx], ymm0
	vmovups	ymmword ptr [r11 + 4*rdx + 32], ymm1
.LBB0_16:
	cmp	r14, rbx
	jne	.LBB0_17
.LBB0_20:
	cmp	dword ptr [rbp + 48], 0
	jle	.LBB0_37
# %bb.21:
	mov	ecx, dword ptr [rbp + 16]
	test	ecx, ecx
	jle	.LBB0_37
# %bb.22:
	mov	r10, qword ptr [rbp + 24]
	mov	ecx, ecx
	mov	qword ptr [rsp + 40], rcx       # 8-byte Spill
	lea	rsi, [rax + 4*rbx]
	lea	rcx, [r11 + 4*rbx]
	mov	qword ptr [rsp + 64], rcx       # 8-byte Spill
	cmp	rcx, rax
	seta	cl
	lea	rdi, [r10 + 4*rbx]
	mov	qword ptr [rsp + 56], rdi       # 8-byte Spill
	mov	qword ptr [rsp + 72], rsi       # 8-byte Spill
	cmp	rsi, r11
	seta	dl
	and	dl, cl
	mov	byte ptr [rsp + 15], dl         # 1-byte Spill
	mov	r15d, ebx
	and	r15d, -16
	lea	rcx, [r10 + 32]
	mov	qword ptr [rsp + 48], rcx       # 8-byte Spill
	xor	r14d, r14d
	jmp	.LBB0_23
	.p2align	4, 0x90
.LBB0_36:                               #   in Loop: Header=BB0_23 Depth=1
	add	r14, 1
	cmp	r14, qword ptr [rsp + 40]       # 8-byte Folded Reload
	je	.LBB0_37
.LBB0_23:                               # =>This Loop Header: Depth=1
                                        #     Child Loop BB0_32 Depth 2
                                        #     Child Loop BB0_35 Depth 2
	mov	rcx, qword ptr [rsp + 80]       # 8-byte Reload
	movsx	r12, word ptr [rcx + 2*r14]
	movsx	rdx, word ptr [r8 + 2*r14]
	mov	esi, dword ptr [rbp + 48]
	movsxd	rcx, esi
	imul	r12, rcx
	imul	rdx, rcx
	movsx	ecx, byte ptr [r13 + r14]
	vcvtsi2ss	xmm0, xmm6, ecx
	movsx	ecx, byte ptr [r9 + r14]
	vcvtsi2ss	xmm1, xmm6, ecx
	cmp	esi, 16
	jae	.LBB0_25
# %bb.24:                               #   in Loop: Header=BB0_23 Depth=1
	xor	r10d, r10d
.LBB0_34:                               #   in Loop: Header=BB0_23 Depth=1
	mov	rsi, qword ptr [rbp + 24]
	lea	rcx, [rsi + 4*rdx]
	lea	rdx, [rsi + 4*r12]
	.p2align	4, 0x90
.LBB0_35:                               #   Parent Loop BB0_23 Depth=1
                                        # =>  This Inner Loop Header: Depth=2
	vmulss	xmm2, xmm0, dword ptr [rdx + 4*r10]
	vaddss	xmm2, xmm2, dword ptr [rax + 4*r10]
	vmovss	dword ptr [rax + 4*r10], xmm2
	vmulss	xmm2, xmm1, dword ptr [rcx + 4*r10]
	vaddss	xmm2, xmm2, dword ptr [r11 + 4*r10]
	vmovss	dword ptr [r11 + 4*r10], xmm2
	add	r10, 1
	cmp	rbx, r10
	jne	.LBB0_35
	jmp	.LBB0_36
	.p2align	4, 0x90
.LBB0_25:                               #   in Loop: Header=BB0_23 Depth=1
	mov	r8, qword ptr [rbp + 24]
	lea	rsi, [r8 + 4*rdx]
	mov	rcx, qword ptr [rsp + 56]       # 8-byte Reload
	lea	rdi, [rcx + 4*rdx]
	lea	r10, [r8 + 4*r12]
	lea	r9, [rcx + 4*r12]
	cmp	rdi, rax
	seta	r8b
	mov	rcx, qword ptr [rsp + 72]       # 8-byte Reload
	cmp	rsi, rcx
	setb	r13b
	and	r13b, r8b
	or	r13b, byte ptr [rsp + 15]       # 1-byte Folded Reload
	cmp	r9, rax
	seta	r8b
	cmp	r10, rcx
	setb	byte ptr [rsp + 14]             # 1-byte Folded Spill
	cmp	rdi, r11
	seta	dil
	mov	rcx, qword ptr [rsp + 64]       # 8-byte Reload
	cmp	rsi, rcx
	setb	byte ptr [rsp + 13]             # 1-byte Folded Spill
	cmp	r9, r11
	seta	sil
	cmp	r10, rcx
	setb	cl
	test	r13b, r13b
	jne	.LBB0_26
# %bb.27:                               #   in Loop: Header=BB0_23 Depth=1
	and	r8b, byte ptr [rsp + 14]        # 1-byte Folded Reload
	jne	.LBB0_26
# %bb.28:                               #   in Loop: Header=BB0_23 Depth=1
	and	dil, byte ptr [rsp + 13]        # 1-byte Folded Reload
	mov	r8, qword ptr [rsp + 24]        # 8-byte Reload
	mov	r13, qword ptr [rsp + 16]       # 8-byte Reload
	jne	.LBB0_29
# %bb.30:                               #   in Loop: Header=BB0_23 Depth=1
	mov	r10d, 0
	and	sil, cl
	mov	r9, qword ptr [rsp + 32]        # 8-byte Reload
	jne	.LBB0_34
# %bb.31:                               #   in Loop: Header=BB0_23 Depth=1
	vpermilps	xmm2, xmm0, 0           # xmm2 = xmm0[0,0,0,0]
	vinsertf128	ymm2, ymm2, xmm2, 1
	vpermilps	xmm3, xmm1, 0           # xmm3 = xmm1[0,0,0,0]
	vinsertf128	ymm3, ymm3, xmm3, 1
	mov	rcx, qword ptr [rsp + 48]       # 8-byte Reload
	lea	rdi, [rcx + 4*r12]
	lea	rsi, [rcx + 4*rdx]
	xor	ecx, ecx
	.p2align	4, 0x90
.LBB0_32:                               #   Parent Loop BB0_23 Depth=1
                                        # =>  This Inner Loop Header: Depth=2
	vmulps	ymm4, ymm2, ymmword ptr [rdi + 4*rcx - 32]
	vmulps	ymm5, ymm2, ymmword ptr [rdi + 4*rcx]
	vaddps	ymm4, ymm4, ymmword ptr [rax + 4*rcx]
	vaddps	ymm5, ymm5, ymmword ptr [rax + 4*rcx + 32]
	vmovups	ymmword ptr [rax + 4*rcx], ymm4
	vmovups	ymmword ptr [rax + 4*rcx + 32], ymm5
	vmulps	ymm4, ymm3, ymmword ptr [rsi + 4*rcx - 32]
	vmulps	ymm5, ymm3, ymmword ptr [rsi + 4*rcx]
	vaddps	ymm4, ymm4, ymmword ptr [r11 + 4*rcx]
	vaddps	ymm5, ymm5, ymmword ptr [r11 + 4*rcx + 32]
	vmovups	ymmword ptr [r11 + 4*rcx], ymm4
	vmovups	ymmword ptr [r11 + 4*rcx + 32], ymm5
	add	rcx, 16
	cmp	r15, rcx
	jne	.LBB0_32
# %bb.33:                               #   in Loop: Header=BB0_23 Depth=1
	mov	r10, r15
	cmp	r15, rbx
	je	.LBB0_36
	jmp	.LBB0_34
.LBB0_26:                               #   in Loop: Header=BB0_23 Depth=1
	xor	r10d, r10d
	mov	r9, qword ptr [rsp + 32]        # 8-byte Reload
	mov	r8, qword ptr [rsp + 24]        # 8-byte Reload
	mov	r13, qword ptr [rsp + 16]       # 8-byte Reload
	jmp	.LBB0_34
.LBB0_29:                               #   in Loop: Header=BB0_23 Depth=1
	xor	r10d, r10d
	mov	r9, qword ptr [rsp + 32]        # 8-byte Reload
	jmp	.LBB0_34
.LBB0_37:
	lea	rsp, [rbp - 40]
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	vzeroupper
	ret
.LBB0_11:
	xor	edx, edx
	test	r10b, 1
	jne	.LBB0_15
	jmp	.LBB0_16
.LBB0_4:
	mov	r9, qword ptr [rsp + 32]        # 8-byte Reload
.LBB0_6:
	mov	r8, qword ptr [rsp + 24]        # 8-byte Reload
	mov	r13, qword ptr [rsp + 16]       # 8-byte Reload
	jmp	.LBB0_17
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
	mov	eax, ecx
	vxorps	xmm0, xmm0, xmm0
	cmp	ecx, 32
	jae	.LBB1_4
# %bb.3:
	xor	ecx, ecx
	vxorps	xmm1, xmm1, xmm1
	jmp	.LBB1_7
.LBB1_1:
	vxorps	xmm1, xmm1, xmm1
	jmp	.LBB1_8
.LBB1_4:
	mov	ecx, eax
	and	ecx, -32
	vxorps	xmm1, xmm1, xmm1
	xor	esi, esi
	vxorps	xmm2, xmm2, xmm2
	vxorps	xmm3, xmm3, xmm3
	vxorps	xmm4, xmm4, xmm4
	vxorps	xmm5, xmm5, xmm5
	.p2align	4, 0x90
.LBB1_5:                                # =>This Inner Loop Header: Depth=1
	vmaxps	ymm6, ymm1, ymmword ptr [rdi + 4*rsi]
	vmaxps	ymm7, ymm1, ymmword ptr [rdi + 4*rsi + 32]
	vmaxps	ymm8, ymm1, ymmword ptr [rdi + 4*rsi + 64]
	vmaxps	ymm9, ymm1, ymmword ptr [rdi + 4*rsi + 96]
	vmulps	ymm6, ymm6, ymmword ptr [rdx + 4*rsi]
	vaddps	ymm2, ymm6, ymm2
	vmulps	ymm6, ymm7, ymmword ptr [rdx + 4*rsi + 32]
	vaddps	ymm3, ymm6, ymm3
	vmulps	ymm6, ymm8, ymmword ptr [rdx + 4*rsi + 64]
	vmulps	ymm7, ymm9, ymmword ptr [rdx + 4*rsi + 96]
	vaddps	ymm4, ymm6, ymm4
	vaddps	ymm5, ymm7, ymm5
	add	rsi, 32
	cmp	rcx, rsi
	jne	.LBB1_5
# %bb.6:
	vaddps	ymm1, ymm3, ymm2
	vaddps	ymm1, ymm4, ymm1
	vaddps	ymm1, ymm5, ymm1
	vextractf128	xmm2, ymm1, 1
	vaddps	xmm1, xmm1, xmm2
	vpermilpd	xmm2, xmm1, 1           # xmm2 = xmm1[1,0]
	vaddps	xmm1, xmm1, xmm2
	vmovshdup	xmm2, xmm1              # xmm2 = xmm1[1,1,3,3]
	vaddss	xmm1, xmm1, xmm2
	cmp	rcx, rax
	je	.LBB1_8
	.p2align	4, 0x90
.LBB1_7:                                # =>This Inner Loop Header: Depth=1
	vmaxss	xmm2, xmm0, dword ptr [rdi + 4*rcx]
	vmulss	xmm2, xmm2, dword ptr [rdx + 4*rcx]
	vaddss	xmm1, xmm2, xmm1
	add	rcx, 1
	cmp	rax, rcx
	jne	.LBB1_7
.LBB1_8:
	vmovss	dword ptr [r8], xmm1
	mov	rsp, rbp
	pop	rbp
	vzeroupper
	ret
.Lfunc_end1:
	.size	quick_feed, .Lfunc_end1-quick_feed
                                        # -- End function
	.ident	"clang version 12.0.1"
	.section	".note.GNU-stack","",@progbits
	.addrsig
