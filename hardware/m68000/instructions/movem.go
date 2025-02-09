package instructions

import "fmt"

type movem struct {
	transferSize     instSize
	mode             addressMode
	reg              byte
	registerListMask uint16
}

func CreateMOVEM(value uint16, code []uint8, offset int) *movem {
	size := oneBitSize(value, 6)
	mode, reg := effective(value, 0, size)

	return &movem{
		transferSize:     size,
		mode:             mode,
		reg:              reg,
		registerListMask: uint16(readData(code, offset, wordSize)),
	}
}

func (m *movem) Name() string {
	return fmt.Sprintf(
		"movem.%s %s, %s",
		m.transferSize.string(),
		m.registers(),
		addressString(m.mode, m.reg, 0, 0))
}

func (m *movem) Length() int {
	return 2 + wordSize.bytesUsed()
}

func (m *movem) registers() string {
	if m.mode == addressRegisterIndirectWithPreDecrement {
		// D then A
		aFirstReg, aLastReg := m.firstAndLastReg(byte(m.registerListMask))
		dFirstReg, dLastReg := m.firstAndLastReg(byte(m.registerListMask >> 8))
		return fmt.Sprintf("D%d-D%d/A%d-A%d", 7-dFirstReg, 7-dLastReg, 7-aFirstReg, 7-aLastReg)
	}

	if m.mode == addressRegisterIndirectWithPostIncrement {
		// A then D
		aLastReg, aFirstReg := m.firstAndLastReg(byte(m.registerListMask >> 8))
		dLastReg, dFirstReg := m.firstAndLastReg(byte(m.registerListMask))
		return fmt.Sprintf("A%d-A%d/D%d-D%d", aFirstReg, aLastReg, dFirstReg, dLastReg)
	}

	panic(fmt.Sprintf("Unknown MOVEM mode %d", m.mode))
}

func (m *movem) firstAndLastReg(value byte) (left byte, right byte) {
	var rigthMostBitSet byte = 0
	var leftMostBitSet byte = 7

	for x := 0; x < 8; x++ {
		if (value>>x)&0b00000001 == 0b00000001 {
			rigthMostBitSet = byte(x)
			break
		}
	}

	for x := 7; x >= 0; x-- {
		if (value>>x)&0b00000001 == 0b00000001 {
			leftMostBitSet = byte(x)
			break
		}
	}

	return leftMostBitSet, rigthMostBitSet
}
