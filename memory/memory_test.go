package memory

import "testing"

func TestReadWriteByte(t *testing.T) {
	mem := CreateMemory()
	mem.WriteByte(0, 0xAB)
	value := mem.ReadByte(0)
	if value != 0xAB {
		t.Errorf("Expected 0xAB but got 0x%X", value)
	}
}

func TestReadWriteShort(t *testing.T) {
	mem := CreateMemory()
	mem.Write(0, []byte{0xCD, 0xEF})
	value := mem.ReadShort(0)
	// Byte order is reversed
	if value != 0xEFCD {
		t.Errorf("Expected 0xEFCD but got 0x%X", value)
	}
}
