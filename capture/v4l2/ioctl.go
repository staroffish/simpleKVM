package v4l2

/*
#include <linux/videodev2.h>
#include <stdlib.h>
#include <sys/ioctl.h>
#include <sys/mman.h>

int ioctl3arg(int fd, unsigned long request, void *arg)
{
    return ioctl(fd, request, arg);
}

*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	V4L2_BUF_TYPE_VIDEO_CAPTURE = C.V4L2_BUF_TYPE_VIDEO_CAPTURE

	V4L2_CAP_DEVICE_CAPS     = C.V4L2_CAP_DEVICE_CAPS
	V4L2_PIX_FMT_MJPEG       = C.V4L2_PIX_FMT_MJPEG
	V4L2_PIX_FMT_YUYV        = C.V4L2_PIX_FMT_YUYV
	V4L2_FIELD_NONE          = C.V4L2_FIELD_NONE
	V4L2_FIELD_ANY           = C.V4L2_FIELD_ANY
	V4L2_FIELD_INTERLACED_TB = C.V4L2_FIELD_INTERLACED_TB
	V4L2_MEMORY_MMAP         = C.V4L2_MEMORY_MMAP
)

type V4l2Capability struct {
	Driver       string
	Card         string
	BusInfo      string
	Version      uint32
	Capabilities uint32
	DeviceCaps   uint32
}

type V4l2Fmtdesc struct {
	Index       uint32 /* Format number      */
	Type        uint32 /* enum v4l2_buf_type */
	Flags       uint32
	Description string /* Description string */
	Pixelformat uint32 /* Format fourcc      */
}

type V4l2Format struct {
	Type uint32
	Pix  *V4l2PixFormat
}

type V4l2PixFormat struct {
	Width        uint32
	Height       uint32
	Pixelformat  uint32
	Field        uint32 /* enum v4l2_field */
	Bytesperline uint32 /* for padding, zero if unused */
	Sizeimage    uint32
	Colorspace   uint32 /* enum v4l2_colorspace */
	Priv         uint32 /* private data, depends on pixelformat */
	Flags        uint32 /* format flags (V4L2_PIX_FMT_FLAG_*) */
	Uycbcr_enc   uint32 /* enum v4l2_ycbcr_encoding */
	Quantization uint32 /* enum v4l2_quantization */
	Xfer_func    uint32 /* enum v4l2_xfer_func */
}

type V4l2StreamParam struct {
	Type         uint32
	CaptureParam *V4l2CaptureParam
}

type V4l2CaptureParam struct {
	Capability   uint32
	Capturemode  uint32
	Extendedmode uint32
	Readbuffers  uint32
	Fract        *V4l2Fract
}

type V4l2Fract struct {
	Numerator   uint32
	Denominator uint32
}

type V4l2RequestBuffers struct {
	Count    uint32
	Type     uint32
	Memory   uint32
	Reserved [2]uint32
}

type V4l2Buffer struct {
	Index     uint32
	Type      uint32
	BytesUsed uint32
	Flags     uint32
	Field     uint32
	Sequence  uint32
	Memory    uint32
	Offset    uint32
	Length    uint32
	Reserved2 uint32
	Reserved  uint32
}

func QueryCapability(fd uintptr) (*V4l2Capability, error) {
	capability := &V4l2Capability{}
	cCap := &C.struct_v4l2_capability{}

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_QUERYCAP, unsafe.Pointer(cCap))
	if ret < 0 {
		return nil, fmt.Errorf("call VIDIOC_QUERYCAP error, ret=%d, err=%v", ret, err)
	}
	capability.BusInfo = fmt.Sprintf("%s", cCap.bus_info)
	capability.Driver = fmt.Sprintf("%s", cCap.driver)
	capability.Card = fmt.Sprintf("%s", cCap.card)
	capability.Capabilities = uint32(cCap.capabilities)
	capability.DeviceCaps = uint32(cCap.device_caps)
	return capability, nil
}

func EnumFormat(fd uintptr, typ int, index uint32) (*V4l2Fmtdesc, error) {
	cDesc := &C.struct_v4l2_fmtdesc{}
	cDesc.index = C.__u32(index)
	cDesc._type = C.__u32(typ)
	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_ENUM_FMT, unsafe.Pointer(cDesc))
	if ret < 0 {
		return nil, fmt.Errorf("call VIDIOC_ENUM_FMT error, ret=%d, err=%v", ret, err)
	}
	desc := &V4l2Fmtdesc{}
	desc.Type = uint32(typ)
	desc.Index = index
	desc.Flags = uint32(cDesc.flags)
	desc.Description = fmt.Sprintf("%s", cDesc.description)
	desc.Pixelformat = uint32(cDesc.pixelformat)
	return desc, nil
}

