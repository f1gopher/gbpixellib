package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
)

type executionState int

const (
	Idle executionState = iota
	Executing
)

type M68000 struct {
	log *log.Log

	registers Registers
	memory    memory.Memory[uint32]
	parser    instructionParser

	cycle uint64
	state executionState

	prefetchRequested       bool
	prefetchAddress         uint32
	prefetchValue           uint16
	prefetchCyclesRemaining uint8

	currentInstruction instruction
}

func CreateM68000(log *log.Log, memory memory.Memory[uint32]) *M68000 {
	return &M68000{
		log:    log,
		memory: memory,
		// registers: instructions.Registers{},
		state: Idle,
		cycle: 0,
	}
}

func (m *M68000) Init() {
	m.state = Idle
	m.cycle = 0

	m.prefetchRequested = false
	m.prefetchAddress = 0
	m.prefetchValue = 0
	m.prefetchCyclesRemaining = 0

	m.currentInstruction = nil
}

func (m *M68000) Reset() {
	m.memory.Reset()
	m.registers.Reset()

	m.Init()
}

func (m *M68000) requestPrefetch(pc uint32) {
	// If there is already a fetch in progress do nothing
	if m.prefetchRequested {
		return
	}

	// Fetch next instruction from PC address and store in prefetch queue
	m.prefetchRequested = true
	m.prefetchAddress = pc
	m.prefetchCyclesRemaining = 4
	m.prefetchValue = m.memory.ReadShort(pc)
}

func (m *M68000) Step() {

	// Read cycle - fetch value for PC address and store in IR (instruction register)
	// Increment PC
	// IR -> instruction decoder
	// maybe read more bytes for the same instruction so multiple read cycles
	// Step cycle
	// 4 steps/cycles??
	//
	// fetch instruction
	// decode
	// fecth data
	// execute
	//
	// overlap so execute, decode and fetch all happen in the same cycle for different instrucitons
	//
	// databus can trasnfer 2bytes per cycl

	//instruction := m.parser.parse(m.memory.RawBuffer(), int(m.registers.Get16(instructions.PC)))

	//instruction.Step(m.registers, m.memory)

	for !m.ExecuteCycle() {
	}
}

func (m *M68000) ExecuteCycle() (breakpoint bool) {
	// If a prefetch has been requested it takes 4 cycles so update
	if m.prefetchRequested {
		m.prefetchCyclesRemaining--

		// If 4 cycles are up then place the value in IRC
		if m.prefetchCyclesRemaining == 1 {
			m.prefetchRequested = false
			m.prefetchAddress = 0
			m.registers.Set16(
				IRC_Prefetch,
				m.prefetchValue)
			m.prefetchValue = 0
		}
	}

	instructionComplete := false
	irc := m.registers.Get16(IRC_Prefetch)
	ir := m.registers.Get16(IR_Decoding)

	// If the prefetch buffer is completely empty then request a new instruction
	// and do nothing else
	if irc == 0 && ir == 0 {
		// Fetch next instruction from PC address and store in prefetch queue
	} else if irc != 0 && ir == 0 {
		// If the first prefetch buffer has a value but the second doesn't
		// then move the first value along and request the next instruction
		m.registers.Set16(IR_Decoding, irc)
		m.registers.Set16(IRC_Prefetch, 0)
		m.requestPrefetch(m.registers.Get32(PC))
	} else {
		// We must have an instruction for ready to execute

		// If the first prefetch is empty then request an instruction
		if irc == 0 {
			m.requestPrefetch(m.registers.Get32(PC) + 2)
		}

		// Start executing the instruction we have
		if m.state == Idle {
			// Move the next instruction out from prefetch queue into executing
			// and clear the prefetch
			m.registers.Set16(IRD_Executing, ir)
			m.registers.Set16(IR_Decoding, 0)
			// Now move the first prefetch value along into decoding and request
			// another instruction
			m.registers.Set16(IR_Decoding, irc)
			m.registers.Set16(IRC_Prefetch, 0)
			// PC+4 because we already have the next 4 bytes read into the prefetch
			// queue
			m.requestPrefetch(m.registers.Get32(PC) + 4)

			// Increment PC by 2 now we have started executing that instruction
			m.registers.IncrementPC(2)

			// Work out what instruction the next opcode is for and do the first
			// cycle of execution
			opcode := m.registers.Get16(IRD_Executing)
			m.currentInstruction = m.parser.decode2(opcode)

			// Start decoding and executing the instruction
			m.currentInstruction.Step(&m.registers, m.memory)

			m.state = Executing
		} else if m.state == Executing {
			// Execute the next cycle of the instruction
			if m.currentInstruction.Step(&m.registers, m.memory) {
				// Instruction completed
				instructionComplete = true
				m.state = Idle
			}
		} else {
			panic(fmt.Sprintf("Unhandled CPU state %v", m.state))
		}
	}

	m.cycle++

	return instructionComplete
}

func (m *M68000) test_registers() *Registers {
	return &m.registers
}
