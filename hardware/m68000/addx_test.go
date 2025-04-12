package m68000

import "testing"

func TestDisassembleADDX(t *testing.T) {
	testCases := []struct {
		bytes    []uint8
		expected string
	}{
		// ADDX - Add Extended
		// Operation: Source + Destination + X → Destination
		// Sizes: Byte (00), Word (01), Long (10)
		// Modes: Register to Register (Dx,Dy) and Memory to Memory (-(Ax),-(Ay))

		// Register to Register mode (Rx,Ry)
		// ADDX.B Dx,Dy - format: 1101yyy1xx000xxx where yyy is destination, xxx is source
		{[]uint8{0xD1, 0x00}, "ADDX.B D0,D0"},
		{[]uint8{0xD1, 0x01}, "ADDX.B D1,D0"},
		{[]uint8{0xD1, 0x02}, "ADDX.B D2,D0"},
		{[]uint8{0xD1, 0x03}, "ADDX.B D3,D0"},
		{[]uint8{0xD1, 0x04}, "ADDX.B D4,D0"},
		{[]uint8{0xD1, 0x05}, "ADDX.B D5,D0"},
		{[]uint8{0xD1, 0x06}, "ADDX.B D6,D0"},
		{[]uint8{0xD1, 0x07}, "ADDX.B D7,D0"},

		{[]uint8{0xD3, 0x00}, "ADDX.B D0,D1"},
		{[]uint8{0xD5, 0x00}, "ADDX.B D0,D2"},
		{[]uint8{0xD7, 0x00}, "ADDX.B D0,D3"},
		{[]uint8{0xD9, 0x00}, "ADDX.B D0,D4"},
		{[]uint8{0xDB, 0x00}, "ADDX.B D0,D5"},
		{[]uint8{0xDD, 0x00}, "ADDX.B D0,D6"},
		{[]uint8{0xDF, 0x00}, "ADDX.B D0,D7"},

		// Additional register to register combinations (byte)
		{[]uint8{0xD3, 0x04}, "ADDX.B D4,D1"},
		{[]uint8{0xD9, 0x06}, "ADDX.B D6,D4"},
		{[]uint8{0xDF, 0x02}, "ADDX.B D2,D7"},

		// ADDX.W Dx,Dy
		{[]uint8{0xD1, 0x40}, "ADDX.W D0,D0"},
		{[]uint8{0xD1, 0x47}, "ADDX.W D7,D0"},
		{[]uint8{0xD3, 0x45}, "ADDX.W D5,D1"},
		{[]uint8{0xD7, 0x46}, "ADDX.W D6,D3"},
		{[]uint8{0xDF, 0x40}, "ADDX.W D0,D7"},

		// ADDX.L Dx,Dy
		{[]uint8{0xD1, 0x80}, "ADDX.L D0,D0"},
		{[]uint8{0xD1, 0x87}, "ADDX.L D7,D0"},
		{[]uint8{0xD5, 0x81}, "ADDX.L D1,D2"},
		{[]uint8{0xD9, 0x82}, "ADDX.L D2,D4"},
		{[]uint8{0xDF, 0x83}, "ADDX.L D3,D7"},

		// Memory to Memory mode (-(Ax),-(Ay))
		// ADDX.B -(Ax),-(Ay) - format: 1101yyy1xx001xxx where yyy is destination, xxx is source
		{[]uint8{0xD1, 0x08}, "ADDX.B -(A0),-(A0)"},
		{[]uint8{0xD1, 0x09}, "ADDX.B -(A1),-(A0)"},
		{[]uint8{0xD1, 0x0A}, "ADDX.B -(A2),-(A0)"},
		{[]uint8{0xD1, 0x0B}, "ADDX.B -(A3),-(A0)"},
		{[]uint8{0xD1, 0x0C}, "ADDX.B -(A4),-(A0)"},
		{[]uint8{0xD1, 0x0D}, "ADDX.B -(A5),-(A0)"},
		{[]uint8{0xD1, 0x0E}, "ADDX.B -(A6),-(A0)"},
		{[]uint8{0xD1, 0x0F}, "ADDX.B -(A7),-(A0)"},

		{[]uint8{0xD3, 0x08}, "ADDX.B -(A0),-(A1)"},
		{[]uint8{0xD5, 0x08}, "ADDX.B -(A0),-(A2)"},
		{[]uint8{0xD7, 0x08}, "ADDX.B -(A0),-(A3)"},
		{[]uint8{0xD9, 0x08}, "ADDX.B -(A0),-(A4)"},
		{[]uint8{0xDB, 0x08}, "ADDX.B -(A0),-(A5)"},
		{[]uint8{0xDD, 0x08}, "ADDX.B -(A0),-(A6)"},
		{[]uint8{0xDF, 0x08}, "ADDX.B -(A0),-(A7)"},

		// Additional memory to memory combinations (byte)
		{[]uint8{0xD3, 0x0C}, "ADDX.B -(A4),-(A1)"},
		{[]uint8{0xD9, 0x0E}, "ADDX.B -(A6),-(A4)"},
		{[]uint8{0xDF, 0x0A}, "ADDX.B -(A2),-(A7)"},

		// ADDX.W -(Ax),-(Ay)
		{[]uint8{0xD1, 0x48}, "ADDX.W -(A0),-(A0)"},
		{[]uint8{0xD1, 0x4F}, "ADDX.W -(A7),-(A0)"},
		{[]uint8{0xD3, 0x4D}, "ADDX.W -(A5),-(A1)"},
		{[]uint8{0xD7, 0x4E}, "ADDX.W -(A6),-(A3)"},
		{[]uint8{0xDF, 0x48}, "ADDX.W -(A0),-(A7)"},

		// ADDX.L -(Ax),-(Ay)
		{[]uint8{0xD1, 0x88}, "ADDX.L -(A0),-(A0)"},
		{[]uint8{0xD1, 0x8F}, "ADDX.L -(A7),-(A0)"},
		{[]uint8{0xD5, 0x89}, "ADDX.L -(A1),-(A2)"},
		{[]uint8{0xD9, 0x8A}, "ADDX.L -(A2),-(A4)"},
		{[]uint8{0xDF, 0x8B}, "ADDX.L -(A3),-(A7)"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, 0)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}
