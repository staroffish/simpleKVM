package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/staroffish/simpleKVM/capture"
	"github.com/staroffish/simpleKVM/hid/common"
	"github.com/staroffish/simpleKVM/log"
	"github.com/staroffish/simpleKVM/streamer"
)

func StartHttpServer(ctx context.Context, addr string, dev *capture.CaptureDevice, hid common.Hid) {
	gin.SetMode(gin.ReleaseMode)
	httpServer := gin.New()
	// gin.Default().Use()

	httpStreamer, err := streamer.NewStreamer(dev)
	if err != nil {
		log.PrintInfo("new streamer error : %v", err)
		return
	}

	path := httpStreamer.Path()
	htmlElement := httpStreamer.HtmlElement()
	httpHandler := httpStreamer.Handler()

	indexPage := fmt.Sprintf(httpTemplateFormat, htmlElement)

	httpServer.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html", []byte(indexPage))
	})

	httpServer.GET(path, httpHandler)

	httpServer.GET("/pprof", func(ctx *gin.Context) {
		date := dumpPprof()
		ctx.Data(http.StatusOK, "text/html", []byte(fmt.Sprintf("<html>dumpping pprof data. you can find them in the program execution path. date=%s</html>", date)))
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
	httpServer.POST("/mousescroll", func(ctx *gin.Context) {
		scroll, exists := ctx.GetQuery("scroll")
		if !exists {
			ctx.Data(http.StatusBadRequest, "text/text", []byte("scroll code dose not exists"))
			return
		}
		scrollCode, err := strconv.ParseInt(scroll, 10, 16)
		if err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("invalid scroll code:%v", scrollCode)))
			return
		}

		if err := hid.MouseScroll(int(scrollCode)); err != nil {
			ctx.Data(http.StatusBadRequest, "text/text", []byte(fmt.Sprintf("mouse up error")))
		}
		ctx.Data(http.StatusOK, "text/text", []byte("ok"))
	})
	httpServer.Static("/static", "static")

	httpServer.Run(addr)
}

func dumpPprof() string {
	workDir, err := os.Getwd()
	nowStr := time.Now().Format("20060102150405000")
	if err != nil {
		log.PrintInfo("dump pprof error:Getwd:%v", err)
		return ""
	}

	go func() {
		fileName := fmt.Sprintf("%s_%s", "cpupprof", nowStr)
		dumpFile, err := os.OpenFile(filepath.Join(workDir, fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.PrintInfo("dump pprof error:OpenFile:%v", err)
			return
		}
		defer dumpFile.Close()
		if err := pprof.StartCPUProfile(dumpFile); err != nil {
			log.PrintInfo("dump pprof error:StartCPUProfile:%v", err)
			return
		}
		time.Sleep(time.Second * 60)
		pprof.StopCPUProfile()
	}()

	lookFunc := func(dumpType string) {
		fileName := fmt.Sprintf("%spprof_%s", dumpType, nowStr)
		dumpFile, err := os.OpenFile(filepath.Join(workDir, fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.PrintInfo("dump pprof error:OpenFile:%v", err)
			return
		}
		defer dumpFile.Close()
		if err := pprof.Lookup(dumpType).WriteTo(dumpFile, 1); err != nil {
			log.PrintInfo("dump pprof error:Lookup:%v", err)
			return
		}
	}

	lookFunc("heap")
	lookFunc("goroutine")

	return nowStr
}

var httpTemplateFormat = `<html>
<script type="text/javascript" src="/static/keyevent.js"></script>
<script type="text/javascript">
    window.document.oncontextmenu = function () { event.returnValue = false; }//disable right mouse button event  
</script>

<body onkeydown="return onKeyDown(event)" onkeyup="return onKeyUp(event)">
    <button onclick="shortcut([17,18,46])">ctrl+alt+del</button>
    <button onclick="shortcut([91, 76])">Win+L</button>
	<br>
	<table border="1">
	<tr><td>%s</td></tr>	
	</table>
</body>

</html>`
