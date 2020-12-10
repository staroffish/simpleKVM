package ch9329

import (
	"fmt"
	"time"

	"github.com/staroffish/simpleKVM/hid/common"
	"github.com/staroffish/simpleKVM/log"
	"github.com/tarm/serial"
)

const (
	ALTTRLSHIFTOSKEYINDEX = 5
	RESPONSESTATUSINDEX   = 5
	KEYEVENTTIMEOUT       = time.Millisecond * 200
	MOUSEBUTTONINDEX      = 6
	XMOUSEMOVEINDEXLOW    = 7
	XMOUSEMOVEINDEXHIGH   = 8
	YMOUSEMOVEINDEXLOW    = 9
	YMOUSEMOVEINDEXHIGH   = 10
	MOUSESCROLLINDEX      = 11
)

type ch9329 struct {
	*common.BaseHid
	device              *serial.Port
	screenX             int
	screenY             int
	pressedKey          map[byte]byte
	keyDownEventCh      chan *ch9329KeyEvent
	keyUpEventCh        chan *ch9329KeyEvent
	mouseMoveEventCh    chan *ch9329MouseEvent
	mouseKeyDownEventCh chan *ch9329MouseEvent
	mouseKeyUpEventCh   chan *ch9329MouseEvent
	mouseScrollEventCh  chan *ch9329MouseEvent
}

type ch9329KeyEvent struct {
	resultCh chan error
	ch9329Key
}

type ch9329Key struct {
	keyCode      byte
	isControlKey bool
}

type ch9329MouseEvent struct {
	resultCh chan error
	*ch9329MouseMoveEvent
}

type ch9329MouseMoveEvent struct {
	xPoint uint16
	yPoint uint16
	button int
	scroll int
}

func NewCh9329(x, y int) *ch9329 {
	dev := &ch9329{
		BaseHid: &common.BaseHid{},
		screenX: x,
		screenY: y,
	}
	dev.pressedKey = make(map[byte]byte)
	dev.keyDownEventCh = make(chan *ch9329KeyEvent)
	dev.keyUpEventCh = make(chan *ch9329KeyEvent)
	dev.mouseMoveEventCh = make(chan *ch9329MouseEvent)
	dev.mouseKeyDownEventCh = make(chan *ch9329MouseEvent)
	dev.mouseKeyUpEventCh = make(chan *ch9329MouseEvent)
	dev.mouseScrollEventCh = make(chan *ch9329MouseEvent)
	return dev
}

func (c *ch9329) OpenDevice(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("device path is empty.")
	}

	devicePath := args[0]
	port, err := serial.OpenPort(&serial.Config{Name: devicePath, Baud: 9600})
	if err != nil {
		return err
	}
	c.device = port
	log.PrintInfo("opened hid device %s", devicePath)
	go c.keyDownKeyUp()
	go c.mouseOperator()
	return nil
}

func (c *ch9329) CloseDevice() error {
	c.device.Close()
	return nil
}

func (c *ch9329) GetModelName() string {
	return "ch9329"
}

func (c *ch9329) KeyDown(eventKeyCode byte) error {

	log.PrintDebug("get event key down %v", eventKeyCode)
	isControlKey := common.IsAltCtrlShiftOsKey(eventKeyCode)
	devkeyCode, err := c.KeyBoardCodeToDeviceCode(eventKeyCode)
	if err != nil {
		return fmt.Errorf("invalid key code:%v", err)
	}
	downEvent := &ch9329KeyEvent{
		ch9329Key: ch9329Key{keyCode: devkeyCode, isControlKey: isControlKey},
		resultCh:  make(chan error),
	}
	c.keyDownEventCh <- downEvent
	err = <-downEvent.resultCh

	return err
}

