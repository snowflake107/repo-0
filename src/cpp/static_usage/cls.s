
cls.o:     file format elf64-x86-64


Disassembly of section .group:

0000000000000000 <.group>:
   0:	01 00                	add    %eax,(%rax)
   2:	00 00                	add    %al,(%rax)
   4:	07                   	(bad)  
   5:	00 00                	add    %al,(%rax)
   7:	00 08                	add    %cl,(%rax)
   9:	00 00                	add    %al,(%rax)
	...

Disassembly of section .group:

0000000000000000 <.group>:
   0:	01 00                	add    %eax,(%rax)
   2:	00 00                	add    %al,(%rax)
   4:	09 00                	or     %eax,(%rax)
	...

Disassembly of section .text:

0000000000000000 <_Z3barv>:
   0:	55                   	push   %rbp
   1:	48 89 e5             	mov    %rsp,%rbp
   4:	e8 00 00 00 00       	call   9 <_Z3barv+0x9>
			5: R_X86_64_PLT32	_ZN1A11func_staticEv-0x4
   9:	5d                   	pop    %rbp
   a:	c3                   	ret    

000000000000000b <_Z3foov>:
   b:	55                   	push   %rbp
   c:	48 89 e5             	mov    %rsp,%rbp
   f:	48 83 ec 10          	sub    $0x10,%rsp
  13:	64 48 8b 04 25 28 00 	mov    %fs:0x28,%rax
  1a:	00 00 
  1c:	48 89 45 f8          	mov    %rax,-0x8(%rbp)
  20:	31 c0                	xor    %eax,%eax
  22:	c7 45 f4 0a 00 00 00 	movl   $0xa,-0xc(%rbp)
  29:	48 8d 45 f4          	lea    -0xc(%rbp),%rax
  2d:	48 89 c7             	mov    %rax,%rdi
  30:	e8 00 00 00 00       	call   35 <_Z3foov+0x2a>
			31: R_X86_64_PLT32	_ZN1A4funcEv-0x4
  35:	48 8b 55 f8          	mov    -0x8(%rbp),%rdx
  39:	64 48 2b 14 25 28 00 	sub    %fs:0x28,%rdx
  40:	00 00 
  42:	74 05                	je     49 <_Z3foov+0x3e>
  44:	e8 00 00 00 00       	call   49 <_Z3foov+0x3e>
			45: R_X86_64_PLT32	__stack_chk_fail-0x4
  49:	c9                   	leave  
  4a:	c3                   	ret    

Disassembly of section .data:

0000000000000000 <_ZN1A1iE>:
   0:	14 00                	adc    $0x0,%al
	...

Disassembly of section .text._ZN1A11func_staticEv:

0000000000000000 <_ZN1A11func_staticEv>:
   0:	55                   	push   %rbp
   1:	48 89 e5             	mov    %rsp,%rbp
   4:	8b 05 00 00 00 00    	mov    0x0(%rip),%eax        # a <_ZN1A11func_staticEv+0xa>
			6: R_X86_64_PC32	_ZN1A1iE-0x4
   a:	5d                   	pop    %rbp
   b:	c3                   	ret    

Disassembly of section .text._ZN1A4funcEv:

0000000000000000 <_ZN1A4funcEv>:
   0:	55                   	push   %rbp
   1:	48 89 e5             	mov    %rsp,%rbp
   4:	48 89 7d f8          	mov    %rdi,-0x8(%rbp)
   8:	48 8b 45 f8          	mov    -0x8(%rbp),%rax
   c:	8b 00                	mov    (%rax),%eax
   e:	5d                   	pop    %rbp
   f:	c3                   	ret    

Disassembly of section .comment:

0000000000000000 <.comment>:
   0:	00 47 43             	add    %al,0x43(%rdi)
   3:	43 3a 20             	rex.XB cmp (%r8),%spl
   6:	28 47 4e             	sub    %al,0x4e(%rdi)
   9:	55                   	push   %rbp
   a:	29 20                	sub    %esp,(%rax)
   c:	31 31                	xor    %esi,(%rcx)
   e:	2e 32 2e             	cs xor (%rsi),%ch
  11:	30 00                	xor    %al,(%rax)

Disassembly of section .note.gnu.property:

