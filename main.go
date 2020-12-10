package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/staroffish/simpleKVM/capture"
	"github.com/staroffish/simpleKVM/hid"
	"github.com/staroffish/simpleKVM/log"
)

var (
	captureDevice string
	frameFormat   string
	frameRate     int
	height        int
	width         int
	hidDevice     string
	hidModel      string
	logLevel      int
	logFilePath   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "simplekvm",
		Short: "a simple kvm via http",
		Run:   run,
	}
	rootCmd.PersistentFlags().StringVarP(&captureDevice, "capture", "", "/dev/video0", "The path of capture device")
	rootCmd.PersistentFlags().StringVarP(&hidDevice, "hid", "", "/dev/ttyUSB0", "The path of hid device")
	rootCmd.PersistentFlags().StringVarP(&hidModel, "model", "", "ch9329", "The hid device model. supported: ch9329")
	rootCmd.PersistentFlags().StringVarP(&frameFormat, "format", "f", "mjpeg", "The frame format. supported: mjpeg")
	rootCmd.PersistentFlags().IntVarP(&frameRate, "rate", "r", 24, "The frame rate")
	rootCmd.PersistentFlags().IntVarP(&width, "width", "", 1920, "screen width")
	rootCmd.PersistentFlags().IntVarP(&height, "height", "", 1080, "screen height")
	rootCmd.PersistentFlags().IntVarP(&logLevel, "log_level", "", log.LOG_LEVEL_INFO, "log level:\n0: DEBUG\n1: INFO")
	rootCmd.PersistentFlags().StringVarP(&logFilePath, "log_file", "l", "", "log file path. if not present, log will output to stdout")
	rootCmd.Execute()

}

func run(_ *cobra.Command, _ []string) {

	logOutput := os.Stdout
	if logFilePath != "" {
		var err error
		logOutput, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("open log file error:%v\n", err)
			os.Exit(-1)
		}
	}

	log.InitLog(logLevel, logOutput)

	format, err := capture.GetFrameFormatCodeByString(frameFormat)
	if err != nil {
		log.PrintInfo("%v", err)
		os.Exit(-1)
	}

	file, err := os.OpenFile(captureDevice, os.O_RDWR, 0)
	if err != nil {
		log.PrintInfo("open %s error: %v", captureDevice, err)
		os.Exit(-1)
	}
	defer file.Close()

	fd := file.Fd()

	dev, err := capture.NewV4l2Device(fd, format, uint32(frameRate), uint32(width), uint32(height), 3)
	if err != nil {
		log.PrintInfo("Init return error: %v", err)
		return
	}

	err = dev.StartStreaming()
	if err != nil {
		log.PrintInfo("start streaming error:%v", err)
		return
	}

	hid.InitHid(width, height)

	hidDev, err := hid.GetHidDevice(hidModel)
	if err != nil {
		log.PrintInfo("get hid device error:%v", err)
		return
	}

	if err := hidDev.OpenDevice(hidDevice); err != nil {
		log.PrintInfo("open hid device error:%v", err)
		return
	}

	StartHttpServer(context.Background(), "0.0.0.0:8181", dev, hidDev)
}
