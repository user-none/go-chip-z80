package z80

func init() {
	// --- JP nn ---
	baseOps[0xC3] = func(c *CPU, _ uint8) {
		c.reg.PC = c.fetchPC16()
		c.cycles += 10
	}

	// --- JP cc, nn ---
	// 0xC2=NZ, 0xCA=Z, 0xD2=NC, 0xDA=C, 0xE2=PO, 0xEA=PE, 0xF2=P, 0xFA=M
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0xC2
		baseOps[op] = func(c *CPU, op uint8) {
			addr := c.fetchPC16()
			if c.testCC((op >> 3) & 7) {
				c.reg.PC = addr
			}
			c.cycles += 10
		}
	}

	// --- JR e ---
	baseOps[0x18] = func(c *CPU, _ uint8) {
		e := int8(c.fetchPC())
		c.reg.PC = uint16(int32(c.reg.PC) + int32(e))
		c.cycles += 12
	}

	// --- JR NZ, e ---
	baseOps[0x20] = func(c *CPU, _ uint8) {
		e := int8(c.fetchPC())
		if c.getF()&flagZ == 0 {
			c.reg.PC = uint16(int32(c.reg.PC) + int32(e))
			c.cycles += 12
		} else {
			c.cycles += 7
		}
	}

	// --- JR Z, e ---
	baseOps[0x28] = func(c *CPU, _ uint8) {
		e := int8(c.fetchPC())
		if c.getF()&flagZ != 0 {
			c.reg.PC = uint16(int32(c.reg.PC) + int32(e))
			c.cycles += 12
		} else {
			c.cycles += 7
		}
	}

	// --- JR NC, e ---
	baseOps[0x30] = func(c *CPU, _ uint8) {
		e := int8(c.fetchPC())
		if c.getF()&flagC == 0 {
			c.reg.PC = uint16(int32(c.reg.PC) + int32(e))
			c.cycles += 12
		} else {
			c.cycles += 7
		}
	}

	// --- JR C, e ---
	baseOps[0x38] = func(c *CPU, _ uint8) {
		e := int8(c.fetchPC())
		if c.getF()&flagC != 0 {
			c.reg.PC = uint16(int32(c.reg.PC) + int32(e))
			c.cycles += 12
		} else {
			c.cycles += 7
		}
	}

	// --- JP (HL) ---
	baseOps[0xE9] = func(c *CPU, _ uint8) {
		c.reg.PC = *c.ixiyReg
		c.cycles += 4
	}

	// --- DJNZ e ---
	baseOps[0x10] = func(c *CPU, _ uint8) {
		e := int8(c.fetchPC())
		b := c.getB() - 1
		c.setB(b)
		if b != 0 {
			c.reg.PC = uint16(int32(c.reg.PC) + int32(e))
			c.cycles += 13
		} else {
			c.cycles += 8
		}
	}

	// --- CALL nn ---
	baseOps[0xCD] = func(c *CPU, _ uint8) {
		addr := c.fetchPC16()
		c.push16(c.reg.PC)
		c.reg.PC = addr
		c.cycles += 17
	}

	// --- CALL cc, nn ---
	// 0xC4=NZ, 0xCC=Z, 0xD4=NC, 0xDC=C, 0xE4=PO, 0xEC=PE, 0xF4=P, 0xFC=M
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0xC4
		baseOps[op] = func(c *CPU, op uint8) {
			addr := c.fetchPC16()
			if c.testCC((op >> 3) & 7) {
				c.push16(c.reg.PC)
				c.reg.PC = addr
				c.cycles += 17
			} else {
				c.cycles += 10
			}
		}
	}

	// --- RET ---
	baseOps[0xC9] = func(c *CPU, _ uint8) {
		c.reg.PC = c.pop16()
		c.cycles += 10
	}

	// --- RET cc ---
	// 0xC0=NZ, 0xC8=Z, 0xD0=NC, 0xD8=C, 0xE0=PO, 0xE8=PE, 0xF0=P, 0xF8=M
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0xC0
		baseOps[op] = func(c *CPU, op uint8) {
			if c.testCC((op >> 3) & 7) {
				c.reg.PC = c.pop16()
				c.cycles += 11
			} else {
				c.cycles += 5
			}
		}
	}

	// --- RST p ---
	// 0xC7=00, 0xCF=08, 0xD7=10, 0xDF=18, 0xE7=20, 0xEF=28, 0xF7=30, 0xFF=38
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0xC7
		baseOps[op] = func(c *CPU, op uint8) {
			c.push16(c.reg.PC)
			c.reg.PC = uint16(op & 0x38)
			c.cycles += 11
		}
	}
}
