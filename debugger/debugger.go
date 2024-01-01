package debugger

import (
	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
)

type Debugger interface {
	StartCycle()
	HasHitBreakpoint() bool
	BreakpointReason() string

	AddRegisterValueBP(reg cpu.Register, comparison BreakpointComparison, value uint16)
	AddMemoryBP(address uint16, comparison BreakpointComparison, value uint8)

	DisableAllBreakpoints()
}

func CreateDebugger(l *log.Log, debug bool) (Debugger, cpu.RegistersInterface, *memory.Bus) {
	if debug {
		return createRealDebugger(l)
	}

	return createFakeDebugger(l)
}
