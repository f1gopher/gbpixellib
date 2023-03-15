package cpu

import "testing"

func TestGetLowByte(t *testing.T) {
	result := getLowByte(0xABCD)

	if result != 0xCD {
		t.Errorf("Expected 0xCD but got 0x%02X", result)
	}
}

func TestGetHighByte(t *testing.T) {
	result := getHighByte(0xABCD)
	if result != 0xAB {
		t.Errorf("Expected 0xCD but got 0x%02X", result)
	}
}

func TestSetLowByte(t *testing.T) {
	var reg uint16 = 0xABCD
	setLowByte(&reg, 0xEF)
	if reg != 0xABEF {
		t.Errorf("Expected 0xABEF but got 0x%04X", reg)
	}

	reg = 0x0013
	setLowByte(&reg, 0x10)
	if reg != 0x0010 {
		t.Errorf("Expected 0x0010 but got 0x%04X", reg)
	}
}

func TestSetHighByte(t *testing.T) {
	var reg uint16 = 0xABCD
	setHighByte(&reg, 0xEF)
	if reg != 0xEFCD {
		t.Errorf("Expected 0xEFCD but got 0x%04X", reg)
	}

	reg = 0x1300
	setHighByte(&reg, 0x10)
	if reg != 0x1000 {
		t.Errorf("Expected 0x1000 but got 0x%04X", reg)
	}
}

func TestGetFlagZ(t *testing.T) {
	cpu := CPU{}
	cpu.setRegByte(F, 0x00)
	if cpu.getFlagZ() != false {
		t.Errorf("Flag Z expected false")
	}

	cpu.setRegByte(F, 0x80)
	if cpu.getFlagZ() != true {
		t.Errorf("Flag Z expected true")
	}
}

func TestSetFlagZ(t *testing.T) {
	cpu := CPU{}
	cpu.setRegByte(F, 0x00)
	cpu.setFlagZ(true)
	if cpu.getFlagZ() != true {
		t.Error("Flag Z expected true")
	}
	cpu.setFlagZ(false)
	if cpu.getFlagZ() != false {
		t.Error("Flag Z expected false")
	}
}

func TestINC(t *testing.T) {
	cpu := CPU{}
	cpu.setRegByte(E, 0xFF)
	cpu.op_INC_E()

	if cpu.getRegByte(E) != 0x00 {
		t.Errorf("Expected 0x00 but got 0x%02X", cpu.getRegByte(E))
	}

	if cpu.getFlagZ() != true {
		t.Error("Expected Z to be false")
	}
}

func TestSetFlags(t *testing.T) {
	cpu := CPU{}

	cpu.setRegByte(F, 0x00)
	cpu.setFlagZ(true)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	if cpu.getFlagZ() != true {
		t.Error("Z was not true")
	}

	cpu.setRegByte(F, 0x00)
	cpu.setFlagN(true)
	cpu.setFlagZ(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	if cpu.getFlagN() != true {
		t.Error("N was not true")
	}

	cpu.setRegByte(F, 0x00)
	cpu.setFlagH(true)
	cpu.setFlagZ(false)
	cpu.setFlagN(false)
	cpu.setFlagC(false)
	if cpu.getFlagH() != true {
		t.Error("H was not true")
	}

	cpu.setRegByte(F, 0x00)
	cpu.setFlagC(true)
	cpu.setFlagZ(false)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	if cpu.getFlagC() != true {
		t.Error("C was not true")
	}
}
