package z80

import "testing"

func TestLDI(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED A0 = LDI
	bus.mem[0] = 0xED
	bus.mem[1] = 0xA0
	cpu.reg.HL = 0x1000
	cpu.reg.DE = 0x2000
	cpu.reg.BC = 0x0003
	bus.mem[0x1000] = 0x42
	cycles := cpu.Step()
	if bus.mem[0x2000] != 0x42 {
		t.Errorf("LDI: (DE)=%02x want 42", bus.mem[0x2000])
	}
	if cpu.reg.HL != 0x1001 {
		t.Errorf("LDI: HL=%04x want 1001", cpu.reg.HL)
	}
	if cpu.reg.DE != 0x2001 {
		t.Errorf("LDI: DE=%04x want 2001", cpu.reg.DE)
	}
	if cpu.reg.BC != 0x0002 {
		t.Errorf("LDI: BC=%04x want 0002", cpu.reg.BC)
	}
	if cpu.getF()&flagPV == 0 {
		t.Error("PV should be set (BC != 0)")
	}
	if cycles != 16 {
		t.Errorf("cycles=%d want 16", cycles)
	}
}

func TestLDD(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xED
	bus.mem[1] = 0xA8
	cpu.reg.HL = 0x1005
	cpu.reg.DE = 0x2005
	cpu.reg.BC = 0x0001
	bus.mem[0x1005] = 0x77
	cpu.Step()
	if bus.mem[0x2005] != 0x77 {
		t.Errorf("LDD: (DE)=%02x want 77", bus.mem[0x2005])
	}
	if cpu.reg.HL != 0x1004 {
		t.Errorf("LDD: HL=%04x want 1004", cpu.reg.HL)
	}
	if cpu.reg.DE != 0x2004 {
		t.Errorf("LDD: DE=%04x want 2004", cpu.reg.DE)
	}
	if cpu.reg.BC != 0x0000 {
		t.Errorf("LDD: BC=%04x want 0000", cpu.reg.BC)
	}
	if cpu.getF()&flagPV != 0 {
		t.Error("PV should be clear (BC == 0)")
	}
}

func TestLDIR(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.HL = 0x1000
	cpu.reg.DE = 0x2000
	cpu.reg.BC = 0x0003
	bus.mem[0x1000] = 0xAA
	bus.mem[0x1001] = 0xBB
	bus.mem[0x1002] = 0xCC

	// LDIR at address 0x0100
	cpu.reg.PC = 0x0100
	bus.mem[0x0100] = 0xED
	bus.mem[0x0101] = 0xB0

	// First iteration: BC=3->2, repeats
	c1 := cpu.Step()
	if c1 != 21 {
		t.Errorf("LDIR iter 1: cycles=%d want 21", c1)
	}
	if cpu.reg.PC != 0x0100 {
		t.Errorf("LDIR should repeat: PC=%04x want 0100", cpu.reg.PC)
	}
	if bus.mem[0x2000] != 0xAA {
		t.Errorf("byte 0: %02x want AA", bus.mem[0x2000])
	}

	// Second iteration: BC=2->1, repeats
	c2 := cpu.Step()
	if c2 != 21 {
		t.Errorf("LDIR iter 2: cycles=%d want 21", c2)
	}

	// Third iteration: BC=1->0, done
	c3 := cpu.Step()
	if c3 != 16 {
		t.Errorf("LDIR final: cycles=%d want 16", c3)
	}
	if cpu.reg.PC != 0x0102 {
		t.Errorf("LDIR done: PC=%04x want 0102", cpu.reg.PC)
	}

	// Verify all bytes copied
	if bus.mem[0x2000] != 0xAA || bus.mem[0x2001] != 0xBB || bus.mem[0x2002] != 0xCC {
		t.Errorf("LDIR copy: %02x %02x %02x",
			bus.mem[0x2000], bus.mem[0x2001], bus.mem[0x2002])
	}
}

func TestLDDR(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.HL = 0x1002
	cpu.reg.DE = 0x2002
	cpu.reg.BC = 0x0002

	cpu.reg.PC = 0x0100
	bus.mem[0x0100] = 0xED
	bus.mem[0x0101] = 0xB8
	bus.mem[0x1002] = 0x11
	bus.mem[0x1001] = 0x22

	cpu.Step() // BC=2->1, repeat
	cpu.Step() // BC=1->0, done

	if bus.mem[0x2002] != 0x11 || bus.mem[0x2001] != 0x22 {
		t.Errorf("LDDR: %02x %02x", bus.mem[0x2002], bus.mem[0x2001])
	}
}

