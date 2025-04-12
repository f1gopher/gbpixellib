package m68000

import "testing"

func TestDisassembleBcc(t *testing.T) {
	testCases := []struct {
		bytes    []uint8
		expected string
		address  uint32 // PC address for relative branch calculation
	}{
		// Bcc - Branch Conditionally
		// Format: 0110 CCCC xxxx xxxx (8-bit displacement)
		//         0110 CCCC 0000 0000 xxxx xxxx xxxx xxxx (16-bit displacement)
		// Where CCCC is the condition code:
		// 0000 = BRA (Branch Always)
		// 0001 = BSR (Branch to Subroutine)
		// 0010 = BHI (Branch if Higher)
		// 0011 = BLS (Branch if Lower or Same)
		// 0100 = BCC/BHS (Branch if Carry Clear/Higher or Same)
		// 0101 = BCS/BLO (Branch if Carry Set/Lower)
		// 0110 = BNE (Branch if Not Equal)
		// 0111 = BEQ (Branch if Equal)
		// 1000 = BVC (Branch if Overflow Clear)
		// 1001 = BVS (Branch if Overflow Set)
		// 1010 = BPL (Branch if Plus)
		// 1011 = BMI (Branch if Minus)
		// 1100 = BGE (Branch if Greater or Equal)
		// 1101 = BLT (Branch if Less Than)
		// 1110 = BGT (Branch if Greater Than)
		// 1111 = BLE (Branch if Less or Equal)

		// BRA - Branch Always (8-bit displacement)
		{[]uint8{0x60, 0x10}, "BRA.S $0012", 0x0000}, // Branch forward
		{[]uint8{0x60, 0xF0}, "BRA.S $FFF2", 0x0000}, // Branch backward
		{[]uint8{0x60, 0x7F}, "BRA.S $0081", 0x0000}, // Maximum forward (127)
		{[]uint8{0x60, 0x80}, "BRA.S $FF82", 0x0000}, // Maximum backward (-128)
		{[]uint8{0x60, 0x00}, "BRA.S $0002", 0x0000}, // Displacement of 0 means +2

		// BRA - Branch Always (16-bit displacement)
		{[]uint8{0x60, 0x00, 0x01, 0x00}, "BRA.W $0102", 0x0000}, // Branch forward
		{[]uint8{0x60, 0x00, 0xFF, 0x00}, "BRA.W $FF02", 0x0000}, // Branch backward
		{[]uint8{0x60, 0x00, 0x7F, 0xFF}, "BRA.W $8001", 0x0000}, // Large forward
		{[]uint8{0x60, 0x00, 0x80, 0x00}, "BRA.W $8002", 0x0000}, // Large backward

		// BSR - Branch to Subroutine (8-bit displacement)
		{[]uint8{0x61, 0x10}, "BSR.S $0012", 0x0000}, // Branch forward
		{[]uint8{0x61, 0xF0}, "BSR.S $FFF2", 0x0000}, // Branch backward

		// BSR - Branch to Subroutine (16-bit displacement)
		{[]uint8{0x61, 0x00, 0x01, 0x00}, "BSR.W $0102", 0x0000}, // Branch forward
		{[]uint8{0x61, 0x00, 0xFF, 0x00}, "BSR.W $FF02", 0x0000}, // Branch backward

		// Bcc - Various condition codes (8-bit displacement)
		{[]uint8{0x62, 0x10}, "BHI.S $0012", 0x0000}, // BHI - Branch if Higher
		{[]uint8{0x63, 0x10}, "BLS.S $0012", 0x0000}, // BLS - Branch if Lower or Same
		{[]uint8{0x64, 0x10}, "BCC.S $0012", 0x0000}, // BCC - Branch if Carry Clear
		{[]uint8{0x65, 0x10}, "BCS.S $0012", 0x0000}, // BCS - Branch if Carry Set
		{[]uint8{0x66, 0x10}, "BNE.S $0012", 0x0000}, // BNE - Branch if Not Equal
		{[]uint8{0x67, 0x10}, "BEQ.S $0012", 0x0000}, // BEQ - Branch if Equal
		{[]uint8{0x68, 0x10}, "BVC.S $0012", 0x0000}, // BVC - Branch if Overflow Clear
		{[]uint8{0x69, 0x10}, "BVS.S $0012", 0x0000}, // BVS - Branch if Overflow Set
		{[]uint8{0x6A, 0x10}, "BPL.S $0012", 0x0000}, // BPL - Branch if Plus
		{[]uint8{0x6B, 0x10}, "BMI.S $0012", 0x0000}, // BMI - Branch if Minus
		{[]uint8{0x6C, 0x10}, "BGE.S $0012", 0x0000}, // BGE - Branch if Greater or Equal
		{[]uint8{0x6D, 0x10}, "BLT.S $0012", 0x0000}, // BLT - Branch if Less Than
		{[]uint8{0x6E, 0x10}, "BGT.S $0012", 0x0000}, // BGT - Branch if Greater Than
		{[]uint8{0x6F, 0x10}, "BLE.S $0012", 0x0000}, // BLE - Branch if Less or Equal

		// Bcc - Various condition codes (16-bit displacement)
		{[]uint8{0x62, 0x00, 0x01, 0x00}, "BHI.W $0102", 0x0000}, // BHI - Branch if Higher
		{[]uint8{0x63, 0x00, 0x01, 0x00}, "BLS.W $0102", 0x0000}, // BLS - Branch if Lower or Same
		{[]uint8{0x64, 0x00, 0x01, 0x00}, "BCC.W $0102", 0x0000}, // BCC - Branch if Carry Clear
		{[]uint8{0x65, 0x00, 0x01, 0x00}, "BCS.W $0102", 0x0000}, // BCS - Branch if Carry Set
		{[]uint8{0x66, 0x00, 0x01, 0x00}, "BNE.W $0102", 0x0000}, // BNE - Branch if Not Equal
		{[]uint8{0x67, 0x00, 0x01, 0x00}, "BEQ.W $0102", 0x0000}, // BEQ - Branch if Equal
		{[]uint8{0x68, 0x00, 0x01, 0x00}, "BVC.W $0102", 0x0000}, // BVC - Branch if Overflow Clear
		{[]uint8{0x69, 0x00, 0x01, 0x00}, "BVS.W $0102", 0x0000}, // BVS - Branch if Overflow Set
		{[]uint8{0x6A, 0x00, 0x01, 0x00}, "BPL.W $0102", 0x0000}, // BPL - Branch if Plus
		{[]uint8{0x6B, 0x00, 0x01, 0x00}, "BMI.W $0102", 0x0000}, // BMI - Branch if Minus
		{[]uint8{0x6C, 0x00, 0x01, 0x00}, "BGE.W $0102", 0x0000}, // BGE - Branch if Greater or Equal
		{[]uint8{0x6D, 0x00, 0x01, 0x00}, "BLT.W $0102", 0x0000}, // BLT - Branch if Less Than
		{[]uint8{0x6E, 0x00, 0x01, 0x00}, "BGT.W $0102", 0x0000}, // BGT - Branch if Greater Than
		{[]uint8{0x6F, 0x00, 0x01, 0x00}, "BLE.W $0102", 0x0000}, // BLE - Branch if Less or Equal

		// Test with different PC addresses
		{[]uint8{0x60, 0x10}, "BRA.S $1012", 0x1000},             // BRA from address 0x1000
		{[]uint8{0x61, 0xF0}, "BSR.S $0FF0", 0x1000},             // BSR from address 0x1000
		{[]uint8{0x66, 0x7F}, "BNE.S $107F", 0x1000},             // BNE from address 0x1000
		{[]uint8{0x67, 0x80}, "BEQ.S $0F80", 0x1000},             // BEQ from address 0x1000
		{[]uint8{0x60, 0x00, 0x01, 0x00}, "BRA.W $1102", 0x1000}, // BRA.W from address 0x1000
		{[]uint8{0x61, 0x00, 0xFE, 0x00}, "BSR.W $0E02", 0x1000}, // BSR.W from address 0x1000

		// Special cases
		{[]uint8{0x60, 0xFE}, "BRA.S $0000", 0x0000},             // BRA.S -2 (branch to self)
		{[]uint8{0x60, 0x00, 0x00, 0x00}, "BRA.W $0002", 0x0000}, // BRA.W with zero offset
		{[]uint8{0x60, 0x00, 0xFF, 0xFE}, "BRA.W $0000", 0x0000}, // BRA.W -2 (branch to self)
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, tc.address)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s (at address 0x%X)", tc.expected, result, tc.address)
			}
		})
	}
}
