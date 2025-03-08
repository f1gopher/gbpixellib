package instructions

import "fmt"

type lsl_lsr struct {
	// Common
	location  rotateShiftLocation
	direction direction

	// Register shift
	size instSize
	data byte

	// Immediate shift
	src addressModeInfo
}

func CreateLSL_LSR(value uint16, code []uint8, codeOffset int) *lsl_lsr {
	location := register
	if value&0b1111111011000000 == 0b1110001011000000 {
		location = immediate
	}

	direction := right
	if value&0b000000100000000 == 0b0000000100000000 {
		direction = left
	}

	if location == immediate {
		src, _, _ := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

		return &lsl_lsr{
			location:  location,
			direction: direction,
			src:       src,
		}
	}

	rotateCountLocation := immediate
	if value&0b0000000000100000 == 0b0000000000100000 {
		rotateCountLocation = register
	}
	size := twoBitSize(value, 6)
	src := addressModeInfo{
		reg:  valueFrom3Bits(value, 0),
		mode: dataRegisterDirect,
	}
	countReg := valueFrom3Bits(value, 9)
	if rotateCountLocation == immediate && countReg == 0 {
		countReg = 8
	}
	// TODO - else modulo 64 to get the register?

	return &lsl_lsr{
		location:  location,
		direction: direction,
		size:      size,
		src:       src,
		data:      countReg,
	}
}

func (l *lsl_lsr) Name() string {
	name := "lsl"
	if l.direction == right {
		name = "lsr"
	}

	if l.location == register {
		return fmt.Sprintf(
			"%s.%s   #%d, %s",
			name,
			l.size.string(),
			l.data,
			l.src.string(),
		)
	}

	return fmt.Sprintf(
		"%s.   ",
		name,
	)
}

func (l *lsl_lsr) Length() int {
	return 2
}
