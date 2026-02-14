package z80

import "testing"

func TestLD_r_r(t *testing.T) {
	cpu, bus := newTestCPU()
	// LD B, A (0x47)
	bus.mem[0] = 0x47
	cpu.reg.AF = 0x4200
	cycles := cpu.Step()
	if cpu.getB() != 0x42 {
		t.Errorf("LD B,A: B=%02x want 42", cpu.getB())
	}
	if cycles != 4 {
		t.Errorf("LD B,A: cycles=%d want 4", cycles)
	}
}

func TestLD_r_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	// LD A, (HL) (0x7E)
	bus.mem[0] = 0x7E
	cpu.reg.HL = 0x1000
	bus.mem[0x1000] = 0x99
	cycles := cpu.Step()
	if cpu.getA() != 0x99 {
		t.Errorf("LD A,(HL): A=%02x want 99", cpu.getA())
	}
	if cycles != 7 {
		t.Errorf("LD A,(HL): cycles=%d want 7", cycles)
	}
}

func TestLD_HL_r(t *testing.T) {
	cpu, bus := newTestCPU()
	// LD (HL), C (0x71)
	bus.mem[0] = 0x71
	cpu.reg.HL = 0x2000
	cpu.reg.BC = 0x0055
	cycles := cpu.Step()
	if bus.mem[0x2000] != 0x55 {
		t.Errorf("LD (HL),C: (HL)=%02x want 55", bus.mem[0x2000])
	}
	if cycles != 7 {
		t.Errorf("LD (HL),C: cycles=%d want 7", cycles)
	}
}

func TestLD_r_n(t *testing.T) {
	cpu, bus := newTestCPU()
	// LD D, 0xAB (0x16 0xAB)
	bus.mem[0] = 0x16
	bus.mem[1] = 0xAB
	cycles := cpu.Step()
	if cpu.getD() != 0xAB {
		t.Errorf("LD D,n: D=%02x want AB", cpu.getD())
	}
	if cycles != 7 {
		t.Errorf("LD D,n: cycles=%d want 7", cycles)
	}
}

func TestLD_HL_n(t *testing.T) {
	cpu, bus := newTestCPU()
	// LD (HL), 0x77 (0x36 0x77)
	bus.mem[0] = 0x36
	bus.mem[1] = 0x77
	cpu.reg.HL = 0x3000
	cycles := cpu.Step()
	if bus.mem[0x3000] != 0x77 {
		t.Errorf("LD (HL),n: (HL)=%02x want 77", bus.mem[0x3000])
	}
	if cycles != 10 {
		t.Errorf("LD (HL),n: cycles=%d want 10", cycles)
	}
}

func TestLD_A_BC(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x0A
	cpu.reg.BC = 0x4000
	bus.mem[0x4000] = 0xEE
	cycles := cpu.Step()
	if cpu.getA() != 0xEE {
		t.Errorf("LD A,(BC): A=%02x want EE", cpu.getA())
	}
	if cycles != 7 {
		t.Errorf("cycles=%d want 7", cycles)
	}
}

func TestLD_A_DE(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x1A
	cpu.reg.DE = 0x5000
	bus.mem[0x5000] = 0xDD
	cpu.Step()
	if cpu.getA() != 0xDD {
		t.Errorf("LD A,(DE): A=%02x want DD", cpu.getA())
	}
}

func TestLD_BC_A(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x02
	cpu.reg.AF = 0xAA00
	cpu.reg.BC = 0x6000
	cpu.Step()
	if bus.mem[0x6000] != 0xAA {
		t.Errorf("LD (BC),A: got %02x want AA", bus.mem[0x6000])
	}
}

func TestLD_A_nn(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x3A
	bus.mem[1] = 0x00 // low
	bus.mem[2] = 0x80 // high -> 0x8000
	bus.mem[0x8000] = 0xBB
	cycles := cpu.Step()
	if cpu.getA() != 0xBB {
		t.Errorf("LD A,(nn): A=%02x want BB", cpu.getA())
	}
	if cycles != 13 {
		t.Errorf("cycles=%d want 13", cycles)
	}
}

func TestLD_nn_A(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x32
	bus.mem[1] = 0x00
	bus.mem[2] = 0x90
	cpu.reg.AF = 0xCC00
	cpu.Step()
	if bus.mem[0x9000] != 0xCC {
		t.Errorf("LD (nn),A: got %02x want CC", bus.mem[0x9000])
	}
}

func TestLD_rr_nn(t *testing.T) {
	cpu, bus := newTestCPU()
	// LD BC, 0x1234 (0x01 0x34 0x12)
	bus.mem[0] = 0x01
	bus.mem[1] = 0x34
	bus.mem[2] = 0x12
	cycles := cpu.Step()
	if cpu.reg.BC != 0x1234 {
		t.Errorf("LD BC,nn: BC=%04x want 1234", cpu.reg.BC)
	}
	if cycles != 10 {
		t.Errorf("cycles=%d want 10", cycles)
	}
}

func TestLD_nn_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x22
	bus.mem[1] = 0x00
	bus.mem[2] = 0xA0
	cpu.reg.HL = 0xBEEF
	cycles := cpu.Step()
	if bus.mem[0xA000] != 0xEF || bus.mem[0xA001] != 0xBE {
		t.Errorf("LD (nn),HL: got %02x%02x want BEEF", bus.mem[0xA001], bus.mem[0xA000])
	}
	if cycles != 16 {
		t.Errorf("cycles=%d want 16", cycles)
	}
}

func TestLD_HL_nn_indirect(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x2A
	bus.mem[1] = 0x00
	bus.mem[2] = 0xB0
	bus.mem[0xB000] = 0xAD
	bus.mem[0xB001] = 0xDE
	cycles := cpu.Step()
	if cpu.reg.HL != 0xDEAD {
		t.Errorf("LD HL,(nn): HL=%04x want DEAD", cpu.reg.HL)
	}
	if cycles != 16 {
		t.Errorf("cycles=%d want 16", cycles)
	}
}

