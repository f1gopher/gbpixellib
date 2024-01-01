package debugger

import (
	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
)

type fakeDebugger struct {
	log *log.Log
}

func createFakeDebugger(log *log.Log) (Debugger, cpu.RegistersInterface, *memory.Bus) {
	mem := memory.CreateBus(log)
	return &fakeDebugger{}, &cpu.Registers{}, mem
}

func (d *fakeDebugger) StartCycle() {
}

func (d *fakeDebugger) HasHitBreakpoint() bool {
	return false
}

func (d *fakeDebugger) BreakpointReason() string {
	panic("Not Supported")
}

func (d *fakeDebugger) DisableAllBreakpoints() {
	panic("Not supported")
}

func (d *fakeDebugger) AddRegisterValueBP(reg cpu.Register, comparison BreakpointComparison, value uint16) {
	panic("Not supported")
}

func (d *fakeDebugger) AddMemoryBP(address uint16, comparison BreakpointComparison, value uint8) {
	panic("Not supported")
}
