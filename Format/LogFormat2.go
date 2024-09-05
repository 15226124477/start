package Format

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type LogFormatter2 struct{}

func (m *LogFormatter2) Format(entry *log.Entry) ([]byte, error) {
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
		newLog = fmt.Sprintf("%s %s| %s \n",
			timestamp, strings.ToUpper(entry.Level.String())[0:3], entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
