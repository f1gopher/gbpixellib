package debugger

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/cpu"
)

type registerBreakpoint struct {
	enabled    bool
	reg        cpu.Register
	value      uint16
	comparison BreakpointComparison
}

type debugRegisters struct {
	registers cpu.RegistersInterface

	hitBreakpoint bool
	description   string
	breakpoints   map[cpu.Register][]registerBreakpoint
}

func (d *debugRegisters) Reset() {
	d.registers.Reset()
	d.hitBreakpoint = false
	d.description = ""
}

func (d *debugRegisters) startCycle() {
	d.hitBreakpoint = false
	d.description = ""
}

func (d *debugRegisters) hasHitBreakpoint() bool {
	return d.hitBreakpoint
}

func (d *debugRegisters) BreakpointReason() string {
	return d.description
}

func (d *debugRegisters) Get8(source cpu.Register) uint8 {
	return d.registers.Get8(source)
}

func (d *debugRegisters) Get16(source cpu.Register) uint16 {
	return d.registers.Get16(source)
}

func (d *debugRegisters) Get16Msb(source cpu.Register) uint8 {
	return d.registers.Get16Msb(source)
}

func (d *debugRegisters) Get16Lsb(source cpu.Register) uint8 {
	return d.registers.Get16Lsb(source)
}

func (d *debugRegisters) Set8(target cpu.Register, value uint8) {
	// bps := d.hasBP(target)
	//
	// if bps != nil {
	// 	for x := 0; x < len(bps); x++ {
	// 		if d.evaluateBp(value, bps[x].comparison, bps[x].value) {
	// 			d.hitBreakpoint = true
	// 			continue
	// 		}
	// 	}
	// }
	//
	d.registers.Set8(target, value)
}

func (d *debugRegisters) Set16(target cpu.Register, value uint16) {
	bps := d.hasBP(target)

	if bps != nil {
		for x := 0; x < len(bps); x++ {
			if d.evaluateBp(value, bps[x].comparison, bps[x].value) {
				d.hitBreakpoint = true
				d.description = fmt.Sprintf("Setting %s to 0x%02X", target.String(), value)
				continue
			}
		}
	}

	d.registers.Set16(target, value)
}

func (d *debugRegisters) Set16FromTwoBytes(target cpu.Register, msb uint8, lsb uint8) {
	d.registers.Set16FromTwoBytes(target, msb, lsb)
}

func (d *debugRegisters) SetRegBit(target cpu.Register, bit uint8, value bool) {
	d.registers.SetRegBit(target, bit, value)
}

func (d *debugRegisters) GetFlag(flag cpu.RegisterFlags) bool {
	return d.registers.GetFlag(flag)
}

func (d *debugRegisters) SetFlag(flag cpu.RegisterFlags, value bool) {
	d.registers.SetFlag(flag, value)
}

func (d *debugRegisters) SetIME(enabled bool) {
	d.registers.SetIME(enabled)
}

func (d *debugRegisters) GetIME() bool {
	return d.registers.GetIME()
}

func (d *debugRegisters) SetHALT(enabled bool) {
	d.registers.SetHALT(enabled)
}

func (d *debugRegisters) GetHALT() bool {
	return d.registers.GetHALT()
}

func (d *debugRegisters) AddBP(reg cpu.Register, comparison BreakpointComparison, value uint16) {
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

func (d *debugRegisters) evaluateBp(value uint16, comparison BreakpointComparison, bpValue uint16) bool {
	switch comparison {
	case Equals:
		return value == bpValue
	case GreaterThanOrEqual:
		return value >= bpValue
	default:
		panic("Not implemented breakpoint comparison")
	}
}
