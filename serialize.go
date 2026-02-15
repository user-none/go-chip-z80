package z80

import (
	"encoding/binary"
	"errors"
)

const cpuSerializeVersion = 1

// SerializeSize is the number of bytes needed to serialize the CPU state.
const SerializeSize = 47

// Serialize writes the complete CPU state into buf in a compact little-endian
// binary format. Returns an error if len(buf) < SerializeSize. Bus
// references are not included — the caller handles memory and I/O state
// separately.
func (c *CPU) Serialize(buf []byte) error {
	if len(buf) < SerializeSize {
		return errors.New("z80: serialize buffer too small")
	}

	buf[0] = cpuSerializeVersion
	binary.LittleEndian.PutUint16(buf[1:], c.reg.AF)
	binary.LittleEndian.PutUint16(buf[3:], c.reg.BC)
	binary.LittleEndian.PutUint16(buf[5:], c.reg.DE)
	binary.LittleEndian.PutUint16(buf[7:], c.reg.HL)
	binary.LittleEndian.PutUint16(buf[9:], c.reg.AF_)
	binary.LittleEndian.PutUint16(buf[11:], c.reg.BC_)
	binary.LittleEndian.PutUint16(buf[13:], c.reg.DE_)
	binary.LittleEndian.PutUint16(buf[15:], c.reg.HL_)
	binary.LittleEndian.PutUint16(buf[17:], c.reg.IX)
	binary.LittleEndian.PutUint16(buf[19:], c.reg.IY)
	binary.LittleEndian.PutUint16(buf[21:], c.reg.SP)
	binary.LittleEndian.PutUint16(buf[23:], c.reg.PC)
	buf[25] = c.reg.I
	buf[26] = c.reg.R
	buf[27] = boolByte(c.reg.IFF1)
	buf[28] = boolByte(c.reg.IFF2)
	buf[29] = c.reg.IM
	buf[30] = boolByte(c.reg.Halted)
	binary.LittleEndian.PutUint64(buf[31:], c.cycles)
	binary.LittleEndian.PutUint32(buf[39:], uint32(int32(c.deficit)))
	buf[43] = boolByte(c.intLine)
	buf[44] = c.intData
	buf[45] = boolByte(c.nmiPending)
	buf[46] = boolByte(c.afterEI)
	return nil
}

// Deserialize restores the complete CPU state from buf, which must have been
// produced by Serialize. Returns an error if the buffer is too small or was
// produced by an incompatible version. Bus references are not modified — the
// caller handles memory and I/O state separately.
func (c *CPU) Deserialize(buf []byte) error {
	if len(buf) < SerializeSize {
		return errors.New("z80: deserialize buffer too small")
	}
	if buf[0] != cpuSerializeVersion {
		return errors.New("z80: unsupported serialize version")
	}

	c.reg.AF = binary.LittleEndian.Uint16(buf[1:])
	c.reg.BC = binary.LittleEndian.Uint16(buf[3:])
	c.reg.DE = binary.LittleEndian.Uint16(buf[5:])
	c.reg.HL = binary.LittleEndian.Uint16(buf[7:])
	c.reg.AF_ = binary.LittleEndian.Uint16(buf[9:])
	c.reg.BC_ = binary.LittleEndian.Uint16(buf[11:])
	c.reg.DE_ = binary.LittleEndian.Uint16(buf[13:])
	c.reg.HL_ = binary.LittleEndian.Uint16(buf[15:])
	c.reg.IX = binary.LittleEndian.Uint16(buf[17:])
	c.reg.IY = binary.LittleEndian.Uint16(buf[19:])
	c.reg.SP = binary.LittleEndian.Uint16(buf[21:])
	c.reg.PC = binary.LittleEndian.Uint16(buf[23:])
	c.reg.I = buf[25]
	c.reg.R = buf[26]
	c.reg.IFF1 = buf[27] != 0
	c.reg.IFF2 = buf[28] != 0
	c.reg.IM = buf[29]
	c.reg.Halted = buf[30] != 0
	c.cycles = binary.LittleEndian.Uint64(buf[31:])
	c.deficit = int(int32(binary.LittleEndian.Uint32(buf[39:])))
	c.intLine = buf[43] != 0
	c.intData = buf[44]
	c.nmiPending = buf[45] != 0
	c.afterEI = buf[46] != 0

	c.ixiyReg = &c.reg.HL
	return nil
}

func boolByte(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
