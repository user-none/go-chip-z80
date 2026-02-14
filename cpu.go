package z80

// CPU is the Z80 processor.
type CPU struct {
	reg    Registers
	bus    Bus
	cbus   CycledBus // non-nil when bus implements CycledBus
	cycles uint64

	// Interrupt state.
	intLine    bool  // INT line level (active when true)
	intData    uint8 // Data bus value for interrupt acknowledge
	nmiPending bool  // NMI edge latch (consumed on next Step)
	afterEI    bool  // Suppress interrupts for one instruction after EI

	// Cycle deficit from StepCycles when an instruction's cost
	// exceeded the budget.
	deficit int

	// DD/FD prefix support: points to HL, IX, or IY.
	ixiyReg *uint16
	// Pre-computed indexed address for DD CB / FD CB instructions.
	idxAddr uint16
}

// New creates a CPU wired to the given bus and performs a reset.
// If bus implements [CycledBus], the cycle-aware methods are used
// for all memory and I/O access.
func New(bus Bus) *CPU {
	c := &CPU{bus: bus}
	if cb, ok := bus.(CycledBus); ok {
		c.cbus = cb
	}
	c.ixiyReg = &c.reg.HL
	c.Reset()
	return c
}

// Reset reinitializes the CPU to its power-on state:
// PC=0, SP=0xFFFF, AF=0xFFFF, interrupts disabled, IM 0, clears HALT.
// The total cycle counter is reset to 0. Bus state is not affected.
func (c *CPU) Reset() {
	c.reg = Registers{
		AF: 0xFFFF,
		SP: 0xFFFF,
	}
	c.cycles = 0
	c.deficit = 0
	c.intLine = false
	c.intData = 0xFF
	c.nmiPending = false
	c.afterEI = false
	c.ixiyReg = &c.reg.HL
}

// Step executes a single instruction and returns the T-states consumed.
//
// Processing order each call:
//  1. If an NMI is latched, service it (11 T-states).
//  2. If INT is asserted, IFF1 is set, and not suppressed by EI delay,
//     service the maskable interrupt (cycles depend on IM).
//  3. If halted, return 4 T-states (internal NOP).
//  4. Otherwise fetch and execute the next instruction.
func (c *CPU) Step() int {
	before := c.cycles

	// 1. NMI has highest priority.
	if c.nmiPending {
		c.nmiPending = false
		c.serviceNMI()
		return int(c.cycles - before)
	}

	// 2. Maskable interrupt (subject to IFF1 and EI delay).
	if c.intLine && c.reg.IFF1 && !c.afterEI {
		c.serviceINT()
		return int(c.cycles - before)
	}
	c.afterEI = false

	// 3. HALT burns NOP cycles.
	if c.reg.Halted {
		c.cycles += 4
		return int(c.cycles - before)
	}

	// 4. Fetch and execute.
	c.execute()

	return int(c.cycles - before)
}

// StepCycles executes a single instruction within the given cycle budget.
// If a previous instruction's cost exceeded its budget, the deficit is
// paid down first without executing a new instruction. When an instruction's
// cost exceeds the budget, the excess is stored as a deficit to be charged
// on subsequent calls. Returns the number of cycles consumed from the budget.
func (c *CPU) StepCycles(budget int) int {
	// Pay down deficit from a previous instruction that exceeded its budget.
	if c.deficit > 0 {
		if budget >= c.deficit {
			n := c.deficit
			c.deficit = 0
			return n
		}
		c.deficit -= budget
		return budget
	}

	cost := c.Step()

	if cost <= budget {
		return cost
	}

	c.deficit = cost - budget
	return budget
}

// Deficit returns the remaining cycle debt from a previous StepCycles
// call where the instruction cost exceeded the budget.
func (c *CPU) Deficit() int {
	return c.deficit
}

// Cycles returns the total T-state count since the last Reset.
func (c *CPU) Cycles() uint64 {
	return c.cycles
}

// Halted returns true if the CPU is in HALT state, waiting for an interrupt.
func (c *CPU) Halted() bool {
	return c.reg.Halted
}

// Registers returns a snapshot of the current register state.
func (c *CPU) Registers() Registers {
	return c.reg
}

// SetState sets all registers directly without performing a reset.
// Intended for testing and state serialization/deserialization.
func (c *CPU) SetState(regs Registers) {
	c.reg = regs
}

// fetchOpcode reads the byte at PC via an M1 (opcode fetch) bus cycle
// and advances PC by 1. Increments the R register (low 7 bits only).
func (c *CPU) fetchOpcode() uint8 {
	val := c.fetchBus(c.reg.PC)
	c.reg.PC++
	c.reg.R = (c.reg.R & 0x80) | ((c.reg.R + 1) & 0x7F)
	return val
}

// --- Bus dispatch helpers ---
// These route through CycledBus when available, otherwise plain Bus.

func (c *CPU) fetchBus(addr uint16) uint8 {
	if c.cbus != nil {
		return c.cbus.CycledFetch(c.cycles, addr)
	}
	return c.bus.Fetch(addr)
}

func (c *CPU) readBus(addr uint16) uint8 {
	if c.cbus != nil {
		return c.cbus.CycledRead(c.cycles, addr)
	}
	return c.bus.Read(addr)
}

func (c *CPU) writeBus(addr uint16, val uint8) {
	if c.cbus != nil {
		c.cbus.CycledWrite(c.cycles, addr, val)
		return
	}
	c.bus.Write(addr, val)
}

func (c *CPU) inBus(port uint16) uint8 {
	if c.cbus != nil {
		return c.cbus.CycledIn(c.cycles, port)
	}
	return c.bus.In(port)
}

func (c *CPU) outBus(port uint16, val uint8) {
	if c.cbus != nil {
		c.cbus.CycledOut(c.cycles, port, val)
		return
	}
	c.bus.Out(port, val)
}

// --- Memory access helpers ---

// fetchPC reads the byte at PC and advances PC by 1.
func (c *CPU) fetchPC() uint8 {
	val := c.readBus(c.reg.PC)
	c.reg.PC++
	return val
}

// fetchPC16 reads a little-endian 16-bit word at PC and advances PC by 2.
func (c *CPU) fetchPC16() uint16 {
	lo := uint16(c.fetchPC())
	hi := uint16(c.fetchPC())
	return hi<<8 | lo
}

// read16 reads a little-endian 16-bit word from addr.
func (c *CPU) read16(addr uint16) uint16 {
	lo := uint16(c.readBus(addr))
	hi := uint16(c.readBus(addr + 1))
	return hi<<8 | lo
}

// write16 writes a little-endian 16-bit word to addr.
func (c *CPU) write16(addr uint16, val uint16) {
	c.writeBus(addr, uint8(val))
	c.writeBus(addr+1, uint8(val>>8))
}

// push16 pushes a 16-bit value onto the stack.
func (c *CPU) push16(val uint16) {
	c.reg.SP--
	c.writeBus(c.reg.SP, uint8(val>>8))
	c.reg.SP--
	c.writeBus(c.reg.SP, uint8(val))
}

// pop16 pops a 16-bit value from the stack.
func (c *CPU) pop16() uint16 {
	lo := uint16(c.readBus(c.reg.SP))
	c.reg.SP++
	hi := uint16(c.readBus(c.reg.SP))
	c.reg.SP++
	return hi<<8 | lo
}
