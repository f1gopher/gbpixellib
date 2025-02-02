package instructions

import "fmt"

type dbCondition int

const (
	CarryClear dbCondition = iota
	CarrySet
	Equal
	False
	GreaterOrEqual
	GreaterThan
	High
	LessOrEqual
	LowOrSame
	LessThan
	Minus
	NotEqual
	Plus
	True
	OverflowClear
	OverflowSet
)

func (d dbCondition) string() string {
	switch d {
	case True:
		return "t"
	case False:
		return "ra"
	case High:
		return "hi"
	case LowOrSame:
		return "ls"
	case CarryClear:
		return "cc"
	case CarrySet:
		return "cs"
	case NotEqual:
		return "ns"
	case Equal:
		return "eq"
	case OverflowClear:
		return "vc"
	case OverflowSet:
		return "vs"
	case Plus:
		return "pl"
	case Minus:
		return "mi"
	case GreaterOrEqual:
		return "ge"
	case LessThan:
		return "lt"
	case GreaterThan:
		return "gt"
	case LessOrEqual:
		return "le"
	default:
		panic(fmt.Sprintf("Unknown DBcc condition for string: %v", d))
	}
}

type db struct {
	reg          byte
	condition    dbCondition
	displacement int

	// TODO - only needed for debug output, remove later and use PC value
	tmpStartOffset int
}

func CreateDB(value uint16, code []uint8, offset int) *db {
	displacement := twosComplementToInt(uint16(readData(code, offset, wordSize)))

	return &db{
		condition:      condition(value),
		reg:            valueFrom3Bits(value, 0),
		displacement:   displacement,
		tmpStartOffset: offset,
	}
}

func (d *db) Name() string {
	fmt.Printf("%d\n", d.displacement)

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

func condition(value uint16) dbCondition {
	condition := (value >> 8) & 0b000000000001111
	switch condition {
	case 0b0000000000000000:
		return True
	case 0b0000000000000001:
		return False
	case 0b0000000000000010:
		return High
	case 0b0000000000000011:
		return LowOrSame
	case 0b0000000000000100:
		return CarryClear
	case 0b0000000000000101:
		return CarrySet
	case 0b0000000000000110:
		return NotEqual
	case 0b0000000000000111:
		return Equal
	case 0b0000000000001000:
		return OverflowClear
	case 0b0000000000001001:
		return OverflowSet
	case 0b0000000000001010:
		return Plus
	case 0b0000000000001011:
		return Minus
	case 0b0000000000001100:
		return GreaterOrEqual
	case 0b0000000000001101:
		return LessThan
	case 0b0000000000001110:
		return GreaterThan
	case 0b0000000000001111:
		return LessOrEqual
	default:
		panic(fmt.Sprintf("Unknown DBcc condition: %016b", condition))
	}
}
