package instructions

import "fmt"

type ori struct {
	size instSize
	mode addressMode
	reg  byte

	data uint32
}

func CreateORI(value uint16, code []uint8, offset int) *ori {
	mode, reg := effective(value, 0)
	size := twoBitSize(value, 6)

	return &ori{
		size: size,
		mode: mode,
		reg:  reg,
		data: readData(code, offset, size),
	}
}

func (o *ori) Name() string {
	return fmt.Sprintf("ori.%s   #$%x, %s%d", o.size.string(), o.data, o.mode.string(), o.reg)
}

func (o *ori) Length() int {
	return 2 + o.size.bytesUsed()
}