func (c *ch9329) KeyUp(eventKeyCode byte) error {
	log.PrintDebug("get event key up %v", eventKeyCode)
	isControlKey := common.IsAltCtrlShiftOsKey(eventKeyCode)
	devkeyCode, err := c.KeyBoardCodeToDeviceCode(eventKeyCode)
	if err != nil {
		return fmt.Errorf("invalid key code:%v", err)
	}
	upEvent := &ch9329KeyEvent{
		ch9329Key: ch9329Key{keyCode: devkeyCode, isControlKey: isControlKey},
		resultCh:  make(chan error),
	}

	c.keyUpEventCh <- upEvent
	err = <-upEvent.resultCh

	return err
}

func (c *ch9329) MoveTo(x uint16, y uint16) error {
	moveEvent := &ch9329MouseEvent{
		resultCh: make(chan error),
		ch9329MouseMoveEvent: &ch9329MouseMoveEvent{
			xPoint: x,
			yPoint: y,
		},
	}

	var err error
	select {
	case c.mouseMoveEventCh <- moveEvent:
		err = <-moveEvent.resultCh
	default:
	}

	return err
}

func (c *ch9329) MouseDown(button int) error {
	mouseEvent := &ch9329MouseEvent{
		resultCh:             make(chan error),
		ch9329MouseMoveEvent: &ch9329MouseMoveEvent{button: button},
	}

	c.mouseKeyDownEventCh <- mouseEvent
	return <-mouseEvent.resultCh
}

func (c *ch9329) MouseUp(button int) error {
	mouseEvent := &ch9329MouseEvent{
		resultCh:             make(chan error),
		ch9329MouseMoveEvent: &ch9329MouseMoveEvent{button: button},
	}

	c.mouseKeyUpEventCh <- mouseEvent
	return <-mouseEvent.resultCh
}

func (c *ch9329) MouseScroll(scroll int) error {
	mouseEvent := &ch9329MouseEvent{
		resultCh:             make(chan error),
		ch9329MouseMoveEvent: &ch9329MouseMoveEvent{scroll: scroll},
	}

	c.mouseScrollEventCh <- mouseEvent
	return <-mouseEvent.resultCh
}

