package z80

func init() {
	// --- ixOps: instructions that need explicit entries because ---
	// --- (HL) becomes (IX+d)/(IY+d) with different timing ---

	// LD r, (IX+d) — opcodes 0x46,0x4E,0x56,0x5E,0x66,0x6E,0x7E
	for i := uint8(0); i < 8; i++ {
		if i == 6 {
			continue // LD (IX+d),(IX+d) doesn't exist
		}
		op := i<<3 | 0x46
		ixOps[op] = func(c *CPU, op uint8) {
			addr := c.ixiyAddr()
			r := (op >> 3) & 7
			// For DD/FD prefix, indices 4,5 target true H,L not IXH/IXL
			switch r {
			case 4:
				c.reg.HL = uint16(c.readBus(addr))<<8 | c.reg.HL&0xFF
			case 5:
				c.reg.HL = c.reg.HL&0xFF00 | uint16(c.readBus(addr))
			default:
				c.setR8(r, c.readBus(addr))
			}
			c.cycles += 19
		}
	}

	// LD (IX+d), r — opcodes 0x70-0x77 (except 0x76=HALT)
	for i := uint8(0); i < 8; i++ {
		if i == 6 {
			continue // 0x76 is HALT, handled by baseOps
		}
		op := uint8(0x70 + i)
		ixOps[op] = func(c *CPU, op uint8) {
			addr := c.ixiyAddr()
			s := op & 7
			var val uint8
			switch s {
			case 4:
				val = uint8(c.reg.HL >> 8)
			case 5:
				val = uint8(c.reg.HL)
			default:
				val = c.getR8(s)
			}
			c.writeBus(addr, val)
			c.cycles += 19
		}
	}

	// LD (IX+d), n
	ixOps[0x36] = func(c *CPU, _ uint8) {
		addr := c.ixiyAddr()
		n := c.fetchPC()
		c.writeBus(addr, n)
		c.cycles += 19
	}

	// INC (IX+d)
	ixOps[0x34] = func(c *CPU, _ uint8) {
		addr := c.ixiyAddr()
		val := c.readBus(addr)
		f := incFlags8(val)
		val++
		c.writeBus(addr, val)
		c.setF(f | (c.getF() & flagC))
		c.cycles += 23
	}

	// DEC (IX+d)
	ixOps[0x35] = func(c *CPU, _ uint8) {
		addr := c.ixiyAddr()
		val := c.readBus(addr)
		f := decFlags8(val)
		val--
		c.writeBus(addr, val)
		c.setF(f | (c.getF() & flagC))
		c.cycles += 23
	}

	// ALU A, (IX+d) — ADD/ADC/SUB/SBC/AND/XOR/OR/CP
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0x86
		ixOps[op] = func(c *CPU, op uint8) {
			addr := c.ixiyAddr()
			val := c.readBus(addr)
			aluOp8(c, (op>>3)&7, val)
			c.cycles += 19
		}
	}

	// EX (SP), IX/IY — already handled via ixiyReg fallthrough in baseOps
	// JP (IX/IY) — already handled via ixiyReg fallthrough in baseOps
	// LD SP, IX/IY — already handled via ixiyReg fallthrough in baseOps
	// ADD IX,rr — already handled via ixiyReg fallthrough in baseOps
	// PUSH/POP IX/IY — already handled via ixiyReg fallthrough in baseOps
	// LD IX/IY,nn — already handled via ixiyReg fallthrough in baseOps
	// LD (nn),IX/IY — already handled via ixiyReg fallthrough in baseOps
	// LD IX/IY,(nn) — already handled via ixiyReg fallthrough in baseOps
	// INC/DEC IX/IY — already handled via ixiyReg fallthrough in baseOps

	// --- ixcbOps: DD CB d op / FD CB d op ---
	// These use pre-computed c.idxAddr for the indexed address.

	// Rotate/shift (IX+d): 0x00-0x3F
	for i := 0; i < 64; i++ {
		op := uint8(i)
		dst := op & 7
		if dst == 6 {
			// Normal: result stored back to (IX+d)
			ixcbOps[op] = func(c *CPU, op uint8) {
				val := c.readBus(c.idxAddr)
				result, f := cbRotShift((op>>3)&7, val, c.getF())
				c.writeBus(c.idxAddr, result)
				c.setF(f)
				c.cycles += 19
			}
		} else {
			// Undocumented: result also copied to register
			ixcbOps[op] = func(c *CPU, op uint8) {
				val := c.readBus(c.idxAddr)
				result, f := cbRotShift((op>>3)&7, val, c.getF())
				c.writeBus(c.idxAddr, result)
				c.setR8Idx(op&7, result)
				c.setF(f)
				c.cycles += 19
			}
		}
	}

	// BIT b, (IX+d): 0x40-0x7F
	for i := 0; i < 64; i++ {
		op := uint8(0x40 + i)
		ixcbOps[op] = func(c *CPU, op uint8) {
			bit := (op >> 3) & 7
			val := c.readBus(c.idxAddr)
			f := c.getF()&flagC | flagH
			if val&(1<<bit) == 0 {
				f |= flagZ | flagPV
			}
			if bit == 7 && val&0x80 != 0 {
				f |= flagS
			}
			// F3/F5 from high byte of address
			f |= uint8(c.idxAddr>>8) & (flagF3 | flagF5)
			c.setF(f)
			c.cycles += 16
		}
	}

	// RES b, (IX+d): 0x80-0xBF
	for i := 0; i < 64; i++ {
		op := uint8(0x80 + i)
		dst := op & 7
		if dst == 6 {
			ixcbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				val := c.readBus(c.idxAddr)
				val &^= 1 << bit
				c.writeBus(c.idxAddr, val)
				c.cycles += 19
			}
		} else {
			// Undocumented: also stores result in register
			ixcbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				val := c.readBus(c.idxAddr)
				val &^= 1 << bit
				c.writeBus(c.idxAddr, val)
				c.setR8Idx(op&7, val)
				c.cycles += 19
			}
		}
	}

	// SET b, (IX+d): 0xC0-0xFF
	for i := 0; i < 64; i++ {
		op := uint8(0xC0 + i)
		dst := op & 7
		if dst == 6 {
			ixcbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				val := c.readBus(c.idxAddr)
				val |= 1 << bit
				c.writeBus(c.idxAddr, val)
				c.cycles += 19
			}
		} else {
			// Undocumented: also stores result in register
			ixcbOps[op] = func(c *CPU, op uint8) {
				bit := (op >> 3) & 7
				val := c.readBus(c.idxAddr)
				val |= 1 << bit
				c.writeBus(c.idxAddr, val)
				c.setR8Idx(op&7, val)
				c.cycles += 19
			}
		}
	}
}

// setR8Idx sets a register by index without going through ixiyReg.
// Used by undocumented DD CB/FD CB instructions that store results
// into both (IX+d) and a register.
func (c *CPU) setR8Idx(idx uint8, v uint8) {
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
		c.reg.HL = uint16(v)<<8 | c.reg.HL&0xFF
	case 5:
		c.reg.HL = c.reg.HL&0xFF00 | uint16(v)
	case 7:
		c.setA(v)
	}
}
