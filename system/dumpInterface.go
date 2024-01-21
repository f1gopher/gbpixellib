package system

import (
	"fmt"
	"image"

	"github.com/f1gopher/gbpixellib/cpu"
	"github.com/f1gopher/gbpixellib/display"
	"github.com/f1gopher/gbpixellib/memory"
)

const executionHistorySize = 10

type ExecutionInfo struct {
	Name           string
	StartMCycle    uint
	ProgramCounter uint16
	StartCPU       CPUState
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

	BGP_Idx3 display.ScreenColor
	BGP_Idx2 display.ScreenColor
	BGP_Idx1 display.ScreenColor
	BGP_Idx0 display.ScreenColor

	OBJ_0_Idx3 display.ScreenColor
	OBJ_0_Idx2 display.ScreenColor
	OBJ_0_Idx1 display.ScreenColor
	OBJ_0_Idx0 display.ScreenColor

	OBJ_1_Idx3 display.ScreenColor
	OBJ_1_Idx2 display.ScreenColor
	OBJ_1_Idx1 display.ScreenColor
	OBJ_1_Idx0 display.ScreenColor

	OBP1 byte

	Clock uint
}
type Dump interface {
	GetCPUState() (state *CPUState, prevOpcode uint8, isCB bool)
	GetInterruptState() *InterruptState
	GetGPUState() *LCDControlState
	GetCartridgeState() *CartridgeState
	GetDebugState() *DebugState
	DumpTileset() image.Image
	DumpTile(tileNum uint16, palette display.Palette) image.Image
	DumpFirstTileMap() *[1024]byte
	DumpSecondTileMap() *[1024]byte
	DumpWindowTileMap() *[1024]byte
	DumpBackgroundTileMap() *[1024]byte
	DumpCode(area memory.Area, bank uint8) (instructions []string, previousPCIndex int, currentPCIndex int)
	DumpCallstack() []string
	GetExecutionHistory() []ExecutionInfo
	DumpMemory(area memory.Area, bank uint8) (data []uint8, startAddress uint16)
	DumpMemoryValue(address uint16) uint8
}

type dumpInterface struct {
	regs             cpu.RegistersInterface
	cpu              *cpu.Cpu
	memory           *memory.Bus
	screen           *display.Screen
	cartridge        memory.Cartridge
	executionHistory []ExecutionInfo
	mCycle           uint
}

func (d *dumpInterface) reset() {
	d.executionHistory = make([]ExecutionInfo, 0)
	d.mCycle = 0
}

func (d *dumpInterface) appendExecutionHistory(action *ExecutionInfo) {
	if len(d.executionHistory) == executionHistorySize {
		d.executionHistory = d.executionHistory[1:]
	}

	d.executionHistory = append(d.executionHistory, *action)
}

func (d *dumpInterface) getCPUStateOnly() *CPUState {
	lsb := d.memory.ReadByte(d.regs.Get16(cpu.SP))
	msb := d.memory.ReadByte(d.regs.Get16(cpu.SP) + 1)
	spMem := uint16(msb)
	spMem = spMem << 8
	spMem = spMem | uint16(lsb)

	return &CPUState{
		A:     d.regs.Get8(cpu.A),
		F:     d.regs.Get8(cpu.F),
		B:     d.regs.Get8(cpu.B),
		C:     d.regs.Get8(cpu.C),
		D:     d.regs.Get8(cpu.D),
		E:     d.regs.Get8(cpu.E),
		H:     d.regs.Get8(cpu.H),
		L:     d.regs.Get8(cpu.L),
		SP:    d.regs.Get16(cpu.SP),
		PC:    d.regs.Get16(cpu.PC),
		ZFlag: d.regs.GetFlag(cpu.ZFlag),
		NFlag: d.regs.GetFlag(cpu.NFlag),
		HFlag: d.regs.GetFlag(cpu.HFlag),
		CFlag: d.regs.GetFlag(cpu.CFlag),
		SPMem: spMem,
		Cycle: d.mCycle,
	}
}

func (d *dumpInterface) GetCPUState() (state *CPUState, prevOpcode uint8, isCB bool) {
	_, isCB = d.cpu.GetNextOpcode()
	return d.getCPUStateOnly(), d.cpu.GetPrevOpcode(), isCB
}

