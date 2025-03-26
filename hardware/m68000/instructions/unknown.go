package instructions

import "fmt"

type unknown struct {
	value uint16
}

func CreateUNKNOWN(value uint16) *unknown {
	return &unknown{value: value}
}

func (u *unknown) Name() string {
	return fmt.Sprintf("unknown instruction: %x", u.value)
}

func (u *unknown) Length() int {
	return 2
}
