package instructions

import "fmt"

type clr struct {
	size     instSize
	destMode addressMode
	destReg  byte

	displacement int
}

func CreateCLR(value uint16, code []uint8, offset int) *clr {
	size := twoBitSize(value, 6)
	mode, reg := effective(value, 0, size)
	displacement := twosComplementToInt16Bit(uint16(readData(code, offset, wordSize)))

	return &clr{
		size:         size,
		destMode:     mode,
		destReg:      reg,
		displacement: displacement,
	}
}

func (c *clr) Name() string {
	return fmt.Sprintf(
		"clr.%s   %s",
		c.size.string(),
		addressString(c.destMode, c.destReg, 0, c.displacement))
}

func (c *clr) Length() int {
	return 2 + addressModeBytesUsed(c.destMode).bytesUsed()
}
