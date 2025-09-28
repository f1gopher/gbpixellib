package m68000

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type loggingMemory struct {
	mem          memory.Memory[uint32]
	transactions []jsonTransaction
}

func (l *loggingMemory) Size() int {
	return l.mem.Size()
}
func (l *loggingMemory) Offset() uint32 {
	return l.mem.Offset()
}
func (l *loggingMemory) LastAddress() uint32 {
	return l.mem.LastAddress()
}
func (l *loggingMemory) Reset() {
	l.mem.Reset()
}
func (l *loggingMemory) ReadBit(address uint32, bit uint8) bool {
	return l.mem.ReadBit(address, bit)
}
func (l *loggingMemory) ReadByte(address uint32) byte {
	result := l.mem.ReadByte(address)
	l.transactions = append(l.transactions, jsonTransaction{op: "r", address: address, size: ".b", value: uint32(result)})
	return result
}
func (l *loggingMemory) ReadShort(address uint32) uint16 {
	result := l.mem.ReadShort(address)
	l.transactions = append(l.transactions, jsonTransaction{op: "r", address: address, size: ".w", value: uint32(result)})
	return result
}
func (l *loggingMemory) ReadLong(address uint32) uint32 {
	result := l.mem.ReadLong(address)
	l.transactions = append(l.transactions, jsonTransaction{op: "r", address: address, size: ".l", value: result})
	return result
}
func (l *loggingMemory) WriteBit(address uint32, bit uint8, value bool) {
	l.mem.WriteBit(address, bit, value)
}
func (l *loggingMemory) WriteByte(address uint32, value byte) {
	l.mem.WriteByte(address, value)
	l.transactions = append(l.transactions, jsonTransaction{op: "w", address: address, size: ".b", value: uint32(value)})
}
func (l *loggingMemory) WriteShort(address uint32, value uint16) {
	l.mem.WriteShort(address, value)
	l.transactions = append(l.transactions, jsonTransaction{op: "w", address: address, size: ".w", value: uint32(value)})
}
func (l *loggingMemory) WriteLong(address uint32, value uint32) {
	l.mem.WriteLong(address, value)
	l.transactions = append(l.transactions, jsonTransaction{op: "w", address: address, size: ".l", value: value})
}
func (l *loggingMemory) DumpCode() (data []uint8, startAddress uint32) {
	return l.mem.DumpCode()
}
func (l *loggingMemory) DumpRegionsCode(address uint32) (data []uint8, err error) {
	return l.mem.DumpRegionsCode(address)
}
func (l *loggingMemory) LoadData(offset uint32, data []uint8) error {
	return l.mem.LoadData(offset, data)
}
func (l *loggingMemory) RawBuffer() []uint8 {
	return l.mem.RawBuffer()
}

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

type jsonTransaction struct {
	op       string
	unknown1 int
	unknown2 int
	address  uint32
	size     string
	value    uint32
}

func (j *jsonTransaction) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &[]interface{}{&j.op, &j.unknown1, &j.unknown2, &j.address, &j.size, &j.value})
}

type jsonTestCase struct {
	Name           string            `json:"name"`
	Initial        jsonTestRegisters `json:"initial"`
	Final          jsonTestRegisters `json:"final"`
	NumClockCycles uint64            `json:"length"`
	Transactions   []jsonTransaction `json:"transactions"`
}

func TestExecution(t *testing.T) {
	testFiles := []string{
		"../../test_data/m68000/CLR.b.json",
		"../../test_data/m68000/MOVE.q.json",
		"../../test_data/m68000/NOP.json",
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
		logMem := loggingMemory{mem: mem, transactions: []jsonTransaction{}}
		cpu := CreateM68000(nil, &logMem)
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

				clockCyclesRan := cpu.Step()

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

				if test.NumClockCycles != clockCyclesRan {
					t.Errorf("Expected %d clock cycles but took %d", test.NumClockCycles, clockCyclesRan)
				}

				if len(test.Transactions) != len(logMem.transactions) {
					t.Errorf("Transaction count mismatch. Expected %d but found %d", len(test.Transactions), len(logMem.transactions))
				}
				for x := range test.Transactions {
					if x >= len(logMem.transactions) {
						continue
					}
					expected := test.Transactions[x]
					result := logMem.transactions[x]

					if expected.op != result.op {
						t.Errorf("Txn %d op - expected %s but found %s", x, expected.op, result.op)
					}

					if expected.address != result.address {
						t.Errorf("Txn %d address - expected 0x%08X but found 0x%08X", x, expected.address, result.address)
					}

					if expected.size != result.size {
						t.Errorf("Txn %d size - expected %s but found %s", x, expected.size, result.size)
					}

					if expected.value != result.value {
						t.Errorf("Txn %d value - expected 0x%08X but found 0x%08X", x, expected.value, result.value)
					}
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
