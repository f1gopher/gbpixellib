package m68000

import (
	"testing"
)

func DisassembleInstruction(code []uint8, address uint32) string {
	parser := instructionParser{}
	inst := parser.parse(code, int(address))
	return inst.Name()
}

func TestDisassembleDataMovement(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []uint8
		expected string
	}{
		// MOVE
		{"MOVE.B D0,D1", []uint8{0x12, 0x00}, "MOVE.B D0,D1"},
		{"MOVE.W D3,D4", []uint8{0x38, 0x03}, "MOVE.W D3,D4"},
		{"MOVE.L D5,D6", []uint8{0x2C, 0x05}, "MOVE.L D5,D6"},
		{"MOVE.W A0,D1", []uint8{0x32, 0x08}, "MOVE.W A0,D1"},
		{"MOVE.L (A0),D1", []uint8{0x22, 0x10}, "MOVE.L (A0),D1"},
		{"MOVE.W (A1)+,D2", []uint8{0x34, 0x19}, "MOVE.W (A1)+,D2"},
		{"MOVE.B -(A2),D3", []uint8{0x16, 0x22}, "MOVE.B -(A2),D3"},
		{"MOVE.L 4(A3),D4", []uint8{0x28, 0x2B, 0x00, 0x04}, "MOVE.L 4(A3),D4"},
		{"MOVE.W #$1234,D0", []uint8{0x30, 0x3C, 0x12, 0x34}, "MOVE.W #$1234,D0"},
		{"MOVE.L #$12345678,D1", []uint8{0x22, 0x3C, 0x12, 0x34, 0x56, 0x78}, "MOVE.L #$12345678,D1"},

		// MOVEA
		{"MOVEA.W D0,A1", []uint8{0x32, 0x40}, "MOVEA.W D0,A1"},
		{"MOVEA.L A2,A3", []uint8{0x26, 0x8A}, "MOVEA.L A2,A3"},
		{"MOVEA.W (A4),A5", []uint8{0x3A, 0x94}, "MOVEA.W (A4),A5"},

		// MOVEQ
		{"MOVEQ #0,D0", []uint8{0x70, 0x00}, "MOVEQ #0,D0"},
		{"MOVEQ #-1,D7", []uint8{0x7E, 0xFF}, "MOVEQ #-1,D7"},
		{"MOVEQ #127,D3", []uint8{0x76, 0x7F}, "MOVEQ #127,D3"},

		// MOVEM
		{"MOVEM.W D0-D7,-(A7)", []uint8{0x48, 0xA7, 0xFF, 0x00}, "MOVEM.W D0-D7,-(A7)"},
		{"MOVEM.L (A0)+,D0/D2/A1/A3", []uint8{0x4C, 0xD8, 0xA5, 0x05}, "MOVEM.L (A0)+,D0/D2/A1/A3"},

		// LEA
		{"LEA (A0),A1", []uint8{0x43, 0xD0}, "LEA (A0),A1"},
		{"LEA 8(A2),A3", []uint8{0x47, 0xEA, 0x00, 0x08}, "LEA 8(A2),A3"},

		// PEA
		{"PEA (A0)", []uint8{0x48, 0x50}, "PEA (A0)"},
		{"PEA 4(A1,D0.W)", []uint8{0x48, 0x71, 0x00, 0x04}, "PEA 4(A1,D0.W)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, 0)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestDisassembleIntegerArithmetic(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []uint8
		expected string
	}{
		// ADD
		{"ADD.B D0,D1", []uint8{0xD2, 0x00}, "ADD.B D0,D1"},
		{"ADD.W D2,D3", []uint8{0xD6, 0x42}, "ADD.W D2,D3"},
		{"ADD.L D4,D5", []uint8{0xDA, 0x84}, "ADD.L D4,D5"},
		{"ADD.W (A0),D1", []uint8{0xD2, 0x50}, "ADD.W (A0),D1"},
		{"ADD.L D5,(A2)", []uint8{0xDB, 0x92}, "ADD.L D5,(A2)"},

		// ADDA
		{"ADDA.W D0,A1", []uint8{0xD2, 0xC0}, "ADDA.W D0,A1"},
		{"ADDA.L (A2),A3", []uint8{0xD7, 0xD2}, "ADDA.L (A2),A3"},

		// ADDQ
		{"ADDQ.B #1,D0", []uint8{0x52, 0x00}, "ADDQ.B #1,D0"},
		{"ADDQ.W #8,A1", []uint8{0x50, 0x49}, "ADDQ.W #8,A1"},

		// SUB
		{"SUB.B D0,D1", []uint8{0x92, 0x00}, "SUB.B D0,D1"},
		{"SUB.W D2,D3", []uint8{0x96, 0x42}, "SUB.W D2,D3"},
		{"SUB.L D4,D5", []uint8{0x9A, 0x84}, "SUB.L D4,D5"},
		{"SUB.W (A0),D1", []uint8{0x92, 0x50}, "SUB.W (A0),D1"},
		{"SUB.L D5,(A2)", []uint8{0x9B, 0x92}, "SUB.L D5,(A2)"},

		// SUBA
		{"SUBA.W D0,A1", []uint8{0x92, 0xC0}, "SUBA.W D0,A1"},
		{"SUBA.L (A2),A3", []uint8{0x97, 0xD2}, "SUBA.L (A2),A3"},

		// SUBQ
		{"SUBQ.B #1,D0", []uint8{0x53, 0x00}, "SUBQ.B #1,D0"},
		{"SUBQ.W #8,A1", []uint8{0x51, 0x49}, "SUBQ.W #8,A1"},

		// MULU
		{"MULU.W D0,D1", []uint8{0xC2, 0xC0}, "MULU.W D0,D1"},
		{"MULU.W (A0),D2", []uint8{0xC4, 0xD0}, "MULU.W (A0),D2"},

		// MULS
		{"MULS.W D0,D1", []uint8{0xC3, 0xC0}, "MULS.W D0,D1"},
		{"MULS.W (A0),D2", []uint8{0xC5, 0xD0}, "MULS.W (A0),D2"},

		// DIVU
		{"DIVU.W D0,D1", []uint8{0x82, 0xC0}, "DIVU.W D0,D1"},
		{"DIVU.W (A0),D2", []uint8{0x84, 0xD0}, "DIVU.W (A0),D2"},

		// DIVS
		{"DIVS.W D0,D1", []uint8{0x83, 0xC0}, "DIVS.W D0,D1"},
		{"DIVS.W (A0),D2", []uint8{0x85, 0xD0}, "DIVS.W (A0),D2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, 0)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestDisassembleLogical(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []uint8
		expected string
	}{
		// AND
		{"AND.B D0,D1", []uint8{0xC2, 0x00}, "AND.B D0,D1"},
		{"AND.W D2,D3", []uint8{0xC6, 0x42}, "AND.W D2,D3"},
		{"AND.L D4,D5", []uint8{0xCA, 0x84}, "AND.L D4,D5"},
		{"AND.W (A0),D1", []uint8{0xC2, 0x50}, "AND.W (A0),D1"},
		{"AND.L D5,(A2)", []uint8{0xCB, 0x92}, "AND.L D5,(A2)"},

		// OR
		{"OR.B D0,D1", []uint8{0x82, 0x00}, "OR.B D0,D1"},
		{"OR.W D2,D3", []uint8{0x86, 0x42}, "OR.W D2,D3"},
		{"OR.L D4,D5", []uint8{0x8A, 0x84}, "OR.L D4,D5"},
		{"OR.W (A0),D1", []uint8{0x82, 0x50}, "OR.W (A0),D1"},
		{"OR.L D5,(A2)", []uint8{0x8B, 0x92}, "OR.L D5,(A2)"},

		// EOR
		{"EOR.B D0,D1", []uint8{0xB1, 0x00}, "EOR.B D0,D1"},
		{"EOR.W D2,D3", []uint8{0xB5, 0x42}, "EOR.W D2,D3"},
		{"EOR.L D4,D5", []uint8{0xB9, 0x84}, "EOR.L D4,D5"},
		{"EOR.W D1,(A0)", []uint8{0xB3, 0x50}, "EOR.W D1,(A0)"},

		// NOT
		{"NOT.B D0", []uint8{0x46, 0x00}, "NOT.B D0"},
		{"NOT.W D1", []uint8{0x46, 0x41}, "NOT.W D1"},
		{"NOT.L D2", []uint8{0x46, 0x82}, "NOT.L D2"},
		{"NOT.W (A0)", []uint8{0x46, 0x50}, "NOT.W (A0)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, 0)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestDisassembleShiftsRotates(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []uint8
		expected string
	}{
		// ASL/ASR register
		{"ASL.B D1,D0", []uint8{0xE3, 0x20}, "ASL.B D1,D0"},
		{"ASL.W #8,D2", []uint8{0xE1, 0x4A}, "ASL.W #8,D2"},
		{"ASR.L #1,D3", []uint8{0xE2, 0x83}, "ASR.L #1,D3"},
		{"ASR.W (A0)", []uint8{0xE0, 0xD0}, "ASR.W (A0)"},

		// LSL/LSR
		{"LSL.B D1,D0", []uint8{0xE3, 0x28}, "LSL.B D1,D0"},
		{"LSL.W #8,D2", []uint8{0xE1, 0x4A}, "LSL.W #8,D2"},
		{"LSR.L #1,D3", []uint8{0xE2, 0x8B}, "LSR.L #1,D3"},
		{"LSR.W (A0)", []uint8{0xE2, 0xD0}, "LSR.W (A0)"},

		// ROL/ROR
		{"ROL.B D1,D0", []uint8{0xE3, 0x38}, "ROL.B D1,D0"},
		{"ROL.W #8,D2", []uint8{0xE1, 0x5A}, "ROL.W #8,D2"},
		{"ROR.L #1,D3", []uint8{0xE2, 0x9B}, "ROR.L #1,D3"},
		{"ROR.W (A0)", []uint8{0xE6, 0xD0}, "ROR.W (A0)"},

		// ROXL/ROXR
		{"ROXL.B D1,D0", []uint8{0xE3, 0x30}, "ROXL.B D1,D0"},
		{"ROXL.W #8,D2", []uint8{0xE1, 0x52}, "ROXL.W #8,D2"},
		{"ROXR.L #1,D3", []uint8{0xE2, 0x93}, "ROXR.L #1,D3"},
		{"ROXR.W (A0)", []uint8{0xE4, 0xD0}, "ROXR.W (A0)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, 0)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestDisassembleBitManipulation(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []uint8
		expected string
	}{
		// BTST
		{"BTST D0,D1", []uint8{0x01, 0x01}, "BTST D0,D1"},
		{"BTST #7,D2", []uint8{0x08, 0x3C, 0x00, 0x07}, "BTST #7,D2"},
		{"BTST D3,(A4)", []uint8{0x07, 0x14}, "BTST D3,(A4)"},

		// BSET
		{"BSET D0,D1", []uint8{0x01, 0xC1}, "BSET D0,D1"},
		{"BSET #7,D2", []uint8{0x08, 0xFC, 0x00, 0x07}, "BSET #7,D2"},
		{"BSET D3,(A4)", []uint8{0x07, 0xD4}, "BSET D3,(A4)"},

		// BCLR
		{"BCLR D0,D1", []uint8{0x01, 0x81}, "BCLR D0,D1"},
		{"BCLR #7,D2", []uint8{0x08, 0xBC, 0x00, 0x07}, "BCLR #7,D2"},
		{"BCLR D3,(A4)", []uint8{0x07, 0x94}, "BCLR D3,(A4)"},

		// BCHG
		{"BCHG D0,D1", []uint8{0x01, 0x41}, "BCHG D0,D1"},
		{"BCHG #7,D2", []uint8{0x08, 0x7C, 0x00, 0x07}, "BCHG #7,D2"},
		{"BCHG D3,(A4)", []uint8{0x07, 0x54}, "BCHG D3,(A4)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, 0)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestDisassembleBCDOps(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []uint8
		expected string
	}{
		// ABCD
		{"ABCD D0,D1", []uint8{0xC3, 0x00}, "ABCD D0,D1"},
		{"ABCD -(A0),-(A1)", []uint8{0xC3, 0x08}, "ABCD -(A0),-(A1)"},

		// SBCD
		{"SBCD D0,D1", []uint8{0x83, 0x00}, "SBCD D0,D1"},
		{"SBCD -(A0),-(A1)", []uint8{0x83, 0x08}, "SBCD -(A0),-(A1)"},

		// NBCD
		{"NBCD D0", []uint8{0x48, 0x00}, "NBCD D0"},
		{"NBCD (A0)", []uint8{0x48, 0x10}, "NBCD (A0)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, 0)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestDisassembleProgramControl(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []uint8
		expected string
		address  uint32
	}{
		// Branches
		{"BRA.S +10", []uint8{0x60, 0x0A}, "BRA.S $000C", 0x0000},
		{"BEQ.S -8", []uint8{0x67, 0xF8}, "BEQ.S $0FF8", 0x1000},
		{"BRA.W +1000", []uint8{0x60, 0x00, 0x03, 0xE8}, "BRA.W $03EA", 0x0000},
		{"BGT.W -1000", []uint8{0x6E, 0x00, 0xFC, 0x18}, "BGT.W $EC1A", 0xF000},

		// JMP
		{"JMP (A0)", []uint8{0x4E, 0xD0}, "JMP (A0)", 0x0000},
		{"JMP 100(A0)", []uint8{0x4E, 0xE8, 0x00, 0x64}, "JMP 100(A0)", 0x0000},

		// JSR
		{"JSR (A0)", []uint8{0x4E, 0x90}, "JSR (A0)", 0x0000},
		{"JSR 100(A0)", []uint8{0x4E, 0xA8, 0x00, 0x64}, "JSR 100(A0)", 0x0000},

		// RTS
		{"RTS", []uint8{0x4E, 0x75}, "RTS", 0x0000},

		// RTR
		{"RTR", []uint8{0x4E, 0x77}, "RTR", 0x0000},

		// BSR
		{"BSR.S +10", []uint8{0x61, 0x0A}, "BSR.S $000C", 0x0000},
		{"BSR.W +1000", []uint8{0x61, 0x00, 0x03, 0xE8}, "BSR.W $03EA", 0x0000},

		// Scc
		{"SEQ D0", []uint8{0x57, 0xC0}, "SEQ D0", 0x0000},
		{"SNE (A0)", []uint8{0x56, 0xD0}, "SNE (A0)", 0x0000},

		// DBcc
		{"DBEQ D0,-10", []uint8{0x57, 0xC8, 0xFF, 0xF6}, "DBEQ D0,$FFF8", 0x0000},
		{"DBNE D1,+20", []uint8{0x56, 0xC9, 0x00, 0x14}, "DBNE D1,$0016", 0x0000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, tc.address)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestDisassembleSystemControl(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []uint8
		expected string
	}{
		// TRAP
		{"TRAP #0", []uint8{0x4E, 0x40}, "TRAP #0"},
		{"TRAP #15", []uint8{0x4E, 0x4F}, "TRAP #15"},

		// RESET
		{"RESET", []uint8{0x4E, 0x70}, "RESET"},

		// NOP
		{"NOP", []uint8{0x4E, 0x71}, "NOP"},

		// STOP
		{"STOP #$2000", []uint8{0x4E, 0x72, 0x20, 0x00}, "STOP #$2000"},

		// RTE
		{"RTE", []uint8{0x4E, 0x73}, "RTE"},

		// LINK
		{"LINK A0,#-4", []uint8{0x4E, 0x50, 0xFF, 0xFC}, "LINK A0,#-4"},

		// UNLK
		{"UNLK A0", []uint8{0x4E, 0x58}, "UNLK A0"},

		// MOVE to/from SR
		{"MOVE SR,D0", []uint8{0x40, 0xC0}, "MOVE SR,D0"},
		{"MOVE D0,SR", []uint8{0x46, 0xC0}, "MOVE D0,SR"},

		// MOVE to/from USP
		{"MOVE USP,A0", []uint8{0x4E, 0x68}, "MOVE USP,A0"},
		{"MOVE A0,USP", []uint8{0x4E, 0x60}, "MOVE A0,USP"},

		// CHK
		{"CHK D0,D1", []uint8{0x43, 0x00}, "CHK D0,D1"},
		{"CHK (A0),D2", []uint8{0x45, 0x10}, "CHK (A0),D2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, 0)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}
