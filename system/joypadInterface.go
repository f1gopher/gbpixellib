package system

type Joypad interface {
	PressStart()
	ReleaseStart()
	PressSelect()
	ReleaseSelect()
	PressA()
	ReleaseA()
	PressB()
	ReleaseB()
	PressUp()
	ReleaseUp()
	PressDown()
	ReleaseDown()
	PressLeft()
	ReleaseLeft()
	PressRight()
	ReleaseRight()
}
