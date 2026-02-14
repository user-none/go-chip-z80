package z80

import "testing"

func TestDD_LD_r_IXd(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 7E 05 = LD A, (IX+5)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x7E
	bus.mem[2] = 0x05
	cpu.reg.IX = 0x1000
	bus.mem[0x1005] = 0x42
	cycles := cpu.Step()
	if cpu.getA() != 0x42 {
		t.Errorf("LD A,(IX+5): A=%02x want 42", cpu.getA())
	}
	if cycles != 19 {
		t.Errorf("cycles=%d want 19", cycles)
	}
}

func TestDD_LD_IXd_r(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 77 03 = LD (IX+3), A
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x77
	bus.mem[2] = 0x03
	cpu.reg.IX = 0x2000
	cpu.reg.AF = 0xBB00
	cycles := cpu.Step()
	if bus.mem[0x2003] != 0xBB {
		t.Errorf("LD (IX+3),A: got %02x want BB", bus.mem[0x2003])
	}
	if cycles != 19 {
		t.Errorf("cycles=%d want 19", cycles)
	}
}

func TestDD_LD_IXd_n(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 36 02 AA = LD (IX+2), 0xAA
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x36
	bus.mem[2] = 0x02
	bus.mem[3] = 0xAA
	cpu.reg.IX = 0x3000
	cycles := cpu.Step()
	if bus.mem[0x3002] != 0xAA {
		t.Errorf("LD (IX+2),n: got %02x want AA", bus.mem[0x3002])
	}
	if cycles != 19 {
		t.Errorf("cycles=%d want 19", cycles)
	}
}

func TestDD_INC_IXd(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 34 01 = INC (IX+1)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x34
	bus.mem[2] = 0x01
	cpu.reg.IX = 0x4000
	bus.mem[0x4001] = 0x0F
	cycles := cpu.Step()
	if bus.mem[0x4001] != 0x10 {
		t.Errorf("INC (IX+1): got %02x want 10", bus.mem[0x4001])
	}
	if cycles != 23 {
		t.Errorf("cycles=%d want 23", cycles)
	}
}

func TestDD_DEC_IXd(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 35 02 = DEC (IX+2)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x35
	bus.mem[2] = 0x02
	cpu.reg.IX = 0x4000
	bus.mem[0x4002] = 0x10
	cycles := cpu.Step()
	if bus.mem[0x4002] != 0x0F {
		t.Errorf("DEC (IX+2): got %02x want 0F", bus.mem[0x4002])
	}
	if cycles != 23 {
		t.Errorf("cycles=%d want 23", cycles)
	}
}

func TestDD_ALU_IXd(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 86 03 = ADD A, (IX+3)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x86
	bus.mem[2] = 0x03
	cpu.reg.IX = 0x5000
	cpu.reg.AF = 0x1000
	bus.mem[0x5003] = 0x20
	cycles := cpu.Step()
	if cpu.getA() != 0x30 {
		t.Errorf("ADD A,(IX+3): A=%02x want 30", cpu.getA())
	}
	if cycles != 19 {
		t.Errorf("cycles=%d want 19", cycles)
	}
}

func TestDD_NegativeDisplacement(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 7E FE = LD A, (IX-2)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x7E
	bus.mem[2] = 0xFE // -2
	cpu.reg.IX = 0x1005
	bus.mem[0x1003] = 0x77
	cpu.Step()
	if cpu.getA() != 0x77 {
		t.Errorf("LD A,(IX-2): A=%02x want 77", cpu.getA())
	}
}

func TestFD_LD_r_IYd(t *testing.T) {
	cpu, bus := newTestCPU()
	// FD 46 04 = LD B, (IY+4)
	bus.mem[0] = 0xFD
	bus.mem[1] = 0x46
	bus.mem[2] = 0x04
	cpu.reg.IY = 0x6000
	bus.mem[0x6004] = 0x55
	cycles := cpu.Step()
	if cpu.getB() != 0x55 {
		t.Errorf("LD B,(IY+4): B=%02x want 55", cpu.getB())
	}
	if cycles != 19 {
		t.Errorf("cycles=%d want 19", cycles)
	}
}

func TestDD_Fallthrough_Register(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 47 = LD B, A (falls through to baseOps since no (HL) involved)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x47
	cpu.reg.AF = 0x9900
	cpu.Step()
	if cpu.getB() != 0x99 {
		t.Errorf("DD fallthrough LD B,A: B=%02x want 99", cpu.getB())
	}
}

