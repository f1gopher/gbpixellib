package instructions

import "fmt"

type clr struct {
	size     instSize
	destMode addressMode
	destReg  byte
}

func CreateCLR(value uint16) *clr {
	size := twoBitSize(value, 6)
	mode, reg := effective(value, 0, size)

	return &clr{
		size:     size,
		destMode: mode,
		destReg:  reg,
	}
}

func (c *clr) Name() string {
	return fmt.Sprintf("clr.%s   %s", c.size.string(), addressString(c.destMode, c.destReg, 0, 0))
}

func (c *clr) Length() int {
	return 2
}
