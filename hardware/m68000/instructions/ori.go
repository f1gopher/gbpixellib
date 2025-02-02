package instructions

import "fmt"

type ori struct {
	size instSize
	mode addressMode
	reg  byte

	data uint32
}

func CreateORI(value uint16, code []uint8, offset int) *ori {
	size := twoBitSize(value, 6)
	mode, reg := effective(value, 0, size)

	return &ori{
		size: size,
		mode: mode,
		reg:  reg,
		data: readData(code, offset, size),
	}
}

func (o *ori) Name() string {
	return fmt.Sprintf(
		"ori.%s   #$%x, %s",
		o.size.string(),
		o.data,
		addressString(o.mode, o.reg, 0, 0))
}

func (o *ori) Length() int {
	return 2 + o.size.bytesUsed()
}