func GetStandard(fd uintptr) error {
	index := 0
	for {
		standard := C.struct_v4l2_standard{}
		standard.index = C.__u32(index)
		ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_ENUMSTD, unsafe.Pointer(&standard))
		if ret < 0 {
			return fmt.Errorf("call VIDIOC_ENUMSTD  error: %v", err)
		}
		fmt.Printf("index=%d, id=%d, name=%s, Numerator=%d, Denominator=%d", index, standard.id, standard.name, standard.frameperiod.numerator, standard.frameperiod.denominator)
		index++
	}
}

func GetFrameFormat(fd uintptr, typ uint32) (*V4l2Format, error) {
	cFormat := &C.struct_v4l2_format{}
	cFormat._type = C.__u32(typ)
	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_G_FMT, unsafe.Pointer(cFormat))
	if ret < 0 {
		return nil, fmt.Errorf("call VIDIOC_G_FMT error, ret=%d, err=%v", ret, err)
	}
	pix := (*C.struct_v4l2_pix_format)(unsafe.Pointer(&cFormat.fmt[0]))

	format := &V4l2Format{Pix: &V4l2PixFormat{}}

	format.Type = typ
	format.Pix.Width = uint32(pix.width)
	format.Pix.Height = uint32(pix.height)
	format.Pix.Pixelformat = uint32(pix.pixelformat)
	format.Pix.Field = uint32(pix.field)
	format.Pix.Bytesperline = uint32(pix.bytesperline)
	format.Pix.Sizeimage = uint32(pix.sizeimage)
	format.Pix.Colorspace = uint32(pix.colorspace)
	format.Pix.Priv = uint32(pix.priv)
	format.Pix.Flags = uint32(pix.flags)
	// format.Pix.Uycbcr_enc = uint32(pix.ycbcr_enc)
	format.Pix.Quantization = uint32(pix.quantization)
	format.Pix.Xfer_func = uint32(pix.xfer_func)

	fmt.Printf("%v\n", format)

	return format, nil
}

func SetFrameFormat(fd uintptr, format *V4l2Format) error {
	cFormat := &C.struct_v4l2_format{}
	cFormat._type = C.__u32(format.Type)

	pix := (*C.struct_v4l2_pix_format)(unsafe.Pointer(&cFormat.fmt[0]))
	pix.width = C.__u32(format.Pix.Width)
	pix.height = C.__u32(format.Pix.Height)
	pix.pixelformat = C.__u32(format.Pix.Pixelformat)
	pix.field = C.__u32(format.Pix.Field)

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_S_FMT, unsafe.Pointer(cFormat))
	if ret < 0 {
		return fmt.Errorf("call VIDIOC_S_FMT error, ret=%d, err=%v", ret, err)
	}

	return nil
}

func GetStreamParam(fd uintptr, typ uint32) (*V4l2StreamParam, error) {
	cParam := &C.struct_v4l2_streamparm{}
	cParam._type = C.__u32(typ)

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_G_PARM, unsafe.Pointer(cParam))
	if ret < 0 {
		return nil, fmt.Errorf("call VIDIOC_G_PARM error, ret=%d, err=%v", ret, err)
	}

	param := &V4l2StreamParam{CaptureParam: &V4l2CaptureParam{Fract: &V4l2Fract{}}}
	param.Type = uint32(cParam._type)
	cCapParam := (*C.struct_v4l2_captureparm)(unsafe.Pointer(&cParam.parm))
	param.CaptureParam.Capability = uint32(cCapParam.capability)
	param.CaptureParam.Capturemode = uint32(cCapParam.capturemode)
	param.CaptureParam.Extendedmode = uint32(cCapParam.extendedmode)
	param.CaptureParam.Readbuffers = uint32(cCapParam.readbuffers)
	param.CaptureParam.Fract.Numerator = uint32(cCapParam.timeperframe.numerator)
	param.CaptureParam.Fract.Denominator = uint32(cCapParam.timeperframe.denominator)

	return param, nil
}

