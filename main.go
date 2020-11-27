package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/staroffish/simpleKVM/capture"
)

var (
	captureDevice string
	frameFormat   string
	frameRate     int
	height        int
	width         int
	hidDevice     string
	hidModel      string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "simplekvm",
		Short: "a simple kvm via http",
		Run:   run,
	}
	rootCmd.PersistentFlags().StringVarP(&captureDevice, "capture", "", "/dev/video0", "The path of capture device")
	rootCmd.PersistentFlags().StringVarP(&hidDevice, "hid", "", "/dev/ttyUSB0", "The path of hid device")
	rootCmd.PersistentFlags().StringVarP(&hidModel, "model", "", "ch9329", "The hid device model")
	rootCmd.PersistentFlags().StringVarP(&frameFormat, "format", "f", "mjpeg", "The frame format. supported: mjpeg")
	rootCmd.PersistentFlags().IntVarP(&frameRate, "rate", "r", 24, "The frame rate")
	rootCmd.PersistentFlags().IntVarP(&height, "height", "", 1920, "height")
	rootCmd.PersistentFlags().IntVarP(&width, "width", "", 1080, "width")
	rootCmd.Execute()

}

func run(_ *cobra.Command, _ []string) {

	format, err := capture.GetFrameFormatCodeByString(frameFormat)
	if err != nil {
		fmt.Printf("%v", err)
	}

	file, err := os.OpenFile(captureDevice, os.O_RDWR, 0)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("open %s error: %v", captureDevice, err)
		}
		return
	}
	defer file.Close()

	fd := file.Fd()

	dev, err := capture.NewV4l2Device(fd, format, uint32(frameRate), uint32(height), uint32(width), 3)
	if err != nil {
		fmt.Printf("Init return error: %v", err)
		return
	}

	err = dev.StartStreaming()
	if err != nil {
		fmt.Printf("start streaming error:%v", err)
		return
	}
	StartHttpServer(context.Background(), "0.0.0.0:8181", dev)
}