package z80

// opFunc is the handler signature for all Z80 instructions.
// The opcode byte is passed so handlers can decode register fields.
type opFunc func(c *CPU, op uint8)

// Opcode dispatch tables.
var (
	baseOps [256]opFunc // Unprefixed opcodes
	cbOps   [256]opFunc // CB prefix (bit/rotate)
	edOps   [256]opFunc // ED prefix (extended)
	ixOps   [256]opFunc // DD/FD shared (IX/IY via pointer)
	ixcbOps [256]opFunc // DD CB / FD CB (indexed bit ops)
)

// execute fetches and runs the instruction at PC.
func (c *CPU) execute() {
	op := c.fetchOpcode()
	baseOps[op](c, op)
}

func init() {
	// NOP
	baseOps[0x00] = func(c *CPU, op uint8) {
		c.cycles += 4
	}

	// CB prefix
	baseOps[0xCB] = prefixCB

	// DD prefix (IX)
	baseOps[0xDD] = prefixDD

	// FD prefix (IY)
	baseOps[0xFD] = prefixFD

	// ED prefix
	baseOps[0xED] = prefixED
}

// prefixCB handles the CB prefix: fetch second opcode, dispatch through cbOps.
func prefixCB(c *CPU, _ uint8) {
	op := c.fetchOpcode()
	if h := cbOps[op]; h != nil {
		h(c, op)
	} else {
		c.cycles += 8 // undocumented NOP
	}
}

// prefixDD handles the DD prefix (IX register).
func prefixDD(c *CPU, _ uint8) { prefixIXIY(c, &c.reg.IX) }

// prefixFD handles the FD prefix (IY register).
func prefixFD(c *CPU, _ uint8) { prefixIXIY(c, &c.reg.IY) }

// prefixIXIY is the shared logic for DD (IX) and FD (IY) prefixes.
func prefixIXIY(c *CPU, reg *uint16) {
	prev := c.ixiyReg
	c.ixiyReg = reg
	op := c.fetchOpcode()
	if op == 0xCB {
		c.idxAddr = c.ixiyAddr()
		op2 := c.fetchPC()
		c.cycles += 4 // extra prefix timing
		if h := ixcbOps[op2]; h != nil {
			h(c, op2)
		} else {
			c.cycles += 8
		}
	} else if h := ixOps[op]; h != nil {
		h(c, op)
	} else if h := baseOps[op]; h != nil {
		c.cycles += 4 // prefix cost
		h(c, op)
	} else {
		c.cycles += 4
	}
	c.ixiyReg = prev
}

// prefixED handles the ED prefix: fetch second opcode, dispatch through edOps.
func prefixED(c *CPU, _ uint8) {
	op := c.fetchOpcode()
	if h := edOps[op]; h != nil {
		h(c, op)
	} else {
		c.cycles += 8 // ED + NOP
	}
}
