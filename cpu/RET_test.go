package cpu

import (
	"github.com/f1gopher/gbpixellib/memory"
	"testing"
)

func Test_RET(t *testing.T) {
	memory := memory.CreateMemory()
	cpu := CreateCPU(memory)

	var stackStart uint16 = 0x1234
	var expectedSP = stackStart + 2
	var expectedLSB byte = 0xCD
	var expectedMSB byte = 0xAB
	var expected uint16 = 0xABCD

	memory.WriteByte(stackStart, expectedLSB)
	memory.WriteByte(stackStart+1, expectedMSB)
	cpu.setRegShort(SP, stackStart)

	executor := cpu.getOpcode(0xC9)
	executor()

	result := cpu.getRegShort(PC)
	if result != expected {
		t.Errorf("Expected PC to be 0x%04X but was 0x%04X", expected, result)
	}

	resultSP := cpu.getRegShort(SP)
	if resultSP != expectedSP {
		t.Errorf("Expected SP to 0x%04X but was 0x%04X", expectedSP, resultSP)
	}
}
