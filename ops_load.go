package z80

func init() {
	// --- LD r,r' (0x40-0x7F except 0x76 which is HALT) ---
	for i := 0; i < 64; i++ {
		op := uint8(0x40 + i)
		if op == 0x76 {
			continue // HALT
		}
		dst := (op >> 3) & 7
		src := op & 7
		if dst == 6 || src == 6 {
			// (HL) involved: 7 cycles
			baseOps[op] = func(c *CPU, op uint8) {
				d := (op >> 3) & 7
				s := op & 7
				c.setR8(d, c.getR8(s))
				c.cycles += 7
			}
		} else {
			// Register-to-register: 4 cycles
			baseOps[op] = func(c *CPU, op uint8) {
				d := (op >> 3) & 7
				s := op & 7
				c.setR8(d, c.getR8(s))
				c.cycles += 4
			}
		}
	}

	// --- LD r, n (immediate) ---
	for i := uint8(0); i < 8; i++ {
		op := i<<3 | 0x06
		if i == 6 {
			// LD (HL), n: 10 cycles
			baseOps[op] = func(c *CPU, _ uint8) {
				n := c.fetchPC()
				c.writeBus(c.reg.HL, n)
				c.cycles += 10
			}
		} else {
			baseOps[op] = func(c *CPU, op uint8) {
				r := (op >> 3) & 7
				n := c.fetchPC()
				c.setR8(r, n)
				c.cycles += 7
			}
		}
	}

	// --- LD A, (BC) ---
	baseOps[0x0A] = func(c *CPU, _ uint8) {
		c.setA(c.readBus(c.reg.BC))
		c.cycles += 7
	}
	// --- LD A, (DE) ---
	baseOps[0x1A] = func(c *CPU, _ uint8) {
		c.setA(c.readBus(c.reg.DE))
		c.cycles += 7
	}
	// --- LD (BC), A ---
	baseOps[0x02] = func(c *CPU, _ uint8) {
		c.writeBus(c.reg.BC, c.getA())
		c.cycles += 7
	}
	// --- LD (DE), A ---
	baseOps[0x12] = func(c *CPU, _ uint8) {
		c.writeBus(c.reg.DE, c.getA())
		c.cycles += 7
	}
	// --- LD A, (nn) ---
	baseOps[0x3A] = func(c *CPU, _ uint8) {
		addr := c.fetchPC16()
		c.setA(c.readBus(addr))
		c.cycles += 13
	}
	// --- LD (nn), A ---
	baseOps[0x32] = func(c *CPU, _ uint8) {
		addr := c.fetchPC16()
		c.writeBus(addr, c.getA())
		c.cycles += 13
	}

	// --- LD rr, nn (16-bit immediate) ---
	// 0x01=BC, 0x11=DE, 0x21=HL, 0x31=SP
	for i := uint8(0); i < 4; i++ {
		op := i << 4 | 0x01
		baseOps[op] = func(c *CPU, op uint8) {
			rr := c.getRR((op >> 4) & 3)
			*rr = c.fetchPC16()
			c.cycles += 10
		}
	}

	// --- LD (nn), HL ---
	baseOps[0x22] = func(c *CPU, _ uint8) {
		addr := c.fetchPC16()
		c.write16(addr, *c.ixiyReg)
		c.cycles += 16
	}
	// --- LD HL, (nn) ---
	baseOps[0x2A] = func(c *CPU, _ uint8) {
		addr := c.fetchPC16()
		*c.ixiyReg = c.read16(addr)
		c.cycles += 16
	}

	// --- LD SP, HL ---
	baseOps[0xF9] = func(c *CPU, _ uint8) {
		c.reg.SP = *c.ixiyReg
		c.cycles += 6
	}

	// --- PUSH rr ---
	// 0xC5=BC, 0xD5=DE, 0xE5=HL, 0xF5=AF
	for i := uint8(0); i < 4; i++ {
		op := i<<4 | 0xC5
		baseOps[op] = func(c *CPU, op uint8) {
			rr := c.getRRPush((op >> 4) & 3)
			c.push16(*rr)
			c.cycles += 11
		}
	}

	// --- POP rr ---
	// 0xC1=BC, 0xD1=DE, 0xE1=HL, 0xF1=AF
	for i := uint8(0); i < 4; i++ {
		op := i<<4 | 0xC1
		baseOps[op] = func(c *CPU, op uint8) {
			rr := c.getRRPush((op >> 4) & 3)
			*rr = c.pop16()
			c.cycles += 10
		}
	}

	// --- EX DE, HL ---
	baseOps[0xEB] = func(c *CPU, _ uint8) {
		c.reg.DE, c.reg.HL = c.reg.HL, c.reg.DE
		c.cycles += 4
	}

	// --- EX AF, AF' ---
	baseOps[0x08] = func(c *CPU, _ uint8) {
		c.reg.AF, c.reg.AF_ = c.reg.AF_, c.reg.AF
		c.cycles += 4
	}

	// --- EXX ---
	baseOps[0xD9] = func(c *CPU, _ uint8) {
		c.reg.BC, c.reg.BC_ = c.reg.BC_, c.reg.BC
		c.reg.DE, c.reg.DE_ = c.reg.DE_, c.reg.DE
		c.reg.HL, c.reg.HL_ = c.reg.HL_, c.reg.HL
		c.cycles += 4
	}

	// --- EX (SP), HL ---
	baseOps[0xE3] = func(c *CPU, _ uint8) {
		lo := uint16(c.readBus(c.reg.SP))
		hi := uint16(c.readBus(c.reg.SP + 1))
		val := hi<<8 | lo
		c.writeBus(c.reg.SP, uint8(*c.ixiyReg))
		c.writeBus(c.reg.SP+1, uint8(*c.ixiyReg>>8))
		*c.ixiyReg = val
		c.cycles += 19
	}

	// --- HALT ---
	// PC already points past the HALT opcode (incremented by fetchOpcode).
	// Leave it there so that when an interrupt pushes PC, the return
	// address is the instruction AFTER HALT, not HALT itself.
	// The halt state is maintained by Step() returning 4-cycle NOPs.
	baseOps[0x76] = func(c *CPU, _ uint8) {
		c.reg.Halted = true
		c.cycles += 4
	}

	// --- DI ---
	baseOps[0xF3] = func(c *CPU, _ uint8) {
		c.reg.IFF1 = false
		c.reg.IFF2 = false
		c.cycles += 4
	}

	// --- EI ---
	baseOps[0xFB] = func(c *CPU, _ uint8) {
		c.reg.IFF1 = true
		c.reg.IFF2 = true
		c.afterEI = true
		c.cycles += 4
	}
}
