package instructions

import (
	"fmt"
)

type addi struct {
	size     instSize
	destMode addressMode
	destReg  byte

	data uint32
}

func CreateADDI(value uint16, code []uint8, offset int) *addi {
	size := twoBitSize(value, 6)
	destMode, destReg := effective(value, 0, size)
	data := readData(code, offset, size)

	return &addi{
		size:     size,
		destMode: destMode,
		destReg:  destReg,
		data:     data,
	}
}

func (a *addi) Name() string {
	return fmt.Sprintf(
		"addi.%s  #$%x, %s",
		a.size.string(),
		a.data,
		addressString(a.destMode, a.destReg, a.data, 0),
	)
}

func (a *addi) Length() int {
	return 2 + a.size.bytesUsed()
}
