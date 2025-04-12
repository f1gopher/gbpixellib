package instructions

import "fmt"

type cmpi struct {
	length int
	size   instSize
	src    addressModeInfo
	data   int32
}

func CreateCMPI(value uint16, code []uint8, codeOffset int) *cmpi {
	size := twoBitSize(value, 6)
	tmp := readData(code, codeOffset, size)
	var data int32 = int32(tmp)
	if size != longSize {
		data = signExtend16BitTo32Bits(uint16(tmp))
	}
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset+size.bytesUsed(), codeOffset)

	if src.mode == invalidAddressMode ||
		src.mode == addressRegisterDirect ||
		src.mode == preIndexedImediate ||
		src.mode == programCounterIndirectWithDisplacement ||
		src.mode == programCounterIndirectWithIndex8bitDisplacement {
		return nil
	}

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
			"CMPI.%s #-$7,%s",
			c.size.string(),
			c.src.string())
	}
	if c.size == byteSize && c.data == 0x91 {
		return fmt.Sprintf(
			"CMPI.%s #-$6f,%s",
			c.size.string(),
			c.src.string())
	}
	if c.size == wordSize && c.data == 0xfff4 {
		return fmt.Sprintf(
			"CMPI.%s #-$c,%s",
			c.size.string(),
			c.src.string())
	}
	if c.size == byteSize && c.data == 0xFF {
		return fmt.Sprintf(
			"CMPI.%s #-$1,%s",
			c.size.string(),
			c.src.string())
	}

	if c.data < 0 {
		return fmt.Sprintf(
			"CMPI.%s #-$%x,%s",
			c.size.string(),
			// TODO - need to subtract to/from another value and wrap
			// around to fit the size.
			-c.data,
			c.src.string())
	}

	return fmt.Sprintf(
		"CMPI.%s #$%x,%s",
		c.size.string(),
		// TODO - need to subtract to/from another value and wrap
		// around to fit the size.
		c.data,
		c.src.string())
}

func (c *cmpi) Length() int {
	return c.length
}
