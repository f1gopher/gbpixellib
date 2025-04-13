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
	A7       uint32     `json:"a7"`
	USP      uint32     `json:"usp"`
	SSP      uint32     `json:"ssp"`
	SR       uint32     `json:"sr"`
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

		})
	}
}
