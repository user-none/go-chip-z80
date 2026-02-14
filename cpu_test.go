package z80

import "testing"

// testBus is a simple Bus implementation for testing.
type testBus struct {
	mem [65536]uint8
}

func (b *testBus) Fetch(addr uint16) uint8        { return b.mem[addr] }
func (b *testBus) Read(addr uint16) uint8         { return b.mem[addr] }
func (b *testBus) Write(addr uint16, val uint8)   { b.mem[addr] = val }
func (b *testBus) In(port uint16) uint8           { return 0xFF }
func (b *testBus) Out(port uint16, val uint8)     {}

func newTestCPU() (*CPU, *testBus) {
	bus := &testBus{}
	cpu := New(bus)
	return cpu, bus
}

// testCycleBus implements CycleBus and records cycle values.
type testCycleBus struct {
	testBus
	lastFetchCycle uint64
	lastReadCycle  uint64
	lastWriteCycle uint64
	lastInCycle    uint64
	lastOutCycle   uint64
}

func (b *testCycleBus) CycleFetch(cycle uint64, addr uint16) uint8 {
	b.lastFetchCycle = cycle
	return b.mem[addr]
}
func (b *testCycleBus) CycleRead(cycle uint64, addr uint16) uint8 {
	b.lastReadCycle = cycle
	return b.mem[addr]
}
func (b *testCycleBus) CycleWrite(cycle uint64, addr uint16, val uint8) {
	b.lastWriteCycle = cycle
	b.mem[addr] = val
}
func (b *testCycleBus) CycleIn(cycle uint64, port uint16) uint8 {
	b.lastInCycle = cycle
	return 0xFF
}
func (b *testCycleBus) CycleOut(cycle uint64, port uint16, val uint8) {
	b.lastOutCycle = cycle
}

func newTestCycleCPU() (*CPU, *testCycleBus) {
	bus := &testCycleBus{}
	cpu := New(bus)
	return cpu, bus
}

func TestNew(t *testing.T) {
	cpu, _ := newTestCPU()
	regs := cpu.Registers()

	if regs.PC != 0 {
		t.Errorf("PC = %04x, want 0000", regs.PC)
	}
	if regs.SP != 0xFFFF {
		t.Errorf("SP = %04x, want FFFF", regs.SP)
	}
	if regs.AF != 0xFFFF {
		t.Errorf("AF = %04x, want FFFF", regs.AF)
	}
	if regs.IFF1 || regs.IFF2 {
		t.Error("interrupts should be disabled after reset")
	}
	if regs.IM != 0 {
		t.Errorf("IM = %d, want 0", regs.IM)
	}
	if regs.Halted {
		t.Error("should not be halted after reset")
	}
	if cpu.Cycles() != 0 {
		t.Errorf("Cycles = %d, want 0", cpu.Cycles())
	}
}

func TestReset(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.reg.PC = 0x1234
	cpu.reg.IFF1 = true
	cpu.cycles = 999

	cpu.Reset()

	if cpu.reg.PC != 0 {
		t.Errorf("PC = %04x after reset, want 0000", cpu.reg.PC)
	}
	if cpu.Cycles() != 0 {
		t.Errorf("Cycles = %d after reset, want 0", cpu.Cycles())
	}
}

func TestStep_ReturnsPositiveCycles(t *testing.T) {
	cpu, _ := newTestCPU()
	// Memory is zero-filled (NOP = 0x00 on Z80)
	cycles := cpu.Step()
	if cycles <= 0 {
		t.Errorf("Step returned %d cycles, want > 0", cycles)
	}
}

func TestStep_AdvancesPC(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.Step()
	if cpu.reg.PC == 0 {
		t.Error("PC should advance after Step")
	}
}

func TestCycles_Accumulates(t *testing.T) {
	cpu, _ := newTestCPU()
	c1 := cpu.Step()
	c2 := cpu.Step()
	if cpu.Cycles() != uint64(c1+c2) {
		t.Errorf("Cycles = %d, want %d", cpu.Cycles(), c1+c2)
	}
}

