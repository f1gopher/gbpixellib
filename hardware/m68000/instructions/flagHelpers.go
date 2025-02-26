package instructions

import (
	"encoding/binary"
	"fmt"
	"math"
)

type instSize int

const (
	zeroSize instSize = iota
	byteSize
	wordSize
	longSize
)

func (i instSize) string() string {
	switch i {
	case byteSize:
		return "b"
	case wordSize:
		return "w"
	case longSize:
		return "l"
	}
	panic("")
}

func (i instSize) bytesUsed() int {
	switch i {
	case zeroSize:
		return 0
	case byteSize, wordSize:
		return 2
	case longSize:
		return 4
	}
	panic("")
}

func oneBitSize(value uint16, offset byte) instSize {
	bit := (value >> offset) & 0b0000000000000001
	if bit == 0 {
		return wordSize
	}

	return longSize
}

func twoBitSize(value uint16, offset byte) instSize {
	bits := (value >> offset) & 0b0000000000000011
	switch bits {
	case 0b0000000000000000:
		return byteSize
	case 0b0000000000000001:
		return wordSize
	case 0b0000000000000010:
		return longSize
	}

	panic(fmt.Sprintf("Unknown size: %016b", bits))
}

func moveSize(value uint16, offset byte) instSize {
	bits := (value >> offset) & 0b0000000000000011
	switch bits {
	case 0b0000000000000001:
		return byteSize
	case 0b0000000000000011:
		return wordSize
	case 0b0000000000000010:
		return longSize
	}

	panic(fmt.Sprintf("Unknown size: %016b", bits))
}

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
)

func addressModeBytesUsed(mode addressMode) instSize {
	switch mode {
	case absoluteShort, addressRegisterIndirectWithIndexAnd8BitDisplacement, addressRegisterIndirectWithDisplacement:
		return wordSize
	case absoluteLong:
		return longSize
	case dataRegisterDirect:
		return zeroSize
	}

	panic(fmt.Sprintf("Unknow address mode for size: %d", mode))
}

func addressString(mode addressMode, reg uint8, value uint32, displacement int) string {
	switch mode {
	case dataRegisterDirect:
		return fmt.Sprintf("D%d", reg)
	case addressRegisterDirect:
		return fmt.Sprintf("A%d", reg)
	case addressRegisterIndirect:
		return fmt.Sprintf("(A%d)", reg)
	case addressRegisterIndirectWithPostIncrement:
		return fmt.Sprintf("(A%d)+", reg)
	case absoluteShort:
		return fmt.Sprintf("$%x.w", value)
	case absoluteLong:
		return fmt.Sprintf("$%x.l", value)
	case preIndexedImediate:
		return fmt.Sprintf("#$%x", value)
	case addressRegisterIndirectWithDisplacement:
		if displacement < 0 {
			return fmt.Sprintf("(-$%x,A%d)", int(math.Abs(float64(displacement))), reg)
		}
		return fmt.Sprintf("($%x,A%d)", value, reg)
	case programCounterIndirectWithDisplacement:
		if displacement < 0 {
			return fmt.Sprintf("(-$%x,PC) ; ($%x)", value, displacement)
		}
		return fmt.Sprintf("($%x,PC) ; ($%x)", value, displacement)
	case addressRegisterIndirectWithIndexAnd8BitDisplacement:
		return fmt.Sprintf("%d", value)
	case addressRegisterIndirectWithPreDecrement:
		return fmt.Sprintf("-(A%d)", reg)
	}

	panic(fmt.Sprintf("Unknown address mode %d", mode))
}

func valueFrom3Bits(value uint16, offset byte) byte {
	return byte((value >> offset) & 0b0000000000000111)
}

func effective(value uint16, offset byte, size instSize) (mode addressMode, register byte) {
	return source(value, offset, size)
}

func source(value uint16, offset byte, size instSize) (mode addressMode, register byte) {
	if value&0b0000000000111111 == 0b000000000011110 {
		return preIndexedImediate, 0
	}

	modeBits := valueFrom3Bits(value, offset+3)
	register = valueFrom3Bits(value, offset)

	return addressingMode(modeBits, register), register
}

func destination(value uint16, offset byte, size instSize) (mode addressMode, register byte) {
	modeBits := valueFrom3Bits(value, offset)
	register = valueFrom3Bits(value, offset+3)

	return addressingMode(modeBits, register), register
}

func addressingMode(modeBits byte, reg byte) addressMode {
	switch modeBits & 0b000000111 {
	case 0b00000000:
		return dataRegisterDirect
	case 0b00000001:
		return addressRegisterDirect
	case 0b00000010:
		return addressRegisterIndirect
	case 0b00000011:
		return addressRegisterIndirectWithPostIncrement
	case 0b000000100:
		return addressRegisterIndirectWithPreDecrement
	case 0b000000101:
		return addressRegisterIndirectWithDisplacement
	case 0b000000110:
		return addressRegisterIndirectWithIndexAnd8BitDisplacement
		// return addressRegisterIndirectWithIndexAndBaseDisplacement
		// return memoryIndirectPostIndexed
		// return memoryIndirectPreIndexed
	case 0b000000111:
		switch reg {
		case 0b00000000:
			return absoluteShort
		case 0b00000001:
			return absoluteLong
		case 0b00000010:
			return programCounterIndirectWithDisplacement
		case 0b00000100:
			return preIndexedImediate
		}
		// return programCounterIndirectWithDisplacement
		panic("")
	}

	panic(fmt.Sprintf("Unknown address mode: %016b", modeBits))
}

func readData(code []uint8, offset int, size instSize) uint32 {
	switch size {
	case byteSize:
		// We want the lower byte
		return uint32(code[offset+1])
	case wordSize:
		return uint32(binary.BigEndian.Uint16(code[offset:]))
	case longSize:
		return binary.BigEndian.Uint32(code[offset:])
	default:
		panic("")
	}
}

func twosComplementToInt8Bit(value uint8) int {
	isNegative := (value & 0b10000000) == 0b10000000

	// For negative invert the bits and add one
	if isNegative {
		value = value ^ 0xFF
		value += 1

		return -int(value)
	}

	return int(value)
}

func twosComplementToInt16Bit(value uint16) int {
	isNegative := (value & 0b1000000000000000) == 0b1000000000000000

	// For negative invert the bits and add one
	if isNegative {
		value = value ^ 0xFFFF
		value += 1

		return -int(value)
	}

	return int(value)
}