func TestCPI(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xED
	bus.mem[1] = 0xA1
	cpu.reg.AF = 0x4200
	cpu.reg.HL = 0x5000
	cpu.reg.BC = 0x0003
	bus.mem[0x5000] = 0x42
	cycles := cpu.Step()
	if cpu.reg.HL != 0x5001 {
		t.Errorf("CPI: HL=%04x want 5001", cpu.reg.HL)
	}
	if cpu.reg.BC != 0x0002 {
		t.Errorf("CPI: BC=%04x want 0002", cpu.reg.BC)
	}
	if cpu.getF()&flagZ == 0 {
		t.Error("Z should be set (match)")
	}
	if cpu.getF()&flagN == 0 {
		t.Error("N should be set")
	}
	if cycles != 16 {
		t.Errorf("cycles=%d want 16", cycles)
	}
}

func TestCPD(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xED
	bus.mem[1] = 0xA9
	cpu.reg.AF = 0x4200
	cpu.reg.HL = 0x5005
	cpu.reg.BC = 0x0001
	bus.mem[0x5005] = 0x99
	cpu.Step()
	if cpu.reg.HL != 0x5004 {
		t.Errorf("CPD: HL=%04x want 5004", cpu.reg.HL)
	}
	if cpu.getF()&flagZ != 0 {
		t.Error("Z should be clear (no match)")
	}
}

func TestCPIR(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.AF = 0x4200
	cpu.reg.HL = 0x5000
	cpu.reg.BC = 0x0003
	bus.mem[0x5000] = 0x11
	bus.mem[0x5001] = 0x42 // match at second position

	cpu.reg.PC = 0x0100
	bus.mem[0x0100] = 0xED
	bus.mem[0x0101] = 0xB1

	// First iteration: no match, BC=2, repeat
	cpu.Step()
	if cpu.reg.PC != 0x0100 {
		t.Errorf("CPIR should repeat: PC=%04x", cpu.reg.PC)
	}

	// Second iteration: match found
	cpu.Step()
	if cpu.reg.PC != 0x0102 {
		t.Errorf("CPIR should stop on match: PC=%04x", cpu.reg.PC)
	}
	if cpu.getF()&flagZ == 0 {
		t.Error("Z should be set (match found)")
	}
}

