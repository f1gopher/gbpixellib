package memory

import "fmt"

func SetBit(input byte, bit int, value bool) byte {

	if bit > 7 || bit < 0 {
		panic(fmt.Sprintf("Invalid bit for setRegBit: %d", bit))
	}

	if value {
		input = input | 0x01<<bit
	} else {
		input = input &^ (0x01 << bit)
	}

	return input
}

func GetBit(input byte, bit int) bool {
	if bit > 7 || bit < 0 {
		panic(fmt.Sprintf("Invalid bit for getRegBit: %d", bit))
	}

	if (input>>bit)&0x01 == 0x01 {
		return true
	}

	return false
}
