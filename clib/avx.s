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
	sub	rsp, 40
	mov	qword ptr [rsp], r8             # 8-byte Spill
	mov	qword ptr [rsp + 32], rdx       # 8-byte Spill
	mov	eax, dword ptr [rbp + 16]
	test	eax, eax
	jle	.LBB0_33
# %bb.1:
	mov	r11d, ecx
	mov	r12d, eax
	cmp	eax, 32
	jb	.LBB0_2
# %bb.3:
	lea	rax, [rdi + 4*r12]
	cmp	rax, r9
	jbe	.LBB0_6
# %bb.4:
	lea	rax, [r9 + 4*r12]
	cmp	rax, rdi
	jbe	.LBB0_6
.LBB0_2:
	xor	eax, eax
.LBB0_12:
	mov	rbx, rax
	not	rbx
	add	rbx, r12
	mov	rcx, r12
	and	rcx, 3
	je	.LBB0_14
	.p2align	4, 0x90
.LBB0_13:                               # =>This Inner Loop Header: Depth=1
	vmovss	xmm0, dword ptr [rdi + 4*rax]   # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [r9 + 4*rax], xmm0
	add	rax, 1
	add	rcx, -1
	jne	.LBB0_13
.LBB0_14:
	cmp	rbx, 3
	mov	ecx, dword ptr [rbp + 16]
	jb	.LBB0_16
	.p2align	4, 0x90
.LBB0_15:                               # =>This Inner Loop Header: Depth=1
	vmovss	xmm0, dword ptr [rdi + 4*rax]   # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [r9 + 4*rax], xmm0
	vmovss	xmm0, dword ptr [rdi + 4*rax + 4] # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [r9 + 4*rax + 4], xmm0
	vmovss	xmm0, dword ptr [rdi + 4*rax + 8] # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [r9 + 4*rax + 8], xmm0
	vmovss	xmm0, dword ptr [rdi + 4*rax + 12] # xmm0 = mem[0],zero,zero,zero
	vmovss	dword ptr [r9 + 4*rax + 12], xmm0
	add	rax, 4
	cmp	r12, rax
	jne	.LBB0_15
	jmp	.LBB0_16
.LBB0_6:
	mov	eax, r12d
	and	eax, -32
	lea	rcx, [rax - 32]
	mov	r10, rcx
	shr	r10, 5
	add	r10, 1
	test	rcx, rcx
	je	.LBB0_34
# %bb.7:
	mov	r14, r10
	and	r14, -2
	neg	r14
	xor	ebx, ebx
	mov	ecx, dword ptr [rbp + 16]
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
	vmovups	ymm0, ymmword ptr [rdi + 4*rbx + 128]
	vmovups	ymm1, ymmword ptr [rdi + 4*rbx + 160]
	vmovups	ymm2, ymmword ptr [rdi + 4*rbx + 192]
	vmovups	ymm3, ymmword ptr [rdi + 4*rbx + 224]
	vmovups	ymmword ptr [r9 + 4*rbx + 128], ymm0
	vmovups	ymmword ptr [r9 + 4*rbx + 160], ymm1
	vmovups	ymmword ptr [r9 + 4*rbx + 192], ymm2
	vmovups	ymmword ptr [r9 + 4*rbx + 224], ymm3
	add	rbx, 64
	add	r14, 2
	jne	.LBB0_8
# %bb.9:
	test	r10b, 1
	je	.LBB0_11
.LBB0_10:
	vmovups	ymm0, ymmword ptr [rdi + 4*rbx]
	vmovups	ymm1, ymmword ptr [rdi + 4*rbx + 32]
	vmovups	ymm2, ymmword ptr [rdi + 4*rbx + 64]
	vmovups	ymm3, ymmword ptr [rdi + 4*rbx + 96]
	vmovups	ymmword ptr [r9 + 4*rbx], ymm0
	vmovups	ymmword ptr [r9 + 4*rbx + 32], ymm1
	vmovups	ymmword ptr [r9 + 4*rbx + 64], ymm2
	vmovups	ymmword ptr [r9 + 4*rbx + 96], ymm3
.LBB0_11:
	cmp	rax, r12
	jne	.LBB0_12
.LBB0_16:
	test	ecx, ecx
	jle	.LBB0_33
# %bb.17:
	test	r11d, r11d
	jle	.LBB0_33
# %bb.18:
	mov	r8d, r11d
	lea	rax, [r9 + 4*r12]
	mov	qword ptr [rsp + 16], rax       # 8-byte Spill
	mov	rax, qword ptr [rsp]            # 8-byte Reload
	lea	rdx, [rax + 4*r12]
	mov	r11d, r12d
	and	r11d, -32
	mov	r13, r12
	neg	r13
	lea	rdi, [rax + 96]
	mov	qword ptr [rsp + 8], rdi        # 8-byte Spill
	add	rax, 4
	mov	qword ptr [rsp + 24], rax       # 8-byte Spill
	xor	edi, edi
	movsxd	r15, ecx
	jmp	.LBB0_20
	.p2align	4, 0x90
