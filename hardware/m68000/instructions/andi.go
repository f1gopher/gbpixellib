package instructions

import "fmt"

type andi struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	data    uint32
}

func CreateANDI(value uint16, code []uint8, offset int) *andi {
	size := twoBitSize(value, 6)
	srcMode, srcReg := effective(value, 0, size)
	return &andi{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
		data:    readData(code, offset, size),
	}
}

func (a *andi) Name() string {
	return fmt.Sprintf(
		"andi.%s  #$%d, %s",
		a.size.string(),
		a.data,
		addressString(a.srcMode, a.srcReg, 0, 0),
	)
}

func (a *andi) Length() int {
	return 2 + a.size.bytesUsed()
}