0000000000000000 <.note.gnu.property>:
   0:	04 00                	add    $0x0,%al
   2:	00 00                	add    %al,(%rax)
   4:	20 00                	and    %al,(%rax)
   6:	00 00                	add    %al,(%rax)
   8:	05 00 00 00 47       	add    $0x47000000,%eax
   d:	4e 55                	rex.WRX push %rbp
   f:	00 02                	add    %al,(%rdx)
  11:	00 01                	add    %al,(%rcx)
  13:	c0 04 00 00          	rolb   $0x0,(%rax,%rax,1)
	...
  1f:	00 01                	add    %al,(%rcx)
  21:	00 01                	add    %al,(%rcx)
  23:	c0 04 00 00          	rolb   $0x0,(%rax,%rax,1)
  27:	00 01                	add    %al,(%rcx)
  29:	00 00                	add    %al,(%rax)
  2b:	00 00                	add    %al,(%rax)
  2d:	00 00                	add    %al,(%rax)
	...

Disassembly of section .eh_frame:

0000000000000000 <.eh_frame>:
   0:	14 00                	adc    $0x0,%al
   2:	00 00                	add    %al,(%rax)
   4:	00 00                	add    %al,(%rax)
   6:	00 00                	add    %al,(%rax)
   8:	01 7a 52             	add    %edi,0x52(%rdx)
   b:	00 01                	add    %al,(%rcx)
   d:	78 10                	js     1f <.eh_frame+0x1f>
   f:	01 1b                	add    %ebx,(%rbx)
  11:	0c 07                	or     $0x7,%al
  13:	08 90 01 00 00 1c    	or     %dl,0x1c000001(%rax)
  19:	00 00                	add    %al,(%rax)
  1b:	00 1c 00             	add    %bl,(%rax,%rax,1)
  1e:	00 00                	add    %al,(%rax)
  20:	00 00                	add    %al,(%rax)
			20: R_X86_64_PC32	.text._ZN1A11func_staticEv
  22:	00 00                	add    %al,(%rax)
  24:	0c 00                	or     $0x0,%al
  26:	00 00                	add    %al,(%rax)
  28:	00 41 0e             	add    %al,0xe(%rcx)
  2b:	10 86 02 43 0d 06    	adc    %al,0x60d4302(%rsi)
  31:	47 0c 07             	rex.RXB or $0x7,%al
  34:	08 00                	or     %al,(%rax)
  36:	00 00                	add    %al,(%rax)
  38:	1c 00                	sbb    $0x0,%al
  3a:	00 00                	add    %al,(%rax)
  3c:	3c 00                	cmp    $0x0,%al
  3e:	00 00                	add    %al,(%rax)
  40:	00 00                	add    %al,(%rax)
			40: R_X86_64_PC32	.text._ZN1A4funcEv
  42:	00 00                	add    %al,(%rax)
  44:	10 00                	adc    %al,(%rax)
  46:	00 00                	add    %al,(%rax)
  48:	00 41 0e             	add    %al,0xe(%rcx)
  4b:	10 86 02 43 0d 06    	adc    %al,0x60d4302(%rsi)
  51:	4b 0c 07             	rex.WXB or $0x7,%al
  54:	08 00                	or     %al,(%rax)
  56:	00 00                	add    %al,(%rax)
  58:	1c 00                	sbb    $0x0,%al
  5a:	00 00                	add    %al,(%rax)
  5c:	5c                   	pop    %rsp
  5d:	00 00                	add    %al,(%rax)
  5f:	00 00                	add    %al,(%rax)
			60: R_X86_64_PC32	.text
  61:	00 00                	add    %al,(%rax)
  63:	00 0b                	add    %cl,(%rbx)
  65:	00 00                	add    %al,(%rax)
  67:	00 00                	add    %al,(%rax)
  69:	41 0e                	rex.B (bad) 
  6b:	10 86 02 43 0d 06    	adc    %al,0x60d4302(%rsi)
  71:	46 0c 07             	rex.RX or $0x7,%al
  74:	08 00                	or     %al,(%rax)
  76:	00 00                	add    %al,(%rax)
  78:	1c 00                	sbb    $0x0,%al
  7a:	00 00                	add    %al,(%rax)
  7c:	7c 00                	jl     7e <.eh_frame+0x7e>
  7e:	00 00                	add    %al,(%rax)
  80:	00 00                	add    %al,(%rax)
			80: R_X86_64_PC32	.text+0xb
  82:	00 00                	add    %al,(%rax)
  84:	40 00 00             	rex add %al,(%rax)
  87:	00 00                	add    %al,(%rax)
  89:	41 0e                	rex.B (bad) 
  8b:	10 86 02 43 0d 06    	adc    %al,0x60d4302(%rsi)
  91:	7b 0c                	jnp    9f <_Z3foov+0x94>
  93:	07                   	(bad)  
  94:	08 00                	or     %al,(%rax)
	...
