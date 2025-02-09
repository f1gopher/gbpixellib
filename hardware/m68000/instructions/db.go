package instructions

import "fmt"

type db struct {
	reg          byte
	condition    testCondition
	displacement int

	// TODO - only needed for debug output, remove later and use PC value
	tmpStartOffset int
}

func CreateDB(value uint16, code []uint8, offset int) *db {
	displacement := twosComplementToInt16Bit(uint16(readData(code, offset, wordSize)))

	return &db{
		condition:      condition(value),
		reg:            valueFrom3Bits(value, 0),
		displacement:   displacement,
		tmpStartOffset: offset,
	}
}

func (d *db) Name() string {
	return fmt.Sprintf(
		"db%s    %s, $%x",
		d.condition.string(),
		addressString(dataRegisterDirect, d.reg, 0, 0),
		d.tmpStartOffset+d.displacement,
	)
}

func (d *db) Length() int {
	return 4
}
