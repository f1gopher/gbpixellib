package system

import (
	"errors"
	"fmt"
	"image"
	"os"
	"sync"

	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/debugger"
	"github.com/f1gopher/gbpixellib/display"
	"github.com/f1gopher/gbpixellib/input"
	"github.com/f1gopher/gbpixellib/interupt"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
	"github.com/f1gopher/gbpixellib/timer"
)

type ExecutionInfo struct {
	Name           string
	StartCycle     uint
	ProgramCounter uint16
}

type DebugState struct {
	NextInstruction     string
	ValueReferencedByPC uint8
}

type CPUState struct {
	A  uint8
	F  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	H  uint8
	L  uint8
	SP uint16
	PC uint16

	ZFlag bool
	NFlag bool
	HFlag bool
	CFlag bool

	SPMem uint16

	Cycle uint
}

type InterruptState struct {
	IME             bool
	VBlankRequested bool
	VBlankEnabled   bool
	LCDRequested    bool
	LCDEnabled      bool
	TimeRequested   bool
	TimeEnabled     bool
	SerialRequested bool
	SerialEnabled   bool
	JoypadRequested bool
	JoypadEnabled   bool
}

type LCDControlState struct {
	LCDEnabled        bool
	WindowTileMapArea uint16
	WindowEnabled     bool
	BGWindowTileData  uint16
	BGTileMap         uint16
	OBJSize           byte
	OBJEnabled        bool
	BGWindowEnabled   bool

	LCDY byte

	LYLYCCompare byte

	LCDStatus_LYCLYInterrupt       bool
	LCDStatus_Mode2OAMInterrupt    bool
	LCDStatus_Mode1VBlankInterrupt bool
	LCDStatus_Mode0HBlankInterrupt bool
	LCDStatus_LYCLY                bool
	LCDStatus_Mode                 byte
	LCDStatus                      uint8

	SCY byte
	SCX byte
	WY  byte
	WX  byte

	BGP_Idx3 byte
	BGP_Idx2 byte
	BGP_Idx1 byte
	BGP_Idx0 byte

	OBP0 byte
	OBP1 byte

	Clock uint
}

type CartridgeState struct {
	CurrentBank uint8
}

type System struct {
	bios      string
	rom       string
	isTestROM bool

	debugger        *debugger.Debugger
	log             *log.Log
	screen          *display.Screen
	memory          *memory.Bus
	regs            *cpu.Registers
	cpu             *cpu.Cpu
	interuptHandler *interupt.Handler
	controller      *input.Input
	timer           *timer.Timer
	cartridgeHeader *CartridgeHeader
	cartridge       memory.Cartridge

	currentDisplay string
	displayLock    sync.Mutex

	pcBreakpoint uint16

	cycle            uint
	executionHistory []ExecutionInfo
}

func CreateSystem(bios string, rom string) *System {
	l := log.CreateLog("./log.txt")
	debugger, registers, memory := debugger.CreateDebugger(l)
	system := System{
		debugger:         debugger,
		log:              l,
		isTestROM:        false,
		bios:             bios,
		rom:              rom,
		memory:           memory,
		regs:             registers,
		pcBreakpoint:     0x0000,
		executionHistory: make([]ExecutionInfo, 10),
	}
	system.cpu = cpu.CreateCPU(l, system.regs, system.memory)
	system.interuptHandler = interupt.CreateHandler(system.memory, system.regs)
	system.screen = display.CreateScreen(system.memory, system.interuptHandler)
	system.controller = input.CreateInput(system.memory, system.interuptHandler)
	system.memory.SetIO(system.controller, system.interuptHandler)
	system.timer = timer.CreateTimer(system.memory, system.interuptHandler)

	system.Reset()

	return &system
}

func (s *System) LoadGame(bios string, rom string) {
	s.isTestROM = false
	s.bios = bios
	s.rom = rom
	s.Reset()
}

func (s *System) LoadTestROM(rom string) {
	s.isTestROM = true
	s.bios = ""
	s.rom = rom
	s.Reset()

	// Disable bios because we load as a ROM
	s.memory.WriteByte(0xFF50, 0xFF)
}

func (s *System) IsCartridgeSupported() bool {
	return s.cartridgeHeader.CartridgeType == 0x00 ||
		s.cartridgeHeader.CartridgeType == 0x01
}

