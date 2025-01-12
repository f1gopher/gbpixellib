package instructions

type rte struct{}

func CreateRTE() *rte {
	return &rte{}
}

func (r *rte) Name() string {
	return "rte"
}

func (r *rte) Length() int {
	return 2
}
