package cpu

import "testing"

type testMemory_NoAccess struct {
	test *testing.T
}

func (t *testMemory_NoAccess) ReadBit(address uint16, bit uint8) bool {
	t.test.FailNow()
	return false
}

func (t *testMemory_NoAccess) ReadByte(address uint16) uint8 {
	t.test.FailNow()
	return 0
}

func (t *testMemory_NoAccess) ReadShort(address uint16) uint16 {
	t.test.FailNow()
	return 0
}

func (t *testMemory_NoAccess) WriteByte(address uint16, value uint8) {
	t.test.FailNow()
}

func (t *testMemory_NoAccess) WriteShort(address uint16, value uint16) {
	t.test.FailNow()
}

func (t *testMemory_NoAccess) Write(address uint16, data []uint8) error {
	t.test.FailNow()
	return nil
}
