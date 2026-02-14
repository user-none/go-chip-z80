package z80

// Register accessors for the Z80 instruction set.
// These provide both indexed access (for instructions encoding registers
// in 3-bit fields) and named access for specific-register instructions.

// --- Named accessors ---

func (c *CPU) getA() uint8    { return uint8(c.reg.AF >> 8) }
func (c *CPU) setA(v uint8)   { c.reg.AF = uint16(v)<<8 | c.reg.AF&0xFF }
func (c *CPU) getF() uint8    { return uint8(c.reg.AF) }
func (c *CPU) setF(v uint8)   { c.reg.AF = c.reg.AF&0xFF00 | uint16(v) }
func (c *CPU) getB() uint8    { return uint8(c.reg.BC >> 8) }
func (c *CPU) setB(v uint8)   { c.reg.BC = uint16(v)<<8 | c.reg.BC&0xFF }
func (c *CPU) getC() uint8    { return uint8(c.reg.BC) }
func (c *CPU) setC(v uint8)   { c.reg.BC = c.reg.BC&0xFF00 | uint16(v) }
func (c *CPU) getD() uint8    { return uint8(c.reg.DE >> 8) }
func (c *CPU) setD(v uint8)   { c.reg.DE = uint16(v)<<8 | c.reg.DE&0xFF }
func (c *CPU) getE() uint8    { return uint8(c.reg.DE) }
func (c *CPU) setE(v uint8)   { c.reg.DE = c.reg.DE&0xFF00 | uint16(v) }
func (c *CPU) getH() uint8    { return uint8(*c.ixiyReg >> 8) }
func (c *CPU) setH(v uint8)   { *c.ixiyReg = uint16(v)<<8 | *c.ixiyReg&0xFF }
func (c *CPU) getL() uint8    { return uint8(*c.ixiyReg) }
func (c *CPU) setL(v uint8)   { *c.ixiyReg = *c.ixiyReg&0xFF00 | uint16(v) }

// --- Indexed register access ---
// 3-bit field: 0=B, 1=C, 2=D, 3=E, 4=H, 5=L, 6=(HL), 7=A
// Indices 4,5 route through ixiyReg for DD/FD prefix support.
// Index 6 accesses memory at (HL); DD/FD handlers use separate paths.

func (c *CPU) getR8(idx uint8) uint8 {
	switch idx {
	case 0:
		return c.getB()
	case 1:
		return c.getC()
	case 2:
		return c.getD()
	case 3:
		return c.getE()
	case 4:
		return c.getH()
	case 5:
		return c.getL()
	case 6:
		return c.readBus(c.reg.HL)
	case 7:
		return c.getA()
	}
	return 0
}

func (c *CPU) setR8(idx uint8, v uint8) {
	switch idx {
	case 0:
		c.setB(v)
	case 1:
		c.setC(v)
	case 2:
		c.setD(v)
	case 3:
		c.setE(v)
	case 4:
		c.setH(v)
	case 5:
		c.setL(v)
	case 6:
		c.writeBus(c.reg.HL, v)
	case 7:
		c.setA(v)
	}
}

// getRR returns a pointer to a 16-bit register pair by 2-bit index.
// 0=BC, 1=DE, 2=HL(or IX/IY), 3=SP
func (c *CPU) getRR(idx uint8) *uint16 {
	switch idx {
	case 0:
		return &c.reg.BC
	case 1:
		return &c.reg.DE
	case 2:
		return c.ixiyReg
	case 3:
		return &c.reg.SP
	}
	return &c.reg.HL
}

// getRRPush returns a pointer to a 16-bit register pair for PUSH/POP.
// 0=BC, 1=DE, 2=HL(or IX/IY), 3=AF
func (c *CPU) getRRPush(idx uint8) *uint16 {
	switch idx {
	case 0:
		return &c.reg.BC
	case 1:
		return &c.reg.DE
	case 2:
		return c.ixiyReg
	case 3:
		return &c.reg.AF
	}
	return &c.reg.AF
}

// ixiyAddr fetches a displacement byte from PC and returns *ixiyReg + sign_extend(d).
func (c *CPU) ixiyAddr() uint16 {
	d := int8(c.fetchPC())
	return uint16(int32(*c.ixiyReg) + int32(d))
}
