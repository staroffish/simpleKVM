package ch9329

import (
	"fmt"
	"testing"

	"github.com/staroffish/simpleKVM/hid/common"
)

func TestCh9329KeyDownUP(t *testing.T) {
	hid := NewCh9329()
	fmt.Printf("new ch9329\n")
	if err := hid.OpenDevice("/dev/ttyUSB0"); err != nil {
		t.Fatalf("open device error: %v", err)
	}
	username := []byte{'R', 'O', 'O', 'T', common.ENTER}
	for _, keyCode := range []byte(username) {
		fmt.Printf(" keyCode=%v\n", keyCode)
		if err := hid.KeyDown(keyCode); err != nil {
			t.Fatalf("key down error: %v", err)
		}
		if err := hid.KeyUp(keyCode); err != nil {
			t.Fatalf("key down error: %v", err)
		}
	}

	hid.CloseDevice()
}

func TestKeyBoardEventCodeToCh9329KeyCode(t *testing.T) {
	hid := &ch9329{}
	keyCode, err := hid.KeyBoardCodeToDeviceCode(common.ALPHABETBASE)
	if err != nil {
		t.Fatalf("call KeyBoardEventCodeToCh9329KeyCode('a') error: %v", err)
	}
	if keyCode != 0x04 {
		t.Fatalf("get ch9329 key code error. 'a' must return 0x04, but returned %x", keyCode)
	}

	keyCode, err = hid.KeyBoardCodeToDeviceCode(common.ALPHABETBASE + common.ALPHABETLENGTH - 1)
	if err != nil {
		t.Fatalf("call KeyBoardEventCodeToCh9329KeyCode('z') error: %v", err)
	}
	if keyCode != 0x1d {
		t.Fatalf("get ch9329 key code error. 'z must return 0x1d, but returned %x", keyCode)
	}

	keyCode, err = hid.KeyBoardCodeToDeviceCode(common.NUMBERBASE + 1)
	if err != nil {
		t.Fatalf("call KeyBoardEventCodeToCh9329KeyCode('1') error: %v", err)
	}
	if keyCode != 0x1E {
		t.Fatalf("get ch9329 key code error. '1' must return 0x1E, but returned %x", keyCode)
	}

	keyCode, err = hid.KeyBoardCodeToDeviceCode(common.NUMBERBASE)
	if err != nil {
		t.Fatalf("call KeyBoardEventCodeToCh9329KeyCode('0') error: %v", err)
	}
	if keyCode != 0x27 {
		t.Fatalf("get ch9329 key code error. '0' must return 0x27, but returned %x", keyCode)
	}

	keyCode, err = hid.KeyBoardCodeToDeviceCode(common.NUMBERBASE + common.NUMBERBASELENGHT - 1)
	if err != nil {
		t.Fatalf("call KeyBoardEventCodeToCh9329KeyCode('9') error: %v", err)
	}
	if keyCode != 0x26 {
		t.Fatalf("get ch9329 key code error. '9' must return 0x26, but returned %x", keyCode)
	}

	keyCode, err = hid.KeyBoardCodeToDeviceCode(common.FUNCTIONBASE)
	if err != nil {
		t.Fatalf("call KeyBoardEventCodeToCh9329KeyCode('f1') error: %v", err)
	}
	if keyCode != 0x3A {
		t.Fatalf("get ch9329 key code error. 'f1' must return 0x3A, but returned %x", keyCode)
	}

	keyCode, err = hid.KeyBoardCodeToDeviceCode(common.FUNCTIONBASE + common.FUNCTIONBASELENGTH - 1)
	if err != nil {
		t.Fatalf("call KeyBoardEventCodeToCh9329KeyCode('f12') error: %v", err)
	}
	if keyCode != 0x45 {
		t.Fatalf("get ch9329 key code error. 'f12' must return 0x45, but returned %x", keyCode)
	}

	keyCode, err = hid.KeyBoardCodeToDeviceCode(common.COMMA)
	if err != nil {
		t.Fatalf("call KeyBoardEventCodeToCh9329KeyCode(',') error: %v", err)
	}
	if keyCode != COMMA {
		t.Fatalf("get ch9329 key code error. ',' must return 0x36, but returned %x", keyCode)
	}
}
