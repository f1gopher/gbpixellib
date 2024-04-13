package debugger

import (
	"github.com/f1gopher/gbpixellib/hardware/z80"
	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
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

	AddRegisterValueBP(
		reg z80.Register,
		comparison BreakpointComparison,
		value uint16,
		hitCount uint) (id int, err error)
	DeleteRegisterBP(id int)
	SetEnabledRegisterBP(id int, enabled bool)
	AddMemoryBP(
		address uint16,
		comparison BreakpointComparison,
		value uint8,
		hitCount uint) (id int, err error)
	DeleteMemoryBP(id int)
	SetEnabledMemoryBP(id int, enabled bool)

	DisableAllBreakpoints()

	AddMemoryRecorder(address uint16)
	DeleteMemoryRecorder(address uint16)
	MemoryRecordValues(address uint16) []MemoryRecordEntry
}

func CreateDebugger(l *log.Log, debug bool) (Debugger, z80.RegistersInterface, z80.MemoryInterface, *memory.Bus) {
	if debug {
		return createRealDebugger(l)
	}

	return createFakeDebugger(l)
}
