package m68000

import (
	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
)

type M68000 struct {
	log *log.Log

	registers registers
	memory    memory.Memory[uint32]
}

func CreateM68000(log *log.Log, memory memory.Memory[uint32]) *M68000 {
	return &M68000{
		log:       log,
		memory:    memory,
		registers: registers{},
	}
}

func (m *M68000) Init() {

}

func (m *M68000) Reset() {

}

func (m *M68000) ExecuteCycle() (breakpoint bool) {
	return false
}
