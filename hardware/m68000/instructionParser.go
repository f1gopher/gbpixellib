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
	btstMask           uint16 = 0b1111000111000000
	btstPrefix         uint16 = 0b0000000100000000
	btstMask2          uint16 = 0b1111111111000000
	btstPrefix2        uint16 = 0b0000100000000000
	bclrMask           uint16 = 0b1111000111000000
	bclrPrefix         uint16 = 0b0000000110000000
	bclrMask2          uint16 = 0b1111111111000000
	bclrPrefix2        uint16 = 0b0000100010000000
	cmpiMask           uint16 = 0b1111111100000000
	cmpiPrefix         uint16 = 0b0000110000000000
	eoriMask           uint16 = 0b1111111100000000
	eoriPrefix         uint16 = 0b0000101000000000
	tstMask            uint16 = 0b1111111100000000
	tstPrefix          uint16 = 0b0100101000000000
	bccMask            uint16 = 0b1111000000000000
	bccPrefix          uint16 = 0b0110000000000000
	addiMask           uint16 = 0b1111111100000000
	addiPrefix         uint16 = 0b0000011000000000
	rtsMask            uint16 = 0b1111111111111111
	rtsPrefix          uint16 = 0b0100111001110101
	movemMask          uint16 = 0b1111101110000000
	movemPrefix        uint16 = 0b0100100010000000
	addqMask           uint16 = 0b1111000100000000
	addqPrefix         uint16 = 0b0101000000000000
	bsrMask            uint16 = 0b1111111100000000
	bsrPrefix          uint16 = 0b0110000100000000
	sccMask            uint16 = 0b1111111111000000
	sccPrefix          uint16 = 0b0101000011000000
	subqMask           uint16 = 0b1111000100000000
	subqPrefix         uint16 = 0b0101000100000000
	addMask            uint16 = 0b1111000000000000
	addPrefix          uint16 = 0b1101000000000000
	subMask            uint16 = 0b1111000000000000
	subPrefix          uint16 = 0b1001000000000000
	eorMask            uint16 = 0b1111000000000000
	eorPrefix          uint16 = 0b1011000000000000
	ror_rolMask        uint16 = 0b1111000000011000
	ror_rolPrefix      uint16 = 0b1110000000011000
	orMask             uint16 = 0b1111000000000000
	orPrefix           uint16 = 0b1000000000000000
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
		} else if instruction&cmpiMask == cmpiPrefix {
			return instructions.CreateCMPI(instruction, code, offset)
		} else if instruction&bclrMask == bclrPrefix || instruction&bclrMask2 == bclrPrefix2 {
			return instructions.CreateBCLR(instruction, code, offset)
		} else if instruction&eoriMask == eoriPrefix {
			return instructions.CreateEORI(instruction, code, offset)
		} else if instruction&btstMask == btstPrefix || instruction&btstMask2 == btstPrefix2 {
			return instructions.CreateBTST(instruction, code, offset)
		} else if instruction&addiMask == addiPrefix {
			return instructions.CreateADDI(instruction, code, offset)
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
		} else if instruction&rtsMask == rtsPrefix {
			return instructions.CreateRTS(instruction)
		} else if instruction&clrMask == clrPrefix {
			return instructions.CreateCLR(instruction)
		} else if instruction&jsrMask == jsrPrefix {
			return instructions.CreateJSR(instruction, code, offset)
		} else if instruction&moveToSRMask == moveToSRPrefix {
			return instructions.CreateMOVEToSR(instruction, code, offset)
		} else if instruction&tstMask == tstPrefix {
			return instructions.CreateTST(instruction, code, offset)
		} else if instruction&movemMask == movemPrefix {
			return instructions.CreateMOVEM(instruction, code, offset)
		} else if instruction&leaMask == leaPrefix {
			return instructions.CreateLEA(instruction, code, offset)
		}
	case 0b0000000000000101:
		// ADDQ/SUBQ/Scc/DBcc/TRAPc c
		if instruction&dbccMask == dbccPrefix {
			return instructions.CreateDB(instruction, code, offset)
		} else if instruction&sccMask == sccPrefix {
			return instructions.CreateScc(instruction, code, offset)
		} else if instruction&addqMask == addqPrefix {
			return instructions.CreateADDQ(instruction, code, offset)
		} else if instruction&subqMask == subqPrefix {
			return instructions.CreateSUBQ(instruction, code, offset)
		}
	case 0b0000000000000110:
		// Bcc/BSR/BRA
		if instruction&bsrMask == bsrPrefix {
			return instructions.CreateBSR(instruction, code, offset)
		} else if instruction&braMask == braPrefix {
			return instructions.CreateBRA(instruction, code, offset)
		} else if instruction&bccMask == bccPrefix {
			return instructions.CreateBcc(instruction, code, offset)
		}
	case 0b0000000000000111:
		// MOVEQ
		if instruction&moveqMask == moveqPrefix {
			return instructions.CreateMOVEQ(instruction)
		}
	case 0b0000000000001000:
		// OR/DIV/SBCD
		if instruction&orMask == orPrefix {
			return instructions.CreateOR(instruction, code, offset)
		}
	case 0b0000000000001001:
		// SUB/SUBX
		if instruction&subMask == subPrefix {
			return instructions.CreateSUB(instruction, code, offset)
		}
	case 0b0000000000001010:
		// (Unassigned, Reserved)
	case 0b0000000000001011:
		// CMP/EOR
		if instruction&eorMask == eorPrefix {
			return instructions.CreateEOR(instruction, code, offset)
		}
	case 0b0000000000001100:
		// AND/MUL/ABCD/EXG
	case 0b0000000000001101:
		// ADD/ADDX
		if instruction&addMask == addPrefix {
			return instructions.CreateADD(instruction, code, offset)
		}
	case 0b0000000000001110:
		// Shift/Rotate/Bit Field
		if instruction&ror_rolMask == ror_rolPrefix {
			return instructions.CreateROR_ROL(instruction, code, offset)
		}
	case 0b0000000000001111:
		// Coprocessor Interface/MC68040 and CPU32 Extensions

		// TODO - don't think this is right
		if instruction&dcPrefix == dcPrefix {
			return instructions.CreateDC(instruction)
		}
	}

	return &instructions.Unknown{}
}
