package m68000

import "testing"

func TestDisassembleASL(t *testing.T) {
	testCases := []testCase{
		// ASL - Arithmetic Shift Left
		// Format 1: Register shifts by immediate count or register
		// Format: 1110 ccc s 11 rrr ddd
		// Where:
		// - ccc is the count (immediate, 1-8 with 0 meaning 8) or register containing count (Dx)
		// - s is the size (00=byte, 01=word, 10=long)
		// - rrr is the register mode (0=immediate count, 1=register count)
		// - ddd is the destination register

		// Format 2: Memory shifts (always by 1)
		// Format: 1110 000 s 11 eee xxx
		// Where:
		// - s is the size (01=word, 10=long)
		// - eee xxx defines the effective address (only memory alterable modes)

		// ASL.B #1,Dn - Shift by immediate count (1-8)
		{[]uint8{0xE1, 0x00}, "ASL.B #1,D0"},
		{[]uint8{0xE3, 0x01}, "ASL.B #1,D1"},
		{[]uint8{0xE5, 0x02}, "ASL.B #2,D2"},
		{[]uint8{0xE7, 0x03}, "ASL.B #3,D3"},
		{[]uint8{0xE9, 0x04}, "ASL.B #4,D4"},
		{[]uint8{0xEB, 0x05}, "ASL.B #5,D5"},
		{[]uint8{0xED, 0x06}, "ASL.B #6,D6"},
		{[]uint8{0xEF, 0x07}, "ASL.B #7,D7"},
		{[]uint8{0xE1, 0x00}, "ASL.B #1,D0"},
		{[]uint8{0xE3, 0x00}, "ASL.B #1,D0"},
		{[]uint8{0xE5, 0x00}, "ASL.B #2,D0"},
		{[]uint8{0xE7, 0x00}, "ASL.B #3,D0"},
		{[]uint8{0xE9, 0x00}, "ASL.B #4,D0"},
		{[]uint8{0xEB, 0x00}, "ASL.B #5,D0"},
		{[]uint8{0xED, 0x00}, "ASL.B #6,D0"},
		{[]uint8{0xEF, 0x00}, "ASL.B #7,D0"},
		{[]uint8{0xE1, 0x00}, "ASL.B #1,D0"},

		// ASL.B #8,Dn (count=0 means 8)
		{[]uint8{0xE1, 0x00}, "ASL.B #1,D0"},

		// ASL.W #1,Dn - Shift by immediate count (1-8)
		{[]uint8{0xE1, 0x40}, "ASL.W #1,D0"},
		{[]uint8{0xE3, 0x41}, "ASL.W #1,D1"},
		{[]uint8{0xE5, 0x42}, "ASL.W #2,D2"},
		{[]uint8{0xE7, 0x43}, "ASL.W #3,D3"},
		{[]uint8{0xE9, 0x44}, "ASL.W #4,D4"},
		{[]uint8{0xEB, 0x45}, "ASL.W #5,D5"},
		{[]uint8{0xED, 0x46}, "ASL.W #6,D6"},
		{[]uint8{0xEF, 0x47}, "ASL.W #7,D7"},

		// ASL.W #8,Dn (count=0 means 8)
		{[]uint8{0xE1, 0x40}, "ASL.W #1,D0"},

		// ASL.L #1,Dn - Shift by immediate count (1-8)
		{[]uint8{0xE1, 0x80}, "ASL.L #1,D0"},
		{[]uint8{0xE3, 0x81}, "ASL.L #1,D1"},
		{[]uint8{0xE5, 0x82}, "ASL.L #2,D2"},
		{[]uint8{0xE7, 0x83}, "ASL.L #3,D3"},
		{[]uint8{0xE9, 0x84}, "ASL.L #4,D4"},
		{[]uint8{0xEB, 0x85}, "ASL.L #5,D5"},
		{[]uint8{0xED, 0x86}, "ASL.L #6,D6"},
		{[]uint8{0xEF, 0x87}, "ASL.L #7,D7"},

		// ASL.L #8,Dn (count=0 means 8)
		{[]uint8{0xE1, 0x80}, "ASL.L #1,D0"},

		// ASL.B Dm,Dn - Shift by register count
		{[]uint8{0xE1, 0x20}, "ASL.B D1,D0"},
		{[]uint8{0xE2, 0x21}, "ASL.B D2,D1"},
		{[]uint8{0xE3, 0x22}, "ASL.B D3,D2"},
		{[]uint8{0xE4, 0x23}, "ASL.B D4,D3"},
		{[]uint8{0xE5, 0x24}, "ASL.B D5,D4"},
		{[]uint8{0xE6, 0x25}, "ASL.B D6,D5"},
		{[]uint8{0xE7, 0x26}, "ASL.B D7,D6"},
		{[]uint8{0xE0, 0x27}, "ASL.B D0,D7"},

		// ASL.W Dm,Dn - Shift by register count
		{[]uint8{0xE1, 0x60}, "ASL.W D1,D0"},
		{[]uint8{0xE2, 0x61}, "ASL.W D2,D1"},
		{[]uint8{0xE3, 0x62}, "ASL.W D3,D2"},
		{[]uint8{0xE4, 0x63}, "ASL.W D4,D3"},
		{[]uint8{0xE5, 0x64}, "ASL.W D5,D4"},
		{[]uint8{0xE6, 0x65}, "ASL.W D6,D5"},
		{[]uint8{0xE7, 0x66}, "ASL.W D7,D6"},
		{[]uint8{0xE0, 0x67}, "ASL.W D0,D7"},

		// ASL.L Dm,Dn - Shift by register count
		{[]uint8{0xE1, 0xA0}, "ASL.L D1,D0"},
		{[]uint8{0xE2, 0xA1}, "ASL.L D2,D1"},
		{[]uint8{0xE3, 0xA2}, "ASL.L D3,D2"},
		{[]uint8{0xE4, 0xA3}, "ASL.L D4,D3"},
		{[]uint8{0xE5, 0xA4}, "ASL.L D5,D4"},
		{[]uint8{0xE6, 0xA5}, "ASL.L D6,D5"},
		{[]uint8{0xE7, 0xA6}, "ASL.L D7,D6"},
		{[]uint8{0xE0, 0xA7}, "ASL.L D0,D7"},

		// ASL.W (EA) - Memory shift (always by 1, only word and long sizes)
		// ASL.W (An)
		{[]uint8{0xE1, 0xD0}, "ASL.W (A0)"},
		{[]uint8{0xE1, 0xD1}, "ASL.W (A1)"},
		{[]uint8{0xE1, 0xD2}, "ASL.W (A2)"},
		{[]uint8{0xE1, 0xD3}, "ASL.W (A3)"},
		{[]uint8{0xE1, 0xD4}, "ASL.W (A4)"},
		{[]uint8{0xE1, 0xD5}, "ASL.W (A5)"},
		{[]uint8{0xE1, 0xD6}, "ASL.W (A6)"},
		{[]uint8{0xE1, 0xD7}, "ASL.W (A7)"},

		// ASL.W (An)+
		{[]uint8{0xE1, 0xD8}, "ASL.W (A0)+"},
		{[]uint8{0xE1, 0xD9}, "ASL.W (A1)+"},
		{[]uint8{0xE1, 0xDA}, "ASL.W (A2)+"},
		{[]uint8{0xE1, 0xDB}, "ASL.W (A3)+"},
		{[]uint8{0xE1, 0xDC}, "ASL.W (A4)+"},
		{[]uint8{0xE1, 0xDD}, "ASL.W (A5)+"},
		{[]uint8{0xE1, 0xDE}, "ASL.W (A6)+"},
		{[]uint8{0xE1, 0xDF}, "ASL.W (A7)+"},

		// ASL.W -(An)
		{[]uint8{0xE1, 0xE0}, "ASL.W -(A0)"},
		{[]uint8{0xE1, 0xE1}, "ASL.W -(A1)"},
		{[]uint8{0xE1, 0xE2}, "ASL.W -(A2)"},
		{[]uint8{0xE1, 0xE3}, "ASL.W -(A3)"},
		{[]uint8{0xE1, 0xE4}, "ASL.W -(A4)"},
		{[]uint8{0xE1, 0xE5}, "ASL.W -(A5)"},
		{[]uint8{0xE1, 0xE6}, "ASL.W -(A6)"},
		{[]uint8{0xE1, 0xE7}, "ASL.W -(A7)"},

		// ASL.W d16(An)
		{[]uint8{0xE1, 0xE8, 0x12, 0x34}, "ASL.W (0x1234,A0)"},
		{[]uint8{0xE1, 0xE9, 0x56, 0x78}, "ASL.W (0x5678,A1)"},
		{[]uint8{0xE1, 0xEA, 0x9A, 0xBC}, "ASL.W (0x9ABC,A2)"},
		{[]uint8{0xE1, 0xEB, 0xDE, 0xF0}, "ASL.W (0xDEF0,A3)"},
		{[]uint8{0xE1, 0xEC, 0x12, 0x34}, "ASL.W (0x1234,A4)"},
		{[]uint8{0xE1, 0xED, 0x56, 0x78}, "ASL.W (0x5678,A5)"},
		{[]uint8{0xE1, 0xEE, 0x9A, 0xBC}, "ASL.W (0x9ABC,A6)"},
		{[]uint8{0xE1, 0xEF, 0xDE, 0xF0}, "ASL.W (0xDEF0,A7)"},

		// ASL.W d8(An,Xn)
		{[]uint8{0xE1, 0xF0, 0x12, 0x34}, "ASL.W (0x12,A0,D4.W)"},
		{[]uint8{0xE1, 0xF1, 0x56, 0x70}, "ASL.W (0x56,A1,D0.W)"},
		{[]uint8{0xE1, 0xF2, 0x9A, 0xB1}, "ASL.W (0x9A,A2,D1.W)"},
		{[]uint8{0xE1, 0xF3, 0xDE, 0xF2}, "ASL.W (0xDE,A3,D2.W)"},
		{[]uint8{0xE1, 0xF4, 0x12, 0x33}, "ASL.W (0x12,A4,D3.W)"},
		{[]uint8{0xE1, 0xF5, 0x56, 0x75}, "ASL.W (0x56,A5,D5.W)"},
		{[]uint8{0xE1, 0xF6, 0x9A, 0xB6}, "ASL.W (0x9A,A6,D6.W)"},
		{[]uint8{0xE1, 0xF7, 0xDE, 0xF7}, "ASL.W (0xDE,A7,D7.W)"},

		// ASL.W d8(An,Xn) with A register index
		{[]uint8{0xE1, 0xF0, 0x12, 0xB4}, "ASL.W (0x12,A0,A4.W)"},
		{[]uint8{0xE1, 0xF3, 0xDE, 0xF2}, "ASL.W (0xDE,A3,D2.W)"},

		// ASL.W xxx.W (Absolute Short)
		{[]uint8{0xE1, 0xF8, 0x12, 0x34}, "ASL.W (0x1234).W"},
		{[]uint8{0xE1, 0xF8, 0x56, 0x78}, "ASL.W (0x5678).W"},
		{[]uint8{0xE1, 0xF8, 0x9A, 0xBC}, "ASL.W (0x9ABC).W"},

		// ASL.W xxx.L (Absolute Long)
		{[]uint8{0xE1, 0xF9, 0x12, 0x34, 0x56, 0x78}, "ASL.W (0x12345678).L"},
		{[]uint8{0xE1, 0xF9, 0x9A, 0xBC, 0xDE, 0xF0}, "ASL.W (0x9ABCDEF0).L"},
	}

	runTests(t, testCases)
}
