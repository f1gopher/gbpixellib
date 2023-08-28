package interupt

import (
	"github.com/f1gopher/gbpixellib/memory"
)

type Interupt int

const (
	VBlank Interupt = iota
	LCD
	Time
	Joypad
)

func (i Interupt) String() string {
	return [...]string{"V-Blank", "LCD", "Timer", "Joypad"}[i]
}

type cpuHandler interface {
	PushAndReplacePC(newPC uint16)
}

type Handler struct {
	memory *memory.Memory
	cpu    cpuHandler

	interuptMaster bool
}

func CreateHandler(memory *memory.Memory, cpu cpuHandler) *Handler {
	return &Handler{
		memory:         memory,
		cpu:            cpu,
		interuptMaster: false,
	}
}

func (h *Handler) Reset() {
	h.interuptMaster = false
}

func (h *Handler) Request(i Interupt) {
	value := h.memory.ReadByte(0xFF0F)

	bit := 0
	switch i {
	case VBlank:
		bit = 0
	case LCD:
		bit = 1
	case Time:
		bit = 2
	case Joypad:
		bit = 4
	default:
		panic("Unhandled interupt type")
	}

	value = memory.SetBit(value, bit, true)
	h.memory.WriteByte(0xFF0F, value)
}

func (h *Handler) Update() {
	if h.interuptMaster {
		req := h.memory.ReadByte(0xFF0F)
		enabled := h.memory.ReadByte(0xFFFF)

		if req > 0 {
			for i := 0; i < 5; i++ {
				if memory.GetBit(req, i) {
					if memory.GetBit(enabled, i) {
						h.serviceInterupt(i)
					}
				}
			}
		}
	}
}

func (h *Handler) serviceInterupt(interupt int) {
	h.interuptMaster = false
	req := h.memory.ReadByte(0xFF0F)
	req = memory.SetBit(req, interupt, false)
	h.memory.WriteByte(0xFF0F, req)

	// TODO - push PC onto stack

	var programCounter uint16 = 0
	switch interupt {
	case 0:
		programCounter = 0x0040
	case 1:
		programCounter = 0x0048
	case 2:
		programCounter = 0x0050
	case 4:
		programCounter = 0x0060
	default:
		panic("Unhandled service interupt")
	}

	// TODO - set PC
	h.cpu.PushAndReplacePC(programCounter)
}
