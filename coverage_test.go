package z80

import (
	"fmt"
	"testing"
)

// TestBaseOps_NilCheck verifies every baseOps entry is non-nil.
func TestBaseOps_NilCheck(t *testing.T) {
	for i := 0; i < 256; i++ {
		if baseOps[i] == nil {
			t.Errorf("baseOps[0x%02X] is nil", i)
		}
	}
}

// TestCbOps_NilCheck verifies every cbOps entry is non-nil.
func TestCbOps_NilCheck(t *testing.T) {
	for i := 0; i < 256; i++ {
		if cbOps[i] == nil {
			t.Errorf("cbOps[0x%02X] is nil", i)
		}
	}
}

// TestIxcbOps_NilCheck verifies every ixcbOps entry is non-nil.
func TestIxcbOps_NilCheck(t *testing.T) {
	for i := 0; i < 256; i++ {
		if ixcbOps[i] == nil {
			t.Errorf("ixcbOps[0x%02X] is nil", i)
		}
	}
}

// TestEdOps_DocumentedCheck verifies all documented ED opcodes are registered.
func TestEdOps_DocumentedCheck(t *testing.T) {
	documented := map[uint8]string{
		// IN r,(C)
		0x40: "IN B,(C)", 0x48: "IN C,(C)", 0x50: "IN D,(C)", 0x58: "IN E,(C)",
		0x60: "IN H,(C)", 0x68: "IN L,(C)", 0x70: "IN F,(C)", 0x78: "IN A,(C)",
		// OUT (C),r
		0x41: "OUT (C),B", 0x49: "OUT (C),C", 0x51: "OUT (C),D", 0x59: "OUT (C),E",
		0x61: "OUT (C),H", 0x69: "OUT (C),L", 0x71: "OUT (C),0", 0x79: "OUT (C),A",
		// SBC HL,rr
		0x42: "SBC HL,BC", 0x52: "SBC HL,DE", 0x62: "SBC HL,HL", 0x72: "SBC HL,SP",
		// ADC HL,rr
		0x4A: "ADC HL,BC", 0x5A: "ADC HL,DE", 0x6A: "ADC HL,HL", 0x7A: "ADC HL,SP",
		// LD (nn),rr
		0x43: "LD (nn),BC", 0x53: "LD (nn),DE", 0x63: "LD (nn),HL", 0x73: "LD (nn),SP",
		// LD rr,(nn)
		0x4B: "LD BC,(nn)", 0x5B: "LD DE,(nn)", 0x6B: "LD HL,(nn)", 0x7B: "LD SP,(nn)",
		// NEG
		0x44: "NEG", 0x4C: "NEG*", 0x54: "NEG*", 0x5C: "NEG*",
		0x64: "NEG*", 0x6C: "NEG*", 0x74: "NEG*", 0x7C: "NEG*",
		// RETN/RETI
		0x45: "RETN", 0x55: "RETN*", 0x65: "RETN*", 0x75: "RETN*",
		0x4D: "RETI", 0x5D: "RETI*", 0x6D: "RETI*", 0x7D: "RETI*",
		// IM
		0x46: "IM 0", 0x4E: "IM 0*", 0x66: "IM 0*", 0x6E: "IM 0*",
		0x56: "IM 1", 0x76: "IM 1*",
		0x5E: "IM 2", 0x7E: "IM 2*",
		// LD I/R
		0x47: "LD I,A", 0x4F: "LD R,A", 0x57: "LD A,I", 0x5F: "LD A,R",
		// Rotate
		0x67: "RRD", 0x6F: "RLD",
		// Block transfer
		0xA0: "LDI", 0xA8: "LDD", 0xB0: "LDIR", 0xB8: "LDDR",
		// Block compare
		0xA1: "CPI", 0xA9: "CPD", 0xB1: "CPIR", 0xB9: "CPDR",
		// Block I/O
		0xA2: "INI", 0xAA: "IND", 0xB2: "INIR", 0xBA: "INDR",
		0xA3: "OUTI", 0xAB: "OUTD", 0xB3: "OTIR", 0xBB: "OTDR",
	}

	for op, name := range documented {
		if edOps[op] == nil {
			t.Errorf("edOps[0x%02X] (%s) is nil", op, name)
		}
	}
}

