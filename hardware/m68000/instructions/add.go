package instructions

import "fmt"

type add struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	mode    byte
	destReg byte
	data    uint32
}

func CreateADD(value uint16, code []uint8, offset int) *add {
	mode := valueFrom3Bits(value, 6)
	size := byteSize
	if mode == 0b001 || mode == 0b101 {
		size = wordSize
	} else if mode == 0b010 || mode == 0b110 {
		size = longSize
	}
	srcMode, srcReg := effective(value, 0, size)
	destReg := valueFrom3Bits(value, 9)

	return &add{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
		mode:    mode,
		destReg: destReg,
		data:    readData(code, offset, size),
	}
}

func (a *add) Name() string {
	return fmt.Sprintf(
		"add.%s   %s, D%d",
		a.size.string(),
		addressString(a.srcMode, a.srcReg, a.data, 0),
		a.destReg,
	)
}

func (a *add) Length() int {
	return 2 + addressModeBytesUsed(a.srcMode).bytesUsed()
}
