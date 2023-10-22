package debugger

import "github.com/f1gopher/gbpixellib/cpu"

type debugRegisters struct {
	registers *cpu.Registers
}

func (d *debugRegisters) reset() {
	d.reset()
}

func (d *debugRegisters) Get8(source cpu.Register) uint8 {
	return d.registers.Get8(source)
}

func (d *debugRegisters) Get16(source cpu.Register) uint16 {
	return d.registers.Get16(source)
}

func (d *debugRegisters) Set8(target cpu.Register, value uint8) {
	d.registers.Set8(target, value)
}

func (d *debugRegisters) Set16(target cpu.Register, value uint16) {
	d.registers.Set16(target, value)
}
