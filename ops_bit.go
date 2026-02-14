package z80

func init() {
	// --- CB prefix: BIT b, r (0x40-0x7F) ---
	for i := 0; i < 64; i++ {
		op := uint8(0x40 + i)
		src := op & 7
		if src == 6 {
			// BIT b, (HL): 12 cycles
			cbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				val := c.readBus(c.reg.HL)
				f := c.getF()&flagC | flagH
				if val&(1<<bit) == 0 {
					f |= flagZ | flagPV
				}
				if bit == 7 && val&0x80 != 0 {
					f |= flagS
				}
				// F3/F5 from high byte of address for (HL) variant
				f |= uint8(c.reg.HL>>8) & (flagF3 | flagF5)
				c.setF(f)
				c.cycles += 12
			}
		} else {
			// BIT b, r: 8 cycles
			cbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				val := c.getR8(op & 7)
				f := c.getF()&flagC | flagH
				if val&(1<<bit) == 0 {
					f |= flagZ | flagPV
				}
				if bit == 7 && val&0x80 != 0 {
					f |= flagS
				}
				f |= val & (flagF3 | flagF5)
				c.setF(f)
				c.cycles += 8
			}
		}
	}

	// --- CB prefix: RES b, r (0x80-0xBF) ---
	for i := 0; i < 64; i++ {
		op := uint8(0x80 + i)
		src := op & 7
		if src == 6 {
			cbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				val := c.readBus(c.reg.HL)
				val &^= 1 << bit
				c.writeBus(c.reg.HL, val)
				c.cycles += 15
			}
		} else {
			cbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				s := op & 7
				val := c.getR8(s)
				val &^= 1 << bit
				c.setR8(s, val)
				c.cycles += 8
			}
		}
	}

	// --- CB prefix: SET b, r (0xC0-0xFF) ---
	for i := 0; i < 64; i++ {
		op := uint8(0xC0 + i)
		src := op & 7
		if src == 6 {
			cbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				val := c.readBus(c.reg.HL)
				val |= 1 << bit
				c.writeBus(c.reg.HL, val)
				c.cycles += 15
			}
		} else {
			cbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				s := op & 7
				val := c.getR8(s)
				val |= 1 << bit
				c.setR8(s, val)
				c.cycles += 8
			}
		}
	}
}
