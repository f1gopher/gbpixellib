package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/log"
)

type RegistersInterface interface {
	Get8(source Register) uint8
	Get16(source Register) uint16
	Get16Msb(source Register) uint8
	Get16Lsb(source Register) uint8
	Set8(target Register, value uint8)
	Set16(target Register, value uint16)
	Set16FromTwoBytes(target Register, msb uint8, lsb uint8)
	SetRegBit(target Register, bit uint8, value bool)

	GetFlag(flag RegisterFlags) bool
	SetFlag(flag RegisterFlags, value bool)

	SetIME(enabled bool)
	GetIME() bool

	SetHALT(enabled bool)
	GetHALT() bool

	Reset()
}

type MemoryInterface interface {
	Reset()

	ReadBit(address uint16, bit uint8) bool
	ReadByte(address uint16) uint8
	ReadShort(address uint16) uint16

	WriteByte(address uint16, value uint8)
	WriteShort(address uint16, value uint16)

	DisplaySetScanline(value uint8)
	DisplaySetStatus(value uint8)

	WriteDividerRegister(value uint8)
	ExecuteDMAIfPending() bool
}

type Cpu struct {
	log    *log.Log
	reg    RegistersInterface
	memory MemoryInterface

	executeOpcodePC      uint16
	prevOpcodePC         uint16
	executeOpcode        opcode
	executeOpcodesMCycle int
	prevOpcode           uint8
	isCB                 bool

	opcodes   [256]opcode
	cbOpcodes [256]opcode

	interruptHappened bool
}

func CreateCPU(log *log.Log, regs RegistersInterface, memory MemoryInterface) *Cpu {
	return &Cpu{
		log:       log,
		reg:       regs,
		memory:    memory,
		opcodes:   createOpcodesTable(),
		cbOpcodes: createCBOpcodesTable(),
	}
}

func (c *Cpu) InitForTestROM() {
	c.executeOpcode = nil

	// Test ROM recommended intial values
	// c.reg.Set8(A, 0x01)
	// c.reg.Set8(F, 0xB0)
	// c.reg.Set8(B, 0x00)
	// c.reg.Set8(C, 0x13)
	// c.reg.Set8(D, 0x00)
	// c.reg.Set8(E, 0xD8)
	// c.reg.Set8(H, 0x01)
	// c.reg.Set8(L, 0x4D)
	// c.reg.Set16(SP, 0xFFFE)
	// c.reg.Set16(PC, 0x0100)

	// Initial values for comparison system
	c.reg.Set8(A, 0x11)
	c.reg.Set8(F, 0x80)
	c.reg.Set8(B, 0x00)
	c.reg.Set8(C, 0x00)
	c.reg.Set8(D, 0xFF)
	c.reg.Set8(E, 0x56)
	c.reg.Set8(H, 0x00)
	c.reg.Set8(L, 0x0D)
	c.reg.Set16(SP, 0xFFFE)
	c.reg.Set16(PC, 0x0100)

	c.executeOpcodePC = c.reg.Get16(PC)
}

func (c *Cpu) Init() {
	c.executeOpcode = nil

	c.reg.Set16(PC, 0x0000)
	c.reg.Set16(AF, 0x01B0)
	c.reg.Set16(BC, 0x0013)
	c.reg.Set16(DE, 0x00D8)
	c.reg.Set16(HL, 0x014D)
	c.reg.Set16(SP, 0xFFFE)
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

	c.executeOpcodePC = c.reg.Get16(PC)
}

func (c *Cpu) Reset() {
	c.executeOpcode = nil
	c.executeOpcodesMCycle = 0
	c.reg.Reset()
	c.executeOpcodePC = 0
	c.prevOpcodePC = 0
	c.prevOpcode = 0
	c.isCB = false
	c.interruptHappened = false
}

