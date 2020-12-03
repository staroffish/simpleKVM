package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/staroffish/simpleKVM/capture"
	"github.com/staroffish/simpleKVM/hid/common"
	"github.com/staroffish/simpleKVM/streamer"
)

func StartHttpServer(ctx context.Context, addr string, dev *capture.CaptureDevice, hid common.Hid) {
	gin.SetMode(gin.ReleaseMode)
	httpServer := gin.New()
	// gin.Default().Use()

	httpServer.GET("/mjpeg", func(ctx *gin.Context) {
		boundary := "frame"
		ctx.Writer.Header().Add("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))
		ctx.Writer.WriteHeader(http.StatusOK)
		mjpegStreamer := streamer.NewMjpegStreamer(dev, boundary)
		mjpegStreamer.Streaming(ctx.Writer)
	})
	httpServer.POST("/keydown", func(ctx *gin.Context) {
		keyCodeStr, exists := ctx.GetQuery("key_code")
		if !exists {
			ctx.Data(http.StatusBadRequest, "text/text", []byte("key code dose not exists"))
			return
		}
		keyCode, err := strconv.ParseUint(keyCodeStr, 10, 8)
		if err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("invalid key code:%v", err)))
			return
		}

		if err := hid.KeyDown(byte(keyCode)); err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("key down error: %v", err)))
			return
		}
		ctx.Data(http.StatusOK, "text/text", []byte("ok"))
	})
	httpServer.POST("/keyup", func(ctx *gin.Context) {
		keyCodeStr, exists := ctx.GetQuery("key_code")
		if !exists {
			ctx.Data(http.StatusBadRequest, "text/text", []byte("key code dose not exists"))
			return
		}
		keyCode, err := strconv.ParseUint(keyCodeStr, 10, 8)
		if err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("invalid key code:%v", err)))
			return
		}
		if err := hid.KeyUp(byte(keyCode)); err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("key up error: %v", err)))
			return
		}
		ctx.Data(http.StatusOK, "text/text", []byte("ok"))
	})
	httpServer.Static("/static", "static")

	httpServer.Run(addr)
}
