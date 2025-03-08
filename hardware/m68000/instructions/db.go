package instructions

import "fmt"

type db struct {
	count        addressModeInfo
	condition    testCondition
	displacement int

	// TODO - only needed for debug output, remove later and use PC value
	tmpStartOffset int
}

func CreateDB(value uint16, code []uint8, codeOffset int) *db {
	displacement := twosComplementToInt16Bit(uint16(readData(code, codeOffset, wordSize)))
	count := addressModeInfo{
		reg:  valueFrom3Bits(value, 0),
		mode: dataRegisterDirect,
	}

	return &db{
		count:          count,
		condition:      condition(value),
		displacement:   displacement,
		tmpStartOffset: codeOffset,
	}
}

func (d *db) Name() string {
	return fmt.Sprintf(
		"db%s    %s, $%x",
		d.condition.string(),
		d.count.string(),
		d.tmpStartOffset+d.displacement,
	)
}

func (d *db) Length() int {
	return 4
}
