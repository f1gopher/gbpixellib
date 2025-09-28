package m68000

import "github.com/f1gopher/gbpixellib/hardware/memory"

type illegal struct{}

func CreateILLEGAL() *illegal {
	return &illegal{}
}

func (i *illegal) Name() string {
	return ""
}

func (i *illegal) Length() int {
	return 2
}

func (i *illegal) StepClockCycle(registers *Registers, memory memory.Memory[uint32], pc uint32, prefetch prefetchQueue) (finished bool, newPC uint32) {
	return true, pc
}
