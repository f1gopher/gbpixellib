package m68000

import "testing"

func TestDisassembleABCD(t *testing.T) {
	testCases := []testCase{
		// Register to Register mode (Rx,Ry)
		// Format: 1100xxx10000yyyy where xxx is destination register, yyyy is source register
		{[]uint8{0xC1, 0x00}, "ABCD D0,D0"},
		{[]uint8{0xC1, 0x01}, "ABCD D1,D0"},
		{[]uint8{0xC1, 0x02}, "ABCD D2,D0"},
		{[]uint8{0xC1, 0x03}, "ABCD D3,D0"},
		{[]uint8{0xC1, 0x04}, "ABCD D4,D0"},
		{[]uint8{0xC1, 0x05}, "ABCD D5,D0"},
		{[]uint8{0xC1, 0x06}, "ABCD D6,D0"},
		{[]uint8{0xC1, 0x07}, "ABCD D7,D0"},

		{[]uint8{0xC3, 0x00}, "ABCD D0,D1"},
		{[]uint8{0xC5, 0x00}, "ABCD D0,D2"},
		{[]uint8{0xC7, 0x00}, "ABCD D0,D3"},
		{[]uint8{0xC9, 0x00}, "ABCD D0,D4"},
		{[]uint8{0xCB, 0x00}, "ABCD D0,D5"},
		{[]uint8{0xCD, 0x00}, "ABCD D0,D6"},
		{[]uint8{0xCF, 0x00}, "ABCD D0,D7"},

		// All possible register to register combinations
		{[]uint8{0xCD, 0x07}, "ABCD D7,D6"},
		{[]uint8{0xC9, 0x03}, "ABCD D3,D4"},
		{[]uint8{0xC5, 0x01}, "ABCD D1,D2"},
		{[]uint8{0xC7, 0x05}, "ABCD D5,D3"},

		// Memory to Memory mode (-(Ax),-(Ay))
		// Format: 1100xxx10000yyyy where xxx is destination register, yyyy is source register + 8
		{[]uint8{0xC1, 0x08}, "ABCD -(A0),-(A0)"},
		{[]uint8{0xC1, 0x09}, "ABCD -(A1),-(A0)"},
		{[]uint8{0xC1, 0x0A}, "ABCD -(A2),-(A0)"},
		{[]uint8{0xC1, 0x0B}, "ABCD -(A3),-(A0)"},
		{[]uint8{0xC1, 0x0C}, "ABCD -(A4),-(A0)"},
		{[]uint8{0xC1, 0x0D}, "ABCD -(A5),-(A0)"},
		{[]uint8{0xC1, 0x0E}, "ABCD -(A6),-(A0)"},
		{[]uint8{0xC1, 0x0F}, "ABCD -(A7),-(A0)"},

		{[]uint8{0xC3, 0x08}, "ABCD -(A0),-(A1)"},
		{[]uint8{0xC5, 0x08}, "ABCD -(A0),-(A2)"},
		{[]uint8{0xC7, 0x08}, "ABCD -(A0),-(A3)"},
		{[]uint8{0xC9, 0x08}, "ABCD -(A0),-(A4)"},
		{[]uint8{0xCB, 0x08}, "ABCD -(A0),-(A5)"},
		{[]uint8{0xCD, 0x08}, "ABCD -(A0),-(A6)"},
		{[]uint8{0xCF, 0x08}, "ABCD -(A0),-(A7)"},

		// All possible memory to memory combinations
		{[]uint8{0xCD, 0x0F}, "ABCD -(A7),-(A6)"},
		{[]uint8{0xC9, 0x0B}, "ABCD -(A3),-(A4)"},
		{[]uint8{0xC5, 0x09}, "ABCD -(A1),-(A2)"},
		{[]uint8{0xC7, 0x0D}, "ABCD -(A5),-(A3)"},
	}

	runTests(t, testCases)
}
