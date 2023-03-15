package main

import (
	"go-boy/cpu"
	"go-boy/memory"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	log.Info("Hello World")

	//screen := display.CreateScreen()
	memory := memory.CreateMemory()
	cpu := cpu.CreateCPU(memory)

	err := memory.LoadBios("./rom/test/06-ld r,r.gb")
	if err != nil {
		log.Error(err)
	}

	//	memory.DumpBios()
	cpu.SetDebugLog("./log.txt")
	cpu.InitForTestROM()

	for {
		cpu.Tick()
		time.Sleep(time.Millisecond * 10)
	}
}
