package common

// js keyboard event key code definition
const (
	BACKSPACE    = 0x8
	ENTER        = 0xd
	SPACE        = 0x20
	TAB          = 0x9
	DELETE       = 0x2e
	END          = 0x23
	HOME         = 0x24
	INSERT       = 0x2d
	PAGEDOWN     = 0x22
	PAGEUP       = 0x21
	DOWN         = 0x28
	LEFT         = 0x25
	RIGHT        = 0x27
	UP           = 0x26
	ESC          = 0x1b
	PRINTSCREEN  = 0x2A
	ALT          = 0x12
	SHIFT        = 0x10
	CAPSLOCK     = 0x14
	CTRL         = 0x11
	OS           = 0x5b
	PAUSEBRAKE   = 0x13
	MINUS        = 0xAD
	PERIOD       = 0xbe
	SLASH        = 0xbf
	EQUAL        = 0x3D
	BACKSLASH    = 0xDC
	SEMICOLON    = 0x3B
	SINGLEQUOTE  = 0xDE
	BACKQUOTE    = 0xc0
	COMMA        = 0xbc
	BRACKETLEFT  = 0xdb
	BRACKETRIGHT = 0xdd
	SCROLLLOCK   = 0x91
)

// js keyboard event key code definition - key code base of alphabet, number, function key
const (
	// alphabet key code base
	ALPHABETBASE = 0x41

	ALPHABETLENGTH = 26

	// number key code base
	NUMBERBASE = 0x30

	NUMBERBASELENGHT = 10

	// function key code base
	FUNCTIONBASE = 0x70

	FUNCTIONBASELENGTH = 12
)

// js mouse event key code
const (
	MouseLeftButton = iota
	MouseMiddleButton
	MouseRightButton
)

func IsAltCtrlShiftOsKey(keyCode byte) bool {
	return keyCode == ALT || keyCode == CTRL ||
		keyCode == SHIFT || keyCode == OS
}
