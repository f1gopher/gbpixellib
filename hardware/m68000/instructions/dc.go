package instructions

type DC struct{}

func (d *DC) Name() string {
	return "dc.w    $00ff; ILLEGAL"
}

func (d *DC) Length() int {
	return 2
}
