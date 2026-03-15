package z80

import "testing"

func BenchmarkStep_PlainBus(b *testing.B) {
	bus := &testBus{}
	cpu := New(bus)
	b.ResetTimer()
	for b.Loop() {
		cpu.Step()
	}
}
