package instructions

import "fmt"

type stop struct {
	immediate uint16
}

func CreateSTOP(value uint16, code []uint8, offset int) *stop {
	return &stop{
		immediate: uint16(readData(code, offset, wordSize)),
	}
}

func (s *stop) Name() string {
	return fmt.Sprintf("stop    #$%x", s.immediate)
}

func (s *stop) Length() int {
	return 4
}
