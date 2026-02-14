package z80

import "testing"

func TestRLCA(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x07
	cpu.reg.AF = 0x8500 // A=0x85 (10000101)
	cpu.Step()
	// Result: 00001011 = 0x0B, carry set
	if cpu.getA() != 0x0B {
		t.Errorf("RLCA: A=%02x want 0B", cpu.getA())
	}
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set")
	}
}

func TestRRCA(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x0F
	cpu.reg.AF = 0x8500 // A=0x85 (10000101)
	cpu.Step()
	// bit0=1, result: 1_1000010 = 0xC2
	if cpu.getA() != 0xC2 {
		t.Errorf("RRCA: A=%02x want C2", cpu.getA())
	}
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set")
	}
}

func TestRLA(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x17
	cpu.reg.AF = 0x8501 // A=0x85, carry=1
	cpu.Step()
	// shift left, old carry in: 00001011 = 0x0B, new carry set
	if cpu.getA() != 0x0B {
		t.Errorf("RLA: A=%02x want 0B", cpu.getA())
	}
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set")
	}
}

func TestRRA(t *testing.T) {
	cpu, bus := newTestCPU()
	bus.mem[0] = 0x1F
	cpu.reg.AF = 0x8501 // A=0x85, carry=1
	cpu.Step()
	// shift right, old carry in bit7: 11000010 = 0xC2, bit0 was 1 -> carry set
	if cpu.getA() != 0xC2 {
		t.Errorf("RRA: A=%02x want C2", cpu.getA())
	}
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set")
	}
}

func TestCB_RLC_r(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 00 = RLC B
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x00
	cpu.reg.BC = 0x8500
	cycles := cpu.Step()
	if cpu.getB() != 0x0B {
		t.Errorf("RLC B: B=%02x want 0B", cpu.getB())
	}
	if cycles != 8 {
		t.Errorf("cycles=%d want 8", cycles)
	}
	f := cpu.getF()
	if f&flagC == 0 {
		t.Error("C should be set")
	}
	if f&flagZ != 0 {
		t.Error("Z should be clear")
	}
}

func TestCB_RLC_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 06 = RLC (HL)
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x06
	cpu.reg.HL = 0x5000
	bus.mem[0x5000] = 0x80
	cycles := cpu.Step()
	if bus.mem[0x5000] != 0x01 {
		t.Errorf("RLC (HL): got %02x want 01", bus.mem[0x5000])
	}
	if cycles != 15 {
		t.Errorf("cycles=%d want 15", cycles)
	}
}

func TestCB_SLA(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 20 = SLA B
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x20
	cpu.reg.BC = 0x8000
	cpu.Step()
	if cpu.getB() != 0x00 {
		t.Errorf("SLA B: B=%02x want 00", cpu.getB())
	}
	f := cpu.getF()
	if f&flagC == 0 {
		t.Error("C should be set")
	}
	if f&flagZ == 0 {
		t.Error("Z should be set")
	}
}

func TestCB_SRA(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 28 = SRA B
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x28
	cpu.reg.BC = 0x8100
	cpu.Step()
	// 0x81 >> 1 with sign preserved = 0xC0, carry set (bit0=1)
	if cpu.getB() != 0xC0 {
		t.Errorf("SRA B: B=%02x want C0", cpu.getB())
	}
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set")
	}
	if cpu.getF()&flagS == 0 {
		t.Error("S should be set")
	}
}

func TestCB_SRL(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 38 = SRL B
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x38
	cpu.reg.BC = 0x8100
	cpu.Step()
	// 0x81 >> 1 = 0x40, carry set
	if cpu.getB() != 0x40 {
		t.Errorf("SRL B: B=%02x want 40", cpu.getB())
	}
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set")
	}
	if cpu.getF()&flagS != 0 {
		t.Error("S should be clear")
	}
}

func TestCB_SLL(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 30 = SLL B (undocumented)
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x30
	cpu.reg.BC = 0x8000
	cpu.Step()
	// 0x80 << 1 | 1 = 0x01, carry set
	if cpu.getB() != 0x01 {
		t.Errorf("SLL B: B=%02x want 01", cpu.getB())
	}
	if cpu.getF()&flagC == 0 {
		t.Error("C should be set")
	}
}

