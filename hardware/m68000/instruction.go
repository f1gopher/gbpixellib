package m68000

type instruction interface {
	Name() string
	Length() int
}
