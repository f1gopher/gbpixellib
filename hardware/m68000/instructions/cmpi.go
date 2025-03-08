package instructions

import "fmt"

type cmpi struct {
	length int
	size   instSize
	src    addressModeInfo
	data   uint32
}

func CreateCMPI(value uint16, code []uint8, codeOffset int) *cmpi {
	size := twoBitSize(value, 6)
	data := readData(code, codeOffset, size)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset+size.bytesUsed(), codeOffset)
	return &cmpi{
		length: 2 + instSize + size.bytesUsed(),
		size:   size,
		src:    src,
		data:   data,
	}
}

func (c *cmpi) Name() string {
	return fmt.Sprintf(
		"cmpi.%s  #$%x, %s",
		c.size.string(),
		c.data,
		c.src.string())
}

func (c *cmpi) Length() int {
	return c.length
}
