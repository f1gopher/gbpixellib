package system

import (
	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/debugger"
)

type Debug interface {
	AddRegisterValueBP(reg cpu.Register, comparison debugger.BreakpointComparison, value uint16) int
	DeleteRegisterBP(id int)
	SetEnabledRegisterBP(id int, enabled bool)
	AddMemoryBP(address uint16, comparison debugger.BreakpointComparison, value uint8) int
	DeleteMemoryBP(id int)
	SetEnabledMemoryBP(id int, enabled bool)

	DisableAllBreakpoints()

	BreakpointReason() string

	AddMemoryRecorder(address uint16)
	DeleteMemoryRecorder(address uint16)
	MemoryRecordValues(address uint16) []debugger.MemoryRecordEntry
}
