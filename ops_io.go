package z80

func init() {
	// --- IN A, (n) ---
	baseOps[0xDB] = func(c *CPU, _ uint8) {
		port := uint16(c.fetchPC()) | uint16(c.getA())<<8
		c.setA(c.inBus(port))
		c.cycles += 11
	}

	// --- OUT (n), A ---
	baseOps[0xD3] = func(c *CPU, _ uint8) {
		port := uint16(c.fetchPC()) | uint16(c.getA())<<8
		c.outBus(port, c.getA())
		c.cycles += 11
	}

	// --- INI ---
	edOps[0xA2] = func(c *CPU, _ uint8) {
		c.blockIN(1)
		c.cycles += 16
	}

	// --- IND ---
	edOps[0xAA] = func(c *CPU, _ uint8) {
		c.blockIN(-1)
		c.cycles += 16
	}

	// --- INIR ---
	edOps[0xB2] = func(c *CPU, _ uint8) {
		c.blockIN(1)
		c.blockIORepeat(c.getB() != 0)
	}

	// --- INDR ---
	edOps[0xBA] = func(c *CPU, _ uint8) {
		c.blockIN(-1)
		c.blockIORepeat(c.getB() != 0)
	}

	// --- OUTI ---
	edOps[0xA3] = func(c *CPU, _ uint8) {
		c.blockOUT(1)
		c.cycles += 16
	}

	// --- OUTD ---
	edOps[0xAB] = func(c *CPU, _ uint8) {
		c.blockOUT(-1)
		c.cycles += 16
	}

	// --- OTIR ---
	edOps[0xB3] = func(c *CPU, _ uint8) {
		c.blockOUT(1)
		c.blockIORepeat(c.getB() != 0)
	}

	// --- OTDR ---
	edOps[0xBB] = func(c *CPU, _ uint8) {
		c.blockOUT(-1)
		c.blockIORepeat(c.getB() != 0)
	}
}

// blockIORepeat handles the repeat-or-finish logic for block I/O instructions.
func (c *CPU) blockIORepeat(repeat bool) {
	if repeat {
		c.reg.PC -= 2
		c.cycles += 21
	} else {
		c.cycles += 16
	}
}

// blockIN performs the core of INI/IND/INIR/INDR.
func (c *CPU) blockIN(dir int) {
	val := c.inBus(c.reg.BC)
	c.writeBus(c.reg.HL, val)
	b := c.getB() - 1
	c.setB(b)
	if dir > 0 {
		c.reg.HL++
	} else {
		c.reg.HL--
	}
	f := szFlags(b)
	if val&0x80 != 0 {
		f |= flagN
	}
	// Undocumented flag behavior for block I/O
	k := uint16(val) + uint16(uint8(c.getC()+uint8(dir)))
	if k > 255 {
		f |= flagH | flagC
	}
	f |= parityTable[uint8(k&7)^b]
	c.setF(f)
}

// blockOUT performs the core of OUTI/OUTD/OTIR/OTDR.
func (c *CPU) blockOUT(dir int) {
	val := c.readBus(c.reg.HL)
	b := c.getB() - 1
	c.setB(b)
	c.outBus(c.reg.BC, val)
	if dir > 0 {
		c.reg.HL++
	} else {
		c.reg.HL--
	}
	f := szFlags(b)
	if val&0x80 != 0 {
		f |= flagN
	}
	k := uint16(val) + uint16(c.getL())
	if k > 255 {
		f |= flagH | flagC
	}
	f |= parityTable[uint8(k&7)^b]
	c.setF(f)
}