func TestDD_IXH_IXL(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 44 = LD B, IXH (undocumented, via ixiyReg routing)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x44
	cpu.reg.IX = 0xAB00
	cpu.Step()
	if cpu.getB() != 0xAB {
		t.Errorf("LD B,IXH: B=%02x want AB", cpu.getB())
	}
}

func TestDD_LD_IXL(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 45 = LD B, IXL (undocumented)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x45
	cpu.reg.IX = 0x00CD
	cpu.Step()
	if cpu.getB() != 0xCD {
		t.Errorf("LD B,IXL: B=%02x want CD", cpu.getB())
	}
}

func TestDD_ADD_IX_rr(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 09 = ADD IX, BC
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x09
	cpu.reg.IX = 0x1000
	cpu.reg.BC = 0x2000
	cycles := cpu.Step()
	if cpu.reg.IX != 0x3000 {
		t.Errorf("ADD IX,BC: IX=%04x want 3000", cpu.reg.IX)
	}
	if cycles != 15 {
		t.Errorf("cycles=%d want 15", cycles)
	}
}

func TestDD_PUSH_POP_IX(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.SP = 0xFFFE
	cpu.reg.IX = 0xBEEF
	// PUSH IX (DD E5)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0xE5
	cpu.Step()
	if cpu.reg.SP != 0xFFFC {
		t.Errorf("PUSH IX: SP=%04x want FFFC", cpu.reg.SP)
	}
	// POP IY (FD E1)
	bus.mem[2] = 0xFD
	bus.mem[3] = 0xE1
	cpu.Step()
	if cpu.reg.IY != 0xBEEF {
		t.Errorf("POP IY: IY=%04x want BEEF", cpu.reg.IY)
	}
}

func TestDD_LD_IX_nn(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD 21 34 12 = LD IX, 0x1234
	bus.mem[0] = 0xDD
	bus.mem[1] = 0x21
	bus.mem[2] = 0x34
	bus.mem[3] = 0x12
	cpu.Step()
	if cpu.reg.IX != 0x1234 {
		t.Errorf("LD IX,nn: IX=%04x want 1234", cpu.reg.IX)
	}
}

func TestDDCB_BIT(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD CB 05 46 = BIT 0, (IX+5)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0xCB
	bus.mem[2] = 0x05 // displacement
	bus.mem[3] = 0x46 // BIT 0
	cpu.reg.IX = 0x1000
	bus.mem[0x1005] = 0xFE // bit 0 clear
	cycles := cpu.Step()
	if cpu.getF()&flagZ == 0 {
		t.Error("BIT 0,(IX+5): Z should be set")
	}
	if cycles != 20 {
		t.Errorf("cycles=%d want 20", cycles)
	}
}

func TestDDCB_RLC(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD CB 02 06 = RLC (IX+2)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0xCB
	bus.mem[2] = 0x02
	bus.mem[3] = 0x06
	cpu.reg.IX = 0x2000
	bus.mem[0x2002] = 0x80
	cycles := cpu.Step()
	if bus.mem[0x2002] != 0x01 {
		t.Errorf("RLC (IX+2): got %02x want 01", bus.mem[0x2002])
	}
	if cycles != 23 {
		t.Errorf("cycles=%d want 23", cycles)
	}
}

func TestDDCB_SET(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD CB 01 C6 = SET 0, (IX+1)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0xCB
	bus.mem[2] = 0x01
	bus.mem[3] = 0xC6
	cpu.reg.IX = 0x3000
	bus.mem[0x3001] = 0x00
	cycles := cpu.Step()
	if bus.mem[0x3001] != 0x01 {
		t.Errorf("SET 0,(IX+1): got %02x want 01", bus.mem[0x3001])
	}
	if cycles != 23 {
		t.Errorf("cycles=%d want 23", cycles)
	}
}

func TestDDCB_RES(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD CB 03 86 = RES 0, (IX+3)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0xCB
	bus.mem[2] = 0x03
	bus.mem[3] = 0x86
	cpu.reg.IX = 0x4000
	bus.mem[0x4003] = 0xFF
	cycles := cpu.Step()
	if bus.mem[0x4003] != 0xFE {
		t.Errorf("RES 0,(IX+3): got %02x want FE", bus.mem[0x4003])
	}
	if cycles != 23 {
		t.Errorf("cycles=%d want 23", cycles)
	}
}

