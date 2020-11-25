package streamer

import (
	"fmt"
	"io"

	"github.com/staroffish/simpleKVM/capture"
)

const (
	HTTP_STREAMER = iota
)

type Streamer interface {
	Streaming(io.Writer) error
}

type baseStreamer struct {
	dev *capture.CaptureDevice
}

func (b *baseStreamer) Streaming() error { return nil }

type MjpegStreamer struct {
	*baseStreamer
	boundary string
}

func NewMjpegStreamer(dev *capture.CaptureDevice, boundary string) *MjpegStreamer {
	return &MjpegStreamer{
		baseStreamer: &baseStreamer{dev},
		boundary:     boundary,
	}
}

func (m *MjpegStreamer) Streaming(writer io.Writer) error {
	subPartHeader := fmt.Sprintf("--%s\nContent-Type: image/jpeg\r\n", m.boundary)
	for {
		_, err := writer.Write([]byte(subPartHeader))
		if err != nil {
			return fmt.Errorf("write subPartHeader error: header=%s, error=%v", subPartHeader, err)
		}

		data := m.dev.GetFrame()

		writeTotal := len(data)

		_, err = writer.Write([]byte(fmt.Sprintf("Content-length: %d\r\n\r\n", writeTotal)))
		if err != nil {
			return fmt.Errorf("write subPartHeader error: header=%s, error=%v", subPartHeader, err)
		}
		for writeTotal != 0 {
			n, err := writer.Write(data)
			if err != nil {
				return fmt.Errorf("write mjpeg error: error=%v", err)
			}
			writeTotal -= n
		}
		writer.Write([]byte("\r\n"))
	}

}
