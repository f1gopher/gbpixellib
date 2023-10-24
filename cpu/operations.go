package cpu

func add8BitWithCarry(original uint8, add uint8) (result uint8, bit3Carry bool, bit7carry bool) {
	result = original + add

	//abc := result ^ 0x01 ^ original
	abc := result ^ add ^ original
	bit3Carry = abc&0x10 == 0x10
	bit7carry = result < original
	return
}

func add8BitAndCarryWithCarry(original uint8, add uint8, carry bool) (result uint8, bit3Carry bool, bit7Carry bool) {

	if !carry {
		return add8BitWithCarry(original, add)
	}

	result = original + add + 1

	//abc := result ^ 0x01 ^ original

	bit3Carry = (original&0xF)+(add&0xF)+0x1 > 0xF
	bit7Carry = uint16(original&0xFF)+(uint16(add)&0xFF)+0x01 > 0xFF
	//bit3Carry = abc&0x10 == 0x10
	//bit7carry = result < original
	return
}

func subtract8BitWithCarry(original uint8, subtract uint8) (result uint8, bit3Carry bool, bit7carry bool) {
	result = original - subtract

	//abc := result ^ 0x01 ^ original
	//bit3Carry = abc&0x10 == 0x10

	// TODO - no idea about this!!
	bit3Carry = original&0xF < subtract&0xF
	bit7carry = original < subtract
	return
}

func subtract8BitAndCarryWithCarry(original uint8, subtract uint8, carry bool) (result uint8, bit3Carry bool, bit7carry bool) {

	if !carry {
		return subtract8BitWithCarry(original, subtract)
	}

	result = original - 1 - subtract

	//abc := result ^ 0x01 ^ original
	//bit3Carry = abc&0x10 == 0x10
	//bit7carry = false // TODO

	// TODO - no idea about this!
	bit3Carry = original&0xF < subtract&0xF
	bit7carry = original < subtract
	return
}

func add16BitWithCarry(original uint16, add uint16) (result uint16, bit11Carry bool, bit15carry bool) {
	result = original + add

	abc := result ^ original ^ add
	bit11Carry = abc&0x1000 == 0x1000
	bit15carry = result < original
	return
}

func add16Bit(original uint16, add uint16) (result uint16) {
	result = original + add
	return
}

func subtract16BitWithCarry(original uint16, subtract uint16) (result uint16, bit3Carry bool, bit7carry bool) {
	result = original - subtract

	abc := result ^ 0x01 ^ original
	bit3Carry = abc&0x10 == 0x10
	bit7carry = false // TODO
	return
}

func subtract16Bit(original uint16, subtract uint16) (result uint16) {
	result = original - subtract
	return
}

func readAndIncPC(reg registersInterface, mem memoryInterface) uint8 {
	pc := reg.Get16(PC)
	result := mem.ReadByte(pc)
	pc++
	reg.Set16(PC, pc)

	return result
}

func readAndIncSP(reg registersInterface, mem memoryInterface) uint8 {
	sp := reg.Get16(SP)
	result := mem.ReadByte(sp)
	sp++
	reg.Set16(SP, sp)

	return result
}

func readAndDecSP(reg registersInterface, mem memoryInterface) uint8 {
	sp := reg.Get16(SP)
	result := mem.ReadByte(sp)
	sp--
	reg.Set16(SP, sp)

	return result
}

func DecAndWriteSP(reg registersInterface, mem memoryInterface, value uint8) {
	sp := reg.Get16(SP)
	sp--
	mem.WriteByte(sp, value)
	reg.Set16(SP, sp)
}

func incPC(reg registersInterface) {
	pc := reg.Get16(PC)
	reg.Set16(PC, pc+1)
}

func decSP(reg registersInterface) {
	sp := reg.Get16(SP)
	reg.Set16(SP, sp-1)
}

func combineBytes(msb uint8, lsb uint8) uint16 {
	value := uint16(msb)
	value = value << 8
	value = value | uint16(lsb)
	return value
}

func Msb(value uint16) uint8 {
	return uint8(value >> 8)
}

func Lsb(value uint16) uint8 {
	return uint8(value)
}

func adds8Tou16(value uint16, add int8) uint16 {
	if add < 0 {
		return value - uint16(-add)
	}

	return value + uint16(add)
}
