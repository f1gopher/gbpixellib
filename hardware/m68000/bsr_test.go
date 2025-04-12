package m68000

import "testing"

func TestDisassembleBSR(t *testing.T) {
	testCases := []testCaseWithAddress{
		// BSR - Branch to Subroutine
		// Format: 0110 0001 xxxx xxxx (8-bit displacement)
		//         0110 0001 0000 0000 xxxx xxxx xxxx xxxx (16-bit displacement)
		// BSR pushes the return address (the instruction following the BSR) onto the stack
		// and then branches to the target address

		// BSR with 8-bit displacement (short form)
		{[]uint8{0x61, 0x10}, "BSR.S $0012", 0x0000}, // Branch forward by 16 bytes
		{[]uint8{0x61, 0xF0}, "BSR.S $FFF2", 0x0000}, // Branch backward by 16 bytes
		{[]uint8{0x61, 0x7F}, "BSR.S $0081", 0x0000}, // Maximum forward branch (127 bytes)
		{[]uint8{0x61, 0x80}, "BSR.S $FF82", 0x0000}, // Maximum backward branch (-128 bytes)
		{[]uint8{0x61, 0x00}, "BSR.S $0002", 0x0000}, // Displacement of 0 (branch to PC+2)
		{[]uint8{0x61, 0xFE}, "BSR.S $0000", 0x0000}, // Displacement of -2 (branch to PC)

		// BSR with 16-bit displacement (word form)
		{[]uint8{0x61, 0x00, 0x01, 0x00}, "BSR.W $0102", 0x0000}, // Forward branch by 256 bytes
		{[]uint8{0x61, 0x00, 0xFF, 0x00}, "BSR.W $FF02", 0x0000}, // Backward branch by 256 bytes
		{[]uint8{0x61, 0x00, 0x7F, 0xFF}, "BSR.W $8001", 0x0000}, // Large forward branch
		{[]uint8{0x61, 0x00, 0x80, 0x00}, "BSR.W $8002", 0x0000}, // Large backward branch
		{[]uint8{0x61, 0x00, 0x00, 0x00}, "BSR.W $0002", 0x0000}, // Zero displacement (branch to PC+2)
		{[]uint8{0x61, 0x00, 0xFF, 0xFE}, "BSR.W $0000", 0x0000}, // Displacement of -2 (branch to PC)

		// BSR from different PC addresses
		{[]uint8{0x61, 0x10}, "BSR.S $1012", 0x1000},             // Short forward from address 0x1000
		{[]uint8{0x61, 0xF0}, "BSR.S $0FF2", 0x1000},             // Short backward from address 0x1000
		{[]uint8{0x61, 0x00, 0x01, 0x00}, "BSR.W $1102", 0x1000}, // Word forward from address 0x1000
		{[]uint8{0x61, 0x00, 0xFF, 0x00}, "BSR.W $0F02", 0x1000}, // Word backward from address 0x1000

		// Special cases
		{[]uint8{0x61, 0x02}, "BSR.S $0006", 0x0002},             // Branch to PC+6 (from address 0x0002)
		{[]uint8{0x61, 0x00, 0x00, 0x02}, "BSR.W $0006", 0x0002}, // Branch to PC+6 (from address 0x0002)

		// Very large displacements
		{[]uint8{0x61, 0x00, 0x7F, 0xFF}, "BSR.W $8001", 0x0000}, // +32767 bytes (max positive)
		{[]uint8{0x61, 0x00, 0x80, 0x00}, "BSR.W $8002", 0x0000}, // -32768 bytes (max negative)

		// Common subroutine targets
		{[]uint8{0x61, 0xFE}, "BSR.S $1000", 0x1000},             // Call to self
		{[]uint8{0x61, 0x00, 0xFF, 0xFE}, "BSR.W $1000", 0x1000}, // Call to self (word form)
		{[]uint8{0x61, 0x06}, "BSR.S $000A", 0x0002},             // Call to PC+8
		{[]uint8{0x61, 0xFA}, "BSR.S $0000", 0x0004},             // Call to address 0x0000
	}

	runAddressTests(t, testCases)
}
