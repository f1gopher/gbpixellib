package instructions

import (
	"fmt"
)

type eori struct {
	size     instSize
	destMode addressMode
	destReg  byte
	data     uint32
}

func CreateEORI(value uint16, code []uint8, offset int) *eori {
	size := twoBitSize(value, 6)
	destMode, destReg := effective(value, 0, size)
	return &eori{
		size:     size,
		destMode: destMode,
		destReg:  destReg,
		data:     readData(code, offset, size),
	}
}

func (e *eori) Name() string {
	return fmt.Sprintf(
		"eori.%s  %s, %s",
		e.size.string(),
		addressString(preIndexedImediate, 0, e.data, 0),
		addressString(e.destMode, e.destReg, 0, 0))
}

func (e *eori) Length() int {
	return 2 + e.size.bytesUsed()
}
