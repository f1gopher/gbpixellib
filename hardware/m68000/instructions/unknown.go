package instructions

type Unknown struct{}

func (u *Unknown) Name() string {
	return "Unknown instruction"
}

func (u *Unknown) Length() int {
	return 2
}
