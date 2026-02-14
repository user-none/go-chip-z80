package z80

import "testing"

func TestIM(t *testing.T) {
	cpu, bus := newTestCPU()
	// IM 1 (ED 56)
	bus.mem[0] = 0xED
	bus.mem[1] = 0x56
	cycles := cpu.Step()
	if cpu.reg.IM != 1 {
		t.Errorf("IM 1: IM=%d want 1", cpu.reg.IM)
	}
	if cycles != 8 {
		t.Errorf("cycles=%d want 8", cycles)
	}
}

func TestRETI(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFC
	bus.mem[0xFFFC] = 0x34
	bus.mem[0xFFFD] = 0x12
	cpu.reg.IFF2 = true
	bus.mem[0] = 0xED
	bus.mem[1] = 0x4D
	cycles := cpu.Step()
	if cpu.reg.PC != 0x1234 {
		t.Errorf("RETI: PC=%04x want 1234", cpu.reg.PC)
	}
	if !cpu.reg.IFF1 {
		t.Error("RETI should restore IFF1 from IFF2")
	}
	if cycles != 14 {
		t.Errorf("cycles=%d want 14", cycles)
	}
}

func TestRETN(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFC
	bus.mem[0xFFFC] = 0x78
	bus.mem[0xFFFD] = 0x56
	cpu.reg.IFF2 = true
	bus.mem[0] = 0xED
	bus.mem[1] = 0x45
	cycles := cpu.Step()
	if cpu.reg.PC != 0x5678 {
		t.Errorf("RETN: PC=%04x want 5678", cpu.reg.PC)
	}
	if !cpu.reg.IFF1 {
		t.Error("RETN should restore IFF1 from IFF2")
	}
	if cycles != 14 {
		t.Errorf("cycles=%d want 14", cycles)
	}
}

func TestLD_I_A(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xED
	bus.mem[1] = 0x47
	cpu.reg.AF = 0x4200
	cycles := cpu.Step()
	if cpu.reg.I != 0x42 {
		t.Errorf("LD I,A: I=%02x want 42", cpu.reg.I)
	}
	if cycles != 9 {
		t.Errorf("cycles=%d want 9", cycles)
	}
}

func TestLD_A_I(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xED
	bus.mem[1] = 0x57
	cpu.reg.I = 0x80
	cpu.reg.IFF2 = true
	cpu.reg.AF = 0x0001 // carry set
	cycles := cpu.Step()
	if cpu.getA() != 0x80 {
		t.Errorf("LD A,I: A=%02x want 80", cpu.getA())
	}
	f := cpu.getF()
	if f&flagS == 0 {
		t.Error("S should be set")
	}
	if f&flagPV == 0 {
		t.Error("PV should reflect IFF2")
	}
	if f&flagC == 0 {
		t.Error("C should be preserved")
	}
	if cycles != 9 {
		t.Errorf("cycles=%d want 9", cycles)
	}
}

func TestLD_R_A(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xED
	bus.mem[1] = 0x4F
	cpu.reg.AF = 0x5500
	cpu.Step()
	if cpu.reg.R != 0x55 {
		t.Errorf("LD R,A: R=%02x want 55", cpu.reg.R)
	}
}

func TestED_LD_nn_rr(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED 43 = LD (nn), BC
	bus.mem[0] = 0xED
	bus.mem[1] = 0x43
	bus.mem[2] = 0x00
	bus.mem[3] = 0xC0
	cpu.reg.BC = 0xBEEF
	cycles := cpu.Step()
	if bus.mem[0xC000] != 0xEF || bus.mem[0xC001] != 0xBE {
		t.Errorf("LD (nn),BC: got %02x%02x want BEEF", bus.mem[0xC001], bus.mem[0xC000])
	}
	if cycles != 20 {
		t.Errorf("cycles=%d want 20", cycles)
	}
}

func TestED_LD_rr_nn(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED 4B = LD BC, (nn)
	bus.mem[0] = 0xED
	bus.mem[1] = 0x4B
	bus.mem[2] = 0x00
	bus.mem[3] = 0xD0
	bus.mem[0xD000] = 0xAD
	bus.mem[0xD001] = 0xDE
	cycles := cpu.Step()
	if cpu.reg.BC != 0xDEAD {
		t.Errorf("LD BC,(nn): BC=%04x want DEAD", cpu.reg.BC)
	}
	if cycles != 20 {
		t.Errorf("cycles=%d want 20", cycles)
	}
}

func TestADC_HL_rr(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED 4A = ADC HL, BC
	bus.mem[0] = 0xED
	bus.mem[1] = 0x4A
	cpu.reg.HL = 0x1000
	cpu.reg.BC = 0x2000
	cpu.reg.AF = 0x0001 // carry
	cycles := cpu.Step()
	if cpu.reg.HL != 0x3001 {
		t.Errorf("ADC HL,BC: HL=%04x want 3001", cpu.reg.HL)
	}
	if cycles != 15 {
		t.Errorf("cycles=%d want 15", cycles)
	}
}

