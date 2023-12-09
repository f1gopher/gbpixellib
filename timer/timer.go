package timer

import (
	"github.com/f1gopher/gbpixellib/interupt"
	"github.com/f1gopher/gbpixellib/memory"
)

type memoryInterface interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)
	WriteDividerRegister(value uint8)
}

type interruptInterface interface {
	Request(i interupt.Interupt)
}

type Timer struct {
	mem      memoryInterface
	interupt interruptInterface

	counter int
}

const timerCounter = 0xFF05
const timerModulo = 0xFF06
const timerControl = 0xFF07

func CreateTimer(mem memoryInterface, interupt interruptInterface) *Timer {
	return &Timer{
		mem:      mem,
		interupt: interupt,
	}
}

func (t *Timer) Update(cycles uint8) {
	// Divider Register
	current := t.mem.ReadByte(memory.DividerRegister)
	current += cycles / 2
	t.mem.WriteDividerRegister(current)

	// Timer Counter
	control := t.mem.ReadByte(timerControl)
	// Increment if enabled
	if memory.GetBit(control, 2) {

		//frequency := 0
		//switch control & 0x3 {
		//case 0x00:
		//	frequency = 1024
		//case 0x01:
		//	frequency = 16
		//case 0x10:
		//	frequency = 64
		//case 0x11:
		//	frequency = 256
		//}

		counter := t.mem.ReadByte(timerCounter)

		var x uint8
		for x = 0; x < cycles; x++ {
			var newValue uint8

			if uint16(counter)+1 > 0x00FF {
				newValue = t.mem.ReadByte(timerModulo)

				// Request interrupt
				t.interupt.Request(interupt.Time)
			} else {
				newValue = counter + 1
			}

			t.mem.WriteByte(timerCounter, newValue)
		}
	}
}
