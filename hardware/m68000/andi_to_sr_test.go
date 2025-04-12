package m68000

import "testing"

func TestDisassembleANDI_to_SR(t *testing.T) {
	testCases := []struct {
		bytes    []uint8
		expected string
	}{
		// ANDI to SR (Status Register)
		// Format: 0000 0010 0111 1100 xxxx xxxx xxxx xxxx
		// Where xxxx xxxx xxxx xxxx is a 16-bit immediate value that will be ANDed with the SR
		// This is a privileged instruction - can only be executed in supervisor mode

		// Test various bit patterns targeting different sections of the SR
		// SR bits: T--S---XNZVC (T=trace, S=supervisor, X/N/Z/V/C=condition codes)

		// Full SR manipulation
		{[]uint8{0x02, 0x7C, 0x00, 0x00}, "ANDI #0x0000,SR"}, // Clear all bits
		{[]uint8{0x02, 0x7C, 0xFF, 0xFF}, "ANDI #0xFFFF,SR"}, // Preserve all bits

		// Condition code manipulation
		{[]uint8{0x02, 0x7C, 0xFF, 0xE0}, "ANDI #0xFFE0,SR"}, // Clear all CCR bits
		{[]uint8{0x02, 0x7C, 0xFF, 0x1F}, "ANDI #0xFF1F,SR"}, // Preserve all CCR bits

		// Individual flag manipulation
		{[]uint8{0x02, 0x7C, 0xFF, 0xFE}, "ANDI #0xFFFE,SR"}, // Clear C flag
		{[]uint8{0x02, 0x7C, 0xFF, 0xFD}, "ANDI #0xFFFD,SR"}, // Clear V flag
		{[]uint8{0x02, 0x7C, 0xFF, 0xFB}, "ANDI #0xFFFB,SR"}, // Clear Z flag
		{[]uint8{0x02, 0x7C, 0xFF, 0xF7}, "ANDI #0xFFF7,SR"}, // Clear N flag
		{[]uint8{0x02, 0x7C, 0xFF, 0xEF}, "ANDI #0xFFEF,SR"}, // Clear X flag

		// System flag manipulation
		{[]uint8{0x02, 0x7C, 0xF7, 0xFF}, "ANDI #0xF7FF,SR"}, // Clear Supervisor bit
		{[]uint8{0x02, 0x7C, 0x7F, 0xFF}, "ANDI #0x7FFF,SR"}, // Clear Trace bit

		// Interrupt mask manipulation
		{[]uint8{0x02, 0x7C, 0xF8, 0xFF}, "ANDI #0xF8FF,SR"}, // Clear interrupt mask (set to level 0)
		{[]uint8{0x02, 0x7C, 0xFF, 0xFF}, "ANDI #0xFFFF,SR"}, // Preserve interrupt mask
		{[]uint8{0x02, 0x7C, 0xF9, 0xFF}, "ANDI #0xF9FF,SR"}, // Set interrupt mask to max level 1
		{[]uint8{0x02, 0x7C, 0xFA, 0xFF}, "ANDI #0xFAFF,SR"}, // Set interrupt mask to max level 2
		{[]uint8{0x02, 0x7C, 0xFB, 0xFF}, "ANDI #0xFBFF,SR"}, // Set interrupt mask to max level 3
		{[]uint8{0x02, 0x7C, 0xFC, 0xFF}, "ANDI #0xFCFF,SR"}, // Set interrupt mask to max level 4
		{[]uint8{0x02, 0x7C, 0xFD, 0xFF}, "ANDI #0xFDFF,SR"}, // Set interrupt mask to max level 5
		{[]uint8{0x02, 0x7C, 0xFE, 0xFF}, "ANDI #0xFEFF,SR"}, // Set interrupt mask to max level 6
		{[]uint8{0x02, 0x7C, 0xFF, 0xFF}, "ANDI #0xFFFF,SR"}, // Set interrupt mask to max level 7

		// Combined operations
		{[]uint8{0x02, 0x7C, 0x70, 0x00}, "ANDI #0x7000,SR"}, // Clear all except interrupt mask bits
		{[]uint8{0x02, 0x7C, 0x80, 0x1F}, "ANDI #0x801F,SR"}, // Preserve trace bit and condition codes only
		{[]uint8{0x02, 0x7C, 0x00, 0x1F}, "ANDI #0x001F,SR"}, // Keep only condition codes, clear all system flags
		{[]uint8{0x02, 0x7C, 0xF0, 0x00}, "ANDI #0xF000,SR"}, // Keep only upper nibble
		{[]uint8{0x02, 0x7C, 0x0F, 0xFF}, "ANDI #0x0FFF,SR"}, // Clear upper nibble

		// Some "real-world" common usages
		{[]uint8{0x02, 0x7C, 0xD8, 0xFF}, "ANDI #0xD8FF,SR"}, // Clear supervisor and interrupt mask
		{[]uint8{0x02, 0x7C, 0x7F, 0xFF}, "ANDI #0x7FFF,SR"}, // Disable trace mode
		{[]uint8{0x02, 0x7C, 0xF8, 0xFF}, "ANDI #0xF8FF,SR"}, // Enable all interrupts (clear mask)
		{[]uint8{0x02, 0x7C, 0xFF, 0xE0}, "ANDI #0xFFE0,SR"}, // Clear all condition codes
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
