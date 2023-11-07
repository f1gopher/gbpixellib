package memory

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/log"
)

const memorySize = 0xFFFF

type inputOutput interface {
	ReadDirectional() uint8
	ReadStandard() uint8
}

type Memory struct {
	buffer []uint8
	log    *log.Log

	abc []uint8

	breakAddress uint16

	dmaPending bool
	dmaAddress uint16

	io inputOutput
}

func CreateMemory(log *log.Log) *Memory {
	return &Memory{
		log:          log,
		buffer:       make([]uint8, memorySize+1),
		abc:          make([]uint8, 256),
		breakAddress: 0x0000,
		io:           nil,
	}
}

func (m *Memory) SetIO(io inputOutput) {
	m.io = io
}

func (m *Memory) Reset() {
	for x := 0; x < len(m.buffer); x++ {
		m.buffer[x] = 0x00
	}

	for x := 0; x < len(m.abc); x++ {
		m.abc[x] = 0x00
	}

	// For controller
	m.buffer[0xFF00] = 0x3F
}

func (m *Memory) ReadBit(address uint16, bit uint8) bool {
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

func (m *Memory) WriteBit(address uint16, bit uint8, value bool) {
	if bit < 0 || bit > 7 {
		panic(fmt.Sprintf("Invalid bit: %d", bit))
	}

	if address == 0xFFFF {
		fmt.Println("asd")
	}

	currentByte := m.ReadByte(address)
	SetBit(currentByte, bit, value)
	m.WriteByte(address, currentByte)
}

func (m *Memory) WriteByte(address uint16, value byte) {
	// This is cartridge ROM and we can't write to it
	if address <= 0x7FFF {
		return
	}

	// If the CPU tries to write to the address (regardless of value) we
	// set the value to zero
	if address == 0xFF44 {
		//m.buffer[address] = 0
		m.write(address, value)
		return
	}

	// TODO - Probably should be handled haigher up in the system code

	// Trigger DMA transfer
	if address == 0xFF46 {
		m.dmaPending = true
		m.dmaAddress = uint16(value) << 8
		return
	}

	// Controller
	if address == 0xFF00 {
		P14 := (value >> 4) & 0x01
		P15 := (value >> 5) & 0x01

		current := 0xFF & (value | 0b11001111)

		if P14 == 0 {
			current &= m.io.ReadDirectional()
		}

		if P15 == 0 {
			current &= m.io.ReadStandard()
		}

		m.write(address, current)
		return
	}

	// Timer register
	if address == 0xFF05 {
		// When timer overflow from FF to 00 triogger interrupt
		if m.ReadByte(address) == 0xFF && value == 0x00 {
		}
	}

	//m.buffer[address] = value
	m.write(address, value)

	// Echo RAM - might not need this
	//
	// 0xE000 - 0xFDFF: Echo RAM
	//This section of memory directly mirrors the working RAM section - meaning if you write into the first address of working RAM (0xC000), the same value will appear in the first spot of echo RAM (0xE000). Nintendo actively discouraged developers from using this area of memory and as such we can just pretend it doesn't exist.
	//if address >= 0xC000 && address <= 0xDFFF {
	//	m.write(address+0x2000, value)
	//}
}

func (m *Memory) ExecuteDMAIfPending() bool {
	if !m.dmaPending {
		return false
	}

	// TODO - this needs to take 160 cycles
	//dmaAddress := uint16(value) << 8
	var i uint16 = 0
	for i = 0; i < 0x9F; i++ {
		//m.Writeuint8(0xFE00+i, m.ReadByte(dmaAddress+i))
		m.write(0xFE00+i, m.ReadByte(m.dmaAddress+i))
	}

	// TODO - is this right?

	m.dmaPending = false

	return true
}

func (m *Memory) DisplaySetScanline(value uint8) {
	// Only used by the display
	//m.buffer[0xFF44] = value
	m.write(0xFF44, value)
}

func (m *Memory) DumpTiles() {
	copy(m.abc, m.buffer[0x9800:0x9BFF])
}

func (m *Memory) WriteShort(address uint16, value uint16) {
	lsb := uint8(value)
	msb := uint8(value >> 8)
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

func (m *Memory) Write(address uint16, data []uint8) error {
	if int(address)+len(data) > memorySize {
		return errors.New("Write buffer will exceed memory range")
	}

	if address >= m.breakAddress && address+uint16(len(data)) <= m.breakAddress {
		panic("breakpoint")
	}

	copy(m.buffer[address:len(data)], data)

	return nil

}

func (m *Memory) write(address uint16, value uint8) error {
	m.buffer[address] = value
	return nil
}

func (m *Memory) LoadBios(data *[]byte) error {
	return m.Write(0, *data)
}

func (m *Memory) LoadRom(data *[]byte) error {
	// check not too big
	if len(*data) > memorySize {
		return errors.New(fmt.Sprintf("ROM size is bigger than memory: %d", len(*data)))
	}

	return m.Write(0, *data)
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

func (m *Memory) DumpCode() []uint8 {
	code := make([]uint8, 0x8000)
	copy(code, m.buffer[:0x8000])
	return code
}

func dumpChar(value uint8) string {
	if value < 32 || value > 126 {
		return "."
	}

	return string(value)
}
