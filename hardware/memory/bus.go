package memory

import (
	"github.com/f1gopher/gbpixellib/log"
)

type Area int

const (
	BIOSROM Area = iota
	GPUMemory
	ConsoleRAM
	CartridgeRAMBank
	CartridgeRAM
	CartridgeROM
	CartridgeROMBank
)

func (a Area) String() string {
	return [...]string{
		"BIOS ROM",
		"GPU Memory",
		"Console RAM",
		"Cartridge RAM Bank",
		"Cartridge RAM",
		"Cartridge ROM",
		"Cartridge ROM Bank"}[a]
}

type timerDivide interface {
	TimerDivideWrite()
}

type Bus struct {
	log *log.Log

	bios      *Memory
	video     *videoRam
	ram       *ram
	cartridge Cartridge
	timer     timerDivide

	dmaPending bool
	dmaAddress uint16
}

const ramSize = 0x4000
const ramOffset = 0xC000

func CreateBus(log *log.Log) *Bus {
	return &Bus{
		log:       log,
		bios:      nil,
		video:     CreateVideoRam(),
		ram:       CreateRam(),
		cartridge: nil,
		timer:     nil,
	}
}

func (b *Bus) SetTimer(timer timerDivide) {
	b.timer = timer
}

func (b *Bus) Load(bios *[]byte, cartridge Cartridge) {
	if bios != nil {
		b.bios = CreateReadOnlyMemory("bios", bios, 0)
	} else {
		b.bios = nil
	}
	b.cartridge = cartridge
	//b.ram.Reset()

	// TODO - need to reset system when loading?
}

func (b *Bus) Reset() {
	b.dmaPending = false
	b.dmaAddress = 0x0000
	b.video.Reset()
	b.ram.Reset()
	if b.cartridge != nil {
		b.cartridge.Reset()
	}
}

func (b *Bus) ReadBit(address uint16, bit uint8) bool {
	return b.target(address).ReadBit(address, bit)
}

func (b *Bus) ReadByte(address uint16) byte {
	return b.target(address).ReadByte(address)
}

func (b *Bus) ReadShort(address uint16) uint16 {
	return b.target(address).ReadShort(address)
}
func (b *Bus) WriteBit(address uint16, bit uint8, value bool) {
	b.target(address).WriteBit(address, bit, value)
}
func (b *Bus) WriteByte(address uint16, value byte) {
	// Trigger DMA transfer
	if address == 0xFF46 {
		b.dmaPending = true
		b.dmaAddress = uint16(value) << 8
		return
	}

	if address == DividerRegister {
		b.timer.TimerDivideWrite()
		return
	}

	b.target(address).WriteByte(address, value)
}
func (b *Bus) WriteShort(address uint16, value uint16) {
	b.target(address).WriteShort(address, value)
}

func (b *Bus) SetIO(io inputOutput, interupt interupt) {
	b.ram.SetIO(io, interupt)
}

func (b *Bus) WriteDividerRegister(value uint8) {
	b.ram.WriteDividerRegister(value)
}

func (b *Bus) ExecuteDMAIfPending() bool {
	if !b.dmaPending {
		return false
	}

	var i uint16 = 0
	for i = 0; i < 0x9F; i++ {
		b.WriteByte(0xFE00+i, b.ReadByte(b.dmaAddress+i))
	}

	b.dmaPending = false

	return true
}

func (b *Bus) DisplaySetScanline(value uint8) {
	b.ram.DisplaySetScanline(value)
}

func (b *Bus) DisplaySetStatus(value uint8) {
	b.ram.DisplaySetStatus(value)
}

func (b *Bus) DumpCode(area Area, bank uint8) (data []uint8, startAddress uint16) {
	switch area {
	case BIOSROM:
		return b.bios.DumpCode()
	case GPUMemory:
		return b.video.mem.DumpCode()
	case ConsoleRAM:
		return b.ram.mem.DumpCode()
	case CartridgeRAMBank:
		return b.cartridge.DumpRAMBankCode(bank)
	case CartridgeRAM:
		return b.cartridge.DumpRAMCode()
	case CartridgeROMBank:
		return b.cartridge.DumpROMBankCode(bank)
	case CartridgeROM:
		return b.cartridge.DumpROMCode()
	default:
		panic("Unhandled memory area for dumping")
	}
}

func (b *Bus) target(address uint16) RWMemory {
	if address <= 0xBFFF {

		// Video ram
		if address >= 0x8000 && address <= 0x9FFF {
			return b.video
		}

		if address <= 0x00FF && b.ram.ReadByte(0xFF50) == 0x00 {
			return b.bios
		}

		return b.cartridge
	}

	return b.ram
}
