package z80

import "testing"

func TestCB_BIT(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 47 = BIT 0, A
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x47
	cpu.reg.AF = 0xFE00 // bit 0 clear
	cycles := cpu.Step()
	if cycles != 8 {
		t.Errorf("cycles=%d want 8", cycles)
	}
	f := cpu.getF()
	if f&flagZ == 0 {
		t.Error("Z should be set when bit is 0")
	}
	if f&flagH == 0 {
		t.Error("H should be set")
	}
	if f&flagN != 0 {
		t.Error("N should be clear")
	}
}

func TestCB_BIT_Set(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 47 = BIT 0, A
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x47
	cpu.reg.AF = 0xFF00 // bit 0 set
	cpu.Step()
	if cpu.getF()&flagZ != 0 {
		t.Error("Z should be clear when bit is 1")
	}
}

func TestCB_BIT_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 7E = BIT 7, (HL)
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x7E
	cpu.reg.HL = 0x5000
	bus.mem[0x5000] = 0x80
	cycles := cpu.Step()
	if cycles != 12 {
		t.Errorf("cycles=%d want 12", cycles)
	}
	if cpu.getF()&flagZ != 0 {
		t.Error("Z should be clear (bit 7 is set)")
	}
	if cpu.getF()&flagS == 0 {
		t.Error("S should be set for BIT 7")
	}
}

func TestCB_RES(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 87 = RES 0, A
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x87
	cpu.reg.AF = 0xFF00
	cycles := cpu.Step()
	if cpu.getA() != 0xFE {
		t.Errorf("RES 0,A: A=%02x want FE", cpu.getA())
	}
	if cycles != 8 {
		t.Errorf("cycles=%d want 8", cycles)
	}
}

func TestCB_RES_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB 86 = RES 0, (HL)
	bus.mem[0] = 0xCB
	bus.mem[1] = 0x86
	cpu.reg.HL = 0x6000
	bus.mem[0x6000] = 0xFF
	cycles := cpu.Step()
	if bus.mem[0x6000] != 0xFE {
		t.Errorf("RES 0,(HL): got %02x want FE", bus.mem[0x6000])
	}
	if cycles != 15 {
		t.Errorf("cycles=%d want 15", cycles)
	}
}

func TestCB_SET(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB C7 = SET 0, A
	bus.mem[0] = 0xCB
	bus.mem[1] = 0xC7
	cpu.reg.AF = 0x0000
	cycles := cpu.Step()
	if cpu.getA() != 0x01 {
		t.Errorf("SET 0,A: A=%02x want 01", cpu.getA())
	}
	if cycles != 8 {
		t.Errorf("cycles=%d want 8", cycles)
	}
}

func TestCB_SET_HL(t *testing.T) {
	cpu, bus := newTestCPU()
	// CB FE = SET 7, (HL)
	bus.mem[0] = 0xCB
	bus.mem[1] = 0xFE
	cpu.reg.HL = 0x7000
	bus.mem[0x7000] = 0x00
	cycles := cpu.Step()
	if bus.mem[0x7000] != 0x80 {
		t.Errorf("SET 7,(HL): got %02x want 80", bus.mem[0x7000])
	}
	if cycles != 15 {
		t.Errorf("cycles=%d want 15", cycles)
	}
}