func TestSBC_HL_rr(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED 42 = SBC HL, BC
	bus.mem[0] = 0xED
	bus.mem[1] = 0x42
	cpu.reg.HL = 0x3000
	cpu.reg.BC = 0x1000
	cpu.reg.AF = 0x0001 // carry
	cycles := cpu.Step()
	if cpu.reg.HL != 0x1FFF {
		t.Errorf("SBC HL,BC: HL=%04x want 1FFF", cpu.reg.HL)
	}
	if cycles != 15 {
		t.Errorf("cycles=%d want 15", cycles)
	}
	if cpu.getF()&flagN == 0 {
		t.Error("N should be set")
	}
}

func TestSBC_HL_rr_Zero(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0xED
	bus.mem[1] = 0x42
	cpu.reg.HL = 0x1000
	cpu.reg.BC = 0x1000
	cpu.reg.AF = 0x0000 // no carry
	cpu.Step()
	if cpu.reg.HL != 0 {
		t.Errorf("SBC HL,BC: HL=%04x want 0000", cpu.reg.HL)
	}
	if cpu.getF()&flagZ == 0 {
		t.Error("Z should be set")
	}
}

func TestRLD(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED 6F = RLD
	bus.mem[0] = 0xED
	bus.mem[1] = 0x6F
	cpu.reg.AF = 0x1200
	cpu.reg.HL = 0x5000
	bus.mem[0x5000] = 0x34
	cycles := cpu.Step()
	// A low nibble=2 -> (HL) low nibble; (HL) high nibble=3 -> A low nibble
	// (HL) becomes: old_low<<4 | A_low = 0x42
	// A becomes: A_high | old_high = 0x13
	if cpu.getA() != 0x13 {
		t.Errorf("RLD: A=%02x want 13", cpu.getA())
	}
	if bus.mem[0x5000] != 0x42 {
		t.Errorf("RLD: (HL)=%02x want 42", bus.mem[0x5000])
	}
	if cycles != 18 {
		t.Errorf("cycles=%d want 18", cycles)
	}
}

func TestRRD(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED 67 = RRD
	bus.mem[0] = 0xED
	bus.mem[1] = 0x67
	cpu.reg.AF = 0x1200
	cpu.reg.HL = 0x5000
	bus.mem[0x5000] = 0x34
	cycles := cpu.Step()
	// (HL) low nibble=4 -> A low nibble
	// A low nibble=2 -> (HL) high nibble
	// (HL) high nibble=3 -> (HL) low nibble
	// A = 0x14, (HL) = 0x23
	if cpu.getA() != 0x14 {
		t.Errorf("RRD: A=%02x want 14", cpu.getA())
	}
	if bus.mem[0x5000] != 0x23 {
		t.Errorf("RRD: (HL)=%02x want 23", bus.mem[0x5000])
	}
	if cycles != 18 {
		t.Errorf("cycles=%d want 18", cycles)
	}
}

func TestIN_r_C(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED 40 = IN B, (C)
	bus.mem[0] = 0xED
	bus.mem[1] = 0x40
	cpu.reg.BC = 0x0042
	// testBus.In returns 0xFF
	cycles := cpu.Step()
	if cpu.getB() != 0xFF {
		t.Errorf("IN B,(C): B=%02x want FF", cpu.getB())
	}
	if cycles != 12 {
		t.Errorf("cycles=%d want 12", cycles)
	}
}

func TestOUT_C_r(t *testing.T) {
	cpu, bus := newTestCPU()
	_ = bus
	// ED 41 = OUT (C), B
	bus.mem[0] = 0xED
	bus.mem[1] = 0x41
	cpu.reg.BC = 0xAA42
	cycles := cpu.Step()
	// Just verify it doesn't crash and has right cycles
	if cycles != 12 {
		t.Errorf("cycles=%d want 12", cycles)
	}
}

func TestED_NOP(t *testing.T) {
	cpu, bus := newTestCPU()
	// ED 00 is an undefined ED opcode = 8-cycle NOP
	bus.mem[0] = 0xED
	bus.mem[1] = 0x00
	cycles := cpu.Step()
	if cycles != 8 {
		t.Errorf("ED NOP: cycles=%d want 8", cycles)
	}
}

