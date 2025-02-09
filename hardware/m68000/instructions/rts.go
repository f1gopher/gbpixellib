package instructions

type rts struct{}

func CreateRTS(value uint16) *rts {
	return &rts{}
}

func (r *rts) Name() string {
	return "rts"
}

func (r *rts) Length() int {
	return 2
}
