package capture

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	// "time"

	"github.com/staroffish/simpleKVM/capture/v4l2"
)

func TestCapture(t *testing.T) {
	file, err := os.OpenFile("/dev/video0", os.O_RDWR, 0)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("open /dev/video0 error: %v", err)
		}
		return
	}
	defer file.Close()

	fd := file.Fd()

	dev, err := NewV4l2Device(fd, v4l2.V4L2_PIX_FMT_MJPEG, 24, 1920, 1080, 3)
	if err != nil {
		t.Fatalf("Init return error: %v", err)
	}

	var data []byte

	err = dev.StartStreaming()
	if err != nil {
		t.Fatalf("start streaming error:%v", err)
	}

	data = dev.GetFrame()
	fmt.Printf("data len= %d", len(data))
	if err := dev.StopStreaming(); err != nil {
		t.Fatalf("stop streaming error:%v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("get home dir error:%v", err)
	}

	filePath := fmt.Sprintf("%s/test.jpeg", homeDir)
	if err := ioutil.WriteFile(filePath, data, 0777); err != nil {
		t.Fatalf("WriteFile error:%v", err)
	}

	_, err = ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile error:%v", err)
	}

}
