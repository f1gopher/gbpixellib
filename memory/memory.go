package memory

import (
	"fmt"
)

type Memory struct {
	name          string
	buffer        []uint8
	addressOffset uint16
	readonly      bool
}

func CreateReadOnlyMemory(name string, data *[]byte, addressOffset uint16) *Memory {
	return &Memory{
		name:          name,
		buffer:        *data,
		addressOffset: addressOffset,
		readonly:      true,
	}
}

func CreateMemory(name string, data *[]byte, addressOffset uint16) *Memory {
	return &Memory{
		name:          name,
		buffer:        *data,
		addressOffset: addressOffset,
		readonly:      false,
	}
}

func (m *Memory) size() int {
	return len(m.buffer)
}

func (m *Memory) offset(address uint16) uint16 {
	return address - m.addressOffset
}

func (m *Memory) Reset() {
	for x := 0; x < len(m.buffer); x++ {
		m.buffer[x] = 0x00
	}
}

func (m *Memory) ReadBit(address uint16, bit uint8) bool {
	if bit < 0 || bit > 7 {
		panic(fmt.Sprintf("Invalid bit: %d", bit))
	}

	value := m.ReadByte(address)
	return (value>>bit)&0x01 == 0x01
}

func (m *Memory) ReadByte(address uint16) byte {
	return m.buffer[m.offset(address)]
}

func (m *Memory) ReadShort(address uint16) uint16 {
	var result uint16
	lsb := m.ReadByte(address)
	msb := m.ReadByte(address + 1)
	// Little endian, lsb stored first
	result = uint16(msb)
	result = result << 8
	result = result | uint16(lsb)
	return result
}

func (m *Memory) WriteBit(address uint16, bit uint8, value bool) {
	if m.readonly {
		return
	}

	if bit < 0 || bit > 7 {
		panic(fmt.Sprintf("Invalid bit: %d", bit))
	}

	currentByte := m.ReadByte(address)
	SetBit(currentByte, bit, value)
	m.WriteByte(address, currentByte)
}

func (m *Memory) WriteByte(address uint16, value byte) {
	if m.readonly {
		return
	}

	m.buffer[m.offset(address)] = value
}

func (m *Memory) WriteShort(address uint16, value uint16) {
	if m.readonly {
		return
	}

	lsb := uint8(value)
	msb := uint8(value >> 8)
	// Little endian - lsb stored first
	m.WriteByte(address, lsb)
	m.WriteByte(address+1, msb)
}

func (m *Memory) DumpCode() (data []uint8, startAddress uint16) {
	code := make([]uint8, len(m.buffer))
	copy(code, m.buffer)
	return code, m.addressOffset
}
