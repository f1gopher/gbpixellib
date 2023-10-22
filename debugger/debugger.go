package debugger

import (
	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
)

type Debugger struct {
	log    *log.Log
	regs   debugRegisters
	memory debugMemory
}

func CreateDebugger(log *log.Log) (*Debugger, *cpu.Registers, *memory.Memory) {
	r := debugRegisters{registers: &cpu.Registers{}}
	m := debugMemory{memory: memory.CreateMemory(log)}

	d := &Debugger{
		regs:   r,
		memory: m,
	}

	return d, r.registers, m.memory
}

func (d *Debugger) StartCycle() {

}

func (d *Debugger) HasHitBreakpoint() bool {
	return false
}

func (d *Debugger) AddRegisterValueBP(reg cpu.Register, value uint16) {

}
