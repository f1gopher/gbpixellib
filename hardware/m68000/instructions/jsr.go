package instructions

import "fmt"

type jsr struct {
	size    instSize
	srcMode addressMode
	srcReg  byte

	data uint32
}

func CreateJSR(value uint16, code []uint8, offset int) *jsr {
	srcMode, srcReg := effective(value, 0, wordSize)
	size := addressModeBytesUsed(srcMode)
	return &jsr{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
		data:    readData(code, offset, size),
	}
}

func (j *jsr) Name() string {
	return fmt.Sprintf("jsr     %s", addressString(j.srcMode, j.srcReg, uint32(j.data), 0))
}

func (j *jsr) Length() int {
	return 2 + j.size.bytesUsed()
}
