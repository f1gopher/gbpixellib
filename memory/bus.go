package memory

import "github.com/f1gopher/gbpixellib/log"

type Bus struct {
	log *log.Log

	bios      *Memory2
	video     *videoRam
	ram       *ram
	cartridge *Cartridge

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
	}
}

func (b *Bus) Load(bios *[]byte, cartridge *Cartridge) {
	if bios != nil {
		b.bios = CreateMemory2("bios", bios, 0)
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
	if address == 0x2000 {
		address = address * 1
	}

	return b.target(address).ReadByte(address)
}

func (b *Bus) ReadShort(address uint16) uint16 {
	return b.target(address).ReadShort(address)
}
func (b *Bus) WriteBit(address uint16, bit uint8, value bool) {
	b.target(address).WriteBit(address, bit, value)
}
func (b *Bus) WriteByte(address uint16, value byte) {
	// Can't write to this address range
	if address <= 0x7FFF {
		return
	}

	// Trigger DMA transfer
	if address == 0xFF46 {
		b.dmaPending = true
		b.dmaAddress = uint16(value) << 8
		return
	}

	b.target(address).WriteByte(address, value)
}
func (b *Bus) WriteShort(address uint16, value uint16) {
	// Can't write to this address range
	if address <= 0x7FFF {
		return
	}

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

func (b *Bus) DumpCode() []uint8 {
	return b.cartridge.DumpCode()
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