func TestSST_Rotate(t *testing.T) {
	tests := []struct {
		name       string
		init, want z80State
	}{
		{ // 07 0000
			name: "07 0000",
			init: z80State{
				A: 0xA5, F: 0x08,
				B: 0x14, C: 0x4A,
				D: 0x06, E: 0xB3,
				H: 0xEE, L: 0x5D,
				I: 0x46, R: 0x67,
				PC: 0x24F3, SP: 0xD509,
				IX: 0x8B6B, IY: 0xD362,
				AF_: 0xF13A, BC_: 0xAD0E,
				DE_: 0x8EA9, HL_: 0xE247,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{9459, 7}},
			},
			want: z80State{
				A: 0x4B, F: 0x09,
				B: 0x14, C: 0x4A,
				D: 0x06, E: 0xB3,
				H: 0xEE, L: 0x5D,
				I: 0x46, R: 0x68,
				PC: 0x24F4, SP: 0xD509,
				IX: 0x8B6B, IY: 0xD362,
				AF_: 0xF13A, BC_: 0xAD0E,
				DE_: 0x8EA9, HL_: 0xE247,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{9459, 7}},
				Cycles: 4,
			},
		},
		{ // 0F 0000
			name: "0F 0000",
			init: z80State{
				A: 0x34, F: 0x25,
				B: 0xDC, C: 0xE1,
				D: 0xC2, E: 0xFC,
				H: 0x58, L: 0x9D,
				I: 0x9C, R: 0x19,
				PC: 0xD304, SP: 0xA4AA,
				IX: 0xE430, IY: 0x3411,
				AF_: 0x4487, BC_: 0xC39E,
				DE_: 0xDDF3, HL_: 0x06A3,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{54020, 15}},
			},
			want: z80State{
				A: 0x1A, F: 0x0C,
				B: 0xDC, C: 0xE1,
				D: 0xC2, E: 0xFC,
				H: 0x58, L: 0x9D,
				I: 0x9C, R: 0x1A,
				PC: 0xD305, SP: 0xA4AA,
				IX: 0xE430, IY: 0x3411,
				AF_: 0x4487, BC_: 0xC39E,
				DE_: 0xDDF3, HL_: 0x06A3,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{54020, 15}},
				Cycles: 4,
			},
		},
		{ // 17 0000
			name: "17 0000",
			init: z80State{
				A: 0x8A, F: 0x6E,
				B: 0xAD, C: 0x63,
				D: 0xC0, E: 0x1B,
				H: 0x9E, L: 0xA6,
				I: 0xD0, R: 0x15,
				PC: 0xA8E0, SP: 0x12CE,
				IX: 0x3B4D, IY: 0x8377,
				AF_: 0xBA4F, BC_: 0xBA04,
				DE_: 0xE6D2, HL_: 0x00A9,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{43232, 23}},
			},
			want: z80State{
				A: 0x14, F: 0x45,
				B: 0xAD, C: 0x63,
				D: 0xC0, E: 0x1B,
				H: 0x9E, L: 0xA6,
				I: 0xD0, R: 0x16,
				PC: 0xA8E1, SP: 0x12CE,
				IX: 0x3B4D, IY: 0x8377,
				AF_: 0xBA4F, BC_: 0xBA04,
				DE_: 0xE6D2, HL_: 0x00A9,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{43232, 23}},
				Cycles: 4,
			},
		},
		{ // 1F 0000
			name: "1F 0000",
			init: z80State{
				A: 0xCE, F: 0x7A,
				B: 0x16, C: 0x76,
				D: 0xC6, E: 0xB5,
				H: 0x4E, L: 0xBA,
				I: 0x70, R: 0x0C,
				PC: 0x06CF, SP: 0x92F9,
				IX: 0xA18A, IY: 0xC75D,
				AF_: 0xA585, BC_: 0x70ED,
				DE_: 0x9069, HL_: 0xE936,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1743, 31}},
			},
			want: z80State{
				A: 0x67, F: 0x60,
				B: 0x16, C: 0x76,
				D: 0xC6, E: 0xB5,
				H: 0x4E, L: 0xBA,
				I: 0x70, R: 0x0D,
				PC: 0x06D0, SP: 0x92F9,
				IX: 0xA18A, IY: 0xC75D,
				AF_: 0xA585, BC_: 0x70ED,
				DE_: 0x9069, HL_: 0xE936,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{1743, 31}},
				Cycles: 4,
			},
		},
		{ // CB 00 0000
			name: "CB 00 0000",
			init: z80State{
				A: 0x21, F: 0xEA,
				B: 0xBD, C: 0x27,
				D: 0x37, E: 0x6D,
				H: 0x50, L: 0xA3,
				I: 0xA0, R: 0x45,
				PC: 0xC817, SP: 0xE460,
				IX: 0x89F7, IY: 0xCE76,
				AF_: 0xDBFF, BC_: 0x8F1E,
				DE_: 0xEAA9, HL_: 0xF321,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{51223, 203}, {51224, 0}},
			},
			want: z80State{
				A: 0x21, F: 0x2D,
				B: 0x7B, C: 0x27,
				D: 0x37, E: 0x6D,
				H: 0x50, L: 0xA3,
				I: 0xA0, R: 0x47,
				PC: 0xC819, SP: 0xE460,
				IX: 0x89F7, IY: 0xCE76,
				AF_: 0xDBFF, BC_: 0x8F1E,
				DE_: 0xEAA9, HL_: 0xF321,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{51223, 203}, {51224, 0}},
				Cycles: 8,
			},
		},
		{ // CB 06 0000
			name: "CB 06 0000",
			init: z80State{
				A: 0x20, F: 0x6C,
				B: 0xC5, C: 0xD6,
				D: 0x62, E: 0x6A,
				H: 0x8C, L: 0x9B,
				I: 0x2E, R: 0x10,
				PC: 0xD792, SP: 0x230F,
				IX: 0xD6F9, IY: 0xC574,
				AF_: 0xFA16, BC_: 0x55F6,
				DE_: 0xEA05, HL_: 0x3559,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{35995, 21}, {55186, 203}, {55187, 6}},
			},
			want: z80State{
				A: 0x20, F: 0x28,
				B: 0xC5, C: 0xD6,
				D: 0x62, E: 0x6A,
				H: 0x8C, L: 0x9B,
				I: 0x2E, R: 0x12,
				PC: 0xD794, SP: 0x230F,
				IX: 0xD6F9, IY: 0xC574,
				AF_: 0xFA16, BC_: 0x55F6,
				DE_: 0xEA05, HL_: 0x3559,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{35995, 42}, {55186, 203}, {55187, 6}},
				Cycles: 15,
			},
		},
		{ // CB 08 0000
			name: "CB 08 0000",
			init: z80State{
				A: 0xF2, F: 0xCE,
				B: 0x92, C: 0x94,
				D: 0x72, E: 0xF3,
				H: 0xF7, L: 0x15,
				I: 0x86, R: 0x4F,
				PC: 0xCA47, SP: 0x768F,
				IX: 0x3F57, IY: 0x1463,
				AF_: 0x8A17, BC_: 0x84C0,
				DE_: 0x861D, HL_: 0x1242,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{51783, 203}, {51784, 8}},
			},
			want: z80State{
				A: 0xF2, F: 0x08,
				B: 0x49, C: 0x94,
				D: 0x72, E: 0xF3,
				H: 0xF7, L: 0x15,
				I: 0x86, R: 0x51,
				PC: 0xCA49, SP: 0x768F,
				IX: 0x3F57, IY: 0x1463,
				AF_: 0x8A17, BC_: 0x84C0,
				DE_: 0x861D, HL_: 0x1242,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{51783, 203}, {51784, 8}},
				Cycles: 8,
			},
		},
		{ // CB 0E 0000
			name: "CB 0E 0000",
			init: z80State{
				A: 0xE8, F: 0x96,
				B: 0x09, C: 0x5A,
				D: 0x63, E: 0x96,
				H: 0x2F, L: 0xC5,
				I: 0x52, R: 0x54,
				PC: 0xEF1C, SP: 0x548F,
				IX: 0xEC98, IY: 0x7155,
				AF_: 0x0815, BC_: 0x94B2,
				DE_: 0xDD5B, HL_: 0x36DD,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{12229, 8}, {61212, 203}, {61213, 14}},
			},
			want: z80State{
				A: 0xE8, F: 0x00,
				B: 0x09, C: 0x5A,
				D: 0x63, E: 0x96,
				H: 0x2F, L: 0xC5,
				I: 0x52, R: 0x56,
				PC: 0xEF1E, SP: 0x548F,
				IX: 0xEC98, IY: 0x7155,
				AF_: 0x0815, BC_: 0x94B2,
				DE_: 0xDD5B, HL_: 0x36DD,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{12229, 4}, {61212, 203}, {61213, 14}},
				Cycles: 15,
			},
		},
		{ // CB 10 0000
			name: "CB 10 0000",
			init: z80State{
				A: 0x15, F: 0x82,
				B: 0x1F, C: 0x62,
				D: 0x64, E: 0x76,
				H: 0xB3, L: 0xE7,
				I: 0x95, R: 0x79,
				PC: 0x98F9, SP: 0x2B26,
				IX: 0xAC59, IY: 0xDA0C,
				AF_: 0xDA4E, BC_: 0x630E,
				DE_: 0x4594, HL_: 0x6A31,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{39161, 203}, {39162, 16}},
			},
			want: z80State{
				A: 0x15, F: 0x28,
				B: 0x3E, C: 0x62,
				D: 0x64, E: 0x76,
				H: 0xB3, L: 0xE7,
				I: 0x95, R: 0x7B,
				PC: 0x98FB, SP: 0x2B26,
				IX: 0xAC59, IY: 0xDA0C,
				AF_: 0xDA4E, BC_: 0x630E,
				DE_: 0x4594, HL_: 0x6A31,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{39161, 203}, {39162, 16}},
				Cycles: 8,
			},
		},
		{ // CB 16 0000
			name: "CB 16 0000",
			init: z80State{
				A: 0x8A, F: 0xAB,
				B: 0x14, C: 0x19,
				D: 0x76, E: 0xD4,
				H: 0x0F, L: 0xDC,
				I: 0xE7, R: 0x62,
				PC: 0xED42, SP: 0x7DAF,
				IX: 0xED93, IY: 0x24B8,
				AF_: 0x9F78, BC_: 0xFCBB,
				DE_: 0x8DF1, HL_: 0xFCAC,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{4060, 207}, {60738, 203}, {60739, 22}},
			},
			want: z80State{
				A: 0x8A, F: 0x8D,
				B: 0x14, C: 0x19,
				D: 0x76, E: 0xD4,
				H: 0x0F, L: 0xDC,
				I: 0xE7, R: 0x64,
				PC: 0xED44, SP: 0x7DAF,
				IX: 0xED93, IY: 0x24B8,
				AF_: 0x9F78, BC_: 0xFCBB,
				DE_: 0x8DF1, HL_: 0xFCAC,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{4060, 159}, {60738, 203}, {60739, 22}},
				Cycles: 15,
			},
		},
		{ // CB 18 0000
			name: "CB 18 0000",
			init: z80State{
				A: 0x94, F: 0x2C,
				B: 0x00, C: 0xBB,
				D: 0x4D, E: 0x02,
				H: 0xA8, L: 0xF0,
				I: 0xB4, R: 0x55,
				PC: 0xAD53, SP: 0xA383,
				IX: 0x8539, IY: 0x944B,
				AF_: 0x5931, BC_: 0x6B4D,
				DE_: 0x5CC3, HL_: 0x19FC,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{44371, 203}, {44372, 24}},
			},
			want: z80State{
				A: 0x94, F: 0x44,
				B: 0x00, C: 0xBB,
				D: 0x4D, E: 0x02,
				H: 0xA8, L: 0xF0,
				I: 0xB4, R: 0x57,
				PC: 0xAD55, SP: 0xA383,
				IX: 0x8539, IY: 0x944B,
				AF_: 0x5931, BC_: 0x6B4D,
				DE_: 0x5CC3, HL_: 0x19FC,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{44371, 203}, {44372, 24}},
				Cycles: 8,
			},
		},
		{ // CB 1E 0000
			name: "CB 1E 0000",
			init: z80State{
				A: 0x75, F: 0xFD,
				B: 0x80, C: 0xA1,
				D: 0xA8, E: 0x2E,
				H: 0xC8, L: 0x50,
				I: 0xA1, R: 0x48,
				PC: 0xF8B0, SP: 0x3742,
				IX: 0xCFEF, IY: 0x01DC,
				AF_: 0x0190, BC_: 0xD870,
				DE_: 0x9AD2, HL_: 0xB7DA,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{51280, 249}, {63664, 203}, {63665, 30}},
			},
			want: z80State{
				A: 0x75, F: 0xAD,
				B: 0x80, C: 0xA1,
				D: 0xA8, E: 0x2E,
				H: 0xC8, L: 0x50,
				I: 0xA1, R: 0x4A,
				PC: 0xF8B2, SP: 0x3742,
				IX: 0xCFEF, IY: 0x01DC,
				AF_: 0x0190, BC_: 0xD870,
				DE_: 0x9AD2, HL_: 0xB7DA,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{51280, 252}, {63664, 203}, {63665, 30}},
				Cycles: 15,
			},
		},
		{ // CB 20 0000
			name: "CB 20 0000",
			init: z80State{
				A: 0x61, F: 0x70,
				B: 0x78, C: 0x33,
				D: 0xE1, E: 0x4B,
				H: 0xFF, L: 0x00,
				I: 0xE9, R: 0x19,
				PC: 0x7B9B, SP: 0x998A,
				IX: 0x3064, IY: 0x12CB,
				AF_: 0xE053, BC_: 0x3462,
				DE_: 0x3DF9, HL_: 0x45C7,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{31643, 203}, {31644, 32}},
			},
			want: z80State{
				A: 0x61, F: 0xA4,
				B: 0xF0, C: 0x33,
				D: 0xE1, E: 0x4B,
				H: 0xFF, L: 0x00,
				I: 0xE9, R: 0x1B,
				PC: 0x7B9D, SP: 0x998A,
				IX: 0x3064, IY: 0x12CB,
				AF_: 0xE053, BC_: 0x3462,
				DE_: 0x3DF9, HL_: 0x45C7,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{31643, 203}, {31644, 32}},
				Cycles: 8,
			},
		},
		{ // CB 26 0000
			name: "CB 26 0000",
			init: z80State{
				A: 0xD4, F: 0xD2,
				B: 0xC1, C: 0xA6,
				D: 0xD0, E: 0x90,
				H: 0xDB, L: 0x3D,
				I: 0x66, R: 0x7F,
				PC: 0x50B7, SP: 0xFBA1,
				IX: 0x0905, IY: 0x30F0,
				AF_: 0xB0CE, BC_: 0xB555,
				DE_: 0x54F0, HL_: 0xC1B6,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{20663, 203}, {20664, 38}, {56125, 111}},
			},
			want: z80State{
				A: 0xD4, F: 0x8C,
				B: 0xC1, C: 0xA6,
				D: 0xD0, E: 0x90,
				H: 0xDB, L: 0x3D,
				I: 0x66, R: 0x01,
				PC: 0x50B9, SP: 0xFBA1,
				IX: 0x0905, IY: 0x30F0,
				AF_: 0xB0CE, BC_: 0xB555,
				DE_: 0x54F0, HL_: 0xC1B6,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{20663, 203}, {20664, 38}, {56125, 222}},
				Cycles: 15,
			},
		},
		{ // CB 28 0000
			name: "CB 28 0000",
			init: z80State{
				A: 0x78, F: 0x91,
				B: 0xB2, C: 0x1F,
				D: 0x86, E: 0xB3,
				H: 0x58, L: 0xA6,
				I: 0xC1, R: 0x51,
				PC: 0x7087, SP: 0xC453,
				IX: 0x5ED5, IY: 0x6793,
				AF_: 0xAF0D, BC_: 0xF0B5,
				DE_: 0x2CCB, HL_: 0x7DC8,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{28807, 203}, {28808, 40}},
			},
			want: z80State{
				A: 0x78, F: 0x88,
				B: 0xD9, C: 0x1F,
				D: 0x86, E: 0xB3,
				H: 0x58, L: 0xA6,
				I: 0xC1, R: 0x53,
				PC: 0x7089, SP: 0xC453,
				IX: 0x5ED5, IY: 0x6793,
				AF_: 0xAF0D, BC_: 0xF0B5,
				DE_: 0x2CCB, HL_: 0x7DC8,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{28807, 203}, {28808, 40}},
				Cycles: 8,
			},
		},
		{ // CB 30 0000
			name: "CB 30 0000",
			init: z80State{
				A: 0xDD, F: 0xFC,
				B: 0x1B, C: 0xF5,
				D: 0x35, E: 0x65,
				H: 0x41, L: 0x4D,
				I: 0xA3, R: 0x62,
				PC: 0x3066, SP: 0x6FB5,
				IX: 0xD6B3, IY: 0xC483,
				AF_: 0x15CC, BC_: 0xCC96,
				DE_: 0x267A, HL_: 0x80A6,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{12390, 203}, {12391, 48}},
			},
			want: z80State{
				A: 0xDD, F: 0x20,
				B: 0x37, C: 0xF5,
				D: 0x35, E: 0x65,
				H: 0x41, L: 0x4D,
				I: 0xA3, R: 0x64,
				PC: 0x3068, SP: 0x6FB5,
				IX: 0xD6B3, IY: 0xC483,
				AF_: 0x15CC, BC_: 0xCC96,
				DE_: 0x267A, HL_: 0x80A6,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{12390, 203}, {12391, 48}},
				Cycles: 8,
			},
		},
		{ // CB 38 0000
			name: "CB 38 0000",
			init: z80State{
				A: 0x7D, F: 0x23,
				B: 0x90, C: 0x33,
				D: 0x08, E: 0x1B,
				H: 0x02, L: 0x22,
				I: 0xBA, R: 0x64,
				PC: 0xDEA5, SP: 0xA836,
				IX: 0x934C, IY: 0xC5D6,
				AF_: 0x11B0, BC_: 0x16C4,
				DE_: 0x6938, HL_: 0x0D60,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{56997, 203}, {56998, 56}},
			},
			want: z80State{
				A: 0x7D, F: 0x0C,
				B: 0x48, C: 0x33,
				D: 0x08, E: 0x1B,
				H: 0x02, L: 0x22,
				I: 0xBA, R: 0x66,
				PC: 0xDEA7, SP: 0xA836,
				IX: 0x934C, IY: 0xC5D6,
				AF_: 0x11B0, BC_: 0x16C4,
				DE_: 0x6938, HL_: 0x0D60,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{56997, 203}, {56998, 56}},
				Cycles: 8,
			},
		},
		{ // CB 3E 0000
			name: "CB 3E 0000",
			init: z80State{
				A: 0x1A, F: 0xB3,
				B: 0xB6, C: 0x35,
				D: 0x68, E: 0x9F,
				H: 0x2F, L: 0xF4,
				I: 0x87, R: 0x30,
				PC: 0x4E2D, SP: 0xBE69,
				IX: 0xA36A, IY: 0x97AC,
				AF_: 0xD7D8, BC_: 0xCBDA,
				DE_: 0xA68E, HL_: 0x467B,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{12276, 37}, {20013, 203}, {20014, 62}},
			},
			want: z80State{
				A: 0x1A, F: 0x05,
				B: 0xB6, C: 0x35,
				D: 0x68, E: 0x9F,
				H: 0x2F, L: 0xF4,
				I: 0x87, R: 0x32,
				PC: 0x4E2F, SP: 0xBE69,
				IX: 0xA36A, IY: 0x97AC,
				AF_: 0xD7D8, BC_: 0xCBDA,
				DE_: 0xA68E, HL_: 0x467B,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{12276, 18}, {20013, 203}, {20014, 62}},
				Cycles: 15,
			},
		},
	}
	for _, tt := range tests {
		runSSTTest(t, tt)
	}
}
