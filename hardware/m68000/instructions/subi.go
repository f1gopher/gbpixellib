package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type subi struct {
	length int
	size   instSize
	src    addressModeInfo
	data   int32
}

func CreateSUBI(value uint16, code []uint8, codeOffset int) *subi {
	size := twoBitSize(value, 6)

	if size == zeroSize {
		return nil
	}

	tmp := readData(code, codeOffset, size)
	var data int32
	if size == wordSize {
		data = signExtend16BitTo32Bits(uint16(tmp))
	} else if size == longSize {
		data = twosComplementToInt32Bit(tmp)
	} else {
		data = int32(twosComplementToInt8Bit(uint8(tmp)))
	}
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

func (s *subi) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