func TestLD_SP_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xF9
	cpu.reg.HL = 0x4321
	cycles := cpu.Step()
	if cpu.reg.SP != 0x4321 {
		t.Errorf("LD SP,HL: SP=%04x want 4321", cpu.reg.SP)
	}
	if cycles != 6 {
		t.Errorf("cycles=%d want 6", cycles)
	}
}

func TestPUSH_POP(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFE
	cpu.reg.BC = 0x1234
	// PUSH BC (0xC5)
	bus.mem[0] = 0xC5
	c1 := cpu.Step()
	if c1 != 11 {
		t.Errorf("PUSH BC: cycles=%d want 11", c1)
	}
	if cpu.reg.SP != 0xFFFC {
		t.Errorf("PUSH BC: SP=%04x want FFFC", cpu.reg.SP)
	}
	// POP DE (0xD1)
	bus.mem[1] = 0xD1
	c2 := cpu.Step()
	if c2 != 10 {
		t.Errorf("POP DE: cycles=%d want 10", c2)
	}
	if cpu.reg.DE != 0x1234 {
		t.Errorf("POP DE: DE=%04x want 1234", cpu.reg.DE)
	}
}

func TestEX_DE_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xEB
	cpu.reg.DE = 0x1111
	cpu.reg.HL = 0x2222
	cpu.Step()
	if cpu.reg.DE != 0x2222 || cpu.reg.HL != 0x1111 {
		t.Errorf("EX DE,HL: DE=%04x HL=%04x", cpu.reg.DE, cpu.reg.HL)
	}
}

func TestEX_AF_AF(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x08
	cpu.reg.AF = 0xAAAA
	cpu.reg.AF_ = 0x5555
	cpu.Step()
	if cpu.reg.AF != 0x5555 || cpu.reg.AF_ != 0xAAAA {
		t.Errorf("EX AF,AF': AF=%04x AF'=%04x", cpu.reg.AF, cpu.reg.AF_)
	}
}

func TestEXX(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xD9
	cpu.reg.BC = 0x1111
	cpu.reg.BC_ = 0x2222
	cpu.Step()
	if cpu.reg.BC != 0x2222 || cpu.reg.BC_ != 0x1111 {
		t.Errorf("EXX: BC=%04x BC'=%04x", cpu.reg.BC, cpu.reg.BC_)
	}
}

func TestEX_SP_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xE3
	cpu.reg.SP = 0xFFF0
	cpu.reg.HL = 0xAAAA
	bus.mem[0xFFF0] = 0x34
	bus.mem[0xFFF1] = 0x12
	cycles := cpu.Step()
	if cpu.reg.HL != 0x1234 {
		t.Errorf("EX (SP),HL: HL=%04x want 1234", cpu.reg.HL)
	}
	if bus.mem[0xFFF0] != 0xAA || bus.mem[0xFFF1] != 0xAA {
		t.Errorf("EX (SP),HL: stack=%02x%02x want AAAA", bus.mem[0xFFF1], bus.mem[0xFFF0])
	}
	if cycles != 19 {
		t.Errorf("cycles=%d want 19", cycles)
	}
}

func TestHALT_Instruction(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x76
	cycles := cpu.Step()
	if !cpu.reg.Halted {
		t.Error("HALT should set Halted")
	}
	if cycles != 4 {
		t.Errorf("cycles=%d want 4", cycles)
	}
	// PC should point past HALT (to address 1)
	if cpu.reg.PC != 1 {
		t.Errorf("HALT: PC=%04x want 0001 (past HALT)", cpu.reg.PC)
	}
}

func TestHALT_NMI_ReturnContinuesPastHALT(t *testing.T) {
	// Regression: HALT must not rewind PC. After NMI handler returns
	// via RETN, execution should continue at the instruction AFTER HALT,
	// not re-execute HALT (which would trap the CPU forever).
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFE

	// 0x0100: HALT
	// 0x0101: NOP (should execute after returning from NMI handler)
	bus.mem[0x0100] = 0x76 // HALT
	bus.mem[0x0101] = 0x00 // NOP (target after RETN)

	// NMI handler at 0x0066: RETN (ED 45)
	bus.mem[0x0066] = 0xED
	bus.mem[0x0067] = 0x45

	cpu.reg.PC = 0x0100
	cpu.reg.IFF1 = true

	// Execute HALT
	cpu.Step()
	if !cpu.reg.Halted {
		t.Fatal("should be halted")
	}
	if cpu.reg.PC != 0x0101 {
		t.Fatalf("HALT: PC=%04x want 0101", cpu.reg.PC)
	}

	// Trigger NMI — should wake HALT and jump to 0x0066
	cpu.NMI()
	cpu.Step()
	if cpu.reg.PC != 0x0066 {
		t.Fatalf("after NMI: PC=%04x want 0066", cpu.reg.PC)
	}

	// Verify return address on stack is 0x0101 (past HALT)
	retAddr := cpu.read16(cpu.reg.SP)
	if retAddr != 0x0101 {
		t.Fatalf("return addr=%04x want 0101 (past HALT)", retAddr)
	}

	// Execute RETN — should return to 0x0101
	cpu.Step()
	if cpu.reg.PC != 0x0101 {
		t.Fatalf("after RETN: PC=%04x want 0101", cpu.reg.PC)
	}
	if cpu.reg.Halted {
		t.Fatal("should NOT be halted after returning from NMI")
	}

	// Execute NOP at 0x0101 — should succeed normally
	cpu.Step()
	if cpu.reg.PC != 0x0102 {
		t.Errorf("after NOP: PC=%04x want 0102", cpu.reg.PC)
	}
}

