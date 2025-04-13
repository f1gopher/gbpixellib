package m68000

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type jsonTestRegisters struct {
	D0       uint32     `json:"d0"`
	D1       uint32     `json:"d1"`
	D2       uint32     `json:"d2"`
	D3       uint32     `json:"d3"`
	D4       uint32     `json:"d4"`
	D5       uint32     `json:"d5"`
	D6       uint32     `json:"d6"`
	D7       uint32     `json:"d7"`
	A0       uint32     `json:"a0"`
	A1       uint32     `json:"a1"`
	A2       uint32     `json:"a2"`
	A3       uint32     `json:"a3"`
	A4       uint32     `json:"a4"`
	A5       uint32     `json:"a5"`
	A6       uint32     `json:"a6"`
	USP      uint32     `json:"usp"`
	SSP      uint32     `json:"ssp"`
	SR       uint16     `json:"sr"`
	PC       uint32     `json:"pc"`
	Prefetch []uint32   `json:"prefetch"`
	RAM      [][]uint32 `json:"ram"`
}

type jsonTestCase struct {
	Name    string            `json:"name"`
	Initial jsonTestRegisters `json:"initial"`
	Final   jsonTestRegisters `json:"final"`
}

func TestExecution(t *testing.T) {
	testFile := "../../test_data/m68000/MOVE.q.json"
	testData, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	var testCases []jsonTestCase
	if err := json.Unmarshal(testData, &testCases); err != nil {
		t.Fatal(err)
	}

	data := make([]byte, 5000)
	mem := memory.CreateMemory[uint32]("main", &data, 0)
	cpu := CreateM68000(nil, mem)
	cpu.Reset()

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			cpu.test_registers().d0 = test.Initial.D0
			cpu.test_registers().d1 = test.Initial.D1
			cpu.test_registers().d2 = test.Initial.D2
			cpu.test_registers().d3 = test.Initial.D3
			cpu.test_registers().d4 = test.Initial.D4
			cpu.test_registers().d5 = test.Initial.D5
			cpu.test_registers().d6 = test.Initial.D6
			cpu.test_registers().d7 = test.Initial.D7
			cpu.test_registers().a0 = test.Initial.A0
			cpu.test_registers().a1 = test.Initial.A1
			cpu.test_registers().a2 = test.Initial.A2
			cpu.test_registers().a3 = test.Initial.A3
			cpu.test_registers().a4 = test.Initial.A4
			cpu.test_registers().a5 = test.Initial.A5
			cpu.test_registers().a6 = test.Initial.A6
			cpu.test_registers().userStack = test.Initial.USP
			cpu.test_registers().supervisorStack = test.Initial.SSP
			cpu.test_registers().sr = test.Initial.SR
			cpu.test_registers().pc = test.Initial.PC

			verifyRegister(t, test.Final.D0, cpu.test_registers().d0, "D0")
			verifyRegister(t, test.Final.D1, cpu.test_registers().d1, "D1")
			verifyRegister(t, test.Final.D2, cpu.test_registers().d2, "D2")
			verifyRegister(t, test.Final.D3, cpu.test_registers().d3, "D3")
			verifyRegister(t, test.Final.D4, cpu.test_registers().d4, "D4")
			verifyRegister(t, test.Final.D5, cpu.test_registers().d5, "D5")
			verifyRegister(t, test.Final.D6, cpu.test_registers().d6, "D6")
			verifyRegister(t, test.Final.D7, cpu.test_registers().d7, "D7")
			verifyRegister(t, test.Final.A0, cpu.test_registers().a0, "A0")
			verifyRegister(t, test.Final.A1, cpu.test_registers().a1, "A1")
			verifyRegister(t, test.Final.A2, cpu.test_registers().a2, "A2")
			verifyRegister(t, test.Final.A3, cpu.test_registers().a3, "A3")
			verifyRegister(t, test.Final.A4, cpu.test_registers().a4, "A4")
			verifyRegister(t, test.Final.A5, cpu.test_registers().a5, "A5")
			verifyRegister(t, test.Final.A6, cpu.test_registers().a6, "A6")
			verifyRegister(t, test.Final.USP, cpu.test_registers().userStack, "USP")
			verifyRegister(t, test.Final.SSP, cpu.test_registers().supervisorStack, "SSP")
			verifyRegister(t, test.Final.SR, cpu.test_registers().sr, "SR")
			verifyRegister(t, test.Final.PC, cpu.test_registers().pc, "PC")
		})
	}
}

func verifyRegister(t *testing.T, expected uint32, result uint32, name string) {
	if result != expected {
		t.Errorf("Register %s should be 0x%04X but is 0x%04X", name, expected, result)
	}
}
