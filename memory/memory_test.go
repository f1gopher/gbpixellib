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
	if value != 0xCDEF {
		t.Errorf("Expected 0xCDEF but got 0x%X", value)
	}
}

func TestWriteShort(t *testing.T) {
	mem := CreateMemory()
	mem.WriteShort(0, 0x1234)

	lsb := mem.ReadByte(0)
	msb := mem.ReadByte(1)

	if lsb != 0x12 {
		t.Errorf("Expected 0x12 for LSB but got 0x%02X", lsb)
	}

	if msb != 0x12 {
		t.Errorf("Expected 0x34 for MSB but got 0x%02X", msb)
	}
}

func TestReadShort(t *testing.T) {
	mem := CreateMemory()
	mem.WriteByte(0, 0x12)
	mem.WriteByte(1, 0x34)

	result := mem.ReadShort(0)

	if result != 0x1234 {
		t.Errorf("Expected 0x1234 but got 0x%04X", result)
	}
}