func (d *dumpInterface) GetInterruptState() *InterruptState {
	requested := d.memory.ReadByte(0xFF0F)
	enabled := d.memory.ReadByte(0xFFFF)

	return &InterruptState{
		IME:             d.regs.GetIME(),
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

func (d *dumpInterface) GetGPUState() *LCDControlState {
	return &LCDControlState{
		LCDEnabled:        d.screen.LCDEnable(),
		WindowTileMapArea: d.screen.WindowTileMapStart(),
		WindowEnabled:     d.screen.WindowEnable(),
		BGWindowTileData:  d.screen.BgWindowTileDataArea(),
		BGTileMap:         d.screen.BackgroundTileMapStart(),
		OBJSize:           d.screen.ObjSize(),
		OBJEnabled:        d.screen.ObjEnable(),
		BGWindowEnabled:   d.screen.BgWindowEnablePriority(),

		LCDY: d.screen.LY(),

		LYLYCCompare: d.screen.LYC(),

		LCDStatus_LYCLYInterrupt:       d.screen.LCDStatusStatInterruptLycLy(),
		LCDStatus_Mode2OAMInterrupt:    d.screen.LCDStatusStatInterruptMode2Oam(),
		LCDStatus_Mode1VBlankInterrupt: d.screen.LCDStatusStatInterruptMode1Vblank(),
		LCDStatus_Mode0HBlankInterrupt: d.screen.LCDStatusStatInterruptMode0Hblank(),
		LCDStatus_LYCLY:                d.screen.LCDStatusLycLy(),
		LCDStatus_Mode:                 byte(d.screen.LCDStatusMode()),
		LCDStatus:                      d.memory.ReadByte(0xFF41),

		SCY: d.screen.SCY(),
		SCX: d.screen.SCX(),
		WY:  d.screen.WY(),
		WX:  d.screen.WX(),

		Clock: d.screen.Cycles(),

		BGP_Idx0: d.screen.BGPIndex0Color(),
		BGP_Idx1: d.screen.BGPIndex1Color(),
		BGP_Idx2: d.screen.BGPIndex2Color(),
		BGP_Idx3: d.screen.BGPIndex3Color(),

		OBJ_0_Idx3: d.screen.ObjPalette0Index3Color(),
		OBJ_0_Idx2: d.screen.ObjPalette0Index2Color(),
		OBJ_0_Idx1: d.screen.ObjPalette0Index1Color(),
		OBJ_0_Idx0: d.screen.ObjPalette0Index0Color(),

		OBJ_1_Idx3: d.screen.ObjPalette1Index3Color(),
		OBJ_1_Idx2: d.screen.ObjPalette1Index2Color(),
		OBJ_1_Idx1: d.screen.ObjPalette1Index1Color(),
		OBJ_1_Idx0: d.screen.ObjPalette1Index0Color(),
	}
}

func (d *dumpInterface) GetCartridgeState() *CartridgeState {
	return &CartridgeState{
		CurrentROMBank: d.cartridge.CurrentROMBank(),
		CurrentRAMBank: d.cartridge.CurrentRAMBank(),
	}
}

func (d *dumpInterface) GetDebugState() *DebugState {
	return &DebugState{
		NextInstruction:     d.cpu.GetOpcode(),
		ValueReferencedByPC: d.memory.ReadByte(d.regs.Get16(cpu.PC)),
	}
}

func (d *dumpInterface) GetExecutionHistory() []ExecutionInfo {
	return d.executionHistory
}

func (d *dumpInterface) DumpTileset() image.Image {
	return d.screen.DumpTileset()
}

func (d *dumpInterface) DumpTile(tileNum uint16, palette display.Palette) image.Image {
	return d.screen.DumpTile(tileNum, palette)
}

func (d *dumpInterface) DumpFirstTileMap() *[1024]byte {
	return d.screen.DumpFirstTileMap()
}

func (d *dumpInterface) DumpSecondTileMap() *[1024]byte {
	return d.screen.DumpSecondTileMap()
}

func (d *dumpInterface) DumpWindowTileMap() *[1024]byte {
	return d.screen.DumpWindowTileMap()
}

func (d *dumpInterface) DumpBackgroundTileMap() *[1024]byte {
	return d.screen.DumpBackgroundTileMap()
}

func (d *dumpInterface) DumpCode(area memory.Area, bank uint8) (instructions []string, previousPCIndex int, currentPCIndex int) {
	bios, _ := d.memory.DumpCode(area, bank)
	current := d.cpu.GetOpcodePC()
	previous := d.cpu.GetPrevOpcodePC()

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
		name, opcodeLength := d.cpu.GetOpcodeInfo(opcode, cbOpcode)
		extraInfo := ""

		for y := uint16(1); y < uint16(opcodeLength); y++ {
			extraInfo += fmt.Sprintf(" %02X", bios[x+y])
		}

		instructions = append(instructions, fmt.Sprintf("0x%04X - %-20s%s\n", x, name, extraInfo))

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

func (d *dumpInterface) DumpCallstack() []string {
	var stackStart uint16 = 0xFFFE
	stackEnd := d.regs.Get16(cpu.SP)
	result := make([]string, 0)

	for x := stackEnd; x < stackStart; x += 2 {
		value := d.memory.ReadShort(x)

		result = append(result, fmt.Sprintf("0x%04X => 0x%04X", x, value))
	}
	result = append(result, fmt.Sprintf("0x%04X => 0x%04X", stackStart, d.memory.ReadShort(stackStart)))

	return result
}

func (d *dumpInterface) DumpMemory(area memory.Area, bank uint8) (data []uint8, startAddress uint16) {
	return d.memory.DumpCode(area, bank)
}

func (d *dumpInterface) DumpMemoryValue(address uint16) uint8 {
	return d.memory.ReadByte(address)
}