.LBB0_19:                               #   in Loop: Header=BB0_20 Depth=1
	add	rdi, 1
	cmp	rdi, r8
	je	.LBB0_33
.LBB0_20:                               # =>This Loop Header: Depth=1
                                        #     Child Loop BB0_26 Depth 2
                                        #     Child Loop BB0_32 Depth 2
	movsx	r10, word ptr [rsi + 2*rdi]
	imul	r10, r15
	mov	rax, qword ptr [rsp + 32]       # 8-byte Reload
	movsx	eax, byte ptr [rax + rdi]
	vcvtsi2ss	xmm0, xmm6, eax
	cmp	ecx, 32
	jb	.LBB0_21
# %bb.22:                               #   in Loop: Header=BB0_20 Depth=1
	lea	rax, [rdx + 4*r10]
	cmp	rax, r9
	jbe	.LBB0_25
# %bb.23:                               #   in Loop: Header=BB0_20 Depth=1
	mov	rax, qword ptr [rsp]            # 8-byte Reload
	lea	rax, [rax + 4*r10]
	cmp	rax, qword ptr [rsp + 16]       # 8-byte Folded Reload
	jae	.LBB0_25
.LBB0_21:                               #   in Loop: Header=BB0_20 Depth=1
	xor	ebx, ebx
.LBB0_28:                               #   in Loop: Header=BB0_20 Depth=1
	mov	r14, rbx
	test	r12b, 1
	je	.LBB0_30
# %bb.29:                               #   in Loop: Header=BB0_20 Depth=1
	lea	rax, [rbx + r10]
	mov	rcx, qword ptr [rsp]            # 8-byte Reload
	vmulss	xmm1, xmm0, dword ptr [rcx + 4*rax]
	mov	ecx, dword ptr [rbp + 16]
	vaddss	xmm1, xmm1, dword ptr [r9 + 4*rbx]
	vmovss	dword ptr [r9 + 4*rbx], xmm1
	mov	r14, rbx
	or	r14, 1
.LBB0_30:                               #   in Loop: Header=BB0_20 Depth=1
	not	rbx
	cmp	rbx, r13
	je	.LBB0_19
# %bb.31:                               #   in Loop: Header=BB0_20 Depth=1
	mov	rax, qword ptr [rsp + 24]       # 8-byte Reload
	lea	rax, [rax + 4*r10]
	.p2align	4, 0x90
.LBB0_32:                               #   Parent Loop BB0_20 Depth=1
                                        # =>  This Inner Loop Header: Depth=2
	vmulss	xmm1, xmm0, dword ptr [rax + 4*r14 - 4]
	vaddss	xmm1, xmm1, dword ptr [r9 + 4*r14]
	vmovss	dword ptr [r9 + 4*r14], xmm1
	vmulss	xmm1, xmm0, dword ptr [rax + 4*r14]
	vaddss	xmm1, xmm1, dword ptr [r9 + 4*r14 + 4]
	vmovss	dword ptr [r9 + 4*r14 + 4], xmm1
	add	r14, 2
	cmp	r12, r14
	jne	.LBB0_32
	jmp	.LBB0_19
	.p2align	4, 0x90
.LBB0_25:                               #   in Loop: Header=BB0_20 Depth=1
	vpermilps	xmm1, xmm0, 0           # xmm1 = xmm0[0,0,0,0]
	vinsertf128	ymm1, ymm1, xmm1, 1
	mov	rax, qword ptr [rsp + 8]        # 8-byte Reload
	lea	rbx, [rax + 4*r10]
	xor	eax, eax
	.p2align	4, 0x90
.LBB0_26:                               #   Parent Loop BB0_20 Depth=1
                                        # =>  This Inner Loop Header: Depth=2
	vmulps	ymm2, ymm1, ymmword ptr [rbx + 4*rax - 96]
	vmulps	ymm3, ymm1, ymmword ptr [rbx + 4*rax - 64]
	vmulps	ymm4, ymm1, ymmword ptr [rbx + 4*rax - 32]
	vmulps	ymm5, ymm1, ymmword ptr [rbx + 4*rax]
	vaddps	ymm2, ymm2, ymmword ptr [r9 + 4*rax]
	vaddps	ymm3, ymm3, ymmword ptr [r9 + 4*rax + 32]
	vaddps	ymm4, ymm4, ymmword ptr [r9 + 4*rax + 64]
	vaddps	ymm5, ymm5, ymmword ptr [r9 + 4*rax + 96]
	vmovups	ymmword ptr [r9 + 4*rax], ymm2
	vmovups	ymmword ptr [r9 + 4*rax + 32], ymm3
	vmovups	ymmword ptr [r9 + 4*rax + 64], ymm4
	vmovups	ymmword ptr [r9 + 4*rax + 96], ymm5
	add	rax, 32
	cmp	r11, rax
	jne	.LBB0_26
# %bb.27:                               #   in Loop: Header=BB0_20 Depth=1
	mov	rbx, r11
	cmp	r11, r12
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
	mov	ecx, dword ptr [rbp + 16]
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
