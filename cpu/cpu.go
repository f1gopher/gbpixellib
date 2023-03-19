package cpu

import (
	"errors"
	"fmt"
	"go-boy/memory"
	"os"
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

	opcodes map[byte]func()

	debugLogFile *os.File

	count int
}

func CreateCPU(memory *memory.Memory) *CPU {
	chip := CPU{
		memory: memory,
		count:  0,
	}

	// TODO - opcodes range from 0-256 so we can use an array instead
	chip.opcodes = map[byte]func(){
		0x40: chip.op_LD_B_B,
		0x41: chip.op_LD_B_C,
		0x42: chip.op_LD_B_D,
		0x43: chip.op_LD_B_E,
		0x44: chip.op_LD_B_H,
		0x45: chip.op_LD_B_L,
		0x47: chip.op_LD_B_A,
		0x48: chip.op_LD_C_B,
		0x49: chip.op_LD_C_C,
		0x4A: chip.op_LD_C_D,
		0x4B: chip.op_LD_C_E,
		0x4C: chip.op_LD_C_H,
		0x4D: chip.op_LD_C_L,
		0x4F: chip.op_LD_C_A,
		0x50: chip.op_LD_D_B,
		0x51: chip.op_LD_D_C,
		0x52: chip.op_LD_D_D,
		0x53: chip.op_LD_D_E,
		0x54: chip.op_LD_D_H,
		0x55: chip.op_LD_D_L,
		0x57: chip.op_LD_D_A,
		0x58: chip.op_LD_E_B,
		0x59: chip.op_LD_E_C,
		0x5A: chip.op_LD_E_D,
		0x5B: chip.op_LD_E_E,
		0x5C: chip.op_LD_E_H,
		0x5D: chip.op_LD_E_L,
		0x5F: chip.op_LD_E_A,
		0x60: chip.op_LD_H_B,
		0x61: chip.op_LD_H_C,
		0x62: chip.op_LD_H_D,
		0x63: chip.op_LD_H_E,
		0x64: chip.op_LD_H_H,
		0x65: chip.op_LD_H_L,
		0x67: chip.op_LD_H_A,
		0x68: chip.op_LD_L_B,
		0x69: chip.op_LD_L_C,
		0x6A: chip.op_LD_L_D,
		0x6B: chip.op_LD_L_E,
		0x6C: chip.op_LD_L_H,
		0x6D: chip.op_LD_L_L,
		0x6F: chip.op_LD_L_A,
		0x78: chip.op_LD_A_B,
		0x79: chip.op_LD_A_C,
		0x7A: chip.op_LD_A_D,
		0x7B: chip.op_LD_A_E,
		0x7C: chip.op_LD_A_H,
		0x7D: chip.op_LD_A_L,
		0x7F: chip.op_LD_A_A,
		0x00: chip.op_NOP,
		0xC3: chip.op_JP,
		0x01: chip.op_LD_BC_nn,
		0x11: chip.op_LD_DE_nn,
		0x21: chip.op_LD_HL_nn,
		0x31: chip.op_LD_SP_nn,
		0x2A: chip.op_LD_A_HL_plus,
		0x12: chip.op_LD_DE_A,
		0x0C: chip.op_INC_C,
		0x1C: chip.op_INC_E,
		0x2C: chip.op_INC_L,
		0x3C: chip.op_INC_A,
		0x04: chip.op_INC_B,
		0x14: chip.op_INC_D,
		0x24: chip.op_INC_H,
		0x0D: chip.op_DEC_C,
		0x1D: chip.op_DEC_E,
		0x2D: chip.op_DEC_L,
		0x3D: chip.op_DEC_A,
		0x05: chip.op_DEC_B,
		0x15: chip.op_DEC_D,
		0x25: chip.op_DEC_H,
		0x03: chip.op_INC_BC,
		0x13: chip.op_INC_DE,
		0x23: chip.op_INC_HL,
		0x33: chip.op_INC_HL,
		0x0B: chip.op_DEC_BC,
		0x1B: chip.op_DEC_DE,
		0x2B: chip.op_DEC_HL,
		0x3B: chip.op_DEC_SP,
		0x20: chip.op_JR_NZ_e,
		0x18: chip.op_JR_e,
		0x30: chip.op_JR_NC_e,
		0x28: chip.op_JR_Z_e,
		0x38: chip.op_JR_C_e,
		0xF3: chip.op_DI,
		0xEA: chip.op_LD_nn_A,
		0x06: chip.op_LD_B_n,
		0x16: chip.op_LD_D_n,
		0x26: chip.op_LD_H_n,
		0x0E: chip.op_LD_C_n,
		0x1E: chip.op_LD_E_n,
		0x2E: chip.op_LD_L_n,
		0x3E: chip.op_LD_A_n,
		0xE0: chip.op_LDH_n_A,
		0xCD: chip.op_CALL_nn,
		0xC9: chip.op_RET,
		0xC5: chip.op_PUSH_BC,
		0xD5: chip.op_PUSH_DE,
		0xE5: chip.op_PUSH_HL,
		0xF5: chip.op_PUSH_AF,
		0xC1: chip.op_POP_BC,
		0xD1: chip.op_POP_DE,
		0xE1: chip.op_POP_HL,
		0xF1: chip.op_POP_AF,
		0xB0: chip.op_OR_B,
		0xB1: chip.op_OR_C,
		0xB2: chip.op_OR_D,
		0xB3: chip.op_OR_E,
		0xB4: chip.op_OR_H,
		0xB5: chip.op_OR_L,
		0xB7: chip.op_OR_A,
		0xA0: chip.op_AND_B,
		0xA1: chip.op_AND_C,
		0xA2: chip.op_AND_D,
		0xA3: chip.op_AND_E,
		0xA4: chip.op_AND_H,
		0xA5: chip.op_AND_L,
		0xA7: chip.op_AND_A,
		0xA8: chip.op_XOR_B,
		0xA9: chip.op_XOR_C,
		0xAA: chip.op_XOR_D,
		0xAB: chip.op_XOR_E,
		0xAC: chip.op_XOR_H,
		0xAD: chip.op_XOR_L,
		0xAF: chip.op_XOR_A,
		0xF0: chip.op_LDH_A_n,
		0xFE: chip.op_CP_n,
		0xFA: chip.op_LD_A_nn,
		0xE6: chip.op_AND_n,
		0xC4: chip.op_CALL_NZ_nn,
		0xCC: chip.op_CALL_Z_nn,
		0xD4: chip.op_CALL_NC_nn,
		0xDC: chip.op_CALL_C_nn,
		0x70: chip.op_LD_HL_B,
		0x71: chip.op_LD_HL_C,
		0x72: chip.op_LD_HL_D,
		0x73: chip.op_LD_HL_E,
		0x74: chip.op_LD_HL_H,
		0x75: chip.op_LD_HL_L,
		0x77: chip.op_LD_HL_A,
		0x1A: chip.op_LD_A_DE,
		0x22: chip.op_LD_HL_plus_A,
		0x32: chip.op_LD_HL_sub_A,
		0xC6: chip.op_ADD_n,
		0xD6: chip.op_SUB_n,
		0x46: chip.op_LD_B_HL,
		0x56: chip.op_LD_D_HL,
		0x66: chip.op_LD_H_HL,
		0x4E: chip.op_LD_C_HL,
		0x5E: chip.op_LD_E_HL,
		0x6E: chip.op_LD_L_HL,
		0x7E: chip.op_LD_A_HL,
		0xAE: chip.op_XOR_HL,
		0xCB: chip.op_CB_op,
		0xE2: chip.op_LDH_C_A,
		0x17: chip.op_RL_A,
	}

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

	c.debugLog()
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

func (c *CPU) SetDebugLog(file string) error {
	var err error
	c.debugLogFile, err = os.Create(file)

	if err != nil {
		return errors.Join(errors.New("Unable to set debug log"), err)
	}

	return nil
}

func (c *CPU) Tick() {
	opcode := c.memory.ReadByte(c.regPC)

	//fmt.Printf("%d - Op: 0x%02X ", c.count, opcode)

	executor, exists := c.opcodes[opcode]
	if !exists {
		panic(fmt.Sprintf("Unsupported opcode: 0x%02X", opcode))
	}

	// TODO - is this the best place to increment the program counter?
	c.regPC++
	executor()

	c.debugLog()

	c.count++
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
	if reg != F {
		panic(fmt.Sprintf("Unexpected register for get bit: %s", reg.String()))
	}

	if bit > 7 || bit < 0 {
		panic(fmt.Sprintf("Invalid bit for getRegBit: %d", bit))
	}

	if c.getRegByte(reg)>>bit == 1 {
		return true
	}

	return false
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

	if bit > 7 || bit < 0 {
		panic(fmt.Sprintf("Invalid bit for setRegBit: %d", bit))
	}

	current := c.getRegByte(reg)

	if value {
		current = current | 0x01<<bit
	} else {
		current = current &^ (0x01 << bit)
	}

	c.setRegByte(reg, current)
}

func (c *CPU) debugLog() {
	msg := c.Debug()
	//fmt.Print(msg)
	if c.debugLogFile != nil {
		c.debugLogFile.WriteString(msg)
	}
}

func (c *CPU) Debug() string {

	// TODO - handle the PC address being 0xFFFF so trying to read would go past the end
	p1 := c.memory.ReadByte(c.regPC)
	p2 := c.memory.ReadByte(c.regPC + 1)
	p3 := c.memory.ReadByte(c.regPC + 2)
	p4 := c.memory.ReadByte(c.regPC + 3)

	return fmt.Sprintf(
		"%d=>Next_Opcode:0x%02X A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X\n",
		c.count,
		p1,
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
}

func (c *CPU) op_LD_B_B() { c.loadByte(B, B) }
func (c *CPU) op_LD_B_C() { c.loadByte(B, C) }
func (c *CPU) op_LD_B_D() { c.loadByte(B, D) }
func (c *CPU) op_LD_B_E() { c.loadByte(B, E) }
func (c *CPU) op_LD_B_H() { c.loadByte(B, H) }
func (c *CPU) op_LD_B_L() { c.loadByte(B, L) }
func (c *CPU) op_LD_B_A() { c.loadByte(B, A) }
func (c *CPU) op_LD_C_B() { c.loadByte(C, D) }
func (c *CPU) op_LD_C_C() { c.loadByte(C, C) }
func (c *CPU) op_LD_C_D() { c.loadByte(C, D) }
func (c *CPU) op_LD_C_E() { c.loadByte(C, E) }
func (c *CPU) op_LD_C_H() { c.loadByte(C, H) }
func (c *CPU) op_LD_C_L() { c.loadByte(C, L) }
func (c *CPU) op_LD_C_A() { c.loadByte(C, A) }
func (c *CPU) op_LD_D_B() { c.loadByte(D, B) }
func (c *CPU) op_LD_D_C() { c.loadByte(D, C) }
func (c *CPU) op_LD_D_D() { c.loadByte(D, D) }
func (c *CPU) op_LD_D_E() { c.loadByte(D, E) }
func (c *CPU) op_LD_D_H() { c.loadByte(D, H) }
func (c *CPU) op_LD_D_L() { c.loadByte(D, L) }
func (c *CPU) op_LD_D_A() { c.loadByte(D, A) }
func (c *CPU) op_LD_E_B() { c.loadByte(E, B) }
func (c *CPU) op_LD_E_C() { c.loadByte(E, C) }
func (c *CPU) op_LD_E_D() { c.loadByte(E, D) }
func (c *CPU) op_LD_E_E() { c.loadByte(E, E) }
func (c *CPU) op_LD_E_H() { c.loadByte(E, H) }
func (c *CPU) op_LD_E_L() { c.loadByte(E, L) }
func (c *CPU) op_LD_E_A() { c.loadByte(E, A) }
func (c *CPU) op_LD_H_B() { c.loadByte(H, B) }
func (c *CPU) op_LD_H_C() { c.loadByte(H, C) }
func (c *CPU) op_LD_H_D() { c.loadByte(H, D) }
func (c *CPU) op_LD_H_E() { c.loadByte(H, E) }
func (c *CPU) op_LD_H_H() { c.loadByte(H, H) }
func (c *CPU) op_LD_H_L() { c.loadByte(H, L) }
func (c *CPU) op_LD_H_A() { c.loadByte(H, A) }
func (c *CPU) op_LD_L_B() { c.loadByte(L, B) }
func (c *CPU) op_LD_L_C() { c.loadByte(L, C) }
func (c *CPU) op_LD_L_D() { c.loadByte(L, D) }
func (c *CPU) op_LD_L_E() { c.loadByte(L, E) }
func (c *CPU) op_LD_L_H() { c.loadByte(L, H) }
func (c *CPU) op_LD_L_L() { c.loadByte(L, L) }
func (c *CPU) op_LD_L_A() { c.loadByte(L, A) }
func (c *CPU) op_LD_A_B() { c.loadByte(A, B) }
func (c *CPU) op_LD_A_C() { c.loadByte(A, C) }
func (c *CPU) op_LD_A_D() { c.loadByte(A, D) }
func (c *CPU) op_LD_A_E() { c.loadByte(A, E) }
func (c *CPU) op_LD_A_H() { c.loadByte(A, H) }
func (c *CPU) op_LD_A_L() { c.loadByte(A, L) }
func (c *CPU) op_LD_A_A() { c.loadByte(A, A) }

func (c *CPU) loadByte(dest register, src register) {
	c.setRegByte(dest, c.getRegByte(src))
}

func (c *CPU) op_LD_B_n() { c.setRegFromMemory(B) }
func (c *CPU) op_LD_D_n() { c.setRegFromMemory(D) }
func (c *CPU) op_LD_H_n() { c.setRegFromMemory(H) }
func (c *CPU) op_LD_C_n() { c.setRegFromMemory(C) }
func (c *CPU) op_LD_E_n() { c.setRegFromMemory(E) }
func (c *CPU) op_LD_L_n() { c.setRegFromMemory(L) }
func (c *CPU) op_LD_A_n() { c.setRegFromMemory(A) }

func (c *CPU) setRegFromMemory(reg register) {
	c.setRegByte(reg, c.memory.ReadByte(c.regPC))
	c.regPC++
}

func (c *CPU) op_LD_BC_nn() {
	c.setRegShort(BC, c.memory.ReadShort(c.regPC))
	c.regPC += 2
}
func (c *CPU) op_LD_DE_nn() {
	c.setRegShort(DE, c.memory.ReadShort(c.regPC))
	c.regPC += 2
}
func (c *CPU) op_LD_HL_nn() {
	c.setRegShort(HL, c.memory.ReadShort(c.regPC))
	c.regPC += 2
}
func (c *CPU) op_LD_SP_nn() {
	c.setRegShort(SP, c.memory.ReadShort(c.regPC))
	c.regPC += 2
}
func (c *CPU) op_LD_A_HL_plus() {
	hl := c.getRegShort(HL)
	c.setRegByte(A, c.memory.ReadByte(hl))
	hl++
	c.setRegShort(HL, hl)
}
func (c *CPU) op_LD_DE_A() {
	c.memory.WriteByte(c.getRegShort(DE), c.getRegByte(A))
}
func (c *CPU) op_LD_nn_A() {
	c.memory.WriteByte(c.regPC, c.getRegByte(A))
	c.regPC += 2
}
func (c *CPU) op_LDH_n_A() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	var addr uint16 = 0xFF00 | uint16(n)
	c.memory.WriteByte(addr, c.getRegByte(A))
}

func (c *CPU) op_LD_HL_B() { c.LD_HL_x(B) }
func (c *CPU) op_LD_HL_C() { c.LD_HL_x(C) }
func (c *CPU) op_LD_HL_D() { c.LD_HL_x(D) }
func (c *CPU) op_LD_HL_E() { c.LD_HL_x(E) }
func (c *CPU) op_LD_HL_H() { c.LD_HL_x(H) }
func (c *CPU) op_LD_HL_L() { c.LD_HL_x(L) }
func (c *CPU) op_LD_HL_A() { c.LD_HL_x(A) }

func (c *CPU) LD_HL_x(reg register) {
	c.memory.WriteByte(c.regHL, c.getRegByte(reg))
}

func (c *CPU) op_LD_A_DE() {
	c.setRegByte(A, c.memory.ReadByte(c.getRegShort(DE)))
}

func (c *CPU) op_LD_HL_plus_A() {
	c.memory.WriteByte(c.regHL, c.getRegByte(A))
	c.regHL++
}

func (c *CPU) op_LD_HL_sub_A() {
	c.memory.WriteByte(c.regHL, c.getRegByte(A))
	c.regHL--
}

func (c *CPU) op_LDH_C_A() {
	addr := 0xFF00 | uint16(c.getRegByte(C))
	c.memory.WriteByte(addr, c.getRegByte(A))
}

func (c *CPU) op_LD_B_HL() { c.LD_x_HL(B) }
func (c *CPU) op_LD_D_HL() { c.LD_x_HL(D) }
func (c *CPU) op_LD_H_HL() { c.LD_x_HL(H) }
func (c *CPU) op_LD_C_HL() { c.LD_x_HL(C) }
func (c *CPU) op_LD_E_HL() { c.LD_x_HL(E) }
func (c *CPU) op_LD_L_HL() { c.LD_x_HL(L) }
func (c *CPU) op_LD_A_HL() { c.LD_x_HL(A) }

func (c *CPU) LD_x_HL(reg register) {
	c.setRegByte(reg, c.memory.ReadByte(c.regHL))
}

func (c *CPU) op_INC_C()  { c.incrementByteRegister(C) }
func (c *CPU) op_INC_E()  { c.incrementByteRegister(E) }
func (c *CPU) op_INC_L()  { c.incrementByteRegister(L) }
func (c *CPU) op_INC_A()  { c.incrementByteRegister(A) }
func (c *CPU) op_INC_B()  { c.incrementByteRegister(B) }
func (c *CPU) op_INC_D()  { c.incrementByteRegister(D) }
func (c *CPU) op_INC_H()  { c.incrementByteRegister(H) }
func (c *CPU) op_DEC_C()  { c.decrementByteRegister(C) }
func (c *CPU) op_DEC_E()  { c.decrementByteRegister(E) }
func (c *CPU) op_DEC_L()  { c.decrementByteRegister(L) }
func (c *CPU) op_DEC_A()  { c.decrementByteRegister(A) }
func (c *CPU) op_DEC_B()  { c.decrementByteRegister(B) }
func (c *CPU) op_DEC_D()  { c.decrementByteRegister(D) }
func (c *CPU) op_DEC_H()  { c.decrementByteRegister(H) }
func (c *CPU) op_INC_BC() { c.incrementShortRegister(BC) }
func (c *CPU) op_INC_DE() { c.incrementShortRegister(DE) }
func (c *CPU) op_INC_HL() { c.incrementShortRegister(HL) }
func (c *CPU) op_INC_SP() { c.incrementShortRegister(SP) }
func (c *CPU) op_DEC_BC() { c.decrementShortRegister(BC) }
func (c *CPU) op_DEC_DE() { c.decrementShortRegister(DE) }
func (c *CPU) op_DEC_HL() { c.decrementShortRegister(HL) }
func (c *CPU) op_DEC_SP() { c.decrementShortRegister(SP) }

func (c *CPU) incrementByteRegister(reg register) {
	current := c.getRegByte(reg) + 1
	c.setFlagZ(current == 0)
	c.setFlagN(false)
	c.setFlagH(((current - 1) & 0x0F) == 0x0F)
	c.setFlagC(true) // TODO - why?
	c.setRegByte(reg, current)
}

func (c *CPU) decrementByteRegister(reg register) {
	current := c.getRegByte(reg) - 1
	c.setFlagZ(current == 0)
	c.setFlagN(true)
	c.setFlagH(((current + 1) & 0x10) == 0x10)
	c.setRegByte(reg, current)
}

func (c *CPU) incrementShortRegister(reg register) {
	current := c.getRegShort(reg) + 1
	//	c.setFlagZ(current == 0)
	//	c.setFlagN(true) // TODO - why?
	//	c.setFlagH(current&0x10 == 0x10)
	//	c.setFlagC(true) // TODO - why?
	c.setRegShort(reg, current)
}

func (c *CPU) decrementShortRegister(reg register) {
	current := c.getRegShort(reg) + 1
	//	c.setFlagZ(current == 0)
	//	c.setFlagN(true)
	//	c.setFlagH(current&0x10 == 0x10)
	c.setRegShort(reg, current)
}

func (c *CPU) op_ADD_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	a := c.getRegByte(A)
	result := a + n
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH((a & 0x10) == 0x10)
	c.setFlagC((a & 0x0F) == 0x0F)
}

