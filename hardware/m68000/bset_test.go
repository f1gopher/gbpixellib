package m68000

import "testing"

func TestDisassembleBSET(t *testing.T) {
	testCases := []struct {
		bytes    []uint8
		expected string
	}{
		// BSET - Bit Set
		// Format 1: BSET Dn,<ea> (dynamic bit number from register)
		// Format: 0000 rrr 111 eee mmm (where rrr is the register number)

		// Format 2: BSET #<data>,<ea> (static bit number from immediate)
		// Format: 0000 1000 11 eee mmm (followed by bit number as immediate word)

		// BSET Dn,Dm (dynamic bit number, data register destination)
		{[]uint8{0x01, 0xC0}, "BSET D0,D0"},
		{[]uint8{0x03, 0xC1}, "BSET D1,D1"},
		{[]uint8{0x05, 0xC2}, "BSET D2,D2"},
		{[]uint8{0x07, 0xC3}, "BSET D3,D3"},
		{[]uint8{0x09, 0xC4}, "BSET D4,D4"},
		{[]uint8{0x0B, 0xC5}, "BSET D5,D5"},
		{[]uint8{0x0D, 0xC6}, "BSET D6,D6"},
		{[]uint8{0x0F, 0xC7}, "BSET D7,D7"},
		{[]uint8{0x01, 0xC7}, "BSET D0,D7"},
		{[]uint8{0x03, 0xC6}, "BSET D1,D6"},
		{[]uint8{0x05, 0xC5}, "BSET D2,D5"},
		{[]uint8{0x07, 0xC4}, "BSET D3,D4"},
		{[]uint8{0x09, 0xC3}, "BSET D4,D3"},
		{[]uint8{0x0B, 0xC2}, "BSET D5,D2"},
		{[]uint8{0x0D, 0xC1}, "BSET D6,D1"},
		{[]uint8{0x0F, 0xC0}, "BSET D7,D0"},

		// BSET Dn,(Am) (dynamic bit number, memory destination)
		{[]uint8{0x01, 0xD0}, "BSET D0,(A0)"},
		{[]uint8{0x03, 0xD1}, "BSET D1,(A1)"},
		{[]uint8{0x05, 0xD2}, "BSET D2,(A2)"},
		{[]uint8{0x07, 0xD3}, "BSET D3,(A3)"},
		{[]uint8{0x09, 0xD4}, "BSET D4,(A4)"},
		{[]uint8{0x0B, 0xD5}, "BSET D5,(A5)"},
		{[]uint8{0x0D, 0xD6}, "BSET D6,(A6)"},
		{[]uint8{0x0F, 0xD7}, "BSET D7,(A7)"},

		// BSET Dn,(Am)+ (dynamic bit number, memory destination with postincrement)
		{[]uint8{0x01, 0xD8}, "BSET D0,(A0)+"},
		{[]uint8{0x03, 0xD9}, "BSET D1,(A1)+"},
		{[]uint8{0x05, 0xDA}, "BSET D2,(A2)+"},
		{[]uint8{0x07, 0xDB}, "BSET D3,(A3)+"},
		{[]uint8{0x09, 0xDC}, "BSET D4,(A4)+"},
		{[]uint8{0x0B, 0xDD}, "BSET D5,(A5)+"},
		{[]uint8{0x0D, 0xDE}, "BSET D6,(A6)+"},
		{[]uint8{0x0F, 0xDF}, "BSET D7,(A7)+"},

		// BSET Dn,-(Am) (dynamic bit number, memory destination with predecrement)
		{[]uint8{0x01, 0xE0}, "BSET D0,-(A0)"},
		{[]uint8{0x03, 0xE1}, "BSET D1,-(A1)"},
		{[]uint8{0x05, 0xE2}, "BSET D2,-(A2)"},
		{[]uint8{0x07, 0xE3}, "BSET D3,-(A3)"},
		{[]uint8{0x09, 0xE4}, "BSET D4,-(A4)"},
		{[]uint8{0x0B, 0xE5}, "BSET D5,-(A5)"},
		{[]uint8{0x0D, 0xE6}, "BSET D6,-(A6)"},
		{[]uint8{0x0F, 0xE7}, "BSET D7,-(A7)"},

		// BSET Dn,d16(Am) (dynamic bit number, memory destination with displacement)
		{[]uint8{0x01, 0xE8, 0x12, 0x34}, "BSET D0,(0x1234,A0)"},
		{[]uint8{0x03, 0xE9, 0x56, 0x78}, "BSET D1,(0x5678,A1)"},
		{[]uint8{0x05, 0xEA, 0x9A, 0xBC}, "BSET D2,(0x9ABC,A2)"},
		{[]uint8{0x07, 0xEB, 0xDE, 0xF0}, "BSET D3,(0xDEF0,A3)"},
		{[]uint8{0x09, 0xEC, 0x11, 0x22}, "BSET D4,(0x1122,A4)"},
		{[]uint8{0x0B, 0xED, 0x33, 0x44}, "BSET D5,(0x3344,A5)"},
		{[]uint8{0x0D, 0xEE, 0x55, 0x66}, "BSET D6,(0x5566,A6)"},
		{[]uint8{0x0F, 0xEF, 0x77, 0x88}, "BSET D7,(0x7788,A7)"},

		// BSET Dn,d8(Am,Xn) (dynamic bit number, memory destination with index)
		{[]uint8{0x01, 0xF0, 0x12, 0x34}, "BSET D0,(0x12,A0,D4.W)"},
		{[]uint8{0x03, 0xF1, 0x56, 0x78}, "BSET D1,(0x56,A1,D7.W)"},
		{[]uint8{0x05, 0xF2, 0x9A, 0xB0}, "BSET D2,(0x9A,A2,D0.W)"},
		{[]uint8{0x07, 0xF3, 0xDE, 0xF3}, "BSET D3,(0xDE,A3,D3.W)"},
		{[]uint8{0x01, 0xF0, 0x12, 0xB4}, "BSET D0,(0x12,A0,A4.W)"},
		{[]uint8{0x03, 0xF1, 0x56, 0xF7}, "BSET D1,(0x56,A1,A7.W)"},

		// BSET Dn,xxx.W (dynamic bit number, absolute short address)
		{[]uint8{0x01, 0xF8, 0x12, 0x34}, "BSET D0,(0x1234).W"},
		{[]uint8{0x03, 0xF8, 0x56, 0x78}, "BSET D1,(0x5678).W"},
		{[]uint8{0x05, 0xF8, 0x9A, 0xBC}, "BSET D2,(0x9ABC).W"},

		// BSET Dn,xxx.L (dynamic bit number, absolute long address)
		{[]uint8{0x01, 0xF9, 0x12, 0x34, 0x56, 0x78}, "BSET D0,(0x12345678).L"},
		{[]uint8{0x03, 0xF9, 0x9A, 0xBC, 0xDE, 0xF0}, "BSET D1,(0x9ABCDEF0).L"},

		// BSET #imm,Dn (static bit number, data register destination)
		{[]uint8{0x08, 0xC0, 0x00, 0x00}, "BSET #0,D0"},
		{[]uint8{0x08, 0xC1, 0x00, 0x07}, "BSET #7,D1"},
		{[]uint8{0x08, 0xC2, 0x00, 0x0F}, "BSET #15,D2"},
		{[]uint8{0x08, 0xC3, 0x00, 0x1F}, "BSET #31,D3"},  // Maximum bit number for long word
		{[]uint8{0x08, 0xC4, 0x00, 0x20}, "BSET #32,D4"},  // Will be modulo 32
		{[]uint8{0x08, 0xC5, 0x00, 0xFF}, "BSET #255,D5"}, // Will be modulo 32

		// BSET #imm,(Am) (static bit number, memory destination)
		{[]uint8{0x08, 0xD0, 0x00, 0x00}, "BSET #0,(A0)"},
		{[]uint8{0x08, 0xD1, 0x00, 0x07}, "BSET #7,(A1)"},
		{[]uint8{0x08, 0xD2, 0x00, 0x08}, "BSET #8,(A2)"},   // Will be modulo 8 for memory ops
		{[]uint8{0x08, 0xD3, 0x00, 0xFF}, "BSET #255,(A3)"}, // Will be modulo 8 for memory ops

		// BSET #imm,(Am)+ (static bit number, memory destination with postincrement)
		{[]uint8{0x08, 0xD8, 0x00, 0x00}, "BSET #0,(A0)+"},
		{[]uint8{0x08, 0xD9, 0x00, 0x07}, "BSET #7,(A1)+"},

		// BSET #imm,-(Am) (static bit number, memory destination with predecrement)
		{[]uint8{0x08, 0xE0, 0x00, 0x00}, "BSET #0,-(A0)"},
		{[]uint8{0x08, 0xE1, 0x00, 0x07}, "BSET #7,-(A1)"},

		// BSET #imm,d16(Am) (static bit number, memory destination with displacement)
		{[]uint8{0x08, 0xE8, 0x12, 0x34, 0x00, 0x00}, "BSET #0,(0x1234,A0)"},
		{[]uint8{0x08, 0xE9, 0x56, 0x78, 0x00, 0x07}, "BSET #7,(0x5678,A1)"},

		// BSET #imm,d8(Am,Xn) (static bit number, memory destination with index)
		{[]uint8{0x08, 0xF0, 0x12, 0x34, 0x00, 0x00}, "BSET #0,(0x12,A0,D4.W)"},
		{[]uint8{0x08, 0xF1, 0x56, 0x78, 0x00, 0x07}, "BSET #7,(0x56,A1,D7.W)"},
		{[]uint8{0x08, 0xF0, 0x12, 0xB4, 0x00, 0x00}, "BSET #0,(0x12,A0,A4.W)"},

		// BSET #imm,xxx.W (static bit number, absolute short address)
		{[]uint8{0x08, 0xF8, 0x12, 0x34, 0x00, 0x00}, "BSET #0,(0x1234).W"},
		{[]uint8{0x08, 0xF8, 0x56, 0x78, 0x00, 0x07}, "BSET #7,(0x5678).W"},

		// BSET #imm,xxx.L (static bit number, absolute long address)
		{[]uint8{0x08, 0xF9, 0x12, 0x34, 0x56, 0x78, 0x00, 0x00}, "BSET #0,(0x12345678).L"},
		{[]uint8{0x08, 0xF9, 0x9A, 0xBC, 0xDE, 0xF0, 0x00, 0x07}, "BSET #7,(0x9ABCDEF0).L"},
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
