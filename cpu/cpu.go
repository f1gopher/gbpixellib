package cpu

import (
	"errors"
	"fmt"
	"os"

	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
)

type register int

const (
	AF register = iota
	BC
	DE
	HL
	SP
	PC
	A
	F
	B
	C
	D
	E
	H
	L
)

func (r register) String() string {
	return [...]string{"AF", "BC", "DE", "HL", "SP", "PC", "A", "F", "B", "C", "D", "E", "H", "L"}[r]
}

type CPU struct {
	memory *memory.Memory
	regAF  uint16
	regBC  uint16
	regDE  uint16
	regHL  uint16
	regSP  uint16
	regPC  uint16

	opcodes     [256]func() int
	opcodesUsed [256]bool

	//debugLogFile *os.File

	log *log.Log

	count int
}

func CreateCPU(log *log.Log, memory *memory.Memory) *CPU {
	chip := CPU{
		log:    log,
		memory: memory,
		count:  0,
	}
	chip.opcodes = opcodes(&chip)

	return &chip
}

func (c *CPU) InitForTestROM() {
	c.setRegByte(A, 0x01)
	c.setRegByte(F, 0xB0)
	c.setRegByte(B, 0x00)
	c.setRegByte(C, 0x13)
	c.setRegByte(D, 0x00)
	c.setRegByte(E, 0xD8)
	c.setRegByte(H, 0x01)
	c.setRegByte(L, 0x4D)
	c.setRegShort(SP, 0xFFFE)
	c.setRegShort(PC, 0x0100)
}

