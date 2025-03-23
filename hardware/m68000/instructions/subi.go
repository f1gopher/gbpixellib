package instructions

import "fmt"

type subi struct {
	length int
	size   instSize
	src    addressModeInfo
	data   int32
}

func CreateSUBI(value uint16, code []uint8, codeOffset int) *subi {
	size := twoBitSize(value, 6)
	tmp := readData(code, codeOffset, size)
	data := signExtend16BitTo32Bits(uint16(tmp))
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset+size.bytesUsed(), codeOffset)
	return &subi{
		length: 2 + instSize + size.bytesUsed(),
		size:   size,
		src:    src,
		data:   data,
	}
}

func (s *subi) Name() string {
	if s.data < 0 {
		return fmt.Sprintf(
			"subi.%s  #-$%x, %s",
			s.size.string(),
			-s.data,
			s.src.string())
	}

	return fmt.Sprintf(
		"subi.%s  #$%x, %s",
		s.size.string(),
		s.data,
		s.src.string())
}

func (s *subi) Length() int {
	return s.length
}
