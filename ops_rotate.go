package z80

func init() {
	// --- Base rotates (accumulator) ---

	// RLCA
	baseOps[0x07] = func(c *CPU, _ uint8) {
		a := c.getA()
		bit7 := a >> 7
		a = a<<1 | bit7
		c.setA(a)
		f := c.getF() & (flagS | flagZ | flagPV)
		f |= a & (flagF3 | flagF5)
		if bit7 != 0 {
			f |= flagC
		}
		c.setF(f)
		c.cycles += 4
	}

	// RRCA
	baseOps[0x0F] = func(c *CPU, _ uint8) {
		a := c.getA()
		bit0 := a & 1
		a = a>>1 | bit0<<7
		c.setA(a)
		f := c.getF() & (flagS | flagZ | flagPV)
		f |= a & (flagF3 | flagF5)
		if bit0 != 0 {
			f |= flagC
		}
		c.setF(f)
		c.cycles += 4
	}

	// RLA
	baseOps[0x17] = func(c *CPU, _ uint8) {
		a := c.getA()
		oldC := (c.getF() & flagC)
		bit7 := a >> 7
		a = a << 1
		if oldC != 0 {
			a |= 1
		}
		c.setA(a)
		f := c.getF() & (flagS | flagZ | flagPV)
		f |= a & (flagF3 | flagF5)
		if bit7 != 0 {
			f |= flagC
		}
		c.setF(f)
		c.cycles += 4
	}

	// RRA
	baseOps[0x1F] = func(c *CPU, _ uint8) {
		a := c.getA()
		oldC := (c.getF() & flagC)
		bit0 := a & 1
		a = a >> 1
		if oldC != 0 {
			a |= 0x80
		}
		c.setA(a)
		f := c.getF() & (flagS | flagZ | flagPV)
		f |= a & (flagF3 | flagF5)
		if bit0 != 0 {
			f |= flagC
		}
		c.setF(f)
		c.cycles += 4
	}

	// --- CB prefix: rotate/shift operations (0x00-0x3F) ---
	for i := 0; i < 64; i++ {
		op := uint8(i)
		src := op & 7
		if src == 6 {
			cbOps[op] = func(c *CPU, op uint8) {
				val := c.readBus(c.reg.HL)
				result, f := cbRotShift((op>>3)&7, val, c.getF())
				c.writeBus(c.reg.HL, result)
				c.setF(f)
				c.cycles += 15
			}
		} else {
			cbOps[op] = func(c *CPU, op uint8) {
				s := op & 7
				val := c.getR8(s)
				result, f := cbRotShift((op>>3)&7, val, c.getF())
				c.setR8(s, result)
				c.setF(f)
				c.cycles += 8
			}
		}
	}
}

// cbRotShift performs one of 8 rotate/shift operations.
// rot: 0=RLC, 1=RRC, 2=RL, 3=RR, 4=SLA, 5=SRA, 6=SLL, 7=SRL
func cbRotShift(rot, val, oldF uint8) (result uint8, f uint8) {
	switch rot {
	case 0: // RLC
		bit7 := val >> 7
		result = val<<1 | bit7
		f = szFlags(result) | parityTable[result]
		if bit7 != 0 {
			f |= flagC
		}
	case 1: // RRC
		bit0 := val & 1
		result = val>>1 | bit0<<7
		f = szFlags(result) | parityTable[result]
		if bit0 != 0 {
			f |= flagC
		}
	case 2: // RL
		bit7 := val >> 7
		result = val << 1
		if oldF&flagC != 0 {
			result |= 1
		}
		f = szFlags(result) | parityTable[result]
		if bit7 != 0 {
			f |= flagC
		}
	case 3: // RR
		bit0 := val & 1
		result = val >> 1
		if oldF&flagC != 0 {
			result |= 0x80
		}
		f = szFlags(result) | parityTable[result]
		if bit0 != 0 {
			f |= flagC
		}
	case 4: // SLA
		bit7 := val >> 7
		result = val << 1
		f = szFlags(result) | parityTable[result]
		if bit7 != 0 {
			f |= flagC
		}
	case 5: // SRA
		bit0 := val & 1
		result = (val >> 1) | (val & 0x80) // preserve sign bit
		f = szFlags(result) | parityTable[result]
		if bit0 != 0 {
			f |= flagC
		}
	case 6: // SLL (undocumented: shift left, bit 0 set to 1)
		bit7 := val >> 7
		result = val<<1 | 1
		f = szFlags(result) | parityTable[result]
		if bit7 != 0 {
			f |= flagC
		}
	case 7: // SRL
		bit0 := val & 1
		result = val >> 1
		f = szFlags(result) | parityTable[result]
		if bit0 != 0 {
			f |= flagC
		}
	}
	return
}
