package gameboy

import (
	"github.com/f1gopher/gbpixellib/hardware/z80"
)

type Debug interface {
	AddRegisterValueBP(
		reg z80.Register,
		ccomparison BreakpointComparison,
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

	BreakpointReason() string

	AddMemoryRecorder(address uint16)
	DeleteMemoryRecorder(address uint16)
	MemoryRecordValues(address uint16) []MemoryRecordEntry
}
