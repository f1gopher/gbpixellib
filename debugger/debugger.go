package debugger

import (
	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
)

type BreakpointComparison int

const (
	Equals BreakpointComparison = iota
	GreaterThanOrEqual
)

type Debugger struct {
	log    *log.Log
	regs   debugRegisters
	memory debugMemory
}

func CreateDebugger(log *log.Log) (*Debugger, *cpu.Registers, *memory.Bus) {
	r := debugRegisters{
		registers:   &cpu.Registers{},
		breakpoints: make(map[cpu.Register][]registerBreakpoint, 0),
	}
	m := debugMemory{memory: memory.CreateBus(log)}

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

func (d *Debugger) AddRegisterValueBP(reg cpu.Register, comparison BreakpointComparison, value uint8) {
	d.regs.AddBP(reg, comparison, value)
}
