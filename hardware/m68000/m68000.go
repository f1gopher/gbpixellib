package m68000

import (
	"fmt"

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

	registers Registers
	memory    memory.Memory[uint32]
	parser    instructionParser

	cycle uint64
	state executionState

	prefetchRequested       bool
	prefetchAddress         uint32
	prefetchCyclesRemaining uint8

	currentInstruction instruction
}

func CreateM68000(log *log.Log, memory memory.Memory[uint32]) *M68000 {
	return &M68000{
		log:    log,
		memory: memory,
		state:  Idle,
		cycle:  0,
	}
}

func (m *M68000) Init() {
	m.state = Idle
	m.cycle = 0

	m.prefetchRequested = false
	m.prefetchAddress = 0
	m.prefetchCyclesRemaining = 0

	m.currentInstruction = nil

	// Power on initial state
	m.registers.Set16(SR, 0x2700)
	m.registers.Set16(IR_Decoding, 0)
	m.registers.Set16(IRC_Prefetch, 0)
	m.registers.Set16(IRD_Executing, 0)

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
			prefetchValue := m.memory.ReadShort(m.prefetchAddress)
			m.prefetchRequested = false
			m.prefetchAddress = 0
			m.registers.Set16(
				IRC_Prefetch,
				prefetchValue)
		}
	}

	irc := m.registers.Get16(IRC_Prefetch)
	ir := m.registers.Get16(IR_Decoding)

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
		m.requestPrefetch(m.registers.Get32(PC))

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
		opcode := m.registers.Get16(IRD_Executing)
		m.currentInstruction = m.parser.decode2(opcode)

		m.currentInstruction.Step(&m.registers, m.memory)

		m.state = Idle
	} else {
		panic(fmt.Sprintf("Unhandled CPU state %v", m.state))
	}

	m.cycle++

	return false
}

func (m *M68000) test_registers() *Registers {
	return &m.registers
}
