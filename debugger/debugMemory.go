package debugger

import (
	"fmt"
	"slices"
	"sync"

	"github.com/f1gopher/gbpixellib/memory"
)

type memoryBreakpoint struct {
	id         int
	enabled    bool
	address    uint16
	value      uint8
	comparison BreakpointComparison
}

type memoryRecord struct {
	address uint16
	history []MemoryRecordEntry
}

type debugMemory struct {
	memory       *memory.Bus
	currentCycle uint
	currentPC    uint16

	nextId        int
	hitBreakpoint bool
	description   string
	breakpoints   map[uint16][]memoryBreakpoint
	bpLock        sync.RWMutex

	records map[uint16]*memoryRecord
}

func (d *debugMemory) Reset() {
	d.currentCycle = 0
	d.memory.Reset()
	d.hitBreakpoint = false
	d.description = ""

	for x := range d.records {
		tmp := d.records[x]
		tmp.history = []MemoryRecordEntry{}
	}
}

func (d *debugMemory) startCycle(cycle uint, pc uint16) {
	d.hitBreakpoint = false
	d.description = ""
	d.currentCycle = cycle
	d.currentPC = pc
}

func (d *debugMemory) hasHitBreakpoint() bool {
	return d.hitBreakpoint
}

func (d *debugMemory) BreakpointReason() string {
	return d.description
}

func (d *debugMemory) addBP(address uint16, comparison BreakpointComparison, value uint8) int {
	bp := memoryBreakpoint{
		id:         d.nextId,
		enabled:    true,
		address:    address,
		value:      value,
		comparison: comparison,
	}
	d.nextId++

	d.bpLock.Lock()
	if d.breakpoints[address] == nil {
		d.breakpoints[address] = make([]memoryBreakpoint, 0)
	}

	d.breakpoints[address] = append(d.breakpoints[address], bp)
	d.bpLock.Unlock()

	return bp.id
}

func (d *debugMemory) deleteBP(id int) {
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

func (d *debugMemory) setEnabledBP(id int, enabled bool) {
	d.bpLock.Lock()
	defer d.bpLock.Unlock()

	for key := range d.breakpoints {
		for x := range d.breakpoints[key] {
			if d.breakpoints[key][x].id == id {
				d.breakpoints[key][id].enabled = enabled
				return
			}
		}
	}
}

func (d *debugMemory) hasBP(address uint16) []memoryBreakpoint {

	d.bpLock.RLock()
	defer d.bpLock.RUnlock()
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

func (d *debugMemory) addRecorder(address uint16) {

	// If a recorder already exists do nothing
	if _, exists := d.records[address]; exists {
		return
	}

	// Store current value as MCycle 0 as the current/starting value
	d.records[address] = &memoryRecord{
		address: address,
		history: []MemoryRecordEntry{
			{
				MCycle: 0,
				Value:  d.memory.ReadByte(address),
				PC:     0,
			},
		},
	}
}

func (d *debugMemory) deleteRecorder(address uint16) {
	delete(d.records, address)
}

func (d *debugMemory) recordValues(address uint16) []MemoryRecordEntry {
	entry, exists := d.records[address]
	if !exists {
		return nil
	}

	return entry.history
}

func (d *debugMemory) ReadBit(address uint16, bit uint8) bool {
	return d.memory.ReadBit(address, bit)
}

func (d *debugMemory) ReadByte(address uint16) uint8 {
	return d.memory.ReadByte(address)
}

func (d *debugMemory) ReadShort(address uint16) uint16 {
	return d.memory.ReadShort(address)
}

func (d *debugMemory) WriteByte(address uint16, value uint8) {
	bps := d.hasBP(address)

	if bps != nil {
		d.bpLock.RLock()
		for _, bp := range bps {
			if evaluateBp(value, bp.comparison, bp.value) {
				d.hitBreakpoint = true
				d.description = fmt.Sprintf("Settings 0x%04X to 0x%02X", address, value)
				continue
			}
		}
		d.bpLock.RUnlock()
	}

	recorder, exists := d.records[address]
	if exists {
		recorder.history = append(recorder.history, MemoryRecordEntry{
			MCycle: d.currentCycle,
			PC:     d.currentPC,
			Value:  value,
		})
	}

	d.memory.WriteByte(address, value)
}

func (d *debugMemory) WriteShort(address uint16, value uint16) {
	d.memory.WriteShort(address, value)
}

func (d *debugMemory) DisplaySetScanline(value uint8) {
	d.memory.DisplaySetScanline(value)
}

func (d *debugMemory) DisplaySetStatus(value uint8) {
	d.memory.DisplaySetStatus(value)
}

func (d *debugMemory) DumpCode(area memory.Area, bank uint8) (data []uint8, startAddress uint16) {
	return d.memory.DumpCode(area, bank)
}

func (d *debugMemory) WriteDividerRegister(value uint8) {
	d.memory.WriteDividerRegister(value)
}

func (d *debugMemory) ExecuteDMAIfPending() bool {
	return d.memory.ExecuteDMAIfPending()
}
