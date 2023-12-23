package debugger

import "github.com/f1gopher/gbpixellib/cpu"

type registerBreakpoint struct {
	enabled    bool
	reg        cpu.Register
	value      uint8
	comparison BreakpointComparison
}

type debugRegisters struct {
	registers *cpu.Registers

	hitBreakpoint bool
	breakpoints   map[cpu.Register][]registerBreakpoint
}

func (d *debugRegisters) reset() {
	// d.registers.reset()
}

func (d *debugRegisters) Get8(source cpu.Register) uint8 {
	return d.registers.Get8(source)
}

func (d *debugRegisters) Get16(source cpu.Register) uint16 {
	return d.registers.Get16(source)
}

func (d *debugRegisters) Set8(target cpu.Register, value uint8) {
	bps := d.hasBP(target)

	if bps != nil {
		for x := 0; x < len(bps); x++ {
			if d.evaluateBp(value, bps[x].comparison, bps[x].value) {
				d.hitBreakpoint = true
				continue
			}
		}
	}

	d.registers.Set8(target, value)
}

func (d *debugRegisters) Set16(target cpu.Register, value uint16) {
	d.registers.Set16(target, value)
}

func (d *debugRegisters) AddBP(reg cpu.Register, comparison BreakpointComparison, value uint8) {
	// TODO - validate value against register (8bit or 16bit)

	bp := registerBreakpoint{
		enabled:    true,
		reg:        reg,
		value:      value,
		comparison: comparison,
	}

	if d.breakpoints[reg] == nil {
		d.breakpoints[reg] = make([]registerBreakpoint, 0)
	}

	d.breakpoints[reg] = append(d.breakpoints[reg], bp)
}

func (d *debugRegisters) hasBP(reg cpu.Register) []registerBreakpoint {
	bps := d.breakpoints[reg]

	if bps == nil {
		return nil
	}

	enabledBps := make([]registerBreakpoint, 0)

	for x := 0; x < len(bps); x++ {
		if bps[x].enabled {
			enabledBps = append(enabledBps, bps[x])
		}
	}

	if len(enabledBps) == 0 {
		return nil
	}

	return enabledBps
}

func (d *debugRegisters) evaluateBp(value uint8, comparison BreakpointComparison, bpValue uint8) bool {
	switch comparison {
	case Equals:
		return value == bpValue
	case GreaterThanOrEqual:
		return value >= bpValue
	default:
		panic("Not implemented breakpoint comparison")
	}
}
