package hid

import (
	"fmt"

	"github.com/staroffish/simpleKVM/hid/ch9329"
	"github.com/staroffish/simpleKVM/hid/common"
)

const (
	MouseLeftButton = iota
	MouseRightButton
	MouseMiddleButton
)

var deviceMap map[string]common.Hid

func init() {
	deviceMap = make(map[string]common.Hid)
	ch9329Dev := ch9329.NewCh9329()
	deviceMap[ch9329Dev.GetModelName()] = ch9329Dev
}

func GetHidDevice(model string) (common.Hid, error) {
	dev, ok := deviceMap[model]
	if !ok {
		return nil, fmt.Errorf("unknown device %s", model)
	}
	return dev, nil
}
