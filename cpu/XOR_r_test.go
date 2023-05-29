package cpu

import (
	"go-boy/memory"
	"testing"
)

func Test_XOR_B(t *testing.T) {
	test_XOR_r(t, 0xA8, B, 0x93, 0x1F)
	test_XOR_r(t, 0xA8, B, 0x93, 0x93)
}

func Test_XOR_C(t *testing.T) {
	test_XOR_r(t, 0xA9, C, 0x93, 0x1F)
	test_XOR_r(t, 0xA9, C, 0x93, 0x93)
}

func Test_XOR_D(t *testing.T) {
	test_XOR_r(t, 0xAA, D, 0x93, 0x1F)
	test_XOR_r(t, 0xAA, D, 0x93, 0x93)
}

func Test_XOR_E(t *testing.T) {
	test_XOR_r(t, 0xAB, E, 0x93, 0x1F)
	test_XOR_r(t, 0xAB, E, 0x93, 0x93)
}

func Test_XOR_H(t *testing.T) {
	test_XOR_r(t, 0xAC, H, 0x93, 0x1F)
	test_XOR_r(t, 0xAC, H, 0x93, 0x93)
}

func Test_XOR_L(t *testing.T) {
	test_XOR_r(t, 0xAD, L, 0x93, 0x1F)
	test_XOR_r(t, 0xAD, L, 0x93, 0x93)
}

func Test_XOR_HL(t *testing.T) {
	t.Error("Not Implemented")
}

func Test_XOR_A(t *testing.T) {
	test_XOR_r(t, 0xAF, A, 0x93, 0x93)
}

func test_XOR_r(t *testing.T, opcode byte, reg register, aValue byte, bValue byte) {
	memory := memory.CreateMemory()
	cpu := CreateCPU(memory)
	var pc uint16 = 0x4321
	var expectedPC = pc

	expected := aValue ^ bValue
	expectedZ := expected == 0x00

	cpu.setRegShort(PC, pc)
	cpu.setRegByte(A, aValue)
	cpu.setRegByte(reg, bValue)
	cpu.setFlagZ(!expectedZ)
	cpu.setFlagN(true)
	cpu.setFlagH(true)
	cpu.setFlagC(true)

	executor := cpu.getOpcode(opcode)
	executor()

	result := cpu.getRegByte(A)

	if result != expected {
		t.Errorf("Expected result value 0x%02X but got 0x%02X", expected, result)
	}

	resultPC := cpu.getRegShort(PC)

	if resultPC != expectedPC {
		t.Errorf("Expected SP 0x%04X but got 0x%04X", expectedPC, resultPC)
	}

	if cpu.getFlagZ() != expectedZ {
		t.Errorf("Expected Z to be %t but was %t", expectedZ, cpu.getFlagZ())
	}

	if cpu.getFlagN() != false {
		t.Errorf("Expected N to be false but was %t", cpu.getFlagN())
	}

	if cpu.getFlagH() != false {
		t.Errorf("Expected H to be false but was %t", cpu.getFlagH())
	}

	if cpu.getFlagC() != false {
		t.Errorf("Expected C to be false but was %t", cpu.getFlagC())
	}
}
