package m68000

import "fmt"

func nzvc_NegZeroClrClr(registers *Registers, size instSize, value uint32) {
	// N - Set if negative, check sign bit
	// Z - Set if zero
	switch size {
	case byteSize:
		registers.setFlag(N, value&0b00000000000000000000000010000000 == 0b00000000000000000000000010000000)
		registers.setFlag(Z, value & 0x000000FF == 0)
	case wordSize:
		registers.setFlag(N, value&0b00000000000000001000000000000000 == 0b00000000000000001000000000000000)
		registers.setFlag(Z, value & 0x0000FFFF == 0)
	case longSize:
		registers.setFlag(N, value&0b10000000000000000000000000000000 == 0b10000000000000000000000000000000)
		registers.setFlag(Z, value == 0)
	default:
		panic(fmt.Errorf("Unhandled size setting status flag %d &08b", size))
	}

	// Clear both
	registers.setFlag(V, false)
	registers.setFlag(C, false)
}
