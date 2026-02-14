package z80

import "testing"

func TestJP_nn(t *testing.T) {
	cpu, bus := newTestCPU()
	// JP 0x1234 (C3 34 12)
	bus.mem[0] = 0xC3
	bus.mem[1] = 0x34
	bus.mem[2] = 0x12
	cycles := cpu.Step()
	if cpu.reg.PC != 0x1234 {
		t.Errorf("JP nn: PC=%04x want 1234", cpu.reg.PC)
	}
	if cycles != 10 {
		t.Errorf("cycles=%d want 10", cycles)
	}
}

func TestJP_cc_nn_Taken(t *testing.T) {
	cpu, bus := newTestCPU()
	// JP Z, 0x5678 (CA 78 56) with Z flag set
	bus.mem[0] = 0xCA
	bus.mem[1] = 0x78
	bus.mem[2] = 0x56
	cpu.reg.AF = 0x0040 // Z flag set
	cycles := cpu.Step()
	if cpu.reg.PC != 0x5678 {
		t.Errorf("JP Z taken: PC=%04x want 5678", cpu.reg.PC)
	}
	if cycles != 10 {
		t.Errorf("cycles=%d want 10", cycles)
	}
}

func TestJP_cc_nn_NotTaken(t *testing.T) {
	cpu, bus := newTestCPU()
	// JP Z, 0x5678 with Z flag clear
	bus.mem[0] = 0xCA
	bus.mem[1] = 0x78
	bus.mem[2] = 0x56
	cpu.reg.AF = 0x0000 // Z flag clear
	cycles := cpu.Step()
	if cpu.reg.PC != 3 { // Past the 3-byte instruction
		t.Errorf("JP Z not taken: PC=%04x want 0003", cpu.reg.PC)
	}
	if cycles != 10 {
		t.Errorf("cycles=%d want 10", cycles)
	}
}

func TestJR_e(t *testing.T) {
	cpu, bus := newTestCPU()
	// JR +5 (18 05) at PC=0x1000
	cpu.reg.PC = 0x1000
	bus.mem[0x1000] = 0x18
	bus.mem[0x1001] = 0x05
	cycles := cpu.Step()
	// PC after fetch = 0x1002, then +5 = 0x1007
	if cpu.reg.PC != 0x1007 {
		t.Errorf("JR: PC=%04x want 1007", cpu.reg.PC)
	}
	if cycles != 12 {
		t.Errorf("cycles=%d want 12", cycles)
	}
}

func TestJR_e_Backward(t *testing.T) {
	cpu, bus := newTestCPU()
	// JR -3 (18 FD) at PC=0x1000
	cpu.reg.PC = 0x1000
	bus.mem[0x1000] = 0x18
	bus.mem[0x1001] = 0xFD // -3
	cpu.Step()
	// PC after fetch = 0x1002, then -3 = 0x0FFF
	if cpu.reg.PC != 0x0FFF {
		t.Errorf("JR back: PC=%04x want 0FFF", cpu.reg.PC)
	}
}

func TestJR_NZ_Taken(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x20 // JR NZ
	bus.mem[1] = 0x10
	cpu.reg.AF = 0x0000 // Z clear
	cycles := cpu.Step()
	if cpu.reg.PC != 0x0012 {
		t.Errorf("JR NZ taken: PC=%04x want 0012", cpu.reg.PC)
	}
	if cycles != 12 {
		t.Errorf("cycles=%d want 12", cycles)
	}
}

func TestJR_NZ_NotTaken(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x20 // JR NZ
	bus.mem[1] = 0x10
	cpu.reg.AF = 0x0040 // Z set
	cycles := cpu.Step()
	if cpu.reg.PC != 2 {
		t.Errorf("JR NZ not taken: PC=%04x want 0002", cpu.reg.PC)
	}
	if cycles != 7 {
		t.Errorf("cycles=%d want 7", cycles)
	}
}

func TestJP_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xE9 // JP (HL)
	cpu.reg.HL = 0xABCD
	cycles := cpu.Step()
	if cpu.reg.PC != 0xABCD {
		t.Errorf("JP (HL): PC=%04x want ABCD", cpu.reg.PC)
	}
	if cycles != 4 {
		t.Errorf("cycles=%d want 4", cycles)
	}
}

func TestDJNZ_Taken(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x10 // DJNZ
	bus.mem[1] = 0x05
	cpu.reg.BC = 0x0200 // B=2
	cycles := cpu.Step()
	if cpu.getB() != 1 {
		t.Errorf("DJNZ: B=%02x want 01", cpu.getB())
	}
	if cpu.reg.PC != 0x0007 {
		t.Errorf("DJNZ taken: PC=%04x want 0007", cpu.reg.PC)
	}
	if cycles != 13 {
		t.Errorf("cycles=%d want 13", cycles)
	}
}

func TestDJNZ_NotTaken(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x10
	bus.mem[1] = 0x05
	cpu.reg.BC = 0x0100 // B=1
	cycles := cpu.Step()
	if cpu.getB() != 0 {
		t.Errorf("DJNZ: B=%02x want 00", cpu.getB())
	}
	if cpu.reg.PC != 2 {
		t.Errorf("DJNZ not taken: PC=%04x want 0002", cpu.reg.PC)
	}
	if cycles != 8 {
		t.Errorf("cycles=%d want 8", cycles)
	}
}

