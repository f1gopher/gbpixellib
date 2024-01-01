package debugger

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type memoryBreakpoint struct {
	enabled    bool
	address    uint16
	value      uint8
	comparison BreakpointComparison
}

type debugMemory struct {
	memory *memory.Bus

	hitBreakpoint bool
	description   string
	breakpoints   map[uint16][]memoryBreakpoint
}

func (d *debugMemory) Reset() {
	d.memory.Reset()
	d.hitBreakpoint = false
	d.description = ""
}

func (d *debugMemory) startCycle() {
	d.hitBreakpoint = false
	d.description = ""
}

func (d *debugMemory) hasHitBreakpoint() bool {
	return d.hitBreakpoint
}

func (d *debugMemory) BreakpointReason() string {
	return d.description
}

func (d *debugMemory) addBP(address uint16, comparison BreakpointComparison, value uint8) {
	bp := memoryBreakpoint{
		enabled:    true,
		address:    address,
		value:      value,
		comparison: comparison,
	}

	if d.breakpoints[address] == nil {
		d.breakpoints[address] = make([]memoryBreakpoint, 0)
	}

	d.breakpoints[address] = append(d.breakpoints[address], bp)
}

func (d *debugMemory) hasBP(address uint16) []memoryBreakpoint {
	bps := d.breakpoints[address]

	if bps == nil {
		return nil
	}

	enabledBps := make([]memoryBreakpoint, 0)

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

func (d *debugMemory) ReadBit(address uint16, bit uint8) bool {
	return d.memory.ReadBit(address, bit)
}

func (d *debugMemory) ReadByte(address uint16) uint8 {
	return d.memory.ReadByte(address)
}

func (d *debugMemory) ReadShort(address uint16) uint16 {
	return d.ReadShort(address)
}

func (d *debugMemory) WriteByte(address uint16, value uint8) {
	bps := d.hasBP(address)

	if bps != nil {
		for _, bp := range bps {
			if evaluateBp(value, bp.comparison, bp.value) {
				d.hitBreakpoint = true
				d.description = fmt.Sprintf("Settings 0x%04X to 0x%02X", address, value)
				continue
			}
		}
	}

	d.memory.WriteByte(address, value)
}

func (d *debugMemory) WriteShort(address uint16, value uint16) {
	d.WriteShort(address, value)
}
