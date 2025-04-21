package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type direction int

const (
	right direction = iota
	left
)

type rotateShiftLocation int

const (
	immediate rotateShiftLocation = iota
	register
)

type ror_rol struct {
	// Common
	location  rotateShiftLocation
	direction direction

	// Register format
	size instSize
	data byte

	// Immediate format
	src addressModeInfo
}

func CreateROL_ROR(value uint16, code []uint8, codeOffset int) *ror_rol {
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

		return &ror_rol{
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

	return &ror_rol{
		location:  location,
		direction: direction,
		size:      size,
		src:       src,
		data:      countReg,
	}
}

func (r *ror_rol) Name() string {
	name := "rol"
	if r.direction == right {
		name = "ror"
	}

	if r.location == register {
		return fmt.Sprintf(
			"%s.%s   #%d, %s",
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

func (r *ror_rol) Length() int {
	return 2
}

func (r *ror_rol) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
