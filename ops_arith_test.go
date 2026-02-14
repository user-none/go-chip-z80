package z80

import "testing"

func TestADD_A_r(t *testing.T) {
	cpu, bus := newTestCPU()
	// ADD A, B (0x80)
	bus.mem[0] = 0x80
	cpu.reg.AF = 0x1000
	cpu.reg.BC = 0x2000
	cycles := cpu.Step()
	if cpu.getA() != 0x30 {
		t.Errorf("ADD A,B: A=%02x want 30", cpu.getA())
	}
	if cycles != 4 {
		t.Errorf("cycles=%d want 4", cycles)
	}
	// Check flags: no carry, no half-carry, no overflow
	f := cpu.getF()
	if f&flagZ != 0 {
		t.Error("Z should be clear")
	}
	if f&flagN != 0 {
		t.Error("N should be clear")
	}
}

func TestADD_A_r_Overflow(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x80 // ADD A,B
	cpu.reg.AF = 0x7F00
	cpu.reg.BC = 0x0100
	cpu.Step()
	if cpu.getA() != 0x80 {
		t.Errorf("A=%02x want 80", cpu.getA())
	}
	f := cpu.getF()
	if f&flagPV == 0 {
		t.Error("overflow should be set")
	}
	if f&flagS == 0 {
		t.Error("sign should be set")
	}
	if f&flagH == 0 {
		t.Error("half-carry should be set")
	}
}

func TestADD_A_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x86 // ADD A,(HL)
	cpu.reg.AF = 0x1000
	cpu.reg.HL = 0x5000
	bus.mem[0x5000] = 0x20
	cycles := cpu.Step()
	if cpu.getA() != 0x30 {
		t.Errorf("ADD A,(HL): A=%02x want 30", cpu.getA())
	}
	if cycles != 7 {
		t.Errorf("cycles=%d want 7", cycles)
	}
}

func TestADC_A_r(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x89 // ADC A,C
	cpu.reg.AF = 0x1001 // A=0x10, F=carry set
	cpu.reg.BC = 0x0020
	cpu.Step()
	if cpu.getA() != 0x31 {
		t.Errorf("ADC A,C: A=%02x want 31", cpu.getA())
	}
}

func TestSUB_A_r(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x90 // SUB B
	cpu.reg.AF = 0x3000
	cpu.reg.BC = 0x1000
	cpu.Step()
	if cpu.getA() != 0x20 {
		t.Errorf("SUB B: A=%02x want 20", cpu.getA())
	}
	f := cpu.getF()
	if f&flagN == 0 {
		t.Error("N should be set for SUB")
	}
}

func TestSUB_A_Zero(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x90 // SUB B
	cpu.reg.AF = 0x4200
	cpu.reg.BC = 0x4200
	cpu.Step()
	if cpu.getA() != 0 {
		t.Errorf("A=%02x want 0", cpu.getA())
	}
	if cpu.getF()&flagZ == 0 {
		t.Error("Z should be set")
	}
}

func TestAND_r(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xA0 // AND B
	cpu.reg.AF = 0xF000
	cpu.reg.BC = 0x0F00
	cpu.Step()
	if cpu.getA() != 0x00 {
		t.Errorf("AND B: A=%02x want 00", cpu.getA())
	}
	f := cpu.getF()
	if f&flagH == 0 {
		t.Error("H should be set for AND")
	}
	if f&flagZ == 0 {
		t.Error("Z should be set")
	}
}

func TestXOR_r(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xA8 // XOR B
	cpu.reg.AF = 0xFF00
	cpu.reg.BC = 0x0F00
	cpu.Step()
	if cpu.getA() != 0xF0 {
		t.Errorf("XOR B: A=%02x want F0", cpu.getA())
	}
}

func TestOR_r(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xB0 // OR B
	cpu.reg.AF = 0xF000
	cpu.reg.BC = 0x0F00
	cpu.Step()
	if cpu.getA() != 0xFF {
		t.Errorf("OR B: A=%02x want FF", cpu.getA())
	}
}

func TestCP_r(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xB8 // CP B
	cpu.reg.AF = 0x4200
	cpu.reg.BC = 0x4200
	cpu.Step()
	// CP doesn't change A
	if cpu.getA() != 0x42 {
		t.Errorf("CP B: A=%02x want 42 (unchanged)", cpu.getA())
	}
	if cpu.getF()&flagZ == 0 {
		t.Error("Z should be set for equal operands")
	}
}

func TestALU_A_n(t *testing.T) {
	cpu, bus := newTestCPU()
	// ADD A, 0x05 (0xC6 0x05)
	bus.mem[0] = 0xC6
	bus.mem[1] = 0x05
	cpu.reg.AF = 0x0A00
	cycles := cpu.Step()
	if cpu.getA() != 0x0F {
		t.Errorf("ADD A,n: A=%02x want 0F", cpu.getA())
	}
	if cycles != 7 {
		t.Errorf("cycles=%d want 7", cycles)
	}
}

func TestINC_r(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x04 // INC B
	cpu.reg.BC = 0xFF00
	cpu.reg.AF = 0x0001 // carry set
	cycles := cpu.Step()
	if cpu.getB() != 0x00 {
		t.Errorf("INC B: B=%02x want 00", cpu.getB())
	}
	if cycles != 4 {
		t.Errorf("cycles=%d want 4", cycles)
	}
	f := cpu.getF()
	if f&flagZ == 0 {
		t.Error("Z should be set")
	}
	if f&flagH == 0 {
		t.Error("H should be set")
	}
	// Carry should be preserved
	if f&flagC == 0 {
		t.Error("C should be preserved")
	}
}

