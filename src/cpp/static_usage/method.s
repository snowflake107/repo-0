	.file	"method.cpp"
	.text
	.data
	.align 4
	.type	_ZZ6methodvE1i, @object
	.size	_ZZ6methodvE1i, 4
_ZZ6methodvE1i:
	.long	10
	.text
	.globl	_Z6methodv
	.type	_Z6methodv, @function
_Z6methodv:
.LFB0:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	movl	_ZZ6methodvE1i(%rip), %eax
	leal	1(%rax), %edx
	movl	%edx, _ZZ6methodvE1i(%rip)
	popq	%rbp
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE0:
	.size	_Z6methodv, .-_Z6methodv
	.ident	"GCC: (GNU) 11.2.0"
	.section	.note.GNU-stack,"",@progbits
