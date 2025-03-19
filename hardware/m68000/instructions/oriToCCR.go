package instructions

import "fmt"

type oriToCCR struct {
	data byte
}

func CreateORIToCCR(value uint16, code []uint8, codeOffset int) *oriToCCR {
	return &oriToCCR{
		data: byte(readData(code, codeOffset, byteSize)),
	}
}

func (o *oriToCCR) Name() string {
	return fmt.Sprintf("ori     #$%x, CCR", o.data)
}

func (o *oriToCCR) Length() int {
	return 4
}
