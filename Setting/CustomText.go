package Setting

import (
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strconv"
	"strings"
)

type CustomTextFormatter1 struct {
	log.TextFormatter
	ForceColors   bool
	ColorInfo     *color.Color
	ColorDebug    *color.Color
	ColorWarning  *color.Color
	ColorError    *color.Color
	ColorCritical *color.Color
}

func (f *CustomTextFormatter1) Format(entry *log.Entry) ([]byte, error) {
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

func (f *CustomTextFormatter1) PrintColored(entry *log.Entry) {
	levelColor := color.New(color.FgCyan, color.Bold)             // 定义蓝色和粗体样式
	levelText := levelColor.Sprintf("%-6s", entry.Level.String()) // 格式化日志级别文本

	msg := fmt.Sprintf("%s %s", entry.Message, levelText)
	if entry.HasCaller() {
		msg += " (" + entry.Caller.File + ":" + strconv.Itoa(entry.Caller.Line) + ")" // 添加调用者信息
	}

	fmt.Fprintln(color.Output, msg) // 使用有颜色的方式打印消息到终端
}

type CustomTextFormatter2 struct {
	log.TextFormatter
	ForceColors   bool
	ColorInfo     *color.Color
	ColorDebug    *color.Color
	ColorWarning  *color.Color
	ColorError    *color.Color
	ColorCritical *color.Color
}

func (f *CustomTextFormatter2) Format(entry *log.Entry) ([]byte, error) {
	if f.ForceColors {
		txt := fmt.Sprintf("%s %s:| %s", entry.Time.Format("2006-01-02 15:04:05.000"), strings.ToUpper(entry.Level.String())[0:3], entry.Message)
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

func (f *CustomTextFormatter2) PrintColored(entry *log.Entry) {
	levelColor := color.New(color.FgCyan, color.Bold)             // 定义蓝色和粗体样式
	levelText := levelColor.Sprintf("%-6s", entry.Level.String()) // 格式化日志级别文本

	msg := fmt.Sprintf("%s %s", entry.Message, levelText)
	if entry.HasCaller() {
		msg += " (" + entry.Caller.File + ":" + strconv.Itoa(entry.Caller.Line) + ")" // 添加调用者信息
	}

	fmt.Fprintln(color.Output, msg) // 使用有颜色的方式打印消息到终端
}