func (c *CPU) op_SUB_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	a := c.getRegByte(A)
	result := a - n
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(true)
	c.setFlagH((a & 0x10) == 0x10)
	c.setFlagC((a & 0x0F) == 0x0F)
}

func (c *CPU) op_NOP() {}

func (c *CPU) op_DI() {
	// TODO - set IME?
}

func (c *CPU) op_JP() {
	newAddr := c.memory.ReadShort(c.regPC)
	c.regPC = newAddr
}

func (c *CPU) op_JR_e() {
	c.conditionalJumpToOffset(true)
}

func (c *CPU) op_JR_NZ_e() {
	c.conditionalJumpToOffset(!c.getFlagZ())
}

func (c *CPU) op_JR_NC_e() {
	c.conditionalJumpToOffset(!c.getFlagC())
}

func (c *CPU) op_JR_Z_e() {
	c.conditionalJumpToOffset(c.getFlagZ())
}

func (c *CPU) op_JR_C_e() {
	c.conditionalJumpToOffset(c.getFlagC())
}

func (c *CPU) conditionalJumpToOffset(condition bool) {
	offset := c.memory.ReadByte((c.regPC))
	c.regPC++

	if !condition {
		return
	}

	//	fmt.Printf("JMP PC: 0x%04X, offset: %d, signed: %d", c.regPC, offset, int8(offset))

	// TODO - is this right?
	c.regPC = uint16(int16(c.regPC) + int16(int8(offset)))
}

