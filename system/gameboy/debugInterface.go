package gameboy

import (
	"github.com/f1gopher/gbpixellib/debugger"
	"github.com/f1gopher/gbpixellib/hardware/cpu"
)

type Debug interface {
	AddRegisterValueBP(
		reg cpu.Register,
		ccomparison debugger.BreakpointComparison,
		value uint16,
		hitCount uint) (id int, err error)
	DeleteRegisterBP(id int)
	SetEnabledRegisterBP(id int, enabled bool)
	AddMemoryBP(
		address uint16,
		comparison debugger.BreakpointComparison,
		value uint8,
		hitCount uint) (id int, err error)
	DeleteMemoryBP(id int)
	SetEnabledMemoryBP(id int, enabled bool)

	DisableAllBreakpoints()

	BreakpointReason() string

	AddMemoryRecorder(address uint16)
	DeleteMemoryRecorder(address uint16)
	MemoryRecordValues(address uint16) []debugger.MemoryRecordEntry
}
