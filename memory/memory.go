package memory

import (
	"errors"
	"fmt"
	"os"

	"github.com/f1gopher/gbpixellib/log"
)

const memorySize = 0xFFFF

type Memory struct {
	buffer []byte
	log    *log.Log

	abc []byte

	breakAddress uint16
}

func CreateMemory(log *log.Log) *Memory {
	return &Memory{
		log:          log,
		buffer:       make([]byte, memorySize+1),
		abc:          make([]byte, 256),
		breakAddress: 0x0000,
	}
}

func (m *Memory) Reset() {
	for x := 0; x < len(m.buffer); x++ {
		m.buffer[x] = 0x00
	}

	for x := 0; x < len(m.abc); x++ {
		m.abc[x] = 0x00
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
	// TODO - Only used for debug ROMs
	//if address == 0xFF44 {
	//	return 0x90
	//}
	return m.buffer[address]
}

func (m *Memory) ReadShort(address uint16) uint16 {
	var result uint16
	lsb := m.buffer[address]
	msb := m.buffer[address+1]
	// Little endian, lsb stored first
	result = uint16(msb)
	result = result << 8
	result = result | uint16(lsb)
	return result
}

func (m *Memory) WriteByte(address uint16, value byte) {
	if address != 0xFF41 {
		m.log.Debug(fmt.Sprintf("[RAM: 0x%04X 0x%02X]", address, value))
	}

	// If the CPU tries to write to the address (regardless of value) we
	// set the value to zero
	if address == 0xFF44 {
		//m.buffer[address] = 0
		m.write(address, value)
		return
	}

	// Trigger DMA transfer
	if address == 0xFF46 {
		// TODO - this needs to take 160 cycles
		dmaAddress := uint16(value) << 8
		var i uint16 = 0
		for i = 0; i < 0x9F; i++ {
			//m.WriteByte(0xFE00+i, m.ReadByte(dmaAddress+i))
			m.write(0xFE00+i, m.ReadByte(dmaAddress+i))
		}

		// TODO - is this right?
		return
	}

	//m.buffer[address] = value
	m.write(address, value)

	if address >= 0xE00 && address < 0xFE00 {
		//m.buffer[address-0x2000] = value
		m.write(address-0x2000, value)
	}
}

func (m *Memory) DisplaySetScanline(value byte) {
	// Only used by the display
	//m.buffer[0xFF44] = value
	m.write(0xFF44, value)
}

func (m *Memory) DumpTiles() {
	copy(m.abc, m.buffer[0x9800:0x9BFF])
}

func (m *Memory) WriteShort(address uint16, value uint16) {
	lsb := byte(value)
	msb := byte(value >> 8)
	// Little endian - lsb stored first
	//m.buffer[address] = lsb
	//m.buffer[address+1] = msb
	m.write(address, lsb)
	m.write(address+1, msb)

	if address >= 0xE00 && address < 0xFE00 {
		//m.buffer[address-0x2000] = lsb
		//m.buffer[(address+1)-0x2000] = msb
		m.write(address-0x2000, lsb)
		m.write((address+1)-0x2000, msb)
	}

	m.log.Debug(fmt.Sprintf("[RAM: 0x%04X 0x%02X]", address+1, msb))
	m.log.Debug(fmt.Sprintf("[RAM: 0x%04X 0x%02X]", address, lsb))
}

func (m *Memory) Write(address uint16, data []byte) error {
	if int(address)+len(data) > memorySize {
		return errors.New("Write buffer will exceed memory range")
	}

	if address >= m.breakAddress && address+uint16(len(data)) <= m.breakAddress {
		panic("breakpoint")
	}

	copy(m.buffer[address:len(data)], data)

	return nil

}

func (m *Memory) write(address uint16, value byte) error {
	if address > memorySize {
		return errors.New("Write buffer will exceed memory range")
	}

	//if address == m.breakAddress && value != 0 {
	//	panic("breakpoint")
	//}

	m.buffer[address] = value

	return nil
}

func (m *Memory) LoadBios(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Join(errors.New("Failed to load bios"), err)
	}

	return m.Write(0, data)
}

func (m *Memory) LoadRom(path string) error {
	if len(path) == 0 {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Join(errors.New("Failed to load ROM"), err)
	}

	// TODO - check not too big

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

func (m *Memory) DumpCode() []byte {
	code := make([]byte, 0x8000)
	copy(code, m.buffer[:0x8000])
	return code
}

func dumpChar(value byte) string {
	if value < 32 || value > 126 {
		return "."
	}

	return string(value)
}
