package instructions

import "fmt"

type roxl_roxr struct {
	// Common
	location  rotateShiftLocation
	direction direction

	// Register format
	size instSize
	data byte

	// Immediate format
	src addressModeInfo
}

func CreateROXL_ROXR(value uint16, code []uint8, codeOffset int) *roxl_roxr {
	location := register
	if value&0b1111111011000000 == 0b1110011011000000 {
		location = immediate
	}

	direction := right
	if value&0b0000000100000000 == 0b0000000100000000 {
		direction = left
	}

	if location == immediate {
		src, _, _ := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

		return &roxl_roxr{
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

	return &roxl_roxr{
		location:  location,
		direction: direction,
		size:      size,
		src:       src,
		data:      countReg,
	}
}

func (r *roxl_roxr) Name() string {
	name := "roxl"
	if r.direction == right {
		name = "roxr"
	}

	if r.location == register {
		return fmt.Sprintf(
			"%s.%s  #%d, %s",
			name,
			r.size.string(),
			r.data,
			r.src.string(),
		)
	}

	return fmt.Sprintf(
		"%s.   ",
		name,
	)
}

func (r *roxl_roxr) Length() int {
	return 2
}
