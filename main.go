package main

import (
	"go-boy/system"
	"go-boy/ui"
)

func main() {
	//log.Info("Hello World")

	//screen := display.CreateScreen()
	//memory := memory.CreateMemory()
	//cpu := cpu.CreateCPU(memory)

	//err := memory.LoadBios("./rom/test/08-misc instrs.gb")
	//if err != nil {
	//	log.Error(err)
	//}

	//	memory.DumpBios()
	//cpu.SetDebugLog("./log.txt")
	//cpu.InitForTestROM()

	//for {
	//	cpu.Tick()
	// time.Sleep(time.Millisecond * 10)
	//}

	system := system.CreateSystem("./bios/dmg.bin", "")
	ui.Main(system)
}
