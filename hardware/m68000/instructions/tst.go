package instructions

import "fmt"

type tst struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	data    uint32

	displacement int
}

func CreateTST(value uint16, code []uint8, offset int) *tst {
	size := twoBitSize(value, 6)
	srcMode, srcReg := effective(value, 0, size)
	displacement := twosComplementToInt16Bit(uint16(readData(code, offset, wordSize)))
	return &tst{
		size:         size,
		srcMode:      srcMode,
		srcReg:       srcReg,
		data:         readData(code, offset, wordSize),
		displacement: displacement,
	}
}

func (t *tst) Name() string {
	return fmt.Sprintf(
		"tst.%s   %s",
		t.size.string(),
		addressString(t.srcMode, t.srcReg, t.data, t.displacement))
}

func (t *tst) Length() int {
	return 2 + t.size.bytesUsed()
}
