package m68000

import "fmt"

// TODO - do we need to know the size of the value so we know which is the MSB?
func nzvc_NegZeroClrClr(registers *Registers, size instSize, value uint32) {
	// Set if negative, check sign bit
	switch size {
	case byteSize:
		registers.setFlag(N, value&0b00000000000000000000000010000000 == 0b00000000000000000000000010000000)
	case wordSize:
		registers.setFlag(N, value&0b00000000000000001000000000000000 == 0b00000000000000001000000000000000)
	case longSize:
		registers.setFlag(N, value&0b10000000000000000000000000000000 == 0b10000000000000000000000000000000)
	default:
		panic(fmt.Errorf("Unhandled size setting status flag %d &08b", size))
	}
	// Set if zero
	registers.setFlag(Z, value == 0)
	// Clear both
	registers.setFlag(V, false)
	registers.setFlag(C, false)
}
