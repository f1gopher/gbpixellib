package instructions

type pmove struct {
	name string
}

func CreatePMOVE(value uint16) *pmove {
	return &pmove{name: "pmove [unknown form] D0"}
}

func (p *pmove) Name() string {
	return p.name
}

func (p *pmove) Length() int {
	return 4
}
