package instructions

import "fmt"

type subq struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	data    byte
}

func CreateSUBQ(value uint16, code []uint8, offset int) *subq {
	size := twoBitSize(value, 6)
	srcMode, srcReg := effective(value, 0, size)
	// TODO - do we do this in multiple places so need a helper method?
	data := valueFrom3Bits(value, 9)
	if data == 0 {
		data = 8
	}
	return &subq{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
		data:    data,
	}
}

func (s *subq) Name() string {
	return fmt.Sprintf(
		"subq.%s  #%d, %s",
		s.size.string(),
		s.data,
		addressString(s.srcMode, s.srcReg, uint32(s.data), 0))
}

func (s *subq) Length() int {
	return 2 + s.size.bytesUsed()
}
