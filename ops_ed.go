package z80

func init() {
	// --- IM 0/1/2 ---
	edOps[0x46] = func(c *CPU, _ uint8) { c.reg.IM = 0; c.cycles += 8 }
	edOps[0x56] = func(c *CPU, _ uint8) { c.reg.IM = 1; c.cycles += 8 }
	edOps[0x5E] = func(c *CPU, _ uint8) { c.reg.IM = 2; c.cycles += 8 }
	// Undocumented IM mirrors
	edOps[0x4E] = edOps[0x46]
	edOps[0x66] = edOps[0x46]
	edOps[0x6E] = edOps[0x46]
	edOps[0x76] = edOps[0x56]
	edOps[0x7E] = edOps[0x5E]

	// --- RETI / RETN (identical behavior in emulation) ---
	retnHandler := opFunc(func(c *CPU, _ uint8) {
		c.reg.PC = c.pop16()
		c.reg.IFF1 = c.reg.IFF2
		c.cycles += 14
	})
	edOps[0x45] = retnHandler // RETN
	edOps[0x4D] = retnHandler // RETI
	// Undocumented mirrors
	edOps[0x55] = retnHandler
	edOps[0x5D] = retnHandler
	edOps[0x65] = retnHandler
	edOps[0x6D] = retnHandler
	edOps[0x75] = retnHandler
	edOps[0x7D] = retnHandler

	// --- LD I, A ---
	edOps[0x47] = func(c *CPU, _ uint8) {
		c.reg.I = c.getA()
		c.cycles += 9
	}

	// --- LD R, A ---
	edOps[0x4F] = func(c *CPU, _ uint8) {
		c.reg.R = c.getA()
		c.cycles += 9
	}

	// --- LD A, I ---
	edOps[0x57] = func(c *CPU, _ uint8) { c.ldAIR(c.reg.I) }

	// --- LD A, R ---
	edOps[0x5F] = func(c *CPU, _ uint8) { c.ldAIR(c.reg.R) }

	// --- LD (nn), rr (ED) ---
	// 0x43=BC, 0x53=DE, 0x63=HL, 0x73=SP
	for i := uint8(0); i < 4; i++ {
		op := i<<4 | 0x43
		edOps[op] = func(c *CPU, op uint8) {
			addr := c.fetchPC16()
			rr := c.getRR((op >> 4) & 3)
			c.write16(addr, *rr)
			c.cycles += 20
		}
	}

	// --- LD rr, (nn) (ED) ---
	// 0x4B=BC, 0x5B=DE, 0x6B=HL, 0x7B=SP
	for i := uint8(0); i < 4; i++ {
		op := i<<4 | 0x4B
		edOps[op] = func(c *CPU, op uint8) {
			addr := c.fetchPC16()
			rr := c.getRR((op >> 4) & 3)
			*rr = c.read16(addr)
			c.cycles += 20
		}
	}

	// --- NEG ---
	edOps[0x44] = negHandler
	// Undocumented NEG mirrors
	edOps[0x4C] = negHandler
	edOps[0x54] = negHandler
	edOps[0x5C] = negHandler
	edOps[0x64] = negHandler
	edOps[0x6C] = negHandler
	edOps[0x74] = negHandler
	edOps[0x7C] = negHandler

	// --- ADC HL, rr ---
	// 0x4A=BC, 0x5A=DE, 0x6A=HL, 0x7A=SP
	for i := uint8(0); i < 4; i++ {
		op := i<<4 | 0x4A
		edOps[op] = func(c *CPU, op uint8) {
			rr := c.getRR((op >> 4) & 3)
			hl := c.reg.HL
			val := *rr
			carry := uint32(0)
			if c.getF()&flagC != 0 {
				carry = 1
			}
			result := uint32(hl) + uint32(val) + carry
			r16 := uint16(result)
			f := uint8(r16>>8) & (flagS | flagF5 | flagF3)
			if r16 == 0 {
				f |= flagZ
			}
			if result > 0xFFFF {
				f |= flagC
			}
			if (hl^val^r16)&0x1000 != 0 {
				f |= flagH
			}
			// Overflow: both same sign, result different
			if (hl^val)&0x8000 == 0 && (hl^r16)&0x8000 != 0 {
				f |= flagPV
			}
			c.setF(f)
			c.reg.HL = r16
			c.cycles += 15
		}
	}

	// --- SBC HL, rr ---
	// 0x42=BC, 0x52=DE, 0x62=HL, 0x72=SP
	for i := uint8(0); i < 4; i++ {
		op := i<<4 | 0x42
		edOps[op] = func(c *CPU, op uint8) {
			rr := c.getRR((op >> 4) & 3)
			hl := c.reg.HL
			val := *rr
			carry := uint32(0)
			if c.getF()&flagC != 0 {
				carry = 1
			}
			result := uint32(hl) - uint32(val) - carry
			r16 := uint16(result)
			f := uint8(r16>>8)&(flagS|flagF5|flagF3) | flagN
			if r16 == 0 {
				f |= flagZ
			}
			if result > 0xFFFF {
				f |= flagC
			}
			if (hl^val^r16)&0x1000 != 0 {
				f |= flagH
			}
			// Overflow: operands different sign, result sign matches val
			if (hl^val)&0x8000 != 0 && (hl^r16)&0x8000 != 0 {
				f |= flagPV
			}
			c.setF(f)
			c.reg.HL = r16
			c.cycles += 15
		}
	}

	// --- RLD ---
	edOps[0x6F] = func(c *CPU, _ uint8) {
		a := c.getA()
		val := c.readBus(c.reg.HL)
		// (HL) = (HL low nibble << 4) | (A low nibble)
		// A = (A high nibble) | (HL high nibble)
		newVal := (val << 4) | (a & 0x0F)
		newA := (a & 0xF0) | (val >> 4)
		c.writeBus(c.reg.HL, newVal)
		c.setA(newA)
		f := szFlags(newA) | parityTable[newA] | (c.getF() & flagC)
		c.setF(f)
		c.cycles += 18
	}

	// --- RRD ---
	edOps[0x67] = func(c *CPU, _ uint8) {
		a := c.getA()
		val := c.readBus(c.reg.HL)
		// (HL) = (A low nibble << 4) | (HL high nibble >> 4... no:)
		// RRD: low nibble of (HL) -> low nibble of A
		//      low nibble of A -> high nibble of (HL)
		//      high nibble of (HL) -> low nibble of (HL)
		newVal := (a << 4) | (val >> 4)
		newA := (a & 0xF0) | (val & 0x0F)
		c.writeBus(c.reg.HL, newVal)
		c.setA(newA)
		f := szFlags(newA) | parityTable[newA] | (c.getF() & flagC)
		c.setF(f)
		c.cycles += 18
	}

	// --- IN r, (C) ---
	// 0x40=B, 0x48=C, 0x50=D, 0x58=E, 0x60=H, 0x68=L, 0x78=A
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0x40
		if i == 6 {
			// IN (C) - undocumented: reads port, sets flags, discards result
			edOps[op] = func(c *CPU, _ uint8) {
				val := c.inBus(c.reg.BC)
				f := szFlags(val) | parityTable[val] | (c.getF() & flagC)
				c.setF(f)
				c.cycles += 12
			}
		} else {
			edOps[op] = func(c *CPU, op uint8) {
				r := (op >> 3) & 7
				val := c.inBus(c.reg.BC)
				c.setR8(r, val)
				f := szFlags(val) | parityTable[val] | (c.getF() & flagC)
				c.setF(f)
				c.cycles += 12
			}
		}
	}

	// --- OUT (C), r ---
	// 0x41=B, 0x49=C, 0x51=D, 0x59=E, 0x61=H, 0x69=L, 0x79=A
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0x41
		if i == 6 {
			// OUT (C), 0 - undocumented
			edOps[op] = func(c *CPU, _ uint8) {
				c.outBus(c.reg.BC, 0)
				c.cycles += 12
			}
		} else {
			edOps[op] = func(c *CPU, op uint8) {
				r := (op >> 3) & 7
				c.outBus(c.reg.BC, c.getR8(r))
				c.cycles += 12
			}
		}
	}
}

// ldAIR implements LD A,I and LD A,R: load val into A, set flags.
func (c *CPU) ldAIR(val uint8) {
	c.setA(val)
	f := szFlags(val)
	if c.reg.IFF2 {
		f |= flagPV
	}
	f |= c.getF() & flagC
	c.setF(f)
	c.cycles += 9
}

func negHandler(c *CPU, _ uint8) {
	a := c.getA()
	f := subFlags8(0, a, 0)
	c.setA(0 - a)
	c.setF(f)
	c.cycles += 8
}
