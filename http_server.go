package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/staroffish/simpleKVM/capture"
	"github.com/staroffish/simpleKVM/streamer"
)

func StartHttpServer(ctx context.Context, addr string, dev *capture.CaptureDevice) {
	httpServer := gin.Default()

	httpServer.GET("/mjpeg", func(ctx *gin.Context) {
		boundary := "frame"
		ctx.Writer.Header().Add("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))
		ctx.Writer.WriteHeader(http.StatusOK)
		mjpegStreamer := streamer.NewMjpegStreamer(dev, boundary)
		mjpegStreamer.Streaming(ctx.Writer)
	})
	httpServer.Static("/static", "static")

	httpServer.Run(addr)
}
