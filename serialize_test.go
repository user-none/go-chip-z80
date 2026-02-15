package z80

import "testing"

func TestSerializeSize(t *testing.T) {
	if SerializeSize != 47 {
		t.Errorf("SerializeSize = %d, want 47", SerializeSize)
	}
}

func TestSerializeRoundTrip(t *testing.T) {
	cpu, _ := newTestCPU()

	// Set all fields to distinct non-zero/non-default values.
	cpu.reg = Registers{
		AF:     0x1234,
		BC:     0x5678,
		DE:     0x9ABC,
		HL:     0xDEF0,
		AF_:    0x1111,
		BC_:    0x2222,
		DE_:    0x3333,
		HL_:    0x4444,
		IX:     0x5555,
		IY:     0x6666,
		SP:     0x7777,
		PC:     0x8888,
		I:      0x42,
		R:      0x73,
		IFF1:   true,
		IFF2:   true,
		IM:     2,
		Halted: true,
	}
	cpu.cycles = 123456789
	cpu.deficit = 7
	cpu.intLine = true
	cpu.intData = 0xCF
	cpu.nmiPending = true
	cpu.afterEI = true

	buf := make([]byte, SerializeSize)
	if err := cpu.Serialize(buf); err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	cpu2, _ := newTestCPU()
	if err := cpu2.Deserialize(buf); err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	// Compare registers.
	r1 := cpu.reg
	r2 := cpu2.reg
	if r1 != r2 {
		t.Errorf("registers mismatch:\n got %+v\nwant %+v", r2, r1)
	}

	// Compare internal state.
	if cpu2.cycles != cpu.cycles {
		t.Errorf("cycles = %d, want %d", cpu2.cycles, cpu.cycles)
	}
	if cpu2.deficit != cpu.deficit {
		t.Errorf("deficit = %d, want %d", cpu2.deficit, cpu.deficit)
	}
	if cpu2.intLine != cpu.intLine {
		t.Errorf("intLine = %v, want %v", cpu2.intLine, cpu.intLine)
	}
	if cpu2.intData != cpu.intData {
		t.Errorf("intData = %02x, want %02x", cpu2.intData, cpu.intData)
	}
	if cpu2.nmiPending != cpu.nmiPending {
		t.Errorf("nmiPending = %v, want %v", cpu2.nmiPending, cpu.nmiPending)
	}
	if cpu2.afterEI != cpu.afterEI {
		t.Errorf("afterEI = %v, want %v", cpu2.afterEI, cpu.afterEI)
	}

	// Verify ixiyReg is reset to HL.
	if cpu2.ixiyReg != &cpu2.reg.HL {
		t.Error("ixiyReg should point to HL after Deserialize")
	}
}

func TestSerializeRoundTripZero(t *testing.T) {
	cpu, _ := newTestCPU()
	cpu.Reset()

	buf := make([]byte, SerializeSize)
	if err := cpu.Serialize(buf); err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	cpu2, _ := newTestCPU()
	// Put cpu2 into a non-default state so we can verify Deserialize resets it.
	cpu2.reg.PC = 0xBEEF
	cpu2.cycles = 999
	if err := cpu2.Deserialize(buf); err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	r1 := cpu.reg
	r2 := cpu2.reg
	if r1 != r2 {
		t.Errorf("registers mismatch:\n got %+v\nwant %+v", r2, r1)
	}
	if cpu2.cycles != cpu.cycles {
		t.Errorf("cycles = %d, want %d", cpu2.cycles, cpu.cycles)
	}
	if cpu2.deficit != cpu.deficit {
		t.Errorf("deficit = %d, want %d", cpu2.deficit, cpu.deficit)
	}
}

func TestSerializePreservesExecution(t *testing.T) {
	cpu1, bus1 := newTestCPU()
	// LD A, 0x42 (3E 42) then INC A (3C) then NOP (00)
	bus1.mem[0] = 0x3E
	bus1.mem[1] = 0x42
	bus1.mem[2] = 0x3C
	bus1.mem[3] = 0x00

	// Run 2 instructions: LD A,0x42 and INC A.
	cpu1.Step()
	cpu1.Step()

	// Serialize after 2 instructions.
	buf := make([]byte, SerializeSize)
	if err := cpu1.Serialize(buf); err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	// Set up cpu2 with the same memory.
	cpu2, bus2 := newTestCPU()
	copy(bus2.mem[:], bus1.mem[:])
	if err := cpu2.Deserialize(buf); err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	// Run 1 more instruction on both (NOP at PC=2).
	cpu1.Step()
	cpu2.Step()

	r1 := cpu1.reg
	r2 := cpu2.reg
	if r1 != r2 {
		t.Errorf("registers diverged after resumed execution:\n got %+v\nwant %+v", r2, r1)
	}
	if cpu1.cycles != cpu2.cycles {
		t.Errorf("cycles diverged: %d vs %d", cpu2.cycles, cpu1.cycles)
	}
}

func TestSerializeErrorShortBuffer(t *testing.T) {
	cpu, _ := newTestCPU()
	buf := make([]byte, 46)

	if err := cpu.Serialize(buf); err == nil {
		t.Error("Serialize should return error with short buffer")
	}
}

func TestDeserializeErrorShortBuffer(t *testing.T) {
	cpu, _ := newTestCPU()
	buf := make([]byte, 46)

	if err := cpu.Deserialize(buf); err == nil {
		t.Error("Deserialize should return error with short buffer")
	}
}

func TestDeserializeErrorBadVersion(t *testing.T) {
	cpu, _ := newTestCPU()
	buf := make([]byte, SerializeSize)
	if err := cpu.Serialize(buf); err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	// Corrupt the version byte.
	buf[0] = 0xFF

	cpu2, _ := newTestCPU()
	if err := cpu2.Deserialize(buf); err == nil {
		t.Error("Deserialize should return error for unsupported version")
	}
}

func BenchmarkSerialize(b *testing.B) {
	bus := &testBus{}
	cpu := New(bus)
	buf := make([]byte, SerializeSize)
	b.ResetTimer()
	for b.Loop() {
		cpu.Serialize(buf)
	}
}

func BenchmarkDeserialize(b *testing.B) {
	bus := &testBus{}
	cpu := New(bus)
	buf := make([]byte, SerializeSize)
	cpu.Serialize(buf)
	cpu2 := New(bus)
	b.ResetTimer()
	for b.Loop() {
		cpu2.Deserialize(buf)
	}
}
