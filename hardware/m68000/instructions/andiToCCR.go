package instructions

import "fmt"

type andiToCCR struct {
	data byte
}

func CreateANDIToCCR(value uint16, code []uint8, codeOffset int) *andiToCCR {
	return &andiToCCR{
		data: byte(readData(code, codeOffset, byteSize)),
	}
}

func (a *andiToCCR) Name() string {
	return fmt.Sprintf("ANDI #0x%02X,CCR", a.data)
}

func (a *andiToCCR) Length() int {
	return 4
}
