package instructions

import "fmt"

type cmpi struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	data    uint32
}

func CreateCMPI(value uint16, code []uint8, offset int) *cmpi {
	size := twoBitSize(value, 6)
	srcMode, srcReg := effective(value, 0, size)
	return &cmpi{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
		data:    readData(code, offset, size),
	}
}

func (c *cmpi) Name() string {
	return fmt.Sprintf(
		"cmpi.%s  #$%x, %s",
		c.size.string(),
		c.data,
		addressString(c.srcMode, c.srcReg, 0, 0))
}

func (c *cmpi) Length() int {
	return 2 + addressModeBytesUsed(c.srcMode).bytesUsed() + c.size.bytesUsed()
}
