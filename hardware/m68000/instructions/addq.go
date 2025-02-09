package instructions

import "fmt"

type addq struct {
	size  instSize
	mode  addressMode
	reg   byte
	value byte
	data  uint32
}

func CreateADDQ(value uint16, code []uint8, offset int) *addq {
	size := twoBitSize(value, 6)
	mode, reg := effective(value, 0, wordSize)
	imediateValue := (value >> 9) & 0b00000111
	// 0 means 8
	if imediateValue == 0 {
		imediateValue = 8
	}

	return &addq{
		size:  size,
		mode:  mode,
		reg:   reg,
		value: byte(imediateValue),
		data:  readData(code, offset, wordSize),
	}
}

func (a *addq) Name() string {
	return fmt.Sprintf(
		"addq.%s  #%d, %s",
		a.size.string(),
		a.value,
		addressString(a.mode, a.reg, a.data, 0))
}

func (a *addq) Length() int {
	// TODO - not sure thi sis right but with long size it only read an extra word
	return 2 + wordSize.bytesUsed()
}
