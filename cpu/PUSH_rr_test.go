package cpu

import (
	"github.com/f1gopher/gbpixellib/memory"
	"testing"
)

func Test_PUSH_BC(t *testing.T) {
	test_PUSH(t, 0xC5, BC)
}

func Test_PUSH_DE(t *testing.T) {
	test_PUSH(t, 0xD5, DE)
}

func Test_PUSH_HL(t *testing.T) {
	test_PUSH(t, 0xE5, HL)
}

func Test_PUSH_AF(t *testing.T) {
	test_PUSH(t, 0xF5, AF)
}

func test_PUSH(t *testing.T, opcode byte, reg register) {
	memory := memory.CreateMemory()
	cpu := CreateCPU(memory)

	var expected uint16 = 0x1234
	var expectedLSB byte = 0x34
	var expectedMSB byte = 0x12
	var stackStart uint16 = 0xAABB
	var expectedSP = stackStart - 2

	cpu.setRegShort(SP, stackStart)
	cpu.setRegShort(reg, expected)

	executer := cpu.getOpcode(opcode)
	executer()

	resultMSB := memory.ReadByte(stackStart - 1)
	if resultMSB != expectedMSB {
		t.Errorf("Expected MSB to be 0x%02X but was 0x%02X", expectedMSB, resultMSB)
	}

	resultLSB := memory.ReadByte(stackStart - 2)
	if resultLSB != expectedLSB {
		t.Errorf("Expected LSB to be 0x%02X but was 0x%02X", expectedLSB, resultLSB)
	}

	resultSP := cpu.getRegShort(SP)
	if resultSP != expectedSP {
		t.Errorf("Expected SP to be 0x%04X but was 0x%04X", expectedSP, resultSP)
	}
}