func TestHALT_EI_INT_Pattern(t *testing.T) {
	// The classic EI; HALT pattern: interrupt should wake HALT
	// and the return address should be past HALT.
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFE
	cpu.reg.IM = 1

	// 0x0000: EI (0xFB)
	// 0x0001: HALT (0x76)
	// 0x0002: NOP (target after RETI)
	bus.mem[0x0000] = 0xFB // EI
	bus.mem[0x0001] = 0x76 // HALT
	bus.mem[0x0002] = 0x00 // NOP

	// IM1 handler at 0x0038: RETI (ED 4D)
	bus.mem[0x0038] = 0xED
	bus.mem[0x0039] = 0x4D

	// Assert INT before EI
	cpu.INT(true, 0xFF)

	// Execute EI — sets IFF1, afterEI suppresses immediate INT
	cpu.Step()
	if !cpu.reg.IFF1 {
		t.Fatal("IFF1 should be set after EI")
	}

	// Execute HALT — afterEI cleared, but INT not checked until next Step
	cpu.Step()
	if !cpu.reg.Halted {
		t.Fatal("should be halted")
	}

	// Next Step: INT is asserted, IFF1 is set, services interrupt
	cpu.Step()
	if cpu.reg.PC != 0x0038 {
		t.Fatalf("INT should jump to 0038, PC=%04x", cpu.reg.PC)
	}

	// Verify return address is 0x0002 (past HALT)
	retAddr := cpu.read16(cpu.reg.SP)
	if retAddr != 0x0002 {
		t.Fatalf("return addr=%04x want 0002", retAddr)
	}

	// Execute RETI
	cpu.Step()
	if cpu.reg.PC != 0x0002 {
		t.Fatalf("after RETI: PC=%04x want 0002", cpu.reg.PC)
	}
	if cpu.reg.Halted {
		t.Fatal("should NOT be halted after returning from INT handler")
	}
}

func TestDI_EI(t *testing.T) {
	cpu, bus := newTestCPU()
	// DI
	cpu.reg.IFF1 = true
	cpu.reg.IFF2 = true
	bus.mem[0] = 0xF3
	cpu.Step()
	if cpu.reg.IFF1 || cpu.reg.IFF2 {
		t.Error("DI should clear IFF1 and IFF2")
	}
	// EI
	bus.mem[1] = 0xFB
	cpu.Step()
	if !cpu.reg.IFF1 || !cpu.reg.IFF2 {
		t.Error("EI should set IFF1 and IFF2")
	}
	if !cpu.afterEI {
		t.Error("EI should set afterEI")
	}
}

