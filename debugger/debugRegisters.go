package debugger

import (
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/f1gopher/gbpixellib/cpu"
)

type registerBreakpoint struct {
	id              int
	enabled         bool
	reg             cpu.Register
	value           uint16
	comparison      BreakpointComparison
	targetHitCount  uint
	currentHitCount uint
}

type debugRegisters struct {
	registers cpu.RegistersInterface

	nextId        int
	hitBreakpoint bool
	description   string
	breakpoints   map[cpu.Register][]registerBreakpoint
	bpLock        sync.RWMutex
}

func (d *debugRegisters) Reset() {
	d.registers.Reset()
	d.hitBreakpoint = false
	d.description = ""

	d.bpLock.Lock()
	for bpsForReg := range d.breakpoints {
		for bpIndex := range d.breakpoints[bpsForReg] {
			d.breakpoints[bpsForReg][bpIndex].currentHitCount = 1
		}
	}
	d.bpLock.Unlock()
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
	bps := d.hasBP(target)

	if bps != nil {
		d.bpLock.Lock()
		for x := 0; x < len(bps); x++ {
			lsbValue := cpu.Lsb(bps[x].value)
			if evaluateBp(value, bps[x].comparison, lsbValue) {
				if bps[x].targetHitCount != bps[x].currentHitCount {
					bps[x].currentHitCount++
					continue
				}

				d.hitBreakpoint = true
				d.description = fmt.Sprintf(
					"Setting %s to 0x%02X for the %d time",
					target.String(),
					lsbValue,
					bps[x].currentHitCount)
				continue
			}
		}
		d.bpLock.Unlock()
	}

	d.registers.Set8(target, value)
}

func (d *debugRegisters) Set16(target cpu.Register, value uint16) {
	bps := d.hasBP(target)

	if bps != nil {
		d.bpLock.Lock()
		for x := 0; x < len(bps); x++ {
			if evaluateBp(value, bps[x].comparison, bps[x].value) {
				if bps[x].targetHitCount != bps[x].currentHitCount {
					bps[x].currentHitCount++
					continue
				}

				d.hitBreakpoint = true
				d.description = fmt.Sprintf(
					"Setting %s to 0x%04X for the %d time",
					target.String(),
					value,
					bps[x].currentHitCount)
				continue
			}
		}
		d.bpLock.Unlock()
	}

	d.registers.Set16(target, value)
}

func (d *debugRegisters) Set16FromTwoBytes(target cpu.Register, msb uint8, lsb uint8) {
	// Do this so we check for breakpoints
	d.Set16(target, cpu.CombineBytes(msb, lsb))
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

func (d *debugRegisters) addBP(
	reg cpu.Register,
	comparison BreakpointComparison,
	value uint16,
	hitCount uint) (id int, err error) {

	if hitCount < 1 {
		return -1, errors.New("hitCount must be >= 1")
	}

	// We accept a 16bit value for the breakpoint but some registers are 8bit.
	// When this happens clear the unused 8bits by setting them to 0
	if reg.Is8Bit() {
		value = value | 0b0000000011111111
	}

	bp := registerBreakpoint{
		id:              d.nextId,
		enabled:         true,
		reg:             reg,
		value:           value,
		comparison:      comparison,
		targetHitCount:  hitCount,
		currentHitCount: 1,
	}
	d.nextId++

	d.bpLock.Lock()
	if d.breakpoints[reg] == nil {
		d.breakpoints[reg] = make([]registerBreakpoint, 0)
	}

	d.breakpoints[reg] = append(d.breakpoints[reg], bp)
	d.bpLock.Unlock()

	return bp.id, nil
}

func (d *debugRegisters) deleteBP(id int) {
	d.bpLock.Lock()
	defer d.bpLock.Unlock()

	for key := range d.breakpoints {
		for x := range d.breakpoints[key] {
			if d.breakpoints[key][x].id == id {
				d.breakpoints[key] = slices.Delete(d.breakpoints[key], x, +1)
				return
			}
		}
	}
}

func (d *debugRegisters) setEnabledBP(id int, enabled bool) {
	d.bpLock.Lock()
	defer d.bpLock.Unlock()

	for key := range d.breakpoints {
		for x := range d.breakpoints[key] {
			if d.breakpoints[key][x].id == id {
				d.breakpoints[key][x].enabled = enabled
				return
			}
		}
	}
}

func (d *debugRegisters) hasBP(reg cpu.Register) []registerBreakpoint {
	d.bpLock.RLock()
	defer d.bpLock.RUnlock()

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

func evaluateBp[T uint8 | uint16](value T, comparison BreakpointComparison, bpValue T) bool {
	switch comparison {
	case Equal:
		return value == bpValue
	case NotEqual:
		return value != bpValue
	case GreaterThan:
		return value > bpValue
	case LessThan:
		return value < bpValue
	case GreaterThanOrEqual:
		return value >= bpValue
	case LessThanOrEqual:
		return value <= bpValue
	default:
		panic("Not implemented breakpoint comparison")
	}
}
