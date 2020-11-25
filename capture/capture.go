package capture

import (
	"context"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/staroffish/simpleKVM/capture/v4l2"
)

type ImageBuffer struct {
	length  uint32
	pointer unsafe.Pointer
	data    []byte
}

type CaptureDevice struct {
	fd            uintptr
	capability    *v4l2.V4l2Capability
	buffer        []ImageBuffer
	streamingType uint32
	cancel        context.CancelFunc
	frameCh       chan []byte
}

var frameFormatNameToCode = map[string]uint32{
	"mjpeg": v4l2.V4L2_PIX_FMT_MJPEG,
}

func GetFrameFormatCodeByString(name string) (uint32, error) {
	code, ok := frameFormatNameToCode[name]
	if ok {
		return code, nil
	}
	return 0, fmt.Errorf("%s not supported", name)
}

func NewV4l2Device(fileDescription uintptr) *CaptureDevice {
	return &CaptureDevice{
		fd: fileDescription,
	}
}

func (c *CaptureDevice) Init(imageDataFormat, frameRate, Width, Height, bufferCount uint32) (err error) {
	capability, err := v4l2.QueryCapability(c.fd)
	if err != nil {
		return fmt.Errorf("QueryCapability error: %v", err)
	}
	if capability.Capabilities&uint32(v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE) <= 0 {
		return fmt.Errorf("device does not support V4L2_BUF_TYPE_VIDEO_CAPTURE")
	}
	if capability.Capabilities&uint32(v4l2.V4L2_CAP_DEVICE_CAPS) > 0 && capability.DeviceCaps&uint32(v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE) <= 0 {
		return fmt.Errorf("device does not support V4L2_BUF_TYPE_VIDEO_CAPTURE")
	}
	c.capability = capability
	c.streamingType = v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE

	var index uint32 = 0
	formatDescs := map[uint32]*v4l2.V4l2Fmtdesc{}
	for {
		format, err := v4l2.EnumFormat(c.fd, v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE, index)
		if err != nil {
			if len(formatDescs) == 0 {
				return fmt.Errorf("QueryFormat error: err=%v", err)
			}
			break
		}
		fmt.Printf("Format=%s\n", format.Description)
		formatDescs[format.Pixelformat] = format
		index++
	}

	formatDesc, ok := formatDescs[imageDataFormat]
	if !ok {
		return fmt.Errorf("unsupported image format\n")
	}

	fmt.Printf("use format %s\n", formatDesc.Description)
	format := &v4l2.V4l2Format{
		Type: c.streamingType,
		Pix: &v4l2.V4l2PixFormat{
			Width:       Width,
			Height:      Height,
			Field:       v4l2.V4L2_FIELD_ANY,
			Pixelformat: formatDesc.Pixelformat,
		},
	}
	if err := v4l2.SetFrameFormat(c.fd, format); err != nil {
		return fmt.Errorf("SetFrameSize error: %v", err)
	}

	reqBuff := &v4l2.V4l2RequestBuffers{
		Count:  bufferCount,
		Type:   c.streamingType,
		Memory: v4l2.V4L2_MEMORY_MMAP,
	}

	streamParam, err := v4l2.GetStreamParam(c.fd, c.streamingType)
	if err != nil {
		return fmt.Errorf("GetStreamParam error: %v", err)
	}

	streamParam.CaptureParam.Fract.Numerator = 1
	streamParam.CaptureParam.Fract.Denominator = frameRate

	if err := v4l2.SetStreamParam(c.fd, streamParam); err != nil {
		return fmt.Errorf("SetStreamParam error: %v", err)
	}

	fmt.Printf("type=%v CaptureParam=%v Fract=%v\n", c.streamingType, streamParam.CaptureParam, streamParam.CaptureParam.Fract)

	if err := v4l2.RequestBuffer(c.fd, reqBuff); err != nil {
		return err
	}
	c.buffer = make([]ImageBuffer, bufferCount)

	for index := uint32(0); index < reqBuff.Count; index++ {
		v4l2Buffer, err := v4l2.QueryBuffer(c.fd, c.streamingType, index)
		if err != nil {
			return err
		}

		c.buffer[index].data, err = syscall.Mmap(int(c.fd), int64(v4l2Buffer.Offset), int(v4l2Buffer.Length), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			return fmt.Errorf("mmap error: %v", err)
		}
		defer func(index uint32) {
			if err != nil {
				unMaperr := syscall.Munmap(c.buffer[index].data)
				if unMaperr != nil {
					err = fmt.Errorf("%v|munmap error: %v", err, unMaperr)
				}
			}
		}(index)

		if err = v4l2.QueueBuffer(c.fd, c.streamingType, index); err != nil {
			return err
		}
	}

	return nil
}

func (c *CaptureDevice) StartStreaming() (err error) {
	ctx, cancel := context.WithCancel(context.Background())

	c.cancel = cancel

	if err := v4l2.StreamOn(c.fd, c.streamingType); err != nil {
		return err
	}
	c.frameCh = make(chan []byte)

	go func(ctx context.Context) {
		defer func() {
			if err := v4l2.StreamOff(c.fd, c.streamingType); err != nil {
				fmt.Printf("%v\n", err)
			}
		}()
		for {
			buffer, err := v4l2.DeQueueBuffer(c.fd, c.streamingType)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}

			data := make([]byte, buffer.BytesUsed)
			copy(data, c.buffer[buffer.Index].data[:buffer.BytesUsed])

			if err = v4l2.QueueBuffer(c.fd, c.streamingType, buffer.Index); err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			select {
			case c.frameCh <- data:
			case <-ctx.Done():
				return
			default:
			}
		}
	}(ctx)

	return nil
}

func (c *CaptureDevice) StopStreaming() error {
	c.cancel()
	var errMsg string

	for index, buffer := range c.buffer {
		err := syscall.Munmap(buffer.data)
		if err != nil {
			errMsg = fmt.Sprintf("%smunmap error. index=%d, error=%v\n", errMsg, index, err)
		}
	}
	if errMsg != "" {
		return fmt.Errorf("%v", errMsg)
	}

	return nil
}

func (c *CaptureDevice) GetFrame() []byte {
	return <-c.frameCh
}
