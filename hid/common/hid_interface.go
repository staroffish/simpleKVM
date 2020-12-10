package common

type Hid interface {
	OpenDevice(args ...string) error
	CloseDevice() error
	GetModelName() string
	Mouse
	KeyBoard
}

type Mouse interface {
	MoveTo(Xpoint uint16, Ypoint uint16) error
	MouseDown(button int) error
	MouseUp(button int) error
	MouseScroll(scroll int) error
}

type KeyBoard interface {
	KeyDown(keyCode byte) error
	KeyUp(keyCode byte) error
	// KeyBoardCodeToDeviceCode(eventKeyCode byte) (byte, error)
}

type BaseHid struct{}

func (c *BaseHid) OpenDevice(...string) error                { return nil }
func (c *BaseHid) CloseDevice() error                        { return nil }
func (c *BaseHid) GetModelName() string                      { return "base" }
func (c *BaseHid) KeyDown(eventKeyCode byte) error           { return nil }
func (c *BaseHid) KeyUp(eventKeyCode byte) error             { return nil }
func (c *BaseHid) MoveTo(Xpoint uint16, Ypoint uint16) error { return nil }
func (c *BaseHid) MouseDown(button int) error                { return nil }
func (c *BaseHid) MouseUp(button int) error                  { return nil }
func (c *BaseHid) MouseScroll(scroll int) error              { return nil }

// func (c *BaseHid) KeyBoardCodeToDeviceCode(eventKeyCode byte) (byte, error) { return 0, nil }
