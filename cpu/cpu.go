package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/log"
)

type registersInterface interface {
	Get8(source register) uint8
	Get16(source register) uint16
	get16Msb(source register) uint8
	get16Lsb(source register) uint8
	set8(target register, value uint8)
	set16(target register, value uint16)
	set16FromTwoBytes(target register, msb uint8, lsb uint8)
	setRegBit(target register, bit uint8, value bool)

	GetFlag(flag registerFlags) bool
	setFlag(flag registerFlags, value bool)

	SetIME(enabled bool)
	GetIME() bool

	reset()
}

type memoryInterface interface {
	ReadBit(address uint16, bit uint8) bool
	ReadByte(address uint16) uint8
	ReadShort(address uint16) uint16

	WriteByte(address uint16, value uint8)
	WriteShort(address uint16, value uint16)
	Write(address uint16, data []uint8) error
}

type CPU2 struct {
	log    *log.Log
	reg    registersInterface
	memory memoryInterface

	executeOpcodePC     uint16
	prevOpcodePC        uint16
	executeOpcode       opcode
	executeOpcodesCycle int
	prevOpcode          uint8
	isCB                bool

	opcodes   [256]opcode
	cbOpcodes [256]opcode
}

func CreateCPU2(log *log.Log, regs registersInterface, memory memoryInterface) *CPU2 {
	return &CPU2{
		log:       log,
		reg:       regs,
		memory:    memory,
		opcodes:   createOpcodesTable(),
		cbOpcodes: createCBOpcodesTable(),
	}
}

func (c *CPU2) InitForTestROM() {
	c.executeOpcode = nil

	c.reg.set8(A, 0x01)
	c.reg.set8(F, 0xB0)
	c.reg.set8(B, 0x00)
	c.reg.set8(C, 0x13)
	c.reg.set8(D, 0x00)
	c.reg.set8(E, 0xD8)
	c.reg.set8(H, 0x01)
	c.reg.set8(L, 0x4D)
	c.reg.set16(SP, 0xFFFE)
	c.reg.set16(PC, 0x0100)
}

func (c *CPU2) Init() {
	c.executeOpcode = nil

	c.reg.set16(PC, 0x0000)
	c.reg.set16(AF, 0x01B0)
	c.reg.set16(BC, 0x0013)
	c.reg.set16(DE, 0x00D8)
	c.reg.set16(HL, 0x014D)
	c.reg.set16(SP, 0xFFFE)
	c.memory.WriteByte(0xFF05, 0x00)
	c.memory.WriteByte(0xFF06, 0x00)
	c.memory.WriteByte(0xFF07, 0x00)
	c.memory.WriteByte(0xFF10, 0x80)
	c.memory.WriteByte(0xFF11, 0xBF)
	c.memory.WriteByte(0xFF12, 0xF3)
	c.memory.WriteByte(0xFF14, 0xBF)
	c.memory.WriteByte(0xFF16, 0x3F)
	c.memory.WriteByte(0xFF17, 0x00)
	c.memory.WriteByte(0xFF19, 0xBF)
	c.memory.WriteByte(0xFF1A, 0x7F)
	c.memory.WriteByte(0xFF1B, 0xFF)
	c.memory.WriteByte(0xFF1C, 0x9F)
	c.memory.WriteByte(0xFF1E, 0xBF)
	c.memory.WriteByte(0xFF20, 0xFF)
	c.memory.WriteByte(0xFF21, 0x00)
	c.memory.WriteByte(0xFF22, 0x00)
	c.memory.WriteByte(0xFF23, 0xBF)
	c.memory.WriteByte(0xFF24, 0x77)
	c.memory.WriteByte(0xFF25, 0xF3)
	c.memory.WriteByte(0xFF26, 0xF1)
	c.memory.WriteByte(0xFF40, 0x91)
	c.memory.WriteByte(0xFF42, 0x00)
	c.memory.WriteByte(0xFF43, 0x00)
	c.memory.WriteByte(0xFF45, 0x00)
	c.memory.WriteByte(0xFF47, 0xFC)
	c.memory.WriteByte(0xFF48, 0xFF)
	c.memory.WriteByte(0xFF49, 0xFF)
	c.memory.WriteByte(0xFF4A, 0x00)
	c.memory.WriteByte(0xFF4B, 0x00)
	c.memory.WriteByte(0xFFFF, 0x00)
}

