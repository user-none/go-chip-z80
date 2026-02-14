package z80

// INT asserts or deasserts the maskable interrupt line (active low on
// real hardware, active high here for clarity).
//
// This is level-triggered: while asserted, the CPU checks it at the end
// of each instruction. If IFF1 is set and the EI delay has passed, the
// interrupt is serviced.
//
// The data parameter is the byte placed on the data bus during the
// interrupt acknowledge cycle:
//   - IM 0: executed as an instruction (typically RST n, e.g. 0xFF for RST 38h)
//   - IM 1: ignored (always jumps to 0x0038)
//   - IM 2: combined with I register to form a vector table address (I<<8 | data)
func (c *CPU) INT(assert bool, data uint8) {
	c.intLine = assert
	c.intData = data
}

// NMI triggers a non-maskable interrupt (edge-triggered).
// The NMI is latched and processed at the start of the next Step() call.
// Multiple calls before the next Step() have no additional effect.
func (c *CPU) NMI() {
	c.nmiPending = true
}

// serviceNMI processes a non-maskable interrupt.
//
// The NMI response:
//  1. Exits HALT state if active.
//  2. Copies IFF1 to IFF2 (so RETN can restore it).
//  3. Clears IFF1 (disables maskable interrupts during NMI handler).
//  4. Pushes PC onto the stack.
//  5. Jumps to 0x0066.
//  6. Costs 11 T-states.
func (c *CPU) serviceNMI() {
	c.reg.Halted = false
	c.reg.IFF2 = c.reg.IFF1
	c.reg.IFF1 = false
	c.push16(c.reg.PC)
	c.reg.PC = 0x0066
	c.cycles += 11
}

// serviceINT processes a maskable interrupt based on the current IM.
//
// All modes:
//  1. Exit HALT state if active.
//  2. Disable interrupts (IFF1=false, IFF2=false).
//
// Mode-specific behavior:
//   - IM 0: Execute the instruction on the data bus. Typically RST n (11 T-states).
//   - IM 1: Push PC, jump to 0x0038 (13 T-states).
//   - IM 2: Push PC, read vector from (I<<8 | data), jump to that address (19 T-states).
func (c *CPU) serviceINT() {
	c.reg.Halted = false
	c.reg.IFF1 = false
	c.reg.IFF2 = false
	c.afterEI = false

	switch c.reg.IM {
	case 0:
		c.serviceIM0()
	case 1:
		c.serviceIM1()
	case 2:
		c.serviceIM2()
	default:
		// Invalid IM treated as IM 0.
		c.serviceIM0()
	}
}

// serviceIM0 handles IM 0: execute the data bus value as an instruction.
// Typically the device places an RST instruction (single-byte CALL to a
// fixed address). Only RST instructions are supported for now.
func (c *CPU) serviceIM0() {
	// RST instructions are 11xx_x111 in binary.
	// The restart address is bits 5-3 * 8.
	if c.intData&0xC7 == 0xC7 {
		addr := uint16(c.intData & 0x38)
		c.push16(c.reg.PC)
		c.reg.PC = addr
		c.cycles += 11
		return
	}

	// Fallback: treat as IM 1 for unsupported data bus values.
	c.serviceIM1()
}

// serviceIM1 handles IM 1: push PC, jump to 0x0038.
func (c *CPU) serviceIM1() {
	c.push16(c.reg.PC)
	c.reg.PC = 0x0038
	c.cycles += 13
}

// serviceIM2 handles IM 2: push PC, read vector from table, jump to vector.
// The vector table address is formed by (I << 8) | data, and the target
// address is the 16-bit value read from that location.
func (c *CPU) serviceIM2() {
	c.push16(c.reg.PC)
	tableAddr := uint16(c.reg.I)<<8 | uint16(c.intData)
	c.reg.PC = c.read16(tableAddr)
	c.cycles += 19
}
