// Package z80 implements a Zilog Z80 CPU emulator.
//
// The Z80 is an 8-bit CPU with a 16-bit address bus and separate I/O space:
//   - Seven 8-bit general-purpose registers (A, B, C, D, E, H, L) plus flags (F)
//   - A complete shadow register set (AF', BC', DE', HL')
//   - Two 16-bit index registers (IX, IY)
//   - A 16-bit stack pointer (SP) and program counter (PC)
//   - An 8-bit interrupt vector register (I) and refresh counter (R)
//   - Three interrupt modes (IM 0, IM 1, IM 2) and two flip-flops (IFF1, IFF2)
package z80

// Bus provides memory and I/O access for the Z80 CPU.
//
// The Z80 has separate address spaces for memory (64KB, accessed via MREQ)
// and I/O ports (accessed via IORQ). Read/Write handle memory; In/Out
// handle I/O.
type Bus interface {
	// Fetch reads an opcode byte during an M1 (opcode fetch) machine cycle.
	// On real hardware the M1 signal is asserted during this access,
	// which some systems use for wait-state insertion, memory contention
	// timing, or bank switching. For systems that don't distinguish M1
	// from data reads, this can simply delegate to Read.
	Fetch(addr uint16) uint8

	// Read reads a byte from the given memory address (non-M1 data read).
	Read(addr uint16) uint8

	// Write writes a byte to the given memory address.
	Write(addr uint16, val uint8)

	// In reads a byte from the given I/O port.
	// The full 16-bit address bus is provided: the low byte is the port
	// number specified by the instruction, and the high byte is context
	// dependent (register A for single-byte IN/OUT, register B or C for
	// block I/O instructions).
	In(port uint16) uint8

	// Out writes a byte to the given I/O port.
	Out(port uint16, val uint8)
}
