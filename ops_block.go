package z80

func init() {
	// --- LDI ---
	edOps[0xA0] = func(c *CPU, _ uint8) {
		c.blockLD(1)
		c.cycles += 16
	}

	// --- LDD ---
	edOps[0xA8] = func(c *CPU, _ uint8) {
		c.blockLD(-1)
		c.cycles += 16
	}

	// --- LDIR ---
	edOps[0xB0] = func(c *CPU, _ uint8) {
		c.blockLD(1)
		c.blockRepeat(c.reg.BC != 0)
	}

	// --- LDDR ---
	edOps[0xB8] = func(c *CPU, _ uint8) {
		c.blockLD(-1)
		c.blockRepeat(c.reg.BC != 0)
	}

	// --- CPI ---
	edOps[0xA1] = func(c *CPU, _ uint8) {
		c.blockCP(1)
		c.cycles += 16
	}

	// --- CPD ---
	edOps[0xA9] = func(c *CPU, _ uint8) {
		c.blockCP(-1)
		c.cycles += 16
	}

	// --- CPIR ---
	edOps[0xB1] = func(c *CPU, _ uint8) {
		c.blockCP(1)
		c.blockRepeat(c.reg.BC != 0 && c.getF()&flagZ == 0)
	}

	// --- CPDR ---
	edOps[0xB9] = func(c *CPU, _ uint8) {
		c.blockCP(-1)
		c.blockRepeat(c.reg.BC != 0 && c.getF()&flagZ == 0)
	}
}

// blockLD performs the core of LDI/LDD/LDIR/LDDR.
func (c *CPU) blockLD(dir int) {
	val := c.readBus(c.reg.HL)
	c.writeBus(c.reg.DE, val)
	if dir > 0 {
		c.reg.HL++
		c.reg.DE++
	} else {
		c.reg.HL--
		c.reg.DE--
	}
	c.reg.BC--
	n := val + c.getA()
	f := c.getF() & (flagS | flagZ | flagC)
	if n&0x02 != 0 {
		f |= flagF5
	}
	f |= n & flagF3
	if c.reg.BC != 0 {
		f |= flagPV
	}
	c.setF(f)
}

// blockRepeat handles the repeat-or-finish logic for block instructions.
// If repeat is true, PC is rewound and extra cycles are charged.
func (c *CPU) blockRepeat(repeat bool) {
	if repeat {
		c.reg.PC -= 2
		c.blockRepeatF35()
		c.cycles += 21
	} else {
		c.cycles += 16
	}
}

// blockRepeatF35 replaces F3/F5 with high byte of PC+1 (WZ) for repeat block ops.
// Must be called after PC -= 2.
func (c *CPU) blockRepeatF35() {
	wzHi := uint8((c.reg.PC + 1) >> 8)
	f := c.getF() &^ (flagF3 | flagF5)
	f |= wzHi & (flagF3 | flagF5)
	c.setF(f)
}

// blockCP performs the core of CPI/CPD/CPIR/CPDR.
func (c *CPU) blockCP(dir int) {
	val := c.readBus(c.reg.HL)
	a := c.getA()
	result := a - val
	if dir > 0 {
		c.reg.HL++
	} else {
		c.reg.HL--
	}
	c.reg.BC--
	f := szFlags(result) | flagN | (c.getF() & flagC)
	if (a^val^result)&0x10 != 0 {
		f |= flagH
	}
	// F3/F5 from result - H (clear szFlags F3/F5 first)
	f &^= flagF3 | flagF5
	n := result
	if f&flagH != 0 {
		n--
	}
	if n&0x02 != 0 {
		f |= flagF5
	}
	f |= n & flagF3
	if c.reg.BC != 0 {
		f |= flagPV
	}
	c.setF(f)
}
