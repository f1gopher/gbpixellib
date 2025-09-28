package m68000

import "github.com/f1gopher/gbpixellib/hardware/memory"

type nop struct {
	cycle uint8
}

func CreateNOP(value uint16, code []uint8, codeOffset int) *nop {
	return &nop{}
}

func (n *nop) Name() string {
	return "nop"
}

func (n *nop) Length() int {
	return 2
}

func (n *nop) StepClockCycle(registers *Registers, memory memory.Memory[uint32], pc uint32, prefetch prefetchQueue) (finished bool, newPC uint32) {
	// Doesn't do anything but takes 4 cycles
	n.cycle++
	return n.cycle == 4, pc
}
