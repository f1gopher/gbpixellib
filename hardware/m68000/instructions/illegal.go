package instructions

type illegal struct{}

func CreateILLEGAL() *illegal {
	return &illegal{}
}

func (i *illegal) Name() string {
	return "ILLEGAL"
}

func (i *illegal) Length() int {
	return 2
}
