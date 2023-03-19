package system

import (
	"go-boy/cpu"
	"go-boy/display"
	"go-boy/memory"
	"log"
	"strings"
	"sync"
	"time"
)

type System struct {
	bios string
	rom  string

	screen *display.Screen
	memory *memory.Memory
	cpu    *cpu.CPU

	currentDisplay string
	displayLock    sync.Mutex
}

func CreateSystem(bios string, rom string) *System {
	system := System{
		bios:   bios,
		rom:    rom,
		memory: memory.CreateMemory(),
	}
	system.cpu = cpu.CreateCPU(system.memory)
	system.screen = display.CreateScreen(system.memory)

	system.currentDisplay = system.screen.Render()

	return &system
}

func (s *System) Start() {
	s.cpu.Init()

	err := s.memory.LoadRom(s.rom)
	if err != nil {
		log.Panic(err)
	}

	err = s.memory.LoadBios(s.bios)
	if err != nil {
		log.Panic(err)
	}

	//go s.loop()

}

func (s *System) Pixels() string {
	s.displayLock.Lock()
	defer s.displayLock.Unlock()
	return s.currentDisplay
}

func (s *System) loop() {
	for {
		s.Tick()
		time.Sleep(time.Millisecond * 5)
	}
}

func (s *System) Tick() {
	s.cpu.Tick()

	s.displayLock.Lock()
	s.currentDisplay = s.screen.Render()
	s.displayLock.Unlock()
}

func (s *System) State() string {
	cpu := strings.ReplaceAll(s.cpu.Debug(), " ", "\n")
	ppu := s.screen.Debug()
	return cpu + "\n" + ppu
}
