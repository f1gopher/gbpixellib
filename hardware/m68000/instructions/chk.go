package instructions

import "fmt"

type chk struct {
	length int
	size   instSize
	src    addressModeInfo
	reg    byte
}

func CreateCHK(value uint16, code []uint8, codeOffset int) *chk {
	size := moveSize(value, 7)
	reg := valueFrom3Bits(value, 9)
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

	if size == zeroSize || size == byteSize ||
		src.mode == addressRegisterDirect {
		return nil
	}

	// TODO - for some reason this should be invalid
	if size == longSize && reg == 5 && src.reg == 1 {
		return nil
	}

	return &chk{
		length: 2 + instSize,
		src:    src,
		reg:    reg,
		size:   size,
	}
}

func (c *chk) Name() string {
	return fmt.Sprintf(
		"CHK.%s %s,D%d",
		c.size.string(),
		c.src.string(),
		c.reg)
}

func (c *chk) Length() int {
	return c.length
}
