package cps1

import (
	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/system/cps1/cartridge"
	"github.com/f1gopher/gbpixellib/system/cps1/memoryBus"
)

type jamma struct {
	log         *log.Log
	jammaMemory memory.Memory[uint32]
	kickMemory  memory.Memory[uint32]
}

const dipAAddress = memoryBus.JammaDipAddress + 0x2
const dipBAddress = memoryBus.JammaDipAddress + 0x4
const dipCAddress = memoryBus.JammaDipAddress + 0x6

func createJamma(log *log.Log, jammaMemory memory.Memory[uint32], kickMemory memory.Memory[uint32]) *jamma {
	return &jamma{
		log:         log,
		jammaMemory: jammaMemory,
		kickMemory:  kickMemory,
	}
}

func (j *jamma) InitForGame(game cartridge.CPS1Game) {
	switch game {
	case cartridge.FFight:

	default:
		panic("Unhandled JAMMA game configuration")
	}
}

func (j *jamma) SetP1Key3(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress, 6, value)
}

func (j *jamma) SetP1Key2(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress, 5, value)
}

func (j *jamma) SetP1Key1(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress, 4, value)
}

func (j *jamma) SetP1Up(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress, 3, value)
}

func (j *jamma) SetP1Down(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress, 2, value)
}

func (j *jamma) SetP1Left(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress, 1, value)
}

func (j *jamma) SetP1Right(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress, 0, value)
}

func (j *jamma) SetP2Key3(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress+1, 6, value)
}

func (j *jamma) SetP2Key2(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress+1, 5, value)
}

func (j *jamma) SetP2Key1(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress+1, 4, value)
}

func (j *jamma) SetP2Up(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress+1, 3, value)
}

func (j *jamma) SetP2Down(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress+1, 2, value)
}

func (j *jamma) SetP2Left(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress+1, 1, value)
}

func (j *jamma) SetP2Right(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaInputAddress+1, 0, value)
}

func (j *jamma) SetCoinService(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaCoinAddress, 6, value)
}

func (j *jamma) SetCoinP2Start(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaCoinAddress, 5, value)
}

func (j *jamma) SetCoinP1Start(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaCoinAddress, 4, value)
}

func (j *jamma) SetCoinCoin2P2(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaCoinAddress, 1, value)
}

func (j *jamma) SetCoinCoinP1(value bool) {
	j.jammaMemory.WriteBit(memoryBus.JammaCoinAddress, 0, value)
}

func (j *jamma) SetDIP_A(bit uint8, value bool) {
	j.jammaMemory.WriteBit(dipAAddress, bit, value)
}

func (j *jamma) DIP_A(bit uint8) bool {
	return j.jammaMemory.ReadBit(dipAAddress, bit)
}
func (j *jamma) SetDIP_B(bit uint8, value bool) {
	j.jammaMemory.WriteBit(dipBAddress, bit, value)
}

func (j *jamma) DIP_B(bit uint8) bool {
	return j.jammaMemory.ReadBit(dipBAddress, bit)
}
func (j *jamma) SetDIP_C(bit uint8, value bool) {
	j.jammaMemory.WriteBit(dipCAddress, bit, value)
}

func (j *jamma) DIP_C(bit uint8) bool {
	return j.jammaMemory.ReadBit(dipCAddress, bit)
}
