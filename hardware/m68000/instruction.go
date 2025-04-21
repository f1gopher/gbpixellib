package m68000

import (
	"github.com/f1gopher/gbpixellib/hardware/m68000/instructions"
	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type instruction interface {
	Name() string
	Length() int
	Execute(registers *instructions.Registers, memory memory.Memory[uint32])
}