func (c *ch9329) keyDownKeyUp() {
	keyCommandIndexMap := map[byte]uint{}
	emptyIndexMap := map[uint]uint{7: 7, 8: 8, 9: 9, 10: 10, 11: 11, 12: 12}
	keyDownTimeoutMap := map[ch9329Key]time.Time{}

	command := []byte{0x57, 0xAB, 0x00, 0x02, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var resultCh chan error
	for {
		select {
		case key := <-c.keyDownEventCh:
			if key.isControlKey {
				log.PrintDebug("ctrl/alt/shift/win key down %v", key)
				if command[ALTTRLSHIFTOSKEYINDEX]&key.keyCode > 0 {
					keyDownTimeoutMap[key.ch9329Key] = time.Now()
					key.resultCh <- nil
					continue
				}
				command[ALTTRLSHIFTOSKEYINDEX] = command[ALTTRLSHIFTOSKEYINDEX] | key.keyCode
			} else {
				log.PrintDebug("normal key down %v", key)
				if _, ok := keyCommandIndexMap[key.keyCode]; ok {
					keyDownTimeoutMap[key.ch9329Key] = time.Now()
					key.resultCh <- nil
					continue
				}
				if len(emptyIndexMap) == 0 {
					key.resultCh <- nil
					continue
				}
				for index := range emptyIndexMap {
					command[index] = key.keyCode
					delete(emptyIndexMap, index)
					log.PrintDebug("deleted empty map, index:%d", index)
					keyCommandIndexMap[key.keyCode] = index
					break
				}
			}
			keyDownTimeoutMap[key.ch9329Key] = time.Now()
			resultCh = key.resultCh

		case key := <-c.keyUpEventCh:
			if key.isControlKey {
				log.PrintDebug("ctrl/alt/shift/win key up %v", key)
				if command[ALTTRLSHIFTOSKEYINDEX]&key.keyCode == 0 {
					key.resultCh <- nil
					continue
				}
				command[ALTTRLSHIFTOSKEYINDEX] = command[ALTTRLSHIFTOSKEYINDEX] & ^key.keyCode
			} else {
				index, ok := keyCommandIndexMap[key.keyCode]
				if !ok {
					key.resultCh <- nil
					continue
				}
				command[index] = 0x00
				delete(keyCommandIndexMap, key.keyCode)
				emptyIndexMap[index] = 0
				log.PrintDebug("setted empty map, index:%d", index)
			}
			delete(keyDownTimeoutMap, key.ch9329Key)
			resultCh = key.resultCh
		case <-time.Tick(100 * time.Millisecond):
			now := time.Now()
			for key, keyDownTime := range keyDownTimeoutMap {
				log.PrintDebug("keydown time %v, key=%v", now.Sub(keyDownTime).Milliseconds(), key)
				if now.Sub(keyDownTime) > KEYEVENTTIMEOUT {
					log.PrintInfo("key down time out. downtime=%v key=%v", keyDownTime.Format("2006/01/02 15:04:05.000"), key)
					go func() {
						upEvent := &ch9329KeyEvent{
							ch9329Key: key,
							resultCh:  make(chan error),
						}
						c.keyUpEventCh <- upEvent
						<-upEvent.resultCh
					}()
				}
			}
			continue
		}

		resultCh <- c.readWriteDevice(command)
	}
}

func (c *ch9329) mouseOperator() {
	command := []byte{0x57, 0xAB, 0x00, 0x04, 0x07, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var resultCh chan error

	for {
		command[MOUSESCROLLINDEX] = 0
		select {
		case event := <-c.mouseMoveEventCh:
			resultCh = event.resultCh
			xCur := c.getDevicePoint(event.xPoint, uint16(c.screenX))

			command[XMOUSEMOVEINDEXLOW] = xCur[0]
			command[XMOUSEMOVEINDEXHIGH] = xCur[1]
			yCur := c.getDevicePoint(event.yPoint, uint16(c.screenY))
			command[YMOUSEMOVEINDEXLOW] = yCur[0]
			command[YMOUSEMOVEINDEXHIGH] = yCur[1]
		case event := <-c.mouseKeyDownEventCh:
			resultCh = event.resultCh
			command[MOUSEBUTTONINDEX] |= byte(mouseButtonMap[event.button])
		case event := <-c.mouseKeyUpEventCh:
			resultCh = event.resultCh
			command[MOUSEBUTTONINDEX] &= byte(^mouseButtonMap[event.button])
		case event := <-c.mouseScrollEventCh:
			resultCh = event.resultCh
			log.PrintDebug("scroll=%v", event.scroll)
			if event.scroll < 0 {
				command[MOUSESCROLLINDEX] = 0x01
			} else {
				command[MOUSESCROLLINDEX] = 0xFF
			}
		}

		resultCh <- c.readWriteDevice(command)
	}
}

func (c *ch9329) readWriteDevice(command []byte) error {
	sum := byte(0)
	sumIndex := len(command) - 1
	command[sumIndex] = 0
	for _, b := range command {
		sum += b
	}

	log.PrintDebug("sum = %v", sum)
	command[sumIndex] = sum

	log.PrintDebug("commond=%v", command)
	n, err := c.device.Write(command)
	if err != nil {
		return fmt.Errorf("write file error:n=%d, err=%v", n, err)
	}

	response := make([]byte, 100)
	_, err = c.device.Read(response)
	if err != nil || response[RESPONSESTATUSINDEX] != 0 {
		return fmt.Errorf("read file error: err=%v, status code = %x", err, response[RESPONSESTATUSINDEX])
	}

	return nil
}

func (c *ch9329) getDevicePoint(point, screenMax uint16) [2]byte {
	devicePoint := [2]byte{}
	curse := (uint32(point) * 4096) / uint32(screenMax)
	devicePoint[0] = byte(uint16(curse) & 0x00FF)
	devicePoint[1] = byte(uint16(curse) >> 8)
	log.PrintDebug("point=%d max=%d curse=%d devicePoint=%v", point, screenMax, curse, devicePoint)
	return devicePoint
}