func (s *System) Start() {
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

	s.cartridge = memory.CreateCartridge(
		s.cartridgeHeader.CartridgeType,
		s.cartridgeHeader.ROMSizeBytes,
		s.cartridgeHeader.RAMSizeBytes,
		&rom)

	s.memory.Load(&bios, s.cartridge)
}

func (s *System) Reset() {
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
	s.executionHistory = make([]ExecutionInfo, 10)
	s.cycle = 0
	s.Start()
}

func (s *System) Render(callback func(x int, y int, color display.ScreenColor)) {
	s.displayLock.Lock()
	defer s.displayLock.Unlock()

	s.screen.Render(callback)
}

func (s *System) Tick() (breakpoint bool, cyclesCompleted uint, err error) {

	const maxCycles = 69905
	prevCompleted := false
	didDMA := false
	cyclesCompleted = 0
	wasHalted := false
	var x uint
	info := ExecutionInfo{}

	for x = 0; x < maxCycles; {
		cyclesCompleted = 1
		info.StartCycle = s.cycle
		info.ProgramCounter = s.cpu.GetOpcodePC()

		if prevCompleted {
			if didDMA = s.memory.ExecuteDMAIfPending(); didDMA {
				info.Name = "DMA"
				cyclesCompleted += 162
			} else {
				wasHalted = s.regs.GetHALT()
				if s.regs.GetHALT() {
					if s.interuptHandler.HasInterrupt() {
						s.regs.SetHALT(false)
					}
				} else {
					// If handled an interrupt don't process any instructions this cycle
					if s.interuptHandler.Update() {
						if err := s.cpu.DoInterruptCycle(); err != nil {
							return false, cyclesCompleted, err
						}
						cyclesCompleted = 5
						info.Name = "Interupt"
						s.appendExecutionHistory(&info)
						s.screen.UpdateForCycles(cyclesCompleted * 4)
						prevCompleted = false
						x += cyclesCompleted
						s.cycle += cyclesCompleted
						continue
					}
				}
			}

			s.screen.UpdateForCycles(cyclesCompleted * 4)
			prevCompleted = false
		}

		if !didDMA && !wasHalted {
			breakpoint, prevCompleted, err = s.cpu.ExecuteCycle()

			if err != nil {
				return false, x, err
			}

			if breakpoint {
				return true, x, nil
			}
		}

		s.timer.Update(uint8(x))

		info.Name = fmt.Sprintf("0x%02X", s.cpu.GetPrevOpcode())
		s.appendExecutionHistory(&info)

		x += cyclesCompleted
		s.cycle += cyclesCompleted
	}

	return false, maxCycles, nil
}

func (s *System) SingleInstruction() (cyclesCompleted uint, err error) {

	cyclesCompleted = 0
	info := ExecutionInfo{
		StartCycle:     s.cycle,
		ProgramCounter: s.cpu.GetOpcodePC(),
	}

	if s.memory.ExecuteDMAIfPending() {
		cyclesCompleted = 162
		info.Name = "DMA"
	} else {
		if s.regs.GetHALT() {
			if s.interuptHandler.HasInterrupt() {
				s.regs.SetHALT(false)
			}
			cyclesCompleted++

		} else {
			// If handled an interrupt don't process any instructions this cycle
			if s.interuptHandler.Update() {
				if err := s.cpu.DoInterruptCycle(); err != nil {
					return cyclesCompleted, err
				}

				info.Name = "Interupt"
				s.appendExecutionHistory(&info)
				cyclesCompleted = 5
				s.screen.UpdateForCycles(cyclesCompleted * 4)

				return cyclesCompleted, nil
			}

			for {
				_, completed, err := s.cpu.ExecuteCycle()

				cyclesCompleted++

				if err != nil {
					return cyclesCompleted, errors.Join(errors.New("Tick incomplete"), err)
				}

				if completed {
					break
				}
			}
		}
	}

	// Update timers
	s.screen.UpdateForCycles(cyclesCompleted * 4)

	s.timer.Update(uint8(cyclesCompleted))

	info.Name = fmt.Sprintf("0x%02X", s.cpu.GetPrevOpcode())
	s.appendExecutionHistory(&info)

	s.cycle += cyclesCompleted

	return cyclesCompleted, nil
}

