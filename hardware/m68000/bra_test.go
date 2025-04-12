package m68000

import "testing"

func TestDisassembleBRA(t *testing.T) {
	testCases := []testCaseWithAddress{
		// BRA - Branch Always (Unconditional)
		// Format: 0110 0000 xxxx xxxx (8-bit displacement)
		//         0110 0000 0000 0000 xxxx xxxx xxxx xxxx (16-bit displacement)
		// BRA is an unconditional branch that adds the sign-extended displacement to the PC

		// BRA with 8-bit displacement (short form)
		{[]uint8{0x60, 0x10}, "BRA.S $0012", 0x0000}, // Branch forward by 16 bytes
		{[]uint8{0x60, 0xF0}, "BRA.S $FFF2", 0x0000}, // Branch backward by 16 bytes
		{[]uint8{0x60, 0x7F}, "BRA.S $0081", 0x0000}, // Maximum forward branch (127 bytes)
		{[]uint8{0x60, 0x80}, "BRA.S $FF82", 0x0000}, // Maximum backward branch (-128 bytes)
		{[]uint8{0x60, 0x00}, "BRA.S $0002", 0x0000}, // Displacement of 0 (branch to PC+2)
		{[]uint8{0x60, 0xFE}, "BRA.S $0000", 0x0000}, // Displacement of -2 (branch to PC)

		// BRA with 16-bit displacement (word form)
		{[]uint8{0x60, 0x00, 0x01, 0x00}, "BRA.W $0102", 0x0000}, // Forward branch by 256 bytes
		{[]uint8{0x60, 0x00, 0xFF, 0x00}, "BRA.W $FF02", 0x0000}, // Backward branch by 256 bytes
		{[]uint8{0x60, 0x00, 0x7F, 0xFF}, "BRA.W $8001", 0x0000}, // Large forward branch
		{[]uint8{0x60, 0x00, 0x80, 0x00}, "BRA.W $8002", 0x0000}, // Large backward branch
		{[]uint8{0x60, 0x00, 0x00, 0x00}, "BRA.W $0002", 0x0000}, // Zero displacement (branch to PC+2)
		{[]uint8{0x60, 0x00, 0xFF, 0xFE}, "BRA.W $0000", 0x0000}, // Displacement of -2 (branch to PC)

		// BRA from different PC addresses
		{[]uint8{0x60, 0x10}, "BRA.S $1012", 0x1000},             // Short forward from address 0x1000
		{[]uint8{0x60, 0xF0}, "BRA.S $0FF2", 0x1000},             // Short backward from address 0x1000
		{[]uint8{0x60, 0x00, 0x01, 0x00}, "BRA.W $1102", 0x1000}, // Word forward from address 0x1000
		{[]uint8{0x60, 0x00, 0xFF, 0x00}, "BRA.W $0F02", 0x1000}, // Word backward from address 0x1000

		// Special cases
		{[]uint8{0x60, 0x02}, "BRA.S $0006", 0x0002},             // Branch to PC+6 (from address 0x0002)
		{[]uint8{0x60, 0x00, 0x00, 0x02}, "BRA.W $0006", 0x0002}, // Branch to PC+6 (from address 0x0002)

		// Very large displacements
		{[]uint8{0x60, 0x00, 0x7F, 0xFF}, "BRA.W $8001", 0x0000}, // +32767 bytes (max positive)
		{[]uint8{0x60, 0x00, 0x80, 0x00}, "BRA.W $8002", 0x0000}, // -32768 bytes (max negative)

		// Branch to common locations
		{[]uint8{0x60, 0xFE}, "BRA.S $1000", 0x1000},             // Branch to self
		{[]uint8{0x60, 0x00, 0xFF, 0xFE}, "BRA.W $1000", 0x1000}, // Branch to self (word form)
		{[]uint8{0x60, 0x06}, "BRA.S $000A", 0x0002},             // Branch to PC+8
		{[]uint8{0x60, 0xFA}, "BRA.S $0000", 0x0004},             // Branch to address 0x0000
	}

	runAddressTests(t, testCases)
}
