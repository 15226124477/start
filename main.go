package start

import (
	_ "embed"
	"fmt"
	"github.com/15226124477/define"
	"github.com/15226124477/start/Format"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func LiveServer() {
	for {
		runExe, _ := os.Executable()
		_, exec := filepath.Split(runExe)
		filenameWithSuffix := path.Base(exec)
		fileSuffix := path.Ext(filenameWithSuffix)
		filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
		filePath := fmt.Sprintf("./log/%s%s.log", filenameOnly, time.Now().Format("2006-01-02"))
		log.Warning(filePath)
		// 配置tail
		config := tail.Config{
			ReOpen:    true,                                 // 文件被截断后重新打开
			Follow:    true,                                 // 跟随文件，监控新增内容
			Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件末尾开始读取
			MustExist: false,                                // 文件不存在不报错
			Poll:      true,                                 // 使用轮询模式
		}
		// 使用tail.TailFile创建一个Tail对象
		tails, err := tail.TailFile(filePath, config)
		if err != nil {
			log.Error(err)
			return
		}
		u := "ws://172.16.32.244:1024/log"
		log.Printf("连接到服务器 %s", u)
		c, _, err := websocket.DefaultDialer.Dial(u, nil) //试连接到WebSocket服务器
		if err != nil {
			log.Error("失败:", err)
			return
		}
		defer c.Close()
		for line := range tails.Lines {
			err := c.WriteMessage(websocket.TextMessage, []byte(line.Text))
			if err != nil {
				log.Error("向WebSocket服务端发送消息出错:", err)
				return
			}
		}
		// 关闭连接
		c.Close()
		time.Sleep(time.Second * 5)
	}
}

func RotateLogs(mode int, isLive bool) *rotatelogs.RotateLogs {
	log.SetReportCaller(true)
	log.SetLevel(log.TraceLevel)
	if mode == 0 {
		log.SetFormatter(&Format.LogFormatter1{})
		log.SetFormatter(&CustomTextFormatter1{
			ForceColors:   true,
			ColorDebug:    color.New(color.FgHiWhite),
			ColorInfo:     color.New(color.BgBlue, color.FgHiWhite),
			ColorWarning:  color.New(color.BgMagenta, color.FgHiWhite),
			ColorError:    color.New(color.BgRed, color.FgHiWhite),
			ColorCritical: color.New(color.BgRed, color.FgHiWhite),
		})
	} else {
		log.SetFormatter(&Format.LogFormatter2{})
		log.SetFormatter(&CustomTextFormatter2{
			ForceColors:   true,
			ColorDebug:    color.New(color.FgHiWhite),
			ColorInfo:     color.New(color.BgBlue, color.FgHiWhite),
			ColorWarning:  color.New(color.BgMagenta, color.FgHiWhite),
			ColorError:    color.New(color.BgRed, color.FgHiWhite),
			ColorCritical: color.New(color.BgRed, color.FgHiWhite),
		})
	}

	runExe, _ := os.Executable()
	_, exec := filepath.Split(runExe)
	ext := path.Ext(exec)
	name := strings.TrimSuffix(exec, ext)

	writer, err := rotatelogs.New(
		"./log/"+name+"%Y-%m-%d.log",           // 日志文件的命名模式，包括时间戳
		rotatelogs.WithLinkName(name),          // 创建一个软链接指向最新的日志文件
		rotatelogs.WithMaxAge(24*time.Hour),    // 设置日志文件的最大保存时间
		rotatelogs.WithRotationTime(time.Hour), // 设置日志文件的轮转时间间隔
	)
	if err != nil {
		log.Fatalf("Error creating rotate logger: %v", err)
	}

	// 使用自定义的RotateLogs对象作为日志的输出
	log.SetOutput(writer)

	if mode == 0 {
		lfsHook := lfshook.NewHook(lfshook.WriterMap{
			log.DebugLevel: writer,
			log.InfoLevel:  writer,
			log.WarnLevel:  writer,
			log.ErrorLevel: writer,
			log.FatalLevel: writer,
			log.PanicLevel: writer,
		}, &Format.LogFormatter1{})
		log.AddHook(lfsHook)
	} else {
		lfsHook := lfshook.NewHook(lfshook.WriterMap{
			log.DebugLevel: writer,
			log.InfoLevel:  writer,
			log.WarnLevel:  writer,
			log.ErrorLevel: writer,
			log.FatalLevel: writer,
			log.PanicLevel: writer,
		}, &Format.LogFormatter2{})
		log.AddHook(lfsHook)
	}

	log.Warning(exec, " \tProgram Running......")
	log.Debug(define.Setting.ProgramVersion)
	// 打开日志转发
	if isLive {
		go func() {
			LiveServer()
		}()
	}

	return writer
}
