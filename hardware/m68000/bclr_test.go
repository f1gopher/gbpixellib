package m68000

import "testing"

func TestDisassembleBCLR(t *testing.T) {
	testCases := []struct {
		bytes    []uint8
		expected string
	}{
		// BCLR - Bit Clear
		// Format 1: BCLR Dn,<ea> (dynamic bit number from register)
		// Format: 0000 rrr 110 eee mmm (where rrr is the register number)

		// Format 2: BCLR #<data>,<ea> (static bit number from immediate)
		// Format: 0000 1000 10 eee mmm (followed by bit number as immediate word)

		// BCLR Dn,Dm (dynamic bit number, data register destination)
		{[]uint8{0x01, 0x80}, "BCLR D0,D0"},
		{[]uint8{0x03, 0x81}, "BCLR D1,D1"},
		{[]uint8{0x05, 0x82}, "BCLR D2,D2"},
		{[]uint8{0x07, 0x83}, "BCLR D3,D3"},
		{[]uint8{0x09, 0x84}, "BCLR D4,D4"},
		{[]uint8{0x0B, 0x85}, "BCLR D5,D5"},
		{[]uint8{0x0D, 0x86}, "BCLR D6,D6"},
		{[]uint8{0x0F, 0x87}, "BCLR D7,D7"},
		{[]uint8{0x01, 0x87}, "BCLR D0,D7"},
		{[]uint8{0x03, 0x86}, "BCLR D1,D6"},
		{[]uint8{0x05, 0x85}, "BCLR D2,D5"},
		{[]uint8{0x07, 0x84}, "BCLR D3,D4"},
		{[]uint8{0x09, 0x83}, "BCLR D4,D3"},
		{[]uint8{0x0B, 0x82}, "BCLR D5,D2"},
		{[]uint8{0x0D, 0x81}, "BCLR D6,D1"},
		{[]uint8{0x0F, 0x80}, "BCLR D7,D0"},

		// BCLR Dn,(Am) (dynamic bit number, memory destination)
		{[]uint8{0x01, 0x90}, "BCLR D0,(A0)"},
		{[]uint8{0x03, 0x91}, "BCLR D1,(A1)"},
		{[]uint8{0x05, 0x92}, "BCLR D2,(A2)"},
		{[]uint8{0x07, 0x93}, "BCLR D3,(A3)"},
		{[]uint8{0x09, 0x94}, "BCLR D4,(A4)"},
		{[]uint8{0x0B, 0x95}, "BCLR D5,(A5)"},
		{[]uint8{0x0D, 0x96}, "BCLR D6,(A6)"},
		{[]uint8{0x0F, 0x97}, "BCLR D7,(A7)"},

		// BCLR Dn,(Am)+ (dynamic bit number, memory destination with postincrement)
		{[]uint8{0x01, 0x98}, "BCLR D0,(A0)+"},
		{[]uint8{0x03, 0x99}, "BCLR D1,(A1)+"},
		{[]uint8{0x05, 0x9A}, "BCLR D2,(A2)+"},
		{[]uint8{0x07, 0x9B}, "BCLR D3,(A3)+"},
		{[]uint8{0x09, 0x9C}, "BCLR D4,(A4)+"},
		{[]uint8{0x0B, 0x9D}, "BCLR D5,(A5)+"},
		{[]uint8{0x0D, 0x9E}, "BCLR D6,(A6)+"},
		{[]uint8{0x0F, 0x9F}, "BCLR D7,(A7)+"},

		// BCLR Dn,-(Am) (dynamic bit number, memory destination with predecrement)
		{[]uint8{0x01, 0xA0}, "BCLR D0,-(A0)"},
		{[]uint8{0x03, 0xA1}, "BCLR D1,-(A1)"},
		{[]uint8{0x05, 0xA2}, "BCLR D2,-(A2)"},
		{[]uint8{0x07, 0xA3}, "BCLR D3,-(A3)"},
		{[]uint8{0x09, 0xA4}, "BCLR D4,-(A4)"},
		{[]uint8{0x0B, 0xA5}, "BCLR D5,-(A5)"},
		{[]uint8{0x0D, 0xA6}, "BCLR D6,-(A6)"},
		{[]uint8{0x0F, 0xA7}, "BCLR D7,-(A7)"},

		// BCLR Dn,d16(Am) (dynamic bit number, memory destination with displacement)
		{[]uint8{0x01, 0xA8, 0x12, 0x34}, "BCLR D0,(0x1234,A0)"},
		{[]uint8{0x03, 0xA9, 0x56, 0x78}, "BCLR D1,(0x5678,A1)"},
		{[]uint8{0x05, 0xAA, 0x9A, 0xBC}, "BCLR D2,(0x9ABC,A2)"},
		{[]uint8{0x07, 0xAB, 0xDE, 0xF0}, "BCLR D3,(0xDEF0,A3)"},
		{[]uint8{0x09, 0xAC, 0x11, 0x22}, "BCLR D4,(0x1122,A4)"},
		{[]uint8{0x0B, 0xAD, 0x33, 0x44}, "BCLR D5,(0x3344,A5)"},
		{[]uint8{0x0D, 0xAE, 0x55, 0x66}, "BCLR D6,(0x5566,A6)"},
		{[]uint8{0x0F, 0xAF, 0x77, 0x88}, "BCLR D7,(0x7788,A7)"},

		// BCLR Dn,d8(Am,Xn) (dynamic bit number, memory destination with index)
		{[]uint8{0x01, 0xB0, 0x12, 0x34}, "BCLR D0,(0x12,A0,D4.W)"},
		{[]uint8{0x03, 0xB1, 0x56, 0x78}, "BCLR D1,(0x56,A1,D7.W)"},
		{[]uint8{0x05, 0xB2, 0x9A, 0xB0}, "BCLR D2,(0x9A,A2,D0.W)"},
		{[]uint8{0x07, 0xB3, 0xDE, 0xF3}, "BCLR D3,(0xDE,A3,D3.W)"},
		{[]uint8{0x01, 0xB0, 0x12, 0xB4}, "BCLR D0,(0x12,A0,A4.W)"},
		{[]uint8{0x03, 0xB1, 0x56, 0xF7}, "BCLR D1,(0x56,A1,A7.W)"},

		// BCLR Dn,xxx.W (dynamic bit number, absolute short address)
		{[]uint8{0x01, 0xB8, 0x12, 0x34}, "BCLR D0,(0x1234).W"},
		{[]uint8{0x03, 0xB8, 0x56, 0x78}, "BCLR D1,(0x5678).W"},
		{[]uint8{0x05, 0xB8, 0x9A, 0xBC}, "BCLR D2,(0x9ABC).W"},

		// BCLR Dn,xxx.L (dynamic bit number, absolute long address)
		{[]uint8{0x01, 0xB9, 0x12, 0x34, 0x56, 0x78}, "BCLR D0,(0x12345678).L"},
		{[]uint8{0x03, 0xB9, 0x9A, 0xBC, 0xDE, 0xF0}, "BCLR D1,(0x9ABCDEF0).L"},

		// BCLR #imm,Dn (static bit number, data register destination)
		{[]uint8{0x08, 0x80, 0x00, 0x00}, "BCLR #0,D0"},
		{[]uint8{0x08, 0x81, 0x00, 0x07}, "BCLR #7,D1"},
		{[]uint8{0x08, 0x82, 0x00, 0x0F}, "BCLR #15,D2"},
		{[]uint8{0x08, 0x83, 0x00, 0x1F}, "BCLR #31,D3"},  // Maximum bit number for long word
		{[]uint8{0x08, 0x84, 0x00, 0x20}, "BCLR #32,D4"},  // Will be modulo 32
		{[]uint8{0x08, 0x85, 0x00, 0xFF}, "BCLR #255,D5"}, // Will be modulo 32

		// BCLR #imm,(Am) (static bit number, memory destination)
		{[]uint8{0x08, 0x90, 0x00, 0x00}, "BCLR #0,(A0)"},
		{[]uint8{0x08, 0x91, 0x00, 0x07}, "BCLR #7,(A1)"},
		{[]uint8{0x08, 0x92, 0x00, 0x08}, "BCLR #8,(A2)"},   // Will be modulo 8 for memory ops
		{[]uint8{0x08, 0x93, 0x00, 0xFF}, "BCLR #255,(A3)"}, // Will be modulo 8 for memory ops

		// BCLR #imm,(Am)+ (static bit number, memory destination with postincrement)
		{[]uint8{0x08, 0x98, 0x00, 0x00}, "BCLR #0,(A0)+"},
		{[]uint8{0x08, 0x99, 0x00, 0x07}, "BCLR #7,(A1)+"},

		// BCLR #imm,-(Am) (static bit number, memory destination with predecrement)
		{[]uint8{0x08, 0xA0, 0x00, 0x00}, "BCLR #0,-(A0)"},
		{[]uint8{0x08, 0xA1, 0x00, 0x07}, "BCLR #7,-(A1)"},

		// BCLR #imm,d16(Am) (static bit number, memory destination with displacement)
		{[]uint8{0x08, 0xA8, 0x12, 0x34, 0x00, 0x00}, "BCLR #0,(0x1234,A0)"},
		{[]uint8{0x08, 0xA9, 0x56, 0x78, 0x00, 0x07}, "BCLR #7,(0x5678,A1)"},

		// BCLR #imm,d8(Am,Xn) (static bit number, memory destination with index)
		{[]uint8{0x08, 0xB0, 0x12, 0x34, 0x00, 0x00}, "BCLR #0,(0x12,A0,D4.W)"},
		{[]uint8{0x08, 0xB1, 0x56, 0x78, 0x00, 0x07}, "BCLR #7,(0x56,A1,D7.W)"},
		{[]uint8{0x08, 0xB0, 0x12, 0xB4, 0x00, 0x00}, "BCLR #0,(0x12,A0,A4.W)"},

		// BCLR #imm,xxx.W (static bit number, absolute short address)
		{[]uint8{0x08, 0xB8, 0x12, 0x34, 0x00, 0x00}, "BCLR #0,(0x1234).W"},
		{[]uint8{0x08, 0xB8, 0x56, 0x78, 0x00, 0x07}, "BCLR #7,(0x5678).W"},

		// BCLR #imm,xxx.L (static bit number, absolute long address)
		{[]uint8{0x08, 0xB9, 0x12, 0x34, 0x56, 0x78, 0x00, 0x00}, "BCLR #0,(0x12345678).L"},
		{[]uint8{0x08, 0xB9, 0x9A, 0xBC, 0xDE, 0xF0, 0x00, 0x07}, "BCLR #7,(0x9ABCDEF0).L"},
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
