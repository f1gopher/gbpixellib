package gameboy

import (
	"errors"
	"fmt"
	"github.com/f1gopher/gbpixellib/system/gameboy/cartridge"
	"github.com/f1gopher/gbpixellib/system/gameboy/display"
	"github.com/f1gopher/gbpixellib/system/gameboy/input"
	"github.com/f1gopher/gbpixellib/system/gameboy/interupt"
	"os"
	"sync"

	"github.com/f1gopher/gbpixellib/hardware/z80"
	"github.com/f1gopher/gbpixellib/log"
)

const cyclesPerSecond = 4194304
const framesPerSecond = 60
const cyclesPerFrame = cyclesPerSecond / framesPerSecond
const cyclesPerMCycle = 4
const mCyclesPerFrame = cyclesPerFrame / cyclesPerMCycle
const dmaMCycles = 160
const handleInterruptMCycles = 5

const haltExecutionName = "**HALTED**"
const interruptExecutionName = "**INTERRUPT** - "

type CartridgeState struct {
	CurrentROMBank uint8
	CurrentRAMBank uint8
}

type Gameboy struct {
	bios      string
	rom       string
	isTestROM bool

	debugger        Debugger
	log             *log.Log
	screen          *display.Screen
	memory          z80.MemoryInterface
	bus             *Bus
	regs            z80.RegistersInterface
	cpu             *z80.Cpu
	interuptHandler *interupt.Handler
	controller      *input.Input
	timer           *Timer
	cartridgeHeader *CartridgeHeader
	cartridge       cartridge.Cartridge

	currentDisplay string
	displayLock    sync.Mutex

	dump dumpInterface
}

func CreateGameboy(bios string, rom string, useDebugger bool) *Gameboy {
	l := log.CreateLog("./log.txt")
	debugger, registers, memory, memoryBus := CreateDebugger(l, useDebugger)
	system := Gameboy{
		debugger:  debugger,
		log:       l,
		isTestROM: false,
		bios:      bios,
		rom:       rom,
		memory:    memory,
		bus:       memoryBus,
		regs:      registers,
	}
	system.cpu = z80.CreateCPU(l, system.regs, system.memory)
	system.interuptHandler = interupt.CreateHandler(system.memory, system.regs)
	system.screen = display.CreateScreen(system.memory, system.interuptHandler)
	system.controller = input.CreateInput(system.bus, system.interuptHandler)
	memoryBus.SetIO(system.controller, system.interuptHandler)
	system.timer = CreateTimer(system.memory, system.interuptHandler)
	memoryBus.SetTimer(system.timer)

	system.dump = dumpInterface{
		regs:             system.regs,
		cpu:              system.cpu,
		memory:           system.bus,
		screen:           system.screen,
		cartridge:        system.cartridge,
		executionHistory: make([]ExecutionInfo, 0),
	}

	system.Reset()

	return &system
}

func (s *Gameboy) LoadGame(bios string, rom string) {
	s.isTestROM = false
	s.bios = bios
	s.rom = rom
	s.Reset()
}

func (s *Gameboy) LoadTestROM(rom string) {
	s.isTestROM = true
	s.bios = ""
	s.rom = rom
	s.Reset()

	// Disable bios because we load as a ROM
	s.memory.WriteByte(0xFF50, 0xFF)
}

func (s *Gameboy) IsCartridgeSupported() bool {
	return cartridge.IsCartridgeSupported(s.cartridgeHeader.CartridgeType)
}

func (s *Gameboy) Start() {
	var rom []byte
	var bios []byte
	var err error

	if !s.isTestROM {
		bios, err = os.ReadFile(s.bios)
		if err != nil {
			panic(errors.Join(err, errors.New("Failed to load bios")))
		}
	}

	rom, err = os.ReadFile(s.rom)
	if err != nil {
		panic(errors.Join(err, errors.New("Failed to load ROM")))
	}
	s.cartridgeHeader = readHeader(&rom)

	if !s.IsCartridgeSupported() {
		s.log.Debug(fmt.Sprintf("Unsupported cartridge type: %s", s.cartridgeHeader.CartridgeType))
	}

	s.cartridge = cartridge.CreateCartridge(
		s.cartridgeHeader.CartridgeType,
		s.cartridgeHeader.ROMSizeBytes,
		s.cartridgeHeader.RAMSizeBytes,
		&rom)
	s.dump.cartridge = s.cartridge

	s.bus.Load(&bios, s.cartridge)
}

func (s *Gameboy) Reset() {
	s.memory.Reset()
	s.cpu.Reset()
	s.interuptHandler.Reset()
	s.screen.Reset()
	if s.isTestROM {
		s.cpu.InitForTestROM()
		// Disable bios because we load as a ROM
		s.memory.WriteByte(0xFF50, 0xFF)
	} else {
		s.cpu.Init()
	}
	s.controller.Reset()
	s.timer.Reset()
	s.dump.reset()
	s.debugger.StartCycle(0, 0)
	s.Start()
}

func (s *Gameboy) Render(callback func(x int, y int, color display.ScreenColor)) {
	s.displayLock.Lock()
	defer s.displayLock.Unlock()

	s.screen.Render(callback)
}

