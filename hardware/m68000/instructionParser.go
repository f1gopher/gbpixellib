package m68000

import (
	"encoding/binary"
	"reflect"

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
	leaMask            uint16 = 0b1111000111000000
	leaPrefix          uint16 = 0b0100000111000000
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
	sccMask            uint16 = 0b1111000011000000
	sccPrefix          uint16 = 0b0101000011000000
	subqMask           uint16 = 0b1111000100000000
	subqPrefix         uint16 = 0b0101000100000000
	addMask            uint16 = 0b1111000000000000
	addPrefix          uint16 = 0b1101000000000000
	subMask            uint16 = 0b1111000000000000
	subPrefix          uint16 = 0b1001000000000000
	rol_rorMask        uint16 = 0b1111000000011000
	rol_rorPrefix      uint16 = 0b1110000000011000
	orMask             uint16 = 0b1111000000000000
	orPrefix           uint16 = 0b1000000000000000
	notMask            uint16 = 0b1111111100000000
	notPrefix          uint16 = 0b0100011000000000
	lsl_lsrMask        uint16 = 0b1111000000000000
	lsl_lsrPrefix      uint16 = 0b1110000000000000
	andiMask           uint16 = 0b1111111100000000
	andiPrefix         uint16 = 0b0000001000000000
	stopMask           uint16 = 0b1111111111111111
	stopPrefix         uint16 = 0b0100111001110010
	peaMask            uint16 = 0b1111111111000000
	peaPrefix          uint16 = 0b0100100001000000
	extMask            uint16 = 0b1111111000111000
	extPrefix          uint16 = 0b0100100000000000
	jmpMask            uint16 = 0b1111111111000000
	jmpPrefix          uint16 = 0b0100111011000000
	moveUSPMask        uint16 = 0b1111111111100000
	moveUSPPrefix      uint16 = 0b0100111001100000
	trapMask           uint16 = 0b1111111111110000
	trapPrefix         uint16 = 0b0100111001000000
	moveCCRMask        uint16 = 0b1111111111000000
	moveCCRPrefix      uint16 = 0b0100010011000000
	subiMask           uint16 = 0b1111111100000000
	subiPrefix         uint16 = 0b0000010000000000
	andMask            uint16 = 0b1111000000000000
	andPrefix          uint16 = 0b1100000000000000
	roxl_roxrMask      uint16 = 0b1111000000011000
	roxl_roxrPrefix    uint16 = 0b1110000000010000
	bchgMask           uint16 = 0b1111000111000000
	bchgPrefix         uint16 = 0b0000000101000000
	movepMask          uint16 = 0b1111000000111000
	movepPrefix        uint16 = 0b0000000000001000
	bsetMask           uint16 = 0b1111000011000000
	bsetPrefix         uint16 = 0b0000000011000000
	divuMask           uint16 = 0b1111000111000000
	divuPrefix         uint16 = 0b1000000011000000
	abcdMask           uint16 = 0b1111000111110000
	abcdPrefix         uint16 = 0b1100000100000000
	andiToCCRMask      uint16 = 0b1111111111111111
	andiToCCRPrefix    uint16 = 0b0000001000111100
	sbcdMask           uint16 = 0b1111000111110000
	sbcdPrefix         uint16 = 0b1000000100000000
)