func (c *CPU) Init() {
	c.setRegShort(PC, 0x0000) //0x100)
	c.setRegShort(AF, 0x01B0)
	c.setRegShort(BC, 0x0013)
	c.setRegShort(DE, 0x00D8)
	c.setRegShort(HL, 0x014D)
	c.setRegShort(SP, 0xFFFE)
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

//func (c *CPU) SetDebugLog(file string) error {
//	var err error
//	c.debugLogFile, err = os.Create(file)
//
//	if err != nil {
//		return errors.Join(errors.New("Unable to set debug log"), err)
//	}
//
//	return nil
//}

func (c *CPU) PushAndReplacePC(newPC uint16) {

	currentPC := c.getRegShort(PC)

	c.regSP -= 2
	c.memory.WriteShort(c.regSP, currentPC)

	c.setRegShort(PC, newPC)
}

func (c *CPU) Tick() (cycles int, err error) {
	//c.debugLog()

	opcode := c.memory.ReadByte(c.regPC)

	//fmt.Printf("%d - Op: 0x%02X ", c.count, opcode)

	c.opcodesUsed[opcode] = true

	executor, err := c.getOpcode(opcode)

	if err != nil {
		return 0, err
	}

	// TODO - is this the best place to increment the program counter?
	c.regPC++
	cycles = executor()

	c.count++

	return cycles, nil
}

func (c *CPU) getOpcode(opcode byte) (executer func() int, err error) {

	executor := c.opcodes[opcode]
	if executor == nil {
		return nil, errors.New(fmt.Sprintf("Unsupported opcode: 0x%02X %s", opcode, opcodeNames[opcode]))
	}

	name := opcodeNames[opcode]
	if opcode == 0xCB {
		name = name + " - " + cbNames[0]
	}

	return executor, nil
}

func (c *CPU) getRegShort(reg register) uint16 {
	switch reg {
	case AF:
		return c.regAF
	case BC:
		return c.regBC
	case DE:
		return c.regDE
	case HL:
		return c.regHL
	case SP:
		return c.regSP
	case PC:
		return c.regPC
	default:
		panic(fmt.Sprintf("Invalid register for get short : %s", reg.String()))
	}
}

func (c *CPU) getRegByte(reg register) byte {
	switch reg {
	case A:
		return getHighByte(c.regAF)
	case F:
		return getLowByte(c.regAF)
	case B:
		return getHighByte(c.regBC)
	case C:
		return getLowByte(c.regBC)
	case D:
		return getHighByte(c.regDE)
	case E:
		return getLowByte(c.regDE)
	case H:
		return getHighByte(c.regHL)
	case L:
		return getLowByte(c.regHL)
	default:
		panic(fmt.Sprintf("Invalid register for byte: %s", reg.String()))
	}
}

func (c *CPU) getFlagZ() bool { return c.getRegBit(F, 7) }
func (c *CPU) getFlagN() bool { return c.getRegBit(F, 6) }
func (c *CPU) getFlagH() bool { return c.getRegBit(F, 5) }
func (c *CPU) getFlagC() bool { return c.getRegBit(F, 4) }

func (c *CPU) getRegBit(reg register, bit int) bool {
	return memory.GetBit(c.getRegByte(reg), bit)
}

func (c *CPU) setRegShort(reg register, value uint16) {
	switch reg {
	case AF:
		c.regAF = value
	case BC:
		c.regBC = value
	case DE:
		c.regDE = value
	case HL:
		c.regHL = value
	case SP:
		c.regSP = value
	case PC:
		c.regPC = value
	default:
		panic(fmt.Sprintf("Invalid register for set short: %s", reg.String()))
	}
}

func (c *CPU) setRegByte(reg register, value byte) {
	switch reg {
	case A:
		setHighByte(&c.regAF, value)
	case F:
		setLowByte(&c.regAF, value)
	case B:
		setHighByte(&c.regBC, value)
	case C:
		setLowByte(&c.regBC, value)
	case D:
		setHighByte(&c.regDE, value)
	case E:
		setLowByte(&c.regDE, value)
	case H:
		setHighByte(&c.regHL, value)
	case L:
		setLowByte(&c.regHL, value)
	default:
		panic(fmt.Sprintf("Invalid register for set byte: %s", reg.String()))
	}
}

func getHighByte(reg uint16) byte {
	return byte(reg >> 8)
}

func getLowByte(reg uint16) byte {
	return byte(reg)
}

func setHighByte(reg *uint16, value byte) {
	*reg = *reg &^ 0xFF00
	var final uint16 = uint16(value) << 8
	*reg = *reg | final
}

func setLowByte(reg *uint16, value byte) {
	*reg = *reg &^ 0x00FF
	*reg = *reg | uint16(value)
}

func (c *CPU) setFlagZ(value bool) { c.setRegBit(F, 7, value) }
func (c *CPU) setFlagN(value bool) { c.setRegBit(F, 6, value) }
func (c *CPU) setFlagH(value bool) { c.setRegBit(F, 5, value) }
func (c *CPU) setFlagC(value bool) { c.setRegBit(F, 4, value) }

func (c *CPU) setRegBit(reg register, bit int, value bool) {
	if reg != F {
		panic(fmt.Sprintf("Unexpected register for get set bit: %s", reg.String()))
	}

	c.setRegByte(reg, memory.SetBit(c.getRegByte(reg), bit, value))
}

func (c *CPU) debugLog() {
	//msg := c.Debug()
	//	fmt.Print(msg)
	//if c.debugLogFile != nil {
	// TODO - handle the PC address being 0xFFFF so trying to read would go past the end
	p1 := c.memory.ReadByte(c.regPC)
	p2 := c.memory.ReadByte(c.regPC + 1)
	p3 := c.memory.ReadByte(c.regPC + 2)
	p4 := c.memory.ReadByte(c.regPC + 3)

	msg := fmt.Sprintf(
		"[A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X]",
		c.getRegByte(A),
		c.getRegByte(F),
		c.getRegByte(B),
		c.getRegByte(C),
		c.getRegByte(D),
		c.getRegByte(E),
		c.getRegByte(H),
		c.getRegByte(L),
		c.getRegShort(SP),
		c.getRegShort(PC),
		p1,
		p2,
		p3,
		p4)

	c.log.Debug(msg)
	//}
}

func (c *CPU) Debug() string {

	// TODO - handle the PC address being 0xFFFF so trying to read would go past the end
	p1 := c.memory.ReadByte(c.regPC)
	p2 := c.memory.ReadByte(c.regPC + 1)
	p3 := c.memory.ReadByte(c.regPC + 2)
	p4 := c.memory.ReadByte(c.regPC + 3)

	name := opcodeNames[p1]
	if p1 == 0xCB {
		// TODO - is p2 right?1
		name = name + " - " + cbNames[p2]
	}

	return fmt.Sprintf(
		"%d=>Next_Opcode:0x%02X %12s -> A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X Z:%t N:%t H:%t C:%t\n",
		c.count,
		p1,
		name,
		c.getRegByte(A),
		c.getRegByte(F),
		c.getRegByte(B),
		c.getRegByte(C),
		c.getRegByte(D),
		c.getRegByte(E),
		c.getRegByte(H),
		c.getRegByte(L),
		c.getRegShort(SP),
		c.getRegShort(PC),
		p1,
		p2,
		p3,
		p4,
		c.getFlagZ(),
		c.getFlagN(),
		c.getFlagH(),
		c.getFlagC())
}

func (c *CPU) DumpOpcodesUsed() {
	f, _ := os.Create("./opcodes_used.txt")
	for x := 0; x < len(c.opcodesUsed); x++ {
		if !c.opcodesUsed[x] {
			continue
		}

		f.WriteString(fmt.Sprintf("0x%02X\n - %s", x, opcodeNames[x]))
	}
	f.Close()
}

func (c *CPU) op_LD_B_B() int { return c.loadByte(B, B) }
func (c *CPU) op_LD_B_C() int { return c.loadByte(B, C) }
func (c *CPU) op_LD_B_D() int { return c.loadByte(B, D) }
func (c *CPU) op_LD_B_E() int { return c.loadByte(B, E) }
func (c *CPU) op_LD_B_H() int { return c.loadByte(B, H) }
func (c *CPU) op_LD_B_L() int { return c.loadByte(B, L) }
func (c *CPU) op_LD_B_A() int { return c.loadByte(B, A) }
func (c *CPU) op_LD_C_B() int { return c.loadByte(C, D) }
func (c *CPU) op_LD_C_C() int { return c.loadByte(C, C) }
func (c *CPU) op_LD_C_D() int { return c.loadByte(C, D) }
func (c *CPU) op_LD_C_E() int { return c.loadByte(C, E) }
func (c *CPU) op_LD_C_H() int { return c.loadByte(C, H) }
func (c *CPU) op_LD_C_L() int { return c.loadByte(C, L) }
func (c *CPU) op_LD_C_A() int { return c.loadByte(C, A) }
func (c *CPU) op_LD_D_B() int { return c.loadByte(D, B) }
func (c *CPU) op_LD_D_C() int { return c.loadByte(D, C) }
func (c *CPU) op_LD_D_D() int { return c.loadByte(D, D) }
func (c *CPU) op_LD_D_E() int { return c.loadByte(D, E) }
func (c *CPU) op_LD_D_H() int { return c.loadByte(D, H) }
func (c *CPU) op_LD_D_L() int { return c.loadByte(D, L) }
func (c *CPU) op_LD_D_A() int { return c.loadByte(D, A) }
func (c *CPU) op_LD_E_B() int { return c.loadByte(E, B) }
func (c *CPU) op_LD_E_C() int { return c.loadByte(E, C) }
func (c *CPU) op_LD_E_D() int { return c.loadByte(E, D) }
func (c *CPU) op_LD_E_E() int { return c.loadByte(E, E) }
func (c *CPU) op_LD_E_H() int { return c.loadByte(E, H) }
func (c *CPU) op_LD_E_L() int { return c.loadByte(E, L) }
func (c *CPU) op_LD_E_A() int { return c.loadByte(E, A) }
func (c *CPU) op_LD_H_B() int { return c.loadByte(H, B) }
func (c *CPU) op_LD_H_C() int { return c.loadByte(H, C) }
func (c *CPU) op_LD_H_D() int { return c.loadByte(H, D) }
func (c *CPU) op_LD_H_E() int { return c.loadByte(H, E) }
func (c *CPU) op_LD_H_H() int { return c.loadByte(H, H) }
func (c *CPU) op_LD_H_L() int { return c.loadByte(H, L) }
func (c *CPU) op_LD_H_A() int { return c.loadByte(H, A) }
func (c *CPU) op_LD_L_B() int { return c.loadByte(L, B) }
func (c *CPU) op_LD_L_C() int { return c.loadByte(L, C) }
func (c *CPU) op_LD_L_D() int { return c.loadByte(L, D) }
func (c *CPU) op_LD_L_E() int { return c.loadByte(L, E) }
func (c *CPU) op_LD_L_H() int { return c.loadByte(L, H) }
func (c *CPU) op_LD_L_L() int { return c.loadByte(L, L) }
func (c *CPU) op_LD_L_A() int { return c.loadByte(L, A) }
func (c *CPU) op_LD_A_B() int { return c.loadByte(A, B) }
func (c *CPU) op_LD_A_C() int { return c.loadByte(A, C) }
func (c *CPU) op_LD_A_D() int { return c.loadByte(A, D) }
func (c *CPU) op_LD_A_E() int { return c.loadByte(A, E) }
func (c *CPU) op_LD_A_H() int { return c.loadByte(A, H) }
func (c *CPU) op_LD_A_L() int { return c.loadByte(A, L) }
func (c *CPU) op_LD_A_A() int { return c.loadByte(A, A) }

func (c *CPU) loadByte(dest register, src register) int {
	c.setRegByte(dest, c.getRegByte(src))
	return 1
}

func (c *CPU) op_LD_B_n() int { return c.setRegFromMemory(B) }
func (c *CPU) op_LD_D_n() int { return c.setRegFromMemory(D) }
func (c *CPU) op_LD_H_n() int { return c.setRegFromMemory(H) }
func (c *CPU) op_LD_C_n() int { return c.setRegFromMemory(C) }
func (c *CPU) op_LD_E_n() int { return c.setRegFromMemory(E) }
func (c *CPU) op_LD_L_n() int { return c.setRegFromMemory(L) }
func (c *CPU) op_LD_A_n() int { return c.setRegFromMemory(A) }

func (c *CPU) setRegFromMemory(reg register) int {
	c.setRegByte(reg, c.memory.ReadByte(c.regPC))
	c.regPC++
	return 1
}

func (c *CPU) op_LD_BC_nn() int {
	c.setRegShort(BC, c.memory.ReadShort(c.regPC))
	c.regPC += 2
	return 1
}
func (c *CPU) op_LD_DE_nn() int {
	c.setRegShort(DE, c.memory.ReadShort(c.regPC))
	c.regPC += 2
	return 1
}
func (c *CPU) op_LD_HL_nn() int {
	c.setRegShort(HL, c.memory.ReadShort(c.regPC))
	c.regPC += 2
	return 1
}
func (c *CPU) op_LD_SP_nn() int {
	c.setRegShort(SP, c.memory.ReadShort(c.regPC))
	c.regPC += 2
	return 1
}
func (c *CPU) op_LD_A_HL_plus() int {
	hl := c.getRegShort(HL)
	c.setRegByte(A, c.memory.ReadByte(hl))
	hl++
	c.setRegShort(HL, hl)
	return 1
}
func (c *CPU) op_LD_DE_A() int {
	c.memory.WriteByte(c.getRegShort(DE), c.getRegByte(A))
	return 2
}
func (c *CPU) op_LD_BC_A() int {
	c.memory.WriteByte(c.getRegShort(BC), c.getRegByte(A))
	return 2
}
func (c *CPU) op_LD_nn_A() int {
	c.memory.WriteByte(c.regPC, c.getRegByte(A))
	c.regPC += 2
	return 1
}
func (c *CPU) op_LDH_n_A() int {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	var addr uint16 = 0xFF00 | uint16(n)
	c.memory.WriteByte(addr, c.getRegByte(A))
	return 1
}

func (c *CPU) op_LD_HL_B() int { return c.LD_HL_x(B) }
func (c *CPU) op_LD_HL_C() int { return c.LD_HL_x(C) }
func (c *CPU) op_LD_HL_D() int { return c.LD_HL_x(D) }
func (c *CPU) op_LD_HL_E() int { return c.LD_HL_x(E) }
func (c *CPU) op_LD_HL_H() int { return c.LD_HL_x(H) }
func (c *CPU) op_LD_HL_L() int { return c.LD_HL_x(L) }
func (c *CPU) op_LD_HL_A() int { return c.LD_HL_x(A) }

func (c *CPU) LD_HL_x(reg register) int {
	c.memory.WriteByte(c.regHL, c.getRegByte(reg))
	return 1
}

func (c *CPU) op_LD_A_DE() int {
	c.setRegByte(A, c.memory.ReadByte(c.getRegShort(DE)))
	return 1
}

func (c *CPU) op_LD_HL_plus_A() int {
	c.memory.WriteByte(c.regHL, c.getRegByte(A))
	c.regHL++
	return 1
}

func (c *CPU) op_LD_HL_sub_A() int {
	c.memory.WriteByte(c.regHL, c.getRegByte(A))
	c.regHL--
	return 1
}

func (c *CPU) op_LDH_C_A() int {
	addr := 0xFF00 | uint16(c.getRegByte(C))
	c.memory.WriteByte(addr, c.getRegByte(A))
	return 1
}

func (c *CPU) op_LD_B_HL() int { return c.LD_x_HL(B) }
func (c *CPU) op_LD_D_HL() int { return c.LD_x_HL(D) }
func (c *CPU) op_LD_H_HL() int { return c.LD_x_HL(H) }
func (c *CPU) op_LD_C_HL() int { return c.LD_x_HL(C) }
func (c *CPU) op_LD_E_HL() int { return c.LD_x_HL(E) }
func (c *CPU) op_LD_L_HL() int { return c.LD_x_HL(L) }
func (c *CPU) op_LD_A_HL() int { return c.LD_x_HL(A) }

func (c *CPU) LD_x_HL(reg register) int {
	c.setRegByte(reg, c.memory.ReadByte(c.regHL))
	return 1
}

func (c *CPU) op_INC_C() int  { return c.incrementByteRegister(C) }
func (c *CPU) op_INC_E() int  { return c.incrementByteRegister(E) }
func (c *CPU) op_INC_L() int  { return c.incrementByteRegister(L) }
func (c *CPU) op_INC_A() int  { return c.incrementByteRegister(A) }
func (c *CPU) op_INC_B() int  { return c.incrementByteRegister(B) }
func (c *CPU) op_INC_D() int  { return c.incrementByteRegister(D) }
func (c *CPU) op_INC_H() int  { return c.incrementByteRegister(H) }
func (c *CPU) op_DEC_C() int  { return c.decrementByteRegister(C) }
func (c *CPU) op_DEC_E() int  { return c.decrementByteRegister(E) }
func (c *CPU) op_DEC_L() int  { return c.decrementByteRegister(L) }
func (c *CPU) op_DEC_A() int  { return c.decrementByteRegister(A) }
func (c *CPU) op_DEC_B() int  { return c.decrementByteRegister(B) }
func (c *CPU) op_DEC_D() int  { return c.decrementByteRegister(D) }
func (c *CPU) op_DEC_H() int  { return c.decrementByteRegister(H) }
func (c *CPU) op_INC_BC() int { return c.incrementShortRegister(BC) }
func (c *CPU) op_INC_DE() int { return c.incrementShortRegister(DE) }
func (c *CPU) op_INC_HL() int { return c.incrementShortRegister(HL) }
func (c *CPU) op_INC_SP() int { return c.incrementShortRegister(SP) }
func (c *CPU) op_DEC_BC() int { return c.decrementShortRegister(BC) }
func (c *CPU) op_DEC_DE() int { return c.decrementShortRegister(DE) }
func (c *CPU) op_DEC_HL() int { return c.decrementShortRegister(HL) }
func (c *CPU) op_DEC_SP() int { return c.decrementShortRegister(SP) }

func (c *CPU) incrementByteRegister(reg register) int {
	value := c.getRegByte(reg)
	current := value + 1
	c.setFlagZ(current == 0)
	c.setFlagN(false)

	abc := current ^ 0x01 ^ value
	c.setFlagH(abc&0x10 == 0x10) //((current - 1) & 0x10) == 0x10)
	//c.setFlagC(true) // TODO - why?
	c.setRegByte(reg, current)
	return 1
}

func (c *CPU) decrementByteRegister(reg register) int {
	value := c.getRegByte(reg)
	current := value - 1
	c.setFlagZ(current == 0)
	c.setFlagN(true)

	abc := current ^ 0x01 ^ value
	c.setFlagH(abc&0x10 == 0x10) //((current + 1) & 0x10) == 0x10)
	c.setRegByte(reg, current)
	return 1
}

func (c *CPU) incrementShortRegister(reg register) int {
	current := c.getRegShort(reg) + 1
	//	c.setFlagZ(current == 0)
	//	c.setFlagN(true) // TODO - why?
	//	c.setFlagH(current&0x10 == 0x10)
	//	c.setFlagC(true) // TODO - why?
	c.setRegShort(reg, current)
	return 1
}

func (c *CPU) decrementShortRegister(reg register) int {
	current := c.getRegShort(reg) + 1
	//	c.setFlagZ(current == 0)
	//	c.setFlagN(true)
	//	c.setFlagH(current&0x10 == 0x10)
	c.setRegShort(reg, current)
	return 1
}

func (c *CPU) op_ADD_n() int {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	c.op_ADD(n)
	return 1
}

func (c *CPU) op_ADD(n byte) int {
	a := c.getRegByte(A)
	result := a + n
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH((a & 0x10) == 0x10)
	c.setFlagC((a & 0x0F) == 0x0F)
	return 1
}

func (c *CPU) op_SUB_n() int {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	c.op_SUB(n)
	return 1
}

func (c *CPU) op_SUB(n byte) int {
	a := c.getRegByte(A)
	result := a - n
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(true)
	c.setFlagH((a & 0x10) == 0x10)
	c.setFlagC((a & 0x0F) == 0x0F)
	return 1
}

func (c *CPU) op_NOP() int { return 1 }

func (c *CPU) op_DI() int {
	// TODO - set IME?
	return 1
}

func (c *CPU) op_JP_nn() int {
	newAddr := c.memory.ReadShort(c.regPC)
	c.regPC = newAddr
	return 1
}

func (c *CPU) op_JR_e() int {
	c.conditionalJumpToOffset(true)
	return 1
}

func (c *CPU) op_JR_NZ_e() int {
	return c.conditionalJumpToOffset(!c.getFlagZ())
}

func (c *CPU) op_JR_NC_e() int {
	return c.conditionalJumpToOffset(!c.getFlagC())
}

func (c *CPU) op_JR_Z_e() int {
	return c.conditionalJumpToOffset(c.getFlagZ())
}

func (c *CPU) op_JR_C_e() int {
	return c.conditionalJumpToOffset(c.getFlagC())
}

func (c *CPU) conditionalJumpToOffset(condition bool) int {
	offset := c.memory.ReadByte((c.regPC))
	c.regPC++

	if !condition {
		return 1
	}

	//	fmt.Printf("JMP PC: 0x%04X, offset: %d, signed: %d", c.regPC, offset, int8(offset))

	// TODO - is this right?
	c.regPC = uint16(int16(c.regPC) + int16(int8(offset)))
	return 1
}

func (c *CPU) op_CALL_nn() int {
	nn := c.memory.ReadShort(c.regPC)
	c.regPC += 2
	c.regSP -= 2
	c.memory.WriteShort(c.regSP, c.regPC)
	c.regPC = nn
	return 1
}

func (c *CPU) op_CALL_NZ_nn() int { return c.callCondition(!c.getFlagZ()) }
func (c *CPU) op_CALL_Z_nn() int  { return c.callCondition(c.getFlagZ()) }
func (c *CPU) op_CALL_NC_nn() int { return c.callCondition(!c.getFlagC()) }
func (c *CPU) op_CALL_C_nn() int  { return c.callCondition(c.getFlagC()) }

func (c *CPU) callCondition(condition bool) int {
	nn := c.memory.ReadShort(c.regPC)
	c.regPC += 2
	if condition {
		c.regSP -= 2
		c.memory.WriteShort(c.regSP, c.regPC)
		c.setRegShort(PC, nn)
	}

	return 1
}

func (c *CPU) op_RET() int {

	c.setRegShort(PC, c.memory.ReadShort(c.getRegShort(SP)))
	c.regSP += 2
	return 1
}

func (c *CPU) op_PUSH_BC() int { return c.push(BC) }
func (c *CPU) op_PUSH_DE() int { return c.push(DE) }
func (c *CPU) op_PUSH_HL() int { return c.push(HL) }
func (c *CPU) op_PUSH_AF() int { return c.push(AF) }
func (c *CPU) op_POP_BC() int  { return c.pop(BC) }
func (c *CPU) op_POP_DE() int  { return c.pop(DE) }
func (c *CPU) op_POP_HL() int  { return c.pop(HL) }
func (c *CPU) op_POP_AF() int  { return c.pop(AF) }

func (c *CPU) push(reg register) int {
	c.regSP -= 2
	c.memory.WriteShort(c.regSP, c.getRegShort(reg))
	return 1
}

func (c *CPU) pop(reg register) int {
	c.setRegShort(reg, c.memory.ReadShort(c.regSP))
	c.regSP += 2
	return 1
}

func (c *CPU) op_OR_B() int { return c.or(B) }
func (c *CPU) op_OR_C() int { return c.or(C) }
func (c *CPU) op_OR_D() int { return c.or(D) }
func (c *CPU) op_OR_E() int { return c.or(E) }
func (c *CPU) op_OR_H() int { return c.or(H) }
func (c *CPU) op_OR_L() int { return c.or(L) }
func (c *CPU) op_OR_A() int { return c.or(A) }

func (c *CPU) or(reg register) int {
	result := c.getRegByte(A) | c.getRegByte(reg)
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
	return 1
}

func (c *CPU) op_AND_B() int { return c.and(B) }
func (c *CPU) op_AND_C() int { return c.and(C) }
func (c *CPU) op_AND_D() int { return c.and(D) }
func (c *CPU) op_AND_E() int { return c.and(E) }
func (c *CPU) op_AND_H() int { return c.and(H) }
func (c *CPU) op_AND_L() int { return c.and(L) }
func (c *CPU) op_AND_A() int { return c.and(A) }

func (c *CPU) and(reg register) int {
	result := c.getRegByte(A) & c.getRegByte(reg)
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
	return 1
}

func (c *CPU) op_XOR_B() int { return c.xor(B) }
func (c *CPU) op_XOR_C() int { return c.xor(C) }
func (c *CPU) op_XOR_D() int { return c.xor(D) }
func (c *CPU) op_XOR_E() int { return c.xor(E) }
func (c *CPU) op_XOR_H() int { return c.xor(H) }
func (c *CPU) op_XOR_L() int { return c.xor(L) }
func (c *CPU) op_XOR_A() int { return c.xor(A) }

func (c *CPU) xor(reg register) int {
	result := c.getRegByte(A) ^ c.getRegByte(reg)
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
	return 1
}

func (c *CPU) op_XOR_HL() int {
	data := c.memory.ReadByte(c.getRegShort(HL))
	result := c.getRegByte(A) ^ data
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
	return 1
}

func (c *CPU) op_CP_n() int {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	c.op_CP(n)
	return 1
}

func (c *CPU) op_CP(n byte) int {
	a := c.getRegByte(A)
	result := a - n
	c.setFlagZ(result == 0)
	c.setFlagN(true)
	c.setFlagH((a & 0x0F) < (n & 0x0F))
	c.setFlagC(a < n)
	return 1
}

func (c *CPU) op_LDH_A_n() int {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	c.setRegByte(A, c.memory.ReadByte(0xFF00|uint16(n)))
	return 1
}

func (c *CPU) op_LD_A_nn() int {
	nn := c.memory.ReadShort(c.regPC)
	c.regPC += 2
	c.setRegByte(A, c.memory.ReadByte(nn))
	return 1
}

func (c *CPU) op_AND_n() int {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	result := c.getRegByte(A) & n
	c.setRegByte(A, n)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(true)
	c.setFlagC(false)
	return 1
}

func (c *CPU) op_CB_op() int {
	op := c.memory.ReadByte(c.regPC)
	c.regPC++

	if op >= 0x40 && op <= 0x7F {
		return c.op_BIT(op)
	} else if op >= 0x30 && op <= 0x37 {
		return c.op_SWAP(op)
	} else if op >= 0x10 && op <= 0x17 {
		return c.op_RL(op)
	} else {
		panic(fmt.Sprintf("Unknown op for CB: 0x%02X", op))
	}
}

func (c *CPU) op_BIT(op byte) int {
	// TODO - no idea if this is right
	var reg register
	switch op {
	case 0x40:
	case 0x48:
	case 0x50:
	case 0x58:
	case 0x60:
	case 0x68:
	case 0x70:
	case 0x78:
		reg = B
	case 0x41:
	case 0x49:
	case 0x51:
	case 0x59:
	case 0x61:
	case 0x69:
	case 0x71:
	case 0x79:
		reg = C
	case 0x42:
	case 0x4A:
	case 0x52:
	case 0x5A:
	case 0x62:
	case 0x6A:
	case 0x72:
	case 0x7A:
		reg = D
	case 0x43:
	case 0x4B:
	case 0x53:
	case 0x5B:
	case 0x63:
	case 0x6B:
	case 0x73:
	case 0x7B:
		reg = E
	case 0x44:
	case 0x4C:
	case 0x54:
	case 0x5C:
	case 0x64:
	case 0x6C:
	case 0x74:
	case 0x7C:
		reg = H
	case 0x45:
	case 0x4D:
	case 0x55:
	case 0x5D:
	case 0x65:
	case 0x6D:
	case 0x75:
	case 0x7D:
		reg = L
	case 0x46:
	case 0x4E:
	case 0x56:
	case 0x5E:
	case 0x66:
	case 0x6E:
	case 0x76:
	case 0x7E:
		reg = HL
	case 0x47:
	case 0x4F:
	case 0x57:
	case 0x5F:
	case 0x67:
	case 0x6F:
	case 0x77:
	case 0x7F:
		reg = A
	default:
		panic(fmt.Sprintf("Unknown register op for BIT: 0x%02X", op))
		// TODO - implement loads of things
	}

	var bit int
	switch {
	case op >= 0x40 && op <= 0x47:
		bit = 0
	case op >= 0x48 && op <= 0x4F:
		bit = 1
	case op >= 0x50 && op <= 0x57:
		bit = 2
	case op >= 0x58 && op <= 0x5F:
		bit = 3
	case op >= 0x60 && op <= 0x67:
		bit = 4
	case op >= 0x68 && op <= 0x6F:
		bit = 5
	case op >= 0x70 && op <= 0x77:
		bit = 6
	case op >= 0x78 && op <= 0x7F:
		bit = 7
	default:
		panic(fmt.Sprintf("Unknown bit op for BIT: 0x%02X", op))
	}

	value := c.getRegBit(reg, bit)
	// TODO - set Z if false. Not clear if that means set do nothing if
	// value is true though
	c.setFlagZ(!value)
	c.setFlagN(false)
	c.setFlagH(true)
	return 1
}

func (c *CPU) op_SWAP(op byte) int {
	// TODO - no idea if this is right
	var reg register
	switch op {
	case 0x37:
		reg = A
	case 0x30:
		reg = B
	case 0x31:
		reg = C
	case 0x32:
		reg = D
	case 0x33:
		reg = E
	case 0x34:
		reg = H
	case 0x35:
		reg = L
	case 0x36:
		reg = HL
	default:
		panic(fmt.Sprintf("Unknown reg op for SWAP: 0x%02X", op))
	}

	n := c.getRegByte(reg)
	// TODO - swap upper and lower nibbles
	c.setRegByte(reg, n)

	c.setFlagZ(n == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
	return 1
}

func (c *CPU) op_RL(op byte) int {
	var reg register
	switch op {
	case 0x10:
		reg = B
	case 0x11:
		reg = C
	case 0x12:
		reg = D
	case 0x13:
		reg = E
	case 0x14:
		reg = H
	case 0x15:
		reg = L
	case 0x16:
		reg = HL
	case 0x17:
		reg = A
	default:
		panic(fmt.Sprintf("Unknown reg op for RL: 0x%02X", op))
	}

	//oldBit7 := c.getRegBit(reg, 7)

	c.rotateLeft(reg)
	//value := c.getRegByte(reg)

	//c.setFlagZ(value == 0)
	//c.setFlagN(false)
	//c.setFlagH(false)
	//c.setFlagC(oldBit7)
	return 1
}

func (c *CPU) op_RL_A() int {
	result := c.getRegByte(A)

	carry := result & 0x80
	//carry = carry >> 7
	result = result << 1
	//result = result & 0xFE
	//result = result | carry

	if c.getFlagC() {
		result = result ^ 0x01
	}

	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(carry == 0x80) //carry&0x01 == 0x01)
	return 1
}

func (c *CPU) rotateLeft(reg register) {
	result := c.getRegByte(reg)

	carry := result & 0x80
	//carry = carry >> 7
	result = result << 1
	//result = result & 0xFE
	//result = result | carry
	if c.getFlagC() {
		result = result ^ 0x01
	}

	c.setRegByte(reg, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(carry == 0x80) //carry&0x01 == 0x01)
}

func (c *CPU) op_SUB_B() int  { return c.op_SUB(c.getRegByte(B)) }
func (c *CPU) op_SUB_C() int  { return c.op_SUB(c.getRegByte(C)) }
func (c *CPU) op_SUB_D() int  { return c.op_SUB(c.getRegByte(D)) }
func (c *CPU) op_SUB_E() int  { return c.op_SUB(c.getRegByte(E)) }
func (c *CPU) op_SUB_H() int  { return c.op_SUB(c.getRegByte(H)) }
func (c *CPU) op_SUB_L() int  { return c.op_SUB(c.getRegByte(L)) }
func (c *CPU) op_SUB_HL() int { return c.op_SUB(c.memory.ReadByte(c.getRegShort(HL))) }
func (c *CPU) op_SUB_A() int  { return c.op_SUB(c.getRegByte(A)) }

func (c *CPU) op_CP_B() int  { return c.op_CP(c.getRegByte(B)) }
func (c *CPU) op_CP_C() int  { return c.op_CP(c.getRegByte(C)) }
func (c *CPU) op_CP_D() int  { return c.op_CP(c.getRegByte(D)) }
func (c *CPU) op_CP_E() int  { return c.op_CP(c.getRegByte(E)) }
func (c *CPU) op_CP_H() int  { return c.op_CP(c.getRegByte(H)) }
func (c *CPU) op_CP_L() int  { return c.op_CP(c.getRegByte(L)) }
func (c *CPU) op_CP_HL() int { return c.op_CP(c.memory.ReadByte(c.getRegShort(HL))) }
func (c *CPU) op_CP_A() int  { return c.op_CP(c.getRegByte(A)) }

func (c *CPU) op_ADD_B() int  { return c.op_ADD(c.getRegByte(B)) }
func (c *CPU) op_ADD_C() int  { return c.op_ADD(c.getRegByte(C)) }
func (c *CPU) op_ADD_D() int  { return c.op_ADD(c.getRegByte(D)) }
func (c *CPU) op_ADD_E() int  { return c.op_ADD(c.getRegByte(E)) }
func (c *CPU) op_ADD_H() int  { return c.op_ADD(c.getRegByte(H)) }
func (c *CPU) op_ADD_L() int  { return c.op_ADD(c.getRegByte(L)) }
func (c *CPU) op_ADD_HL() int { return c.op_ADD(c.memory.ReadByte(c.getRegShort(HL))) }
func (c *CPU) op_ADD_A() int  { return c.op_ADD(c.getRegByte(A)) }

func (c *CPU) op_LD_HL_n() int {
	pc := c.getRegShort(PC)
	n := c.memory.ReadByte(pc)
	c.setRegShort(PC, pc+1)
	c.memory.WriteByte(c.getRegShort(HL), n)
	return 1
}

func (c *CPU) op_CPL() int {
	value := c.getRegByte(A)
	// Flip all bits
	value = value ^ 0xFF
	c.setRegByte(A, value)
	c.setFlagN(true)
	c.setFlagH(true)
	return 1
}

func (c *CPU) op_CCF() int {
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(!c.getFlagC())
	return 1
}

func (c *CPU) op_RST_0x00() int { return c.RST_x(0x00) }
func (c *CPU) op_RST_0x08() int { return c.RST_x(0x08) }
func (c *CPU) op_RST_0x10() int { return c.RST_x(0x10) }
func (c *CPU) op_RST_0x18() int { return c.RST_x(0x18) }
func (c *CPU) op_RST_0x20() int { return c.RST_x(0x20) }
func (c *CPU) op_RST_0x28() int { return c.RST_x(0x28) }
func (c *CPU) op_RST_0x30() int { return c.RST_x(0x30) }
func (c *CPU) op_RST_0x38() int { return c.RST_x(0x38) }

func (c *CPU) RST_x(n byte) int {
	c.push(PC)
	c.setRegShort(PC, uint16(n))
	return 4
}

func (c *CPU) op_LD_nn_SP() int {
	pc := c.getRegShort(PC)
	nn := c.memory.ReadShort(pc)
	c.setRegShort(PC, pc+2)
	c.memory.WriteShort(nn, c.getRegShort(SP))
	return 5
}
