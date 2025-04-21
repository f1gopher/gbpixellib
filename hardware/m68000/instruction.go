package m68000

import (
	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type instruction interface {
	Name() string
	Length() int
	Execute(registers *Registers, memory memory.Memory[uint32])
}
