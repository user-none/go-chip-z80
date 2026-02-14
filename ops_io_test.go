package z80

import "testing"

// ioBus records I/O operations.
type ioBus struct {
	testBus
	lastInPort  uint16
	lastOutPort uint16
	lastOutVal  uint8
	inVal       uint8
}

func (b *ioBus) In(port uint16) uint8  { b.lastInPort = port; return b.inVal }
func (b *ioBus) Out(port uint16, val uint8) {
	b.lastOutPort = port
	b.lastOutVal = val
}

func TestIN_A_n(t *testing.T) {
	bus := &ioBus{inVal: 0x42}
	cpu := New(bus)
	bus.mem[0] = 0xDB // IN A, (n)
	bus.mem[1] = 0x10 // port 0x10
	cpu.reg.AF = 0xFF00
	cycles := cpu.Step()
	if cpu.getA() != 0x42 {
		t.Errorf("IN A,(n): A=%02x want 42", cpu.getA())
	}
	// Port should be A<<8 | n = 0xFF10
	if bus.lastInPort != 0xFF10 {
		t.Errorf("port=%04x want FF10", bus.lastInPort)
	}
	if cycles != 11 {
		t.Errorf("cycles=%d want 11", cycles)
	}
}

func TestOUT_n_A(t *testing.T) {
	bus := &ioBus{}
	cpu := New(bus)
	bus.mem[0] = 0xD3 // OUT (n), A
	bus.mem[1] = 0x20
	cpu.reg.AF = 0xAB00
	cycles := cpu.Step()
	if bus.lastOutVal != 0xAB {
		t.Errorf("OUT (n),A: val=%02x want AB", bus.lastOutVal)
	}
	if bus.lastOutPort != 0xAB20 {
		t.Errorf("port=%04x want AB20", bus.lastOutPort)
	}
	if cycles != 11 {
		t.Errorf("cycles=%d want 11", cycles)
	}
}

func TestINI(t *testing.T) {
	bus := &ioBus{inVal: 0x55}
	cpu := New(bus)
	bus.mem[0] = 0xED
	bus.mem[1] = 0xA2 // INI
	cpu.reg.BC = 0x0310 // B=3, C=0x10
	cpu.reg.HL = 0x8000
	cycles := cpu.Step()
	if bus.mem[0x8000] != 0x55 {
		t.Errorf("INI: (HL)=%02x want 55", bus.mem[0x8000])
	}
	if cpu.getB() != 2 {
		t.Errorf("INI: B=%02x want 02", cpu.getB())
	}
	if cpu.reg.HL != 0x8001 {
		t.Errorf("INI: HL=%04x want 8001", cpu.reg.HL)
	}
	if cycles != 16 {
		t.Errorf("cycles=%d want 16", cycles)
	}
}

func TestIND(t *testing.T) {
	bus := &ioBus{inVal: 0x77}
	cpu := New(bus)
	bus.mem[0] = 0xED
	bus.mem[1] = 0xAA // IND
	cpu.reg.BC = 0x0210
	cpu.reg.HL = 0x8005
	cpu.Step()
	if bus.mem[0x8005] != 0x77 {
		t.Errorf("IND: (HL)=%02x want 77", bus.mem[0x8005])
	}
	if cpu.reg.HL != 0x8004 {
		t.Errorf("IND: HL=%04x want 8004", cpu.reg.HL)
	}
}

func TestINIR(t *testing.T) {
	bus := &ioBus{inVal: 0xAA}
	cpu := New(bus)
	cpu.reg.PC = 0x0100
	bus.mem[0x0100] = 0xED
	bus.mem[0x0101] = 0xB2 // INIR
	cpu.reg.BC = 0x0210
	cpu.reg.HL = 0x8000

	c1 := cpu.Step() // B=2->1, repeat
	if c1 != 21 {
		t.Errorf("INIR iter 1: cycles=%d want 21", c1)
	}
	if cpu.reg.PC != 0x0100 {
		t.Errorf("INIR should repeat: PC=%04x", cpu.reg.PC)
	}

	c2 := cpu.Step() // B=1->0, done
	if c2 != 16 {
		t.Errorf("INIR final: cycles=%d want 16", c2)
	}
	if cpu.getB() != 0 {
		t.Errorf("INIR done: B=%02x want 0", cpu.getB())
	}
	if cpu.getF()&flagZ == 0 {
		t.Error("Z should be set (B=0)")
	}
}

