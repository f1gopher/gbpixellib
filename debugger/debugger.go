package debugger

import (
	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
)

type MemoryRecordEntry struct {
	Value  uint8
	MCycle uint
	PC     uint16
}

type Debugger interface {
	StartCycle(cycle uint, pc uint16)
	HasHitBreakpoint() bool
	BreakpointReason() string

	AddRegisterValueBP(reg cpu.Register, comparison BreakpointComparison, value uint16) int
	DeleteRegisterBP(id int)
	SetEnabledRegisterBP(id int, enabled bool)
	AddMemoryBP(address uint16, comparison BreakpointComparison, value uint8) int
	DeleteMemoryBP(id int)
	SetEnabledMemoryBP(id int, enabled bool)

	DisableAllBreakpoints()

	AddMemoryRecorder(address uint16)
	DeleteMemoryRecorder(address uint16)
	MemoryRecordValues(address uint16) []MemoryRecordEntry
}

func CreateDebugger(l *log.Log, debug bool) (Debugger, cpu.RegistersInterface, cpu.MemoryInterface, *memory.Bus) {
	if debug {
		return createRealDebugger(l)
	}

	return createFakeDebugger(l)
}