func TestCALL_nn(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFE
	bus.mem[0] = 0xCD
	bus.mem[1] = 0x00
	bus.mem[2] = 0x50
	cycles := cpu.Step()
	if cpu.reg.PC != 0x5000 {
		t.Errorf("CALL: PC=%04x want 5000", cpu.reg.PC)
	}
	if cycles != 17 {
		t.Errorf("cycles=%d want 17", cycles)
	}
	// Return address (3) should be on stack
	retAddr := cpu.read16(cpu.reg.SP)
	if retAddr != 3 {
		t.Errorf("return addr=%04x want 0003", retAddr)
	}
}

func TestCALL_cc_Taken(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFE
	// CALL Z, 0x6000 (CC 00 60)
	bus.mem[0] = 0xCC
	bus.mem[1] = 0x00
	bus.mem[2] = 0x60
	cpu.reg.AF = 0x0040 // Z set
	cycles := cpu.Step()
	if cpu.reg.PC != 0x6000 {
		t.Errorf("CALL Z taken: PC=%04x want 6000", cpu.reg.PC)
	}
	if cycles != 17 {
		t.Errorf("cycles=%d want 17", cycles)
	}
}

func TestCALL_cc_NotTaken(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFE
	bus.mem[0] = 0xCC
	bus.mem[1] = 0x00
	bus.mem[2] = 0x60
	cpu.reg.AF = 0x0000 // Z clear
	cycles := cpu.Step()
	if cpu.reg.PC != 3 {
		t.Errorf("CALL Z not taken: PC=%04x want 0003", cpu.reg.PC)
	}
	if cycles != 10 {
		t.Errorf("cycles=%d want 10", cycles)
	}
}

func TestRET(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFC
	bus.mem[0xFFFC] = 0x34
	bus.mem[0xFFFD] = 0x12
	bus.mem[0] = 0xC9 // RET
	cycles := cpu.Step()
	if cpu.reg.PC != 0x1234 {
		t.Errorf("RET: PC=%04x want 1234", cpu.reg.PC)
	}
	if cycles != 10 {
		t.Errorf("cycles=%d want 10", cycles)
	}
}

func TestRET_cc_Taken(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFC
	bus.mem[0xFFFC] = 0x78
	bus.mem[0xFFFD] = 0x56
	bus.mem[0] = 0xC8 // RET Z
	cpu.reg.AF = 0x0040 // Z set
	cycles := cpu.Step()
	if cpu.reg.PC != 0x5678 {
		t.Errorf("RET Z taken: PC=%04x want 5678", cpu.reg.PC)
	}
	if cycles != 11 {
		t.Errorf("cycles=%d want 11", cycles)
	}
}

func TestRET_cc_NotTaken(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xC8 // RET Z
	cpu.reg.AF = 0x0000 // Z clear
	cycles := cpu.Step()
	if cpu.reg.PC != 1 {
		t.Errorf("RET Z not taken: PC=%04x want 0001", cpu.reg.PC)
	}
	if cycles != 5 {
		t.Errorf("cycles=%d want 5", cycles)
	}
}

func TestRST(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFE
	cpu.reg.PC = 0x1000
	bus.mem[0x1000] = 0xDF // RST 18h
	cycles := cpu.Step()
	if cpu.reg.PC != 0x0018 {
		t.Errorf("RST 18h: PC=%04x want 0018", cpu.reg.PC)
	}
	if cycles != 11 {
		t.Errorf("cycles=%d want 11", cycles)
	}
	retAddr := cpu.read16(cpu.reg.SP)
	if retAddr != 0x1001 {
		t.Errorf("return addr=%04x want 1001", retAddr)
	}
}

