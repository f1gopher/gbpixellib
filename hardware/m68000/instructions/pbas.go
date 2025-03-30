package instructions

type pbas struct {
	name string
}

func CreatePBAS(value uint16) *pbas {
	return &pbas{name: "pbas e584"}
}

func (p *pbas) Name() string {
	return p.name
}

func (p *pbas) Length() int {
	return 6
}
