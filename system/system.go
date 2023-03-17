package system

import (
	"go-boy/cpu"
	"go-boy/display"
	"go-boy/memory"
	"log"
	"time"
)

type System struct {
	bios string
	rom  string

	screen *display.Screen
	memory *memory.Memory
	cpu    *cpu.CPU
}

func CreateSystem(bios string, rom string) *System {
	system := System{
		bios:   bios,
		rom:    rom,
		memory: memory.CreateMemory(),
	}
	system.cpu = cpu.CreateCPU(system.memory)
	system.screen = display.CreateScreen(system.memory)

	return &system
}

func (s *System) Start() {
	err := s.memory.LoadBios(s.bios)
	if err != nil {
		log.Panic(err)
	}

	go s.loop()
}

func (s *System) Pixels() string {
	return s.screen.Render()
}

func (s *System) loop() {
	for {
		s.cpu.Tick()
		time.Sleep(time.Millisecond * 5)
	}
}