func (i *instructionParser) parse(code []uint8, offset int) instruction {
	// Read two bytes (the minimum for all instructions) and move the offset
	opcode := binary.BigEndian.Uint16(code[offset:])
	offset += 2
	var result instruction = nil

	operation := opcode >> 12 & 0b0000000000001111
	switch operation {
	case 0b0000000000000000:
		// Bit Manipulation/MOVEP/Immed iate
		if opcode&dcPrefix == dcPrefix || opcode&dc2Mask == dc2Prefix {
			result = instructions.CreateDC(opcode)
		} else if opcode&oriMask == oriPrefix {
			result = instructions.CreateORI(opcode, code, offset)
		} else if opcode&cmpiMask == cmpiPrefix {
			result = instructions.CreateCMPI(opcode, code, offset)
		} else if opcode&bclrMask == bclrPrefix || opcode&bclrMask2 == bclrPrefix2 {
			result = instructions.CreateBCLR(opcode, code, offset)
		} else if opcode&eoriMask == eoriPrefix {
			result = instructions.CreateEORI(opcode, code, offset)
		} else if opcode&movepMask == movepPrefix {
			result = instructions.CreateMOVEP(opcode, code, offset)
		} else if opcode&btstMask == btstPrefix || opcode&btstMask2 == btstPrefix2 {
			result = instructions.CreateBTST(opcode, code, offset)
		} else if opcode&addiMask == addiPrefix {
			result = instructions.CreateADDI(opcode, code, offset)
		} else if opcode&andiToCCRMask == andiToCCRPrefix {
			result = instructions.CreateANDIToCCR(opcode, code, offset)
		} else if opcode&andiMask == andiPrefix {
			result = instructions.CreateANDI(opcode, code, offset)
		} else if opcode&subiMask == subiPrefix {
			result = instructions.CreateSUBI(opcode, code, offset)
		} else if opcode&bchgMask == bchgPrefix {
			result = instructions.CreateBCHG(opcode, code, offset)
		} else if opcode&bsetMask == bsetPrefix {
			result = instructions.CreateBSET(opcode, code, offset)
		}
	case 0b0000000000000001:
		// Move Byte
		result = instructions.CreateMOVE(opcode, code, offset)
	case 0b0000000000000010:
		// Move Long
		if opcode&moveaMask == moveaPrefix {
			result = instructions.CreateMOVEA(opcode, code, offset)
		} else if opcode&moveMask == moveLPrefix {
			result = instructions.CreateMOVE(opcode, code, offset)
		}
	case 0b0000000000000011:
		// Move Word
		result = instructions.CreateMOVE(opcode, code, offset)
	case 0b0000000000000100:
		// Miscellaneous
		if opcode == illegalInstruction {
			result = instructions.CreateILLEGAL()
		} else if opcode == rteInstruction {
			result = instructions.CreateRTE()
		} else if opcode&rtsMask == rtsPrefix {
			result = instructions.CreateRTS(opcode)
		} else if opcode&clrMask == clrPrefix {
			result = instructions.CreateCLR(opcode, code, offset)
		} else if opcode&jsrMask == jsrPrefix {
			result = instructions.CreateJSR(opcode, code, offset)
		} else if opcode&moveToSRMask == moveToSRPrefix {
			result = instructions.CreateMOVEToSR(opcode, code, offset)
		} else if opcode&tstMask == tstPrefix {
			result = instructions.CreateTST(opcode, code, offset)
		} else if opcode&extMask == extPrefix {
			result = instructions.CreateEXT(opcode)
		} else if opcode&movemMask == movemPrefix {
			result = instructions.CreateMOVEM(opcode, code, offset)
		} else if opcode&notMask == notPrefix {
			result = instructions.CreateNOT(opcode, code, offset)
		} else if opcode&stopMask == stopPrefix {
			result = instructions.CreateSTOP(opcode, code, offset)
		} else if opcode&peaMask == peaPrefix {
			result = instructions.CreatePEA(opcode, code, offset)
		} else if opcode&jmpMask == jmpPrefix {
			result = instructions.CreateJMP(opcode, code, offset)
		} else if opcode&moveUSPMask == moveUSPPrefix {
			result = instructions.CreateMOVEToUSP(opcode)
		} else if opcode&moveCCRMask == moveCCRPrefix {
			result = instructions.CreateMOVEToCCR(opcode, code, offset)
		} else if opcode&trapMask == trapPrefix {
			result = instructions.CreateTRAP(opcode)
		} else if opcode&leaMask == leaPrefix {
			result = instructions.CreateLEA(opcode, code, offset)
		}
	case 0b0000000000000101:
		// ADDQ/SUBQ/Scc/DBcc/TRAPc c
		if opcode&dbccMask == dbccPrefix {
			result = instructions.CreateDB(opcode, code, offset)
		} else if opcode&sccMask == sccPrefix {
			result = instructions.CreateScc(opcode, code, offset)
		} else if opcode&addqMask == addqPrefix {
			result = instructions.CreateADDQ(opcode, code, offset)
		} else if opcode&subqMask == subqPrefix {
			result = instructions.CreateSUBQ(opcode, code, offset)
		}
	case 0b0000000000000110:
		// Bcc/BSR/BRA
		if opcode&bsrMask == bsrPrefix {
			result = instructions.CreateBSR(opcode, code, offset)
		} else if opcode&braMask == braPrefix {
			result = instructions.CreateBRA(opcode, code, offset)
		} else if opcode&bccMask == bccPrefix {
			result = instructions.CreateBcc(opcode, code, offset)
		}
	case 0b0000000000000111:
		// MOVEQ
		if opcode&moveqMask == moveqPrefix {
			result = instructions.CreateMOVEQ(opcode)
		}
	case 0b0000000000001000:
		// OR/DIV/SBCD
		if opcode&divuMask == divuPrefix {
			result = instructions.CreateDIVU(opcode, code, offset)
		} else if opcode&sbcdMask == sbcdPrefix {
			result = instructions.CreateSBCD(opcode, code, offset)
		} else if opcode&orMask == orPrefix {
			result = instructions.CreateOR(opcode, code, offset)
		}
	case 0b0000000000001001:
		// SUB/SUBX
		if opcode&subMask == subPrefix {
			result = instructions.CreateSUB(opcode, code, offset)
		}
	case 0b0000000000001010:
		// (Unassigned, Reserved)
	case 0b0000000000001011:
		// CMP/EOR
		opmode := (opcode >> 6) & 0b000000000000111

		if opmode == 0 || opmode == 0b00000001 || opmode == 0b00000010 {
			result = instructions.CreateCMP(opcode, code, offset)
		} else {
			result = instructions.CreateEOR(opcode, code, offset)
		}
	case 0b0000000000001100:
		// AND/MUL/ABCD/EXG
		if opcode&abcdMask == abcdPrefix {
			result = instructions.CreateABCD(opcode, code, offset)
		} else if opcode&andMask == andPrefix {
			result = instructions.CreateAND(opcode, code, offset)
		}
	case 0b0000000000001101:
		// ADD/ADDX
		if opcode&addMask == addPrefix {
			result = instructions.CreateADD(opcode, code, offset)
		}
	case 0b0000000000001110:
		// Shift/Rotate/Bit Field
		if opcode&rol_rorMask == rol_rorPrefix {
			result = instructions.CreateROL_ROR(opcode, code, offset)
		} else if opcode&roxl_roxrMask == roxl_roxrPrefix {
			result = instructions.CreateROXL_ROXR(opcode, code, offset)
		} else if opcode&lsl_lsrMask == lsl_lsrPrefix {
			result = instructions.CreateLSL_LSR(opcode, code, offset)
		}
	case 0b0000000000001111:
		// Coprocessor Interface/MC68040 and CPU32 Extensions

		// // TODO - don't think this is right
		// if opcode&dcPrefix == dcPrefix {
		// 	result = instructions.CreateDC(opcode)
		// }
	}

	if result != nil && !reflect.ValueOf(result).IsNil() {
		return result
	}

	// Ilegal instructions are called dc.w for some reason
	return instructions.CreateDC(opcode)
}
