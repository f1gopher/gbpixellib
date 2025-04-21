package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type addi struct {
	length int
	opSize instSize
	dest   addressModeInfo
	data   int32
}

func CreateADDI(value uint16, code []uint8, codeOffset int) *addi {
	size := twoBitSize(value, 6)
	dest, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset+size.bytesUsed(), codeOffset)

	if dest.mode == addressRegisterDirect ||
		dest.mode == programCounterIndirectWithDisplacement ||
		dest.mode == preIndexedImediate ||
		dest.mode == programCounterIndirectWithIndex8bitDisplacement ||
		dest.mode == invalidAddressMode ||
		size == zeroSize {
		return nil
	}
	var data int32
	tmp := readData(code, codeOffset, size)
	if size == wordSize {
		data = signExtend16BitTo32Bits(uint16(tmp))
	} else if size == longSize {
		data = twosComplementToInt32Bit(tmp)
	} else {
		data = int32(twosComplementToInt8Bit(uint8(tmp)))
	}

	return &addi{
		length: 2 + instSize + size.bytesUsed(),
		opSize: size,
		dest:   dest,
		data:   data,
	}
}

func (a *addi) Name() string {
	if a.data < 0 {
		return fmt.Sprintf(
			"addi.%s  #-$%x, %s",
			a.opSize.string(),
			-a.data,
			a.dest.string())
	}

	return fmt.Sprintf(
		"addi.%s  #$%x, %s",
		a.opSize.string(),
		a.data,
		a.dest.string())
}

func (a *addi) Length() int {
	return a.length
}

func (a *addi) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
