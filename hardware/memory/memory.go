package memory

import (
	"errors"
	"fmt"
)

type Memory interface {
	Size() int
	Reset()
	ReadBit(address uint16, bit uint8) bool
	ReadByte(address uint16) byte
	ReadShort(address uint16) uint16
	WriteBit(address uint16, bit uint8, value bool)
	WriteByte(address uint16, value byte)
	WriteShort(address uint16, value uint16)
	DumpCode() (data []uint8, startAddress uint16)
	LoadData(offset uint16, data []uint8) error
}

type memory struct {
	name          string
	buffer        []uint8
	addressOffset uint16
	readonly      bool
}

func CreateReadOnlyMemory(name string, data *[]byte, addressOffset uint16) Memory {
	return &memory{
		name:          name,
		buffer:        *data,
		addressOffset: addressOffset,
		readonly:      true,
	}
}

func CreateMemory(name string, data *[]byte, addressOffset uint16) Memory {
	return &memory{
		name:          name,
		buffer:        *data,
		addressOffset: addressOffset,
		readonly:      false,
	}
}

func (m *memory) Size() int {
	return len(m.buffer)
}

func (m *memory) offset(address uint16) uint16 {
	return address - m.addressOffset
}

func (m *memory) Reset() {
	for x := 0; x < len(m.buffer); x++ {
		m.buffer[x] = 0x00
	}
}

func (m *memory) ReadBit(address uint16, bit uint8) bool {
	if bit < 0 || bit > 7 {
		panic(fmt.Sprintf("Invalid bit: %d", bit))
	}

	value := m.ReadByte(address)
	return (value>>bit)&0x01 == 0x01
}

func (m *memory) ReadByte(address uint16) byte {
	return m.buffer[m.offset(address)]
}

func (m *memory) ReadShort(address uint16) uint16 {
	var result uint16
	lsb := m.ReadByte(address)
	msb := m.ReadByte(address + 1)
	// Little endian, lsb stored first
	result = uint16(msb)
	result = result << 8
	result = result | uint16(lsb)
	return result
}

func (m *memory) WriteBit(address uint16, bit uint8, value bool) {
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

func (m *memory) WriteByte(address uint16, value byte) {
	if m.readonly {
		return
	}

	m.buffer[m.offset(address)] = value
}

func (m *memory) WriteShort(address uint16, value uint16) {
	if m.readonly {
		return
	}

	lsb := uint8(value)
	msb := uint8(value >> 8)
	// Little endian - lsb stored first
	m.WriteByte(address, lsb)
	m.WriteByte(address+1, msb)
}

func (m *memory) DumpCode() (data []uint8, startAddress uint16) {
	code := make([]uint8, len(m.buffer))
	copy(code, m.buffer)
	return code, m.addressOffset
}

func (m *memory) LoadData(offset uint16, data []uint8) error {
	if int(offset)+len(data) > len(m.buffer) {
		return errors.New(fmt.Sprintf("Unable to LoadData(%d, %d) when buffer is %d", offset, len(data), len(m.buffer)))
	}

	copy(m.buffer[offset:], data)
	return nil
}
