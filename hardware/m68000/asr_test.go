package m68000

import "testing"

func TestDisassembleASR(t *testing.T) {
	testCases := []struct {
		bytes    []uint8
		expected string
	}{
		// ASR - Arithmetic Shift Right
		// Format 1: Register shifts by immediate count or register
		// Format: 1110 ccc s 00 rrr ddd
		// Where:
		// - ccc is the count (immediate, 1-8 with 0 meaning 8) or register containing count (Dx)
		// - s is the size (00=byte, 01=word, 10=long)
		// - rrr is the register mode (0=immediate count, 1=register count)
		// - ddd is the destination register

		// Format 2: Memory shifts (always by 1)
		// Format: 1110 000 s 00 eee xxx
		// Where:
		// - s is the size (01=word, 10=long)
		// - eee xxx defines the effective address (only memory alterable modes)

		// ASR.B #1,Dn - Shift by immediate count (1-8)
		{[]uint8{0xE0, 0x00}, "ASR.B #1,D0"},
		{[]uint8{0xE2, 0x01}, "ASR.B #1,D1"},
		{[]uint8{0xE4, 0x02}, "ASR.B #2,D2"},
		{[]uint8{0xE6, 0x03}, "ASR.B #3,D3"},
		{[]uint8{0xE8, 0x04}, "ASR.B #4,D4"},
		{[]uint8{0xEA, 0x05}, "ASR.B #5,D5"},
		{[]uint8{0xEC, 0x06}, "ASR.B #6,D6"},
		{[]uint8{0xEE, 0x07}, "ASR.B #7,D7"},
		{[]uint8{0xE0, 0x00}, "ASR.B #1,D0"},
		{[]uint8{0xE2, 0x00}, "ASR.B #1,D0"},
		{[]uint8{0xE4, 0x00}, "ASR.B #2,D0"},
		{[]uint8{0xE6, 0x00}, "ASR.B #3,D0"},
		{[]uint8{0xE8, 0x00}, "ASR.B #4,D0"},
		{[]uint8{0xEA, 0x00}, "ASR.B #5,D0"},
		{[]uint8{0xEC, 0x00}, "ASR.B #6,D0"},
		{[]uint8{0xEE, 0x00}, "ASR.B #7,D0"},
		{[]uint8{0xE0, 0x00}, "ASR.B #1,D0"},

		// ASR.B #8,Dn (count=0 means 8)
		{[]uint8{0xE0, 0x00}, "ASR.B #1,D0"},

		// ASR.W #1,Dn - Shift by immediate count (1-8)
		{[]uint8{0xE0, 0x40}, "ASR.W #1,D0"},
		{[]uint8{0xE2, 0x41}, "ASR.W #1,D1"},
		{[]uint8{0xE4, 0x42}, "ASR.W #2,D2"},
		{[]uint8{0xE6, 0x43}, "ASR.W #3,D3"},
		{[]uint8{0xE8, 0x44}, "ASR.W #4,D4"},
		{[]uint8{0xEA, 0x45}, "ASR.W #5,D5"},
		{[]uint8{0xEC, 0x46}, "ASR.W #6,D6"},
		{[]uint8{0xEE, 0x47}, "ASR.W #7,D7"},

		// ASR.W #8,Dn (count=0 means 8)
		{[]uint8{0xE0, 0x40}, "ASR.W #1,D0"},

		// ASR.L #1,Dn - Shift by immediate count (1-8)
		{[]uint8{0xE0, 0x80}, "ASR.L #1,D0"},
		{[]uint8{0xE2, 0x81}, "ASR.L #1,D1"},
		{[]uint8{0xE4, 0x82}, "ASR.L #2,D2"},
		{[]uint8{0xE6, 0x83}, "ASR.L #3,D3"},
		{[]uint8{0xE8, 0x84}, "ASR.L #4,D4"},
		{[]uint8{0xEA, 0x85}, "ASR.L #5,D5"},
		{[]uint8{0xEC, 0x86}, "ASR.L #6,D6"},
		{[]uint8{0xEE, 0x87}, "ASR.L #7,D7"},

		// ASR.L #8,Dn (count=0 means 8)
		{[]uint8{0xE0, 0x80}, "ASR.L #1,D0"},

		// ASR.B Dm,Dn - Shift by register count
		{[]uint8{0xE0, 0x20}, "ASR.B D1,D0"},
		{[]uint8{0xE1, 0x21}, "ASR.B D2,D1"},
		{[]uint8{0xE2, 0x22}, "ASR.B D3,D2"},
		{[]uint8{0xE3, 0x23}, "ASR.B D4,D3"},
		{[]uint8{0xE4, 0x24}, "ASR.B D5,D4"},
		{[]uint8{0xE5, 0x25}, "ASR.B D6,D5"},
		{[]uint8{0xE6, 0x26}, "ASR.B D7,D6"},
		{[]uint8{0xE7, 0x27}, "ASR.B D0,D7"},

		// ASR.W Dm,Dn - Shift by register count
		{[]uint8{0xE0, 0x60}, "ASR.W D1,D0"},
		{[]uint8{0xE1, 0x61}, "ASR.W D2,D1"},
		{[]uint8{0xE2, 0x62}, "ASR.W D3,D2"},
		{[]uint8{0xE3, 0x63}, "ASR.W D4,D3"},
		{[]uint8{0xE4, 0x64}, "ASR.W D5,D4"},
		{[]uint8{0xE5, 0x65}, "ASR.W D6,D5"},
		{[]uint8{0xE6, 0x66}, "ASR.W D7,D6"},
		{[]uint8{0xE7, 0x67}, "ASR.W D0,D7"},

		// ASR.L Dm,Dn - Shift by register count
		{[]uint8{0xE0, 0xA0}, "ASR.L D1,D0"},
		{[]uint8{0xE1, 0xA1}, "ASR.L D2,D1"},
		{[]uint8{0xE2, 0xA2}, "ASR.L D3,D2"},
		{[]uint8{0xE3, 0xA3}, "ASR.L D4,D3"},
		{[]uint8{0xE4, 0xA4}, "ASR.L D5,D4"},
		{[]uint8{0xE5, 0xA5}, "ASR.L D6,D5"},
		{[]uint8{0xE6, 0xA6}, "ASR.L D7,D6"},
		{[]uint8{0xE7, 0xA7}, "ASR.L D0,D7"},

		// ASR.W (EA) - Memory shift (always by 1, only word and long sizes)
		// ASR.W (An)
		{[]uint8{0xE0, 0xD0}, "ASR.W (A0)"},
		{[]uint8{0xE0, 0xD1}, "ASR.W (A1)"},
		{[]uint8{0xE0, 0xD2}, "ASR.W (A2)"},
		{[]uint8{0xE0, 0xD3}, "ASR.W (A3)"},
		{[]uint8{0xE0, 0xD4}, "ASR.W (A4)"},
		{[]uint8{0xE0, 0xD5}, "ASR.W (A5)"},
		{[]uint8{0xE0, 0xD6}, "ASR.W (A6)"},
		{[]uint8{0xE0, 0xD7}, "ASR.W (A7)"},

		// ASR.W (An)+
		{[]uint8{0xE0, 0xD8}, "ASR.W (A0)+"},
		{[]uint8{0xE0, 0xD9}, "ASR.W (A1)+"},
		{[]uint8{0xE0, 0xDA}, "ASR.W (A2)+"},
		{[]uint8{0xE0, 0xDB}, "ASR.W (A3)+"},
		{[]uint8{0xE0, 0xDC}, "ASR.W (A4)+"},
		{[]uint8{0xE0, 0xDD}, "ASR.W (A5)+"},
		{[]uint8{0xE0, 0xDE}, "ASR.W (A6)+"},
		{[]uint8{0xE0, 0xDF}, "ASR.W (A7)+"},

		// ASR.W -(An)
		{[]uint8{0xE0, 0xE0}, "ASR.W -(A0)"},
		{[]uint8{0xE0, 0xE1}, "ASR.W -(A1)"},
		{[]uint8{0xE0, 0xE2}, "ASR.W -(A2)"},
		{[]uint8{0xE0, 0xE3}, "ASR.W -(A3)"},
		{[]uint8{0xE0, 0xE4}, "ASR.W -(A4)"},
		{[]uint8{0xE0, 0xE5}, "ASR.W -(A5)"},
		{[]uint8{0xE0, 0xE6}, "ASR.W -(A6)"},
		{[]uint8{0xE0, 0xE7}, "ASR.W -(A7)"},

		// ASR.W d16(An)
		{[]uint8{0xE0, 0xE8, 0x12, 0x34}, "ASR.W (0x1234,A0)"},
		{[]uint8{0xE0, 0xE9, 0x56, 0x78}, "ASR.W (0x5678,A1)"},
		{[]uint8{0xE0, 0xEA, 0x9A, 0xBC}, "ASR.W (0x9ABC,A2)"},
		{[]uint8{0xE0, 0xEB, 0xDE, 0xF0}, "ASR.W (0xDEF0,A3)"},
		{[]uint8{0xE0, 0xEC, 0x12, 0x34}, "ASR.W (0x1234,A4)"},
		{[]uint8{0xE0, 0xED, 0x56, 0x78}, "ASR.W (0x5678,A5)"},
		{[]uint8{0xE0, 0xEE, 0x9A, 0xBC}, "ASR.W (0x9ABC,A6)"},
		{[]uint8{0xE0, 0xEF, 0xDE, 0xF0}, "ASR.W (0xDEF0,A7)"},

		// ASR.W d8(An,Xn)
		{[]uint8{0xE0, 0xF0, 0x12, 0x34}, "ASR.W (0x12,A0,D4.W)"},
		{[]uint8{0xE0, 0xF1, 0x56, 0x70}, "ASR.W (0x56,A1,D0.W)"},
		{[]uint8{0xE0, 0xF2, 0x9A, 0xB1}, "ASR.W (0x9A,A2,D1.W)"},
		{[]uint8{0xE0, 0xF3, 0xDE, 0xF2}, "ASR.W (0xDE,A3,D2.W)"},
		{[]uint8{0xE0, 0xF4, 0x12, 0x33}, "ASR.W (0x12,A4,D3.W)"},
		{[]uint8{0xE0, 0xF5, 0x56, 0x75}, "ASR.W (0x56,A5,D5.W)"},
		{[]uint8{0xE0, 0xF6, 0x9A, 0xB6}, "ASR.W (0x9A,A6,D6.W)"},
		{[]uint8{0xE0, 0xF7, 0xDE, 0xF7}, "ASR.W (0xDE,A7,D7.W)"},

		// ASR.W d8(An,Xn) with A register index
		{[]uint8{0xE0, 0xF0, 0x12, 0xB4}, "ASR.W (0x12,A0,A4.W)"},
		{[]uint8{0xE0, 0xF3, 0xDE, 0xF2}, "ASR.W (0xDE,A3,D2.W)"},

		// ASR.W xxx.W (Absolute Short)
		{[]uint8{0xE0, 0xF8, 0x12, 0x34}, "ASR.W (0x1234).W"},
		{[]uint8{0xE0, 0xF8, 0x56, 0x78}, "ASR.W (0x5678).W"},
		{[]uint8{0xE0, 0xF8, 0x9A, 0xBC}, "ASR.W (0x9ABC).W"},

		// ASR.W xxx.L (Absolute Long)
		{[]uint8{0xE0, 0xF9, 0x12, 0x34, 0x56, 0x78}, "ASR.W (0x12345678).L"},
		{[]uint8{0xE0, 0xF9, 0x9A, 0xBC, 0xDE, 0xF0}, "ASR.W (0x9ABCDEF0).L"},
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