func SetStreamParam(fd uintptr, param *V4l2StreamParam) error {
	cParam := &C.struct_v4l2_streamparm{}
	cParam._type = C.__u32(param.Type)

	cCapParam := (*C.struct_v4l2_captureparm)(unsafe.Pointer(&cParam.parm))
	cCapParam.timeperframe.numerator = C.__u32(param.CaptureParam.Fract.Numerator)
	cCapParam.timeperframe.denominator = C.__u32(param.CaptureParam.Fract.Denominator)

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_S_PARM, unsafe.Pointer(cParam))
	if ret < 0 {
		return fmt.Errorf("call VIDIOC_S_PARM error, ret=%d, err=%v", ret, err)
	}
	return nil
}

func RequestBuffer(fd uintptr, requestBuffer *V4l2RequestBuffers) error {

	cRequestBuffer := &C.struct_v4l2_requestbuffers{}

	cRequestBuffer.count = C.__u32(requestBuffer.Count)
	cRequestBuffer._type = C.__u32(requestBuffer.Type)
	cRequestBuffer.memory = C.__u32(requestBuffer.Memory)

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_REQBUFS, unsafe.Pointer(cRequestBuffer))
	if ret < 0 {
		return fmt.Errorf("call VIDIOC_REQBUFS error, ret=%d, err=%v", ret, err)
	}

	return nil
}

func QueryBuffer(fd uintptr, typ, index uint32) (*V4l2Buffer, error) {

	cBuffer := &C.struct_v4l2_buffer{}
	cBuffer.index = C.__u32(index)
	cBuffer._type = C.__u32(typ)

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_QUERYBUF, unsafe.Pointer(cBuffer))
	if ret < 0 {
		return nil, fmt.Errorf("call VIDIOC_QUERYBUF error, ret=%d, err=%v", ret, err)
	}

	return CBufferToGoBuffer(cBuffer), nil
}

func QueueBuffer(fd uintptr, typ, index uint32) error {

	cBuffer := &C.struct_v4l2_buffer{}
	cBuffer.index = C.__u32(index)
	cBuffer._type = C.__u32(typ)
	cBuffer.memory = C.V4L2_MEMORY_MMAP

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_QBUF, unsafe.Pointer(cBuffer))
	if ret < 0 {
		return fmt.Errorf("call VIDIOC_QBUF error, ret=%d, err=%v", ret, err)
	}

	return nil
}

func DeQueueBuffer(fd uintptr, typ uint32) (*V4l2Buffer, error) {

	cBuffer := &C.struct_v4l2_buffer{}
	cBuffer._type = C.__u32(typ)
	cBuffer.memory = C.V4L2_MEMORY_MMAP

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_DQBUF, unsafe.Pointer(cBuffer))
	if ret < 0 {
		return nil, fmt.Errorf("call VIDIOC_DQBUF error, ret=%d, err=%v", ret, err)
	}

	return CBufferToGoBuffer(cBuffer), nil
}

func StreamOn(fd uintptr, typ uint32) error {
	cRequestBuffer := &C.struct_v4l2_requestbuffers{}
	cRequestBuffer._type = C.__u32(typ)

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_STREAMON, unsafe.Pointer(&cRequestBuffer._type))
	if ret < 0 {
		return fmt.Errorf("call VIDIOC_STREAMON error, ret=%d, err=%v", ret, err)
	}

	return nil
}

func StreamOff(fd uintptr, typ uint32) error {
	cRequestBuffer := &C.struct_v4l2_requestbuffers{}
	cRequestBuffer._type = C.__u32(typ)

	ret, err := C.ioctl3arg(C.int(fd), C.VIDIOC_STREAMOFF, unsafe.Pointer(&cRequestBuffer._type))
	if ret < 0 {
		return fmt.Errorf("call VIDIOC_STREAMOFF error, ret=%d, err=%v", ret, err)
	}

	return nil
}

func CBufferToGoBuffer(cBuffer *C.struct_v4l2_buffer) *V4l2Buffer {

	buffer := &V4l2Buffer{}
	buffer.Index = uint32(cBuffer.index)
	buffer.Type = uint32(cBuffer._type)
	buffer.BytesUsed = uint32(cBuffer.bytesused)
	buffer.Flags = uint32(cBuffer.flags)
	buffer.Field = uint32(cBuffer.field)
	buffer.Sequence = uint32(cBuffer.sequence)
	buffer.Memory = uint32(cBuffer.memory)
	buffer.Length = uint32(cBuffer.length)

	offset := (*C.__u32)(unsafe.Pointer(&cBuffer.m))
	buffer.Offset = uint32(*offset)

	return buffer
}