func TestOUTI(t *testing.T) {
	bus := &ioBus{}
	cpu := New(bus)
	bus.mem[0] = 0xED
	bus.mem[1] = 0xA3 // OUTI
	cpu.reg.BC = 0x0310
	cpu.reg.HL = 0x8000
	bus.mem[0x8000] = 0x99
	cycles := cpu.Step()
	if bus.lastOutVal != 0x99 {
		t.Errorf("OUTI: val=%02x want 99", bus.lastOutVal)
	}
	if cpu.getB() != 2 {
		t.Errorf("OUTI: B=%02x want 02", cpu.getB())
	}
	if cpu.reg.HL != 0x8001 {
		t.Errorf("OUTI: HL=%04x want 8001", cpu.reg.HL)
	}
	if cycles != 16 {
		t.Errorf("cycles=%d want 16", cycles)
	}
}

func TestOTIR(t *testing.T) {
	bus := &ioBus{}
	cpu := New(bus)
	cpu.reg.PC = 0x0200
	bus.mem[0x0200] = 0xED
	bus.mem[0x0201] = 0xB3 // OTIR
	cpu.reg.BC = 0x0210
	cpu.reg.HL = 0x8000
	bus.mem[0x8000] = 0x11
	bus.mem[0x8001] = 0x22

	cpu.Step() // B=2->1, repeat
	cpu.Step() // B=1->0, done
	if cpu.getB() != 0 {
		t.Errorf("OTIR: B=%02x want 0", cpu.getB())
	}
}

