package m68000

import (
	"encoding/binary"

	"github.com/f1gopher/gbpixellib/hardware/m68000/instructions"
)

type instructionParser struct{}

func createInstructionParser() instructionParser {
	return instructionParser{}
}

func (i *instructionParser) dump(code []uint8) []string {
	result := []string{}

	for x := 0; x < len(code); {
		inst := i.parse(code, x)
		x += inst.Length()

		result = append(result, inst.Name())
	}

	return result
}

const (
	illegalInstruction uint16 = 0b0100101011111100
	rteInstruction     uint16 = 0b0100111001110011
	dcPrefix           uint16 = 0x00FF
	moveMask           uint16 = 0b1111000000000000
	moveBPrefix        uint16 = 0b0001000000000000
	moveWPrefix        uint16 = 0b0011000000000000
	moveLPrefix        uint16 = 0b0010000000000000
	oriMask            uint16 = 0b1111111100000000
	oriPrefix          uint16 = 0b0000000000000000
	leaMask            uint16 = 0b1111000000000000
	leaPrefix          uint16 = 0b0100000000000000
	moveaMask          uint16 = 0b1110000111000000
	moveaPrefix        uint16 = 0b0010000001000000
)

func (i *instructionParser) parse(code []uint8, offset int) instruction {
	// Read two bytes (the minimum for all instructions) and move the offset
	instruction := binary.BigEndian.Uint16(code[offset:])
	offset += 2

	if instruction == illegalInstruction {
		return instructions.CreateILLEGAL()
	} else if instruction == rteInstruction {
		return instructions.CreateRTE()
	} else if instruction&dcPrefix == dcPrefix {
		return instructions.CreateDC(instruction)
	} else if instruction&oriMask == oriPrefix {
		return instructions.CreateORI(instruction, code, offset)
	} else if instruction&moveaMask == moveaPrefix {
		return instructions.CreateMOVEA(instruction, code, offset)
	} else if instruction&moveMask == moveBPrefix ||
		instruction&moveMask == moveWPrefix ||
		instruction&moveMask == moveLPrefix {
		return instructions.CreateMOVE(instruction, code, offset)
	} else if instruction&leaMask == leaPrefix {
		return instructions.CreateLEA(instruction, code, offset)
	}

	return &instructions.Unknown{}
}
