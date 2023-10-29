package timer

type memoryInterface interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)
}

type Timer struct {
	mem memoryInterface
}

const DividerRegister = 0xFF04

func CreateTimer(mem memoryInterface) *Timer {
	return &Timer{
		mem: mem,
	}
}

func (t *Timer) Update(cycles uint8) {
	current := t.mem.ReadByte(DividerRegister)
	current += cycles / 2
	t.mem.WriteByte(DividerRegister, current)
}
