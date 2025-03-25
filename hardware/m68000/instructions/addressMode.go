package instructions

import (
	"fmt"
	"math"
)

type addressMode int

const (
	invalidAddressMode addressMode = iota
	dataRegisterDirect
	addressRegisterDirect
	addressRegisterIndirect
	addressRegisterIndirectWithPostIncrement
	addressRegisterIndirectWithPreDecrement
	addressRegisterIndirectWithDisplacement
	addressRegisterIndirectWithIndexAnd8BitDisplacement
	addressRegisterIndirectWithIndexAndBaseDisplacement
	memoryIndirectPostIndexed
	memoryIndirectPreIndexed
	absoluteWord
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
	invalid      string
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
	// if value&0b0000000000111111 == 0b000000000011110 {
	// 	return addressModeInfo{
	// 		size:         zeroSize,
	// 		mode:         preIndexedImediate,
	// 		reg:          0,
	// 		data:         0,
	// 		displacement: 0,
	// 		pc:           0,
	// 	}, codeOffset, 0
	// }

	modeBits := valueFrom3Bits(value, modeOffset)
	register := valueFrom3Bits(value, registerOffset)
	bytesUsed = 0
	invalid := ""
	paramSize := zeroSize

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
	} else if mode == absoluteWord {
		data = readData(code, codeOffset, wordSize)
		bytesUsed = 2
	} else if mode == absoluteLong {
		data = readData(code, codeOffset, longSize)
		bytesUsed = 4
	} else if mode == programCounterIndirectWithIndex8bitDisplacement {
		extraData := uint16(readData(code, codeOffset, wordSize))
		// Only want the first 8bits of info
		displacement = int32(twosComplementToInt8Bit(byte(extraData)))
		data = uint32(valueFrom3Bits(extraData, 12))
		bytesUsed = 2
	} else if mode == addressRegisterIndirectWithIndexAnd8BitDisplacement {
		extraData := uint16(readData(code, codeOffset, wordSize))

		if ((extraData >> 8) & 0b0000000000000111) != 0 {

			if modeOffset > registerOffset {
				invalid = fmt.Sprintf("INVALID %x", value&0x3f)
			} else {
				// Need to swap mode and register bits around so they match the expected output
				t := ((value >> uint16(registerOffset)) & 0b111) | ((value >> (uint16(modeOffset) - 3)) & 0b111000)
				invalid = fmt.Sprintf("INVALID %x", t&0x3f)
			}
		}

		// Only want the first 8bits of info
		displacement = int32(twosComplementToInt8Bit(uint8(extraData)))

		data = uint32(valueFrom3Bits(extraData, 12))
		bytesUsed = 2
		paramSize = wordSize
		if extraData&0b0000100000000000 == 0b0000100000000000 {
			paramSize = longSize
		}

		// TODO - bit 10 of extraData might be used to chose Address ot Data register (1 for data)
	}

	return addressModeInfo{
		size:         paramSize,
		mode:         mode,
		reg:          register,
		data:         data,
		displacement: displacement,
		pc:           int32(programCounter),
		invalid:      invalid,
	}, codeOffset + bytesUsed, bytesUsed
}

func addressingMode(modeBits byte, reg byte) addressMode {
	switch modeBits & 0b00000111 {
	case 0b00000000:
		return dataRegisterDirect
	case 0b00000001:
		return addressRegisterDirect
	case 0b00000010:
		return addressRegisterIndirect
	case 0b00000011:
		return addressRegisterIndirectWithPostIncrement
	case 0b00000100:
		return addressRegisterIndirectWithPreDecrement
	case 0b00000101:
		return addressRegisterIndirectWithDisplacement
	case 0b00000110:
		return addressRegisterIndirectWithIndexAnd8BitDisplacement
		// return addressRegisterIndirectWithIndexAndBaseDisplacement
		// return memoryIndirectPostIndexed
		// return memoryIndirectPreIndexed
	case 0b00000111:
		switch reg {
		case 0b00000000:
			return absoluteWord
		case 0b00000001:
			return absoluteLong
		case 0b00000010:
			return programCounterIndirectWithDisplacement
		case 0b00000100:
			return preIndexedImediate
		case 0b00000011:
			return programCounterIndirectWithIndex8bitDisplacement
		}
		// return programCounterIndirectWithDisplacement
		// panic(fmt.Sprintf("reg: %03b", reg))
		return invalidAddressMode
	}

	panic(fmt.Sprintf("Unknown address mode: %016b", modeBits))
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
	case absoluteWord:
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
		if a.invalid != "" {
			return a.invalid
		}

		if a.displacement == 0 {
			return fmt.Sprintf("(A%d,D%d.%s)", a.reg, a.data, a.size.string())
		}

		if a.displacement < 0 {
			return fmt.Sprintf("(-$%x,A%d,D%d.%s)", -a.displacement, a.reg, a.data, a.size.string())
		}
		return fmt.Sprintf("($%x,A%d,D%d.%s)", a.displacement, a.reg, a.data, a.size.string())
	case addressRegisterIndirectWithPreDecrement:
		return fmt.Sprintf("-(A%d)", a.reg)
	case programCounterIndirectWithIndex8bitDisplacement:
		if a.displacement == 0 {
			return fmt.Sprintf("(PC,D%d.w)", a.data)
		}
		if a.displacement < 0 {
			return fmt.Sprintf("(-$%x,PC,D%d.w)", -a.displacement, a.data)
		}

		return fmt.Sprintf("($%x,PC,D%d.w)", a.displacement, a.data)
	}

	// panic(fmt.Sprintf("Unknown address mode %d", a.mode))
	return "unknown address mode"
}
