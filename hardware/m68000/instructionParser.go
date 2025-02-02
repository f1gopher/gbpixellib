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
	dc2Mask            uint16 = 0b1111111100000000
	dc2Prefix          uint16 = 0b0000111000000000
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
	moveqMask          uint16 = 0b1111000000000000
	moveqPrefix        uint16 = 0b0111000000000000
	dbccMask           uint16 = 0b1111000011111000
	dbccPrefix         uint16 = 0b0101000011001000
	clrMask            uint16 = 0b1111111100000000
	clrPrefix          uint16 = 0b0100001000000000
	jsrMask            uint16 = 0b1111111111000000
	jsrPrefix          uint16 = 0b0100111010000000
	moveToSRMask       uint16 = 0b1111111111000000
	moveToSRPrefix     uint16 = 0b0100011011000000
	braMask            uint16 = 0b1111111100000000
	braPrefix          uint16 = 0b0110000000000000
	btstMask           uint16 = 0b0000000000000000
	btstPrefix         uint16 = 0b0000000000000000
)

func (i *instructionParser) parse(code []uint8, offset int) instruction {
	// Read two bytes (the minimum for all instructions) and move the offset
	instruction := binary.BigEndian.Uint16(code[offset:])
	offset += 2

	operation := instruction >> 12 & 0b0000000000001111
	switch operation {
	case 0b0000000000000000:
		// Bit Manipulation/MOVEP/Immed iate
		if instruction&dcPrefix == dcPrefix || instruction&dc2Mask == dc2Prefix {
			return instructions.CreateDC(instruction)
		} else if instruction&oriMask == oriPrefix {
			return instructions.CreateORI(instruction, code, offset)
		} else if instruction&btstMask == btstPrefix {
			return instructions.CreateBTST(instruction, code, offset)
		}
	case 0b0000000000000001:
		// Move Byte
		return instructions.CreateMOVE(instruction, code, offset)
	case 0b0000000000000010:
		// Move Long
		if instruction&moveaMask == moveaPrefix {
			return instructions.CreateMOVEA(instruction, code, offset)
		} else if instruction&moveMask == moveLPrefix {
			return instructions.CreateMOVE(instruction, code, offset)
		}
	case 0b0000000000000011:
		// Move Word
		return instructions.CreateMOVE(instruction, code, offset)
	case 0b0000000000000100:
		// Miscellaneous
		if instruction == illegalInstruction {
			return instructions.CreateILLEGAL()
		} else if instruction == rteInstruction {
			return instructions.CreateRTE()
		} else if instruction&clrMask == clrPrefix {
			return instructions.CreateCLR(instruction)
		} else if instruction&jsrMask == jsrPrefix {
			return instructions.CreateJSR(instruction, code, offset)
		} else if instruction&moveToSRMask == moveToSRPrefix {
			return instructions.CreateMOVEToSR(instruction, code, offset)
		} else if instruction&leaMask == leaPrefix {
			return instructions.CreateLEA(instruction, code, offset)
		}
	case 0b0000000000000101:
		// ADDQ/SUBQ/Scc/DBcc/TRAPc c
		if instruction&dbccMask == dbccPrefix {
			return instructions.CreateDB(instruction, code, offset)
		}
	case 0b0000000000000110:
		// Bcc/BSR/BRA
		if instruction&braMask == braPrefix {
			return instructions.CreateBRA(instruction, code, offset)
		}
	case 0b0000000000000111:
		// MOVEQ
		if instruction&moveqMask == moveqPrefix {
			return instructions.CreateMOVEQ(instruction)
		}
	case 0b0000000000001000:
		// OR/DIV/SBCD
	case 0b0000000000001001:
		// SUB/SUBX
	case 0b0000000000001010:
		// (Unassigned, Reserved)
	case 0b0000000000001011:
		// CMP/EOR
	case 0b0000000000001100:
		// AND/MUL/ABCD/EXG
	case 0b0000000000001101:
		// ADD/ADDX
	case 0b0000000000001110:
		// Shift/Rotate/Bit Field
	case 0b0000000000001111:
		// Coprocessor Interface/MC68040 and CPU32 Extensions

		// TODO - don't think this is right
		if instruction&dcPrefix == dcPrefix {
			return instructions.CreateDC(instruction)
		}
	}

	return &instructions.Unknown{}
}
