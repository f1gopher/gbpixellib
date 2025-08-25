package sound

import (
	"fmt"

	"github.com/ebitengine/oto/v3"
	"github.com/f1gopher/gbpixellib/hardware/z80"
)

type Sound struct {
	memory z80.MemoryInterface
	ctx    *oto.Context
}

func CreateSound(memory z80.MemoryInterface) *Sound {
	options := oto.NewContextOptions{
		ChannelCount: 2,
		SampleRate:   44100,
		Format:       oto.FormatSignedInt16LE,
	}
	// TODO - need to wait for the channel to close for it to be ready
	ctx, ready, err := oto.NewContext(&options)
	if err != nil {
		panic(err)
	}
	// Wait for Oto to be ready
	<-ready

	return &Sound{
		memory: memory,
		ctx:    ctx,
	}
}

func (s *Sound) Shutdown() {
}

func (s *Sound) Reset() {
	// Default wave memory values. R-type will need this becaue it
	// doesn't initialise the memory first
	s.memory.WriteByte(0xFF30, 0x84)
	s.memory.WriteByte(0xFF31, 0x40)
	s.memory.WriteByte(0xFF32, 0x43)
	s.memory.WriteByte(0xFF33, 0xAA)
	s.memory.WriteByte(0xFF34, 0x2D)
	s.memory.WriteByte(0xFF35, 0x78)
	s.memory.WriteByte(0xFF36, 0x92)
	s.memory.WriteByte(0xFF37, 0x3C)
	s.memory.WriteByte(0xFF38, 0x60)
	s.memory.WriteByte(0xFF39, 0x59)
	s.memory.WriteByte(0xFF3A, 0x59)
	s.memory.WriteByte(0xFF3B, 0xB0)
	s.memory.WriteByte(0xFF3C, 0x34)
	s.memory.WriteByte(0xFF3D, 0xB8)
	s.memory.WriteByte(0xFF3E, 0x2E)
	s.memory.WriteByte(0xFF3F, 0xDa)
}

func (s *Sound) Update() {
	if s.ControlStatus_AllOn() || s.ControlStatus_Channel1Enabled() {
		if s.Square1_Trigger() {
			fmt.Println("Play")
		}
	}
}