func (c *CPU2) Reset() {
	c.executeOpcode = nil
	c.executeOpcodesCycle = 0
	c.reg.reset()
	c.executeOpcodePC = 0
	c.prevOpcodePC = 0
	c.prevOpcode = 0
	c.isCB = false
}

func (c *CPU2) ExecuteCycle() (breakpoint bool, instructionCompleted bool, err error) {
	// On startup we will have no fetched opcode so
	// fetch a nw opcode then end the cycle
	if c.executeOpcode == nil {
		c.executeOpcodePC = c.reg.Get16(PC)
		opcode := readAndIncPC(c.reg, c.memory)
		var cbOpcode uint8 = 0x00
		if opcode == 0xCB {
			cbOpcode = readAndIncPC(c.reg, c.memory)
		}

		var err error
		// Assign it straight to execute so it will be ready for the next cycle
		c.executeOpcode, err = c.getOpcode(opcode, cbOpcode)

		// The first cycle for an opcode is always the fetch
		c.executeOpcodesCycle = 1

		return false, false, err
	}

	// We already have an opcode so do an excute on that opcode
	completed, err := c.executeOpcode.doCycle(
		c.executeOpcodesCycle,
		c.reg,
		c.memory)

	if err != nil {
		return false, completed, errors.Join(errors.New(fmt.Sprintf("Opcode %s cycle %d", c.executeOpcode.name(), c.executeOpcodesCycle)), err)
	}

	c.executeOpcodesCycle++

	// TODO - directly check registers/memory for breakpoints here (and clear
	// before execution)

	// If the current opcode has finished then do a fetch for the next opcode
	// in the same cycle (overlapping execute and fetch)
	if completed {
		// The first cycle for an copcode is always the fetch so execute starts
		// on cycle 1
		c.executeOpcodesCycle = 1

		c.prevOpcodePC = c.executeOpcodePC
		c.executeOpcodePC = c.reg.Get16(PC)
		c.prevOpcode = c.executeOpcode.opcode()
		opcode := readAndIncPC(c.reg, c.memory)
		var cbOpcode uint8 = 0x00
		if opcode == 0xCB {
			cbOpcode = readAndIncPC(c.reg, c.memory)
		}
		c.isCB = opcode == 0xCB

		var err error
		c.executeOpcode, err = c.getOpcode(opcode, cbOpcode)

		// If we had an error fetching the opcode fail
		if err != nil {
			return false, completed, err
		}
	}

	return false, completed, nil
}

func (c *CPU2) GetOpcode() string {
	if c.executeOpcode == nil {
		return "N/A"
	}

	return c.executeOpcode.name()
}

func (c *CPU2) GetNextOpcode() (opcode uint8, isCB bool) {
	if c.executeOpcode == nil {
		return 0, false
	}

	return c.executeOpcode.opcode(), c.isCB
}

func (c *CPU2) GetPrevOpcode() uint8 {
	return c.prevOpcode
}

func (c *CPU2) GetOpcodePC() uint16 {
	return c.executeOpcodePC
}

func (c *CPU2) GetPrevOpcodePC() uint16 {
	return c.prevOpcodePC
}

func (c *CPU2) GetOpcodeInfo(opcodeId uint8, cbOpcode uint8) (name string, length uint8) {
	opcode, err := c.getOpcode(opcodeId, cbOpcode)

	if err != nil {
		return "Unknown", 1
	}

	if opcode.length() == 0 {
		panic(fmt.Sprintf("Opcode %s doesn't have a length set", opcode.name()))
	}

	return opcode.name(), opcode.length()
}

func (c *CPU2) getNextOpcode() string {
	number := c.memory.ReadByte(c.reg.Get16(PC))
	cbOpcode := c.memory.ReadByte(c.reg.Get16(PC) + 1)

	opcode, err := c.getOpcode(number, cbOpcode)

	if err != nil {
		return err.Error()
	}

	return opcode.name()
}

func (c *CPU2) getOpcode(opcode uint8, cbOpcode uint8) (executer opcode, err error) {

	if opcode == 0xCB {
		executer := c.cbOpcodes[cbOpcode]

		if executer == nil {
			return nil, errors.New(fmt.Sprintf("Unsupported CB opcode: 0x%02X", cbOpcode))
		}

		return executer, nil
	}

	executor := c.opcodes[opcode]
	if executor == nil {
		return nil, errors.New(fmt.Sprintf("Unsupported opcode: 0x%02X", opcode))
	}

	return executor, nil
}
