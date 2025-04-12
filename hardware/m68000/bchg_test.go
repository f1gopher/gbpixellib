package m68000

import "testing"

func TestDisassembleBCHG(t *testing.T) {
	testCases := []struct {
		bytes    []uint8
		expected string
	}{
		// BCHG - Bit Change
		// Format 1: BCHG Dn,<ea> (dynamic bit number from register)
		// Format: 0000 rrr 101 eee mmm (where rrr is the register number)

		// Format 2: BCHG #<data>,<ea> (static bit number from immediate)
		// Format: 0000 1000 01 eee mmm (followed by bit number as immediate word)

		// BCHG Dn,Dm (dynamic bit number, data register destination)
		{[]uint8{0x01, 0x40}, "BCHG D0,D0"},
		{[]uint8{0x03, 0x41}, "BCHG D1,D1"},
		{[]uint8{0x05, 0x42}, "BCHG D2,D2"},
		{[]uint8{0x07, 0x43}, "BCHG D3,D3"},
		{[]uint8{0x09, 0x44}, "BCHG D4,D4"},
		{[]uint8{0x0B, 0x45}, "BCHG D5,D5"},
		{[]uint8{0x0D, 0x46}, "BCHG D6,D6"},
		{[]uint8{0x0F, 0x47}, "BCHG D7,D7"},
		{[]uint8{0x01, 0x47}, "BCHG D0,D7"},
		{[]uint8{0x03, 0x46}, "BCHG D1,D6"},
		{[]uint8{0x05, 0x45}, "BCHG D2,D5"},
		{[]uint8{0x07, 0x44}, "BCHG D3,D4"},
		{[]uint8{0x09, 0x43}, "BCHG D4,D3"},
		{[]uint8{0x0B, 0x42}, "BCHG D5,D2"},
		{[]uint8{0x0D, 0x41}, "BCHG D6,D1"},
		{[]uint8{0x0F, 0x40}, "BCHG D7,D0"},

		// BCHG Dn,(Am) (dynamic bit number, memory destination)
		{[]uint8{0x01, 0x50}, "BCHG D0,(A0)"},
		{[]uint8{0x03, 0x51}, "BCHG D1,(A1)"},
		{[]uint8{0x05, 0x52}, "BCHG D2,(A2)"},
		{[]uint8{0x07, 0x53}, "BCHG D3,(A3)"},
		{[]uint8{0x09, 0x54}, "BCHG D4,(A4)"},
		{[]uint8{0x0B, 0x55}, "BCHG D5,(A5)"},
		{[]uint8{0x0D, 0x56}, "BCHG D6,(A6)"},
		{[]uint8{0x0F, 0x57}, "BCHG D7,(A7)"},

		// BCHG Dn,(Am)+ (dynamic bit number, memory destination with postincrement)
		{[]uint8{0x01, 0x58}, "BCHG D0,(A0)+"},
		{[]uint8{0x03, 0x59}, "BCHG D1,(A1)+"},
		{[]uint8{0x05, 0x5A}, "BCHG D2,(A2)+"},
		{[]uint8{0x07, 0x5B}, "BCHG D3,(A3)+"},
		{[]uint8{0x09, 0x5C}, "BCHG D4,(A4)+"},
		{[]uint8{0x0B, 0x5D}, "BCHG D5,(A5)+"},
		{[]uint8{0x0D, 0x5E}, "BCHG D6,(A6)+"},
		{[]uint8{0x0F, 0x5F}, "BCHG D7,(A7)+"},

		// BCHG Dn,-(Am) (dynamic bit number, memory destination with predecrement)
		{[]uint8{0x01, 0x60}, "BCHG D0,-(A0)"},
		{[]uint8{0x03, 0x61}, "BCHG D1,-(A1)"},
		{[]uint8{0x05, 0x62}, "BCHG D2,-(A2)"},
		{[]uint8{0x07, 0x63}, "BCHG D3,-(A3)"},
		{[]uint8{0x09, 0x64}, "BCHG D4,-(A4)"},
		{[]uint8{0x0B, 0x65}, "BCHG D5,-(A5)"},
		{[]uint8{0x0D, 0x66}, "BCHG D6,-(A6)"},
		{[]uint8{0x0F, 0x67}, "BCHG D7,-(A7)"},

		// BCHG Dn,d16(Am) (dynamic bit number, memory destination with displacement)
		{[]uint8{0x01, 0x68, 0x12, 0x34}, "BCHG D0,(0x1234,A0)"},
		{[]uint8{0x03, 0x69, 0x56, 0x78}, "BCHG D1,(0x5678,A1)"},
		{[]uint8{0x05, 0x6A, 0x9A, 0xBC}, "BCHG D2,(0x9ABC,A2)"},
		{[]uint8{0x07, 0x6B, 0xDE, 0xF0}, "BCHG D3,(0xDEF0,A3)"},
		{[]uint8{0x09, 0x6C, 0x11, 0x22}, "BCHG D4,(0x1122,A4)"},
		{[]uint8{0x0B, 0x6D, 0x33, 0x44}, "BCHG D5,(0x3344,A5)"},
		{[]uint8{0x0D, 0x6E, 0x55, 0x66}, "BCHG D6,(0x5566,A6)"},
		{[]uint8{0x0F, 0x6F, 0x77, 0x88}, "BCHG D7,(0x7788,A7)"},

		// BCHG Dn,d8(Am,Xn) (dynamic bit number, memory destination with index)
		{[]uint8{0x01, 0x70, 0x12, 0x34}, "BCHG D0,(0x12,A0,D4.W)"},
		{[]uint8{0x03, 0x71, 0x56, 0x78}, "BCHG D1,(0x56,A1,D7.W)"},
		{[]uint8{0x05, 0x72, 0x9A, 0xB0}, "BCHG D2,(0x9A,A2,D0.W)"},
		{[]uint8{0x07, 0x73, 0xDE, 0xF3}, "BCHG D3,(0xDE,A3,D3.W)"},
		{[]uint8{0x01, 0x70, 0x12, 0xB4}, "BCHG D0,(0x12,A0,A4.W)"},
		{[]uint8{0x03, 0x71, 0x56, 0xF7}, "BCHG D1,(0x56,A1,A7.W)"},

		// BCHG Dn,xxx.W (dynamic bit number, absolute short address)
		{[]uint8{0x01, 0x78, 0x12, 0x34}, "BCHG D0,(0x1234).W"},
		{[]uint8{0x03, 0x78, 0x56, 0x78}, "BCHG D1,(0x5678).W"},
		{[]uint8{0x05, 0x78, 0x9A, 0xBC}, "BCHG D2,(0x9ABC).W"},

		// BCHG Dn,xxx.L (dynamic bit number, absolute long address)
		{[]uint8{0x01, 0x79, 0x12, 0x34, 0x56, 0x78}, "BCHG D0,(0x12345678).L"},
		{[]uint8{0x03, 0x79, 0x9A, 0xBC, 0xDE, 0xF0}, "BCHG D1,(0x9ABCDEF0).L"},

		// BCHG #imm,Dn (static bit number, data register destination)
		{[]uint8{0x08, 0x40, 0x00, 0x00}, "BCHG #0,D0"},
		{[]uint8{0x08, 0x41, 0x00, 0x07}, "BCHG #7,D1"},
		{[]uint8{0x08, 0x42, 0x00, 0x0F}, "BCHG #15,D2"},
		{[]uint8{0x08, 0x43, 0x00, 0x1F}, "BCHG #31,D3"},  // Maximum bit number for long word
		{[]uint8{0x08, 0x44, 0x00, 0x20}, "BCHG #32,D4"},  // Will be modulo 32
		{[]uint8{0x08, 0x45, 0x00, 0xFF}, "BCHG #255,D5"}, // Will be modulo 32

		// BCHG #imm,(Am) (static bit number, memory destination)
		{[]uint8{0x08, 0x50, 0x00, 0x00}, "BCHG #0,(A0)"},
		{[]uint8{0x08, 0x51, 0x00, 0x07}, "BCHG #7,(A1)"},
		{[]uint8{0x08, 0x52, 0x00, 0x08}, "BCHG #8,(A2)"},   // Will be modulo 8 for memory ops
		{[]uint8{0x08, 0x53, 0x00, 0xFF}, "BCHG #255,(A3)"}, // Will be modulo 8 for memory ops

		// BCHG #imm,(Am)+ (static bit number, memory destination with postincrement)
		{[]uint8{0x08, 0x58, 0x00, 0x00}, "BCHG #0,(A0)+"},
		{[]uint8{0x08, 0x59, 0x00, 0x07}, "BCHG #7,(A1)+"},

		// BCHG #imm,-(Am) (static bit number, memory destination with predecrement)
		{[]uint8{0x08, 0x60, 0x00, 0x00}, "BCHG #0,-(A0)"},
		{[]uint8{0x08, 0x61, 0x00, 0x07}, "BCHG #7,-(A1)"},

		// BCHG #imm,d16(Am) (static bit number, memory destination with displacement)
		{[]uint8{0x08, 0x68, 0x12, 0x34, 0x00, 0x00}, "BCHG #0,(0x1234,A0)"},
		{[]uint8{0x08, 0x69, 0x56, 0x78, 0x00, 0x07}, "BCHG #7,(0x5678,A1)"},

		// BCHG #imm,d8(Am,Xn) (static bit number, memory destination with index)
		{[]uint8{0x08, 0x70, 0x12, 0x34, 0x00, 0x00}, "BCHG #0,(0x12,A0,D4.W)"},
		{[]uint8{0x08, 0x71, 0x56, 0x78, 0x00, 0x07}, "BCHG #7,(0x56,A1,D7.W)"},
		{[]uint8{0x08, 0x70, 0x12, 0xB4, 0x00, 0x00}, "BCHG #0,(0x12,A0,A4.W)"},

		// BCHG #imm,xxx.W (static bit number, absolute short address)
		{[]uint8{0x08, 0x78, 0x12, 0x34, 0x00, 0x00}, "BCHG #0,(0x1234).W"},
		{[]uint8{0x08, 0x78, 0x56, 0x78, 0x00, 0x07}, "BCHG #7,(0x5678).W"},

		// BCHG #imm,xxx.L (static bit number, absolute long address)
		{[]uint8{0x08, 0x79, 0x12, 0x34, 0x56, 0x78, 0x00, 0x00}, "BCHG #0,(0x12345678).L"},
		{[]uint8{0x08, 0x79, 0x9A, 0xBC, 0xDE, 0xF0, 0x00, 0x07}, "BCHG #7,(0x9ABCDEF0).L"},
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