func TestSST_Block(t *testing.T) {
	tests := []struct {
		name       string
		init, want z80State
	}{
		{ // ED A0 0000
			name: "ED A0 0000",
			init: z80State{
				A: 0x4A, F: 0xF3,
				B: 0xED, C: 0xCF,
				D: 0xDA, E: 0x43,
				H: 0xE7, L: 0x65,
				I: 0xA2, R: 0x14,
				PC: 0xCE39, SP: 0x095C,
				IX: 0x756E, IY: 0x23EA,
				AF_: 0x9F92, BC_: 0x8CCC,
				DE_: 0xB5B6, HL_: 0xF130,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{52793, 237}, {52794, 160}, {55875, 0}, {59237, 232}},
			},
			want: z80State{
				A: 0x4A, F: 0xE5,
				B: 0xED, C: 0xCE,
				D: 0xDA, E: 0x44,
				H: 0xE7, L: 0x66,
				I: 0xA2, R: 0x16,
				PC: 0xCE3B, SP: 0x095C,
				IX: 0x756E, IY: 0x23EA,
				AF_: 0x9F92, BC_: 0x8CCC,
				DE_: 0xB5B6, HL_: 0xF130,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{52793, 237}, {52794, 160}, {55875, 232}, {59237, 232}},
				Cycles: 16,
			},
		},
		{ // ED A0 0001
			name: "ED A0 0001",
			init: z80State{
				A: 0x3D, F: 0x7A,
				B: 0x4E, C: 0xF5,
				D: 0xDE, E: 0xE9,
				H: 0x5C, L: 0x73,
				I: 0x02, R: 0x04,
				PC: 0x731C, SP: 0x84A8,
				IX: 0x648A, IY: 0x7446,
				AF_: 0xEDFC, BC_: 0x227A,
				DE_: 0x26E8, HL_: 0xF479,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{23667, 54}, {29468, 237}, {29469, 160}, {57065, 0}},
			},
			want: z80State{
				A: 0x3D, F: 0x64,
				B: 0x4E, C: 0xF4,
				D: 0xDE, E: 0xEA,
				H: 0x5C, L: 0x74,
				I: 0x02, R: 0x06,
				PC: 0x731E, SP: 0x84A8,
				IX: 0x648A, IY: 0x7446,
				AF_: 0xEDFC, BC_: 0x227A,
				DE_: 0x26E8, HL_: 0xF479,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{23667, 54}, {29468, 237}, {29469, 160}, {57065, 54}},
				Cycles: 16,
			},
		},
		{ // ED A8 0000
			name: "ED A8 0000",
			init: z80State{
				A: 0x87, F: 0xD9,
				B: 0x4D, C: 0x61,
				D: 0xFF, E: 0x85,
				H: 0x2A, L: 0x3D,
				I: 0xE3, R: 0x6E,
				PC: 0x048D, SP: 0x963D,
				IX: 0x502E, IY: 0xF122,
				AF_: 0x0994, BC_: 0x0D69,
				DE_: 0xD421, HL_: 0x86C4,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1165, 237}, {1166, 168}, {10813, 167}, {65413, 0}},
			},
			want: z80State{
				A: 0x87, F: 0xED,
				B: 0x4D, C: 0x60,
				D: 0xFF, E: 0x84,
				H: 0x2A, L: 0x3C,
				I: 0xE3, R: 0x70,
				PC: 0x048F, SP: 0x963D,
				IX: 0x502E, IY: 0xF122,
				AF_: 0x0994, BC_: 0x0D69,
				DE_: 0xD421, HL_: 0x86C4,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1165, 237}, {1166, 168}, {10813, 167}, {65413, 167}},
				Cycles: 16,
			},
		},
		{ // ED A8 0001
			name: "ED A8 0001",
			init: z80State{
				A: 0x6A, F: 0x27,
				B: 0x52, C: 0x64,
				D: 0x1B, E: 0x76,
				H: 0x9C, L: 0xF8,
				I: 0x38, R: 0x50,
				PC: 0xD09C, SP: 0x4573,
				IX: 0xA6FD, IY: 0xCA66,
				AF_: 0xCC76, BC_: 0xF679,
				DE_: 0xA958, HL_: 0x6297,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{7030, 0}, {40184, 186}, {53404, 237}, {53405, 168}},
			},
			want: z80State{
				A: 0x6A, F: 0x05,
				B: 0x52, C: 0x63,
				D: 0x1B, E: 0x75,
				H: 0x9C, L: 0xF7,
				I: 0x38, R: 0x52,
				PC: 0xD09E, SP: 0x4573,
				IX: 0xA6FD, IY: 0xCA66,
				AF_: 0xCC76, BC_: 0xF679,
				DE_: 0xA958, HL_: 0x6297,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{7030, 186}, {40184, 186}, {53404, 237}, {53405, 168}},
				Cycles: 16,
			},
		},
		{ // ED B0 0000
			name: "ED B0 0000",
			init: z80State{
				A: 0xB2, F: 0x22,
				B: 0x98, C: 0xEA,
				D: 0x5B, E: 0x54,
				H: 0xD0, L: 0x7B,
				I: 0x74, R: 0x73,
				PC: 0x2D25, SP: 0x0FC9,
				IX: 0x0BBB, IY: 0xF96E,
				AF_: 0x591B, BC_: 0xA463,
				DE_: 0x5F3B, HL_: 0xF9ED,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{11557, 237}, {11558, 176}, {23380, 0}, {53371, 137}},
			},
			want: z80State{
				A: 0xB2, F: 0x2C,
				B: 0x98, C: 0xE9,
				D: 0x5B, E: 0x55,
				H: 0xD0, L: 0x7C,
				I: 0x74, R: 0x75,
				PC: 0x2D25, SP: 0x0FC9,
				IX: 0x0BBB, IY: 0xF96E,
				AF_: 0x591B, BC_: 0xA463,
				DE_: 0x5F3B, HL_: 0xF9ED,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{11557, 237}, {11558, 176}, {23380, 137}, {53371, 137}},
				Cycles: 21,
			},
		},
		{ // ED B0 0001
			name: "ED B0 0001",
			init: z80State{
				A: 0xDE, F: 0xDC,
				B: 0xD0, C: 0xE4,
				D: 0x57, E: 0x80,
				H: 0x5D, L: 0x41,
				I: 0x5C, R: 0x58,
				PC: 0x11D4, SP: 0x7B67,
				IX: 0x8E7C, IY: 0x6F58,
				AF_: 0xADC1, BC_: 0x2253,
				DE_: 0x08EE, HL_: 0x2888,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{4564, 237}, {4565, 176}, {22400, 0}, {23873, 46}},
			},
			want: z80State{
				A: 0xDE, F: 0xC4,
				B: 0xD0, C: 0xE3,
				D: 0x57, E: 0x81,
				H: 0x5D, L: 0x42,
				I: 0x5C, R: 0x5A,
				PC: 0x11D4, SP: 0x7B67,
				IX: 0x8E7C, IY: 0x6F58,
				AF_: 0xADC1, BC_: 0x2253,
				DE_: 0x08EE, HL_: 0x2888,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{4564, 237}, {4565, 176}, {22400, 46}, {23873, 46}},
				Cycles: 21,
			},
		},
		{ // ED B8 0000
			name: "ED B8 0000",
			init: z80State{
				A: 0x7A, F: 0x75,
				B: 0x0C, C: 0x14,
				D: 0x50, E: 0x22,
				H: 0x94, L: 0xF7,
				I: 0x6C, R: 0x3B,
				PC: 0x562B, SP: 0x62C3,
				IX: 0x401B, IY: 0xD41C,
				AF_: 0xA592, BC_: 0x9C10,
				DE_: 0x78EC, HL_: 0x5A5A,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{20514, 0}, {22059, 237}, {22060, 184}, {38135, 175}},
			},
			want: z80State{
				A: 0x7A, F: 0x45,
				B: 0x0C, C: 0x13,
				D: 0x50, E: 0x21,
				H: 0x94, L: 0xF6,
				I: 0x6C, R: 0x3D,
				PC: 0x562B, SP: 0x62C3,
				IX: 0x401B, IY: 0xD41C,
				AF_: 0xA592, BC_: 0x9C10,
				DE_: 0x78EC, HL_: 0x5A5A,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{20514, 175}, {22059, 237}, {22060, 184}, {38135, 175}},
				Cycles: 21,
			},
		},
		{ // ED B8 0001
			name: "ED B8 0001",
			init: z80State{
				A: 0x05, F: 0xD5,
				B: 0x5D, C: 0x05,
				D: 0xF9, E: 0xD9,
				H: 0xCA, L: 0xFA,
				I: 0xDD, R: 0x72,
				PC: 0x0862, SP: 0x1614,
				IX: 0xC3E8, IY: 0xA3A5,
				AF_: 0x12FC, BC_: 0x7A66,
				DE_: 0xB49C, HL_: 0xBAF1,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{2146, 237}, {2147, 184}, {51962, 94}, {63961, 0}},
			},
			want: z80State{
				A: 0x05, F: 0xCD,
				B: 0x5D, C: 0x04,
				D: 0xF9, E: 0xD8,
				H: 0xCA, L: 0xF9,
				I: 0xDD, R: 0x74,
				PC: 0x0862, SP: 0x1614,
				IX: 0xC3E8, IY: 0xA3A5,
				AF_: 0x12FC, BC_: 0x7A66,
				DE_: 0xB49C, HL_: 0xBAF1,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{2146, 237}, {2147, 184}, {51962, 94}, {63961, 94}},
				Cycles: 21,
			},
		},
		{ // ED A1 0000
			name: "ED A1 0000",
			init: z80State{
				A: 0x14, F: 0x84,
				B: 0x58, C: 0x77,
				D: 0xAF, E: 0xF9,
				H: 0x72, L: 0x3D,
				I: 0x79, R: 0x49,
				PC: 0xD0B3, SP: 0x8312,
				IX: 0x4489, IY: 0x9B79,
				AF_: 0xCDEC, BC_: 0xB41C,
				DE_: 0x2C33, HL_: 0x81B8,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{29245, 107}, {53427, 237}, {53428, 161}},
			},
			want: z80State{
				A: 0x14, F: 0x9E,
				B: 0x58, C: 0x76,
				D: 0xAF, E: 0xF9,
				H: 0x72, L: 0x3E,
				I: 0x79, R: 0x4B,
				PC: 0xD0B5, SP: 0x8312,
				IX: 0x4489, IY: 0x9B79,
				AF_: 0xCDEC, BC_: 0xB41C,
				DE_: 0x2C33, HL_: 0x81B8,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{29245, 107}, {53427, 237}, {53428, 161}},
				Cycles: 16,
			},
		},
		{ // ED A1 0001
			name: "ED A1 0001",
			init: z80State{
				A: 0xB1, F: 0xA5,
				B: 0x26, C: 0x07,
				D: 0x3D, E: 0x06,
				H: 0x9A, L: 0xA2,
				I: 0x5C, R: 0x5E,
				PC: 0x16CB, SP: 0xBA48,
				IX: 0x7C3D, IY: 0xE892,
				AF_: 0x7CF1, BC_: 0x50C4,
				DE_: 0xD839, HL_: 0xF980,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{5835, 237}, {5836, 161}, {39586, 46}},
			},
			want: z80State{
				A: 0xB1, F: 0xB7,
				B: 0x26, C: 0x06,
				D: 0x3D, E: 0x06,
				H: 0x9A, L: 0xA3,
				I: 0x5C, R: 0x60,
				PC: 0x16CD, SP: 0xBA48,
				IX: 0x7C3D, IY: 0xE892,
				AF_: 0x7CF1, BC_: 0x50C4,
				DE_: 0xD839, HL_: 0xF980,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{5835, 237}, {5836, 161}, {39586, 46}},
				Cycles: 16,
			},
		},
		{ // ED A9 0000
			name: "ED A9 0000",
			init: z80State{
				A: 0xF5, F: 0xF0,
				B: 0x97, C: 0xAC,
				D: 0x76, E: 0x29,
				H: 0xBA, L: 0x2A,
				I: 0x6A, R: 0x50,
				PC: 0x3717, SP: 0x4B4D,
				IX: 0xC133, IY: 0xDC13,
				AF_: 0x46CB, BC_: 0xB86C,
				DE_: 0x651E, HL_: 0x2B5F,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{14103, 237}, {14104, 169}, {47658, 37}},
			},
			want: z80State{
				A: 0xF5, F: 0x86,
				B: 0x97, C: 0xAB,
				D: 0x76, E: 0x29,
				H: 0xBA, L: 0x29,
				I: 0x6A, R: 0x52,
				PC: 0x3719, SP: 0x4B4D,
				IX: 0xC133, IY: 0xDC13,
				AF_: 0x46CB, BC_: 0xB86C,
				DE_: 0x651E, HL_: 0x2B5F,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{14103, 237}, {14104, 169}, {47658, 37}},
				Cycles: 16,
			},
		},
		{ // ED A9 0001
			name: "ED A9 0001",
			init: z80State{
				A: 0x26, F: 0x0E,
				B: 0xA0, C: 0x0D,
				D: 0xD9, E: 0xD5,
				H: 0xD5, L: 0xA2,
				I: 0xD1, R: 0x48,
				PC: 0x1430, SP: 0x0487,
				IX: 0x2367, IY: 0x2FD8,
				AF_: 0x4F64, BC_: 0xC416,
				DE_: 0x7182, HL_: 0x20B8,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{5168, 237}, {5169, 169}, {54690, 202}},
			},
			want: z80State{
				A: 0x26, F: 0x3E,
				B: 0xA0, C: 0x0C,
				D: 0xD9, E: 0xD5,
				H: 0xD5, L: 0xA1,
				I: 0xD1, R: 0x4A,
				PC: 0x1432, SP: 0x0487,
				IX: 0x2367, IY: 0x2FD8,
				AF_: 0x4F64, BC_: 0xC416,
				DE_: 0x7182, HL_: 0x20B8,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{5168, 237}, {5169, 169}, {54690, 202}},
				Cycles: 16,
			},
		},
		{ // ED B1 0000
			name: "ED B1 0000",
			init: z80State{
				A: 0xA3, F: 0xCB,
				B: 0x14, C: 0x95,
				D: 0xCE, E: 0xB3,
				H: 0xBF, L: 0xF6,
				I: 0xFC, R: 0x12,
				PC: 0xE842, SP: 0x9C9E,
				IX: 0x6D73, IY: 0x32C5,
				AF_: 0x01D3, BC_: 0xA700,
				DE_: 0x5010, HL_: 0x4A5E,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{49142, 145}, {59458, 237}, {59459, 177}},
			},
			want: z80State{
				A: 0xA3, F: 0x2F,
				B: 0x14, C: 0x94,
				D: 0xCE, E: 0xB3,
				H: 0xBF, L: 0xF7,
				I: 0xFC, R: 0x14,
				PC: 0xE842, SP: 0x9C9E,
				IX: 0x6D73, IY: 0x32C5,
				AF_: 0x01D3, BC_: 0xA700,
				DE_: 0x5010, HL_: 0x4A5E,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{49142, 145}, {59458, 237}, {59459, 177}},
				Cycles: 21,
			},
		},
		{ // ED B1 0001
			name: "ED B1 0001",
			init: z80State{
				A: 0x5C, F: 0xDC,
				B: 0x0E, C: 0x99,
				D: 0x36, E: 0x1D,
				H: 0xF1, L: 0x40,
				I: 0xE3, R: 0x3F,
				PC: 0x6B51, SP: 0x2773,
				IX: 0x2931, IY: 0xEA85,
				AF_: 0x9309, BC_: 0xDC43,
				DE_: 0x3ADB, HL_: 0xCD5F,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{27473, 237}, {27474, 177}, {61760, 142}},
			},
			want: z80State{
				A: 0x5C, F: 0xBE,
				B: 0x0E, C: 0x98,
				D: 0x36, E: 0x1D,
				H: 0xF1, L: 0x41,
				I: 0xE3, R: 0x41,
				PC: 0x6B51, SP: 0x2773,
				IX: 0x2931, IY: 0xEA85,
				AF_: 0x9309, BC_: 0xDC43,
				DE_: 0x3ADB, HL_: 0xCD5F,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{27473, 237}, {27474, 177}, {61760, 142}},
				Cycles: 21,
			},
		},
		{ // ED B9 0000
			name: "ED B9 0000",
			init: z80State{
				A: 0x2F, F: 0x4F,
				B: 0x2B, C: 0xD9,
				D: 0xEE, E: 0x49,
				H: 0x20, L: 0xAB,
				I: 0x12, R: 0x3C,
				PC: 0xB4F3, SP: 0x054B,
				IX: 0x32A5, IY: 0xD43A,
				AF_: 0x319D, BC_: 0x7EF2,
				DE_: 0xDCFA, HL_: 0xE86D,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{8363, 74}, {46323, 237}, {46324, 185}},
			},
			want: z80State{
				A: 0x2F, F: 0xA7,
				B: 0x2B, C: 0xD8,
				D: 0xEE, E: 0x49,
				H: 0x20, L: 0xAA,
				I: 0x12, R: 0x3E,
				PC: 0xB4F3, SP: 0x054B,
				IX: 0x32A5, IY: 0xD43A,
				AF_: 0x319D, BC_: 0x7EF2,
				DE_: 0xDCFA, HL_: 0xE86D,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{8363, 74}, {46323, 237}, {46324, 185}},
				Cycles: 21,
			},
		},
		{ // ED B9 0001
			name: "ED B9 0001",
			init: z80State{
				A: 0x99, F: 0x35,
				B: 0x1D, C: 0x81,
				D: 0x8C, E: 0xDB,
				H: 0x96, L: 0x0C,
				I: 0xAE, R: 0x26,
				PC: 0xB5A0, SP: 0x3446,
				IX: 0xF5AC, IY: 0x06D0,
				AF_: 0x6466, BC_: 0xACAD,
				DE_: 0x6120, HL_: 0xF085,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{38412, 144}, {46496, 237}, {46497, 185}},
			},
			want: z80State{
				A: 0x99, F: 0x27,
				B: 0x1D, C: 0x80,
				D: 0x8C, E: 0xDB,
				H: 0x96, L: 0x0B,
				I: 0xAE, R: 0x28,
				PC: 0xB5A0, SP: 0x3446,
				IX: 0xF5AC, IY: 0x06D0,
				AF_: 0x6466, BC_: 0xACAD,
				DE_: 0x6120, HL_: 0xF085,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{38412, 144}, {46496, 237}, {46497, 185}},
				Cycles: 21,
			},
		},
	}
	for _, tt := range tests {
		runSSTTest(t, tt)
	}
}
