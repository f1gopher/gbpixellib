package instructions

import "fmt"

type scc struct {
	mode         addressMode
	reg          byte
	condition    testCondition
	displacement int

	tmpStartOffset int
}

func CreateScc(value uint16, code []uint8, offset int) *scc {
	mode, reg := effective(value, 0, wordSize)
	displacement := twosComplementToInt16Bit(uint16(readData(code, offset, addressModeBytesUsed(mode))))

	return &scc{
		mode:           mode,
		reg:            reg,
		condition:      condition(value),
		displacement:   displacement,
		tmpStartOffset: offset,
	}
}

func (s *scc) Name() string {
	return fmt.Sprintf(
		"s%s      %s",
		s.condition.string(),
		addressString(s.mode, s.reg, 0, s.displacement),
	)
}

func (s *scc) Length() int {
	return 2 + addressModeBytesUsed(s.mode).bytesUsed()
}
