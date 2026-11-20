package m68000

import (
	"fmt"
	"strings"
	"testing"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

// These tests state the MC68000 behavior in the M68000 User's Manual
// (Ninth Edition) for the execution engine in m68000.go: power-on via Init
// before the first cycle (5.5, 6.3.1), the four-clock word read (5.1.1),
// big-endian words (2.4), NOP 4(1/0) (7.9), the group 1 stack frame and
// trace (6.2.4, 6.3.8), and address error (6.3.10).
// The processor keeps a two-word prefetch queue. At an instruction boundary
// IR holds the opcode, IRC holds the next word, and PC addresses the word
// after IRC. A failure means the core does not match that manual. Expected
// values stay on the manual even where the core does not.

const nopOpcode uint16 = 0x4E71

// busLog records word reads and writes. Reset is counted because a CPU
// reset must not erase RAM.
type busLog struct {
	memory.Memory[uint32]
	reads      []uint32
	writeAddrs []uint32
	resetCalls int
}

func (b *busLog) ReadShort(address uint32) uint16 {
	b.reads = append(b.reads, address)
	return b.Memory.ReadShort(address)
}

func (b *busLog) WriteByte(address uint32, value byte) {
	b.writeAddrs = append(b.writeAddrs, address)
	b.Memory.WriteByte(address, value)
}

func (b *busLog) WriteShort(address uint32, value uint16) {
	b.writeAddrs = append(b.writeAddrs, address)
	b.Memory.WriteShort(address, value)
}

func (b *busLog) WriteLong(address uint32, value uint32) {
	b.writeAddrs = append(b.writeAddrs, address)
	b.Memory.WriteShort(address, uint16(value>>16))
	b.Memory.WriteShort(address+2, uint16(value))
}

func (b *busLog) WriteBit(address uint32, bit uint8, value bool) {
	b.writeAddrs = append(b.writeAddrs, address)
	b.Memory.WriteBit(address, bit, value)
}

func (b *busLog) Reset() {
	b.resetCalls++
	b.Memory.Reset()
}

func (b *busLog) clear() {
	b.reads = nil
	b.writeAddrs = nil
	b.resetCalls = 0
}

func newCPU() (*M68000, *busLog, []byte) {
	buf := make([]byte, 0x8000)
	mem := memory.CreateMemoryBigEndian[uint32]("main", &buf, 0)
	bus := &busLog{Memory: mem}
	return CreateM68000(nil, bus), bus, buf
}

// writeWords stores big-endian words at consecutive even addresses.
func writeWords(buf []byte, at uint32, words ...uint16) {
	for i, word := range words {
		address := at + uint32(i)*2
		buf[address] = byte(word >> 8)
		buf[address+1] = byte(word)
	}
}

// primeAtBoundary is the state at the start of an instruction. pc addresses
// the next word to prefetch. opcode is the word at pc-4 and extension is the
// word at pc-2.
func primeAtBoundary(cpu *M68000, pc uint32, opcode uint16, extension uint16) {
	cpu.state = Idle
	cpu.prefetchRequested = false
	cpu.prefetchCyclesRemaining = 0
	cpu.cycle = 0
	regs := cpu.test_registers()
	regs.Set32(PC, pc)
	regs.Set16(IR_Decoding, opcode)
	regs.Set16(IRC_Prefetch, extension)
	regs.Set16(IRD_Executing, 0)
}

func fillRegisters(cpu *M68000, sr uint16) {
	regs := cpu.test_registers()
	for n := 0; n < 8; n++ {
		regs.Set32(register68000(int(D0)+n), 0xD1000000+uint32(n)<<16+0x10+uint32(n))
	}
	for n := 0; n < 7; n++ {
		regs.Set32(register68000(int(A0)+n), 0xA1000000+uint32(n)<<16+0x20+uint32(n))
	}
	regs.Set32(UserStack, 0x0000AAAA)
	regs.Set32(SupervisorStack, 0x0000BBBB)
	regs.Set16(SR, sr)
}

func expect16(t *testing.T, name string, got uint16, want uint16) {
	t.Helper()
	if got != want {
		t.Errorf("%s should be 0x%04X but is 0x%04X", name, want, got)
	}
}

func expectBoundary(t *testing.T, cpu *M68000, pc uint32, ir uint16, irc uint16) {
	t.Helper()
	regs := cpu.test_registers()
	verifyRegister(t, pc, regs.Get32(PC), "PC")
	expect16(t, "IR", regs.Get16(IR_Decoding), ir)
	expect16(t, "IRC", regs.Get16(IRC_Prefetch), irc)
}

func expectWord(t *testing.T, buf []byte, address uint32, want uint16) {
	t.Helper()
	got := uint16(buf[address])<<8 | uint16(buf[address+1])
	if got != want {
		t.Errorf("memory at 0x%06X should be 0x%04X but is 0x%04X", address, want, got)
	}
}

func formatAddrs(addrs []uint32) string {
	if len(addrs) == 0 {
		return "[]"
	}
	parts := make([]string, len(addrs))
	for i, address := range addrs {
		parts[i] = fmt.Sprintf("0x%06X", address)
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func expectReads(t *testing.T, got []uint32, want ...uint32) {
	t.Helper()
	if len(got) == len(want) {
		match := true
		for i := range want {
			if got[i] != want[i] {
				match = false
				break
			}
		}
		if match {
			return
		}
	}
	t.Errorf("word reads should be %s but are %s", formatAddrs(want), formatAddrs(got))
}

func expectNoWrites(t *testing.T, bus *busLog) {
	t.Helper()
	if len(bus.writeAddrs) != 0 {
		t.Errorf("memory writes should be none but are %s", formatAddrs(bus.writeAddrs))
	}
}

func expectRegistersKept(t *testing.T, cpu *M68000, sr uint16, a7 uint32) {
	t.Helper()
	regs := cpu.test_registers()
	for n := 0; n < 8; n++ {
		want := 0xD1000000 + uint32(n)<<16 + 0x10 + uint32(n)
		verifyRegister(t, want, regs.Get32(register68000(int(D0)+n)), fmt.Sprintf("D%d", n))
	}
	for n := 0; n < 7; n++ {
		want := 0xA1000000 + uint32(n)<<16 + 0x20 + uint32(n)
		verifyRegister(t, want, regs.Get32(register68000(int(A0)+n)), fmt.Sprintf("A%d", n))
	}
	verifyRegister(t, 0x0000AAAA, regs.Get32(UserStack), "USP")
	verifyRegister(t, 0x0000BBBB, regs.Get32(SupervisorStack), "SSP")
	verifyRegister(t, a7, regs.Get32(A7), "A7")
	expect16(t, "SR", regs.Get16(SR), sr)
}

// Power-on vector in supervisor program space. Init must not load it.
// Those reads are the first bus cycles after RESET and HALT are released.
const (
	powerOnSSP uint32 = 0x00002000
	powerOnPC  uint32 = 0x00004000
)

// srControlMask selects T, S, and the interrupt mask. The condition codes
// and reserved bits 14, 12, 11, 7, 6, and 5 stay outside the mask: section
// 2.1.3 calls the reserved bits undefined, and section 5.5 does not write
// the condition codes.
const (
	srControlMask uint16 = 0xA700
	srPowerOn     uint16 = 0x2700
)

func writePowerOnMemory(buf []byte) {
	writeWords(buf, 0,
		uint16(powerOnSSP>>16), uint16(powerOnSSP),
		uint16(powerOnPC>>16), uint16(powerOnPC))
	writeWords(buf, 0x0100, 0xA5A5)
}

func expectPowerOn(t *testing.T, cpu *M68000, bus *busLog) {
	t.Helper()
	regs := cpu.test_registers()
	sr := regs.Get16(SR)
	if sr&srControlMask != srPowerOn {
		t.Errorf("SR should have T clear, S set, and interrupt mask 7 (control bits 0x%04X) but SR is 0x%04X", srPowerOn, sr)
	}
	if cpu.state != Idle {
		t.Errorf("power-on is idle, before decode or execute, but state is %d", cpu.state)
	}
	if cpu.cycle != 0 {
		t.Errorf("power-on is before the first cycle but the cycle counter is %d", cpu.cycle)
	}
	if cpu.prefetchRequested || cpu.prefetchCyclesRemaining != 0 || cpu.prefetchAddress != 0 {
		t.Errorf("power-on has no prefetch in progress (requested %v, address 0x%06X, clocks left %d)",
			cpu.prefetchRequested, cpu.prefetchAddress, cpu.prefetchCyclesRemaining)
	}
	if cpu.currentInstruction != nil {
		t.Errorf("power-on has no instruction in progress but the current instruction is %s", cpu.currentInstruction.Name())
	}
	if len(bus.reads) != 0 {
		t.Errorf("power-on reads no memory but reads are %s", formatAddrs(bus.reads))
	}
	expectNoWrites(t, bus)
	if bus.resetCalls != 0 {
		t.Errorf("power-on does not reset memory but Reset was called %d times", bus.resetCalls)
	}
	if regs.Get32(PC) == powerOnPC {
		t.Errorf("PC should still be unloaded; the long word at $000004 is read after RESET and HALT are released")
	}
	if regs.Get32(SupervisorStack) == powerOnSSP {
		t.Errorf("SSP should still be unloaded; the long word at $000000 is read after RESET and HALT are released")
	}
}

func expectPowerOnMemory(t *testing.T, buf []byte) {
	t.Helper()
	expectWord(t, buf, 0, uint16(powerOnSSP>>16))
	expectWord(t, buf, 2, uint16(powerOnSSP))
	expectWord(t, buf, 4, uint16(powerOnPC>>16))
	expectWord(t, buf, 6, uint16(powerOnPC))
	expectWord(t, buf, 0x0100, 0xA5A5)
}

func expectPrefetchAborted(t *testing.T, cpu *M68000, ir uint16, irc uint16, ird uint16) {
	t.Helper()
	regs := cpu.test_registers()
	if got := regs.Get16(IR_Decoding); got == ir {
		t.Errorf("IR should drop the opcode from before power-on but still holds 0x%04X", ir)
	}
	if got := regs.Get16(IRC_Prefetch); got == irc {
		t.Errorf("IRC should drop the word prefetched before power-on but still holds 0x%04X", irc)
	}
	if got := regs.Get16(IRD_Executing); got == ird {
		t.Errorf("IRD should drop the instruction that was executing before power-on but still holds 0x%04X", ird)
	}
}

func expectReadOrder(t *testing.T, got []uint32, want ...uint32) {
	t.Helper()
	next := 0
	for _, address := range got {
		if next < len(want) && address == want[next] {
			next++
		}
	}
	if next != len(want) {
		t.Errorf("word reads should include %s in order but are %s", formatAddrs(want), formatAddrs(got))
	}
}

// TestM68000PowerOn checks Init on a processor that has never executed.
// Section 5.5 holds RESET and HALT low for the initial reset, and Figure
// 5-30 shows no bus cycle until those pins are released. Section 6.3.1
// forces supervisor mode, trace off, and interrupt mask 7 before the vector
// reads. SSP, PC, and the prefetch queue are loaded by later cycles, so
// Init leaves them unloaded. The condition codes, D0-D7, A0-A6, and USP
// have no defined power-on value, and this test does not assign them one.
func TestM68000PowerOn(t *testing.T) {
	cpu, bus, buf := newCPU()
	writePowerOnMemory(buf)
	bus.clear()

	cpu.Init()

	expectPowerOn(t, cpu, bus)
	expectPowerOnMemory(t, buf)
}

// TestM68000PowerOnAbortsExecution checks Init after a program has been
// running. Section 6.3.1 aborts whatever was in progress. The processor
// returns to the idle power-on state, drops the old prefetch, and still
// does not read the reset vector or write memory.
func TestM68000PowerOnAbortsExecution(t *testing.T) {
	cpu, bus, buf := newCPU()
	writePowerOnMemory(buf)
	fillRegisters(cpu, 0x801F)
	regs := cpu.test_registers()
	regs.Set32(PC, 0x00ABCDEF)
	regs.Set16(IR_Decoding, nopOpcode)
	regs.Set16(IRC_Prefetch, 0xABCD)
	regs.Set16(IRD_Executing, 0x1111)
	cpu.state = Execute_Cycle2
	cpu.cycle = 41
	cpu.prefetchRequested = true
	cpu.prefetchAddress = powerOnPC
	cpu.prefetchCyclesRemaining = 3
	cpu.currentInstruction = CreateNOP(nopOpcode, nil, 0)
	bus.clear()

	cpu.Init()

	expectPowerOn(t, cpu, bus)
	expectPowerOnMemory(t, buf)
	expectPrefetchAborted(t, cpu, nopOpcode, 0xABCD, 0x1111)
}

// TestM68000NOP checks one NOP from a primed prefetch queue.
// NOP is 4(1/0): four clocks, one word read, no write. The read is the
// prefetch of the word PC addressed, and that word is big-endian.
func TestM68000NOP(t *testing.T) {
	cpu, bus, buf := newCPU()
	const at uint32 = 0x1000
	writeWords(buf, at, nopOpcode, 0xABCD, 0x1234, 0x5678)
	fillRegisters(cpu, 0x271F)
	primeAtBoundary(cpu, at+4, nopOpcode, 0xABCD)
	bus.clear()

	cpu.Step()

	expectBoundary(t, cpu, at+6, 0xABCD, 0x1234)
	expectReads(t, bus.reads, at+4)
	expectNoWrites(t, bus)
	expectRegistersKept(t, cpu, 0x271F, 0x0000BBBB)
	if cpu.cycle != 4 {
		t.Errorf("NOP is 4(1/0), so the cycle counter should advance by 4 but it is %d", cpu.cycle)
	}
	expectWord(t, buf, at, nopOpcode)
	expectWord(t, buf, at+2, 0xABCD)
	expectWord(t, buf, at+4, 0x1234)
	expectWord(t, buf, at+6, 0x5678)
}

// TestM68000UserModeNOP checks NOP while S is clear. A7 is the user stack
// pointer, and NOP does not write either stack or the condition codes.
func TestM68000UserModeNOP(t *testing.T) {
	cpu, bus, buf := newCPU()
	const at uint32 = 0x1500
	writeWords(buf, at, nopOpcode, 0xABCD, 0x1234, 0x5678)
	fillRegisters(cpu, 0x001F)
	primeAtBoundary(cpu, at+4, nopOpcode, 0xABCD)
	bus.clear()

	cpu.Step()

	expectBoundary(t, cpu, at+6, 0xABCD, 0x1234)
	expectReads(t, bus.reads, at+4)
	expectNoWrites(t, bus)
	expectRegistersKept(t, cpu, 0x001F, 0x0000AAAA)
	if cpu.cycle != 4 {
		t.Errorf("NOP is 4(1/0), so the cycle counter should advance by 4 but it is %d", cpu.cycle)
	}
}

// TestM68000ConsecutiveNOPs checks that each NOP runs as its own instruction.
// The second opcode is the word that was in IRC, including when that prefetch
// was the word 0x0000.
func TestM68000ConsecutiveNOPs(t *testing.T) {
	cpu, bus, buf := newCPU()
	const at uint32 = 0x1200
	writeWords(buf, at, nopOpcode, nopOpcode, 0x0000, 0xBEEF)
	fillRegisters(cpu, 0x271F)
	primeAtBoundary(cpu, at+4, nopOpcode, nopOpcode)

	bus.clear()
	cpu.Step()
	expectBoundary(t, cpu, at+6, nopOpcode, 0x0000)
	expectReads(t, bus.reads, at+4)
	expectNoWrites(t, bus)
	if cpu.cycle != 4 {
		t.Errorf("first NOP should take 4 clocks but the cycle counter is %d", cpu.cycle)
	}

	bus.clear()
	cpu.Step()
	expectBoundary(t, cpu, at+8, 0x0000, 0xBEEF)
	expectReads(t, bus.reads, at+6)
	expectNoWrites(t, bus)
	expectRegistersKept(t, cpu, 0x271F, 0x0000BBBB)
	if cpu.cycle != 8 {
		t.Errorf("two NOP instructions should take 8 clocks but the cycle counter is %d", cpu.cycle)
	}
}

// TestM68000ThreeNOPs checks a run of NOPs whose queued words are all
// non-zero. Each instruction takes four clocks and performs one prefetch.
func TestM68000ThreeNOPs(t *testing.T) {
	cpu, bus, buf := newCPU()
	const at uint32 = 0x1400
	writeWords(buf, at, nopOpcode, nopOpcode, nopOpcode, 0x1111, 0x2222)
	fillRegisters(cpu, 0x271F)
	primeAtBoundary(cpu, at+4, nopOpcode, nopOpcode)

	steps := []struct {
		pc   uint32
		ir   uint16
		irc  uint16
		read uint32
	}{
		{at + 6, nopOpcode, nopOpcode, at + 4},
		{at + 8, nopOpcode, 0x1111, at + 6},
		{at + 10, 0x1111, 0x2222, at + 8},
	}
	for i, step := range steps {
		bus.clear()
		cpu.Step()
		expectBoundary(t, cpu, step.pc, step.ir, step.irc)
		expectReads(t, bus.reads, step.read)
		expectNoWrites(t, bus)
		if cpu.cycle != uint64(i+1)*4 {
			t.Errorf("NOP %d should end on clock %d but the cycle counter is %d", i+1, (i+1)*4, cpu.cycle)
		}
	}
	expectRegistersKept(t, cpu, 0x271F, 0x0000BBBB)
}

// TestM68000ZeroPrefetchWord checks that a prefetched 0x0000 is a real queue
// entry. The opcode 0x0000 is ORI, so an empty queue cannot be represented by
// a zero word. The prefetch started for this NOP reads the word at PC.
func TestM68000ZeroPrefetchWord(t *testing.T) {
	cpu, bus, buf := newCPU()
	const at uint32 = 0x1100
	writeWords(buf, at, nopOpcode, 0x0000, 0x1234, 0x5678)
	fillRegisters(cpu, 0x271F)
	primeAtBoundary(cpu, at+4, nopOpcode, 0x0000)
	bus.clear()

	cpu.Step()

	expectBoundary(t, cpu, at+6, 0x0000, 0x1234)
	expectReads(t, bus.reads, at+4)
	expectNoWrites(t, bus)
	expectRegistersKept(t, cpu, 0x271F, 0x0000BBBB)
	if cpu.cycle != 4 {
		t.Errorf("NOP is 4(1/0), so the cycle counter should advance by 4 but it is %d", cpu.cycle)
	}
}

// TestM68000ReadLatch checks when the prefetched word is sampled.
// NOP's only bus cycle is that prefetch. With no wait states the cycle is
// four clocks, and the data is latched on S6 in the fourth clock. A write to
// that address during the first three clocks is the word that enters IRC.
func TestM68000ReadLatch(t *testing.T) {
	cpu, bus, buf := newCPU()
	const at uint32 = 0x1300
	writeWords(buf, at, nopOpcode, 0x1111, 0xAAAA, 0xCCCC)
	fillRegisters(cpu, 0x271F)
	pc := at + 4
	primeAtBoundary(cpu, pc, nopOpcode, 0x1111)
	bus.clear()

	cpu.ExecuteCycle()
	if cpu.cycle != 1 {
		t.Errorf("ExecuteCycle should be one clock but the cycle counter is %d", cpu.cycle)
	}

	for clock := 0; clock < 2; clock++ {
		cpu.ExecuteCycle()
	}
	if cpu.cycle != 3 {
		t.Errorf("three clocks should have run but the cycle counter is %d", cpu.cycle)
	}
	if cpu.test_registers().Get16(IRC_Prefetch) == 0xAAAA {
		t.Errorf("IRC should wait for S6 in the fourth clock but already holds 0xAAAA")
	}

	writeWords(buf, pc, 0xBBBB)
	cpu.ExecuteCycle()

	expectBoundary(t, cpu, pc+2, 0x1111, 0xBBBB)
	expectReads(t, bus.reads, pc)
	expectNoWrites(t, bus)
	expectRegistersKept(t, cpu, 0x271F, 0x0000BBBB)
	if cpu.cycle != 4 {
		t.Errorf("NOP is 4(1/0), so the cycle counter should advance by 4 but it is %d", cpu.cycle)
	}

	// The latch has closed. A later store to that address stays out of IRC.
	writeWords(buf, pc, 0xDDDD)
	expect16(t, "IRC", cpu.test_registers().Get16(IRC_Prefetch), 0xBBBB)
}
