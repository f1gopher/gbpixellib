package instructions

import "fmt"

type eor struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	mode    byte
	destReg byte
	data    uint32
}

func CreateEOR(value uint16, code []uint8, offset int) *eor {
	mode := valueFrom3Bits(value, 6)
	size := byteSize
	if mode == 0b001 || mode == 0b101 {
		size = wordSize
	} else if mode == 0b010 || mode == 0b110 {
		size = longSize
	}
	srcMode, srcReg := effective(value, 0, size)
	destReg := valueFrom3Bits(value, 9)

	return &eor{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
		mode:    mode,
		destReg: destReg,
		data:    readData(code, offset, size),
	}
}

func (e *eor) Name() string {
	return fmt.Sprintf(
		"eor.%s   D%d, %s",
		e.size.string(),
		e.destReg,
		addressString(e.srcMode, e.srcReg, e.data, 0),
	)
}

func (e *eor) Length() int {
	return 2 + addressModeBytesUsed(e.srcMode).bytesUsed()
}
