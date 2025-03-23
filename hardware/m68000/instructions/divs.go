package instructions

import "fmt"

type divs struct {
length int
src addressModeInfo
reg byte
}

func CreateDIVS(value uint16, code []uint8, codeOffset int) *divs {
	src, _, instSize := readEffectiveAddressMode(value, 0, byteSize, code, codeOffset, codeOffset)
	reg := valueFrom3Bits(value, 9)

	return &divs {
length: 2 + instSize,
src: src,
reg: reg,
	}
}

func (d *divs) Name() string {
	return fmt.Sprintf(
	"divs.w  %s, D%d",
	d.src.string(),
	d.reg)
}

func (d *divs) Length() int {
	return d.length
}
