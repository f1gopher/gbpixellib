package gameboy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/hardware/z80"
	"github.com/f1gopher/gbpixellib/system/gameboy/cartridge"
)

type z80_test_registers struct {
	A   int     `json:"a"`
	B   int     `json:"b"`
	C   int     `json:"c"`
	D   int     `json:"d"`
	E   int     `json:"e"`
	F   int     `json:"f"`
	H   int     `json:"h"`
	L   int     `json:"l"`
	PC  int     `json:"pc"`
	SP  int     `json:"sp"`
	RAM [][]int `json:"ram"`
}

type z80_test struct {
	Name    string             `json:"name"`
	Initial z80_test_registers `json:"initial"`
	Final   z80_test_registers `json:"final"`
	//Cycles []uint8 `json:"cycles"`
}

type testBus struct {
	mem memory.Memory[uint16]
}

func (b *testBus) ReadBit(address uint16, bit uint8) bool {
	return b.mem.ReadBit(address, bit)
}

func (b *testBus) ReadByte(address uint16) byte {
	return b.mem.ReadByte(address)
}

func (b *testBus) ReadShort(address uint16) uint16 {
	return b.mem.ReadShort(address)
}
func (b *testBus) WriteBit(address uint16, bit uint8, value bool) {
	b.mem.WriteBit(address, bit, value)
}
func (b *testBus) WriteByte(address uint16, value byte) {
	b.mem.WriteByte(address, value)
}
func (b *testBus) WriteShort(address uint16, value uint16) {
	b.mem.WriteShort(address, value)
}
func (b *testBus) DisplaySetScanline(value uint8) {
}

func (b *testBus) SetIO(io inputOutput, interupt overflowInterupt) {}
func (b *testBus) DisplaySetStatus(value uint8)                    {}
func (b *testBus) ExecuteDMAIfPending() bool                       { return false }
func (b *testBus) Reset()                                          {}
func (b *testBus) WriteDividerRegister(value uint8)                {}
func (b *testBus) SetTimer(timer timerDivide)                      {}
func (b *testBus) LoadTestBIOS(cartridge cartridge.Cartridge)      {}

func TestExecution(t *testing.T) {
	testFilesFolder := "../../test_data/gameboy/v2"
	files, err := os.ReadDir(testFilesFolder)
	if err != nil {
		t.Fatal(err)
	}

	// TODO - temp to work through test one at a time
	//	start := 220
	//files = files[start : start+5]

	for _, testFile := range files {
		if filepath.Ext(testFile.Name()) != ".json" {
			continue
		}

		testData, err := os.ReadFile(testFilesFolder + "/" + testFile.Name())
		if err != nil {
			t.Fatal(err)
		}

		var testCases []z80_test
		if err := json.Unmarshal(testData, &testCases); err != nil {
			t.Fatal(err)
		}

		regs := z80.Registers{}
		testMem := make([]byte, 0xFFFF)
		bus := &testBus{
			mem: memory.CreateMemoryLittleEndian[uint16](
				"test",
				&testMem,
				0),
		}
		cpu := z80.CreateCPU(nil, &regs, bus)
		cpu.Reset()

		for _, test := range testCases {
			cpu.Reset()
			bus.Reset()

			regs.Set8(z80.A, uint8(test.Initial.A))
			regs.Set8(z80.B, uint8(test.Initial.B))
			regs.Set8(z80.C, uint8(test.Initial.C))
			regs.Set8(z80.D, uint8(test.Initial.D))
			regs.Set8(z80.E, uint8(test.Initial.E))
			regs.Set8(z80.F, uint8(test.Initial.F))
			regs.Set8(z80.H, uint8(test.Initial.H))
			regs.Set8(z80.L, uint8(test.Initial.L))
			regs.Set16(z80.SP, uint16(test.Initial.SP))
			regs.Set16(z80.PC, uint16(test.Initial.PC))

			for x := range test.Initial.RAM {
				bus.WriteByte(uint16(test.Initial.RAM[x][0]), uint8(test.Initial.RAM[x][1]))
			}

			cpu.TestInit()
			opcode, _ := cpu.GetNextOpcode()
			opcodeName, _ := cpu.GetOpcodeInfo(opcode, 0x00)
			fullName := fmt.Sprintf("%s - %s", opcodeName, test.Name)

			t.Run(fullName, func(t *testing.T) {

				for {
					finished, err := cpu.ExecuteUnitTestMCycle()
					if err != nil {
						t.Error(err)
						break
					}
					if finished {
						break
					}
				}

				verifyRegister(t, uint8(test.Final.A), regs.Get8(z80.A), "A")
				verifyRegister(t, uint8(test.Final.B), regs.Get8(z80.B), "B")
				verifyRegister(t, uint8(test.Final.C), regs.Get8(z80.C), "C")
				verifyRegister(t, uint8(test.Final.D), regs.Get8(z80.D), "D")
				verifyRegister(t, uint8(test.Final.E), regs.Get8(z80.E), "E")
				verifyRegister(t, uint8(test.Final.F), regs.Get8(z80.F), "F")
				verifyRegister(t, uint8(test.Final.H), regs.Get8(z80.H), "H")
				verifyRegister(t, uint8(test.Final.L), regs.Get8(z80.L), "L")
				verifyRegister16(t, uint16(test.Final.SP), regs.Get16(z80.SP), "SP")
				verifyRegister16(t, uint16(test.Final.PC), regs.Get16(z80.PC), "PC")

				for x := range test.Final.RAM {
					actual := bus.ReadByte(uint16(test.Final.RAM[x][0]))
					expected := uint8(test.Final.RAM[x][1])

					if actual != expected {
						t.Errorf("RAM at 0x%04X should be 0x%02X but is 0x%02X", test.Final.RAM[x][0], expected, actual)
					}
				}
			})
		}
	}
}

func verifyRegister(t *testing.T, expected uint8, result uint8, name string) {
	if result != expected {
		t.Errorf("Register %s should be 0x%02X but is 0x%02X", name, expected, result)
	}
}

func verifyRegister16(t *testing.T, expected uint16, result uint16, name string) {
	if result != expected {
		t.Errorf("Register %s should be 0x%04X but is 0x%04X", name, expected, result)
	}
}
