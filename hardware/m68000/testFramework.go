package m68000

import (
	"fmt"
	"testing"
)

type testCase struct {
	bytes    []uint8
	expected string
}

type testCaseWithAddress struct {
	bytes    []uint8
	expected string
	address  uint32
}

func DisassembleInstruction(code []uint8, address uint32) string {
	parser := instructionParser{}
	inst := parser.parse(code, int(address))
	return inst.Name()
}

func runTests(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, 0)
			if result != tc.expected {
				hex := ""
				binary := ""
				for x := range tc.bytes {
					if x%2 == 0 && x != 0 {
						hex += " "
						binary += " "
					}

					hex += fmt.Sprintf("%02X", tc.bytes[x])
					binary += fmt.Sprintf("%08b", tc.bytes[x])
				}

				t.Errorf("Expected: %s -------> %s\n%s\n%s", tc.expected, result, hex, binary)
			}
		})
	}
}

func runAddressTests(t *testing.T, testCases []testCaseWithAddress) {
	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := DisassembleInstruction(tc.bytes, tc.address)
			if result != tc.expected {
				hex := ""
				binary := ""
				for x := range tc.bytes {
					if x%2 == 0 && x != 0 {
						hex += " "
						binary += " "
					}

					hex += fmt.Sprintf("%02X", tc.bytes[x])
					binary += fmt.Sprintf("%08b", tc.bytes[x])
				}

				t.Errorf("Expected: %s -------> %s\n%s\n%s\nAddress=%d", tc.expected, result, hex, binary, tc.address)
			}
		})
	}
}
