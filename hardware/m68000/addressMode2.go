package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

func readEffectiveAddressMode2(
	value uint16,
	modeOffset byte,
	size instSize,
	memory memory.Memory[uint32],
	codeOffset uint32,
	programCounter uint32,
) (addrMode addressModeInfo, newOffset uint32, bytesUsed uint32) {
	return readAddressMode2(
		value,
		size,
		memory,
		codeOffset,
		programCounter,
		modeOffset+3,
		modeOffset)
}

func readSrcAddressMode2(
	value uint16,
	modeOffset byte,
	size instSize,
	memory memory.Memory[uint32],
	codeOffset uint32,
	programCounter uint32,
) (addrMode addressModeInfo, newOffset uint32, bytesUsed uint32) {
	return readAddressMode2(
		value,
		size,
		memory,
		codeOffset,
		programCounter,
		modeOffset+3,
		modeOffset)
}

func readDestAddressMode2(
	value uint16,
	modeOffset byte,
	size instSize,
	memory memory.Memory[uint32],
	codeOffset uint32,
	programCounter uint32,
) (addrMode addressModeInfo, newOffset uint32, bytesUsed uint32) {
	return readAddressMode2(
		value,
		size,
		memory,
		codeOffset,
		programCounter,
		modeOffset,
		modeOffset+3)
}

func readAddressMode2(
	value uint16,
	size instSize,
	code memory.Memory[uint32],
	codeOffset uint32,
	programCounter uint32,
	modeOffset byte,
	registerOffset byte,
) (addrMode addressModeInfo, newOffset uint32, bytesUsed uint32) {
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
		// TODO - this should pull from the prefetch queue instead of memory
		//	displacement = int32(signExtend16BitTo32Bits(code.ReadShort(codeOffset))) //uint16(readData(code, codeOffset, wordSize))))
		bytesUsed = 2
	} else if mode == preIndexedImediate {
		//data = readData(code, codeOffset, size)
		if size.bytesUsed() == 1 {
			data = uint32(code.ReadByte(codeOffset))
		} else if size.bytesUsed() == 2 {
			data = uint32(code.ReadShort(codeOffset))
		} else if size.bytesUsed() == 4 {
			data = code.ReadLong(codeOffset)
		} else {
			panic("")
		}
		bytesUsed = uint32(size.bytesUsed())
	} else if mode == dataRegisterDirect || mode == addressRegisterDirect {
		bytesUsed = 0
	} else if mode == absoluteWord {
		data = uint32(code.ReadShort(codeOffset)) //readData(code, codeOffset, wordSize)
		bytesUsed = 2
	} else if mode == absoluteLong {
		data = code.ReadLong(codeOffset) //readData(code, codeOffset, longSize)
		bytesUsed = 4
	} else if mode == programCounterIndirectWithIndex8bitDisplacement {
		extraData := code.ReadShort(codeOffset) //uint16(readData(code, codeOffset, wordSize))
		// Only want the first 8bits of info
		displacement = int32(twosComplementToInt8Bit(byte(extraData)))
		data = uint32(valueFrom3Bits(extraData, 12))
		bytesUsed = 2
	} else if mode == addressRegisterIndirectWithIndexAnd8BitDisplacement {
		extraData := code.ReadShort(codeOffset) //uint16(readData(code, codeOffset, wordSize))

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