func (s *System) State() string {
	//info := s.cpu.Debug()

	//parts := strings.Split(info, "->")
	//cpu := strings.ReplaceAll(parts[1], " ", "\n")
	//ppu := s.screen.Debug()

	//return parts[0] + "\n" + cpu + "\n" + ppu
	return "<Not implemented>"
}

func (s *System) OpcodesUsed() {
	// s.cpu.DumpOpcodesUsed()
}

func (s *System) PreviousPC() uint16 {
	return s.cpu.GetPrevOpcodePC()
}

func (s *System) GetCPUState() (state *CPUState, prevOpcode uint8, isCB bool) {
	_, isCB = s.cpu.GetNextOpcode()

	lsb := s.memory.ReadByte(s.regs.Get16(cpu.SP))
	msb := s.memory.ReadByte(s.regs.Get16(cpu.SP) + 1)
	spMem := uint16(msb)
	spMem = spMem << 8
	spMem = spMem | uint16(lsb)

	return &CPUState{
		A:     s.regs.Get8(cpu.A),
		F:     s.regs.Get8(cpu.F),
		B:     s.regs.Get8(cpu.B),
		C:     s.regs.Get8(cpu.C),
		D:     s.regs.Get8(cpu.D),
		E:     s.regs.Get8(cpu.E),
		H:     s.regs.Get8(cpu.H),
		L:     s.regs.Get8(cpu.L),
		SP:    s.regs.Get16(cpu.SP),
		PC:    s.regs.Get16(cpu.PC),
		ZFlag: s.regs.GetFlag(cpu.ZFlag),
		NFlag: s.regs.GetFlag(cpu.NFlag),
		HFlag: s.regs.GetFlag(cpu.HFlag),
		CFlag: s.regs.GetFlag(cpu.CFlag),
		SPMem: spMem,
		Cycle: s.cycle,
	}, s.cpu.GetPrevOpcode(), isCB
}

func (s *System) GetInterruptState() *InterruptState {
	requested := s.memory.ReadByte(0xFF0F)
	enabled := s.memory.ReadByte(0xFFFF)

	return &InterruptState{
		IME:             s.regs.GetIME(),
		VBlankRequested: memory.GetBit(requested, 0),
		VBlankEnabled:   memory.GetBit(enabled, 0),
		LCDRequested:    memory.GetBit(requested, 1),
		LCDEnabled:      memory.GetBit(enabled, 1),
		TimeRequested:   memory.GetBit(requested, 2),
		TimeEnabled:     memory.GetBit(enabled, 2),
		SerialRequested: memory.GetBit(requested, 3),
		SerialEnabled:   memory.GetBit(enabled, 3),
		JoypadRequested: memory.GetBit(requested, 4),
		JoypadEnabled:   memory.GetBit(enabled, 4),
	}
}

func (s *System) GetGPUState() *LCDControlState {
	return &LCDControlState{
		LCDEnabled:        s.screen.LCDEnable(),
		WindowTileMapArea: s.screen.WindowTileMapStart(),
		WindowEnabled:     s.screen.WindowEnable(),
		BGWindowTileData:  s.screen.BgWindowTileDataArea(),
		BGTileMap:         s.screen.BgTileMapArea(0), // TODO - is this right
		OBJSize:           s.screen.ObjSize(),
		OBJEnabled:        s.screen.ObjEnable(),
		BGWindowEnabled:   s.screen.BgWindowEnablePriority(),

		LCDY: s.screen.LY(),

		LYLYCCompare: s.screen.LYC(),

		LCDStatus_LYCLYInterrupt:       s.screen.LCDStatusStatInterruptLycLy(),
		LCDStatus_Mode2OAMInterrupt:    s.screen.LCDStatusStatInterruptMode2Oam(),
		LCDStatus_Mode1VBlankInterrupt: s.screen.LCDStatusStatInterruptMode1Vblank(),
		LCDStatus_Mode0HBlankInterrupt: s.screen.LCDStatusStatInterruptMode0Hblank(),
		LCDStatus_LYCLY:                s.screen.LCDStatusLycLy(),
		LCDStatus_Mode:                 byte(s.screen.LCDStatusMode()),
		LCDStatus:                      s.memory.ReadByte(0xFF41),

		SCY: s.screen.SCY(),
		SCX: s.screen.SCX(),
		WY:  s.screen.WY(),
		WX:  s.screen.WX(),

		Clock: s.screen.Cycles(),
	}
}

