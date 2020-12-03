package ch9329

import (
	"fmt"
	"time"

	"github.com/staroffish/simpleKVM/hid/common"
	"github.com/tarm/serial"
)

const (
	ALTTRLSHIFTOSKEYINDEX = 5
	RESPONSESTATUSINDEX   = 5
	SUMINDEX              = 13
	KEYEVENTTIMEOUT       = time.Millisecond * 200
)

type ch9329 struct {
	*common.BaseHid
	device         *serial.Port
	pressedKey     map[byte]byte
	keyDownEventCh chan *ch9329KeyEvent
	keyUpEventCh   chan *ch9329KeyEvent
}

type ch9329KeyEvent struct {
	resultCh chan error
	ch9329Key
}

type ch9329Key struct {
	keyCode      byte
	isControlKey bool
}

func NewCh9329() *ch9329 {
	dev := &ch9329{
		BaseHid: &common.BaseHid{},
	}
	dev.pressedKey = make(map[byte]byte)
	dev.keyDownEventCh = make(chan *ch9329KeyEvent)
	dev.keyUpEventCh = make(chan *ch9329KeyEvent)
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
	fmt.Printf("opened device\n")
	go c.keyDownKeyUp()
	return nil
}
func (c *ch9329) CloseDevice() error {
	c.device.Close()
	return nil
}

func (c *ch9329) MouseDown(button int) error {
	var x uint16
	for x = 1; x <= 2560; x = x + 128 {
		fmt.Printf("x=%d byte1=%x byte2=%x\n", x, byte((x<<8)>>8), byte(x>>8))
		hidData := []byte{0x57, 0xAB, 0x00, 0x04, 0x07, 0x02, 0x00, byte((x << 8) >> 8), byte(x >> 8), 0x01, 0x00, 0x00}
		var totalData uint16
		for _, data := range hidData {
			totalData += uint16(data)
		}
		hidData = append(hidData, byte((totalData<<8)>>8))
		n, err := c.device.Write(hidData)
		if err != nil {
			return fmt.Errorf("write file error:n=%d, err=%v", n, err)
		}
		response := make([]byte, 100)
		n, err = c.device.Read(response)
		if err != nil {
			return fmt.Errorf("read file error: err=%v", err)
		}

	}
	return nil
}

func (c *ch9329) GetModelName() string {
	return "ch9329"
}

func (c *ch9329) KeyDown(eventKeyCode byte) error {

	fmt.Printf("get event key down %v\n", eventKeyCode)
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

func (c *ch9329) keyDownKeyUp() {
	keyCommandIndexMap := map[byte]uint{}
	emptyIndexMap := map[uint]uint{7: 7, 8: 8, 9: 9, 10: 10, 11: 11, 12: 12}
	keyDownTimeoutMap := map[ch9329Key]time.Time{}

	command := []byte{0x57, 0xAB, 0x00, 0x02, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for {
		var resultCh chan error
		command[SUMINDEX] = 0
		select {
		case key := <-c.keyDownEventCh:
			if key.isControlKey {
				if command[ALTTRLSHIFTOSKEYINDEX]&key.keyCode > 0 {
					fmt.Printf("get ctrl key down %v\n", key)
					keyDownTimeoutMap[key.ch9329Key] = time.Now()
					key.resultCh <- nil
					continue
				}
				command[ALTTRLSHIFTOSKEYINDEX] = command[ALTTRLSHIFTOSKEYINDEX] | key.keyCode
			} else {
				fmt.Printf("get normal key down %v\n", key)
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
					fmt.Printf("index=%d\n", index)
					command[index] = key.keyCode
					delete(emptyIndexMap, index)
					keyCommandIndexMap[key.keyCode] = index
					break
				}
			}
			keyDownTimeoutMap[key.ch9329Key] = time.Now()
			resultCh = key.resultCh

		case key := <-c.keyUpEventCh:
			fmt.Printf("get key up %v\n", key)
			if key.isControlKey {
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
				fmt.Printf("index=%d\n", index)
				command[index] = 0x00
				delete(keyCommandIndexMap, key.keyCode)
				emptyIndexMap[index] = 0
			}
			delete(keyDownTimeoutMap, key.ch9329Key)
			resultCh = key.resultCh
		case <-time.Tick(100 * time.Millisecond):
			now := time.Now()
			for key, keyDownTime := range keyDownTimeoutMap {
				fmt.Printf("now sub keydown time %v, key=%v command=%v \n", now.Sub(keyDownTime), key, command)
				if now.Sub(keyDownTime) > KEYEVENTTIMEOUT {
					fmt.Printf("key time out %v:%v\n", keyDownTime, key)
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

		sum := byte(0)
		for _, b := range command {
			sum += b
		}

		fmt.Printf("sum = %v\n", sum)
		command[SUMINDEX] = sum

		fmt.Printf("commond=%v\n", command)
		n, err := c.device.Write(command)
		if err != nil {
			resultCh <- fmt.Errorf("write file error:n=%d, err=%v", n, err)
			continue
		}

		response := make([]byte, 100)
		_, err = c.device.Read(response)
		if err != nil || response[RESPONSESTATUSINDEX] != 0 {
			resultCh <- fmt.Errorf("read file error: err=%v, status code = %x", err, response[RESPONSESTATUSINDEX])
			continue
		}
		resultCh <- nil
	}
}