func TestSST_Load(t *testing.T) {
	tests := []struct {
		name       string
		init, want z80State
	}{
		{ // 00 0000
			name: "00 0000",
			init: z80State{
				A: 0x6E, F: 0xFA,
				B: 0xB9, C: 0x90,
				D: 0xD0, E: 0xBE,
				H: 0x83, L: 0x93,
				I: 0xA6, R: 0x10,
				PC: 0x4DDF, SP: 0xE82E,
				IX: 0x8C13, IY: 0xB28C,
				AF_: 0x7631, BC_: 0x440B,
				DE_: 0x3612, HL_: 0x6E81,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{19935, 0}},
			},
			want: z80State{
				A: 0x6E, F: 0xFA,
				B: 0xB9, C: 0x90,
				D: 0xD0, E: 0xBE,
				H: 0x83, L: 0x93,
				I: 0xA6, R: 0x11,
				PC: 0x4DE0, SP: 0xE82E,
				IX: 0x8C13, IY: 0xB28C,
				AF_: 0x7631, BC_: 0x440B,
				DE_: 0x3612, HL_: 0x6E81,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{19935, 0}},
				Cycles: 4,
			},
		},
		{ // 47 0000
			name: "47 0000",
			init: z80State{
				A: 0x2D, F: 0x67,
				B: 0xB7, C: 0x38,
				D: 0xBC, E: 0x10,
				H: 0xEB, L: 0xFD,
				I: 0x86, R: 0x51,
				PC: 0x5282, SP: 0xE279,
				IX: 0x72B8, IY: 0x0B21,
				AF_: 0x732F, BC_: 0x85D2,
				DE_: 0x9F11, HL_: 0x820F,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{21122, 71}},
			},
			want: z80State{
				A: 0x2D, F: 0x67,
				B: 0x2D, C: 0x38,
				D: 0xBC, E: 0x10,
				H: 0xEB, L: 0xFD,
				I: 0x86, R: 0x52,
				PC: 0x5283, SP: 0xE279,
				IX: 0x72B8, IY: 0x0B21,
				AF_: 0x732F, BC_: 0x85D2,
				DE_: 0x9F11, HL_: 0x820F,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{21122, 71}},
				Cycles: 4,
			},
		},
		{ // 46 0000
			name: "46 0000",
			init: z80State{
				A: 0x79, F: 0x87,
				B: 0x21, C: 0x17,
				D: 0xE3, E: 0x14,
				H: 0x9B, L: 0xC3,
				I: 0x5C, R: 0x68,
				PC: 0x2474, SP: 0xDBA9,
				IX: 0xA55A, IY: 0xC68C,
				AF_: 0x82BF, BC_: 0x560C,
				DE_: 0x6633, HL_: 0x3E2E,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{9332, 70}, {39875, 72}},
			},
			want: z80State{
				A: 0x79, F: 0x87,
				B: 0x48, C: 0x17,
				D: 0xE3, E: 0x14,
				H: 0x9B, L: 0xC3,
				I: 0x5C, R: 0x69,
				PC: 0x2475, SP: 0xDBA9,
				IX: 0xA55A, IY: 0xC68C,
				AF_: 0x82BF, BC_: 0x560C,
				DE_: 0x6633, HL_: 0x3E2E,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{9332, 70}, {39875, 72}},
				Cycles: 7,
			},
		},
		{ // 70 0000
			name: "70 0000",
			init: z80State{
				A: 0xCB, F: 0x12,
				B: 0xD5, C: 0x6A,
				D: 0xDB, E: 0xEE,
				H: 0x79, L: 0x8D,
				I: 0x92, R: 0x19,
				PC: 0x3BFE, SP: 0x279F,
				IX: 0x558D, IY: 0x8CC0,
				AF_: 0xDD54, BC_: 0xA51A,
				DE_: 0x8D14, HL_: 0xD542,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{15358, 112}, {31117, 0}},
			},
			want: z80State{
				A: 0xCB, F: 0x12,
				B: 0xD5, C: 0x6A,
				D: 0xDB, E: 0xEE,
				H: 0x79, L: 0x8D,
				I: 0x92, R: 0x1A,
				PC: 0x3BFF, SP: 0x279F,
				IX: 0x558D, IY: 0x8CC0,
				AF_: 0xDD54, BC_: 0xA51A,
				DE_: 0x8D14, HL_: 0xD542,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{15358, 112}, {31117, 213}},
				Cycles: 7,
			},
		},
		{ // 3E 0000
			name: "3E 0000",
			init: z80State{
				A: 0x96, F: 0xF0,
				B: 0xC9, C: 0xD8,
				D: 0xB2, E: 0x78,
				H: 0xEC, L: 0xFA,
				I: 0xBC, R: 0x0E,
				PC: 0x2F94, SP: 0xD50A,
				IX: 0x3201, IY: 0x9295,
				AF_: 0x07EE, BC_: 0x8350,
				DE_: 0x0A89, HL_: 0xF8CE,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{12180, 62}, {12181, 169}},
			},
			want: z80State{
				A: 0xA9, F: 0xF0,
				B: 0xC9, C: 0xD8,
				D: 0xB2, E: 0x78,
				H: 0xEC, L: 0xFA,
				I: 0xBC, R: 0x0F,
				PC: 0x2F96, SP: 0xD50A,
				IX: 0x3201, IY: 0x9295,
				AF_: 0x07EE, BC_: 0x8350,
				DE_: 0x0A89, HL_: 0xF8CE,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{12180, 62}, {12181, 169}},
				Cycles: 7,
			},
		},
		{ // 36 0000
			name: "36 0000",
			init: z80State{
				A: 0xD2, F: 0xA0,
				B: 0x4A, C: 0x24,
				D: 0x30, E: 0x69,
				H: 0x0A, L: 0x1A,
				I: 0xA3, R: 0x02,
				PC: 0x36A4, SP: 0x2F8D,
				IX: 0x7C62, IY: 0xB95C,
				AF_: 0x0C98, BC_: 0x1DDD,
				DE_: 0x3899, HL_: 0x2980,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{2586, 0}, {13988, 54}, {13989, 254}},
			},
			want: z80State{
				A: 0xD2, F: 0xA0,
				B: 0x4A, C: 0x24,
				D: 0x30, E: 0x69,
				H: 0x0A, L: 0x1A,
				I: 0xA3, R: 0x03,
				PC: 0x36A6, SP: 0x2F8D,
				IX: 0x7C62, IY: 0xB95C,
				AF_: 0x0C98, BC_: 0x1DDD,
				DE_: 0x3899, HL_: 0x2980,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{2586, 254}, {13988, 54}, {13989, 254}},
				Cycles: 10,
			},
		},
		{ // 0A 0000
			name: "0A 0000",
			init: z80State{
				A: 0x40, F: 0x4F,
				B: 0x5F, C: 0xCD,
				D: 0x93, E: 0xA8,
				H: 0x62, L: 0xFB,
				I: 0x3C, R: 0x42,
				PC: 0x41BA, SP: 0x24A7,
				IX: 0x38C0, IY: 0x4964,
				AF_: 0xBE4E, BC_: 0xB094,
				DE_: 0xA735, HL_: 0xE848,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{16826, 10}, {24525, 174}},
			},
			want: z80State{
				A: 0xAE, F: 0x4F,
				B: 0x5F, C: 0xCD,
				D: 0x93, E: 0xA8,
				H: 0x62, L: 0xFB,
				I: 0x3C, R: 0x43,
				PC: 0x41BB, SP: 0x24A7,
				IX: 0x38C0, IY: 0x4964,
				AF_: 0xBE4E, BC_: 0xB094,
				DE_: 0xA735, HL_: 0xE848,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{16826, 10}, {24525, 174}},
				Cycles: 7,
			},
		},
		{ // 12 0000
			name: "12 0000",
			init: z80State{
				A: 0xE5, F: 0xC0,
				B: 0xC4, C: 0x60,
				D: 0x4E, E: 0x7B,
				H: 0x0D, L: 0xFE,
				I: 0xE2, R: 0x5D,
				PC: 0xE446, SP: 0x70A8,
				IX: 0x9E7C, IY: 0xA811,
				AF_: 0x1D45, BC_: 0x17D6,
				DE_: 0x5EE8, HL_: 0x6DDC,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{20091, 0}, {58438, 18}},
			},
			want: z80State{
				A: 0xE5, F: 0xC0,
				B: 0xC4, C: 0x60,
				D: 0x4E, E: 0x7B,
				H: 0x0D, L: 0xFE,
				I: 0xE2, R: 0x5E,
				PC: 0xE447, SP: 0x70A8,
				IX: 0x9E7C, IY: 0xA811,
				AF_: 0x1D45, BC_: 0x17D6,
				DE_: 0x5EE8, HL_: 0x6DDC,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{20091, 229}, {58438, 18}},
				Cycles: 7,
			},
		},
		{ // 3A 0000
			name: "3A 0000",
			init: z80State{
				A: 0x41, F: 0xAE,
				B: 0xE3, C: 0x3E,
				D: 0x80, E: 0x1F,
				H: 0x56, L: 0x1A,
				I: 0x6E, R: 0x0C,
				PC: 0xC8FC, SP: 0x41D9,
				IX: 0xFB3A, IY: 0x1342,
				AF_: 0x3606, BC_: 0x42F1,
				DE_: 0x7FEF, HL_: 0xC5A1,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{4100, 10}, {51452, 58}, {51453, 4}, {51454, 16}},
			},
			want: z80State{
				A: 0x0A, F: 0xAE,
				B: 0xE3, C: 0x3E,
				D: 0x80, E: 0x1F,
				H: 0x56, L: 0x1A,
				I: 0x6E, R: 0x0D,
				PC: 0xC8FF, SP: 0x41D9,
				IX: 0xFB3A, IY: 0x1342,
				AF_: 0x3606, BC_: 0x42F1,
				DE_: 0x7FEF, HL_: 0xC5A1,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{4100, 10}, {51452, 58}, {51453, 4}, {51454, 16}},
				Cycles: 13,
			},
		},
		{ // 32 0000
			name: "32 0000",
			init: z80State{
				A: 0x97, F: 0xC5,
				B: 0x38, C: 0xEE,
				D: 0x57, E: 0x8F,
				H: 0x98, L: 0xFE,
				I: 0x9F, R: 0x19,
				PC: 0xACE0, SP: 0x7A15,
				IX: 0xCF2B, IY: 0x6F98,
				AF_: 0x5D3C, BC_: 0xEB61,
				DE_: 0x91C8, HL_: 0xCABA,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{44256, 50}, {44257, 133}, {44258, 200}, {51333, 0}},
			},
			want: z80State{
				A: 0x97, F: 0xC5,
				B: 0x38, C: 0xEE,
				D: 0x57, E: 0x8F,
				H: 0x98, L: 0xFE,
				I: 0x9F, R: 0x1A,
				PC: 0xACE3, SP: 0x7A15,
				IX: 0xCF2B, IY: 0x6F98,
				AF_: 0x5D3C, BC_: 0xEB61,
				DE_: 0x91C8, HL_: 0xCABA,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{44256, 50}, {44257, 133}, {44258, 200}, {51333, 151}},
				Cycles: 13,
			},
		},
		{ // 01 0000
			name: "01 0000",
			init: z80State{
				A: 0x51, F: 0xF3,
				B: 0x64, C: 0xA2,
				D: 0xA5, E: 0xBD,
				H: 0x95, L: 0x70,
				I: 0x09, R: 0x56,
				PC: 0xE5FE, SP: 0x3246,
				IX: 0x63CB, IY: 0x76F2,
				AF_: 0xBCCB, BC_: 0x36D7,
				DE_: 0xFB9D, HL_: 0x9050,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{58878, 1}, {58879, 222}, {58880, 16}},
			},
			want: z80State{
				A: 0x51, F: 0xF3,
				B: 0x10, C: 0xDE,
				D: 0xA5, E: 0xBD,
				H: 0x95, L: 0x70,
				I: 0x09, R: 0x57,
				PC: 0xE601, SP: 0x3246,
				IX: 0x63CB, IY: 0x76F2,
				AF_: 0xBCCB, BC_: 0x36D7,
				DE_: 0xFB9D, HL_: 0x9050,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{58878, 1}, {58879, 222}, {58880, 16}},
				Cycles: 10,
			},
		},
		{ // 22 0000
			name: "22 0000",
			init: z80State{
				A: 0xB2, F: 0x3B,
				B: 0x01, C: 0x84,
				D: 0xD7, E: 0x32,
				H: 0xEB, L: 0x0B,
				I: 0x58, R: 0x5D,
				PC: 0xE81E, SP: 0x57E7,
				IX: 0xC3DC, IY: 0xC485,
				AF_: 0xDF5F, BC_: 0x3746,
				DE_: 0x1C8F, HL_: 0xE7DB,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{51553, 0}, {51554, 0}, {59422, 34}, {59423, 97}, {59424, 201}},
			},
			want: z80State{
				A: 0xB2, F: 0x3B,
				B: 0x01, C: 0x84,
				D: 0xD7, E: 0x32,
				H: 0xEB, L: 0x0B,
				I: 0x58, R: 0x5E,
				PC: 0xE821, SP: 0x57E7,
				IX: 0xC3DC, IY: 0xC485,
				AF_: 0xDF5F, BC_: 0x3746,
				DE_: 0x1C8F, HL_: 0xE7DB,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{51553, 11}, {51554, 235}, {59422, 34}, {59423, 97}, {59424, 201}},
				Cycles: 16,
			},
		},
		{ // 2A 0000
			name: "2A 0000",
			init: z80State{
				A: 0xDA, F: 0x65,
				B: 0xDB, C: 0x3D,
				D: 0xE2, E: 0x8A,
				H: 0x11, L: 0xEF,
				I: 0xC8, R: 0x55,
				PC: 0x8493, SP: 0x26F2,
				IX: 0x323C, IY: 0x64F4,
				AF_: 0xAD27, BC_: 0x20AB,
				DE_: 0xA4B5, HL_: 0x4B60,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{19914, 86}, {19915, 104}, {33939, 42}, {33940, 202}, {33941, 77}},
			},
			want: z80State{
				A: 0xDA, F: 0x65,
				B: 0xDB, C: 0x3D,
				D: 0xE2, E: 0x8A,
				H: 0x68, L: 0x56,
				I: 0xC8, R: 0x56,
				PC: 0x8496, SP: 0x26F2,
				IX: 0x323C, IY: 0x64F4,
				AF_: 0xAD27, BC_: 0x20AB,
				DE_: 0xA4B5, HL_: 0x4B60,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{19914, 86}, {19915, 104}, {33939, 42}, {33940, 202}, {33941, 77}},
				Cycles: 16,
			},
		},
		{ // F9 0000
			name: "F9 0000",
			init: z80State{
				A: 0xC5, F: 0x73,
				B: 0x5A, C: 0xCB,
				D: 0x3F, E: 0x1B,
				H: 0x3F, L: 0x35,
				I: 0x1A, R: 0x71,
				PC: 0xEA9A, SP: 0x1382,
				IX: 0x913D, IY: 0x602E,
				AF_: 0x2614, BC_: 0xF81E,
				DE_: 0x0BDB, HL_: 0x7026,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{60058, 249}},
			},
			want: z80State{
				A: 0xC5, F: 0x73,
				B: 0x5A, C: 0xCB,
				D: 0x3F, E: 0x1B,
				H: 0x3F, L: 0x35,
				I: 0x1A, R: 0x72,
				PC: 0xEA9B, SP: 0x3F35,
				IX: 0x913D, IY: 0x602E,
				AF_: 0x2614, BC_: 0xF81E,
				DE_: 0x0BDB, HL_: 0x7026,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{60058, 249}},
				Cycles: 6,
			},
		},
		{ // C5 0000
			name: "C5 0000",
			init: z80State{
				A: 0x31, F: 0xE4,
				B: 0xAF, C: 0xBD,
				D: 0xDD, E: 0x15,
				H: 0xEF, L: 0xCD,
				I: 0x45, R: 0x23,
				PC: 0x079C, SP: 0x6387,
				IX: 0x73CD, IY: 0x5B59,
				AF_: 0x83E9, BC_: 0x516D,
				DE_: 0xE8FA, HL_: 0xEA17,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{1948, 197}, {25477, 0}, {25478, 0}},
			},
			want: z80State{
				A: 0x31, F: 0xE4,
				B: 0xAF, C: 0xBD,
				D: 0xDD, E: 0x15,
				H: 0xEF, L: 0xCD,
				I: 0x45, R: 0x24,
				PC: 0x079D, SP: 0x6385,
				IX: 0x73CD, IY: 0x5B59,
				AF_: 0x83E9, BC_: 0x516D,
				DE_: 0xE8FA, HL_: 0xEA17,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{1948, 197}, {25477, 189}, {25478, 175}},
				Cycles: 11,
			},
		},
		{ // E1 0000
			name: "E1 0000",
			init: z80State{
				A: 0x93, F: 0x64,
				B: 0xC2, C: 0xFD,
				D: 0x29, E: 0xBF,
				H: 0x45, L: 0x92,
				I: 0x9F, R: 0x79,
				PC: 0xA28E, SP: 0x5F59,
				IX: 0x838D, IY: 0xFA80,
				AF_: 0xD38B, BC_: 0x1882,
				DE_: 0xDD6F, HL_: 0x7B3A,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{24409, 224}, {24410, 254}, {41614, 225}},
			},
			want: z80State{
				A: 0x93, F: 0x64,
				B: 0xC2, C: 0xFD,
				D: 0x29, E: 0xBF,
				H: 0xFE, L: 0xE0,
				I: 0x9F, R: 0x7A,
				PC: 0xA28F, SP: 0x5F5B,
				IX: 0x838D, IY: 0xFA80,
				AF_: 0xD38B, BC_: 0x1882,
				DE_: 0xDD6F, HL_: 0x7B3A,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{24409, 224}, {24410, 254}, {41614, 225}},
				Cycles: 10,
			},
		},
		{ // EB 0000
			name: "EB 0000",
			init: z80State{
				A: 0xFF, F: 0x55,
				B: 0x81, C: 0x10,
				D: 0x14, E: 0x2C,
				H: 0x3E, L: 0x81,
				I: 0xE8, R: 0x04,
				PC: 0x7CE6, SP: 0x1966,
				IX: 0x2DF8, IY: 0x88EC,
				AF_: 0x10D0, BC_: 0x8ED3,
				DE_: 0x702A, HL_: 0x3B9F,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{31974, 235}},
			},
			want: z80State{
				A: 0xFF, F: 0x55,
				B: 0x81, C: 0x10,
				D: 0x3E, E: 0x81,
				H: 0x14, L: 0x2C,
				I: 0xE8, R: 0x05,
				PC: 0x7CE7, SP: 0x1966,
				IX: 0x2DF8, IY: 0x88EC,
				AF_: 0x10D0, BC_: 0x8ED3,
				DE_: 0x702A, HL_: 0x3B9F,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{31974, 235}},
				Cycles: 4,
			},
		},
		{ // 08 0000
			name: "08 0000",
			init: z80State{
				A: 0xC7, F: 0x1A,
				B: 0x4F, C: 0x06,
				D: 0x05, E: 0xEF,
				H: 0xE3, L: 0xB5,
				I: 0x1F, R: 0x3C,
				PC: 0x413D, SP: 0x0D16,
				IX: 0x2287, IY: 0x4255,
				AF_: 0xF647, BC_: 0xB54F,
				DE_: 0xA0AE, HL_: 0x7F32,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{16701, 8}},
			},
			want: z80State{
				A: 0xF6, F: 0x47,
				B: 0x4F, C: 0x06,
				D: 0x05, E: 0xEF,
				H: 0xE3, L: 0xB5,
				I: 0x1F, R: 0x3D,
				PC: 0x413E, SP: 0x0D16,
				IX: 0x2287, IY: 0x4255,
				AF_: 0xC71A, BC_: 0xB54F,
				DE_: 0xA0AE, HL_: 0x7F32,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{16701, 8}},
				Cycles: 4,
			},
		},
		{ // D9 0000
			name: "D9 0000",
			init: z80State{
				A: 0x1F, F: 0xEF,
				B: 0x5A, C: 0x9E,
				D: 0x88, E: 0x2B,
				H: 0x5E, L: 0x8A,
				I: 0xB6, R: 0x3D,
				PC: 0x22EC, SP: 0x92F0,
				IX: 0xC99E, IY: 0x9759,
				AF_: 0x0A7E, BC_: 0x7665,
				DE_: 0x52E4, HL_: 0xDE85,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{8940, 217}},
			},
			want: z80State{
				A: 0x1F, F: 0xEF,
				B: 0x76, C: 0x65,
				D: 0x52, E: 0xE4,
				H: 0xDE, L: 0x85,
				I: 0xB6, R: 0x3E,
				PC: 0x22ED, SP: 0x92F0,
				IX: 0xC99E, IY: 0x9759,
				AF_: 0x0A7E, BC_: 0x5A9E,
				DE_: 0x882B, HL_: 0x5E8A,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{8940, 217}},
				Cycles: 4,
			},
		},
		{ // E3 0000
			name: "E3 0000",
			init: z80State{
				A: 0x83, F: 0xA8,
				B: 0x50, C: 0xA2,
				D: 0x7B, E: 0x6A,
				H: 0xD7, L: 0x49,
				I: 0x2D, R: 0x1E,
				PC: 0x69CA, SP: 0x6CFD,
				IX: 0xF92A, IY: 0x93FA,
				AF_: 0x36C4, BC_: 0x8BBA,
				DE_: 0x885C, HL_: 0xD420,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{27082, 227}, {27901, 228}, {27902, 232}},
			},
			want: z80State{
				A: 0x83, F: 0xA8,
				B: 0x50, C: 0xA2,
				D: 0x7B, E: 0x6A,
				H: 0xE8, L: 0xE4,
				I: 0x2D, R: 0x1F,
				PC: 0x69CB, SP: 0x6CFD,
				IX: 0xF92A, IY: 0x93FA,
				AF_: 0x36C4, BC_: 0x8BBA,
				DE_: 0x885C, HL_: 0xD420,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{27082, 227}, {27901, 73}, {27902, 215}},
				Cycles: 19,
			},
		},
		{ // 02 0000
			name: "02 0000",
			init: z80State{
				A: 0xA2, F: 0x9E,
				B: 0x8A, C: 0x1E,
				D: 0xFA, E: 0xFE,
				H: 0x04, L: 0xB5,
				I: 0xDF, R: 0x0F,
				PC: 0x459A, SP: 0x3EFB,
				IX: 0xCC5D, IY: 0xD5CA,
				AF_: 0x031F, BC_: 0x6274,
				DE_: 0xD3A0, HL_: 0x9ECD,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{17818, 2}, {35358, 0}},
			},
			want: z80State{
				A: 0xA2, F: 0x9E,
				B: 0x8A, C: 0x1E,
				D: 0xFA, E: 0xFE,
				H: 0x04, L: 0xB5,
				I: 0xDF, R: 0x10,
				PC: 0x459B, SP: 0x3EFB,
				IX: 0xCC5D, IY: 0xD5CA,
				AF_: 0x031F, BC_: 0x6274,
				DE_: 0xD3A0, HL_: 0x9ECD,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{17818, 2}, {35358, 162}},
				Cycles: 7,
			},
		},
		{ // 1A 0000
			name: "1A 0000",
			init: z80State{
				A: 0xF8, F: 0x95,
				B: 0x31, C: 0x79,
				D: 0x17, E: 0xB7,
				H: 0xE2, L: 0x4F,
				I: 0xE0, R: 0x48,
				PC: 0x7BFB, SP: 0x59FC,
				IX: 0xDDC1, IY: 0x8C60,
				AF_: 0xD94C, BC_: 0x4CF2,
				DE_: 0x97EF, HL_: 0xBEB9,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{6071, 166}, {31739, 26}},
			},
			want: z80State{
				A: 0xA6, F: 0x95,
				B: 0x31, C: 0x79,
				D: 0x17, E: 0xB7,
				H: 0xE2, L: 0x4F,
				I: 0xE0, R: 0x49,
				PC: 0x7BFC, SP: 0x59FC,
				IX: 0xDDC1, IY: 0x8C60,
				AF_: 0xD94C, BC_: 0x4CF2,
				DE_: 0x97EF, HL_: 0xBEB9,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{6071, 166}, {31739, 26}},
				Cycles: 7,
			},
		},
		{ // 11 0000
			name: "11 0000",
			init: z80State{
				A: 0x4C, F: 0xB9,
				B: 0x7A, C: 0x54,
				D: 0xEB, E: 0xD9,
				H: 0xB7, L: 0xAA,
				I: 0x9C, R: 0x64,
				PC: 0x0C55, SP: 0xE875,
				IX: 0x698A, IY: 0x241E,
				AF_: 0xA5B1, BC_: 0x43A7,
				DE_: 0x456D, HL_: 0x4E36,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{3157, 17}, {3158, 182}, {3159, 24}},
			},
			want: z80State{
				A: 0x4C, F: 0xB9,
				B: 0x7A, C: 0x54,
				D: 0x18, E: 0xB6,
				H: 0xB7, L: 0xAA,
				I: 0x9C, R: 0x65,
				PC: 0x0C58, SP: 0xE875,
				IX: 0x698A, IY: 0x241E,
				AF_: 0xA5B1, BC_: 0x43A7,
				DE_: 0x456D, HL_: 0x4E36,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{3157, 17}, {3158, 182}, {3159, 24}},
				Cycles: 10,
			},
		},
		{ // 21 0000
			name: "21 0000",
			init: z80State{
				A: 0xCF, F: 0x0A,
				B: 0x4D, C: 0x34,
				D: 0x56, E: 0x20,
				H: 0xD9, L: 0x25,
				I: 0xCF, R: 0x74,
				PC: 0x1F5F, SP: 0x6BE5,
				IX: 0x08F9, IY: 0x489F,
				AF_: 0x6731, BC_: 0x20CD,
				DE_: 0x5617, HL_: 0xE1C8,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{8031, 33}, {8032, 62}, {8033, 143}},
			},
			want: z80State{
				A: 0xCF, F: 0x0A,
				B: 0x4D, C: 0x34,
				D: 0x56, E: 0x20,
				H: 0x8F, L: 0x3E,
				I: 0xCF, R: 0x75,
				PC: 0x1F62, SP: 0x6BE5,
				IX: 0x08F9, IY: 0x489F,
				AF_: 0x6731, BC_: 0x20CD,
				DE_: 0x5617, HL_: 0xE1C8,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{8031, 33}, {8032, 62}, {8033, 143}},
				Cycles: 10,
			},
		},
		{ // C1 0000
			name: "C1 0000",
			init: z80State{
				A: 0x3D, F: 0x00,
				B: 0xD6, C: 0x78,
				D: 0xE4, E: 0x3B,
				H: 0x80, L: 0x98,
				I: 0x2F, R: 0x7F,
				PC: 0x2A19, SP: 0x4288,
				IX: 0x8BCC, IY: 0x1C6D,
				AF_: 0x3089, BC_: 0x3910,
				DE_: 0x462C, HL_: 0x77E2,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{10777, 193}, {17032, 242}, {17033, 175}},
			},
			want: z80State{
				A: 0x3D, F: 0x00,
				B: 0xAF, C: 0xF2,
				D: 0xE4, E: 0x3B,
				H: 0x80, L: 0x98,
				I: 0x2F, R: 0x00,
				PC: 0x2A1A, SP: 0x428A,
				IX: 0x8BCC, IY: 0x1C6D,
				AF_: 0x3089, BC_: 0x3910,
				DE_: 0x462C, HL_: 0x77E2,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{10777, 193}, {17032, 242}, {17033, 175}},
				Cycles: 10,
			},
		},
		{ // D5 0000
			name: "D5 0000",
			init: z80State{
				A: 0x55, F: 0x78,
				B: 0x47, C: 0xD4,
				D: 0x3A, E: 0x70,
				H: 0xBB, L: 0x6A,
				I: 0x90, R: 0x4B,
				PC: 0x9210, SP: 0x4449,
				IX: 0x585D, IY: 0xE062,
				AF_: 0x53B8, BC_: 0x435A,
				DE_: 0x7405, HL_: 0x5A33,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{17479, 0}, {17480, 0}, {37392, 213}},
			},
			want: z80State{
				A: 0x55, F: 0x78,
				B: 0x47, C: 0xD4,
				D: 0x3A, E: 0x70,
				H: 0xBB, L: 0x6A,
				I: 0x90, R: 0x4C,
				PC: 0x9211, SP: 0x4447,
				IX: 0x585D, IY: 0xE062,
				AF_: 0x53B8, BC_: 0x435A,
				DE_: 0x7405, HL_: 0x5A33,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{17479, 112}, {17480, 58}, {37392, 213}},
				Cycles: 11,
			},
		},
		{ // F3 0000
			name: "F3 0000",
			init: z80State{
				A: 0xDD, F: 0xFF,
				B: 0x4F, C: 0x82,
				D: 0x0A, E: 0xF0,
				H: 0x30, L: 0x8D,
				I: 0x7C, R: 0x32,
				PC: 0x83DF, SP: 0x4E88,
				IX: 0x21A9, IY: 0xA1A3,
				AF_: 0xA672, BC_: 0x4B37,
				DE_: 0xE9ED, HL_: 0x866F,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{33759, 243}},
			},
			want: z80State{
				A: 0xDD, F: 0xFF,
				B: 0x4F, C: 0x82,
				D: 0x0A, E: 0xF0,
				H: 0x30, L: 0x8D,
				I: 0x7C, R: 0x33,
				PC: 0x83E0, SP: 0x4E88,
				IX: 0x21A9, IY: 0xA1A3,
				AF_: 0xA672, BC_: 0x4B37,
				DE_: 0xE9ED, HL_: 0x866F,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{33759, 243}},
				Cycles: 4,
			},
		},
		{ // FB 0000
			name: "FB 0000",
			init: z80State{
				A: 0x1F, F: 0x6B,
				B: 0x25, C: 0x0D,
				D: 0xFC, E: 0x27,
				H: 0xD3, L: 0x33,
				I: 0x24, R: 0x5F,
				PC: 0x96B9, SP: 0xFBA3,
				IX: 0xB2A7, IY: 0x66AB,
				AF_: 0xCDD2, BC_: 0x9519,
				DE_: 0xE235, HL_: 0x9A54,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{38585, 251}},
			},
			want: z80State{
				A: 0x1F, F: 0x6B,
				B: 0x25, C: 0x0D,
				D: 0xFC, E: 0x27,
				H: 0xD3, L: 0x33,
				I: 0x24, R: 0x60,
				PC: 0x96BA, SP: 0xFBA3,
				IX: 0xB2A7, IY: 0x66AB,
				AF_: 0xCDD2, BC_: 0x9519,
				DE_: 0xE235, HL_: 0x9A54,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{38585, 251}},
				Cycles: 4,
			},
		},
	}
	for _, tt := range tests {
		runSSTTest(t, tt)
	}
}
