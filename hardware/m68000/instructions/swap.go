package instructions

import "fmt"

type swap struct {
	reg byte
}

func CreateSWAP(value uint16, code []uint8, codeOffset int) *swap {
	return &swap{
		reg: valueFrom3Bits(value, 0),
	}
}

func (s *swap) Name() string {
	return fmt.Sprintf("swap    D%d", s.reg)
}

func (s *swap) Length() int {
	return 2
}
