package cpu

import (
	"testing"

	"golang.org/x/exp/slices"
)

type testRegisters_UseOneRegister struct {
	Registers

	test *testing.T

	allowedReg   Register
	allowedFlags []RegisterFlags
}

func (t *testRegisters_UseOneRegister) Get8(source Register) uint8 {
	if source != t.allowedReg {
		t.test.FailNow()
	}

	return t.Registers.Get8(source)
}

func (t *testRegisters_UseOneRegister) Get16(source Register) uint16 {
	if source != t.allowedReg {
		t.test.FailNow()
	}

	return t.Registers.Get16(source)
}

func (t *testRegisters_UseOneRegister) Get16Msb(source Register) uint8 {
	if source != t.allowedReg {
		t.test.FailNow()
	}

	return t.Registers.Get16Msb(source)
}

func (t *testRegisters_UseOneRegister) Get16Lsb(source Register) uint8 {
	if source != t.allowedReg {
		t.test.FailNow()
	}

	return t.Registers.Get16Lsb(source)
}

func (t *testRegisters_UseOneRegister) Set8(target Register, value uint8) {
	if target != t.allowedReg {
		t.test.FailNow()
	}

	t.Registers.Set8(target, value)
}

func (t *testRegisters_UseOneRegister) Set16(target Register, value uint16) {
	if target != t.allowedReg {
		t.test.FailNow()
	}

	t.Registers.Set16(target, value)
}

func (t *testRegisters_UseOneRegister) GetFlag(flag RegisterFlags) bool {
	if !slices.Contains(t.allowedFlags, flag) {
		t.test.FailNow()
	}

	return t.Registers.GetFlag(flag)
}

func (t *testRegisters_UseOneRegister) SetFlag(flag RegisterFlags, value bool) {
	if !slices.Contains(t.allowedFlags, flag) {
		t.test.FailNow()
	}

	t.Registers.SetFlag(flag, value)
}
