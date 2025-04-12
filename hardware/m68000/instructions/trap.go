package instructions

import "fmt"

type trap struct {
	vector byte
}

func CreateTRAP(value uint16) *trap {
	return &trap{
		vector: byte(value & 0b0000000000001111),
	}
}

func (t *trap) Name() string {
	return fmt.Sprintf("TRAP #%d", t.vector)
}

func (t *trap) Length() int {
	return 2
}
