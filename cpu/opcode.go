package cpu

type opcode interface {
	doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error)
	name() string
	opcode() uint8
	length() uint8
}

type opcodeBase struct {
	opcodeName   string
	opcodeId     uint8
	opcodeLength uint8
}

func (o *opcodeBase) name() string { return o.opcodeName }

func (o *opcodeBase) opcode() uint8 { return o.opcodeId }

func (o *opcodeBase) length() uint8 { return o.opcodeLength }
