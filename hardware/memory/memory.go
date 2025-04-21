package memory

import (
	"errors"
	"fmt"
)

type Memory[T uint16 | uint32] interface {
	Size() int
	Offset() T
	LastAddress() T
	Reset()
	ReadBit(address T, bit uint8) bool
	ReadByte(address T) byte
	ReadShort(address T) uint16
	WriteBit(address T, bit uint8, value bool)
	WriteByte(address T, value byte)
	WriteShort(address T, value uint16)
	DumpCode() (data []uint8, startAddress T)
	DumpRegionsCode(address T) (data []uint8, err error)
	LoadData(offset T, data []uint8) error
	RawBuffer() []uint8
}

type memory[T uint16 | uint32] struct {
	name          string
	buffer        []uint8
	addressOffset T
	readonly      bool
	isLitleEndian bool
}

func CreateReadOnlyMemoryLittleEndian[T uint16 | uint32](name string, data *[]byte, addressOffset T) Memory[T] {
	return createReadOnlyMemory(name, data, addressOffset, true)
}
func CreateReadOnlyMemoryBigEndian[T uint16 | uint32](name string, data *[]byte, addressOffset T) Memory[T] {
	return createReadOnlyMemory(name, data, addressOffset, false)
}
func CreateMemoryLittleEndian[T uint16 | uint32](name string, data *[]byte, addressOffset T) Memory[T] {
	return createMemory(name, data, addressOffset, true)
}
func CreateMemoryBigEndian[T uint16 | uint32](name string, data *[]byte, addressOffset T) Memory[T] {
	return createMemory(name, data, addressOffset, false)
}

func createReadOnlyMemory[T uint16 | uint32](name string, data *[]byte, addressOffset T, isLittleEndian bool) Memory[T] {
	return &memory[T]{
		name:          name,
		buffer:        *data,
		addressOffset: addressOffset,
		readonly:      true,
		isLitleEndian: isLittleEndian,
	}
}

func createMemory[T uint16 | uint32](name string, data *[]byte, addressOffset T, isLittleEndian bool) Memory[T] {
	return &memory[T]{
		name:          name,
		buffer:        *data,
		addressOffset: addressOffset,
		readonly:      false,
		isLitleEndian: isLittleEndian,
	}
}

func (m *memory[T]) Size() int {
	return len(m.buffer)
}

func (m *memory[T]) Offset() T {
	return m.addressOffset
}

func (m *memory[T]) offset(address T) T {
	return address - m.addressOffset
}

func (m *memory[T]) LastAddress() T {
	return (m.addressOffset + T(m.Size())) - 1
}

func (m *memory[T]) Reset() {
	for x := 0; x < len(m.buffer); x++ {
		m.buffer[x] = 0x00
	}
}

func (m *memory[T]) ReadBit(address T, bit uint8) bool {
	if bit < 0 || bit > 7 {
		panic(fmt.Sprintf("Invalid bit: %d", bit))
	}

	value := m.ReadByte(address)
	return (value>>bit)&0x01 == 0x01
}

func (m *memory[T]) ReadByte(address T) byte {
	return m.buffer[m.offset(address)]
}

func (m *memory[T]) ReadShort(address T) uint16 {
	var result uint16
	var lsb, msb uint8
	if m.isLitleEndian {
		// Little endian, lsb stored first
		lsb = m.ReadByte(address)
		msb = m.ReadByte(address + 1)
	} else {
		msb = m.ReadByte(address)
		lsb = m.ReadByte(address + 1)
	}
	result = uint16(msb)
	result = result << 8
	result = result | uint16(lsb)
	return result
}

func (m *memory[T]) WriteBit(address T, bit uint8, value bool) {
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

func (m *memory[T]) WriteByte(address T, value byte) {
	if m.readonly {
		return
	}

	m.buffer[m.offset(address)] = value
}

func (m *memory[T]) WriteShort(address T, value uint16) {
	if m.readonly {
		return
	}
	var lsb, msb uint8
	if m.isLitleEndian {
		// Little endian - lsb stored first
		lsb = uint8(value)
		msb = uint8(value >> 8)
	} else {
		msb = uint8(value)
		lsb = uint8(value >> 8)
	}
	m.WriteByte(address, lsb)
	m.WriteByte(address+1, msb)
}

func (m *memory[T]) DumpCode() (data []uint8, startAddress T) {
	code := make([]uint8, len(m.buffer))
	copy(code, m.buffer)
	return code, m.addressOffset
}

func (m *memory[T]) DumpRegionsCode(address T) (data []uint8, err error) {
	code, _ := m.DumpCode()
	return code, nil
}

func (m *memory[T]) LoadData(offset T, data []uint8) error {
	if int(offset)+len(data) > len(m.buffer) {
		return errors.New(fmt.Sprintf("Unable to LoadData(%d, %d) when buffer is %d", offset, len(data), len(m.buffer)))
	}

	copy(m.buffer[offset:], data)
	return nil
}

func (m *memory[T]) RawBuffer() []uint8 {
	return m.buffer
}
