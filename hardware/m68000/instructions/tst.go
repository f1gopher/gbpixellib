package instructions

import "fmt"

type tst struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	data    uint32
}

func CreateTST(value uint16, code []uint8, offset int) *tst {
	size := twoBitSize(value, 6)
	srcMode, srcReg := effective(value, 0, size)
	return &tst{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
		data:    readData(code, offset, wordSize),
	}
}

func (t *tst) Name() string {
	return fmt.Sprintf(
		"tst.%s   %s",
		t.size.string(),
		addressString(t.srcMode, t.srcReg, t.data, 0))
}

func (t *tst) Length() int {
	return 2 + t.size.bytesUsed()
}
