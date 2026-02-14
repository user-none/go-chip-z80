package z80

func init() {
	// --- 8-bit ALU: ADD/ADC/SUB/SBC/AND/XOR/OR/CP ---
	// Opcodes 0x80-0xBF: ALU A, r
	for i := 0; i < 64; i++ {
		op := uint8(0x80 + i)
		src := op & 7
		if src == 6 {
			baseOps[op] = func(c *CPU, op uint8) {
				aluOp8(c, (op>>3)&7, c.readBus(c.reg.HL))
				c.cycles += 7
			}
		} else {
			baseOps[op] = func(c *CPU, op uint8) {
				aluOp8(c, (op>>3)&7, c.getR8(op&7))
				c.cycles += 4
			}
		}
	}

	// --- ALU A, n (immediate) ---
	// 0xC6=ADD, 0xCE=ADC, 0xD6=SUB, 0xDE=SBC, 0xE6=AND, 0xEE=XOR, 0xF6=OR, 0xFE=CP
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0xC6
		baseOps[op] = func(c *CPU, op uint8) {
			aluOp8(c, (op>>3)&7, c.fetchPC())
			c.cycles += 7
		}
	}

	// --- INC r ---
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0x04
		if i == 6 {
			// INC (HL): 11 cycles
			baseOps[op] = func(c *CPU, _ uint8) {
				val := c.readBus(c.reg.HL)
				f := incFlags8(val)
				val++
				c.writeBus(c.reg.HL, val)
				c.setF(f | (c.getF() & flagC))
				c.cycles += 11
			}
		} else {
			baseOps[op] = func(c *CPU, op uint8) {
				r := (op >> 3) & 7
				val := c.getR8(r)
				f := incFlags8(val)
				val++
				c.setR8(r, val)
				c.setF(f | (c.getF() & flagC))
				c.cycles += 4
			}
		}
	}

	// --- DEC r ---
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0x05
		if i == 6 {
			// DEC (HL): 11 cycles
			baseOps[op] = func(c *CPU, _ uint8) {
				val := c.readBus(c.reg.HL)
				f := decFlags8(val)
				val--
				c.writeBus(c.reg.HL, val)
				c.setF(f | (c.getF() & flagC))
				c.cycles += 11
			}
		} else {
			baseOps[op] = func(c *CPU, op uint8) {
				r := (op >> 3) & 7
				val := c.getR8(r)
				f := decFlags8(val)
				val--
				c.setR8(r, val)
				c.setF(f | (c.getF() & flagC))
				c.cycles += 4
			}
		}
	}

	// --- INC rr (16-bit) ---
	// 0x03=BC, 0x13=DE, 0x23=HL, 0x33=SP
	for i := uint8(0); i < 4; i++ {
		op := i<<4 | 0x03
		baseOps[op] = func(c *CPU, op uint8) {
			rr := c.getRR((op >> 4) & 3)
			*rr++
			c.cycles += 6
		}
	}

	// --- DEC rr (16-bit) ---
	// 0x0B=BC, 0x1B=DE, 0x2B=HL, 0x3B=SP
	for i := uint8(0); i < 4; i++ {
		op := i<<4 | 0x0B
		baseOps[op] = func(c *CPU, op uint8) {
			rr := c.getRR((op >> 4) & 3)
			*rr--
			c.cycles += 6
		}
	}

	// --- ADD HL, rr ---
	// 0x09=BC, 0x19=DE, 0x29=HL, 0x39=SP
	for i := uint8(0); i < 4; i++ {
		op := i<<4 | 0x09
		baseOps[op] = func(c *CPU, op uint8) {
			rr := c.getRR((op >> 4) & 3)
			hl := *c.ixiyReg
			val := *rr
			result := uint32(hl) + uint32(val)
			f := c.getF() & (flagS | flagZ | flagPV) // preserve S, Z, PV
			if result > 0xFFFF {
				f |= flagC
			}
			if (hl^val^uint16(result))&0x1000 != 0 {
				f |= flagH
			}
			r16 := uint16(result)
			f |= uint8(r16>>8) & (flagF5 | flagF3)
			c.setF(f)
			*c.ixiyReg = r16
			c.cycles += 11
		}
	}

	// --- DAA ---
	baseOps[0x27] = func(c *CPU, _ uint8) {
		a := c.getA()
		f := c.getF()
		correction := uint8(0)

		if f&flagH != 0 || a&0x0F > 9 {
			correction |= 0x06
		}
		if f&flagC != 0 || a > 0x99 {
			correction |= 0x60
		}

		newA := a
		if f&flagN != 0 {
			newA -= correction
		} else {
			newA += correction
		}

		newF := szFlags(newA) | parityTable[newA]
		newF |= f & flagN // preserve N
		if f&flagC != 0 || a > 0x99 {
			newF |= flagC
		}
		if f&flagN != 0 {
			if f&flagH != 0 && a&0x0F < 6 {
				newF |= flagH
			}
		} else {
			if a&0x0F > 9 {
				newF |= flagH
			}
		}

		c.setA(newA)
		c.setF(newF)
		c.cycles += 4
	}

	// --- CPL ---
	baseOps[0x2F] = func(c *CPU, _ uint8) {
		a := c.getA() ^ 0xFF
		c.setA(a)
		f := c.getF() & (flagS | flagZ | flagPV | flagC)
		f |= flagH | flagN
		f |= a & (flagF3 | flagF5)
		c.setF(f)
		c.cycles += 4
	}

	// --- SCF ---
	baseOps[0x37] = func(c *CPU, _ uint8) {
		a := c.getA()
		oldF := c.getF()
		f := oldF & (flagS | flagZ | flagPV)
		f |= flagC
		f |= (a | oldF) & (flagF3 | flagF5)
		c.setF(f)
		c.cycles += 4
	}

	// --- CCF ---
	baseOps[0x3F] = func(c *CPU, _ uint8) {
		a := c.getA()
		oldF := c.getF()
		oldC := oldF & flagC
		f := oldF & (flagS | flagZ | flagPV)
		if oldC != 0 {
			f |= flagH
		} else {
			f |= flagC
		}
		f |= (a | oldF) & (flagF3 | flagF5)
		c.setF(f)
		c.cycles += 4
	}
}

// aluOp8 performs ALU operation on A with operand b.
// op: 0=ADD, 1=ADC, 2=SUB, 3=SBC, 4=AND, 5=XOR, 6=OR, 7=CP
func aluOp8(c *CPU, op, b uint8) {
	a := c.getA()
	carry := uint8(0)
	if c.getF()&flagC != 0 {
		carry = 1
	}

	switch op {
	case 0: // ADD
		f := addFlags8(a, b, 0)
		c.setA(a + b)
		c.setF(f)
	case 1: // ADC
		f := addFlags8(a, b, carry)
		c.setA(a + b + carry)
		c.setF(f)
	case 2: // SUB
		f := subFlags8(a, b, 0)
		c.setA(a - b)
		c.setF(f)
	case 3: // SBC
		f := subFlags8(a, b, carry)
		c.setA(a - b - carry)
		c.setF(f)
	case 4: // AND
		result := a & b
		c.setA(result)
		c.setF(logicFlags(result, true))
	case 5: // XOR
		result := a ^ b
		c.setA(result)
		c.setF(logicFlags(result, false))
	case 6: // OR
		result := a | b
		c.setA(result)
		c.setF(logicFlags(result, false))
	case 7: // CP
		f := subFlags8(a, b, 0)
		// F3/F5 come from the operand b, not the result
		f = (f &^ (flagF3 | flagF5)) | (b & (flagF3 | flagF5))
		c.setF(f)
	}
}
