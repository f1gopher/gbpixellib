package main

import (
	"go-boy/system"
	"go-boy/ui"
	"log"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//log.Info("Hello World")

	//memory := memory.CreateMemory()
	//cpu := cpu.CreateCPU(memory)
	//screen := display.CreateScreen(memory)

	//memory.LoadBios("./rom/test/01-special.gb")
	//if err != nil {
	//	log.Error(err)
	//}

	//	memory.DumpBios()
	//cpu.SetDebugLog("./log.txt")
	//cpu.InitForTestROM()

	//for {
	//	cpu.Tick()
	//	time.Sleep(time.Millisecond * 10)
	//}

	//system := system.CreateTestSystem("./rom/test/06-ld r,r.gb")
	system := system.CreateSystem("./bios/dmg.bin", "./rom/games/Tetris (World) (Rev 1).gb")
	system.Start()
	ui.Main(system)
}
