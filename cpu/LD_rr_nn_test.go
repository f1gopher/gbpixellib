package cpu

import (
	"go-boy/memory"
	"testing"
)

func Test_LD_BC_nn(t *testing.T) {
	test_LD_rr_nn(t, 0x01, BC)
}

func Test_LD_DE_nn(t *testing.T) {
	test_LD_rr_nn(t, 0x11, DE)
}

func Test_LD_HL_nn(t *testing.T) {
	test_LD_rr_nn(t, 0x21, HL)
}

func Test_LD_SP_nn(t *testing.T) {
	test_LD_rr_nn(t, 0x31, SP)
}

func test_LD_rr_nn(t *testing.T, opcode byte, reg register) {
	memory := memory.CreateMemory()
	cpu := CreateCPU(memory)
	var pc uint16 = 0x4321
	var expected uint16 = 0x1234
	var expectedPC uint16 = pc + 2

	cpu.setRegShort(PC, pc)
	memory.WriteShort(pc, expected)

	executor := cpu.getOpcode(opcode)
	executor()

	result := cpu.getRegShort(reg)

	if result != expected {
		t.Errorf("Expected result value 0x%04X but got 0x%04X", expected, result)
	}

	resultPC := cpu.getRegShort(PC)

	if resultPC != expectedPC {
		t.Errorf("Expected SP 0x%04X but got 0x%04X", expectedPC, resultPC)
	}
}
