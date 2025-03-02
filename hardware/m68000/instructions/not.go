package instructions

import (
	"fmt"
)

type not struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
}

func CreateNOT(value uint16) *not {
	size := twoBitSize(value, 6)
	srcMode, srcReg := effective(value, 0, size)

	return &not{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
	}
}

func (n *not) Name() string {
	return fmt.Sprintf(
		"not.%s   %s",
		n.size.string(),
		addressString(n.srcMode, n.srcReg, 0, 0))
}

func (n *not) Length() int {
	return 2
}
