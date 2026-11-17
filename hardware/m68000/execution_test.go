package m68000

import (
	"encoding/json"
	"fmt"
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
	Prefetch []uint16   `json:"prefetch"`
	RAM      [][]uint32 `json:"ram"`
}

type jsonTestCase struct {
	Name    string            `json:"name"`
	Initial jsonTestRegisters `json:"initial"`
	Final   jsonTestRegisters `json:"final"`
}

func TestExecution(t *testing.T) {
	testFiles := []string{
		"../../../external/m68000_tests/v1/NOP.json",
		"../../../external/m68000_tests/v1/MOVE.q.json",
	}

	for _, testFile := range testFiles {
		testData, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatal(err)
		}

		var testCases []jsonTestCase
		if err := json.Unmarshal(testData, &testCases); err != nil {
			t.Fatal(err)
		}

		data := make([]byte, 0xFFFFFF)
		mem := memory.CreateMemoryBigEndian[uint32]("main", &data, 0)
		cpu := CreateM68000(nil, mem)
		cpu.Reset()

		for _, test := range testCases {

			fullName := fmt.Sprintf("%s - 0x%04X %016b", test.Name, test.Initial.Prefetch[0], test.Initial.Prefetch[0])

			t.Run(fullName, func(t *testing.T) {
				cpu.Reset()

				cpu.test_registers().Set32(D0, test.Initial.D0)
				cpu.test_registers().Set32(D1, test.Initial.D1)
				cpu.test_registers().Set32(D2, test.Initial.D2)
				cpu.test_registers().Set32(D3, test.Initial.D3)
				cpu.test_registers().Set32(D4, test.Initial.D4)
				cpu.test_registers().Set32(D5, test.Initial.D5)
				cpu.test_registers().Set32(D6, test.Initial.D6)
				cpu.test_registers().Set32(D7, test.Initial.D7)
				cpu.test_registers().Set32(A0, test.Initial.A0)
				cpu.test_registers().Set32(A1, test.Initial.A1)
				cpu.test_registers().Set32(A2, test.Initial.A2)
				cpu.test_registers().Set32(A3, test.Initial.A3)
				cpu.test_registers().Set32(A4, test.Initial.A4)
				cpu.test_registers().Set32(A5, test.Initial.A5)
				cpu.test_registers().Set32(A6, test.Initial.A6)
				cpu.test_registers().Set32(UserStack, test.Initial.USP)
				cpu.test_registers().Set32(SupervisorStack, test.Initial.SSP)
				cpu.test_registers().Set16(SR, test.Initial.SR)
				cpu.test_registers().Set32(PC, test.Initial.PC)

				for x := range test.Initial.RAM {
					mem.WriteByte(test.Initial.RAM[x][0], uint8(test.Initial.RAM[x][1]))
				}

				cpu.test_registers().Set16(IRC_Prefetch, test.Initial.Prefetch[1])
				cpu.test_registers().Set16(IR_Decoding, test.Initial.Prefetch[0])

				cpu.Step()

				verifyRegister(t, test.Final.D0, cpu.test_registers().Get32(D0), "D0")
				verifyRegister(t, test.Final.D1, cpu.test_registers().Get32(D1), "D1")
				verifyRegister(t, test.Final.D2, cpu.test_registers().Get32(D2), "D2")
				verifyRegister(t, test.Final.D3, cpu.test_registers().Get32(D3), "D3")
				verifyRegister(t, test.Final.D4, cpu.test_registers().Get32(D4), "D4")
				verifyRegister(t, test.Final.D5, cpu.test_registers().Get32(D5), "D5")
				verifyRegister(t, test.Final.D6, cpu.test_registers().Get32(D6), "D6")
				verifyRegister(t, test.Final.D7, cpu.test_registers().Get32(D7), "D7")
				verifyRegister(t, test.Final.A0, cpu.test_registers().Get32(A0), "A0")
				verifyRegister(t, test.Final.A1, cpu.test_registers().Get32(A1), "A1")
				verifyRegister(t, test.Final.A2, cpu.test_registers().Get32(A2), "A2")
				verifyRegister(t, test.Final.A3, cpu.test_registers().Get32(A3), "A3")
				verifyRegister(t, test.Final.A4, cpu.test_registers().Get32(A4), "A4")
				verifyRegister(t, test.Final.A5, cpu.test_registers().Get32(A5), "A5")
				verifyRegister(t, test.Final.A6, cpu.test_registers().Get32(A6), "A6")
				verifyRegister(t, test.Final.USP, cpu.test_registers().Get32(UserStack), "USP")
				verifyRegister(t, test.Final.SSP, cpu.test_registers().Get32(SupervisorStack), "SSP")
				if cpu.test_registers().Get16(SR) != test.Initial.SR {
					t.Errorf("SR should be 0x%02X but is 0x%02X", test.Initial.SR, cpu.test_registers().Get16(SR))
				}
				verifyRegister(t, test.Final.PC, cpu.test_registers().Get32(PC), "PC")

				for x := range test.Final.RAM {
					actual := mem.ReadByte(test.Final.RAM[x][0])
					expected := uint8(test.Final.RAM[x][1])

					if actual != expected {
						t.Errorf("RAM at 0x%04X should be 0x%02X but is 0x%02X", test.Final.RAM[x][0], expected, actual)
					}
				}

				actualIRC := cpu.test_registers().Get16(IRC_Prefetch)
				if actualIRC != test.Final.Prefetch[1] {
					t.Errorf("IRC Prefetch should be 0x%04X but is 0x%04X", test.Final.Prefetch[1], actualIRC)
				}
				actualIR := cpu.test_registers().Get16(IR_Decoding)
				if actualIR != test.Final.Prefetch[0] {
					t.Errorf("IR Decoding should be 0x%04X but is 0x%04X", test.Final.Prefetch[0], actualIR)
				}
			})
		}
	}
}

func verifyRegister(t *testing.T, expected uint32, result uint32, name string) {
	if result != expected {
		t.Errorf("Register %s should be 0x%08X but is 0x%08X", name, expected, result)
	}
}
