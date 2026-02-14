package z80

// Flag bit positions in the F register.
const (
	flagC  uint8 = 1 << 0 // Carry
	flagN  uint8 = 1 << 1 // Subtract
	flagPV uint8 = 1 << 2 // Parity/Overflow
	flagF3 uint8 = 1 << 3 // Undocumented (bit 3 of result)
	flagH  uint8 = 1 << 4 // Half-carry
	flagF5 uint8 = 1 << 5 // Undocumented (bit 5 of result)
	flagZ  uint8 = 1 << 6 // Zero
	flagS  uint8 = 1 << 7 // Sign
)

// parityTable[i] is flagPV if i has even parity, 0 otherwise.
var parityTable [256]uint8

func init() {
	for i := 0; i < 256; i++ {
		bits := 0
		v := i
		for v != 0 {
			bits += v & 1
			v >>= 1
		}
		if bits%2 == 0 {
			parityTable[i] = flagPV
		}
	}
}

// szFlags returns S, Z, F5, F3 flags for the given 8-bit result.
func szFlags(val uint8) uint8 {
	f := val & (flagS | flagF5 | flagF3)
	if val == 0 {
		f |= flagZ
	}
	return f
}

// addFlags8 computes all flags for an 8-bit addition: a + b + carry.
// Returns the full flag byte.
func addFlags8(a, b, carry uint8) uint8 {
	result := uint16(a) + uint16(b) + uint16(carry)
	r8 := uint8(result)
	f := szFlags(r8)
	if result > 0xFF {
		f |= flagC
	}
	// Half-carry: carry out of bit 3
	if (a^b^r8)&0x10 != 0 {
		f |= flagH
	}
	// Overflow: both operands same sign, result different sign
	if (a^b)&0x80 == 0 && (a^r8)&0x80 != 0 {
		f |= flagPV
	}
	return f
}

// subFlags8 computes all flags for an 8-bit subtraction: a - b - carry.
// Returns the full flag byte.
func subFlags8(a, b, carry uint8) uint8 {
	result := uint16(a) - uint16(b) - uint16(carry)
	r8 := uint8(result)
	f := szFlags(r8) | flagN
	if result > 0xFF {
		f |= flagC
	}
	// Half-borrow: borrow from bit 4
	if (a^b^r8)&0x10 != 0 {
		f |= flagH
	}
	// Overflow: operands different sign, result sign matches b
	if (a^b)&0x80 != 0 && (a^r8)&0x80 != 0 {
		f |= flagPV
	}
	return f
}

// incFlags8 computes flags for INC: preserves carry flag.
func incFlags8(val uint8) uint8 {
	result := val + 1
	f := szFlags(result)
	if result&0x0F == 0 {
		f |= flagH
	}
	if val == 0x7F {
		f |= flagPV
	}
	return f
}

// decFlags8 computes flags for DEC: preserves carry, sets N.
func decFlags8(val uint8) uint8 {
	result := val - 1
	f := szFlags(result) | flagN
	if result&0x0F == 0x0F {
		f |= flagH
	}
	if val == 0x80 {
		f |= flagPV
	}
	return f
}

// logicFlags computes flags for AND/OR/XOR. setH is true for AND.
func logicFlags(result uint8, setH bool) uint8 {
	f := szFlags(result) | parityTable[result]
	if setH {
		f |= flagH
	}
	return f
}

// testCC evaluates a 3-bit condition code against the current flags.
// cc: 0=NZ, 1=Z, 2=NC, 3=C, 4=PO, 5=PE, 6=P, 7=M
func (c *CPU) testCC(cc uint8) bool {
	f := uint8(c.reg.AF)
	switch cc {
	case 0:
		return f&flagZ == 0
	case 1:
		return f&flagZ != 0
	case 2:
		return f&flagC == 0
	case 3:
		return f&flagC != 0
	case 4:
		return f&flagPV == 0
	case 5:
		return f&flagPV != 0
	case 6:
		return f&flagS == 0
	case 7:
		return f&flagS != 0
	}
	return false
}
