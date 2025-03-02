package instructions

import "fmt"

type rotateDirection int

const (
	right rotateDirection = iota
	left
)

type rotateLocation int

const (
	immediate rotateLocation = iota
	register
)

type ror_rol struct {
	// Common
	location  rotateLocation
	direction rotateDirection

	// Register format
	size instSize
	data byte

	// Immediate format
	srcMode addressMode
	srcReg  byte
}

func CreateROR_ROL(value uint16, code []uint8, offset int) *ror_rol {
	location := register
	if value&0b1111111011000000 == 0b1110011011000000 {
		location = immediate
	}

	direction := right
	if value&0b0000000010000000 == 0b0000000100000000 {
		direction = left
	}

	if location == immediate {
		srcMode, srcReg := effective(value, 0, wordSize)

		return &ror_rol{
			location:  location,
			direction: direction,
			srcMode:   srcMode,
			srcReg:    srcReg,
		}
	}

	rotateCountLocation := immediate
	if value&0b0000000000100000 == 0b0000000000100000 {
		rotateCountLocation = register
	}
	size := twoBitSize(value, 6)
	register := valueFrom3Bits(value, 0)
	countReg := valueFrom3Bits(value, 9)
	if rotateCountLocation == immediate && countReg == 0 {
		countReg = 8
	}
	// TODO - else modulo 64 to get the register?

	return &ror_rol{
		location:  location,
		direction: direction,
		size:      size,
		srcReg:    register,
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
			addressString(r.srcMode, r.srcReg, 0, 0),
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
