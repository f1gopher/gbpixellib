package instructions

import "fmt"

type clr struct {
	length int
	size   instSize
	dest   addressModeInfo
}

func CreateCLR(value uint16, code []uint8, codeOffset int) *clr {
	size := twoBitSize(value, 6)
	dest, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)

	return &clr{
		length: 2 + instSize,
		size:   size,
		dest:   dest,
	}
}

func (c *clr) Name() string {
	return fmt.Sprintf(
		"clr.%s   %s",
		c.size.string(),
		c.dest.string())
}

func (c *clr) Length() int {
	return c.length
}
