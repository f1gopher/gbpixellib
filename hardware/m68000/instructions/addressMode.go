package instructions

import (
	"fmt"
	"math"
)

type addressMode int

const (
	dataRegisterDirect addressMode = iota
	addressRegisterDirect
	addressRegisterIndirect
	addressRegisterIndirectWithPostIncrement
	addressRegisterIndirectWithPreDecrement
	addressRegisterIndirectWithDisplacement
	addressRegisterIndirectWithIndexAnd8BitDisplacement
	addressRegisterIndirectWithIndexAndBaseDisplacement
	memoryIndirectPostIndexed
	memoryIndirectPreIndexed
	absoluteShort
	absoluteLong
	programCounterIndirectWithDisplacement
	preIndexedImediate
	programCounterIndirectWithIndex8bitDisplacement
)

type addressModeInfo struct {
	size         instSize
	mode         addressMode
	reg          byte
	data         uint32
	displacement int32
	pc           int32
}

func readEffectiveAddressMode(
	value uint16,
	modeOffset byte,
	size instSize,
	code []uint8,
	codeOffset int,
	programCounter int,
) (addrMode addressModeInfo, newOffset int, bytesUsed int) {
	return readAddressMode(
		value,
		size,
		code,
		codeOffset,
		programCounter,
		modeOffset+3,
		modeOffset)
}

func readSrcAddressMode(
	value uint16,
	modeOffset byte,
	size instSize,
	code []uint8,
	codeOffset int,
	programCounter int,
) (addrMode addressModeInfo, newOffset int, bytesUsed int) {
	return readAddressMode(
		value,
		size,
		code,
		codeOffset,
		programCounter,
		modeOffset+3,
		modeOffset)
}

func readDestAddressMode(
	value uint16,
	modeOffset byte,
	size instSize,
	code []uint8,
	codeOffset int,
	programCounter int,
) (addrMode addressModeInfo, newOffset int, bytesUsed int) {
	return readAddressMode(
		value,
		size,
		code,
		codeOffset,
		programCounter,
		modeOffset,
		modeOffset+3)
}

func readAddressMode(
	value uint16,
	size instSize,
	code []uint8,
	codeOffset int,
	programCounter int,
	modeOffset byte,
	registerOffset byte,
) (addrMode addressModeInfo, newOffset int, bytesUsed int) {
	// TODO - do we need to shift by offset for when the mode isn't bit 0?
	if value&0b0000000000111111 == 0b000000000011110 {
		return addressModeInfo{
			size:         zeroSize,
			mode:         preIndexedImediate,
			reg:          0,
			data:         0,
			displacement: 0,
			pc:           0,
		}, codeOffset, 0
	}

	modeBits := valueFrom3Bits(value, modeOffset)
	register := valueFrom3Bits(value, registerOffset)
	bytesUsed = 0

	mode := addressingMode(modeBits, register)
	var displacement int32 = 0
	var data uint32 = 0

	if mode == addressRegisterIndirectWithDisplacement ||
		mode == programCounterIndirectWithDisplacement {
		displacement = int32(signExtend16BitTo32Bits(uint16(readData(code, codeOffset, wordSize))))
		bytesUsed = 2
	} else if mode == preIndexedImediate {
		data = readData(code, codeOffset, size)
		bytesUsed = size.bytesUsed()
	} else if mode == dataRegisterDirect || mode == addressRegisterDirect {
		bytesUsed = 0
	} else if mode == absoluteShort {
		data = readData(code, codeOffset, wordSize)
		bytesUsed = 2
	} else if mode == absoluteLong {
		data = readData(code, codeOffset, longSize)
		bytesUsed = 4
	} else if mode == programCounterIndirectWithIndex8bitDisplacement {
		data = readData(code, codeOffset, wordSize)
		bytesUsed = 2
	}

	return addressModeInfo{
		size:         zeroSize,
		mode:         mode,
		reg:          register,
		data:         data,
		displacement: displacement,
		pc:           int32(programCounter),
	}, codeOffset + bytesUsed, bytesUsed
}

func (a *addressModeInfo) string() string {
	switch a.mode {
	case dataRegisterDirect:
		return fmt.Sprintf("D%d", a.reg)
	case addressRegisterDirect:
		return fmt.Sprintf("A%d", a.reg)
	case addressRegisterIndirect:
		return fmt.Sprintf("(A%d)", a.reg)
	case addressRegisterIndirectWithPostIncrement:
		return fmt.Sprintf("(A%d)+", a.reg)
	case absoluteShort:
		return fmt.Sprintf("$%x.w", a.data)
	case absoluteLong:
		return fmt.Sprintf("$%x.l", a.data)
	case preIndexedImediate:
		return fmt.Sprintf("#$%x", a.data)
	case addressRegisterIndirectWithDisplacement:
		if a.displacement < 0 {
			return fmt.Sprintf("(-$%x,A%d)", int(math.Abs(float64(a.displacement))), a.reg)
		}
		return fmt.Sprintf("($%x,A%d)", a.displacement, a.reg)
	case programCounterIndirectWithDisplacement:
		if a.displacement < 0 {
			return fmt.Sprintf("(-$%x,PC) ; ($%x)", -a.displacement, a.pc+a.displacement)
		}
		return fmt.Sprintf("($%x,PC) ; ($%x)", a.displacement, a.pc+a.displacement)
	case addressRegisterIndirectWithIndexAnd8BitDisplacement:
		return fmt.Sprintf("%d", a.data)
	case addressRegisterIndirectWithPreDecrement:
		return fmt.Sprintf("-(A%d)", a.reg)
	case programCounterIndirectWithIndex8bitDisplacement:
		return fmt.Sprintf("(PC,D1.w)")
	}

	panic(fmt.Sprintf("Unknown address mode %d", a.mode))
}
