package cpu

import (
	"testing"

	"golang.org/x/exp/slices"
)

type testRegisters_UseOneRegister struct {
	Registers

	test *testing.T

	allowedReg   register
	allowedFlags []registerFlags
}

func (t *testRegisters_UseOneRegister) Get8(source register) uint8 {
	if source != t.allowedReg {
		t.test.FailNow()
	}

	return t.Registers.Get8(source)
}

func (t *testRegisters_UseOneRegister) Get16(source register) uint16 {
	if source != t.allowedReg {
		t.test.FailNow()
	}

	return t.Registers.Get16(source)
}

func (t *testRegisters_UseOneRegister) get16Msb(source register) uint8 {
	if source != t.allowedReg {
		t.test.FailNow()
	}

	return t.Registers.get16Msb(source)
}

func (t *testRegisters_UseOneRegister) get16Lsb(source register) uint8 {
	if source != t.allowedReg {
		t.test.FailNow()
	}

	return t.Registers.get16Lsb(source)
}

func (t *testRegisters_UseOneRegister) set8(target register, value uint8) {
	if target != t.allowedReg {
		t.test.FailNow()
	}

	t.Registers.set8(target, value)
}

func (t *testRegisters_UseOneRegister) set16(target register, value uint16) {
	if target != t.allowedReg {
		t.test.FailNow()
	}

	t.Registers.set16(target, value)
}

func (t *testRegisters_UseOneRegister) GetFlag(flag registerFlags) bool {
	if !slices.Contains(t.allowedFlags, flag) {
		t.test.FailNow()
	}

	return t.Registers.GetFlag(flag)
}

func (t *testRegisters_UseOneRegister) setFlag(flag registerFlags, value bool) {
	if !slices.Contains(t.allowedFlags, flag) {
		t.test.FailNow()
	}

	t.Registers.setFlag(flag, value)
}