func (s *System) GetCartridgeState() *CartridgeState {
	return &CartridgeState{
		CurrentBank: s.cartridge.CurrentBank(),
	}
}

func (s *System) GetDebugState() *DebugState {
	return &DebugState{
		NextInstruction:     s.cpu.GetOpcode(),
		ValueReferencedByPC: s.memory.ReadByte(s.regs.Get16(cpu.PC)),
	}
}

func (s *System) GetExecutionHistory() []ExecutionInfo {
	return s.executionHistory
}

func (s *System) DumpTileset() image.Image {
	return s.screen.DumpTileset()
}

func (s *System) DumpTileMap() *[1024]byte {
	return s.screen.DumpBackgroundTileMap()
}

func (s *System) DumpCode() (instructions []string, previousPCIndex int, currentPCIndex int) {
	bios := s.memory.DumpCode()
	current := s.cpu.GetOpcodePC()
	previous := s.cpu.GetPrevOpcodePC()

	//// If we are executing then subtract one because we will inc the PC after getting the opcocde and so will point past
	//// the current instruction
	//if current != 0 {
	//	current--
	//}

	instructions = make([]string, 0)
	currentIndex := 0
	previousIndex := 0

	for x := uint16(0); x < uint16(len(bios)); {
		opcode := bios[x]
		var cbOpcode uint8 = 0
		if x+1 < uint16(len(bios)) {
			cbOpcode = bios[x+1]
		}
		name, opcodeLength := s.cpu.GetOpcodeInfo(opcode, cbOpcode)
		extraInfo := ""

		for y := uint16(1); y < uint16(opcodeLength); y++ {
			extraInfo += fmt.Sprintf(" %02X", bios[x+y])
		}

		if x == s.pcBreakpoint {
			instructions = append(instructions, fmt.Sprintf("0x%04X - %-20s%s **** BREAKPOINT ****\n", x, name, extraInfo))
		} else {
			instructions = append(instructions, fmt.Sprintf("0x%04X - %-20s%s\n", x, name, extraInfo))
		}

		if x == current {
			currentIndex = len(instructions) - 1
		}

		if x == previous {
			previousIndex = len(instructions) - 1
		}

		x += uint16(opcodeLength)
	}

	return instructions, previousIndex, currentIndex
}

func (s *System) SetBreakpoint(pcAddress uint16) {
	s.pcBreakpoint = pcAddress
}

func (s *System) CartridgeHeader() CartridgeHeader {
	return *s.cartridgeHeader
}

func (s *System) PressStart() {
	s.controller.InputStart(true)
}

func (s *System) ReleaseStart() {
	s.controller.InputStart(false)
}

func (s *System) PressSelect() {
	s.controller.InputSelect(true)
}

func (s *System) ReleaseSelect() {
	s.controller.InputSelect(false)
}
func (s *System) PressA() {
	s.controller.InputA(true)
}

func (s *System) ReleaseA() {
	s.controller.InputA(false)
}
func (s *System) PressB() {
	s.controller.InputB(true)
}

func (s *System) ReleaseB() {
	s.controller.InputB(false)
}
func (s *System) PressUp() {
	s.controller.InputUp(true)
}

func (s *System) ReleaseUp() {
	s.controller.InputUp(false)
}
func (s *System) PressDown() {
	s.controller.InputDown(true)
}

func (s *System) ReleaseDown() {
	s.controller.InputDown(false)
}

func (s *System) PressLeft() {
	s.controller.InputLeft(true)
}

func (s *System) ReleaseLeft() {
	s.controller.InputLeft(false)
}
func (s *System) PressRight() {
	s.controller.InputRight(true)
}
func (s *System) ReleaseRight() {
	s.controller.InputRight(false)
}

func (s *System) appendExecutionHistory(action *ExecutionInfo) {
	if len(s.executionHistory) == 10 {
		s.executionHistory = s.executionHistory[1:]
	}

	s.executionHistory = append(s.executionHistory)
}
