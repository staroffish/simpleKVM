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
	httpServer.POST("/mousemove", func(ctx *gin.Context) {
		x, exists := ctx.GetQuery("x")
		if !exists {
			ctx.Data(http.StatusBadRequest, "text/text", []byte("x point dose not exists"))
			return
		}
		y, exists := ctx.GetQuery("y")
		if !exists {
			ctx.Data(http.StatusBadRequest, "text/text", []byte("y point dose not exists"))
			return
		}
		xPoint, err := strconv.ParseUint(x, 10, 16)
		if err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("invalid x point:%v", err)))
			return
		}
		yPoint, err := strconv.ParseUint(y, 10, 16)
		if err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("invalid y point:%v", err)))
			return
		}

		if err := hid.MoveTo(uint16(xPoint), uint16(yPoint)); err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("mouse move error")))
		}
		ctx.Data(http.StatusOK, "text/text", []byte("ok"))
	})
	httpServer.POST("/mousedown", func(ctx *gin.Context) {
		button, exists := ctx.GetQuery("button")
		if !exists {
			ctx.Data(http.StatusBadRequest, "text/text", []byte("button code dose not exists"))
			return
		}
		buttonCode, err := strconv.ParseUint(button, 10, 8)
		if err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("invalid button code:%v", buttonCode)))
			return
		}

		if err := hid.MouseDown(int(buttonCode)); err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("mouse down error")))
		}
		ctx.Data(http.StatusOK, "text/text", []byte("ok"))
	})
	httpServer.POST("/mouseup", func(ctx *gin.Context) {
		button, exists := ctx.GetQuery("button")
		if !exists {
			ctx.Data(http.StatusBadRequest, "text/text", []byte("button code dose not exists"))
			return
		}
		buttonCode, err := strconv.ParseUint(button, 10, 8)
		if err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("invalid button code:%v", buttonCode)))
			return
		}

		if err := hid.MouseUp(int(buttonCode)); err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("mouse up error")))
		}
		ctx.Data(http.StatusOK, "text/text", []byte("ok"))
	})
	httpServer.Static("/static", "static")

	httpServer.Run(addr)
}
