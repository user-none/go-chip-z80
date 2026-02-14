package z80

import (
	"fmt"
	"testing"
)

// z80State describes the full Z80 state for a single-step test case.
type z80State struct {
	A, F, B, C, D, E, H, L uint8
	I, R                    uint8
	PC, SP                  uint16
	IX, IY                  uint16
	AF_, BC_, DE_, HL_      uint16
	IM                      uint8
	IFF1, IFF2              bool
	RAM                     [][2]uint16 // {{addr, val}, ...}
	Ports                   [][2]uint16 // {{port, val}, ...} for input ports
	Cycles                  int         // 0 = don't check
}

// sstBus implements Bus for single-step tests with configurable port reads.
type sstBus struct {
	mem    [65536]uint8
	portIn map[uint16]uint8
}

func (b *sstBus) Fetch(addr uint16) uint8      { return b.mem[addr] }
func (b *sstBus) Read(addr uint16) uint8       { return b.mem[addr] }
func (b *sstBus) Write(addr uint16, val uint8) { b.mem[addr] = val }
func (b *sstBus) In(port uint16) uint8 {
	if v, ok := b.portIn[port]; ok {
		return v
	}
	return 0xFF
}
func (b *sstBus) Out(port uint16, val uint8) {}

// runSSTTest sets up initial state, executes one Step, and checks the result.
func runSSTTest(t *testing.T, tc struct {
	name       string
	init, want z80State
}) {
	t.Helper()
	t.Run(tc.name, func(t *testing.T) {
		bus := &sstBus{portIn: make(map[uint16]uint8)}
		for _, entry := range tc.init.RAM {
			bus.mem[entry[0]] = uint8(entry[1])
		}
		for _, entry := range tc.init.Ports {
			bus.portIn[entry[0]] = uint8(entry[1])
		}

		cpu := New(bus)
		cpu.SetState(Registers{
			AF:   uint16(tc.init.A)<<8 | uint16(tc.init.F),
			BC:   uint16(tc.init.B)<<8 | uint16(tc.init.C),
			DE:   uint16(tc.init.D)<<8 | uint16(tc.init.E),
			HL:   uint16(tc.init.H)<<8 | uint16(tc.init.L),
			AF_:  tc.init.AF_,
			BC_:  tc.init.BC_,
			DE_:  tc.init.DE_,
			HL_:  tc.init.HL_,
			IX:   tc.init.IX,
			IY:   tc.init.IY,
			SP:   tc.init.SP,
			PC:   tc.init.PC,
			I:    tc.init.I,
			R:    tc.init.R,
			IFF1: tc.init.IFF1,
			IFF2: tc.init.IFF2,
			IM:   tc.init.IM,
		})

		cycles := cpu.Step()
		regs := cpu.Registers()

		gotA := uint8(regs.AF >> 8)
		gotF := uint8(regs.AF)
		gotB := uint8(regs.BC >> 8)
		gotC := uint8(regs.BC)
		gotD := uint8(regs.DE >> 8)
		gotE := uint8(regs.DE)
		gotH := uint8(regs.HL >> 8)
		gotL := uint8(regs.HL)

		check := func(name string, got, want uint8) {
			if got != want {
				t.Errorf("%s = 0x%02X, want 0x%02X", name, got, want)
			}
		}
		check16 := func(name string, got, want uint16) {
			if got != want {
				t.Errorf("%s = 0x%04X, want 0x%04X", name, got, want)
			}
		}

		check("A", gotA, tc.want.A)
		if gotF != tc.want.F {
			t.Errorf("F = 0x%02X, want 0x%02X  %s", gotF, tc.want.F, flagDiff(gotF, tc.want.F))
		}
		check("B", gotB, tc.want.B)
		check("C", gotC, tc.want.C)
		check("D", gotD, tc.want.D)
		check("E", gotE, tc.want.E)
		check("H", gotH, tc.want.H)
		check("L", gotL, tc.want.L)
		check("I", regs.I, tc.want.I)
		check("R", regs.R, tc.want.R)
		check16("PC", regs.PC, tc.want.PC)
		check16("SP", regs.SP, tc.want.SP)
		check16("IX", regs.IX, tc.want.IX)
		check16("IY", regs.IY, tc.want.IY)
		check16("AF'", regs.AF_, tc.want.AF_)
		check16("BC'", regs.BC_, tc.want.BC_)
		check16("DE'", regs.DE_, tc.want.DE_)
		check16("HL'", regs.HL_, tc.want.HL_)
		check("IM", regs.IM, tc.want.IM)

		if regs.IFF1 != tc.want.IFF1 {
			t.Errorf("IFF1 = %v, want %v", regs.IFF1, tc.want.IFF1)
		}
		if regs.IFF2 != tc.want.IFF2 {
			t.Errorf("IFF2 = %v, want %v", regs.IFF2, tc.want.IFF2)
		}

		for _, entry := range tc.want.RAM {
			addr, val := entry[0], uint8(entry[1])
			got := bus.mem[addr]
			if got != val {
				t.Errorf("RAM[0x%04X] = 0x%02X, want 0x%02X", addr, got, val)
			}
		}

		if tc.want.Cycles > 0 && cycles != tc.want.Cycles {
			t.Errorf("cycles = %d, want %d", cycles, tc.want.Cycles)
		}
	})
}

// flagDiff returns a human-readable flag-by-flag comparison.
func flagDiff(got, want uint8) string {
	names := [8]string{"C", "N", "PV", "F3", "H", "F5", "Z", "S"}
	var s string
	for i := 0; i < 8; i++ {
		g := (got >> uint(i)) & 1
		w := (want >> uint(i)) & 1
		if g != w {
			if s != "" {
				s += " "
			}
			s += fmt.Sprintf("%s:%d->%d", names[i], w, g)
		}
	}
	return "[" + s + "]"
}
