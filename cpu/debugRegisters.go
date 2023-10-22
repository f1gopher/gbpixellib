package cpu

type debugRegisters struct {
	registers Registers
}

func (d *debugRegisters) reset() {
	d.reset()
}

func (d *debugRegisters) get8(source Register) uint8 {
	return d.registers.Get8(source)
}

func (d *debugRegisters) get16(source Register) uint16 {
	return d.registers.Get16(source)
}

func (d *debugRegisters) set8(target Register, value uint8) {
	d.registers.Set8(target, value)
}

func (d *debugRegisters) set16(target Register, value uint16) {
	d.registers.Set16(target, value)
}
