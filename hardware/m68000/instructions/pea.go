package instructions

import "fmt"

type pea struct {
	length int
	eff    addressModeInfo
}

func CreatePEA(value uint16, code []uint8, codeOffset int) *pea {
	eff, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

	if eff.mode == invalidAddressMode ||
		eff.mode == dataRegisterDirect ||
		eff.mode == addressRegisterDirect ||
		eff.mode == addressRegisterIndirectWithPostIncrement ||
		eff.mode == addressRegisterIndirectWithPreDecrement ||
		eff.mode == preIndexedImediate {
		return nil
	}

	return &pea{
		length: 2 + instSize,
		eff:    eff,
	}
}

func (p *pea) Name() string {
	return fmt.Sprintf(
		"pea     %s",
		p.eff.string())
}

func (p *pea) Length() int {
	return p.length
}
