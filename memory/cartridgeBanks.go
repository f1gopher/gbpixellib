package memory

import "fmt"

func splitDataIntoBanks(bank0StartAddress uint16, otherBankStartAddress uint16, bankSize uint16, data *[]byte, name string) map[uint8]*Memory {
	banks := make(map[uint8]*Memory)
	var currentBank uint8 = 0

	// Iterate the data and split into chunks of 0x4000 bytes per bank
	for x := 0; x < len(*data); x += int(bankSize) {
		bank := make([]byte, bankSize)

		bank = (*data)[x : x+int(bankSize)]

		offset := bank0StartAddress
		if x > 0 {
			offset = otherBankStartAddress
		}

		result := CreateReadOnlyMemory(fmt.Sprintf("cartridge %s bank %d", name, currentBank), &bank, offset)
		banks[currentBank] = result
		currentBank++
	}

	return banks
}