func TestDDCB_Undocumented_RotateToReg(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD CB 00 00 = RLC (IX+0) -> B (undocumented)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0xCB
	bus.mem[2] = 0x00
	bus.mem[3] = 0x00 // RLC, dest=B
	cpu.reg.IX = 0x5000
	bus.mem[0x5000] = 0x80
	cpu.Step()
	if bus.mem[0x5000] != 0x01 {
		t.Errorf("(IX+0)=%02x want 01", bus.mem[0x5000])
	}
	if cpu.getB() != 0x01 {
		t.Errorf("B=%02x want 01 (undocumented copy)", cpu.getB())
	}
}

func TestDD_JP_IX(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD E9 = JP (IX)
	bus.mem[0] = 0xDD
	bus.mem[1] = 0xE9
	cpu.reg.IX = 0x9876
	cpu.Step()
	if cpu.reg.PC != 0x9876 {
		t.Errorf("JP (IX): PC=%04x want 9876", cpu.reg.PC)
	}
}

func TestDD_LD_SP_IX(t *testing.T) {
	cpu, bus := newTestCPU()
	// DD F9 = LD SP, IX
	bus.mem[0] = 0xDD
	bus.mem[1] = 0xF9
	cpu.reg.IX = 0x4444
	cpu.Step()
	if cpu.reg.SP != 0x4444 {
		t.Errorf("LD SP,IX: SP=%04x want 4444", cpu.reg.SP)
	}
}