func TestINC_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x34 // INC (HL)
	cpu.reg.HL = 0x4000
	bus.mem[0x4000] = 0x0F
	cycles := cpu.Step()
	if bus.mem[0x4000] != 0x10 {
		t.Errorf("INC (HL): got %02x want 10", bus.mem[0x4000])
	}
	if cycles != 11 {
		t.Errorf("cycles=%d want 11", cycles)
	}
	if cpu.getF()&flagH == 0 {
		t.Error("H should be set for 0F->10")
	}
}

func TestDEC_r(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x05 // DEC B
	cpu.reg.BC = 0x0100
	cpu.Step()
	if cpu.getB() != 0x00 {
		t.Errorf("DEC B: B=%02x want 00", cpu.getB())
	}
	if cpu.getF()&flagZ == 0 {
		t.Error("Z should be set")
	}
	if cpu.getF()&flagN == 0 {
		t.Error("N should be set for DEC")
	}
}

func TestINC_DEC_rr(t *testing.T) {
	cpu, bus := newTestCPU()
	// INC BC (0x03)
	bus.mem[0] = 0x03
	cpu.reg.BC = 0xFFFF
	cycles := cpu.Step()
	if cpu.reg.BC != 0x0000 {
		t.Errorf("INC BC: BC=%04x want 0000", cpu.reg.BC)
	}
	if cycles != 6 {
		t.Errorf("cycles=%d want 6", cycles)
	}
	// DEC DE (0x1B)
	bus.mem[1] = 0x1B
	cpu.reg.DE = 0x0000
	cpu.Step()
	if cpu.reg.DE != 0xFFFF {
		t.Errorf("DEC DE: DE=%04x want FFFF", cpu.reg.DE)
	}
}

func TestADD_HL_rr(t *testing.T) {
	cpu, bus := newTestCPU()
	// ADD HL, BC (0x09)
	bus.mem[0] = 0x09
	cpu.reg.HL = 0x1000
	cpu.reg.BC = 0x2000
	cycles := cpu.Step()
	if cpu.reg.HL != 0x3000 {
		t.Errorf("ADD HL,BC: HL=%04x want 3000", cpu.reg.HL)
	}
	if cycles != 11 {
		t.Errorf("cycles=%d want 11", cycles)
	}
	if cpu.getF()&flagN != 0 {
		t.Error("N should be clear")
	}
}

func TestADD_HL_rr_Carry(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x09
	cpu.reg.HL = 0x8000
	cpu.reg.BC = 0x8000
	cpu.Step()
	if cpu.reg.HL != 0x0000 {
		t.Errorf("HL=%04x want 0000", cpu.reg.HL)
	}
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set")
	}
}

func TestDAA(t *testing.T) {
	cpu, bus := newTestCPU()
	// 0x15 + 0x27 in BCD = 0x42
	bus.mem[0] = 0x80 // ADD A,B
	bus.mem[1] = 0x27 // DAA
	cpu.reg.AF = 0x1500
	cpu.reg.BC = 0x2700
	cpu.Step() // ADD A,B -> A=0x3C
	cpu.Step() // DAA -> should adjust to 0x42
	if cpu.getA() != 0x42 {
		t.Errorf("DAA: A=%02x want 42", cpu.getA())
	}
}

func TestCPL(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x2F
	cpu.reg.AF = 0xAA00
	cpu.Step()
	if cpu.getA() != 0x55 {
		t.Errorf("CPL: A=%02x want 55", cpu.getA())
	}
	f := cpu.getF()
	if f&flagH == 0 || f&flagN == 0 {
		t.Error("H and N should be set")
	}
}

func TestSCF(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x37
	cpu.reg.AF = 0x0000
	cpu.Step()
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set")
	}
	if cpu.getF()&flagN != 0 || cpu.getF()&flagH != 0 {
		t.Error("N and H should be clear")
	}
}

func TestCCF(t *testing.T) {
	cpu, bus := newTestCPU()
	// Set carry first
	bus.mem[0] = 0x37 // SCF
	bus.mem[1] = 0x3F // CCF
	cpu.Step()
	cpu.Step()
	if cpu.getF()&flagC != 0 {
		t.Error("C should be clear after CCF")
	}
	if cpu.getF()&flagH == 0 {
		t.Error("H should be set (old carry)")
	}
}

func TestNEG(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED 44 = NEG
	bus.mem[0] = 0xED
	bus.mem[1] = 0x44
	cpu.reg.AF = 0x0100
	cycles := cpu.Step()
	if cpu.getA() != 0xFF {
		t.Errorf("NEG: A=%02x want FF", cpu.getA())
	}
	if cycles != 8 {
		t.Errorf("cycles=%d want 8", cycles)
	}
	if cpu.getF()&flagN == 0 {
		t.Error("N should be set")
	}
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set for non-zero operand")
	}
}

