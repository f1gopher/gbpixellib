package debugger

import (
	"errors"
	"fmt"
	"sync"

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

func (b BreakpointComparison) String() string {
	return [...]string{"==", "!=", ">", "<", ">=", "<="}[b]
}

func ParseBreakpointComparison(value string) (bp BreakpointComparison, err error) {
	switch value {
	case Equal.String():
		return Equal, nil
	case NotEqual.String():
		return NotEqual, nil
	case LessThan.String():
		return LessThan, nil
	case GreaterThan.String():
		return GreaterThan, nil
	case LessThanOrEqual.String():
		return LessThanOrEqual, nil
	case GreaterThanOrEqual.String():
		return GreaterThanOrEqual, nil
	default:
		return Equal, errors.New(fmt.Sprintf("Unknown breakpoint condition: '%s'", value))
	}
}

type realDebugger struct {
	log    *log.Log
	regs   debugRegisters
	memory debugMemory
}

func createRealDebugger(log *log.Log) (Debugger, cpu.RegistersInterface, cpu.MemoryInterface, *memory.Bus) {
	r := debugRegisters{
		registers:   &cpu.Registers{},
		breakpoints: make(map[cpu.Register][]registerBreakpoint, 0),
		bpLock:      sync.RWMutex{},
	}
	m := debugMemory{
		memory:      memory.CreateBus(log),
		records:     make(map[uint16]*memoryRecord, 0),
		breakpoints: make(map[uint16][]memoryBreakpoint),
		bpLock:      sync.RWMutex{},
	}

	d := &realDebugger{
		regs:   r,
		memory: m,
	}

	return d, &d.regs, &d.memory, d.memory.memory
}

func (d *realDebugger) StartCycle(cycle uint, pc uint16) {
	d.regs.startCycle()
	d.memory.startCycle(cycle, pc)
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

	if regBP != "" && memBP == "" {
		return regBP
	}

	if regBP == "" && memBP != "" {
		return memBP
	}

	return regBP + " and " + memBP
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

func (d *realDebugger) AddRegisterValueBP(
	reg cpu.Register,
	comparison BreakpointComparison,
	value uint16,
	hitCount uint) (id int, err error) {

	return d.regs.addBP(reg, comparison, value, hitCount)
}

func (d *realDebugger) DeleteRegisterBP(id int) {
	d.regs.deleteBP(id)
}

func (d *realDebugger) SetEnabledRegisterBP(id int, enabled bool) {
	d.regs.setEnabledBP(id, enabled)
}

func (d *realDebugger) AddMemoryBP(
	address uint16,
	comparison BreakpointComparison,
	value uint8,
	hitCount uint) (id int, err error) {

	return d.memory.addBP(address, comparison, value, hitCount)
}

func (d *realDebugger) DeleteMemoryBP(id int) {
	d.memory.deleteBP(id)
}

func (d *realDebugger) SetEnabledMemoryBP(id int, enabled bool) {
	d.memory.setEnabledBP(id, enabled)
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
