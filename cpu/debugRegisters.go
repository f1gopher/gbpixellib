package cpu

type debugRegisters struct {
	registers Registers
}

func (d *debugRegisters) reset() {
	d.reset()
}

func (d *debugRegisters) get8(source register) uint8 {
	return d.registers.Get8(source)
}

func (d *debugRegisters) get16(source register) uint16 {
	return d.registers.Get16(source)
}

func (d *debugRegisters) set8(target register, value uint8) {
	d.registers.set8(target, value)
}

func (d *debugRegisters) set16(target register, value uint16) {
	d.registers.set16(target, value)
}