func TestSST_Arith(t *testing.T) {
	tests := []struct {
		name       string
		init, want z80State
	}{
		{ // 80 0000
			name: "80 0000",
			init: z80State{
				A: 0x51, F: 0xA7,
				B: 0x5C, C: 0x86,
				D: 0xB1, E: 0xE6,
				H: 0xA6, L: 0x7E,
				I: 0x3F, R: 0x70,
				PC: 0xC8C7, SP: 0xEA62,
				IX: 0xAE15, IY: 0x5F63,
				AF_: 0x2945, BC_: 0xC464,
				DE_: 0x1112, HL_: 0x2182,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{51399, 128}},
			},
			want: z80State{
				A: 0xAD, F: 0xAC,
				B: 0x5C, C: 0x86,
				D: 0xB1, E: 0xE6,
				H: 0xA6, L: 0x7E,
				I: 0x3F, R: 0x71,
				PC: 0xC8C8, SP: 0xEA62,
				IX: 0xAE15, IY: 0x5F63,
				AF_: 0x2945, BC_: 0xC464,
				DE_: 0x1112, HL_: 0x2182,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{51399, 128}},
				Cycles: 4,
			},
		},
		{ // 86 0000
			name: "86 0000",
			init: z80State{
				A: 0xB0, F: 0x45,
				B: 0xF7, C: 0xF8,
				D: 0x67, E: 0x73,
				H: 0x43, L: 0xC3,
				I: 0x4B, R: 0x4D,
				PC: 0x8341, SP: 0xA7CB,
				IX: 0x0886, IY: 0x931B,
				AF_: 0x0101, BC_: 0x681C,
				DE_: 0x9FAD, HL_: 0x1842,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{17347, 190}, {33601, 134}},
			},
			want: z80State{
				A: 0x6E, F: 0x2D,
				B: 0xF7, C: 0xF8,
				D: 0x67, E: 0x73,
				H: 0x43, L: 0xC3,
				I: 0x4B, R: 0x4E,
				PC: 0x8342, SP: 0xA7CB,
				IX: 0x0886, IY: 0x931B,
				AF_: 0x0101, BC_: 0x681C,
				DE_: 0x9FAD, HL_: 0x1842,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{17347, 190}, {33601, 134}},
				Cycles: 7,
			},
		},
		{ // C6 0000
			name: "C6 0000",
			init: z80State{
				A: 0xC5, F: 0xDF,
				B: 0x42, C: 0xDF,
				D: 0x83, E: 0xBE,
				H: 0xBB, L: 0x41,
				I: 0x31, R: 0x76,
				PC: 0xE9C6, SP: 0xD5EE,
				IX: 0xFC5D, IY: 0x2C60,
				AF_: 0x4F47, BC_: 0x7667,
				DE_: 0x21B4, HL_: 0x58DF,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{59846, 198}, {59847, 145}},
			},
			want: z80State{
				A: 0x56, F: 0x05,
				B: 0x42, C: 0xDF,
				D: 0x83, E: 0xBE,
				H: 0xBB, L: 0x41,
				I: 0x31, R: 0x77,
				PC: 0xE9C8, SP: 0xD5EE,
				IX: 0xFC5D, IY: 0x2C60,
				AF_: 0x4F47, BC_: 0x7667,
				DE_: 0x21B4, HL_: 0x58DF,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{59846, 198}, {59847, 145}},
				Cycles: 7,
			},
		},
		{ // 89 0000
			name: "89 0000",
			init: z80State{
				A: 0x46, F: 0x60,
				B: 0x66, C: 0xF7,
				D: 0x3B, E: 0xFE,
				H: 0x1F, L: 0x31,
				I: 0x45, R: 0x39,
				PC: 0x84DA, SP: 0x6980,
				IX: 0x46E9, IY: 0xD3F0,
				AF_: 0xCCCD, BC_: 0x02EE,
				DE_: 0x4093, HL_: 0x08E4,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{34010, 137}},
			},
			want: z80State{
				A: 0x3D, F: 0x29,
				B: 0x66, C: 0xF7,
				D: 0x3B, E: 0xFE,
				H: 0x1F, L: 0x31,
				I: 0x45, R: 0x3A,
				PC: 0x84DB, SP: 0x6980,
				IX: 0x46E9, IY: 0xD3F0,
				AF_: 0xCCCD, BC_: 0x02EE,
				DE_: 0x4093, HL_: 0x08E4,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{34010, 137}},
				Cycles: 4,
			},
		},
		{ // 90 0000
			name: "90 0000",
			init: z80State{
				A: 0x3C, F: 0x3C,
				B: 0xB4, C: 0x38,
				D: 0x3D, E: 0xBD,
				H: 0x22, L: 0x9B,
				I: 0x7E, R: 0x1C,
				PC: 0x0926, SP: 0x8708,
				IX: 0xA4A7, IY: 0xBE2F,
				AF_: 0xDE46, BC_: 0x6C5F,
				DE_: 0x8657, HL_: 0xD7BA,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{2342, 144}},
			},
			want: z80State{
				A: 0x88, F: 0x8F,
				B: 0xB4, C: 0x38,
				D: 0x3D, E: 0xBD,
				H: 0x22, L: 0x9B,
				I: 0x7E, R: 0x1D,
				PC: 0x0927, SP: 0x8708,
				IX: 0xA4A7, IY: 0xBE2F,
				AF_: 0xDE46, BC_: 0x6C5F,
				DE_: 0x8657, HL_: 0xD7BA,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{2342, 144}},
				Cycles: 4,
			},
		},
		{ // 99 0000
			name: "99 0000",
			init: z80State{
				A: 0x78, F: 0xAB,
				B: 0xF0, C: 0x23,
				D: 0xBF, E: 0x72,
				H: 0x7B, L: 0x68,
				I: 0x75, R: 0x23,
				PC: 0xA66F, SP: 0x2E01,
				IX: 0x9D38, IY: 0xC129,
				AF_: 0xF49E, BC_: 0xF747,
				DE_: 0x7602, HL_: 0x6DA7,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{42607, 153}},
			},
			want: z80State{
				A: 0x54, F: 0x02,
				B: 0xF0, C: 0x23,
				D: 0xBF, E: 0x72,
				H: 0x7B, L: 0x68,
				I: 0x75, R: 0x24,
				PC: 0xA670, SP: 0x2E01,
				IX: 0x9D38, IY: 0xC129,
				AF_: 0xF49E, BC_: 0xF747,
				DE_: 0x7602, HL_: 0x6DA7,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{42607, 153}},
				Cycles: 4,
			},
		},
		{ // A0 0000
			name: "A0 0000",
			init: z80State{
				A: 0x74, F: 0xBE,
				B: 0x12, C: 0x9E,
				D: 0x37, E: 0xB5,
				H: 0x6A, L: 0xA1,
				I: 0x0A, R: 0x46,
				PC: 0x6FCF, SP: 0x82EA,
				IX: 0xB4ED, IY: 0xCB42,
				AF_: 0x6E82, BC_: 0xC7C1,
				DE_: 0xBE7B, HL_: 0x8C1A,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{28623, 160}},
			},
			want: z80State{
				A: 0x10, F: 0x10,
				B: 0x12, C: 0x9E,
				D: 0x37, E: 0xB5,
				H: 0x6A, L: 0xA1,
				I: 0x0A, R: 0x47,
				PC: 0x6FD0, SP: 0x82EA,
				IX: 0xB4ED, IY: 0xCB42,
				AF_: 0x6E82, BC_: 0xC7C1,
				DE_: 0xBE7B, HL_: 0x8C1A,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{28623, 160}},
				Cycles: 4,
			},
		},
		{ // A8 0000
			name: "A8 0000",
			init: z80State{
				A: 0x5C, F: 0x63,
				B: 0x24, C: 0x1B,
				D: 0x1C, E: 0xC6,
				H: 0x12, L: 0xBE,
				I: 0x2B, R: 0x06,
				PC: 0x641A, SP: 0x4569,
				IX: 0x95A7, IY: 0x4902,
				AF_: 0x9943, BC_: 0x754B,
				DE_: 0x7F38, HL_: 0x78D3,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{25626, 168}},
			},
			want: z80State{
				A: 0x78, F: 0x2C,
				B: 0x24, C: 0x1B,
				D: 0x1C, E: 0xC6,
				H: 0x12, L: 0xBE,
				I: 0x2B, R: 0x07,
				PC: 0x641B, SP: 0x4569,
				IX: 0x95A7, IY: 0x4902,
				AF_: 0x9943, BC_: 0x754B,
				DE_: 0x7F38, HL_: 0x78D3,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{25626, 168}},
				Cycles: 4,
			},
		},
		{ // B0 0000
			name: "B0 0000",
			init: z80State{
				A: 0x84, F: 0xEC,
				B: 0xD8, C: 0xE3,
				D: 0xE2, E: 0x70,
				H: 0xD2, L: 0x7E,
				I: 0x10, R: 0x65,
				PC: 0x8E59, SP: 0x5D95,
				IX: 0x0DAD, IY: 0xD153,
				AF_: 0x3B37, BC_: 0xB57B,
				DE_: 0xE256, HL_: 0xDD68,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{36441, 176}},
			},
			want: z80State{
				A: 0xDC, F: 0x88,
				B: 0xD8, C: 0xE3,
				D: 0xE2, E: 0x70,
				H: 0xD2, L: 0x7E,
				I: 0x10, R: 0x66,
				PC: 0x8E5A, SP: 0x5D95,
				IX: 0x0DAD, IY: 0xD153,
				AF_: 0x3B37, BC_: 0xB57B,
				DE_: 0xE256, HL_: 0xDD68,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{36441, 176}},
				Cycles: 4,
			},
		},
		{ // FE 0000
			name: "FE 0000",
			init: z80State{
				A: 0x49, F: 0xFF,
				B: 0x5F, C: 0x35,
				D: 0x3E, E: 0x2E,
				H: 0x1A, L: 0x9A,
				I: 0xDE, R: 0x2D,
				PC: 0xFD21, SP: 0xEB66,
				IX: 0x404C, IY: 0x4924,
				AF_: 0x8C83, BC_: 0xF584,
				DE_: 0x60AB, HL_: 0x885A,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{64801, 254}, {64802, 192}},
			},
			want: z80State{
				A: 0x49, F: 0x87,
				B: 0x5F, C: 0x35,
				D: 0x3E, E: 0x2E,
				H: 0x1A, L: 0x9A,
				I: 0xDE, R: 0x2E,
				PC: 0xFD23, SP: 0xEB66,
				IX: 0x404C, IY: 0x4924,
				AF_: 0x8C83, BC_: 0xF584,
				DE_: 0x60AB, HL_: 0x885A,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{64801, 254}, {64802, 192}},
				Cycles: 7,
			},
		},
		{ // 04 0000
			name: "04 0000",
			init: z80State{
				A: 0x7D, F: 0x0B,
				B: 0x56, C: 0xB6,
				D: 0x17, E: 0xEC,
				H: 0xA7, L: 0x60,
				I: 0xD7, R: 0x19,
				PC: 0x993B, SP: 0xA0BF,
				IX: 0xFFC8, IY: 0xE6A0,
				AF_: 0x464A, BC_: 0x150F,
				DE_: 0xB105, HL_: 0x05E0,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{39227, 4}},
			},
			want: z80State{
				A: 0x7D, F: 0x01,
				B: 0x57, C: 0xB6,
				D: 0x17, E: 0xEC,
				H: 0xA7, L: 0x60,
				I: 0xD7, R: 0x1A,
				PC: 0x993C, SP: 0xA0BF,
				IX: 0xFFC8, IY: 0xE6A0,
				AF_: 0x464A, BC_: 0x150F,
				DE_: 0xB105, HL_: 0x05E0,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{39227, 4}},
				Cycles: 4,
			},
		},
		{ // 34 0000
			name: "34 0000",
			init: z80State{
				A: 0xF2, F: 0xAB,
				B: 0xBF, C: 0xD4,
				D: 0x64, E: 0xB2,
				H: 0x8D, L: 0x10,
				I: 0x94, R: 0x5E,
				PC: 0x4C94, SP: 0xDF02,
				IX: 0x75F7, IY: 0xC1D5,
				AF_: 0x5A5E, BC_: 0xFC53,
				DE_: 0xDC15, HL_: 0xAEFE,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{19604, 52}, {36112, 94}},
			},
			want: z80State{
				A: 0xF2, F: 0x09,
				B: 0xBF, C: 0xD4,
				D: 0x64, E: 0xB2,
				H: 0x8D, L: 0x10,
				I: 0x94, R: 0x5F,
				PC: 0x4C95, SP: 0xDF02,
				IX: 0x75F7, IY: 0xC1D5,
				AF_: 0x5A5E, BC_: 0xFC53,
				DE_: 0xDC15, HL_: 0xAEFE,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{19604, 52}, {36112, 95}},
				Cycles: 11,
			},
		},
		{ // 3D 0000
			name: "3D 0000",
			init: z80State{
				A: 0x5A, F: 0xF8,
				B: 0xF1, C: 0x9D,
				D: 0x95, E: 0x50,
				H: 0x57, L: 0xDB,
				I: 0xD4, R: 0x29,
				PC: 0x416E, SP: 0x2ABD,
				IX: 0x6C44, IY: 0x45DF,
				AF_: 0xC439, BC_: 0xDF73,
				DE_: 0xAD07, HL_: 0xC204,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{16750, 61}},
			},
			want: z80State{
				A: 0x59, F: 0x0A,
				B: 0xF1, C: 0x9D,
				D: 0x95, E: 0x50,
				H: 0x57, L: 0xDB,
				I: 0xD4, R: 0x2A,
				PC: 0x416F, SP: 0x2ABD,
				IX: 0x6C44, IY: 0x45DF,
				AF_: 0xC439, BC_: 0xDF73,
				DE_: 0xAD07, HL_: 0xC204,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{16750, 61}},
				Cycles: 4,
			},
		},
		{ // 35 0000
			name: "35 0000",
			init: z80State{
				A: 0x77, F: 0xD6,
				B: 0x41, C: 0x4E,
				D: 0xCA, E: 0xE4,
				H: 0x49, L: 0xA5,
				I: 0xF4, R: 0x55,
				PC: 0x292E, SP: 0x98D9,
				IX: 0xEC37, IY: 0xC6AB,
				AF_: 0x6EEE, BC_: 0xBFED,
				DE_: 0x04DB, HL_: 0xB2AB,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{10542, 53}, {18853, 164}},
			},
			want: z80State{
				A: 0x77, F: 0xA2,
				B: 0x41, C: 0x4E,
				D: 0xCA, E: 0xE4,
				H: 0x49, L: 0xA5,
				I: 0xF4, R: 0x56,
				PC: 0x292F, SP: 0x98D9,
				IX: 0xEC37, IY: 0xC6AB,
				AF_: 0x6EEE, BC_: 0xBFED,
				DE_: 0x04DB, HL_: 0xB2AB,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{10542, 53}, {18853, 163}},
				Cycles: 11,
			},
		},
		{ // 03 0000
			name: "03 0000",
			init: z80State{
				A: 0x97, F: 0x0D,
				B: 0xA5, C: 0x64,
				D: 0xA2, E: 0x21,
				H: 0x42, L: 0x02,
				I: 0x1B, R: 0x14,
				PC: 0xBF86, SP: 0xAC41,
				IX: 0xA7FE, IY: 0x245A,
				AF_: 0x4B1D, BC_: 0xA669,
				DE_: 0x27C5, HL_: 0x85D3,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{49030, 3}},
			},
			want: z80State{
				A: 0x97, F: 0x0D,
				B: 0xA5, C: 0x65,
				D: 0xA2, E: 0x21,
				H: 0x42, L: 0x02,
				I: 0x1B, R: 0x15,
				PC: 0xBF87, SP: 0xAC41,
				IX: 0xA7FE, IY: 0x245A,
				AF_: 0x4B1D, BC_: 0xA669,
				DE_: 0x27C5, HL_: 0x85D3,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{49030, 3}},
				Cycles: 6,
			},
		},
		{ // 0B 0000
			name: "0B 0000",
			init: z80State{
				A: 0x8D, F: 0x96,
				B: 0x02, C: 0x04,
				D: 0x30, E: 0x02,
				H: 0xBB, L: 0x53,
				I: 0x35, R: 0x1F,
				PC: 0xD191, SP: 0xED4A,
				IX: 0xF7B6, IY: 0xDDA2,
				AF_: 0x8CCB, BC_: 0xB5F9,
				DE_: 0x690F, HL_: 0x3730,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{53649, 11}},
			},
			want: z80State{
				A: 0x8D, F: 0x96,
				B: 0x02, C: 0x03,
				D: 0x30, E: 0x02,
				H: 0xBB, L: 0x53,
				I: 0x35, R: 0x20,
				PC: 0xD192, SP: 0xED4A,
				IX: 0xF7B6, IY: 0xDDA2,
				AF_: 0x8CCB, BC_: 0xB5F9,
				DE_: 0x690F, HL_: 0x3730,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{53649, 11}},
				Cycles: 6,
			},
		},
		{ // 09 0000
			name: "09 0000",
			init: z80State{
				A: 0x1D, F: 0x4B,
				B: 0x61, C: 0xF3,
				D: 0x72, E: 0xF8,
				H: 0xB0, L: 0x15,
				I: 0x5B, R: 0x7E,
				PC: 0xA6CB, SP: 0x2C91,
				IX: 0xEDEA, IY: 0x2BD5,
				AF_: 0x31C1, BC_: 0x85EB,
				DE_: 0x7B60, HL_: 0x5166,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{42699, 9}},
			},
			want: z80State{
				A: 0x1D, F: 0x41,
				B: 0x61, C: 0xF3,
				D: 0x72, E: 0xF8,
				H: 0x12, L: 0x08,
				I: 0x5B, R: 0x7F,
				PC: 0xA6CC, SP: 0x2C91,
				IX: 0xEDEA, IY: 0x2BD5,
				AF_: 0x31C1, BC_: 0x85EB,
				DE_: 0x7B60, HL_: 0x5166,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{42699, 9}},
				Cycles: 11,
			},
		},
		{ // 27 0000
			name: "27 0000",
			init: z80State{
				A: 0xF3, F: 0x4E,
				B: 0xB0, C: 0xA9,
				D: 0x0A, E: 0x3D,
				H: 0x7D, L: 0xAA,
				I: 0x9F, R: 0x2F,
				PC: 0x3D52, SP: 0x0024,
				IX: 0x7E3F, IY: 0x1584,
				AF_: 0xFC7B, BC_: 0x68F2,
				DE_: 0x616B, HL_: 0x0EF7,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{15698, 39}},
			},
			want: z80State{
				A: 0x93, F: 0x87,
				B: 0xB0, C: 0xA9,
				D: 0x0A, E: 0x3D,
				H: 0x7D, L: 0xAA,
				I: 0x9F, R: 0x30,
				PC: 0x3D53, SP: 0x0024,
				IX: 0x7E3F, IY: 0x1584,
				AF_: 0xFC7B, BC_: 0x68F2,
				DE_: 0x616B, HL_: 0x0EF7,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{15698, 39}},
				Cycles: 4,
			},
		},
		{ // 37 0000
			name: "37 0000",
			init: z80State{
				A: 0x27, F: 0x2C,
				B: 0x05, C: 0xAD,
				D: 0x33, E: 0xE8,
				H: 0x7C, L: 0x1A,
				I: 0x4F, R: 0x2C,
				PC: 0x462C, SP: 0x7368,
				IX: 0x1EAA, IY: 0xFCC5,
				AF_: 0xD7D2, BC_: 0xB569,
				DE_: 0xD6A6, HL_: 0x31C2,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{17964, 55}},
			},
			want: z80State{
				A: 0x27, F: 0x2D,
				B: 0x05, C: 0xAD,
				D: 0x33, E: 0xE8,
				H: 0x7C, L: 0x1A,
				I: 0x4F, R: 0x2D,
				PC: 0x462D, SP: 0x7368,
				IX: 0x1EAA, IY: 0xFCC5,
				AF_: 0xD7D2, BC_: 0xB569,
				DE_: 0xD6A6, HL_: 0x31C2,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{17964, 55}},
				Cycles: 4,
			},
		},
		{ // 3F 0001 (q=0: F3/F5 from A|F)
			name: "3F 0001",
			init: z80State{
				A: 0x9E, F: 0xFC,
				B: 0x44, C: 0x68,
				D: 0x2D, E: 0x91,
				H: 0xCA, L: 0x2D,
				I: 0xC6, R: 0x23,
				PC: 0xCA70, SP: 0xE0A3,
				IX: 0x7DBC, IY: 0x11D8,
				AF_: 0x6034, BC_: 0x020F,
				DE_: 0xBE86, HL_: 0x44E8,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{51824, 63}},
			},
			want: z80State{
				A: 0x9E, F: 0xED,
				B: 0x44, C: 0x68,
				D: 0x2D, E: 0x91,
				H: 0xCA, L: 0x2D,
				I: 0xC6, R: 0x24,
				PC: 0xCA71, SP: 0xE0A3,
				IX: 0x7DBC, IY: 0x11D8,
				AF_: 0x6034, BC_: 0x020F,
				DE_: 0xBE86, HL_: 0x44E8,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{51824, 63}},
				Cycles: 4,
			},
		},
		{ // 2F 0000
			name: "2F 0000",
			init: z80State{
				A: 0xA1, F: 0x17,
				B: 0x1E, C: 0xAD,
				D: 0x10, E: 0x93,
				H: 0xEE, L: 0xA5,
				I: 0x7E, R: 0x0B,
				PC: 0x6611, SP: 0x5929,
				IX: 0x8E2F, IY: 0xC5BC,
				AF_: 0x089B, BC_: 0xD50A,
				DE_: 0x1756, HL_: 0x87C5,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{26129, 47}},
			},
			want: z80State{
				A: 0x5E, F: 0x1F,
				B: 0x1E, C: 0xAD,
				D: 0x10, E: 0x93,
				H: 0xEE, L: 0xA5,
				I: 0x7E, R: 0x0C,
				PC: 0x6612, SP: 0x5929,
				IX: 0x8E2F, IY: 0xC5BC,
				AF_: 0x089B, BC_: 0xD50A,
				DE_: 0x1756, HL_: 0x87C5,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{26129, 47}},
				Cycles: 4,
			},
		},
		{ // CE 0000
			name: "CE 0000",
			init: z80State{
				A: 0xFB, F: 0xD9,
				B: 0x79, C: 0x10,
				D: 0x16, E: 0x80,
				H: 0x76, L: 0x6F,
				I: 0xD3, R: 0x10,
				PC: 0x8DD6, SP: 0xFD46,
				IX: 0x46D1, IY: 0x1A7C,
				AF_: 0x21D6, BC_: 0xD866,
				DE_: 0xD2C1, HL_: 0x5417,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{36310, 206}, {36311, 70}},
			},
			want: z80State{
				A: 0x42, F: 0x11,
				B: 0x79, C: 0x10,
				D: 0x16, E: 0x80,
				H: 0x76, L: 0x6F,
				I: 0xD3, R: 0x11,
				PC: 0x8DD8, SP: 0xFD46,
				IX: 0x46D1, IY: 0x1A7C,
				AF_: 0x21D6, BC_: 0xD866,
				DE_: 0xD2C1, HL_: 0x5417,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{36310, 206}, {36311, 70}},
				Cycles: 7,
			},
		},
		{ // DE 0000
			name: "DE 0000",
			init: z80State{
				A: 0xDA, F: 0x55,
				B: 0x8C, C: 0x2F,
				D: 0x30, E: 0xA0,
				H: 0x23, L: 0xAD,
				I: 0x1D, R: 0x71,
				PC: 0x2368, SP: 0x1209,
				IX: 0xDBE4, IY: 0x6511,
				AF_: 0x81F5, BC_: 0xB10B,
				DE_: 0x7A0E, HL_: 0x6024,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{9064, 222}, {9065, 187}},
			},
			want: z80State{
				A: 0x1E, F: 0x1A,
				B: 0x8C, C: 0x2F,
				D: 0x30, E: 0xA0,
				H: 0x23, L: 0xAD,
				I: 0x1D, R: 0x72,
				PC: 0x236A, SP: 0x1209,
				IX: 0xDBE4, IY: 0x6511,
				AF_: 0x81F5, BC_: 0xB10B,
				DE_: 0x7A0E, HL_: 0x6024,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{9064, 222}, {9065, 187}},
				Cycles: 7,
			},
		},
		{ // D6 0000
			name: "D6 0000",
			init: z80State{
				A: 0xBE, F: 0xC6,
				B: 0x91, C: 0xF9,
				D: 0x5E, E: 0x28,
				H: 0xA1, L: 0x19,
				I: 0xF9, R: 0x4E,
				PC: 0x938B, SP: 0x247B,
				IX: 0x0966, IY: 0xA07E,
				AF_: 0x0A2E, BC_: 0x46F5,
				DE_: 0x6566, HL_: 0xB4F1,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{37771, 214}, {37772, 133}},
			},
			want: z80State{
				A: 0x39, F: 0x2A,
				B: 0x91, C: 0xF9,
				D: 0x5E, E: 0x28,
				H: 0xA1, L: 0x19,
				I: 0xF9, R: 0x4F,
				PC: 0x938D, SP: 0x247B,
				IX: 0x0966, IY: 0xA07E,
				AF_: 0x0A2E, BC_: 0x46F5,
				DE_: 0x6566, HL_: 0xB4F1,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{37771, 214}, {37772, 133}},
				Cycles: 7,
			},
		},
		{ // E6 0000
			name: "E6 0000",
			init: z80State{
				A: 0x78, F: 0x74,
				B: 0x7D, C: 0x68,
				D: 0x3C, E: 0x6D,
				H: 0x12, L: 0xC7,
				I: 0x6F, R: 0x60,
				PC: 0x7F0E, SP: 0x9CB9,
				IX: 0x423B, IY: 0x82E8,
				AF_: 0x293C, BC_: 0xEF31,
				DE_: 0x2AF0, HL_: 0xEDE4,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{32526, 230}, {32527, 39}},
			},
			want: z80State{
				A: 0x20, F: 0x30,
				B: 0x7D, C: 0x68,
				D: 0x3C, E: 0x6D,
				H: 0x12, L: 0xC7,
				I: 0x6F, R: 0x61,
				PC: 0x7F10, SP: 0x9CB9,
				IX: 0x423B, IY: 0x82E8,
				AF_: 0x293C, BC_: 0xEF31,
				DE_: 0x2AF0, HL_: 0xEDE4,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{32526, 230}, {32527, 39}},
				Cycles: 7,
			},
		},
		{ // EE 0000
			name: "EE 0000",
			init: z80State{
				A: 0x96, F: 0x74,
				B: 0xCD, C: 0x1E,
				D: 0x54, E: 0x3D,
				H: 0x56, L: 0x58,
				I: 0xFA, R: 0x78,
				PC: 0x7414, SP: 0x0012,
				IX: 0xB668, IY: 0xE370,
				AF_: 0x8EA4, BC_: 0x872D,
				DE_: 0x286A, HL_: 0xAEA0,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{29716, 238}, {29717, 240}},
			},
			want: z80State{
				A: 0x66, F: 0x24,
				B: 0xCD, C: 0x1E,
				D: 0x54, E: 0x3D,
				H: 0x56, L: 0x58,
				I: 0xFA, R: 0x79,
				PC: 0x7416, SP: 0x0012,
				IX: 0xB668, IY: 0xE370,
				AF_: 0x8EA4, BC_: 0x872D,
				DE_: 0x286A, HL_: 0xAEA0,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{29716, 238}, {29717, 240}},
				Cycles: 7,
			},
		},
		{ // F6 0000
			name: "F6 0000",
			init: z80State{
				A: 0xDA, F: 0xCF,
				B: 0xE4, C: 0x1B,
				D: 0xA8, E: 0x5F,
				H: 0x26, L: 0x69,
				I: 0x58, R: 0x6E,
				PC: 0xF6B4, SP: 0x0D96,
				IX: 0x470E, IY: 0xDFCE,
				AF_: 0xD476, BC_: 0xFD8E,
				DE_: 0xA9C6, HL_: 0x1942,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{63156, 246}, {63157, 214}},
			},
			want: z80State{
				A: 0xDE, F: 0x8C,
				B: 0xE4, C: 0x1B,
				D: 0xA8, E: 0x5F,
				H: 0x26, L: 0x69,
				I: 0x58, R: 0x6F,
				PC: 0xF6B6, SP: 0x0D96,
				IX: 0x470E, IY: 0xDFCE,
				AF_: 0xD476, BC_: 0xFD8E,
				DE_: 0xA9C6, HL_: 0x1942,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{63156, 246}, {63157, 214}},
				Cycles: 7,
			},
		},
		{ // 0C 0000
			name: "0C 0000",
			init: z80State{
				A: 0x93, F: 0xE0,
				B: 0x62, C: 0x1E,
				D: 0xA8, E: 0x4D,
				H: 0x88, L: 0x16,
				I: 0xCB, R: 0x13,
				PC: 0x2A66, SP: 0x0414,
				IX: 0x9FDF, IY: 0x0E82,
				AF_: 0x9369, BC_: 0x2795,
				DE_: 0xAAF9, HL_: 0x8CC9,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{10854, 12}},
			},
			want: z80State{
				A: 0x93, F: 0x08,
				B: 0x62, C: 0x1F,
				D: 0xA8, E: 0x4D,
				H: 0x88, L: 0x16,
				I: 0xCB, R: 0x14,
				PC: 0x2A67, SP: 0x0414,
				IX: 0x9FDF, IY: 0x0E82,
				AF_: 0x9369, BC_: 0x2795,
				DE_: 0xAAF9, HL_: 0x8CC9,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{10854, 12}},
				Cycles: 4,
			},
		},
		{ // 15 0000
			name: "15 0000",
			init: z80State{
				A: 0xD9, F: 0x27,
				B: 0x53, C: 0x38,
				D: 0x66, E: 0x0C,
				H: 0xFE, L: 0xDB,
				I: 0xD5, R: 0x6B,
				PC: 0x7679, SP: 0x197E,
				IX: 0x4BC6, IY: 0x2811,
				AF_: 0xA693, BC_: 0xA0DC,
				DE_: 0x8ADD, HL_: 0xCFED,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{30329, 21}},
			},
			want: z80State{
				A: 0xD9, F: 0x23,
				B: 0x53, C: 0x38,
				D: 0x65, E: 0x0C,
				H: 0xFE, L: 0xDB,
				I: 0xD5, R: 0x6C,
				PC: 0x767A, SP: 0x197E,
				IX: 0x4BC6, IY: 0x2811,
				AF_: 0xA693, BC_: 0xA0DC,
				DE_: 0x8ADD, HL_: 0xCFED,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{30329, 21}},
				Cycles: 4,
			},
		},
		{ // 19 0000
			name: "19 0000",
			init: z80State{
				A: 0xF7, F: 0x71,
				B: 0x15, C: 0xEA,
				D: 0x85, E: 0xA2,
				H: 0xD7, L: 0x43,
				I: 0x24, R: 0x24,
				PC: 0xFD55, SP: 0x474F,
				IX: 0xE6D8, IY: 0x67E6,
				AF_: 0x1321, BC_: 0x2AE2,
				DE_: 0x9716, HL_: 0xBC84,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{64853, 25}},
			},
			want: z80State{
				A: 0xF7, F: 0x49,
				B: 0x15, C: 0xEA,
				D: 0x85, E: 0xA2,
				H: 0x5C, L: 0xE5,
				I: 0x24, R: 0x25,
				PC: 0xFD56, SP: 0x474F,
				IX: 0xE6D8, IY: 0x67E6,
				AF_: 0x1321, BC_: 0x2AE2,
				DE_: 0x9716, HL_: 0xBC84,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{64853, 25}},
				Cycles: 11,
			},
		},
	}
	for _, tt := range tests {
		runSSTTest(t, tt)
	}
}