func TestSST_ED(t *testing.T) {
	tests := []struct {
		name       string
		init, want z80State
	}{
		{ // ED 46 0000
			name: "ED 46 0000",
			init: z80State{
				A: 0x56, F: 0x2E,
				B: 0x4C, C: 0x1D,
				D: 0xA5, E: 0x7C,
				H: 0x8E, L: 0xAD,
				I: 0x28, R: 0x45,
				PC: 0xF898, SP: 0x726A,
				IX: 0xA417, IY: 0x002F,
				AF_: 0x452C, BC_: 0xC51D,
				DE_: 0x4DBF, HL_: 0x179C,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{63640, 237}, {63641, 70}},
			},
			want: z80State{
				A: 0x56, F: 0x2E,
				B: 0x4C, C: 0x1D,
				D: 0xA5, E: 0x7C,
				H: 0x8E, L: 0xAD,
				I: 0x28, R: 0x47,
				PC: 0xF89A, SP: 0x726A,
				IX: 0xA417, IY: 0x002F,
				AF_: 0x452C, BC_: 0xC51D,
				DE_: 0x4DBF, HL_: 0x179C,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{63640, 237}, {63641, 70}},
				Cycles: 8,
			},
		},
		{ // ED 56 0000
			name: "ED 56 0000",
			init: z80State{
				A: 0x85, F: 0x92,
				B: 0xB3, C: 0x62,
				D: 0x8E, E: 0xD2,
				H: 0x09, L: 0xD1,
				I: 0xBD, R: 0x4C,
				PC: 0x8163, SP: 0x8183,
				IX: 0xBD6F, IY: 0x6115,
				AF_: 0x5401, BC_: 0x058A,
				DE_: 0xD30F, HL_: 0xC315,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{33123, 237}, {33124, 86}},
			},
			want: z80State{
				A: 0x85, F: 0x92,
				B: 0xB3, C: 0x62,
				D: 0x8E, E: 0xD2,
				H: 0x09, L: 0xD1,
				I: 0xBD, R: 0x4E,
				PC: 0x8165, SP: 0x8183,
				IX: 0xBD6F, IY: 0x6115,
				AF_: 0x5401, BC_: 0x058A,
				DE_: 0xD30F, HL_: 0xC315,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{33123, 237}, {33124, 86}},
				Cycles: 8,
			},
		},
		{ // ED 5E 0000
			name: "ED 5E 0000",
			init: z80State{
				A: 0x3A, F: 0x2B,
				B: 0xF0, C: 0x92,
				D: 0xF4, E: 0x98,
				H: 0x46, L: 0x73,
				I: 0xCE, R: 0x42,
				PC: 0x046C, SP: 0xBD2C,
				IX: 0x16CB, IY: 0x1FD0,
				AF_: 0x50AA, BC_: 0xB51B,
				DE_: 0x7615, HL_: 0x27B3,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1132, 237}, {1133, 94}},
			},
			want: z80State{
				A: 0x3A, F: 0x2B,
				B: 0xF0, C: 0x92,
				D: 0xF4, E: 0x98,
				H: 0x46, L: 0x73,
				I: 0xCE, R: 0x44,
				PC: 0x046E, SP: 0xBD2C,
				IX: 0x16CB, IY: 0x1FD0,
				AF_: 0x50AA, BC_: 0xB51B,
				DE_: 0x7615, HL_: 0x27B3,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1132, 237}, {1133, 94}},
				Cycles: 8,
			},
		},
		{ // ED 4D 0000
			name: "ED 4D 0000",
			init: z80State{
				A: 0x83, F: 0xB6,
				B: 0x66, C: 0x1E,
				D: 0x11, E: 0xD0,
				H: 0x91, L: 0x2F,
				I: 0xB8, R: 0x47,
				PC: 0xFA82, SP: 0x04CE,
				IX: 0xE951, IY: 0xCB72,
				AF_: 0xE75B, BC_: 0x0326,
				DE_: 0xA207, HL_: 0x8547,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{1230, 188}, {1231, 80}, {64130, 237}, {64131, 77}},
			},
			want: z80State{
				A: 0x83, F: 0xB6,
				B: 0x66, C: 0x1E,
				D: 0x11, E: 0xD0,
				H: 0x91, L: 0x2F,
				I: 0xB8, R: 0x49,
				PC: 0x50BC, SP: 0x04D0,
				IX: 0xE951, IY: 0xCB72,
				AF_: 0xE75B, BC_: 0x0326,
				DE_: 0xA207, HL_: 0x8547,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{1230, 188}, {1231, 80}, {64130, 237}, {64131, 77}},
				Cycles: 14,
			},
		},
		{ // ED 45 0000
			name: "ED 45 0000",
			init: z80State{
				A: 0xAE, F: 0x21,
				B: 0x71, C: 0x9C,
				D: 0x5E, E: 0x5E,
				H: 0x09, L: 0x86,
				I: 0x6E, R: 0x69,
				PC: 0x49BC, SP: 0x28D5,
				IX: 0x147A, IY: 0xCF5D,
				AF_: 0x3527, BC_: 0x9136,
				DE_: 0xFCBB, HL_: 0x2A3D,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{10453, 57}, {10454, 137}, {18876, 237}, {18877, 69}},
			},
			want: z80State{
				A: 0xAE, F: 0x21,
				B: 0x71, C: 0x9C,
				D: 0x5E, E: 0x5E,
				H: 0x09, L: 0x86,
				I: 0x6E, R: 0x6B,
				PC: 0x8939, SP: 0x28D7,
				IX: 0x147A, IY: 0xCF5D,
				AF_: 0x3527, BC_: 0x9136,
				DE_: 0xFCBB, HL_: 0x2A3D,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{10453, 57}, {10454, 137}, {18876, 237}, {18877, 69}},
				Cycles: 14,
			},
		},
		{ // ED 47 0000
			name: "ED 47 0000",
			init: z80State{
				A: 0x11, F: 0xD3,
				B: 0x4D, C: 0x86,
				D: 0x9B, E: 0x25,
				H: 0x4E, L: 0x45,
				I: 0x03, R: 0x18,
				PC: 0x2317, SP: 0x4F9D,
				IX: 0x7DC4, IY: 0x83CD,
				AF_: 0x747C, BC_: 0x19F4,
				DE_: 0xAFF9, HL_: 0x5D3C,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{8983, 237}, {8984, 71}},
			},
			want: z80State{
				A: 0x11, F: 0xD3,
				B: 0x4D, C: 0x86,
				D: 0x9B, E: 0x25,
				H: 0x4E, L: 0x45,
				I: 0x11, R: 0x1A,
				PC: 0x2319, SP: 0x4F9D,
				IX: 0x7DC4, IY: 0x83CD,
				AF_: 0x747C, BC_: 0x19F4,
				DE_: 0xAFF9, HL_: 0x5D3C,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{8983, 237}, {8984, 71}},
				Cycles: 9,
			},
		},
		{ // ED 57 0000
			name: "ED 57 0000",
			init: z80State{
				A: 0x4C, F: 0x13,
				B: 0x23, C: 0x8E,
				D: 0x82, E: 0xD1,
				H: 0x43, L: 0xE0,
				I: 0x98, R: 0x7F,
				PC: 0x68E8, SP: 0xA3F2,
				IX: 0x4C43, IY: 0xC90A,
				AF_: 0x81FA, BC_: 0x5812,
				DE_: 0x1639, HL_: 0xC702,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{26856, 237}, {26857, 87}},
			},
			want: z80State{
				A: 0x98, F: 0x8D,
				B: 0x23, C: 0x8E,
				D: 0x82, E: 0xD1,
				H: 0x43, L: 0xE0,
				I: 0x98, R: 0x01,
				PC: 0x68EA, SP: 0xA3F2,
				IX: 0x4C43, IY: 0xC90A,
				AF_: 0x81FA, BC_: 0x5812,
				DE_: 0x1639, HL_: 0xC702,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{26856, 237}, {26857, 87}},
				Cycles: 9,
			},
		},
		{ // ED 4F 0000
			name: "ED 4F 0000",
			init: z80State{
				A: 0x0B, F: 0xE3,
				B: 0xE4, C: 0xDD,
				D: 0x2F, E: 0x5B,
				H: 0xB0, L: 0x24,
				I: 0xF7, R: 0x3E,
				PC: 0x8256, SP: 0x14F2,
				IX: 0x5639, IY: 0x774C,
				AF_: 0x3976, BC_: 0x0970,
				DE_: 0x6FF2, HL_: 0x5F34,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{33366, 237}, {33367, 79}},
			},
			want: z80State{
				A: 0x0B, F: 0xE3,
				B: 0xE4, C: 0xDD,
				D: 0x2F, E: 0x5B,
				H: 0xB0, L: 0x24,
				I: 0xF7, R: 0x0B,
				PC: 0x8258, SP: 0x14F2,
				IX: 0x5639, IY: 0x774C,
				AF_: 0x3976, BC_: 0x0970,
				DE_: 0x6FF2, HL_: 0x5F34,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{33366, 237}, {33367, 79}},
				Cycles: 9,
			},
		},
		{ // ED 5F 0000
			name: "ED 5F 0000",
			init: z80State{
				A: 0x36, F: 0x12,
				B: 0xD6, C: 0xC7,
				D: 0x6A, E: 0xCF,
				H: 0xD6, L: 0x50,
				I: 0x7A, R: 0x35,
				PC: 0x1044, SP: 0xFF8F,
				IX: 0x01D0, IY: 0x6F70,
				AF_: 0x887B, BC_: 0xB6BC,
				DE_: 0xE736, HL_: 0x3299,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{4164, 237}, {4165, 95}},
			},
			want: z80State{
				A: 0x37, F: 0x20,
				B: 0xD6, C: 0xC7,
				D: 0x6A, E: 0xCF,
				H: 0xD6, L: 0x50,
				I: 0x7A, R: 0x37,
				PC: 0x1046, SP: 0xFF8F,
				IX: 0x01D0, IY: 0x6F70,
				AF_: 0x887B, BC_: 0xB6BC,
				DE_: 0xE736, HL_: 0x3299,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{4164, 237}, {4165, 95}},
				Cycles: 9,
			},
		},
		{ // ED 44 0000
			name: "ED 44 0000",
			init: z80State{
				A: 0x52, F: 0xFF,
				B: 0x30, C: 0x2E,
				D: 0x05, E: 0x47,
				H: 0xAE, L: 0x00,
				I: 0x55, R: 0x31,
				PC: 0xA811, SP: 0xCDC6,
				IX: 0x2B80, IY: 0x1761,
				AF_: 0x13B5, BC_: 0x7449,
				DE_: 0x4620, HL_: 0xCEBB,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{43025, 237}, {43026, 68}},
			},
			want: z80State{
				A: 0xAE, F: 0xBB,
				B: 0x30, C: 0x2E,
				D: 0x05, E: 0x47,
				H: 0xAE, L: 0x00,
				I: 0x55, R: 0x33,
				PC: 0xA813, SP: 0xCDC6,
				IX: 0x2B80, IY: 0x1761,
				AF_: 0x13B5, BC_: 0x7449,
				DE_: 0x4620, HL_: 0xCEBB,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{43025, 237}, {43026, 68}},
				Cycles: 8,
			},
		},
		{ // ED 42 0000
			name: "ED 42 0000",
			init: z80State{
				A: 0x59, F: 0x47,
				B: 0x6D, C: 0xFF,
				D: 0x32, E: 0xD6,
				H: 0xFB, L: 0x45,
				I: 0xCE, R: 0x76,
				PC: 0x5B5F, SP: 0xFD3F,
				IX: 0x6DA6, IY: 0xA440,
				AF_: 0xD08E, BC_: 0xFECF,
				DE_: 0x06E4, HL_: 0x5479,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{23391, 237}, {23392, 66}},
			},
			want: z80State{
				A: 0x59, F: 0x9A,
				B: 0x6D, C: 0xFF,
				D: 0x32, E: 0xD6,
				H: 0x8D, L: 0x45,
				I: 0xCE, R: 0x78,
				PC: 0x5B61, SP: 0xFD3F,
				IX: 0x6DA6, IY: 0xA440,
				AF_: 0xD08E, BC_: 0xFECF,
				DE_: 0x06E4, HL_: 0x5479,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{23391, 237}, {23392, 66}},
				Cycles: 15,
			},
		},
		{ // ED 52 0000
			name: "ED 52 0000",
			init: z80State{
				A: 0x34, F: 0x2C,
				B: 0x86, C: 0x8B,
				D: 0x0A, E: 0xDA,
				H: 0x6A, L: 0xF1,
				I: 0x76, R: 0x50,
				PC: 0x61D4, SP: 0x9D46,
				IX: 0xC148, IY: 0xC68C,
				AF_: 0x2EB5, BC_: 0xBBAC,
				DE_: 0x4A40, HL_: 0xF23F,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{25044, 237}, {25045, 82}},
			},
			want: z80State{
				A: 0x34, F: 0x22,
				B: 0x86, C: 0x8B,
				D: 0x0A, E: 0xDA,
				H: 0x60, L: 0x17,
				I: 0x76, R: 0x52,
				PC: 0x61D6, SP: 0x9D46,
				IX: 0xC148, IY: 0xC68C,
				AF_: 0x2EB5, BC_: 0xBBAC,
				DE_: 0x4A40, HL_: 0xF23F,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{25044, 237}, {25045, 82}},
				Cycles: 15,
			},
		},
		{ // ED 4A 0000
			name: "ED 4A 0000",
			init: z80State{
				A: 0x18, F: 0x1F,
				B: 0xB0, C: 0xCF,
				D: 0x44, E: 0x64,
				H: 0x37, L: 0xAD,
				I: 0x92, R: 0x07,
				PC: 0xA8DD, SP: 0xDEE1,
				IX: 0x4EBD, IY: 0x4D7C,
				AF_: 0xAD44, BC_: 0x433F,
				DE_: 0xCE7A, HL_: 0xFF79,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{43229, 237}, {43230, 74}},
			},
			want: z80State{
				A: 0x18, F: 0xA8,
				B: 0xB0, C: 0xCF,
				D: 0x44, E: 0x64,
				H: 0xE8, L: 0x7D,
				I: 0x92, R: 0x09,
				PC: 0xA8DF, SP: 0xDEE1,
				IX: 0x4EBD, IY: 0x4D7C,
				AF_: 0xAD44, BC_: 0x433F,
				DE_: 0xCE7A, HL_: 0xFF79,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{43229, 237}, {43230, 74}},
				Cycles: 15,
			},
		},
		{ // ED 6F 0000
			name: "ED 6F 0000",
			init: z80State{
				A: 0x9F, F: 0x86,
				B: 0xC3, C: 0xF1,
				D: 0x40, E: 0xDC,
				H: 0xA6, L: 0x81,
				I: 0x28, R: 0x46,
				PC: 0xA11D, SP: 0x55E3,
				IX: 0x550E, IY: 0x5D9F,
				AF_: 0x74D3, BC_: 0x3401,
				DE_: 0xF97D, HL_: 0x1186,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{41245, 237}, {41246, 111}, {42625, 131}},
			},
			want: z80State{
				A: 0x98, F: 0x88,
				B: 0xC3, C: 0xF1,
				D: 0x40, E: 0xDC,
				H: 0xA6, L: 0x81,
				I: 0x28, R: 0x48,
				PC: 0xA11F, SP: 0x55E3,
				IX: 0x550E, IY: 0x5D9F,
				AF_: 0x74D3, BC_: 0x3401,
				DE_: 0xF97D, HL_: 0x1186,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{41245, 237}, {41246, 111}, {42625, 63}},
				Cycles: 18,
			},
		},
		{ // ED 67 0000
			name: "ED 67 0000",
			init: z80State{
				A: 0x42, F: 0x44,
				B: 0xB7, C: 0x10,
				D: 0x19, E: 0x5B,
				H: 0x09, L: 0x8F,
				I: 0x01, R: 0x04,
				PC: 0xF938, SP: 0xBCCF,
				IX: 0x2B97, IY: 0xC662,
				AF_: 0x1FDE, BC_: 0xEF27,
				DE_: 0xD65B, HL_: 0x6998,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{2447, 56}, {63800, 237}, {63801, 103}},
			},
			want: z80State{
				A: 0x48, F: 0x0C,
				B: 0xB7, C: 0x10,
				D: 0x19, E: 0x5B,
				H: 0x09, L: 0x8F,
				I: 0x01, R: 0x06,
				PC: 0xF93A, SP: 0xBCCF,
				IX: 0x2B97, IY: 0xC662,
				AF_: 0x1FDE, BC_: 0xEF27,
				DE_: 0xD65B, HL_: 0x6998,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{2447, 35}, {63800, 237}, {63801, 103}},
				Cycles: 18,
			},
		},
		{ // ED 43 0000
			name: "ED 43 0000",
			init: z80State{
				A: 0x25, F: 0xE0,
				B: 0xC2, C: 0x88,
				D: 0xE9, E: 0x78,
				H: 0xD8, L: 0x66,
				I: 0x95, R: 0x1A,
				PC: 0x8594, SP: 0xD224,
				IX: 0x724E, IY: 0xB466,
				AF_: 0xEAD7, BC_: 0xE556,
				DE_: 0xFAF8, HL_: 0x8CB7,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{34196, 237}, {34197, 67}, {34198, 235}, {34199, 136}, {35051, 0}, {35052, 0}},
			},
			want: z80State{
				A: 0x25, F: 0xE0,
				B: 0xC2, C: 0x88,
				D: 0xE9, E: 0x78,
				H: 0xD8, L: 0x66,
				I: 0x95, R: 0x1C,
				PC: 0x8598, SP: 0xD224,
				IX: 0x724E, IY: 0xB466,
				AF_: 0xEAD7, BC_: 0xE556,
				DE_: 0xFAF8, HL_: 0x8CB7,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{34196, 237}, {34197, 67}, {34198, 235}, {34199, 136}, {35051, 136}, {35052, 194}},
				Cycles: 20,
			},
		},
		{ // ED 4B 0000
			name: "ED 4B 0000",
			init: z80State{
				A: 0x55, F: 0xDB,
				B: 0x2A, C: 0xF8,
				D: 0xE5, E: 0xFD,
				H: 0xD0, L: 0x6A,
				I: 0x4C, R: 0x69,
				PC: 0x646C, SP: 0x4D3E,
				IX: 0x6A74, IY: 0x0CDD,
				AF_: 0x93EA, BC_: 0x0177,
				DE_: 0x61D7, HL_: 0x9F51,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{25708, 237}, {25709, 75}, {25710, 166}, {25711, 188}, {48294, 42}, {48295, 70}},
			},
			want: z80State{
				A: 0x55, F: 0xDB,
				B: 0x46, C: 0x2A,
				D: 0xE5, E: 0xFD,
				H: 0xD0, L: 0x6A,
				I: 0x4C, R: 0x6B,
				PC: 0x6470, SP: 0x4D3E,
				IX: 0x6A74, IY: 0x0CDD,
				AF_: 0x93EA, BC_: 0x0177,
				DE_: 0x61D7, HL_: 0x9F51,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{25708, 237}, {25709, 75}, {25710, 166}, {25711, 188}, {48294, 42}, {48295, 70}},
				Cycles: 20,
			},
		},
		{ // ED 40 0000
			name: "ED 40 0000",
			init: z80State{
				A: 0x14, F: 0xBC,
				B: 0x43, C: 0x81,
				D: 0xFC, E: 0xAB,
				H: 0x25, L: 0xBC,
				I: 0xE8, R: 0x4C,
				PC: 0x3F78, SP: 0xD060,
				IX: 0x6126, IY: 0x27C1,
				AF_: 0xBBB3, BC_: 0xC2CF,
				DE_: 0xBCA3, HL_: 0x33C7,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{16248, 237}, {16249, 64}},
				Ports: [][2]uint16{{17281, 210}},
			},
			want: z80State{
				A: 0x14, F: 0x84,
				B: 0xD2, C: 0x81,
				D: 0xFC, E: 0xAB,
				H: 0x25, L: 0xBC,
				I: 0xE8, R: 0x4E,
				PC: 0x3F7A, SP: 0xD060,
				IX: 0x6126, IY: 0x27C1,
				AF_: 0xBBB3, BC_: 0xC2CF,
				DE_: 0xBCA3, HL_: 0x33C7,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{16248, 237}, {16249, 64}},
				Cycles: 12,
			},
		},
		{ // ED 79 0000
			name: "ED 79 0000",
			init: z80State{
				A: 0x96, F: 0x21,
				B: 0xB2, C: 0x09,
				D: 0x9E, E: 0xA8,
				H: 0x5E, L: 0xD4,
				I: 0x02, R: 0x42,
				PC: 0x600B, SP: 0xA2D7,
				IX: 0xD2D6, IY: 0xB20C,
				AF_: 0xB2C2, BC_: 0x1FF5,
				DE_: 0x2E79, HL_: 0xDAE4,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{24587, 237}, {24588, 121}},
			},
			want: z80State{
				A: 0x96, F: 0x21,
				B: 0xB2, C: 0x09,
				D: 0x9E, E: 0xA8,
				H: 0x5E, L: 0xD4,
				I: 0x02, R: 0x44,
				PC: 0x600D, SP: 0xA2D7,
				IX: 0xD2D6, IY: 0xB20C,
				AF_: 0xB2C2, BC_: 0x1FF5,
				DE_: 0x2E79, HL_: 0xDAE4,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{24587, 237}, {24588, 121}},
				Cycles: 12,
			},
		},
		{ // ED 70 0000
			name: "ED 70 0000",
			init: z80State{
				A: 0x3F, F: 0x06,
				B: 0x02, C: 0x84,
				D: 0xDA, E: 0xA3,
				H: 0x21, L: 0x37,
				I: 0xB8, R: 0x23,
				PC: 0x908E, SP: 0x961A,
				IX: 0xB784, IY: 0x80EA,
				AF_: 0x1C9F, BC_: 0x26DE,
				DE_: 0x3E15, HL_: 0xCA2D,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{37006, 237}, {37007, 112}},
				Ports: [][2]uint16{{644, 176}},
			},
			want: z80State{
				A: 0x3F, F: 0xA0,
				B: 0x02, C: 0x84,
				D: 0xDA, E: 0xA3,
				H: 0x21, L: 0x37,
				I: 0xB8, R: 0x25,
				PC: 0x9090, SP: 0x961A,
				IX: 0xB784, IY: 0x80EA,
				AF_: 0x1C9F, BC_: 0x26DE,
				DE_: 0x3E15, HL_: 0xCA2D,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{37006, 237}, {37007, 112}},
				Cycles: 12,
			},
		},
		{ // ED 48 0000 — IN C,(C)
			name: "ED 48 0000",
			init: z80State{
				A: 0x36, F: 0xDB,
				B: 0x13, C: 0xDC,
				D: 0xA0, E: 0x78,
				H: 0xD0, L: 0xE0,
				I: 0xE0, R: 0x7A,
				PC: 0xD3FF, SP: 0x3C9C,
				IX: 0xC180, IY: 0x1028,
				AF_: 0xCEDF, BC_: 0x458C,
				DE_: 0x0956, HL_: 0x32F5,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{54271, 237}, {54272, 72}},
				Ports: [][2]uint16{{5084, 12}},
			},
			want: z80State{
				A: 0x36, F: 0x0D,
				B: 0x13, C: 0x0C,
				D: 0xA0, E: 0x78,
				H: 0xD0, L: 0xE0,
				I: 0xE0, R: 0x7C,
				PC: 0xD401, SP: 0x3C9C,
				IX: 0xC180, IY: 0x1028,
				AF_: 0xCEDF, BC_: 0x458C,
				DE_: 0x0956, HL_: 0x32F5,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{54271, 237}, {54272, 72}},
				Cycles: 12,
			},
		},
		{ // ED 50 0000 — IN D,(C)
			name: "ED 50 0000",
			init: z80State{
				A: 0x60, F: 0x1D,
				B: 0xC3, C: 0xE9,
				D: 0x29, E: 0xB0,
				H: 0x7F, L: 0xFE,
				I: 0x8C, R: 0x35,
				PC: 0xF250, SP: 0xAC37,
				IX: 0x1673, IY: 0xE794,
				AF_: 0x294C, BC_: 0x14FB,
				DE_: 0x7F5A, HL_: 0xDC59,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{62032, 237}, {62033, 80}},
				Ports: [][2]uint16{{50153, 133}},
			},
			want: z80State{
				A: 0x60, F: 0x81,
				B: 0xC3, C: 0xE9,
				D: 0x85, E: 0xB0,
				H: 0x7F, L: 0xFE,
				I: 0x8C, R: 0x37,
				PC: 0xF252, SP: 0xAC37,
				IX: 0x1673, IY: 0xE794,
				AF_: 0x294C, BC_: 0x14FB,
				DE_: 0x7F5A, HL_: 0xDC59,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{62032, 237}, {62033, 80}},
				Cycles: 12,
			},
		},
		{ // ED 41 0000 — OUT (C),B
			name: "ED 41 0000",
			init: z80State{
				A: 0x5C, F: 0x4B,
				B: 0x98, C: 0x1C,
				D: 0xF1, E: 0x9E,
				H: 0xED, L: 0x99,
				I: 0x35, R: 0x40,
				PC: 0x3BE3, SP: 0xD45F,
				IX: 0xAD73, IY: 0xF3B2,
				AF_: 0xD1C7, BC_: 0x0938,
				DE_: 0x1A5C, HL_: 0x3E4F,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{15331, 237}, {15332, 65}},
			},
			want: z80State{
				A: 0x5C, F: 0x4B,
				B: 0x98, C: 0x1C,
				D: 0xF1, E: 0x9E,
				H: 0xED, L: 0x99,
				I: 0x35, R: 0x42,
				PC: 0x3BE5, SP: 0xD45F,
				IX: 0xAD73, IY: 0xF3B2,
				AF_: 0xD1C7, BC_: 0x0938,
				DE_: 0x1A5C, HL_: 0x3E4F,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{15331, 237}, {15332, 65}},
				Cycles: 12,
			},
		},
		{ // ED 51 0000 — OUT (C),D
			name: "ED 51 0000",
			init: z80State{
				A: 0x16, F: 0x35,
				B: 0xAD, C: 0x57,
				D: 0x19, E: 0x27,
				H: 0xD8, L: 0x7D,
				I: 0x68, R: 0x2C,
				PC: 0xE29C, SP: 0x3B51,
				IX: 0x123C, IY: 0x3EC7,
				AF_: 0x194F, BC_: 0x999B,
				DE_: 0x5F3D, HL_: 0x8CC7,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{58012, 237}, {58013, 81}},
			},
			want: z80State{
				A: 0x16, F: 0x35,
				B: 0xAD, C: 0x57,
				D: 0x19, E: 0x27,
				H: 0xD8, L: 0x7D,
				I: 0x68, R: 0x2E,
				PC: 0xE29E, SP: 0x3B51,
				IX: 0x123C, IY: 0x3EC7,
				AF_: 0x194F, BC_: 0x999B,
				DE_: 0x5F3D, HL_: 0x8CC7,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{58012, 237}, {58013, 81}},
				Cycles: 12,
			},
		},
		{ // ED 5A 0000
			name: "ED 5A 0000",
			init: z80State{
				A: 0x16, F: 0xF2,
				B: 0x18, C: 0xC8,
				D: 0x51, E: 0x3D,
				H: 0x67, L: 0x7B,
				I: 0xBE, R: 0x7A,
				PC: 0x50D2, SP: 0x771C,
				IX: 0x7C6C, IY: 0x497E,
				AF_: 0x9D25, BC_: 0x1926,
				DE_: 0xCAB5, HL_: 0x7D00,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{20690, 237}, {20691, 90}},
			},
			want: z80State{
				A: 0x16, F: 0xAC,
				B: 0x18, C: 0xC8,
				D: 0x51, E: 0x3D,
				H: 0xB8, L: 0xB8,
				I: 0xBE, R: 0x7C,
				PC: 0x50D4, SP: 0x771C,
				IX: 0x7C6C, IY: 0x497E,
				AF_: 0x9D25, BC_: 0x1926,
				DE_: 0xCAB5, HL_: 0x7D00,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{20690, 237}, {20691, 90}},
				Cycles: 15,
			},
		},
		{ // ED 53 0000
			name: "ED 53 0000",
			init: z80State{
				A: 0x19, F: 0x84,
				B: 0xC5, C: 0x33,
				D: 0x0D, E: 0xEC,
				H: 0x9A, L: 0xAB,
				I: 0x98, R: 0x76,
				PC: 0x7E39, SP: 0xD7DE,
				IX: 0xEC0C, IY: 0x0046,
				AF_: 0x131D, BC_: 0x7173,
				DE_: 0x0ABB, HL_: 0x804A,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{32313, 237}, {32314, 83}, {32315, 239}, {32316, 203}, {52207, 0}, {52208, 0}},
			},
			want: z80State{
				A: 0x19, F: 0x84,
				B: 0xC5, C: 0x33,
				D: 0x0D, E: 0xEC,
				H: 0x9A, L: 0xAB,
				I: 0x98, R: 0x78,
				PC: 0x7E3D, SP: 0xD7DE,
				IX: 0xEC0C, IY: 0x0046,
				AF_: 0x131D, BC_: 0x7173,
				DE_: 0x0ABB, HL_: 0x804A,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{32313, 237}, {32314, 83}, {32315, 239}, {32316, 203}, {52207, 236}, {52208, 13}},
				Cycles: 20,
			},
		},
		{ // ED 5B 0000
			name: "ED 5B 0000",
			init: z80State{
				A: 0xFA, F: 0xC5,
				B: 0x27, C: 0x21,
				D: 0x92, E: 0x22,
				H: 0xC2, L: 0x8B,
				I: 0xAA, R: 0x12,
				PC: 0xCEE1, SP: 0xE218,
				IX: 0x10E0, IY: 0xAB83,
				AF_: 0xD74C, BC_: 0x427F,
				DE_: 0x7295, HL_: 0xF204,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{18635, 145}, {18636, 38}, {52961, 237}, {52962, 91}, {52963, 203}, {52964, 72}},
			},
			want: z80State{
				A: 0xFA, F: 0xC5,
				B: 0x27, C: 0x21,
				D: 0x26, E: 0x91,
				H: 0xC2, L: 0x8B,
				I: 0xAA, R: 0x14,
				PC: 0xCEE5, SP: 0xE218,
				IX: 0x10E0, IY: 0xAB83,
				AF_: 0xD74C, BC_: 0x427F,
				DE_: 0x7295, HL_: 0xF204,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{18635, 145}, {18636, 38}, {52961, 237}, {52962, 91}, {52963, 203}, {52964, 72}},
				Cycles: 20,
			},
		},
	}
	for _, tt := range tests {
		runSSTTest(t, tt)
	}
}