func TestHalted(t *testing.T) {
	cpu, _ := newTestCPU()
	if cpu.Halted() {
		t.Error("should not be halted initially")
	}
	cpu.reg.Halted = true
	if !cpu.Halted() {
		t.Error("Halted() should return true")
	}
}

func TestSetState(t *testing.T) {
	cpu, _ := newTestCPU()
	regs := Registers{
		AF: 0x1234,
		BC: 0x5678,
		PC: 0xABCD,
		SP: 0xFFFE,
		IM: 1,
	}
	cpu.SetState(regs)
	got := cpu.Registers()
	if got.AF != 0x1234 || got.BC != 0x5678 || got.PC != 0xABCD {
		t.Errorf("SetState/Registers mismatch: %+v", got)
	}
}

func TestStepCycles_WithinBudget(t *testing.T) {
	cpu, _ := newTestCPU()
	consumed := cpu.StepCycles(100)
	if consumed <= 0 || consumed > 100 {
		t.Errorf("StepCycles consumed %d, want 1..100", consumed)
	}
	if cpu.Deficit() != 0 {
		t.Errorf("Deficit = %d, want 0", cpu.Deficit())
	}
}

func TestStepCycles_ExceedsBudget(t *testing.T) {
	cpu, _ := newTestCPU()
	// Budget of 1 cycle; any instruction will exceed it.
	consumed := cpu.StepCycles(1)
	if consumed != 1 {
		t.Errorf("StepCycles consumed %d, want 1", consumed)
	}
	if cpu.Deficit() <= 0 {
		t.Error("expected positive deficit")
	}
}

func TestStepCycles_DeficitPaydown(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.StepCycles(1) // Create deficit
	deficit := cpu.Deficit()

	// Next call should pay down deficit without executing.
	pcBefore := cpu.reg.PC
	consumed := cpu.StepCycles(deficit + 10)
	if consumed != deficit {
		t.Errorf("deficit paydown consumed %d, want %d", consumed, deficit)
	}
	if cpu.reg.PC != pcBefore {
		t.Error("PC should not advance during deficit paydown")
	}
	if cpu.Deficit() != 0 {
		t.Errorf("Deficit = %d after paydown, want 0", cpu.Deficit())
	}
}

func TestNMI(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.PC = 0x0100
	cpu.reg.SP = 0xFFFE

	// Write a NOP at 0x0066 so the CPU has something to execute there.
	bus.mem[0x0066] = 0x00

	cpu.NMI()
	cycles := cpu.Step()

	if cycles != 11 {
		t.Errorf("NMI cycles = %d, want 11", cycles)
	}
	if cpu.reg.PC != 0x0066 {
		t.Errorf("PC = %04x after NMI, want 0066", cpu.reg.PC)
	}
	if cpu.reg.IFF1 {
		t.Error("IFF1 should be false after NMI")
	}
	// Check return address was pushed.
	retAddr := cpu.read16(cpu.reg.SP)
	if retAddr != 0x0100 {
		t.Errorf("return address on stack = %04x, want 0100", retAddr)
	}
}

func TestNMI_WakesHalt(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.reg.Halted = true
	cpu.reg.SP = 0xFFFE

	cpu.NMI()
	cpu.Step()

	if cpu.Halted() {
		t.Error("NMI should wake from HALT")
	}
}

func TestINT_IM1(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.reg.PC = 0x0200
	cpu.reg.SP = 0xFFFE
	cpu.reg.IFF1 = true
	cpu.reg.IFF2 = true
	cpu.reg.IM = 1

	cpu.INT(true, 0xFF)
	cycles := cpu.Step()

	if cycles != 13 {
		t.Errorf("IM1 interrupt cycles = %d, want 13", cycles)
	}
	if cpu.reg.PC != 0x0038 {
		t.Errorf("PC = %04x after IM1 INT, want 0038", cpu.reg.PC)
	}
	if cpu.reg.IFF1 || cpu.reg.IFF2 {
		t.Error("interrupts should be disabled after INT")
	}
	retAddr := cpu.read16(cpu.reg.SP)
	if retAddr != 0x0200 {
		t.Errorf("return address on stack = %04x, want 0200", retAddr)
	}
}

