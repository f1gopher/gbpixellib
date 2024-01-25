package debugger

import (
	"strings"

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

type realDebugger struct {
	log    *log.Log
	regs   debugRegisters
	memory debugMemory
}

func createRealDebugger(log *log.Log) (Debugger, cpu.RegistersInterface, cpu.MemoryInterface, *memory.Bus) {
	r := debugRegisters{
		registers:   &cpu.Registers{},
		breakpoints: make(map[cpu.Register][]registerBreakpoint, 0),
	}
	m := debugMemory{
		memory:  memory.CreateBus(log),
		records: make(map[uint16]*memoryRecord, 0),
	}

	d := &realDebugger{
		regs:   r,
		memory: m,
	}

	return d, &d.regs, &d.memory, d.memory.memory
}

func (d *realDebugger) StartCycle(cycle uint) {
	d.regs.startCycle()
	d.memory.startCycle(cycle)
}

func (d *realDebugger) HasHitBreakpoint() bool {
	return d.regs.hasHitBreakpoint() || d.memory.hasHitBreakpoint()
}

func (d *realDebugger) BreakpointReason() string {
	regBP := d.regs.BreakpointReason()
	memBP := d.memory.BreakpointReason()

	if regBP == "" && memBP == "" {
		return ""
	}

	return strings.Join([]string{regBP, memBP}, "\n")
}

func (d *realDebugger) DisableAllBreakpoints() {
	for _, bps := range d.regs.breakpoints {
		for x, _ := range bps {
			bps[x].enabled = false
		}
	}
	for _, bps := range d.memory.breakpoints {
		for x, _ := range bps {
			bps[x].enabled = false
		}
	}
}

func (d *realDebugger) AddRegisterValueBP(reg cpu.Register, comparison BreakpointComparison, value uint16) {
	d.regs.addBP(reg, comparison, value)
}

func (d *realDebugger) AddMemoryBP(address uint16, comparison BreakpointComparison, value uint8) {
	d.memory.addBP(address, comparison, value)
}

func (d *realDebugger) AddMemoryRecorder(address uint16) {
	d.memory.addRecorder(address)
}

func (d *realDebugger) DeleteMemoryRecorder(address uint16) {
	d.memory.deleteRecorder(address)
}

func (d *realDebugger) MemoryRecordValues(address uint16) []MemoryRecordEntry {
	return d.memory.recordValues(address)
}
