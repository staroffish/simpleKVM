package ch9329

import (
	"fmt"

	"github.com/staroffish/simpleKVM/hid/common"
)

func (c *ch9329) KeyBoardCodeToDeviceCode(eventCode byte) (byte, error) {

	// alphabet key
	if eventCode >= common.ALPHABETBASE && eventCode < (common.ALPHABETBASE+common.ALPHABETLENGTH) {
		return ALPHABETBASE + (eventCode - common.ALPHABETBASE), nil
	}

	// number key
	if eventCode == common.NUMBERBASE {
		return NUMBERBASE + common.NUMBERBASELENGHT - 1, nil
	} else if eventCode > common.NUMBERBASE && eventCode < common.NUMBERBASE+common.NUMBERBASELENGHT {
		return NUMBERBASE + (eventCode - common.NUMBERBASE - 1), nil
	}

	// function key
	if eventCode >= common.FUNCTIONBASE && eventCode < (common.FUNCTIONBASE+common.FUNCTIONBASELENGTH) {
		return FUNCTIONBASE + (eventCode - common.FUNCTIONBASE), nil
	}

	// other key
	keyCode, ok := keyMapping[eventCode]
	if ok {
		return keyCode, nil
	}

	return 0, fmt.Errorf("known code: %v", eventCode)
}

// ch9329 key definition - key code
const (
	ENTER = iota + 0x28
	ESC
	BACKSPACE
	TAB
	SPACE
	MINUS
	EQUAL
	BRACKETLEFT
	BRACKETRIGHT
	BACKSLASH
	_
	SEMICOLON
	SINGLEQUOTE
	BACKQUOTE
	COMMA
	PERIOD
	SLASH
	CAPSLOCK
)

const (
	PRINTSCREEN = iota + 0x46
	SCORLLLOCK
	PAUSEBRAKE
	INSERT
	HOME
	PAGEUP
	DELETE
	END
	PAGEDOWN
	RIGHT
	LEFT
	DOWN
	UP
)

const (
	CTRLLEFT   = 0x01
	CTRLRIGHT  = 0x10
	SHIFTLEFT  = 0x02
	SHIFTRIGHT = 0x20
	ALTLEFT    = 0x04
	ALTRIGHT   = 0x40
	OSLEFT     = 0x08
	OSRIGHT    = 0x80
)

// ch9329 key definition - key code base of alphabet, number, function key
const (
	// alphabet key code base
	ALPHABETBASE = 0x04

	// number key code base
	NUMBERBASE = 0x1e

	// function key code base
	FUNCTIONBASE = 0x3a
)

var keyMapping map[byte]byte = map[byte]byte{
	common.ENTER:        ENTER,
	common.ESC:          ESC,
	common.BACKSPACE:    BACKSPACE,
	common.TAB:          TAB,
	common.SPACE:        SPACE,
	common.MINUS:        MINUS,
	common.EQUAL:        EQUAL,
	common.BRACKETLEFT:  BRACKETLEFT,
	common.BRACKETRIGHT: BRACKETRIGHT,
	common.BACKSLASH:    BACKSLASH,
	common.SEMICOLON:    SEMICOLON,
	common.SINGLEQUOTE:  SINGLEQUOTE,
	common.BACKQUOTE:    BACKQUOTE,
	common.COMMA:        COMMA,
	common.PERIOD:       PERIOD,
	common.SLASH:        SLASH,
	common.CAPSLOCK:     CAPSLOCK,
	common.PRINTSCREEN:  PRINTSCREEN,
	common.SCORLLLOCK:   SCORLLLOCK,
	common.PAUSEBRAKE:   PAUSEBRAKE,
	common.INSERT:       INSERT,
	common.HOME:         HOME,
	common.PAGEUP:       PAGEUP,
	common.DELETE:       DELETE,
	common.END:          END,
	common.PAGEDOWN:     PAGEDOWN,
	common.RIGHT:        RIGHT,
	common.LEFT:         LEFT,
	common.DOWN:         DOWN,
	common.UP:           UP,
	common.CTRL:         CTRLLEFT,
	common.SHIFT:        SHIFTLEFT,
	common.ALT:          ALTLEFT,
	common.OS:           OSLEFT,
}