func TestINT_NotServicedWhenDisabled(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.reg.IFF1 = false

	cpu.INT(true, 0xFF)
	cpu.Step()

	// Should have executed NOP at 0x0000, not serviced interrupt.
	if cpu.reg.PC == 0x0038 {
		t.Error("INT should not be serviced when IFF1 is false")
	}
}

func TestINT_Deassert(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.reg.IFF1 = true
	cpu.reg.IM = 1

	cpu.INT(true, 0xFF)
	cpu.INT(false, 0) // deassert before Step
	cpu.Step()

	if cpu.reg.PC == 0x0038 {
		t.Error("INT should not be serviced after deassertion")
	}
}

func TestINT_IM2(t *testing.T) {
	cpu, bus := newTestCPU()
	cpu.reg.PC = 0x0300
	cpu.reg.SP = 0xFFFE
	cpu.reg.IFF1 = true
	cpu.reg.IM = 2
	cpu.reg.I = 0x80

	// Vector table at 0x80FE: handler at 0x1234
	bus.mem[0x80FE] = 0x34 // low byte
	bus.mem[0x80FF] = 0x12 // high byte

	cpu.INT(true, 0xFE)
	cycles := cpu.Step()

	if cycles != 19 {
		t.Errorf("IM2 interrupt cycles = %d, want 19", cycles)
	}
	if cpu.reg.PC != 0x1234 {
		t.Errorf("PC = %04x after IM2 INT, want 1234", cpu.reg.PC)
	}
}

func TestINT_IM0_RST(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.reg.PC = 0x0400
	cpu.reg.SP = 0xFFFE
	cpu.reg.IFF1 = true
	cpu.reg.IM = 0

	// RST 38h = 0xFF
	cpu.INT(true, 0xFF)
	cycles := cpu.Step()

	if cycles != 11 {
		t.Errorf("IM0 RST cycles = %d, want 11", cycles)
	}
	if cpu.reg.PC != 0x0038 {
		t.Errorf("PC = %04x after IM0 RST 38h, want 0038", cpu.reg.PC)
	}
}

func TestCycleBus_Detected(t *testing.T) {
	cpu, _ := newTestCycleCPU()
	if cpu.cbus == nil {
		t.Fatal("CycleBus should be detected")
	}
}

func TestCycleBus_FetchPassesCycles(t *testing.T) {
	cpu, bus := newTestCycleCPU()
	// Step will fetchOpcode which calls fetchBus -> CycledFetch at cycle 0.
	cpu.Step()
	if bus.lastFetchCycle != 0 {
		t.Errorf("lastFetchCycle = %d, want 0", bus.lastFetchCycle)
	}
}

func TestCycleBus_WritePassesCycles(t *testing.T) {
	cpu, bus := newTestCycleCPU()
	cpu.reg.SP = 0xFFFE
	cpu.reg.IFF1 = true
	cpu.reg.IM = 1
	cpu.INT(true, 0xFF)
	cpu.Step() // Services IM1 interrupt, which pushes PC via writeBus

	if bus.lastWriteCycle == 0 && cpu.Cycles() == 0 {
		t.Error("expected CycledWrite to be called")
	}
}

func TestPlainBus_NoCycleBus(t *testing.T) {
	cpu, _ := newTestCPU()
	if cpu.cbus != nil {
		t.Fatal("plain Bus should not set cbus")
	}
}

func TestNMI_PriorityOverINT(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.reg.PC = 0x0500
	cpu.reg.SP = 0xFFFE
	cpu.reg.IFF1 = true
	cpu.reg.IM = 1

	// Both pending: NMI should win.
	cpu.INT(true, 0xFF)
	cpu.NMI()
	cpu.Step()

	if cpu.reg.PC != 0x0066 {
		t.Errorf("PC = %04x, want 0066 (NMI should take priority)", cpu.reg.PC)
	}
}
