package instructions

import (
	"encoding/binary"
	"fmt"
)

type instSize int

const (
	byteSize instSize = iota
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
	case byteSize, wordSize:
		return 2
	case longSize:
		return 4
	}
	panic("")
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

func addressString(mode addressMode, reg uint8, value uint32) string {
	switch mode {
	case dataRegisterDirect:
		return fmt.Sprintf("D%d", reg)
	case addressRegisterDirect:
		return fmt.Sprintf("A%d", reg)
	case addressRegisterIndirect:
		return fmt.Sprintf("(A)%d", reg)
	case addressRegisterIndirectWithPostIncrement:
		return fmt.Sprintf("(A%d)+", reg)
	case absoluteShort, absoluteLong:
		return fmt.Sprintf("#$%x", value)
	case preIndexedImediate:
		return fmt.Sprintf("#$%d", value)
	case addressRegisterIndirectWithDisplacement:
		return fmt.Sprintf("($%x,A%d)", value, reg)
	}

	panic(fmt.Sprintf("Unknown address mode %d", mode))
}

func valueFrom3Bits(value uint16, offset byte) byte {
	return byte((value >> offset) & 0b0000000000000111)
}

func effective(value uint16, offset byte) (mode addressMode, register byte) {
	return source(value, offset)
}

func source(value uint16, offset byte) (mode addressMode, register byte) {
	if value&0b0000000000111111 == 0b000000000011110 {
		return preIndexedImediate, 0
	}

	modeBits := valueFrom3Bits(value, offset+3)
	register = valueFrom3Bits(value, offset)

	return addressingMode(modeBits), register
}

func destination(value uint16, offset byte) (mode addressMode, register byte) {
	modeBits := valueFrom3Bits(value, offset)
	register = valueFrom3Bits(value, offset+3)

	return addressingMode(modeBits), register
}

func addressingMode(modeBits byte) addressMode {
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
		return absoluteShort
		// return absoluteLong
		// return programCounterIndirectWithDisplacement
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
