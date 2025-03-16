package instructions

import "fmt"

type moveDirection int

const (
	registerToMemory moveDirection = iota
	memoryToRegister
)

type movem struct {
	effective addressModeInfo

	transferSize     instSize
	registerListMask uint16
	direction        moveDirection
}

func CreateMOVEM(value uint16, code []uint8, codeOffset int) *movem {
	size := oneBitSize(value, 6)
	eff, _, _ := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
	direction := registerToMemory
	if ((value >> 10) & 0b0000000000000001) == 0b0000000000000001 {
		direction = memoryToRegister
	}

	if direction == registerToMemory {
		if eff.mode == dataRegisterDirect ||
			eff.mode == addressRegisterDirect ||
			eff.mode == addressRegisterIndirectWithPostIncrement ||
			eff.mode == programCounterIndirectWithDisplacement ||
			eff.mode == preIndexedImediate ||
			eff.mode == programCounterIndirectWithIndex8bitDisplacement {
			return nil
		}
	} else {
		if eff.mode == dataRegisterDirect ||
			eff.mode == addressRegisterDirect ||
			eff.mode == addressRegisterIndirectWithPreDecrement ||
			eff.mode == programCounterIndirectWithDisplacement ||
			eff.mode == programCounterIndirectWithIndex8bitDisplacement {
			return nil
		}
	}

	return &movem{
		transferSize:     size,
		registerListMask: uint16(readData(code, codeOffset, wordSize)),
		direction:        direction,
		effective:        eff,
	}
}

func (m *movem) Name() string {
	if m.direction == memoryToRegister {
		return fmt.Sprintf(
			"movem.%s %s, %s",
			m.transferSize.string(),
			m.effective.string(),
			m.registers())
	}

	return fmt.Sprintf(
		"movem.%s %s, %s",
		m.transferSize.string(),
		m.registers(),
		m.effective.string())
}

func (m *movem) Length() int {
	return 2 + wordSize.bytesUsed()
}

func (m *movem) registers() string {
	if m.effective.mode == addressRegisterIndirectWithPreDecrement {
		// D then A
		aFirstReg, aLastReg := m.firstAndLastReg(byte(m.registerListMask))
		dFirstReg, dLastReg := m.firstAndLastReg(byte(m.registerListMask >> 8))

		if aFirstReg == aLastReg {
			return fmt.Sprintf("A%d", 7-aLastReg)
		}

		return fmt.Sprintf("D%d-D%d/A%d-A%d", 7-dFirstReg, 7-dLastReg, 7-aFirstReg, 7-aLastReg)
	}

	if m.effective.mode == addressRegisterIndirectWithPostIncrement {
		// A then D
		aLastReg, aFirstReg := m.firstAndLastReg(byte(m.registerListMask >> 8))
		dLastReg, dFirstReg := m.firstAndLastReg(byte(m.registerListMask))

		if aFirstReg == aLastReg {
			return fmt.Sprintf("A%d", aLastReg)
		}

		return fmt.Sprintf("D%d-D%d/A%d-A%d", dFirstReg, dLastReg, aFirstReg, aLastReg)
	}

	if m.effective.mode == addressRegisterIndirect {
		dLastReg, dFirstReg := m.firstAndLastReg(byte(m.registerListMask))
		return fmt.Sprintf("D%d-D%d", dFirstReg, dLastReg)
	}

	panic(fmt.Sprintf("Unknown MOVEM mode %d", m.effective.mode))
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
