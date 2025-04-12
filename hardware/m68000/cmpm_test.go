package m68000

import "testing"

func TestDisassembleCMPM(t *testing.T) {
	testCases := []testCase{
		// CMPM - Compare Memory
		// Format: 1011 rrr 1ss 001 mmm
		// Where:
		// - rrr is the second register number (Ay)
		// - ss is the size (00=byte, 01=word, 10=long)
		// - mmm is the first register number (Ax)
		//
		// This instruction compares memory pointed to by (Ax)+ with (Ay)+
		// A specialized form that auto-increments both address registers

		// CMPM.B (An)+,(Am)+ (Compare byte memory to memory with postincrement)
		{[]uint8{0xB1, 0x08}, "CMPM.B (A0)+,(A0)+"},
		{[]uint8{0xB1, 0x09}, "CMPM.B (A1)+,(A0)+"},
		{[]uint8{0xB1, 0x0A}, "CMPM.B (A2)+,(A0)+"},
		{[]uint8{0xB1, 0x0B}, "CMPM.B (A3)+,(A0)+"},
		{[]uint8{0xB1, 0x0C}, "CMPM.B (A4)+,(A0)+"},
		{[]uint8{0xB1, 0x0D}, "CMPM.B (A5)+,(A0)+"},
		{[]uint8{0xB1, 0x0E}, "CMPM.B (A6)+,(A0)+"},
		{[]uint8{0xB1, 0x0F}, "CMPM.B (A7)+,(A0)+"},

		{[]uint8{0xB3, 0x08}, "CMPM.B (A0)+,(A1)+"},
		{[]uint8{0xB5, 0x08}, "CMPM.B (A0)+,(A2)+"},
		{[]uint8{0xB7, 0x08}, "CMPM.B (A0)+,(A3)+"},
		{[]uint8{0xB9, 0x08}, "CMPM.B (A0)+,(A4)+"},
		{[]uint8{0xBB, 0x08}, "CMPM.B (A0)+,(A5)+"},
		{[]uint8{0xBD, 0x08}, "CMPM.B (A0)+,(A6)+"},
		{[]uint8{0xBF, 0x08}, "CMPM.B (A0)+,(A7)+"},

		// All combinations of registers for byte
		{[]uint8{0xB3, 0x09}, "CMPM.B (A1)+,(A1)+"},
		{[]uint8{0xB5, 0x0A}, "CMPM.B (A2)+,(A2)+"},
		{[]uint8{0xB7, 0x0B}, "CMPM.B (A3)+,(A3)+"},
		{[]uint8{0xB9, 0x0C}, "CMPM.B (A4)+,(A4)+"},
		{[]uint8{0xBB, 0x0D}, "CMPM.B (A5)+,(A5)+"},
		{[]uint8{0xBD, 0x0E}, "CMPM.B (A6)+,(A6)+"},
		{[]uint8{0xBF, 0x0F}, "CMPM.B (A7)+,(A7)+"},

		{[]uint8{0xB3, 0x0F}, "CMPM.B (A7)+,(A1)+"},
		{[]uint8{0xB5, 0x0D}, "CMPM.B (A5)+,(A2)+"},
		{[]uint8{0xB7, 0x0A}, "CMPM.B (A2)+,(A3)+"},
		{[]uint8{0xB9, 0x0E}, "CMPM.B (A6)+,(A4)+"},
		{[]uint8{0xBB, 0x0B}, "CMPM.B (A3)+,(A5)+"},
		{[]uint8{0xBD, 0x09}, "CMPM.B (A1)+,(A6)+"},
		{[]uint8{0xBF, 0x0C}, "CMPM.B (A4)+,(A7)+"},

		// CMPM.W (An)+,(Am)+ (Compare word memory to memory with postincrement)
		{[]uint8{0xB1, 0x48}, "CMPM.W (A0)+,(A0)+"},
		{[]uint8{0xB1, 0x49}, "CMPM.W (A1)+,(A0)+"},
		{[]uint8{0xB1, 0x4A}, "CMPM.W (A2)+,(A0)+"},
		{[]uint8{0xB1, 0x4F}, "CMPM.W (A7)+,(A0)+"},

		{[]uint8{0xB3, 0x48}, "CMPM.W (A0)+,(A1)+"},
		{[]uint8{0xB5, 0x48}, "CMPM.W (A0)+,(A2)+"},
		{[]uint8{0xB7, 0x48}, "CMPM.W (A0)+,(A3)+"},
		{[]uint8{0xB9, 0x48}, "CMPM.W (A0)+,(A4)+"},
		{[]uint8{0xBB, 0x48}, "CMPM.W (A0)+,(A5)+"},
		{[]uint8{0xBD, 0x48}, "CMPM.W (A0)+,(A6)+"},
		{[]uint8{0xBF, 0x48}, "CMPM.W (A0)+,(A7)+"},

		// Cross register tests for word
		{[]uint8{0xB3, 0x4F}, "CMPM.W (A7)+,(A1)+"},
		{[]uint8{0xB5, 0x4D}, "CMPM.W (A5)+,(A2)+"},
		{[]uint8{0xB7, 0x4A}, "CMPM.W (A2)+,(A3)+"},
		{[]uint8{0xBF, 0x4C}, "CMPM.W (A4)+,(A7)+"},

		// CMPM.L (An)+,(Am)+ (Compare long memory to memory with postincrement)
		{[]uint8{0xB1, 0x88}, "CMPM.L (A0)+,(A0)+"},
		{[]uint8{0xB1, 0x89}, "CMPM.L (A1)+,(A0)+"},
		{[]uint8{0xB1, 0x8A}, "CMPM.L (A2)+,(A0)+"},
		{[]uint8{0xB1, 0x8F}, "CMPM.L (A7)+,(A0)+"},

		{[]uint8{0xB3, 0x88}, "CMPM.L (A0)+,(A1)+"},
		{[]uint8{0xB5, 0x88}, "CMPM.L (A0)+,(A2)+"},
		{[]uint8{0xB7, 0x88}, "CMPM.L (A0)+,(A3)+"},
		{[]uint8{0xB9, 0x88}, "CMPM.L (A0)+,(A4)+"},
		{[]uint8{0xBB, 0x88}, "CMPM.L (A0)+,(A5)+"},
		{[]uint8{0xBD, 0x88}, "CMPM.L (A0)+,(A6)+"},
		{[]uint8{0xBF, 0x88}, "CMPM.L (A0)+,(A7)+"},

		// Cross register tests for long
		{[]uint8{0xB3, 0x8F}, "CMPM.L (A7)+,(A1)+"},
		{[]uint8{0xB5, 0x8D}, "CMPM.L (A5)+,(A2)+"},
		{[]uint8{0xB7, 0x8A}, "CMPM.L (A2)+,(A3)+"},
		{[]uint8{0xBF, 0x8C}, "CMPM.L (A4)+,(A7)+"},
	}

	runTests(t, testCases)
}