func (s *Gameboy) SingleFrame() (breakpoint bool, mCyclesCompleted uint, err error) {

	prevCompleted := false
	didDMA := false
	mCyclesCompleted = 0
	wasHalted := false
	var x uint
	info := ExecutionInfo{}

	s.debugger.StartCycle(s.dump.mCycle, info.ProgramCounter)

	for x = 0; x < mCyclesPerFrame; {
		mCyclesCompleted = 1
		info.Opcode = 0
		info.StartMCycle = s.dump.mCycle
		info.ProgramCounter = s.cpu.GetOpcodePC()
		info.StartCPU = *s.dump.getCPUStateOnly()

		if prevCompleted {
			s.debugger.StartCycle(s.dump.mCycle, info.ProgramCounter)

			if didDMA = s.memory.ExecuteDMAIfPending(); didDMA {
				info.Name = "**DMA**"
				mCyclesCompleted += dmaMCycles

				if s.debugger.HasHitBreakpoint() {
					return true, x, nil
				}

			} else {
				wasHalted = s.regs.GetHALT()
				if wasHalted {
					if s.interuptHandler.HasInterrupt() {
						s.regs.SetHALT(false)
						info.Name = "**UNHALT**"
					} else {
						info.Name = "**HALTED**"
					}
					mCyclesCompleted++
					s.dump.appendExecutionHistory(&info)
				} else {
					// If handled an interrupt don't process any instructions this cycle
					if interupted, name := s.interuptHandler.Update(s.cpu.GetOpcodePC()); interupted {
						if err := s.cpu.DoInterruptCycle(); err != nil {
							return false, mCyclesCompleted, err
						}
						mCyclesCompleted = handleInterruptMCycles
						info.Name = interruptExecutionName + name
						s.dump.appendExecutionHistory(&info)
						s.screen.UpdateForCycles(mCyclesCompleted * cyclesPerMCycle)
						x += mCyclesCompleted
						s.dump.mCycle += mCyclesCompleted

						if s.debugger.HasHitBreakpoint() {
							return true, x, nil
						}

						continue
					}
				}
			}

			s.screen.UpdateForCycles(mCyclesCompleted * cyclesPerMCycle)
		}

		if !didDMA && !wasHalted {
			_, prevCompleted, info.Opcode, info.Name, err = s.cpu.ExecuteMCycle()

			if err != nil {
				return false, x, err
			}

			if prevCompleted {
				s.dump.appendExecutionHistory(&info)
			}
		}

		s.timer.Update(uint8(mCyclesCompleted * cyclesPerMCycle))

		x += mCyclesCompleted
		s.dump.mCycle += mCyclesCompleted

		if prevCompleted && s.debugger.HasHitBreakpoint() {
			return true, x, nil
		}
	}

	return false, mCyclesCompleted, nil
}

func (s *Gameboy) SingleInstruction() (breakpoint bool, mCyclesCompleted uint, err error) {

	mCyclesCompleted = 0
	info := ExecutionInfo{
		StartMCycle:    s.dump.mCycle,
		ProgramCounter: s.cpu.GetOpcodePC(),
		StartCPU:       *s.dump.getCPUStateOnly(),
	}
	s.debugger.StartCycle(s.dump.mCycle, info.ProgramCounter)

	if s.memory.ExecuteDMAIfPending() {
		mCyclesCompleted = dmaMCycles
		info.Name = "**DMA**"
	} else {
		if s.regs.GetHALT() {
			if s.interuptHandler.HasInterrupt() {
				s.regs.SetHALT(false)
				info.Name = "**UNHALT**"
			} else {
				info.Name = haltExecutionName
			}
			mCyclesCompleted++

		} else {
			// If handled an interrupt don't process any instructions this cycle
			if interupted, name := s.interuptHandler.Update(s.cpu.GetOpcodePC()); interupted {
				if err := s.cpu.DoInterruptCycle(); err != nil {
					return false, mCyclesCompleted, err
				}

				info.Name = interruptExecutionName + name
				s.dump.appendExecutionHistory(&info)
				mCyclesCompleted = handleInterruptMCycles
				s.screen.UpdateForCycles(mCyclesCompleted * cyclesPerMCycle)

				return s.debugger.HasHitBreakpoint(), mCyclesCompleted, nil
			}

			var completed bool
			for {
				_, completed, info.Opcode, info.Name, err = s.cpu.ExecuteMCycle()

				mCyclesCompleted++

				if err != nil {
					return false, mCyclesCompleted, errors.Join(errors.New("Tick incomplete"), err)
				}

				if completed {
					break
				}
			}
		}
	}

	// Update timers
	s.screen.UpdateForCycles(mCyclesCompleted * cyclesPerMCycle)

	s.timer.Update(uint8(mCyclesCompleted * cyclesPerMCycle))

	s.dump.appendExecutionHistory(&info)

	s.dump.mCycle += mCyclesCompleted

	return s.debugger.HasHitBreakpoint(), mCyclesCompleted, nil
}

func (s *Gameboy) State() string {
	//info := s.cpu.Debug()

	//parts := strings.Split(info, "->")
	//cpu := strings.ReplaceAll(parts[1], " ", "\n")
	//ppu := s.screen.Debug()

	//return parts[0] + "\n" + cpu + "\n" + ppu
	return "<Not implemented>"
}

func (s *Gameboy) OpcodesUsed() {
	// s.cpu.DumpOpcodesUsed()
}

func (s *Gameboy) PreviousPC() uint16 {
	return s.cpu.GetPrevOpcodePC()
}

func (s *Gameboy) CartridgeHeader() CartridgeHeader {
	return *s.cartridgeHeader
}

func (s *Gameboy) Dump() Dump {
	return &s.dump
}

func (s *Gameboy) Debug() Debug {
	return s.debugger
}

func (s *Gameboy) Joypad() Joypad {
	return s.controller
}

func (s *Gameboy) DisplayConfig() display.DisplayConfig {
	return s.screen.DisplayConfig()
}
