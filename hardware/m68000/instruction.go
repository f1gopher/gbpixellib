package m68000

import (
	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type prefetchQueue interface {
	ReadNextValue() uint16
}

type instruction interface {
	Name() string
	Length() int
	StepClockCycle(registers *Registers, memory memory.Memory[uint32], pc uint32, prefetch prefetchQueue) (finished bool, newPC uint32)
}
