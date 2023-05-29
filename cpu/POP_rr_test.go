package cpu

import (
	"go-boy/memory"
	"testing"
)

func Test_POP_BC(t *testing.T) {
	test_POP(t, 0xC1, BC)
}

func Test_POP_DE(t *testing.T) {
	test_POP(t, 0xD1, DE)
}

func Test_POP_HL(t *testing.T) {
	test_POP(t, 0xE1, HL)
}

func Test_POP_AF(t *testing.T) {
	test_POP(t, 0xF1, AF)
}

func test_POP(t *testing.T, opcode byte, reg register) {
	memory := memory.CreateMemory()
	cpu := CreateCPU(memory)

	var expected uint16 = 0x1234
	var expectedLSB byte = 0x34
	var expectedMSB byte = 0x12
	var stackStart uint16 = 0xAABB
	var expectedSP = stackStart + 2

	memory.WriteByte(stackStart, expectedLSB)
	memory.WriteByte(stackStart+1, expectedMSB)
	cpu.setRegShort(SP, stackStart)

	executer := cpu.getOpcode(opcode)
	executer()

	result := cpu.getRegShort(reg)
	if result != expected {
		t.Errorf("Expected %s to be 0x%04X but was 0x%04X", reg.String(), expected, result)
	}

	resultSP := cpu.getRegShort(SP)
	if resultSP != expectedSP {
		t.Errorf("Expected SP to be 0x%04X but was 0x%04X", expectedSP, resultSP)
	}
}
