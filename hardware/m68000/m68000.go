package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/m68000/instructions"
	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
)

type executionState int

const (
	Idle executionState = iota
	Decode_Cycle1
	Execute_Cycle1
	Execute_Cycle2
)

type M68000 struct {
	log *log.Log

	registers instructions.Registers
	memory    memory.Memory[uint32]
	parser    instructionParser

	cycle uint64
	state executionState

	prefetchRequested       bool
	prefetchAddress         uint32
	prefetchValue           uint16
	prefetchCyclesRemaining uint8
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
	// Execute cycle
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

	//instruction.Execute(m.registers, m.memory)

	for x := 0; x < 4; x++ {
		m.ExecuteCycle()
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
				instructions.IRC_Prefetch,
				m.prefetchValue)
			m.prefetchValue = 0
		}
	}

	irc := m.registers.Get16(instructions.IRC_Prefetch)
	ir := m.registers.Get16(instructions.IR_Decoding)

	// If the prefetch buffer is completely empty then request a new instruction
	// and do nothing else
	if irc == 0 && ir == 0 {
		// Fetch next instruction from PC address and store in prefetch queue
	} else if irc != 0 && ir == 0 {
		// If the first prefetch buffer has a value but the second doesn't
		// then move the first value along and request the next instruction
		m.registers.Set16(instructions.IR_Decoding, irc)
		m.registers.Set16(instructions.IRC_Prefetch, 0)
		m.requestPrefetch(m.registers.Get32(instructions.PC))
	} else {
		// We must have an instruction for ready to execute

		// If the first prefetch is empty then request an instruction
		if irc == 0 {
			m.requestPrefetch(m.registers.Get32(instructions.PC) + 2)
		}

		// Start executing the instruction we have
		if m.state == Idle {
			// Move the next instruction out from prefetch queue into executing
			// and clear the prefetch
			m.registers.Set16(instructions.IRD_Executing, ir)
			m.registers.Set16(instructions.IR_Decoding, 0)
			// Now move the first prefetch value along into decoding and request
			// another instruction
			m.registers.Set16(instructions.IR_Decoding, irc)
			m.registers.Set16(instructions.IRC_Prefetch, 0)
			m.requestPrefetch(m.registers.Get32(instructions.PC) + 4)

			// Increment PC by 2 now we have started executing that instruction
			m.registers.IncrementPC(2)

			// TODO - start decoding

			m.state = Decode_Cycle1
		} else if m.state == Decode_Cycle1 {
			// This is the second cycle for decoding
			m.state = Execute_Cycle1
		} else if m.state == Execute_Cycle1 {
			m.state = Execute_Cycle2
		} else if m.state == Execute_Cycle2 {
			opcode := m.registers.Get16(instructions.IRD_Executing)
			instruction := m.parser.decode(
				opcode,
				m.memory.RawBuffer(),
				int(m.registers.Get16(instructions.PC)))

			instruction.Execute(&m.registers, m.memory)

			// Execution finished???? so either idle and wait for the next instruction
			// or start decoding the next intruction
			if m.registers.Get16(instructions.IRC_Prefetch) == 0 {
				m.state = Idle
			} else {
				m.state = Decode_Cycle1
			}
		} else {
			panic(fmt.Sprintf("Unhandled CPU state %v", m.state))
		}
	}

	m.cycle++

	return false
}

func (m *M68000) test_registers() *instructions.Registers {
	return &m.registers
}
