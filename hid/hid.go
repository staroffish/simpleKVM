package hid

import (
	"fmt"

	"github.com/staroffish/simpleKVM/hid/ch9329"
	"github.com/staroffish/simpleKVM/hid/common"
)

var deviceMap map[string]common.Hid

func InitHid(x, y int) {
	deviceMap = make(map[string]common.Hid)
	ch9329Dev := ch9329.NewCh9329(x, y)
	deviceMap[ch9329Dev.GetModelName()] = ch9329Dev
}

func GetHidDevice(model string) (common.Hid, error) {
	dev, ok := deviceMap[model]
	if !ok {
		return nil, fmt.Errorf("unknown device %s", model)
	}
	return dev, nil
}
