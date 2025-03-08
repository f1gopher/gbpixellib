package instructions

import "fmt"

type subq struct {
	length int
	size   instSize
	src    addressModeInfo
	data   byte
}

func CreateSUBQ(value uint16, code []uint8, codeOffset int) *subq {
	size := twoBitSize(value, 6)
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)
	// TODO - do we do this in multiple places so need a helper method?
	data := valueFrom3Bits(value, 9)
	if data == 0 {
		data = 8
	}
	return &subq{
		length: 2 + instSize,
		src:    src,
		size:   size,
		data:   data,
	}
}

func (s *subq) Name() string {
	return fmt.Sprintf(
		"subq.%s  #%d, %s",
		s.size.string(),
		s.data,
		s.src.string())
}

func (s *subq) Length() int {
	return s.length
}
