package memory

import (
	"errors"
	"fmt"
	"os"
)

const memorySize = 0xFFFF

type Memory struct {
	buffer []byte
}

func CreateMemory() *Memory {
	return &Memory{
		buffer: make([]byte, memorySize+1),
	}
}

func (m *Memory) ReadBit(address uint16, bit byte) bool {
	if bit < 0 || bit > 7 {
		panic(fmt.Sprintf("Invalid bit: %d", bit))
	}

	value := m.ReadByte(address)
	return (value>>bit)&0x01 == 0x01
}

func (m *Memory) ReadByte(address uint16) byte {
	// TODO - remove after implementing screen
	// Only used for debug ROMs
	if address == 0xFF44 {
		return 0x90
	}
	return m.buffer[address]
}

func (m *Memory) ReadShort(address uint16) uint16 {
	// TODO - the bytes are in the opposite endian/order
	var result uint16
	result = uint16(m.buffer[address])
	//	result = result << 8
	tmp := uint16(m.buffer[address+1])
	tmp = tmp << 8
	result = result | tmp //uint16(m.buffer[address+1])
	return result
}

func (m *Memory) WriteByte(address uint16, value byte) {
	m.buffer[address] = value
}

func (m *Memory) WriteShort(address uint16, value uint16) {
	// Write bytes in opposite order
	m.buffer[address] = byte(value)
	m.buffer[address+1] = byte(value >> 8)
}

func (m *Memory) Write(address uint16, data []byte) error {
	if int(address)+len(data) > memorySize {
		return errors.New("Write buffer will exceed memory range")
	}

	copy(m.buffer[address:len(data)], data)

	return nil

}

func (m *Memory) LoadBios(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Join(errors.New("Failed to load bios"), err)
	}

	return m.Write(0, data)
}

func (m *Memory) DumpBios() {
	for x := 0; x < 256; x += 8 {
		fmt.Printf("%05d 0x%04X    %02X %02X %02X %02X %02X %02X %02X %02X    %s%s%s%s%s%s%s%s\n",
			x,
			x,
			m.buffer[x],
			m.buffer[x+1],
			m.buffer[x+2],
			m.buffer[x+3],
			m.buffer[x+4],
			m.buffer[x+5],
			m.buffer[x+6],
			m.buffer[x+7],
			dumpChar(m.buffer[x]),
			dumpChar(m.buffer[x+1]),
			dumpChar(m.buffer[x+2]),
			dumpChar(m.buffer[x+3]),
			dumpChar(m.buffer[x+4]),
			dumpChar(m.buffer[x+5]),
			dumpChar(m.buffer[x+6]),
			dumpChar(m.buffer[x+7]))
	}
}

func dumpChar(value byte) string {
	if value < 32 || value > 126 {
		return "."
	}

	return string(value)
}
