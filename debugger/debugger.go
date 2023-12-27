package debugger

import (
	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
)

type BreakpointComparison int

const (
	Equal BreakpointComparison = iota
	NotEqual
	GreaterThan
	LessThan
	GreaterThanOrEqual
	LessThanOrEqual
)

type Debugger struct {
	log    *log.Log
	regs   debugRegisters
	memory debugMemory
}

func CreateDebugger(log *log.Log) (*Debugger, cpu.RegistersInterface, *memory.Bus) {
	r := debugRegisters{
		registers:   &cpu.Registers{},
		breakpoints: make(map[cpu.Register][]registerBreakpoint, 0),
	}
	m := debugMemory{memory: memory.CreateBus(log)}

	d := &Debugger{
		regs:   r,
		memory: m,
	}

	return d, &d.regs, m.memory
}

func (d *Debugger) StartCycle() {
	d.regs.startCycle()
}

func (d *Debugger) HasHitBreakpoint() bool {
	return d.regs.hasHitBreakpoint()
}

func (d *Debugger) BreakpointReason() string {
	return d.regs.BreakpointReason()
}

func (d *Debugger) AddRegisterValueBP(reg cpu.Register, comparison BreakpointComparison, value uint16) {
	d.regs.AddBP(reg, comparison, value)
}
