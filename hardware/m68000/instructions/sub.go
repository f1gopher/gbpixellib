package instructions

import "fmt"

type sub struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	mode    byte
	destReg byte
	data    uint32
}

func CreateSUB(value uint16, code []uint8, offset int) *sub {
	mode := valueFrom3Bits(value, 6)
	size := byteSize
	if mode == 0b001 || mode == 0b101 {
		size = wordSize
	} else if mode == 0b010 || mode == 0b110 {
		size = longSize
	}
	srcMode, srcReg := effective(value, 0, size)
	destReg := valueFrom3Bits(value, 9)

	return &sub{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
		mode:    mode,
		destReg: destReg,
		data:    readData(code, offset, size),
	}
}

func (s *sub) Name() string {
	return fmt.Sprintf(
		"sub.%s   %s, D%d",
		s.size.string(),
		addressString(s.srcMode, s.srcReg, s.data, 0),
		s.destReg,
	)
}

func (s *sub) Length() int {
	return 2 + addressModeBytesUsed(s.srcMode).bytesUsed()
}