func (c *Cpu) DoInterruptCycle() error {
	c.executeOpcodesMCycle++

	// If an interrupt with jump happens reset and fetch the next instruction
	c.executeOpcodesMCycle = 0

	c.prevOpcodePC = c.executeOpcodePC
	//c.executeOpcodePC = c.reg.Get16(PC)
	//c.prevOpcode = c.executeOpcode.opcode()
	//opcode := readAndIncPC(c.reg, c.memory)
	//var cbOpcode uint8 = 0x00
	//if opcode == 0xCB {
	//	cbOpcode = readAndIncPC(c.reg, c.memory)
	//}
	//c.isCB = opcode == 0xCB
	//
	//var err error
	//c.executeOpcode, err = c.getOpcode(opcode, cbOpcode)

	c.interruptHappened = true

	// If we had an error fetching the opcode fail
	return nil
}

func (c *Cpu) ExecuteMCycle() (breakpoint bool, instructionCompleted bool, opcode uint8, opcodeDescription string, err error) {
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

		c.executeOpcodesMCycle = 0

		if err != nil {
			return false, false, 0, "", err
		}

		// TODO - not sure is this is right, hack to match the other app for testing

		//return false, false, err
	}

	// TODO - is this right?
	if c.interruptHappened {
		c.interruptHappened = false

		c.prevOpcodePC = c.executeOpcodePC
		c.executeOpcodePC = c.reg.Get16(PC)
		c.prevOpcode = c.executeOpcode.opcode()

		opcode := readAndIncPC(c.reg, c.memory)
		var cbOpcode uint8 = 0x00
		if opcode == 0xCB {
			cbOpcode = readAndIncPC(c.reg, c.memory)
		}

		var err error
		// Assign it straight to execute so it will be ready for the next cycle
		c.executeOpcode, err = c.getOpcode(opcode, cbOpcode)

		c.executeOpcodesMCycle = 0

		if err != nil {
			return false, false, 0, "", err
		}
	}

	// We already have an opcode so do an excute on that opcode
	completed, err := c.executeOpcode.doCycle(
		c.executeOpcodesMCycle+1,
		c.reg,
		c.memory)

	if err != nil {
		return false, completed, 0, "", errors.Join(errors.New(fmt.Sprintf("Opcode %s cycle %d", c.executeOpcode.name(), c.executeOpcodesMCycle)), err)
	}

	c.executeOpcodesMCycle++
	var description string
	breakpointHit := false

	// TODO - directly check registers/memory for breakpoints here (and clear
	// before execution)

	// If the current opcode has finished then do a fetch for the next opcode
	// in the same cycle (overlapping execute and fetch)
	if completed {
		c.executeOpcodesMCycle = 0

		description = c.executeOpcode.name()

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
			return breakpointHit, completed, 0, "", err
		}
	}

	return breakpointHit, completed, c.executeOpcode.opcode(), description, nil
}

func (c *Cpu) GetOpcode() string {
	if c.executeOpcode == nil {
		return "N/A"
	}

	return c.executeOpcode.name()
}

func (c *Cpu) GetNextOpcode() (opcode uint8, isCB bool) {
	if c.executeOpcode == nil {
		return 0, false
	}

	return c.executeOpcode.opcode(), c.isCB
}

func (c *Cpu) GetPrevOpcode() uint8 {
	return c.prevOpcode
}

func (c *Cpu) GetOpcodePC() uint16 {
	return c.executeOpcodePC
}

func (c *Cpu) GetPrevOpcodePC() uint16 {
	return c.prevOpcodePC
}

func (c *Cpu) GetOpcodeInfo(opcodeId uint8, cbOpcode uint8) (name string, length uint8) {
	opcode, err := c.getOpcode(opcodeId, cbOpcode)

	if err != nil {
		return "Unknown", 1
	}

	if opcode.length() == 0 {
		panic(fmt.Sprintf("Opcode %s doesn't have a length set", opcode.name()))
	}

	return opcode.name(), opcode.length()
}

func (c *Cpu) getNextOpcode() string {
	number := c.memory.ReadByte(c.reg.Get16(PC))
	cbOpcode := c.memory.ReadByte(c.reg.Get16(PC) + 1)

	opcode, err := c.getOpcode(number, cbOpcode)

	if err != nil {
		return err.Error()
	}

	return opcode.name()
}

func (c *Cpu) getOpcode(opcode uint8, cbOpcode uint8) (executer opcode, err error) {

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