func (c *CPU) op_CALL_nn() {
	nn := c.memory.ReadShort(c.regPC)
	c.regPC += 2
	c.regSP -= 2
	c.memory.WriteShort(c.regSP, c.regPC)
	c.regPC = nn
}

func (c *CPU) op_CALL_NZ_nn() { c.callCondition(!c.getFlagZ()) }
func (c *CPU) op_CALL_Z_nn()  { c.callCondition(c.getFlagZ()) }
func (c *CPU) op_CALL_NC_nn() { c.callCondition(!c.getFlagC()) }
func (c *CPU) op_CALL_C_nn()  { c.callCondition(c.getFlagC()) }

func (c *CPU) callCondition(condition bool) {
	nn := c.memory.ReadShort(c.regPC)
	c.regPC += 2
	if !c.getFlagZ() {
		c.regSP -= 2
		c.memory.WriteShort(c.regSP, c.regPC)
		c.setRegShort(PC, nn)
	}
}

func (c *CPU) op_RET() {

	c.setRegShort(PC, c.memory.ReadShort(c.getRegShort(SP)))
	c.regSP += 2
}

func (c *CPU) op_PUSH_BC() { c.push(BC) }
func (c *CPU) op_PUSH_DE() { c.push(DE) }
func (c *CPU) op_PUSH_HL() { c.push(HL) }
func (c *CPU) op_PUSH_AF() { c.push(AF) }
func (c *CPU) op_POP_BC()  { c.pop(BC) }
func (c *CPU) op_POP_DE()  { c.pop(DE) }
func (c *CPU) op_POP_HL()  { c.pop(HL) }
func (c *CPU) op_POP_AF()  { c.pop(AF) }

