package start

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/15226124477/define"
	"github.com/15226124477/method"
	"github.com/fatih/color"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type LogFormatter struct{}

type CustomTextFormatter struct {
	log.TextFormatter
	ForceColors   bool
	ColorInfo     *color.Color
	ColorDebug    *color.Color
	ColorWarning  *color.Color
	ColorError    *color.Color
	ColorCritical *color.Color
}

func (f *CustomTextFormatter) Format(entry *log.Entry) ([]byte, error) {
	if f.ForceColors {

		txt := fmt.Sprintf("%s %s: %-47s | %-20s:%-5d | %s", entry.Time.Format("2006-01-02 15:04:05.000"), strings.ToUpper(entry.Level.String())[0:3], entry.Caller.Function, filepath.Base(entry.Caller.File), entry.Caller.Line, entry.Message)
		switch entry.Level {
		case log.InfoLevel:
			_, err := f.ColorInfo.Println(txt)
			if err != nil {
				return nil, err
			} // 使用蓝色打印信息日志
		case log.WarnLevel:
			_, err := f.ColorWarning.Println(txt)
			if err != nil {
				return nil, err
			} // 使用黄色打印警告日志
		case log.ErrorLevel:
			_, err := f.ColorError.Println(txt)
			if err != nil {
				return nil, err
			} // 使用红色打印错误日志
		case log.DebugLevel:
			_, err := f.ColorDebug.Println(txt)
			if err != nil {
				return nil, err
			}
		case log.FatalLevel, log.PanicLevel:
			_, err := f.ColorCritical.Println(txt)
			if err != nil {
				return nil, err
			} // 使用带有红色背景和白色文本的样式打印严重日志
		default:
			f.PrintColored(entry)
		}
		return nil, nil
	} else {
		return f.TextFormatter.Format(entry)
	}
}

func (f *CustomTextFormatter) PrintColored(entry *log.Entry) {
	levelColor := color.New(color.FgCyan, color.Bold)             // 定义蓝色和粗体样式
	levelText := levelColor.Sprintf("%-6s", entry.Level.String()) // 格式化日志级别文本

	msg := fmt.Sprintf("%s %s", entry.Message, levelText)
	if entry.HasCaller() {
		msg += " (" + entry.Caller.File + ":" + strconv.Itoa(entry.Caller.Line) + ")" // 添加调用者信息
	}

	fmt.Fprintln(color.Output, msg) // 使用有颜色的方式打印消息到终端
}

func (m *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("%s %s: %-47s | %-20s:%-5d | %s \n",
			timestamp, strings.ToUpper(entry.Level.String())[0:3], entry.Caller.Function, fName, entry.Caller.Line, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func CreatFolders(folders []string) {
	method.CreatFolders(folders)
}

func RotateLogs() *rotatelogs.RotateLogs {
	log.SetReportCaller(true)
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&LogFormatter{})
	log.SetFormatter(&CustomTextFormatter{
		ForceColors:   true,
		ColorDebug:    color.New(color.FgHiWhite),
		ColorInfo:     color.New(color.BgBlue, color.FgHiWhite),
		ColorWarning:  color.New(color.BgMagenta, color.FgHiWhite),
		ColorError:    color.New(color.BgRed, color.FgHiWhite),
		ColorCritical: color.New(color.BgRed, color.FgHiWhite),
	})
	runExe, _ := os.Executable()
	_, exec := filepath.Split(runExe)
	ext := path.Ext(exec)
	name := strings.TrimSuffix(exec, ext)

	writer, err := rotatelogs.New(
		"./log/"+name+"%Y-%m-%d.log",           // 日志文件的命名模式，包括时间戳
		rotatelogs.WithLinkName(name),          // 创建一个软链接指向最新的日志文件
		rotatelogs.WithMaxAge(24*time.Hour),    // 设置日志文件的最大保存时间
		rotatelogs.WithRotationTime(time.Hour), // 设置日志文件的轮转时间间隔
		// rotatelogs.WithRotationCount(30),
	)
	if err != nil {
		log.Fatalf("Error creating rotate logger: %v", err)
	}

	// 使用自定义的RotateLogs对象作为日志的输出
	log.SetOutput(writer)

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &LogFormatter{})

	log.AddHook(lfsHook)
	log.Warning(exec, " \tProgram Running......")
	log.Debug(define.Setting.ProgramVersion)
	// 打开日志转发
	go func() {
		method.LiveServer()
	}()
	return writer
}