func TestSST_IX(t *testing.T) {
	tests := []struct {
		name       string
		init, want z80State
	}{
		{ // DD 7E 0000
			name: "DD 7E 0000",
			init: z80State{
				A: 0x9E, F: 0xCA,
				B: 0x3B, C: 0xC1,
				D: 0xF9, E: 0x28,
				H: 0x96, L: 0x54,
				I: 0x64, R: 0x55,
				PC: 0xA263, SP: 0x4499,
				IX: 0x2936, IY: 0x0340,
				AF_: 0xDB17, BC_: 0x5E01,
				DE_: 0x2BCD, HL_: 0x1726,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{10591, 14}, {41571, 221}, {41572, 126}, {41573, 41}},
			},
			want: z80State{
				A: 0x0E, F: 0xCA,
				B: 0x3B, C: 0xC1,
				D: 0xF9, E: 0x28,
				H: 0x96, L: 0x54,
				I: 0x64, R: 0x57,
				PC: 0xA266, SP: 0x4499,
				IX: 0x2936, IY: 0x0340,
				AF_: 0xDB17, BC_: 0x5E01,
				DE_: 0x2BCD, HL_: 0x1726,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{10591, 14}, {41571, 221}, {41572, 126}, {41573, 41}},
				Cycles: 19,
			},
		},
		{ // DD 46 0000
			name: "DD 46 0000",
			init: z80State{
				A: 0x2A, F: 0xAA,
				B: 0x90, C: 0x41,
				D: 0xF8, E: 0x0A,
				H: 0x00, L: 0x24,
				I: 0xA4, R: 0x00,
				PC: 0x7BEC, SP: 0xC284,
				IX: 0xC3CA, IY: 0xEDC4,
				AF_: 0xC60F, BC_: 0x0809,
				DE_: 0x363A, HL_: 0x7506,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{31724, 221}, {31725, 70}, {31726, 75}, {50197, 115}},
			},
			want: z80State{
				A: 0x2A, F: 0xAA,
				B: 0x73, C: 0x41,
				D: 0xF8, E: 0x0A,
				H: 0x00, L: 0x24,
				I: 0xA4, R: 0x02,
				PC: 0x7BEF, SP: 0xC284,
				IX: 0xC3CA, IY: 0xEDC4,
				AF_: 0xC60F, BC_: 0x0809,
				DE_: 0x363A, HL_: 0x7506,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{31724, 221}, {31725, 70}, {31726, 75}, {50197, 115}},
				Cycles: 19,
			},
		},
		{ // DD 77 0000
			name: "DD 77 0000",
			init: z80State{
				A: 0x9B, F: 0x41,
				B: 0x1F, C: 0xDE,
				D: 0x42, E: 0x49,
				H: 0x40, L: 0x92,
				I: 0xE9, R: 0x4D,
				PC: 0xAFC5, SP: 0x48DA,
				IX: 0x4CF3, IY: 0x6037,
				AF_: 0xBC7E, BC_: 0x8813,
				DE_: 0xB0E5, HL_: 0x734D,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{19596, 0}, {44997, 221}, {44998, 119}, {44999, 153}},
			},
			want: z80State{
				A: 0x9B, F: 0x41,
				B: 0x1F, C: 0xDE,
				D: 0x42, E: 0x49,
				H: 0x40, L: 0x92,
				I: 0xE9, R: 0x4F,
				PC: 0xAFC8, SP: 0x48DA,
				IX: 0x4CF3, IY: 0x6037,
				AF_: 0xBC7E, BC_: 0x8813,
				DE_: 0xB0E5, HL_: 0x734D,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{19596, 155}, {44997, 221}, {44998, 119}, {44999, 153}},
				Cycles: 19,
			},
		},
		{ // DD 36 0000
			name: "DD 36 0000",
			init: z80State{
				A: 0x30, F: 0x21,
				B: 0xFC, C: 0xC3,
				D: 0xFB, E: 0x32,
				H: 0xDE, L: 0xC6,
				I: 0xCF, R: 0x7E,
				PC: 0x387F, SP: 0x41B5,
				IX: 0xF602, IY: 0xDB4F,
				AF_: 0x15A2, BC_: 0xFB2D,
				DE_: 0xD1C4, HL_: 0x24AF,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{14463, 221}, {14464, 54}, {14465, 20}, {14466, 47}, {62998, 0}},
			},
			want: z80State{
				A: 0x30, F: 0x21,
				B: 0xFC, C: 0xC3,
				D: 0xFB, E: 0x32,
				H: 0xDE, L: 0xC6,
				I: 0xCF, R: 0x00,
				PC: 0x3883, SP: 0x41B5,
				IX: 0xF602, IY: 0xDB4F,
				AF_: 0x15A2, BC_: 0xFB2D,
				DE_: 0xD1C4, HL_: 0x24AF,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{14463, 221}, {14464, 54}, {14465, 20}, {14466, 47}, {62998, 47}},
				Cycles: 19,
			},
		},
		{ // DD 34 0000
			name: "DD 34 0000",
			init: z80State{
				A: 0xCF, F: 0xF8,
				B: 0xAD, C: 0x2B,
				D: 0x22, E: 0xBC,
				H: 0xF6, L: 0xAE,
				I: 0x68, R: 0x18,
				PC: 0xAE96, SP: 0x7B75,
				IX: 0xEF3D, IY: 0xDCCE,
				AF_: 0x44CE, BC_: 0x015B,
				DE_: 0xF278, HL_: 0xAF45,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{44694, 221}, {44695, 52}, {44696, 166}, {61155, 191}},
			},
			want: z80State{
				A: 0xCF, F: 0x90,
				B: 0xAD, C: 0x2B,
				D: 0x22, E: 0xBC,
				H: 0xF6, L: 0xAE,
				I: 0x68, R: 0x1A,
				PC: 0xAE99, SP: 0x7B75,
				IX: 0xEF3D, IY: 0xDCCE,
				AF_: 0x44CE, BC_: 0x015B,
				DE_: 0xF278, HL_: 0xAF45,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{44694, 221}, {44695, 52}, {44696, 166}, {61155, 192}},
				Cycles: 23,
			},
		},
		{ // DD 35 0000
			name: "DD 35 0000",
			init: z80State{
				A: 0xD7, F: 0x36,
				B: 0x17, C: 0x26,
				D: 0xE2, E: 0x6C,
				H: 0x85, L: 0xCC,
				I: 0x01, R: 0x2F,
				PC: 0x0FF6, SP: 0x5B69,
				IX: 0x2325, IY: 0xD028,
				AF_: 0x3452, BC_: 0x2F85,
				DE_: 0x6E59, HL_: 0xF10B,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{4086, 221}, {4087, 53}, {4088, 56}, {9053, 108}},
			},
			want: z80State{
				A: 0xD7, F: 0x2A,
				B: 0x17, C: 0x26,
				D: 0xE2, E: 0x6C,
				H: 0x85, L: 0xCC,
				I: 0x01, R: 0x31,
				PC: 0x0FF9, SP: 0x5B69,
				IX: 0x2325, IY: 0xD028,
				AF_: 0x3452, BC_: 0x2F85,
				DE_: 0x6E59, HL_: 0xF10B,
				IM: 0, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{4086, 221}, {4087, 53}, {4088, 56}, {9053, 107}},
				Cycles: 23,
			},
		},
		{ // DD 86 0000
			name: "DD 86 0000",
			init: z80State{
				A: 0x0B, F: 0xEE,
				B: 0x69, C: 0x91,
				D: 0x01, E: 0x9A,
				H: 0x4B, L: 0x9A,
				I: 0x68, R: 0x26,
				PC: 0x3A1C, SP: 0x07A0,
				IX: 0x8FF5, IY: 0xABAE,
				AF_: 0xC5D1, BC_: 0x8336,
				DE_: 0xC03B, HL_: 0xE0FB,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{14876, 221}, {14877, 134}, {14878, 88}, {36941, 26}},
			},
			want: z80State{
				A: 0x25, F: 0x30,
				B: 0x69, C: 0x91,
				D: 0x01, E: 0x9A,
				H: 0x4B, L: 0x9A,
				I: 0x68, R: 0x28,
				PC: 0x3A1F, SP: 0x07A0,
				IX: 0x8FF5, IY: 0xABAE,
				AF_: 0xC5D1, BC_: 0x8336,
				DE_: 0xC03B, HL_: 0xE0FB,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{14876, 221}, {14877, 134}, {14878, 88}, {36941, 26}},
				Cycles: 19,
			},
		},
		{ // DD 96 0000
			name: "DD 96 0000",
			init: z80State{
				A: 0xBA, F: 0xE5,
				B: 0xA1, C: 0x10,
				D: 0xB8, E: 0x4B,
				H: 0x83, L: 0x12,
				I: 0x3B, R: 0x22,
				PC: 0xF687, SP: 0xA864,
				IX: 0xC142, IY: 0xCEB6,
				AF_: 0xFC4D, BC_: 0x6736,
				DE_: 0x1FD0, HL_: 0x4683,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{49491, 140}, {63111, 221}, {63112, 150}, {63113, 17}},
			},
			want: z80State{
				A: 0x2E, F: 0x3A,
				B: 0xA1, C: 0x10,
				D: 0xB8, E: 0x4B,
				H: 0x83, L: 0x12,
				I: 0x3B, R: 0x24,
				PC: 0xF68A, SP: 0xA864,
				IX: 0xC142, IY: 0xCEB6,
				AF_: 0xFC4D, BC_: 0x6736,
				DE_: 0x1FD0, HL_: 0x4683,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{49491, 140}, {63111, 221}, {63112, 150}, {63113, 17}},
				Cycles: 19,
			},
		},
		{ // DD BE 0000
			name: "DD BE 0000",
			init: z80State{
				A: 0x94, F: 0x8A,
				B: 0x26, C: 0xDB,
				D: 0x30, E: 0x45,
				H: 0xC0, L: 0x2F,
				I: 0x8C, R: 0x75,
				PC: 0x58B2, SP: 0x63AB,
				IX: 0x8DB2, IY: 0x2074,
				AF_: 0x3DA0, BC_: 0x313E,
				DE_: 0xF67D, HL_: 0x3BF4,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{22706, 221}, {22707, 190}, {22708, 75}, {36349, 137}},
			},
			want: z80State{
				A: 0x94, F: 0x1A,
				B: 0x26, C: 0xDB,
				D: 0x30, E: 0x45,
				H: 0xC0, L: 0x2F,
				I: 0x8C, R: 0x77,
				PC: 0x58B5, SP: 0x63AB,
				IX: 0x8DB2, IY: 0x2074,
				AF_: 0x3DA0, BC_: 0x313E,
				DE_: 0xF67D, HL_: 0x3BF4,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{22706, 221}, {22707, 190}, {22708, 75}, {36349, 137}},
				Cycles: 19,
			},
		},
		{ // FD 7E 0000
			name: "FD 7E 0000",
			init: z80State{
				A: 0x7B, F: 0xF3,
				B: 0xEE, C: 0x9D,
				D: 0x8D, E: 0xB4,
				H: 0xEC, L: 0xC4,
				I: 0x51, R: 0x27,
				PC: 0x4A67, SP: 0xE174,
				IX: 0x9D6B, IY: 0x42DD,
				AF_: 0xDB5C, BC_: 0x26E5,
				DE_: 0x2FBB, HL_: 0x9BA7,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{17106, 194}, {19047, 253}, {19048, 126}, {19049, 245}},
			},
			want: z80State{
				A: 0xC2, F: 0xF3,
				B: 0xEE, C: 0x9D,
				D: 0x8D, E: 0xB4,
				H: 0xEC, L: 0xC4,
				I: 0x51, R: 0x29,
				PC: 0x4A6A, SP: 0xE174,
				IX: 0x9D6B, IY: 0x42DD,
				AF_: 0xDB5C, BC_: 0x26E5,
				DE_: 0x2FBB, HL_: 0x9BA7,
				IM: 2, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{17106, 194}, {19047, 253}, {19048, 126}, {19049, 245}},
				Cycles: 19,
			},
		},
		{ // FD 46 0000
			name: "FD 46 0000",
			init: z80State{
				A: 0x1A, F: 0x01,
				B: 0xFE, C: 0x63,
				D: 0xFF, E: 0xA6,
				H: 0x66, L: 0x49,
				I: 0x6A, R: 0x5D,
				PC: 0x1E13, SP: 0xB7A6,
				IX: 0xE22C, IY: 0x8BE6,
				AF_: 0x0E32, BC_: 0xF61F,
				DE_: 0xC892, HL_: 0xF017,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{7699, 253}, {7700, 70}, {7701, 175}, {35733, 78}},
			},
			want: z80State{
				A: 0x1A, F: 0x01,
				B: 0x4E, C: 0x63,
				D: 0xFF, E: 0xA6,
				H: 0x66, L: 0x49,
				I: 0x6A, R: 0x5F,
				PC: 0x1E16, SP: 0xB7A6,
				IX: 0xE22C, IY: 0x8BE6,
				AF_: 0x0E32, BC_: 0xF61F,
				DE_: 0xC892, HL_: 0xF017,
				IM: 2, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{7699, 253}, {7700, 70}, {7701, 175}, {35733, 78}},
				Cycles: 19,
			},
		},
		{ // DD 21 0000
			name: "DD 21 0000",
			init: z80State{
				A: 0xFA, F: 0x7F,
				B: 0x96, C: 0xC1,
				D: 0x11, E: 0xB9,
				H: 0x72, L: 0x48,
				I: 0x01, R: 0x1F,
				PC: 0x2721, SP: 0x8F99,
				IX: 0x5435, IY: 0xD754,
				AF_: 0xC443, BC_: 0x14B4,
				DE_: 0xD87B, HL_: 0xE99D,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{10017, 221}, {10018, 33}, {10019, 131}, {10020, 191}},
			},
			want: z80State{
				A: 0xFA, F: 0x7F,
				B: 0x96, C: 0xC1,
				D: 0x11, E: 0xB9,
				H: 0x72, L: 0x48,
				I: 0x01, R: 0x21,
				PC: 0x2725, SP: 0x8F99,
				IX: 0xBF83, IY: 0xD754,
				AF_: 0xC443, BC_: 0x14B4,
				DE_: 0xD87B, HL_: 0xE99D,
				IM: 1, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{10017, 221}, {10018, 33}, {10019, 131}, {10020, 191}},
				Cycles: 14,
			},
		},
		{ // DD 09 0000
			name: "DD 09 0000",
			init: z80State{
				A: 0x59, F: 0xA9,
				B: 0xFA, C: 0xAE,
				D: 0xC9, E: 0xD8,
				H: 0xB1, L: 0x1A,
				I: 0xFC, R: 0x0E,
				PC: 0x9857, SP: 0x1B12,
				IX: 0x3A5A, IY: 0xAD4A,
				AF_: 0x410E, BC_: 0xCB8A,
				DE_: 0x528C, HL_: 0xC94E,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{38999, 221}, {39000, 9}},
			},
			want: z80State{
				A: 0x59, F: 0xB1,
				B: 0xFA, C: 0xAE,
				D: 0xC9, E: 0xD8,
				H: 0xB1, L: 0x1A,
				I: 0xFC, R: 0x10,
				PC: 0x9859, SP: 0x1B12,
				IX: 0x3508, IY: 0xAD4A,
				AF_: 0x410E, BC_: 0xCB8A,
				DE_: 0x528C, HL_: 0xC94E,
				IM: 1, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{38999, 221}, {39000, 9}},
				Cycles: 15,
			},
		},
		{ // DD CB __ 06 0000
			name: "DD CB __ 06 0000",
			init: z80State{
				A: 0x22, F: 0x36,
				B: 0x7C, C: 0x45,
				D: 0x01, E: 0xDD,
				H: 0xC7, L: 0x5D,
				I: 0xA3, R: 0x59,
				PC: 0x211F, SP: 0x52A1,
				IX: 0x62D3, IY: 0x5C3B,
				AF_: 0x0860, BC_: 0xF319,
				DE_: 0x4F78, HL_: 0x7E91,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{8479, 221}, {8480, 203}, {8481, 253}, {8482, 6}, {25296, 81}},
			},
			want: z80State{
				A: 0x22, F: 0xA0,
				B: 0x7C, C: 0x45,
				D: 0x01, E: 0xDD,
				H: 0xC7, L: 0x5D,
				I: 0xA3, R: 0x5B,
				PC: 0x2123, SP: 0x52A1,
				IX: 0x62D3, IY: 0x5C3B,
				AF_: 0x0860, BC_: 0xF319,
				DE_: 0x4F78, HL_: 0x7E91,
				IM: 2, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{8479, 221}, {8480, 203}, {8481, 253}, {8482, 6}, {25296, 162}},
				Cycles: 23,
			},
		},
		{ // DD CB __ 0E 0000
			name: "DD CB __ 0E 0000",
			init: z80State{
				A: 0x5E, F: 0xFA,
				B: 0x98, C: 0xEB,
				D: 0xF3, E: 0x1C,
				H: 0x26, L: 0xB7,
				I: 0x23, R: 0x4C,
				PC: 0x4BFF, SP: 0xCFF4,
				IX: 0xBECC, IY: 0x0553,
				AF_: 0x91F3, BC_: 0x47CF,
				DE_: 0x9C76, HL_: 0xFFCC,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{19455, 221}, {19456, 203}, {19457, 32}, {19458, 14}, {48876, 116}},
			},
			want: z80State{
				A: 0x5E, F: 0x2C,
				B: 0x98, C: 0xEB,
				D: 0xF3, E: 0x1C,
				H: 0x26, L: 0xB7,
				I: 0x23, R: 0x4E,
				PC: 0x4C03, SP: 0xCFF4,
				IX: 0xBECC, IY: 0x0553,
				AF_: 0x91F3, BC_: 0x47CF,
				DE_: 0x9C76, HL_: 0xFFCC,
				IM: 2, IFF1: false, IFF2: false,
				RAM: [][2]uint16{{19455, 221}, {19456, 203}, {19457, 32}, {19458, 14}, {48876, 58}},
				Cycles: 23,
			},
		},
		{ // DD CB __ 16 0000
			name: "DD CB __ 16 0000",
			init: z80State{
				A: 0xCA, F: 0x33,
				B: 0x42, C: 0x68,
				D: 0xB4, E: 0x79,
				H: 0x9E, L: 0x3A,
				I: 0xA9, R: 0x61,
				PC: 0xE341, SP: 0xB741,
				IX: 0xAB1A, IY: 0x74EB,
				AF_: 0x3A15, BC_: 0xCC72,
				DE_: 0xC478, HL_: 0x4B74,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{43736, 156}, {58177, 221}, {58178, 203}, {58179, 190}, {58180, 22}},
			},
			want: z80State{
				A: 0xCA, F: 0x2D,
				B: 0x42, C: 0x68,
				D: 0xB4, E: 0x79,
				H: 0x9E, L: 0x3A,
				I: 0xA9, R: 0x63,
				PC: 0xE345, SP: 0xB741,
				IX: 0xAB1A, IY: 0x74EB,
				AF_: 0x3A15, BC_: 0xCC72,
				DE_: 0xC478, HL_: 0x4B74,
				IM: 0, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{43736, 57}, {58177, 221}, {58178, 203}, {58179, 190}, {58180, 22}},
				Cycles: 23,
			},
		},
		{ // DD CB __ 46 0000
			name: "DD CB __ 46 0000",
			init: z80State{
				A: 0xF6, F: 0xB2,
				B: 0xEF, C: 0xC8,
				D: 0xBB, E: 0x6B,
				H: 0xDE, L: 0xCF,
				I: 0xED, R: 0x67,
				PC: 0xFD0E, SP: 0x7616,
				IX: 0xDB03, IY: 0x0545,
				AF_: 0x05C7, BC_: 0x71FB,
				DE_: 0x36A5, HL_: 0xAD94,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{55971, 95}, {64782, 221}, {64783, 203}, {64784, 160}, {64785, 70}},
			},
			want: z80State{
				A: 0xF6, F: 0x18,
				B: 0xEF, C: 0xC8,
				D: 0xBB, E: 0x6B,
				H: 0xDE, L: 0xCF,
				I: 0xED, R: 0x69,
				PC: 0xFD12, SP: 0x7616,
				IX: 0xDB03, IY: 0x0545,
				AF_: 0x05C7, BC_: 0x71FB,
				DE_: 0x36A5, HL_: 0xAD94,
				IM: 0, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{55971, 95}, {64782, 221}, {64783, 203}, {64784, 160}, {64785, 70}},
				Cycles: 20,
			},
		},
		{ // DD CB __ 86 0000
			name: "DD CB __ 86 0000",
			init: z80State{
				A: 0x09, F: 0x75,
				B: 0xE1, C: 0x05,
				D: 0x0F, E: 0xFB,
				H: 0x06, L: 0x97,
				I: 0x7D, R: 0x24,
				PC: 0xE92B, SP: 0x8396,
				IX: 0x5905, IY: 0xC6D8,
				AF_: 0x37C2, BC_: 0xD24C,
				DE_: 0x8618, HL_: 0x1824,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{22726, 247}, {59691, 221}, {59692, 203}, {59693, 193}, {59694, 134}},
			},
			want: z80State{
				A: 0x09, F: 0x75,
				B: 0xE1, C: 0x05,
				D: 0x0F, E: 0xFB,
				H: 0x06, L: 0x97,
				I: 0x7D, R: 0x26,
				PC: 0xE92F, SP: 0x8396,
				IX: 0x5905, IY: 0xC6D8,
				AF_: 0x37C2, BC_: 0xD24C,
				DE_: 0x8618, HL_: 0x1824,
				IM: 1, IFF1: false, IFF2: true,
				RAM: [][2]uint16{{22726, 246}, {59691, 221}, {59692, 203}, {59693, 193}, {59694, 134}},
				Cycles: 23,
			},
		},
		{ // DD CB __ C6 0000
			name: "DD CB __ C6 0000",
			init: z80State{
				A: 0x36, F: 0x79,
				B: 0xE6, C: 0xAB,
				D: 0x60, E: 0x01,
				H: 0x6B, L: 0x09,
				I: 0x52, R: 0x63,
				PC: 0x61E5, SP: 0xD1E1,
				IX: 0x593E, IY: 0xF904,
				AF_: 0xD139, BC_: 0x6558,
				DE_: 0xD259, HL_: 0xB735,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{22902, 76}, {25061, 221}, {25062, 203}, {25063, 56}, {25064, 198}},
			},
			want: z80State{
				A: 0x36, F: 0x79,
				B: 0xE6, C: 0xAB,
				D: 0x60, E: 0x01,
				H: 0x6B, L: 0x09,
				I: 0x52, R: 0x65,
				PC: 0x61E9, SP: 0xD1E1,
				IX: 0x593E, IY: 0xF904,
				AF_: 0xD139, BC_: 0x6558,
				DE_: 0xD259, HL_: 0xB735,
				IM: 1, IFF1: true, IFF2: false,
				RAM: [][2]uint16{{22902, 77}, {25061, 221}, {25062, 203}, {25063, 56}, {25064, 198}},
				Cycles: 23,
			},
		},
		{ // DD CB __ 00 0000
			name: "DD CB __ 00 0000",
			init: z80State{
				A: 0x61, F: 0x4F,
				B: 0xFB, C: 0x65,
				D: 0x3F, E: 0x67,
				H: 0x80, L: 0x72,
				I: 0x60, R: 0x1F,
				PC: 0xBAB2, SP: 0xE22D,
				IX: 0x0CF9, IY: 0x4B05,
				AF_: 0x984A, BC_: 0xCEF5,
				DE_: 0x06AD, HL_: 0x172E,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{3307, 251}, {47794, 221}, {47795, 203}, {47796, 242}, {47797, 0}},
			},
			want: z80State{
				A: 0x61, F: 0xA1,
				B: 0xF7, C: 0x65,
				D: 0x3F, E: 0x67,
				H: 0x80, L: 0x72,
				I: 0x60, R: 0x21,
				PC: 0xBAB6, SP: 0xE22D,
				IX: 0x0CF9, IY: 0x4B05,
				AF_: 0x984A, BC_: 0xCEF5,
				DE_: 0x06AD, HL_: 0x172E,
				IM: 0, IFF1: true, IFF2: true,
				RAM: [][2]uint16{{3307, 247}, {47794, 221}, {47795, 203}, {47796, 242}, {47797, 0}},
				Cycles: 23,
			},
		},
	}
	for _, tt := range tests {
		runSSTTest(t, tt)
	}
}