func (c *CPU) push(reg register) {
	c.regSP -= 2
	c.memory.WriteShort(c.regSP, c.getRegShort(reg))
}

func (c *CPU) pop(reg register) {
	c.setRegShort(reg, c.memory.ReadShort(c.regSP))
	c.regSP += 2
}

func (c *CPU) op_OR_B() { c.or(B) }
func (c *CPU) op_OR_C() { c.or(C) }
func (c *CPU) op_OR_D() { c.or(D) }
func (c *CPU) op_OR_E() { c.or(E) }
func (c *CPU) op_OR_H() { c.or(H) }
func (c *CPU) op_OR_L() { c.or(L) }
func (c *CPU) op_OR_A() { c.or(A) }

func (c *CPU) or(reg register) {
	result := c.getRegByte(A) | c.getRegByte(reg)
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_AND_B() { c.and(B) }
func (c *CPU) op_AND_C() { c.and(C) }
func (c *CPU) op_AND_D() { c.and(D) }
func (c *CPU) op_AND_E() { c.and(E) }
func (c *CPU) op_AND_H() { c.and(H) }
func (c *CPU) op_AND_L() { c.and(L) }
func (c *CPU) op_AND_A() { c.and(A) }

func (c *CPU) and(reg register) {
	result := c.getRegByte(A) & c.getRegByte(reg)
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_XOR_B() { c.xor(B) }
func (c *CPU) op_XOR_C() { c.xor(C) }
func (c *CPU) op_XOR_D() { c.xor(D) }
func (c *CPU) op_XOR_E() { c.xor(E) }
func (c *CPU) op_XOR_H() { c.xor(H) }
func (c *CPU) op_XOR_L() { c.xor(L) }
func (c *CPU) op_XOR_A() { c.xor(A) }

func (c *CPU) xor(reg register) {
	result := c.getRegByte(A) ^ c.getRegByte(reg)
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_XOR_HL() {
	data := c.memory.ReadByte(c.getRegShort(HL))
	result := c.getRegByte(A) ^ data
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_CP_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	a := c.getRegByte(A)
	result := a - n
	c.setFlagZ(result == 0)
	c.setFlagN(true)
	c.setFlagH(((result - 1) & 0x0F) != 0x0F)
	c.setFlagC(a < n)
}

func (c *CPU) op_LDH_A_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	c.setRegByte(A, c.memory.ReadByte((0xFF00 | uint16(n))))
}

func (c *CPU) op_LD_A_nn() {
	nn := c.memory.ReadShort(c.regPC)
	c.regPC += 2
	c.setRegByte(A, c.memory.ReadByte(nn))
}

func (c *CPU) op_AND_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	result := c.getRegByte(A) & n
	c.setRegByte(A, n)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(true)
	c.setFlagC(false)
}

func (c *CPU) op_CB_op() {
	// TODO - no idea if this is right
	op := c.memory.ReadByte(c.regPC)
	c.regPC++

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
		//	panic(fmt.Sprintf("Unknown op for CB: 0x%02X", op))
		// TODO - implement loads of things
		return
	}

	n := c.getRegByte(reg)
	// TODO - swap upper and lower nibbles
	c.setRegByte(reg, n)

	c.setFlagZ(n == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_RL_A() { c.rotateLeft(A) }

func (c *CPU) rotateLeft(reg register) {
	result := c.getRegByte(reg)

	carry := result & 0x80
	carry = carry >> 7
	result = result << 1
	result = result & 0xFE
	result = result | carry

	c.setRegByte(reg, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(carry&0x01 == 0x01)
}
