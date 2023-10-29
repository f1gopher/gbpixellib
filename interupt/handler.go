package interupt

import (
	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/memory"
)

type Interupt int

const (
	VBlank Interupt = iota
	LCD
	Time
	Serial
	Joypad
)

const InteruptEnableRegister = 0xFFFF
const InteruptFlag = 0xFF0F

func (i Interupt) String() string {
	return [...]string{"V-Blank", "LCD", "Timer", "Joypad"}[i]
}

type Handler struct {
	memory *memory.Memory
	regs   *cpu.Registers
}

func CreateHandler(memory *memory.Memory, registers *cpu.Registers) *Handler {
	return &Handler{
		memory: memory,
		regs:   registers,
	}
}

func (h *Handler) Reset() {
}

func (h *Handler) Request(i Interupt) {
	value := h.memory.ReadByte(InteruptFlag)
	// Todo - why??
	value = value | 0xE0

	var bit uint8 = 0
	switch i {
	case VBlank:
		bit = 0
	case LCD:
		bit = 1
	case Time:
		bit = 2
	case Serial:
		bit = 3
	case Joypad:
		bit = 4
	default:
		panic("Unhandled interupt type")
	}

	value = memory.SetBit(value, bit, true)
	h.memory.WriteByte(InteruptFlag, value)
}

func (h *Handler) Update() bool {
	if h.regs.GetIME() {
		req := h.memory.ReadByte(InteruptFlag)
		enabled := h.memory.ReadByte(InteruptEnableRegister)

		if req > 0 {
			for i := 0; i < 5; i++ {
				if memory.GetBit(req, i) {
					if memory.GetBit(enabled, i) {
						h.serviceInterupt(uint8(i))
						return true
					}
				}
			}
		}
	}

	return false
}

func (h *Handler) serviceInterupt(interupt uint8) {
	h.regs.SetIME(false)
	req := h.memory.ReadByte(InteruptFlag)
	req = memory.SetBit(req, interupt, false)
	h.memory.WriteByte(InteruptFlag, req)

	// TODO - push PC onto stack

	var programCounter uint16 = 0
	switch interupt {
	case 0: // Vertical Blank
		programCounter = 0x0040
	case 1: // LCDC Status
		programCounter = 0x0048
	case 2: // Timer Overflow
		programCounter = 0x0050
	case 3: // Serial Transfer
		programCounter = 0x0058
	case 4: // Joypad
		programCounter = 0x0060
	default:
		panic("Unhandled service interupt")
	}

	// TODO - set PC
	//h.cpu.PushAndReplacePC(programCounter)

	// TODO - -1 for PC because we fetch the next opcode but we shouldn't???

	currentPC := h.regs.Get16(cpu.PC) - 1
	cpu.DecAndWriteSP(h.regs, h.memory, cpu.Msb(currentPC))
	cpu.DecAndWriteSP(h.regs, h.memory, cpu.Lsb(currentPC))
	h.regs.SetPC(programCounter)
}
