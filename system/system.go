package system

import (
	"go-boy/cpu"
	"go-boy/display"
	"go-boy/interupt"
	"go-boy/log"
	"go-boy/memory"
	"strings"
	"sync"
)

type System struct {
	bios string
	rom  string

	log             *log.Log
	screen          *display.Screen
	memory          *memory.Memory
	cpu             *cpu.CPU
	interuptHandler *interupt.Handler

	currentDisplay string
	displayLock    sync.Mutex
}

func CreateSystem(bios string, rom string) *System {
	l := log.CreateLog("./log.txt")
	system := System{
		log:    l,
		bios:   bios,
		rom:    rom,
		memory: memory.CreateMemory(l),
	}
	system.cpu = cpu.CreateCPU(l, system.memory)
	system.interuptHandler = interupt.CreateHandler(system.memory, system.cpu)
	system.screen = display.CreateScreen(system.memory, system.interuptHandler)
	//	system.currentDisplay = system.screen.Render()

	system.cpu.Init()

	return &system
}

func CreateTestSystem(testRom string) *System {
	l := log.CreateLog("log.txt")
	system := System{
		bios:   testRom,
		memory: memory.CreateMemory(l),
	}
	system.cpu = cpu.CreateCPU(l, system.memory)
	system.interuptHandler = interupt.CreateHandler(system.memory, system.cpu)
	system.screen = display.CreateScreen(system.memory, system.interuptHandler)
	//	system.currentDisplay = system.screen.Render()

	system.cpu.InitForTestROM()

	return &system
}

func (s *System) Start() {
	err := s.memory.LoadRom(s.rom)
	if err != nil {
		panic(err)
	}

	err = s.memory.LoadBios(s.bios)
	if err != nil {
		panic(err)
	}

	// go s.loop()
}

//func (s *System) Pixels() string {
//	s.displayLock.Lock()
//	defer s.displayLock.Unlock()
//	return s.currentDisplay
//}

func (s *System) Render(callback func(x int, y int, color display.ScreenColor)) {
	s.displayLock.Lock()
	defer s.displayLock.Unlock()

	s.screen.Render(callback)
}

//func (s *System) loop() {
//
//	const maxCycles = 69905
//
//	frameTime := (time.Millisecond * 1000) / 60
//
//	for {
//		start := time.Now().UTC()
//
//		currentCycles := 0
//
//		for currentCycles < maxCycles {
//			currentCycles += s.Tick()
//		}
//
//		//		s.displayLock.Lock()
//		//		s.currentDisplay = s.screen.Render()
//		//		s.displayLock.Unlock()
//
//		// TODO - need to sleep to force 60fps rate
//
//		elapsed := time.Now().UTC().Sub(start)
//
//		sleep := frameTime.Milliseconds() - elapsed.Milliseconds()
//
//		time.Sleep(time.Duration(sleep))
//	}
//}

func (s *System) Tick() int {

	const maxCycles = 69905
	currentCycles := 0

	for currentCycles < maxCycles {
		cyclesCompleted := s.cpu.Tick()

		// Update timers
		s.screen.UpdateForCycles(cyclesCompleted)
		s.interuptHandler.Update()

		currentCycles += cyclesCompleted
	}

	return currentCycles
}

func (s *System) State() string {
	info := s.cpu.Debug()

	parts := strings.Split(info, "->")
	cpu := strings.ReplaceAll(parts[1], " ", "\n")
	ppu := s.screen.Debug()

	return parts[0] + "\n" + cpu + "\n" + ppu
}

func (s *System) OpcodesUsed() {
	s.cpu.DumpOpcodesUsed()
}
