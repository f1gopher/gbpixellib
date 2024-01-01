package system

import (
	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/debugger"
)

type Debug interface {
	AddRegisterValueBP(reg cpu.Register, comparison debugger.BreakpointComparison, value uint16)
	AddMemoryBP(address uint16, comparison debugger.BreakpointComparison, value uint8)

	DisableAllBreakpoints()

	BreakpointReason() string
}
