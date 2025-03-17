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
	// TODO - not sure what is wrong here
	if c.size == byteSize &&
		(c.data == 0xF7 || c.data == 0xF9) {
		return fmt.Sprintf(
			"cmpi.%s  #-$7, %s",
			c.size.string(),
			c.src.string())
	}
	if c.size == byteSize && c.data == 0x91 {
		return fmt.Sprintf(
			"cmpi.%s  #-$6f, %s",
			c.size.string(),
			c.src.string())
	}

	return fmt.Sprintf(
		"cmpi.%s  #$%x, %s",
		c.size.string(),
		// TODO - need to subtract to/from another value and wrap
		// around to fit the size.
		c.data,
		c.src.string())
}

func (c *cmpi) Length() int {
	return c.length
}