func TestSST_IO(t *testing.T) {
	tests := []struct {
		name       string
		init, want z80State
	}{
		{ // DB 0000
			name: "DB 0000",
			init: z80State{
				A: 0xE3, F: 0x8C,
				B: 0xAD, C: 0xF9,
				D: 0x8E, E: 0xEE,
				H: 0x15, L: 0xBE,
				I: 0x3E, R: 0x00,
				PC: 0xC26E, SP: 0x1E55,
				IX: 0x73DD, IY: 0x1558,
				AF_: 0x6D53, BC_: 0x9210,
				DE_: 0xCCB2, HL_: 0xA36A,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{49774, 219}, {49775, 249}},
				Ports: [][2]uint16{{58361, 155}},
			},
			want: z80State{
				A: 0x9B, F: 0x8C,
				B: 0xAD, C: 0xF9,
				D: 0x8E, E: 0xEE,
				H: 0x15, L: 0xBE,
				I: 0x3E, R: 0x01,
				PC: 0xC270, SP: 0x1E55,
				IX: 0x73DD, IY: 0x1558,
				AF_: 0x6D53, BC_: 0x9210,
				DE_: 0xCCB2, HL_: 0xA36A,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{49774, 219}, {49775, 249}},
				Cycles: 11,
			},
		},
		{ // D3 0000
			name: "D3 0000",
			init: z80State{
				A: 0x66, F: 0x54,
				B: 0x78, C: 0xA7,
				D: 0x1F, E: 0xF9,
				H: 0x21, L: 0x02,
				I: 0xE4, R: 0x7E,
				PC: 0x95E3, SP: 0xA2E3,
				IX: 0xCF25, IY: 0xA534,
				AF_: 0x1939, BC_: 0x774D,
				DE_: 0x3F83, HL_: 0x08A0,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{38371, 211}, {38372, 159}},
			},
			want: z80State{
				A: 0x66, F: 0x54,
				B: 0x78, C: 0xA7,
				D: 0x1F, E: 0xF9,
				H: 0x21, L: 0x02,
				I: 0xE4, R: 0x7F,
				PC: 0x95E5, SP: 0xA2E3,
				IX: 0xCF25, IY: 0xA534,
				AF_: 0x1939, BC_: 0x774D,
				DE_: 0x3F83, HL_: 0x08A0,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{38371, 211}, {38372, 159}},
				Cycles: 11,
			},
		},
		{ // ED A2 0000
			name: "ED A2 0000",
			init: z80State{
				A: 0xC7, F: 0x34,
				B: 0x66, C: 0x62,
				D: 0xCF, E: 0xCF,
				H: 0x55, L: 0x0F,
				I: 0x1C, R: 0x52,
				PC: 0xD3A9, SP: 0x31DC,
				IX: 0x1062, IY: 0x2313,
				AF_: 0x681F, BC_: 0x6BDE,
				DE_: 0x0D61, HL_: 0x6D49,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{21775, 0}, {54185, 237}, {54186, 162}},
				Ports: [][2]uint16{{26210, 91}},
			},
			want: z80State{
				A: 0xC7, F: 0x24,
				B: 0x65, C: 0x62,
				D: 0xCF, E: 0xCF,
				H: 0x55, L: 0x10,
				I: 0x1C, R: 0x54,
				PC: 0xD3AB, SP: 0x31DC,
				IX: 0x1062, IY: 0x2313,
				AF_: 0x681F, BC_: 0x6BDE,
				DE_: 0x0D61, HL_: 0x6D49,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{21775, 91}, {54185, 237}, {54186, 162}},
				Cycles: 16,
			},
		},
		{ // ED A2 0001
			name: "ED A2 0001",
			init: z80State{
				A: 0x51, F: 0xDA,
				B: 0x40, C: 0xED,
				D: 0x4F, E: 0x94,
				H: 0x72, L: 0xD5,
				I: 0x80, R: 0x5E,
				PC: 0x5CB7, SP: 0x1FE3,
				IX: 0x70F0, IY: 0x72C7,
				AF_: 0xC8EE, BC_: 0x3BEB,
				DE_: 0x33B8, HL_: 0x4A90,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{23735, 237}, {23736, 162}, {29397, 0}},
				Ports: [][2]uint16{{16621, 117}},
			},
			want: z80State{
				A: 0x51, F: 0x3D,
				B: 0x3F, C: 0xED,
				D: 0x4F, E: 0x94,
				H: 0x72, L: 0xD6,
				I: 0x80, R: 0x60,
				PC: 0x5CB9, SP: 0x1FE3,
				IX: 0x70F0, IY: 0x72C7,
				AF_: 0xC8EE, BC_: 0x3BEB,
				DE_: 0x33B8, HL_: 0x4A90,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{23735, 237}, {23736, 162}, {29397, 117}},
				Cycles: 16,
			},
		},
		{ // ED AA 0000
			name: "ED AA 0000",
			init: z80State{
				A: 0xFA, F: 0x1B,
				B: 0xFA, C: 0x88,
				D: 0x80, E: 0xEA,
				H: 0xAA, L: 0xB0,
				I: 0xDD, R: 0x08,
				PC: 0x90BE, SP: 0xFF2F,
				IX: 0xD070, IY: 0xD3CB,
				AF_: 0x603E, BC_: 0x2120,
				DE_: 0x1689, HL_: 0x400C,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{37054, 237}, {37055, 170}, {43696, 0}},
				Ports: [][2]uint16{{64136, 41}},
			},
			want: z80State{
				A: 0xFA, F: 0xAC,
				B: 0xF9, C: 0x88,
				D: 0x80, E: 0xEA,
				H: 0xAA, L: 0xAF,
				I: 0xDD, R: 0x0A,
				PC: 0x90C0, SP: 0xFF2F,
				IX: 0xD070, IY: 0xD3CB,
				AF_: 0x603E, BC_: 0x2120,
				DE_: 0x1689, HL_: 0x400C,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{37054, 237}, {37055, 170}, {43696, 41}},
				Cycles: 16,
			},
		},
		{ // ED B2 0027 (INIR repeat, formula-compatible)
			name: "ED B2 0027",
			init: z80State{
				A: 0x75, F: 0x43,
				B: 0xEF, C: 0xAD,
				D: 0x7D, E: 0x3F,
				H: 0x1D, L: 0xE1,
				I: 0x8A, R: 0x7C,
				PC: 0x7BFD, SP: 0x46B0,
				IX: 0x5F55, IY: 0x0E0C,
				AF_: 0xBCB2, BC_: 0x8B21,
				DE_: 0x7D1E, HL_: 0xD9D8,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{7649, 0}, {31741, 237}, {31742, 178}},
				Ports: [][2]uint16{{61357, 4}},
			},
			want: z80State{
				A: 0x75, F: 0xA8,
				B: 0xEE, C: 0xAD,
				D: 0x7D, E: 0x3F,
				H: 0x1D, L: 0xE2,
				I: 0x8A, R: 0x7E,
				PC: 0x7BFD, SP: 0x46B0,
				IX: 0x5F55, IY: 0x0E0C,
				AF_: 0xBCB2, BC_: 0x8B21,
				DE_: 0x7D1E, HL_: 0xD9D8,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{7649, 4}, {31741, 237}, {31742, 178}},
				Cycles: 21,
			},
		},
		{ // ED B2 004B (INIR non-repeat B->0)
			name: "ED B2 004B",
			init: z80State{
				A: 0xAB, F: 0xC4,
				B: 0x01, C: 0xD8,
				D: 0x29, E: 0xD4,
				H: 0x33, L: 0x36,
				I: 0x16, R: 0x65,
				PC: 0xB829, SP: 0xCF58,
				IX: 0xF1B6, IY: 0x7B1E,
				AF_: 0xAA16, BC_: 0x285F,
				DE_: 0x660D, HL_: 0x6B1E,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{13110, 0}, {47145, 237}, {47146, 178}},
				Ports: [][2]uint16{{472, 198}},
			},
			want: z80State{
				A: 0xAB, F: 0x53,
				B: 0x00, C: 0xD8,
				D: 0x29, E: 0xD4,
				H: 0x33, L: 0x37,
				I: 0x16, R: 0x67,
				PC: 0xB82B, SP: 0xCF58,
				IX: 0xF1B6, IY: 0x7B1E,
				AF_: 0xAA16, BC_: 0x285F,
				DE_: 0x660D, HL_: 0x6B1E,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{13110, 198}, {47145, 237}, {47146, 178}},
				Cycles: 16,
			},
		},
		{ // ED BA 000D (INDR repeat, formula-compatible)
			name: "ED BA 000D",
			init: z80State{
				A: 0xE1, F: 0xD6,
				B: 0x67, C: 0x89,
				D: 0xDB, E: 0x44,
				H: 0x9A, L: 0x59,
				I: 0xAE, R: 0x5F,
				PC: 0xE4E0, SP: 0x1BBD,
				IX: 0xCC6E, IY: 0x0F44,
				AF_: 0x94E1, BC_: 0x7CB1,
				DE_: 0x1796, HL_: 0xB279,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{39513, 0}, {58592, 237}, {58593, 186}},
				Ports: [][2]uint16{{26505, 107}},
			},
			want: z80State{
				A: 0xE1, F: 0x24,
				B: 0x66, C: 0x89,
				D: 0xDB, E: 0x44,
				H: 0x9A, L: 0x58,
				I: 0xAE, R: 0x61,
				PC: 0xE4E0, SP: 0x1BBD,
				IX: 0xCC6E, IY: 0x0F44,
				AF_: 0x94E1, BC_: 0x7CB1,
				DE_: 0x1796, HL_: 0xB279,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{39513, 107}, {58592, 237}, {58593, 186}},
				Cycles: 21,
			},
		},
		{ // ED A3 0000
			name: "ED A3 0000",
			init: z80State{
				A: 0xBF, F: 0x6F,
				B: 0xD8, C: 0x0A,
				D: 0x88, E: 0x41,
				H: 0x05, L: 0x90,
				I: 0x3C, R: 0x34,
				PC: 0xCA85, SP: 0x5076,
				IX: 0xBC34, IY: 0xA7D5,
				AF_: 0x758C, BC_: 0x5A47,
				DE_: 0xA2AF, HL_: 0x89C3,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{1424, 242}, {51845, 237}, {51846, 163}},
			},
			want: z80State{
				A: 0xBF, F: 0x97,
				B: 0xD7, C: 0x0A,
				D: 0x88, E: 0x41,
				H: 0x05, L: 0x91,
				I: 0x3C, R: 0x36,
				PC: 0xCA87, SP: 0x5076,
				IX: 0xBC34, IY: 0xA7D5,
				AF_: 0x758C, BC_: 0x5A47,
				DE_: 0xA2AF, HL_: 0x89C3,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{1424, 242}, {51845, 237}, {51846, 163}},
				Cycles: 16,
			},
		},
		{ // ED A3 0001
			name: "ED A3 0001",
			init: z80State{
				A: 0xB4, F: 0x86,
				B: 0x38, C: 0x8F,
				D: 0xF2, E: 0xEA,
				H: 0xD3, L: 0x40,
				I: 0x2F, R: 0x04,
				PC: 0x9301, SP: 0xF0AB,
				IX: 0x4611, IY: 0x415D,
				AF_: 0xD6BE, BC_: 0x06BA,
				DE_: 0xF6B6, HL_: 0x41CE,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{37633, 237}, {37634, 163}, {54080, 156}},
			},
			want: z80State{
				A: 0xB4, F: 0x22,
				B: 0x37, C: 0x8F,
				D: 0xF2, E: 0xEA,
				H: 0xD3, L: 0x41,
				I: 0x2F, R: 0x06,
				PC: 0x9303, SP: 0xF0AB,
				IX: 0x4611, IY: 0x415D,
				AF_: 0xD6BE, BC_: 0x06BA,
				DE_: 0xF6B6, HL_: 0x41CE,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{37633, 237}, {37634, 163}, {54080, 156}},
				Cycles: 16,
			},
		},
		{ // ED AB 0000
			name: "ED AB 0000",
			init: z80State{
				A: 0x4A, F: 0x79,
				B: 0xD9, C: 0x93,
				D: 0x72, E: 0x39,
				H: 0xDE, L: 0x11,
				I: 0x56, R: 0x62,
				PC: 0xBA1C, SP: 0xF607,
				IX: 0xE1D3, IY: 0xF267,
				AF_: 0xEDCE, BC_: 0xA478,
				DE_: 0xEFD6, HL_: 0x75E3,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{47644, 237}, {47645, 171}, {56849, 208}},
			},
			want: z80State{
				A: 0x4A, F: 0x8E,
				B: 0xD8, C: 0x93,
				D: 0x72, E: 0x39,
				H: 0xDE, L: 0x10,
				I: 0x56, R: 0x64,
				PC: 0xBA1E, SP: 0xF607,
				IX: 0xE1D3, IY: 0xF267,
				AF_: 0xEDCE, BC_: 0xA478,
				DE_: 0xEFD6, HL_: 0x75E3,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{47644, 237}, {47645, 171}, {56849, 208}},
				Cycles: 16,
			},
		},
		{ // ED B3 0000
			name: "ED B3 0000",
			init: z80State{
				A: 0xCB, F: 0x7D,
				B: 0x49, C: 0xEA,
				D: 0x4F, E: 0x4D,
				H: 0xF6, L: 0x98,
				I: 0xCF, R: 0x5B,
				PC: 0x5E36, SP: 0xD3CF,
				IX: 0xD1BE, IY: 0x4AB3,
				AF_: 0x3849, BC_: 0x953E,
				DE_: 0xE998, HL_: 0x1FB3,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{24118, 237}, {24119, 179}, {63128, 29}},
			},
			want: z80State{
				A: 0xCB, F: 0x0C,
				B: 0x48, C: 0xEA,
				D: 0x4F, E: 0x4D,
				H: 0xF6, L: 0x99,
				I: 0xCF, R: 0x5D,
				PC: 0x5E36, SP: 0xD3CF,
				IX: 0xD1BE, IY: 0x4AB3,
				AF_: 0x3849, BC_: 0x953E,
				DE_: 0xE998, HL_: 0x1FB3,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{24118, 237}, {24119, 179}, {63128, 29}},
				Cycles: 21,
			},
		},
		{ // ED B3 0016 (OTIR repeat, formula-compatible)
			name: "ED B3 0016",
			init: z80State{
				A: 0x17, F: 0xFB,
				B: 0x49, C: 0xC3,
				D: 0xDB, E: 0xC8,
				H: 0x56, L: 0x49,
				I: 0xD1, R: 0x41,
				PC: 0x5D57, SP: 0x926A,
				IX: 0x4547, IY: 0xB66B,
				AF_: 0xEB19, BC_: 0x989B,
				DE_: 0x5D31, HL_: 0x1864,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{22089, 112}, {23895, 237}, {23896, 179}},
			},
			want: z80State{
				A: 0x17, F: 0x08,
				B: 0x48, C: 0xC3,
				D: 0xDB, E: 0xC8,
				H: 0x56, L: 0x4A,
				I: 0xD1, R: 0x43,
				PC: 0x5D57, SP: 0x926A,
				IX: 0x4547, IY: 0xB66B,
				AF_: 0xEB19, BC_: 0x989B,
				DE_: 0x5D31, HL_: 0x1864,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{22089, 112}, {23895, 237}, {23896, 179}},
				Cycles: 21,
			},
		},
		{ // ED BB 0012 (OTDR repeat, formula-compatible)
			name: "ED BB 0012",
			init: z80State{
				A: 0x19, F: 0x86,
				B: 0x3C, C: 0xD1,
				D: 0x88, E: 0x93,
				H: 0xA3, L: 0x4E,
				I: 0xC3, R: 0x01,
				PC: 0x3F3F, SP: 0x82EB,
				IX: 0xCF1F, IY: 0x4E22,
				AF_: 0x5280, BC_: 0x8FB1,
				DE_: 0xF7C8, HL_: 0x13EA,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{16191, 237}, {16192, 187}, {41806, 53}},
			},
			want: z80State{
				A: 0x19, F: 0x2C,
				B: 0x3B, C: 0xD1,
				D: 0x88, E: 0x93,
				H: 0xA3, L: 0x4D,
				I: 0xC3, R: 0x03,
				PC: 0x3F3F, SP: 0x82EB,
				IX: 0xCF1F, IY: 0x4E22,
				AF_: 0x5280, BC_: 0x8FB1,
				DE_: 0xF7C8, HL_: 0x13EA,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{16191, 237}, {16192, 187}, {41806, 53}},
				Cycles: 21,
			},
		},
		{ // ED BB 0037 (OTDR repeat, formula-compatible)
			name: "ED BB 0037",
			init: z80State{
				A: 0xF0, F: 0xE7,
				B: 0x1C, C: 0x32,
				D: 0xF5, E: 0x9B,
				H: 0xE3, L: 0x06,
				I: 0x43, R: 0x3D,
				PC: 0x88A5, SP: 0xC6EE,
				IX: 0x181B, IY: 0x4696,
				AF_: 0x2BCE, BC_: 0x863C,
				DE_: 0x372A, HL_: 0xAF8B,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{34981, 237}, {34982, 187}, {58118, 100}},
			},
			want: z80State{
				A: 0xF0, F: 0x08,
				B: 0x1B, C: 0x32,
				D: 0xF5, E: 0x9B,
				H: 0xE3, L: 0x05,
				I: 0x43, R: 0x3F,
				PC: 0x88A5, SP: 0xC6EE,
				IX: 0x181B, IY: 0x4696,
				AF_: 0x2BCE, BC_: 0x863C,
				DE_: 0x372A, HL_: 0xAF8B,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{34981, 237}, {34982, 187}, {58118, 100}},
				Cycles: 21,
			},
		},
	}
	for _, tt := range tests {
		runSSTTest(t, tt)
	}
}