// TestIxOps_DocumentedCheck verifies all (IX+d) variants are registered.
func TestIxOps_DocumentedCheck(t *testing.T) {
	documented := map[uint8]string{
		// LD r,(IX+d)
		0x46: "LD B,(IX+d)", 0x4E: "LD C,(IX+d)", 0x56: "LD D,(IX+d)",
		0x5E: "LD E,(IX+d)", 0x66: "LD H,(IX+d)", 0x6E: "LD L,(IX+d)",
		0x7E: "LD A,(IX+d)",
		// LD (IX+d),r
		0x70: "LD (IX+d),B", 0x71: "LD (IX+d),C", 0x72: "LD (IX+d),D",
		0x73: "LD (IX+d),E", 0x74: "LD (IX+d),H", 0x75: "LD (IX+d),L",
		0x77: "LD (IX+d),A",
		// LD (IX+d),n
		0x36: "LD (IX+d),n",
		// INC/DEC (IX+d)
		0x34: "INC (IX+d)", 0x35: "DEC (IX+d)",
		// ALU A,(IX+d)
		0x86: "ADD A,(IX+d)", 0x8E: "ADC A,(IX+d)",
		0x96: "SUB (IX+d)", 0x9E: "SBC A,(IX+d)",
		0xA6: "AND (IX+d)", 0xAE: "XOR (IX+d)",
		0xB6: "OR (IX+d)", 0xBE: "CP (IX+d)",
	}

	for op, name := range documented {
		if ixOps[op] == nil {
			t.Errorf("ixOps[0x%02X] (%s) is nil", op, name)
		}
	}
}

// TestAllOpcodes_NoPanic executes every unprefixed opcode to ensure no nil panics.
func TestAllOpcodes_NoPanic(t *testing.T) {
	for op := 0; op < 256; op++ {
		// Skip prefixes — they'd chain into more opcodes
		if op == 0xCB || op == 0xDD || op == 0xED || op == 0xFD {
			continue
		}
		t.Run(fmt.Sprintf("op_%02X", op), func(t *testing.T) {
			bus := &testBus{}
			cpu := New(bus)
			cpu.reg.SP = 0xFFF0 // Avoid stack underflow
			bus.mem[0] = uint8(op)
			bus.mem[1] = 0x00 // immediate/displacement
			bus.mem[2] = 0x00
			bus.mem[3] = 0x00
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic on opcode 0x%02X: %v", op, r)
				}
			}()
			cpu.Step()
		})
	}
}

// TestAllCBOpcodes_NoPanic executes every CB-prefixed opcode.
func TestAllCBOpcodes_NoPanic(t *testing.T) {
	for op := 0; op < 256; op++ {
		t.Run(fmt.Sprintf("CB_%02X", op), func(t *testing.T) {
			bus := &testBus{}
			cpu := New(bus)
			cpu.reg.SP = 0xFFF0
			cpu.reg.HL = 0x8000 // Safe address for (HL) ops
			bus.mem[0] = 0xCB
			bus.mem[1] = uint8(op)
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic on CB 0x%02X: %v", op, r)
				}
			}()
			cpu.Step()
		})
	}
}

// TestAllEDOpcodes_NoPanic executes every ED-prefixed opcode.
func TestAllEDOpcodes_NoPanic(t *testing.T) {
	for op := 0; op < 256; op++ {
		t.Run(fmt.Sprintf("ED_%02X", op), func(t *testing.T) {
			bus := &testBus{}
			cpu := New(bus)
			cpu.reg.SP = 0xFFF0
			cpu.reg.HL = 0x8000
			cpu.reg.BC = 0x0100
			bus.mem[0] = 0xED
			bus.mem[1] = uint8(op)
			bus.mem[2] = 0x00
			bus.mem[3] = 0x00
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic on ED 0x%02X: %v", op, r)
				}
			}()
			cpu.Step()
		})
	}
}

// TestAllDDOpcodes_NoPanic executes every DD-prefixed opcode.
func TestAllDDOpcodes_NoPanic(t *testing.T) {
	for op := 0; op < 256; op++ {
		if op == 0xCB {
			continue // tested separately
		}
		t.Run(fmt.Sprintf("DD_%02X", op), func(t *testing.T) {
			bus := &testBus{}
			cpu := New(bus)
			cpu.reg.SP = 0xFFF0
			cpu.reg.IX = 0x8000
			bus.mem[0] = 0xDD
			bus.mem[1] = uint8(op)
			bus.mem[2] = 0x00
			bus.mem[3] = 0x00
			bus.mem[4] = 0x00
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic on DD 0x%02X: %v", op, r)
				}
			}()
			cpu.Step()
		})
	}
}

// TestAllDDCBOpcodes_NoPanic executes every DD CB d op sequence.
func TestAllDDCBOpcodes_NoPanic(t *testing.T) {
	for op := 0; op < 256; op++ {
		t.Run(fmt.Sprintf("DDCB_%02X", op), func(t *testing.T) {
			bus := &testBus{}
			cpu := New(bus)
			cpu.reg.SP = 0xFFF0
			cpu.reg.IX = 0x8000
			bus.mem[0] = 0xDD
			bus.mem[1] = 0xCB
			bus.mem[2] = 0x00 // displacement
			bus.mem[3] = uint8(op)
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic on DD CB 00 0x%02X: %v", op, r)
				}
			}()
			cpu.Step()
		})
	}
}