func TestSST_Branch(t *testing.T) {
	tests := []struct {
		name       string
		init, want z80State
	}{
		{ // C3 0000
			name: "C3 0000",
			init: z80State{
				A: 0xC1, F: 0xD4,
				B: 0xDF, C: 0x13,
				D: 0x0E, E: 0x7F,
				H: 0xB1, L: 0x3E,
				I: 0xE5, R: 0x65,
				PC: 0x0B56, SP: 0xA40A,
				IX: 0x448A, IY: 0xD77A,
				AF_: 0xFBEB, BC_: 0x8A99,
				DE_: 0xFEC3, HL_: 0x8269,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{2902, 195}, {2903, 166}, {2904, 133}},
			},
			want: z80State{
				A: 0xC1, F: 0xD4,
				B: 0xDF, C: 0x13,
				D: 0x0E, E: 0x7F,
				H: 0xB1, L: 0x3E,
				I: 0xE5, R: 0x66,
				PC: 0x85A6, SP: 0xA40A,
				IX: 0x448A, IY: 0xD77A,
				AF_: 0xFBEB, BC_: 0x8A99,
				DE_: 0xFEC3, HL_: 0x8269,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{2902, 195}, {2903, 166}, {2904, 133}},
				Cycles: 10,
			},
		},
		{ // C2 0000
			name: "C2 0000",
			init: z80State{
				A: 0xFA, F: 0xB7,
				B: 0xA5, C: 0x57,
				D: 0x3B, E: 0x14,
				H: 0x31, L: 0x65,
				I: 0xAD, R: 0x0E,
				PC: 0x0FEC, SP: 0x8B9D,
				IX: 0xC4E4, IY: 0xCE66,
				AF_: 0xE9E3, BC_: 0x8E5E,
				DE_: 0xE476, HL_: 0x74D4,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{4076, 194}, {4077, 135}, {4078, 241}},
			},
			want: z80State{
				A: 0xFA, F: 0xB7,
				B: 0xA5, C: 0x57,
				D: 0x3B, E: 0x14,
				H: 0x31, L: 0x65,
				I: 0xAD, R: 0x0F,
				PC: 0xF187, SP: 0x8B9D,
				IX: 0xC4E4, IY: 0xCE66,
				AF_: 0xE9E3, BC_: 0x8E5E,
				DE_: 0xE476, HL_: 0x74D4,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{4076, 194}, {4077, 135}, {4078, 241}},
				Cycles: 10,
			},
		},
		{ // CA 0000
			name: "CA 0000",
			init: z80State{
				A: 0xA3, F: 0x64,
				B: 0xFE, C: 0xE7,
				D: 0xE1, E: 0x1D,
				H: 0xE5, L: 0xC1,
				I: 0xF9, R: 0x65,
				PC: 0x93D4, SP: 0x35F0,
				IX: 0xEABB, IY: 0x4552,
				AF_: 0xBC77, BC_: 0xE723,
				DE_: 0xCE9F, HL_: 0xD8BF,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{37844, 202}, {37845, 21}, {37846, 62}},
			},
			want: z80State{
				A: 0xA3, F: 0x64,
				B: 0xFE, C: 0xE7,
				D: 0xE1, E: 0x1D,
				H: 0xE5, L: 0xC1,
				I: 0xF9, R: 0x66,
				PC: 0x3E15, SP: 0x35F0,
				IX: 0xEABB, IY: 0x4552,
				AF_: 0xBC77, BC_: 0xE723,
				DE_: 0xCE9F, HL_: 0xD8BF,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{37844, 202}, {37845, 21}, {37846, 62}},
				Cycles: 10,
			},
		},
		{ // 18 0000
			name: "18 0000",
			init: z80State{
				A: 0x37, F: 0x69,
				B: 0x60, C: 0xFB,
				D: 0xFC, E: 0x93,
				H: 0xA0, L: 0x6C,
				I: 0xEC, R: 0x15,
				PC: 0x5F0F, SP: 0xF363,
				IX: 0xC9A1, IY: 0xC0DF,
				AF_: 0xBE26, BC_: 0x488C,
				DE_: 0x5E06, HL_: 0x9A4D,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{24335, 24}, {24336, 53}},
			},
			want: z80State{
				A: 0x37, F: 0x69,
				B: 0x60, C: 0xFB,
				D: 0xFC, E: 0x93,
				H: 0xA0, L: 0x6C,
				I: 0xEC, R: 0x16,
				PC: 0x5F46, SP: 0xF363,
				IX: 0xC9A1, IY: 0xC0DF,
				AF_: 0xBE26, BC_: 0x488C,
				DE_: 0x5E06, HL_: 0x9A4D,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{24335, 24}, {24336, 53}},
				Cycles: 12,
			},
		},
		{ // 20 0000
			name: "20 0000",
			init: z80State{
				A: 0xCA, F: 0x27,
				B: 0x52, C: 0xE9,
				D: 0x3F, E: 0x9C,
				H: 0x0F, L: 0xAB,
				I: 0x5F, R: 0x67,
				PC: 0xCA5B, SP: 0xB6F8,
				IX: 0xC96A, IY: 0x3B26,
				AF_: 0x7DD5, BC_: 0xAE8E,
				DE_: 0xF25E, HL_: 0x264C,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{51803, 32}, {51804, 232}},
			},
			want: z80State{
				A: 0xCA, F: 0x27,
				B: 0x52, C: 0xE9,
				D: 0x3F, E: 0x9C,
				H: 0x0F, L: 0xAB,
				I: 0x5F, R: 0x68,
				PC: 0xCA45, SP: 0xB6F8,
				IX: 0xC96A, IY: 0x3B26,
				AF_: 0x7DD5, BC_: 0xAE8E,
				DE_: 0xF25E, HL_: 0x264C,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{51803, 32}, {51804, 232}},
				Cycles: 12,
			},
		},
		{ // 28 0000
			name: "28 0000",
			init: z80State{
				A: 0xC0, F: 0x84,
				B: 0xAD, C: 0xCB,
				D: 0x0C, E: 0x3C,
				H: 0x52, L: 0x0D,
				I: 0x0F, R: 0x4D,
				PC: 0x21D8, SP: 0x2EC8,
				IX: 0xD0E9, IY: 0xDF02,
				AF_: 0x55AB, BC_: 0x2BE1,
				DE_: 0x73C2, HL_: 0xC005,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{8664, 40}, {8665, 244}},
			},
			want: z80State{
				A: 0xC0, F: 0x84,
				B: 0xAD, C: 0xCB,
				D: 0x0C, E: 0x3C,
				H: 0x52, L: 0x0D,
				I: 0x0F, R: 0x4E,
				PC: 0x21DA, SP: 0x2EC8,
				IX: 0xD0E9, IY: 0xDF02,
				AF_: 0x55AB, BC_: 0x2BE1,
				DE_: 0x73C2, HL_: 0xC005,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{8664, 40}, {8665, 244}},
				Cycles: 7,
			},
		},
		{ // 38 0000
			name: "38 0000",
			init: z80State{
				A: 0x0E, F: 0xA6,
				B: 0x03, C: 0x91,
				D: 0xB4, E: 0xAF,
				H: 0x50, L: 0x8E,
				I: 0x7C, R: 0x07,
				PC: 0x6CB3, SP: 0x98EB,
				IX: 0xD631, IY: 0x5D0A,
				AF_: 0x5BEF, BC_: 0xAD53,
				DE_: 0xC7E1, HL_: 0x9355,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{27827, 56}, {27828, 50}},
			},
			want: z80State{
				A: 0x0E, F: 0xA6,
				B: 0x03, C: 0x91,
				D: 0xB4, E: 0xAF,
				H: 0x50, L: 0x8E,
				I: 0x7C, R: 0x08,
				PC: 0x6CB5, SP: 0x98EB,
				IX: 0xD631, IY: 0x5D0A,
				AF_: 0x5BEF, BC_: 0xAD53,
				DE_: 0xC7E1, HL_: 0x9355,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{27827, 56}, {27828, 50}},
				Cycles: 7,
			},
		},
		{ // E9 0000
			name: "E9 0000",
			init: z80State{
				A: 0x0E, F: 0x69,
				B: 0xD3, C: 0x8C,
				D: 0xA1, E: 0x39,
				H: 0xBC, L: 0x84,
				I: 0x6E, R: 0x1D,
				PC: 0xE49E, SP: 0xE08A,
				IX: 0xB107, IY: 0xF670,
				AF_: 0x1562, BC_: 0xED6B,
				DE_: 0x9B03, HL_: 0xD07F,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{58526, 233}},
			},
			want: z80State{
				A: 0x0E, F: 0x69,
				B: 0xD3, C: 0x8C,
				D: 0xA1, E: 0x39,
				H: 0xBC, L: 0x84,
				I: 0x6E, R: 0x1E,
				PC: 0xBC84, SP: 0xE08A,
				IX: 0xB107, IY: 0xF670,
				AF_: 0x1562, BC_: 0xED6B,
				DE_: 0x9B03, HL_: 0xD07F,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{58526, 233}},
				Cycles: 4,
			},
		},
		{ // 10 0000
			name: "10 0000",
			init: z80State{
				A: 0xD1, F: 0x53,
				B: 0x74, C: 0xAD,
				D: 0xCA, E: 0x34,
				H: 0x5D, L: 0x13,
				I: 0x50, R: 0x06,
				PC: 0xEB22, SP: 0x5A32,
				IX: 0x71BE, IY: 0x3CB1,
				AF_: 0x0FAC, BC_: 0x4C44,
				DE_: 0x9032, HL_: 0x045E,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{60194, 16}, {60195, 194}},
			},
			want: z80State{
				A: 0xD1, F: 0x53,
				B: 0x73, C: 0xAD,
				D: 0xCA, E: 0x34,
				H: 0x5D, L: 0x13,
				I: 0x50, R: 0x07,
				PC: 0xEAE6, SP: 0x5A32,
				IX: 0x71BE, IY: 0x3CB1,
				AF_: 0x0FAC, BC_: 0x4C44,
				DE_: 0x9032, HL_: 0x045E,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{60194, 16}, {60195, 194}},
				Cycles: 13,
			},
		},
		{ // CD 0000
			name: "CD 0000",
			init: z80State{
				A: 0xBE, F: 0x45,
				B: 0xEF, C: 0xD9,
				D: 0xC4, E: 0x47,
				H: 0x9B, L: 0xA6,
				I: 0xF6, R: 0x29,
				PC: 0x7FB5, SP: 0x4671,
				IX: 0x3FC3, IY: 0x455D,
				AF_: 0x90B0, BC_: 0x0DCE,
				DE_: 0x05B2, HL_: 0xD049,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{18031, 0}, {18032, 0}, {32693, 205}, {32694, 156}, {32695, 62}},
			},
			want: z80State{
				A: 0xBE, F: 0x45,
				B: 0xEF, C: 0xD9,
				D: 0xC4, E: 0x47,
				H: 0x9B, L: 0xA6,
				I: 0xF6, R: 0x2A,
				PC: 0x3E9C, SP: 0x466F,
				IX: 0x3FC3, IY: 0x455D,
				AF_: 0x90B0, BC_: 0x0DCE,
				DE_: 0x05B2, HL_: 0xD049,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{18031, 184}, {18032, 127}, {32693, 205}, {32694, 156}, {32695, 62}},
				Cycles: 17,
			},
		},
		{ // C4 0000
			name: "C4 0000",
			init: z80State{
				A: 0xEC, F: 0x3F,
				B: 0xCE, C: 0x6A,
				D: 0xFE, E: 0x7E,
				H: 0x67, L: 0x85,
				I: 0x0E, R: 0x03,
				PC: 0xF416, SP: 0x5DFB,
				IX: 0x95B3, IY: 0x057F,
				AF_: 0xDB64, BC_: 0xC52A,
				DE_: 0xB80F, HL_: 0xA7D8,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{24057, 0}, {24058, 0}, {62486, 196}, {62487, 177}, {62488, 161}},
			},
			want: z80State{
				A: 0xEC, F: 0x3F,
				B: 0xCE, C: 0x6A,
				D: 0xFE, E: 0x7E,
				H: 0x67, L: 0x85,
				I: 0x0E, R: 0x04,
				PC: 0xA1B1, SP: 0x5DF9,
				IX: 0x95B3, IY: 0x057F,
				AF_: 0xDB64, BC_: 0xC52A,
				DE_: 0xB80F, HL_: 0xA7D8,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{24057, 25}, {24058, 244}, {62486, 196}, {62487, 177}, {62488, 161}},
				Cycles: 17,
			},
		},
		{ // CC 0000
			name: "CC 0000",
			init: z80State{
				A: 0x02, F: 0x20,
				B: 0x59, C: 0x3A,
				D: 0x8C, E: 0x61,
				H: 0xCD, L: 0x75,
				I: 0x35, R: 0x21,
				PC: 0x8A19, SP: 0xD84B,
				IX: 0xB6BA, IY: 0xE7F6,
				AF_: 0x1E76, BC_: 0xDC7F,
				DE_: 0x0666, HL_: 0x955B,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{35353, 204}, {35354, 48}, {35355, 112}},
			},
			want: z80State{
				A: 0x02, F: 0x20,
				B: 0x59, C: 0x3A,
				D: 0x8C, E: 0x61,
				H: 0xCD, L: 0x75,
				I: 0x35, R: 0x22,
				PC: 0x8A1C, SP: 0xD84B,
				IX: 0xB6BA, IY: 0xE7F6,
				AF_: 0x1E76, BC_: 0xDC7F,
				DE_: 0x0666, HL_: 0x955B,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{35353, 204}, {35354, 48}, {35355, 112}},
				Cycles: 10,
			},
		},
		{ // C9 0000
			name: "C9 0000",
			init: z80State{
				A: 0xB5, F: 0xDC,
				B: 0x30, C: 0xF8,
				D: 0x4C, E: 0xB8,
				H: 0xA9, L: 0x04,
				I: 0x3E, R: 0x6C,
				PC: 0xB7F5, SP: 0x9F4E,
				IX: 0x23CB, IY: 0x9743,
				AF_: 0x9638, BC_: 0xB4D2,
				DE_: 0x3021, HL_: 0x9B29,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{40782, 196}, {40783, 238}, {47093, 201}},
			},
			want: z80State{
				A: 0xB5, F: 0xDC,
				B: 0x30, C: 0xF8,
				D: 0x4C, E: 0xB8,
				H: 0xA9, L: 0x04,
				I: 0x3E, R: 0x6D,
				PC: 0xEEC4, SP: 0x9F50,
				IX: 0x23CB, IY: 0x9743,
				AF_: 0x9638, BC_: 0xB4D2,
				DE_: 0x3021, HL_: 0x9B29,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{40782, 196}, {40783, 238}, {47093, 201}},
				Cycles: 10,
			},
		},
		{ // C0 0000
			name: "C0 0000",
			init: z80State{
				A: 0x8E, F: 0xC8,
				B: 0xF9, C: 0x97,
				D: 0xFC, E: 0x10,
				H: 0xA0, L: 0x54,
				I: 0xA4, R: 0x67,
				PC: 0xA0C8, SP: 0x2F1E,
				IX: 0xBFC4, IY: 0xC20D,
				AF_: 0x7C19, BC_: 0xB32E,
				DE_: 0x8B83, HL_: 0xB00D,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{41160, 192}},
			},
			want: z80State{
				A: 0x8E, F: 0xC8,
				B: 0xF9, C: 0x97,
				D: 0xFC, E: 0x10,
				H: 0xA0, L: 0x54,
				I: 0xA4, R: 0x68,
				PC: 0xA0C9, SP: 0x2F1E,
				IX: 0xBFC4, IY: 0xC20D,
				AF_: 0x7C19, BC_: 0xB32E,
				DE_: 0x8B83, HL_: 0xB00D,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{41160, 192}},
				Cycles: 5,
			},
		},
		{ // C8 0000
			name: "C8 0000",
			init: z80State{
				A: 0x32, F: 0x43,
				B: 0x0C, C: 0x5E,
				D: 0x8B, E: 0x8B,
				H: 0x49, L: 0xF0,
				I: 0x93, R: 0x14,
				PC: 0xCD0B, SP: 0x0F0D,
				IX: 0xC5D9, IY: 0x38C5,
				AF_: 0xB84E, BC_: 0xC3BE,
				DE_: 0x1E9D, HL_: 0x6403,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{3853, 55}, {3854, 51}, {52491, 200}},
			},
			want: z80State{
				A: 0x32, F: 0x43,
				B: 0x0C, C: 0x5E,
				D: 0x8B, E: 0x8B,
				H: 0x49, L: 0xF0,
				I: 0x93, R: 0x15,
				PC: 0x3337, SP: 0x0F0F,
				IX: 0xC5D9, IY: 0x38C5,
				AF_: 0xB84E, BC_: 0xC3BE,
				DE_: 0x1E9D, HL_: 0x6403,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{3853, 55}, {3854, 51}, {52491, 200}},
				Cycles: 11,
			},
		},
		{ // C7 0000
			name: "C7 0000",
			init: z80State{
				A: 0x1C, F: 0x80,
				B: 0xDD, C: 0xF8,
				D: 0x45, E: 0x23,
				H: 0xCD, L: 0x45,
				I: 0xA6, R: 0x2B,
				PC: 0x62F8, SP: 0x5C41,
				IX: 0xA417, IY: 0x9073,
				AF_: 0x4469, BC_: 0x5D6C,
				DE_: 0x230F, HL_: 0x5F2B,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{23615, 0}, {23616, 0}, {25336, 199}},
			},
			want: z80State{
				A: 0x1C, F: 0x80,
				B: 0xDD, C: 0xF8,
				D: 0x45, E: 0x23,
				H: 0xCD, L: 0x45,
				I: 0xA6, R: 0x2C,
				PC: 0x0000, SP: 0x5C3F,
				IX: 0xA417, IY: 0x9073,
				AF_: 0x4469, BC_: 0x5D6C,
				DE_: 0x230F, HL_: 0x5F2B,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{23615, 249}, {23616, 98}, {25336, 199}},
				Cycles: 11,
			},
		},
		{ // CF 0000
			name: "CF 0000",
			init: z80State{
				A: 0x31, F: 0x8D,
				B: 0xF6, C: 0x46,
				D: 0x28, E: 0x99,
				H: 0x9B, L: 0xEF,
				I: 0x24, R: 0x6A,
				PC: 0x59DD, SP: 0xEC30,
				IX: 0xE099, IY: 0xDF51,
				AF_: 0xFFAE, BC_: 0x72FF,
				DE_: 0x03EB, HL_: 0xAE26,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{23005, 207}, {60462, 0}, {60463, 0}},
			},
			want: z80State{
				A: 0x31, F: 0x8D,
				B: 0xF6, C: 0x46,
				D: 0x28, E: 0x99,
				H: 0x9B, L: 0xEF,
				I: 0x24, R: 0x6B,
				PC: 0x0008, SP: 0xEC2E,
				IX: 0xE099, IY: 0xDF51,
				AF_: 0xFFAE, BC_: 0x72FF,
				DE_: 0x03EB, HL_: 0xAE26,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{23005, 207}, {60462, 222}, {60463, 89}},
				Cycles: 11,
			},
		},
		{ // FF 0000
			name: "FF 0000",
			init: z80State{
				A: 0x78, F: 0xE4,
				B: 0x11, C: 0xA6,
				D: 0xBA, E: 0x3F,
				H: 0xFB, L: 0x09,
				I: 0x09, R: 0x34,
				PC: 0xF548, SP: 0x558B,
				IX: 0x6704, IY: 0xA73D,
				AF_: 0x3F54, BC_: 0x0BF1,
				DE_: 0x647E, HL_: 0xA287,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{21897, 0}, {21898, 0}, {62792, 255}},
			},
			want: z80State{
				A: 0x78, F: 0xE4,
				B: 0x11, C: 0xA6,
				D: 0xBA, E: 0x3F,
				H: 0xFB, L: 0x09,
				I: 0x09, R: 0x35,
				PC: 0x0038, SP: 0x5589,
				IX: 0x6704, IY: 0xA73D,
				AF_: 0x3F54, BC_: 0x0BF1,
				DE_: 0x647E, HL_: 0xA287,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{21897, 73}, {21898, 245}, {62792, 255}},
				Cycles: 11,
			},
		},
		{ // D2 0000 — JP NC,nn (NC not met → not taken)
			name: "D2 0000",
			init: z80State{
				A: 0xD4, F: 0x59,
				B: 0x99, C: 0x21,
				D: 0x26, E: 0x92,
				H: 0x75, L: 0xAF,
				I: 0x4B, R: 0x7B,
				PC: 0xE276, SP: 0x24DA,
				IX: 0x1FAF, IY: 0x7E62,
				AF_: 0xADCE, BC_: 0x82E4,
				DE_: 0x430B, HL_: 0xCA8D,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{57974, 210}, {57975, 154}, {57976, 166}},
			},
			want: z80State{
				A: 0xD4, F: 0x59,
				B: 0x99, C: 0x21,
				D: 0x26, E: 0x92,
				H: 0x75, L: 0xAF,
				I: 0x4B, R: 0x7C,
				PC: 0xE279, SP: 0x24DA,
				IX: 0x1FAF, IY: 0x7E62,
				AF_: 0xADCE, BC_: 0x82E4,
				DE_: 0x430B, HL_: 0xCA8D,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{57974, 210}, {57975, 154}, {57976, 166}},
				Cycles: 10,
			},
		},
		{ // D0 0000 — RET NC (NC not met → not taken)
			name: "D0 0000",
			init: z80State{
				A: 0x4E, F: 0x87,
				B: 0xBC, C: 0x97,
				D: 0xED, E: 0xBC,
				H: 0x4B, L: 0x69,
				I: 0x24, R: 0x0F,
				PC: 0x5FD9, SP: 0x108A,
				IX: 0x0541, IY: 0x284C,
				AF_: 0xE133, BC_: 0x8695,
				DE_: 0x9108, HL_: 0xE458,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{24537, 208}},
			},
			want: z80State{
				A: 0x4E, F: 0x87,
				B: 0xBC, C: 0x97,
				D: 0xED, E: 0xBC,
				H: 0x4B, L: 0x69,
				I: 0x24, R: 0x10,
				PC: 0x5FDA, SP: 0x108A,
				IX: 0x0541, IY: 0x284C,
				AF_: 0xE133, BC_: 0x8695,
				DE_: 0x9108, HL_: 0xE458,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{24537, 208}},
				Cycles: 5,
			},
		},
		{ // E2 0000 — JP PO,nn (PO not met → not taken)
			name: "E2 0000",
			init: z80State{
				A: 0xB5, F: 0xF6,
				B: 0x3F, C: 0x1E,
				D: 0xDA, E: 0x41,
				H: 0x6B, L: 0x98,
				I: 0xB7, R: 0x04,
				PC: 0x067E, SP: 0x19DC,
				IX: 0x0039, IY: 0xC3C7,
				AF_: 0x9156, BC_: 0xC44C,
				DE_: 0x44B2, HL_: 0xF0A2,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{1662, 226}, {1663, 25}, {1664, 95}},
			},
			want: z80State{
				A: 0xB5, F: 0xF6,
				B: 0x3F, C: 0x1E,
				D: 0xDA, E: 0x41,
				H: 0x6B, L: 0x98,
				I: 0xB7, R: 0x05,
				PC: 0x0681, SP: 0x19DC,
				IX: 0x0039, IY: 0xC3C7,
				AF_: 0x9156, BC_: 0xC44C,
				DE_: 0x44B2, HL_: 0xF0A2,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{1662, 226}, {1663, 25}, {1664, 95}},
				Cycles: 10,
			},
		},
		{ // E4 0000 — CALL PO,nn (PO not met → not taken)
			name: "E4 0000",
			init: z80State{
				A: 0x84, F: 0xBF,
				B: 0xAD, C: 0xBD,
				D: 0x67, E: 0x66,
				H: 0x3C, L: 0xBD,
				I: 0x3C, R: 0x52,
				PC: 0x2CC3, SP: 0xA90E,
				IX: 0xB6CF, IY: 0x99FD,
				AF_: 0x04EC, BC_: 0x011B,
				DE_: 0x8B1F, HL_: 0x0E54,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{11459, 228}, {11460, 21}, {11461, 135}},
			},
			want: z80State{
				A: 0x84, F: 0xBF,
				B: 0xAD, C: 0xBD,
				D: 0x67, E: 0x66,
				H: 0x3C, L: 0xBD,
				I: 0x3C, R: 0x53,
				PC: 0x2CC6, SP: 0xA90E,
				IX: 0xB6CF, IY: 0x99FD,
				AF_: 0x04EC, BC_: 0x011B,
				DE_: 0x8B1F, HL_: 0x0E54,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{11459, 228}, {11460, 21}, {11461, 135}},
				Cycles: 10,
			},
		},
		{ // EC 0000 — CALL PE,nn (PE not met → not taken)
			name: "EC 0000",
			init: z80State{
				A: 0x98, F: 0x9A,
				B: 0x93, C: 0x99,
				D: 0x4C, E: 0xAF,
				H: 0x5A, L: 0x50,
				I: 0xA6, R: 0x7B,
				PC: 0x3FE9, SP: 0x806F,
				IX: 0xA291, IY: 0x28A6,
				AF_: 0xCC7D, BC_: 0x19A7,
				DE_: 0x6CEA, HL_: 0xD9FE,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{16361, 236}, {16362, 105}, {16363, 14}},
			},
			want: z80State{
				A: 0x98, F: 0x9A,
				B: 0x93, C: 0x99,
				D: 0x4C, E: 0xAF,
				H: 0x5A, L: 0x50,
				I: 0xA6, R: 0x7C,
				PC: 0x3FEC, SP: 0x806F,
				IX: 0xA291, IY: 0x28A6,
				AF_: 0xCC7D, BC_: 0x19A7,
				DE_: 0x6CEA, HL_: 0xD9FE,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{16361, 236}, {16362, 105}, {16363, 14}},
				Cycles: 10,
			},
		},
		{ // E8 0000 — RET PE (PE met → taken)
			name: "E8 0000",
			init: z80State{
				A: 0xFA, F: 0x27,
				B: 0x45, C: 0xDF,
				D: 0xEA, E: 0x63,
				H: 0x9E, L: 0xA2,
				I: 0x0F, R: 0x3E,
				PC: 0x1FBC, SP: 0x203F,
				IX: 0xFD2E, IY: 0x5EC8,
				AF_: 0x4227, BC_: 0xCD06,
				DE_: 0xA1A3, HL_: 0x0D6A,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{8124, 232}, {8255, 218}, {8256, 243}},
			},
			want: z80State{
				A: 0xFA, F: 0x27,
				B: 0x45, C: 0xDF,
				D: 0xEA, E: 0x63,
				H: 0x9E, L: 0xA2,
				I: 0x0F, R: 0x3F,
				PC: 0xF3DA, SP: 0x2041,
				IX: 0xFD2E, IY: 0x5EC8,
				AF_: 0x4227, BC_: 0xCD06,
				DE_: 0xA1A3, HL_: 0x0D6A,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{8124, 232}, {8255, 218}, {8256, 243}},
				Cycles: 11,
			},
		},
		{ // F2 0000 — JP P,nn (P not met → not taken)
			name: "F2 0000",
			init: z80State{
				A: 0x58, F: 0xCA,
				B: 0x03, C: 0xF7,
				D: 0x77, E: 0xAF,
				H: 0xB3, L: 0x4C,
				I: 0xBC, R: 0x17,
				PC: 0x070D, SP: 0x28F2,
				IX: 0x0084, IY: 0x9DB3,
				AF_: 0xAAE4, BC_: 0x3D40,
				DE_: 0xF1B4, HL_: 0x7B4C,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1805, 242}, {1806, 194}, {1807, 108}},
			},
			want: z80State{
				A: 0x58, F: 0xCA,
				B: 0x03, C: 0xF7,
				D: 0x77, E: 0xAF,
				H: 0xB3, L: 0x4C,
				I: 0xBC, R: 0x18,
				PC: 0x0710, SP: 0x28F2,
				IX: 0x0084, IY: 0x9DB3,
				AF_: 0xAAE4, BC_: 0x3D40,
				DE_: 0xF1B4, HL_: 0x7B4C,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1805, 242}, {1806, 194}, {1807, 108}},
				Cycles: 10,
			},
		},
		{ // F4 0000 — CALL P,nn (P not met → not taken)
			name: "F4 0000",
			init: z80State{
				A: 0x78, F: 0xEC,
				B: 0xD1, C: 0x69,
				D: 0x76, E: 0x86,
				H: 0xD4, L: 0xBF,
				I: 0x45, R: 0x56,
				PC: 0xA144, SP: 0x21ED,
				IX: 0x973D, IY: 0xA1C6,
				AF_: 0xDCBA, BC_: 0xC9E4,
				DE_: 0x5BA4, HL_: 0x7903,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{41284, 244}, {41285, 218}, {41286, 13}},
			},
			want: z80State{
				A: 0x78, F: 0xEC,
				B: 0xD1, C: 0x69,
				D: 0x76, E: 0x86,
				H: 0xD4, L: 0xBF,
				I: 0x45, R: 0x57,
				PC: 0xA147, SP: 0x21ED,
				IX: 0x973D, IY: 0xA1C6,
				AF_: 0xDCBA, BC_: 0xC9E4,
				DE_: 0x5BA4, HL_: 0x7903,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{41284, 244}, {41285, 218}, {41286, 13}},
				Cycles: 10,
			},
		},
		{ // F8 0000 — RET M (M met → taken)
			name: "F8 0000",
			init: z80State{
				A: 0xB6, F: 0xB7,
				B: 0x03, C: 0xA2,
				D: 0x0A, E: 0x13,
				H: 0x24, L: 0xAE,
				I: 0xE5, R: 0x2D,
				PC: 0x0688, SP: 0xA77B,
				IX: 0xEA42, IY: 0x7564,
				AF_: 0xD713, BC_: 0x9579,
				DE_: 0xAC82, HL_: 0xBDF0,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1672, 248}, {42875, 54}, {42876, 13}},
			},
			want: z80State{
				A: 0xB6, F: 0xB7,
				B: 0x03, C: 0xA2,
				D: 0x0A, E: 0x13,
				H: 0x24, L: 0xAE,
				I: 0xE5, R: 0x2E,
				PC: 0x0D36, SP: 0xA77D,
				IX: 0xEA42, IY: 0x7564,
				AF_: 0xD713, BC_: 0x9579,
				DE_: 0xAC82, HL_: 0xBDF0,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1672, 248}, {42875, 54}, {42876, 13}},
				Cycles: 11,
			},
		},
		{ // FC 0000 — CALL M,nn (M met → taken)
			name: "FC 0000",
			init: z80State{
				A: 0xD4, F: 0xCE,
				B: 0x49, C: 0x6B,
				D: 0xCE, E: 0x3E,
				H: 0x30, L: 0x04,
				I: 0xC4, R: 0x1D,
				PC: 0xE52E, SP: 0x460F,
				IX: 0xFA78, IY: 0x569A,
				AF_: 0x0C71, BC_: 0xCA98,
				DE_: 0x09BB, HL_: 0x434E,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{17933, 0}, {17934, 0}, {58670, 252}, {58671, 152}, {58672, 201}},
			},
			want: z80State{
				A: 0xD4, F: 0xCE,
				B: 0x49, C: 0x6B,
				D: 0xCE, E: 0x3E,
				H: 0x30, L: 0x04,
				I: 0xC4, R: 0x1E,
				PC: 0xC998, SP: 0x460D,
				IX: 0xFA78, IY: 0x569A,
				AF_: 0x0C71, BC_: 0xCA98,
				DE_: 0x09BB, HL_: 0x434E,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{17933, 49}, {17934, 229}, {58670, 252}, {58671, 152}, {58672, 201}},
				Cycles: 17,
			},
		},
		{ // 30 0000
			name: "30 0000",
			init: z80State{
				A: 0x12, F: 0x0D,
				B: 0xEF, C: 0x1D,
				D: 0x17, E: 0xC7,
				H: 0x80, L: 0xB3,
				I: 0x11, R: 0x54,
				PC: 0x50D1, SP: 0xCDC7,
				IX: 0xBD30, IY: 0x6D85,
				AF_: 0x6324, BC_: 0x4D72,
				DE_: 0x0AF5, HL_: 0x4AC8,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{20689, 48}, {20690, 102}},
			},
			want: z80State{
				A: 0x12, F: 0x0D,
				B: 0xEF, C: 0x1D,
				D: 0x17, E: 0xC7,
				H: 0x80, L: 0xB3,
				I: 0x11, R: 0x55,
				PC: 0x50D3, SP: 0xCDC7,
				IX: 0xBD30, IY: 0x6D85,
				AF_: 0x6324, BC_: 0x4D72,
				DE_: 0x0AF5, HL_: 0x4AC8,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{20689, 48}, {20690, 102}},
				Cycles: 7,
			},
		},
		{ // D7 0000
			name: "D7 0000",
			init: z80State{
				A: 0xB0, F: 0x5C,
				B: 0xF3, C: 0x14,
				D: 0xAF, E: 0xBD,
				H: 0xC7, L: 0x5B,
				I: 0xD6, R: 0x43,
				PC: 0x1DB0, SP: 0xD515,
				IX: 0x0599, IY: 0xA836,
				AF_: 0x0563, BC_: 0x2090,
				DE_: 0x26B2, HL_: 0xE699,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{7600, 215}, {54547, 0}, {54548, 0}},
			},
			want: z80State{
				A: 0xB0, F: 0x5C,
				B: 0xF3, C: 0x14,
				D: 0xAF, E: 0xBD,
				H: 0xC7, L: 0x5B,
				I: 0xD6, R: 0x44,
				PC: 0x0010, SP: 0xD513,
				IX: 0x0599, IY: 0xA836,
				AF_: 0x0563, BC_: 0x2090,
				DE_: 0x26B2, HL_: 0xE699,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{7600, 215}, {54547, 177}, {54548, 29}},
				Cycles: 11,
			},
		},
	}
	for _, tt := range tests {
		runSSTTest(t, tt)
	}
}
