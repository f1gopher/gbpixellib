package instructions

import "fmt"

type dc struct {
	value uint16
}

func CreateDC(value uint16) *dc {
	return &dc{
		value: value,
	}
}

func (d *dc) Name() string {
	if d.value&0xFF00 == 0xFF00 ||
		d.value == 0xFE01 || d.value == 0xFEE8 ||
		d.value == 0xFEF4 {
		return fmt.Sprintf("dc.w    $%x; opcode 1111", d.value)
	}

	return fmt.Sprintf("dc.w    $%04x; ILLEGAL", d.value)
}

func (d *dc) Length() int {
	return 2
}
