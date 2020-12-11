package streamer

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/staroffish/simpleKVM/capture"
	"github.com/staroffish/simpleKVM/capture/v4l2"
)

type Streamer interface {
	Streaming(io.Writer) error
	Handler() func(ctx *gin.Context)
	Path() string
	HtmlElement() string
}

type baseStreamer struct {
	dev    *capture.CaptureDevice
	height int
	width  int
}

func (b *baseStreamer) Streaming() error                { return nil }
func (b *baseStreamer) Handler() func(ctx *gin.Context) { return nil }
func (b *baseStreamer) Path() string                    { return "" }
func (b *baseStreamer) HtmlElement() string             { return "" }

type MjpegStreamer struct {
	*baseStreamer
	boundary string
}

func NewStreamer(dev *capture.CaptureDevice) (Streamer, error) {
	var streamer Streamer
	switch dev.GetFormat() {
	case v4l2.V4L2_PIX_FMT_MJPEG:
		return NewMjpegStreamer(dev, "frame", int(dev.GetWidth()), int(dev.GetHeight())), nil
	default:
		return nil, fmt.Errorf("unsupported streamer type")
	}
	return streamer, nil
}

func NewMjpegStreamer(dev *capture.CaptureDevice, boundary string, width int, height int) *MjpegStreamer {
	return &MjpegStreamer{
		baseStreamer: &baseStreamer{
			dev:    dev,
			height: height,
			width:  width,
		},
		boundary: boundary,
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

func (m *MjpegStreamer) Handler() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", m.boundary))
		ctx.Writer.WriteHeader(http.StatusOK)
		m.Streaming(ctx.Writer)
	}
}

func (m *MjpegStreamer) Path() string {
	return "/mjpeg"
}

func (m *MjpegStreamer) HtmlElement() string {
	return fmt.Sprintf(`<img src="%s" onmousemove="mouseMove(event)" onmousedown="return mouseDown(event)" onmouseup="return mouseUp(event)" onwheel="return mouseScroll(event)" width="%d" height="%d" />`,
		m.Path(),
		m.width,
		m.height,
	)
}
