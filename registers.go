package z80

// Registers holds the programmer-visible state of the Z80.
//
// Register pairs are stored as uint16 with the high byte first:
// AF has A in bits 15-8 and F in bits 7-0. Individual registers
// can be extracted with shifts and masks (e.g. A = AF >> 8).
type Registers struct {
	AF, BC, DE, HL     uint16 // Main register pairs
	AF_, BC_, DE_, HL_ uint16 // Shadow register pairs
	IX, IY             uint16 // Index registers
	SP, PC             uint16 // Stack pointer, program counter
	I, R               uint8  // Interrupt vector, refresh counter
	IFF1, IFF2         bool   // Interrupt flip-flops
	IM                 uint8  // Interrupt mode (0, 1, or 2)
	Halted             bool   // True if executing HALT instruction
}