func TestSST_Bit(t *testing.T) {
	tests := []struct {
		name       string
		init, want z80State
	}{
		{ // CB 40 0000
			name: "CB 40 0000",
			init: z80State{
				A: 0x5E, F: 0x82,
				B: 0x87, C: 0x57,
				D: 0x98, E: 0x17,
				H: 0x70, L: 0xBB,
				I: 0x25, R: 0x4C,
				PC: 0xEEBD, SP: 0x09EA,
				IX: 0x815E, IY: 0x3BD2,
				AF_: 0x6582, BC_: 0x27EF,
				DE_: 0xA172, HL_: 0x92D5,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{61117, 203}, {61118, 64}},
			},
			want: z80State{
				A: 0x5E, F: 0x10,
				B: 0x87, C: 0x57,
				D: 0x98, E: 0x17,
				H: 0x70, L: 0xBB,
				I: 0x25, R: 0x4E,
				PC: 0xEEBF, SP: 0x09EA,
				IX: 0x815E, IY: 0x3BD2,
				AF_: 0x6582, BC_: 0x27EF,
				DE_: 0xA172, HL_: 0x92D5,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{61117, 203}, {61118, 64}},
				Cycles: 8,
			},
		},
		{ // CB 46 0002 (WZ F3/F5 aligned with H)
			name: "CB 46 0002",
			init: z80State{
				A: 0x47, F: 0xE0,
				B: 0xFB, C: 0x38,
				D: 0xAC, E: 0x0B,
				H: 0xE7, L: 0x03,
				I: 0x83, R: 0x51,
				PC: 0x9C60, SP: 0x2295,
				IX: 0x7B83, IY: 0xA880,
				AF_: 0xC486, BC_: 0x8CD4,
				DE_: 0xBFA5, HL_: 0x3F9B,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{40032, 203}, {40033, 70}, {59139, 16}},
			},
			want: z80State{
				A: 0x47, F: 0x74,
				B: 0xFB, C: 0x38,
				D: 0xAC, E: 0x0B,
				H: 0xE7, L: 0x03,
				I: 0x83, R: 0x53,
				PC: 0x9C62, SP: 0x2295,
				IX: 0x7B83, IY: 0xA880,
				AF_: 0xC486, BC_: 0x8CD4,
				DE_: 0xBFA5, HL_: 0x3F9B,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{40032, 203}, {40033, 70}, {59139, 16}},
				Cycles: 12,
			},
		},
		{ // CB 7F 0000
			name: "CB 7F 0000",
			init: z80State{
				A: 0x73, F: 0x02,
				B: 0x11, C: 0x05,
				D: 0x18, E: 0x6D,
				H: 0x5B, L: 0xF1,
				I: 0x14, R: 0x35,
				PC: 0x159D, SP: 0x8AB7,
				IX: 0xEF77, IY: 0x2D40,
				AF_: 0xABD0, BC_: 0x2DD1,
				DE_: 0x7466, HL_: 0x6FD4,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{5533, 203}, {5534, 127}},
			},
			want: z80State{
				A: 0x73, F: 0x74,
				B: 0x11, C: 0x05,
				D: 0x18, E: 0x6D,
				H: 0x5B, L: 0xF1,
				I: 0x14, R: 0x37,
				PC: 0x159F, SP: 0x8AB7,
				IX: 0xEF77, IY: 0x2D40,
				AF_: 0xABD0, BC_: 0x2DD1,
				DE_: 0x7466, HL_: 0x6FD4,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{5533, 203}, {5534, 127}},
				Cycles: 8,
			},
		},
		{ // CB 7E 0003 (WZ F3/F5 aligned with H)
			name: "CB 7E 0003",
			init: z80State{
				A: 0x40, F: 0x39,
				B: 0xB1, C: 0x8F,
				D: 0x74, E: 0xB0,
				H: 0x8A, L: 0x3F,
				I: 0x20, R: 0x48,
				PC: 0x3B9B, SP: 0x6F3B,
				IX: 0xA25D, IY: 0x9E08,
				AF_: 0xED6D, BC_: 0xF48A,
				DE_: 0x6778, HL_: 0x4E2C,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{15259, 203}, {15260, 126}, {35391, 86}},
			},
			want: z80State{
				A: 0x40, F: 0x5D,
				B: 0xB1, C: 0x8F,
				D: 0x74, E: 0xB0,
				H: 0x8A, L: 0x3F,
				I: 0x20, R: 0x4A,
				PC: 0x3B9D, SP: 0x6F3B,
				IX: 0xA25D, IY: 0x9E08,
				AF_: 0xED6D, BC_: 0xF48A,
				DE_: 0x6778, HL_: 0x4E2C,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{15259, 203}, {15260, 126}, {35391, 86}},
				Cycles: 12,
			},
		},
		{ // CB 50 0000
			name: "CB 50 0000",
			init: z80State{
				A: 0xA6, F: 0x2D,
				B: 0x8D, C: 0x9E,
				D: 0x97, E: 0x85,
				H: 0xA3, L: 0xDD,
				I: 0x73, R: 0x71,
				PC: 0xA5FA, SP: 0x2C51,
				IX: 0x6B43, IY: 0xBC57,
				AF_: 0x5191, BC_: 0x7BF3,
				DE_: 0x73D0, HL_: 0xE146,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{42490, 203}, {42491, 80}},
			},
			want: z80State{
				A: 0xA6, F: 0x19,
				B: 0x8D, C: 0x9E,
				D: 0x97, E: 0x85,
				H: 0xA3, L: 0xDD,
				I: 0x73, R: 0x73,
				PC: 0xA5FC, SP: 0x2C51,
				IX: 0x6B43, IY: 0xBC57,
				AF_: 0x5191, BC_: 0x7BF3,
				DE_: 0x73D0, HL_: 0xE146,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{42490, 203}, {42491, 80}},
				Cycles: 8,
			},
		},
		{ // CB 68 0000
			name: "CB 68 0000",
			init: z80State{
				A: 0xE5, F: 0x56,
				B: 0x64, C: 0x1C,
				D: 0xBC, E: 0xA2,
				H: 0xA3, L: 0x40,
				I: 0xB8, R: 0x08,
				PC: 0xE180, SP: 0x9901,
				IX: 0xCE64, IY: 0xF36B,
				AF_: 0xB286, BC_: 0x1B4B,
				DE_: 0x5937, HL_: 0x0976,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{57728, 203}, {57729, 104}},
			},
			want: z80State{
				A: 0xE5, F: 0x30,
				B: 0x64, C: 0x1C,
				D: 0xBC, E: 0xA2,
				H: 0xA3, L: 0x40,
				I: 0xB8, R: 0x0A,
				PC: 0xE182, SP: 0x9901,
				IX: 0xCE64, IY: 0xF36B,
				AF_: 0xB286, BC_: 0x1B4B,
				DE_: 0x5937, HL_: 0x0976,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{57728, 203}, {57729, 104}},
				Cycles: 8,
			},
		},
		{ // CB 80 0000
			name: "CB 80 0000",
			init: z80State{
				A: 0xE4, F: 0xBA,
				B: 0x38, C: 0x6F,
				D: 0x8A, E: 0xED,
				H: 0xEB, L: 0x8B,
				I: 0xD3, R: 0x14,
				PC: 0xAC63, SP: 0x4055,
				IX: 0x0B89, IY: 0xC5DD,
				AF_: 0x2447, BC_: 0xD00A,
				DE_: 0x05CF, HL_: 0x3F80,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{44131, 203}, {44132, 128}},
			},
			want: z80State{
				A: 0xE4, F: 0xBA,
				B: 0x38, C: 0x6F,
				D: 0x8A, E: 0xED,
				H: 0xEB, L: 0x8B,
				I: 0xD3, R: 0x16,
				PC: 0xAC65, SP: 0x4055,
				IX: 0x0B89, IY: 0xC5DD,
				AF_: 0x2447, BC_: 0xD00A,
				DE_: 0x05CF, HL_: 0x3F80,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{44131, 203}, {44132, 128}},
				Cycles: 8,
			},
		},
		{ // CB 86 0000
			name: "CB 86 0000",
			init: z80State{
				A: 0x19, F: 0x5F,
				B: 0x36, C: 0x12,
				D: 0x0D, E: 0x28,
				H: 0xF3, L: 0x6E,
				I: 0x45, R: 0x62,
				PC: 0x815C, SP: 0xB704,
				IX: 0xCF8D, IY: 0x7733,
				AF_: 0x875C, BC_: 0xB2CE,
				DE_: 0xB24E, HL_: 0xA5DE,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{33116, 203}, {33117, 134}, {62318, 41}},
			},
			want: z80State{
				A: 0x19, F: 0x5F,
				B: 0x36, C: 0x12,
				D: 0x0D, E: 0x28,
				H: 0xF3, L: 0x6E,
				I: 0x45, R: 0x64,
				PC: 0x815E, SP: 0xB704,
				IX: 0xCF8D, IY: 0x7733,
				AF_: 0x875C, BC_: 0xB2CE,
				DE_: 0xB24E, HL_: 0xA5DE,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{33116, 203}, {33117, 134}, {62318, 40}},
				Cycles: 15,
			},
		},
		{ // CB B8 0000
			name: "CB B8 0000",
			init: z80State{
				A: 0xE0, F: 0x7B,
				B: 0x37, C: 0x0C,
				D: 0x35, E: 0x5A,
				H: 0x20, L: 0x69,
				I: 0x0E, R: 0x6E,
				PC: 0xCA41, SP: 0xE675,
				IX: 0xB55C, IY: 0xA9E8,
				AF_: 0xD8AA, BC_: 0x2FD0,
				DE_: 0xA0E4, HL_: 0xCA14,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{51777, 203}, {51778, 184}},
			},
			want: z80State{
				A: 0xE0, F: 0x7B,
				B: 0x37, C: 0x0C,
				D: 0x35, E: 0x5A,
				H: 0x20, L: 0x69,
				I: 0x0E, R: 0x70,
				PC: 0xCA43, SP: 0xE675,
				IX: 0xB55C, IY: 0xA9E8,
				AF_: 0xD8AA, BC_: 0x2FD0,
				DE_: 0xA0E4, HL_: 0xCA14,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{51777, 203}, {51778, 184}},
				Cycles: 8,
			},
		},
		{ // CB BE 0000
			name: "CB BE 0000",
			init: z80State{
				A: 0x75, F: 0xA0,
				B: 0x29, C: 0x6E,
				D: 0x59, E: 0x34,
				H: 0x35, L: 0x9D,
				I: 0x06, R: 0x57,
				PC: 0x8E1B, SP: 0xB6D1,
				IX: 0xC3EB, IY: 0xB745,
				AF_: 0xBF34, BC_: 0x4B4D,
				DE_: 0xCC51, HL_: 0x51D0,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{13725, 239}, {36379, 203}, {36380, 190}},
			},
			want: z80State{
				A: 0x75, F: 0xA0,
				B: 0x29, C: 0x6E,
				D: 0x59, E: 0x34,
				H: 0x35, L: 0x9D,
				I: 0x06, R: 0x59,
				PC: 0x8E1D, SP: 0xB6D1,
				IX: 0xC3EB, IY: 0xB745,
				AF_: 0xBF34, BC_: 0x4B4D,
				DE_: 0xCC51, HL_: 0x51D0,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{13725, 111}, {36379, 203}, {36380, 190}},
				Cycles: 15,
			},
		},
		{ // CB C0 0000
			name: "CB C0 0000",
			init: z80State{
				A: 0x11, F: 0x2F,
				B: 0x32, C: 0x83,
				D: 0x7D, E: 0x1A,
				H: 0x33, L: 0x8C,
				I: 0x6F, R: 0x19,
				PC: 0x950D, SP: 0xE16B,
				IX: 0xA251, IY: 0xE715,
				AF_: 0x7DF6, BC_: 0xAFB2,
				DE_: 0x05F8, HL_: 0xEE51,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{38157, 203}, {38158, 192}},
			},
			want: z80State{
				A: 0x11, F: 0x2F,
				B: 0x33, C: 0x83,
				D: 0x7D, E: 0x1A,
				H: 0x33, L: 0x8C,
				I: 0x6F, R: 0x1B,
				PC: 0x950F, SP: 0xE16B,
				IX: 0xA251, IY: 0xE715,
				AF_: 0x7DF6, BC_: 0xAFB2,
				DE_: 0x05F8, HL_: 0xEE51,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{38157, 203}, {38158, 192}},
				Cycles: 8,
			},
		},
		{ // CB C6 0000
			name: "CB C6 0000",
			init: z80State{
				A: 0xAF, F: 0xD4,
				B: 0xDB, C: 0xB6,
				D: 0xCE, E: 0x6B,
				H: 0xA2, L: 0xB8,
				I: 0x1D, R: 0x4F,
				PC: 0x454E, SP: 0x0B7D,
				IX: 0x95DF, IY: 0x5D6F,
				AF_: 0xFF23, BC_: 0x1EE8,
				DE_: 0x430B, HL_: 0xC871,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{17742, 203}, {17743, 198}, {41656, 179}},
			},
			want: z80State{
				A: 0xAF, F: 0xD4,
				B: 0xDB, C: 0xB6,
				D: 0xCE, E: 0x6B,
				H: 0xA2, L: 0xB8,
				I: 0x1D, R: 0x51,
				PC: 0x4550, SP: 0x0B7D,
				IX: 0x95DF, IY: 0x5D6F,
				AF_: 0xFF23, BC_: 0x1EE8,
				DE_: 0x430B, HL_: 0xC871,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{17742, 203}, {17743, 198}, {41656, 179}},
				Cycles: 15,
			},
		},
		{ // CB F8 0000
			name: "CB F8 0000",
			init: z80State{
				A: 0xC6, F: 0x61,
				B: 0xE2, C: 0xA9,
				D: 0x84, E: 0x3E,
				H: 0xEC, L: 0x2B,
				I: 0xC3, R: 0x01,
				PC: 0x0A40, SP: 0x2AF1,
				IX: 0x8D28, IY: 0x34B4,
				AF_: 0xEB82, BC_: 0x55DF,
				DE_: 0x1510, HL_: 0x6DCB,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{2624, 203}, {2625, 248}},
			},
			want: z80State{
				A: 0xC6, F: 0x61,
				B: 0xE2, C: 0xA9,
				D: 0x84, E: 0x3E,
				H: 0xEC, L: 0x2B,
				I: 0xC3, R: 0x03,
				PC: 0x0A42, SP: 0x2AF1,
				IX: 0x8D28, IY: 0x34B4,
				AF_: 0xEB82, BC_: 0x55DF,
				DE_: 0x1510, HL_: 0x6DCB,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{2624, 203}, {2625, 248}},
				Cycles: 8,
			},
		},
		{ // CB FE 0000
			name: "CB FE 0000",
			init: z80State{
				A: 0xDB, F: 0x14,
				B: 0x07, C: 0xBE,
				D: 0x05, E: 0x0B,
				H: 0x5A, L: 0x67,
				I: 0x60, R: 0x08,
				PC: 0x64D6, SP: 0xC618,
				IX: 0xC2E6, IY: 0xDB00,
				AF_: 0xE2CE, BC_: 0x024E,
				DE_: 0x8F12, HL_: 0x8AD2,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{23143, 199}, {25814, 203}, {25815, 254}},
			},
			want: z80State{
				A: 0xDB, F: 0x14,
				B: 0x07, C: 0xBE,
				D: 0x05, E: 0x0B,
				H: 0x5A, L: 0x67,
				I: 0x60, R: 0x0A,
				PC: 0x64D8, SP: 0xC618,
				IX: 0xC2E6, IY: 0xDB00,
				AF_: 0xE2CE, BC_: 0x024E,
				DE_: 0x8F12, HL_: 0x8AD2,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{23143, 199}, {25814, 203}, {25815, 254}},
				Cycles: 15,
			},
		},
		{ // CB C7 0000
			name: "CB C7 0000",
			init: z80State{
				A: 0xCA, F: 0x77,
				B: 0x21, C: 0xF3,
				D: 0xA4, E: 0xE2,
				H: 0x54, L: 0x4B,
				I: 0x70, R: 0x17,
				PC: 0x21C6, SP: 0x769D,
				IX: 0xD0D3, IY: 0x54C5,
				AF_: 0x21CF, BC_: 0x05CD,
				DE_: 0x500C, HL_: 0x5D33,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{8646, 203}, {8647, 199}},
			},
			want: z80State{
				A: 0xCB, F: 0x77,
				B: 0x21, C: 0xF3,
				D: 0xA4, E: 0xE2,
				H: 0x54, L: 0x4B,
				I: 0x70, R: 0x19,
				PC: 0x21C8, SP: 0x769D,
				IX: 0xD0D3, IY: 0x54C5,
				AF_: 0x21CF, BC_: 0x05CD,
				DE_: 0x500C, HL_: 0x5D33,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{8646, 203}, {8647, 199}},
				Cycles: 8,
			},
		},
	}
	for _, tt := range tests {
		runSSTTest(t, tt)
	}
}
