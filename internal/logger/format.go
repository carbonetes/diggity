package logger

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

type Formatter struct{}

var levelList = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

var (
	levelColor  func(string) string
	timeColor   func(string) string
	callerColor func(string) string

	defaultColorScheme *ColorScheme = &ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		DebugLevelStyle: "blue",
		CallerStyle:     "cyan",
		TimestampStyle:  "magenta",
	}
	defaultCompiledColorScheme *compiledColorScheme = compileColorScheme(defaultColorScheme)
)

type ColorScheme struct {
	InfoLevelStyle  string
	WarnLevelStyle  string
	ErrorLevelStyle string
	FatalLevelStyle string
	PanicLevelStyle string
	DebugLevelStyle string
	CallerStyle     string
	TimestampStyle  string
}

type compiledColorScheme struct {
	InfoLevelColor  func(string) string
	WarnLevelColor  func(string) string
	ErrorLevelColor func(string) string
	FatalLevelColor func(string) string
	PanicLevelColor func(string) string
	DebugLevelColor func(string) string
	CallerColor     func(string) string
	TimestampColor  func(string) string
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = defaultCompiledColorScheme.InfoLevelColor
	case logrus.WarnLevel:
		levelColor = defaultCompiledColorScheme.WarnLevelColor
	case logrus.ErrorLevel:
		levelColor = defaultCompiledColorScheme.ErrorLevelColor
	case logrus.FatalLevel:
		levelColor = defaultCompiledColorScheme.FatalLevelColor
	case logrus.PanicLevel:
		levelColor = defaultCompiledColorScheme.PanicLevelColor
	default:
		levelColor = defaultCompiledColorScheme.DebugLevelColor
	}

	callerColor = defaultCompiledColorScheme.CallerColor
	timeColor = defaultCompiledColorScheme.TimestampColor
	levelText := levelList[int(entry.Level)]
	level := levelColor(fmt.Sprintf("%-5s", levelText))
	filePath := strings.Split(entry.Caller.File, "/")
	file := strings.Split(filePath[len(filePath)-1], ".")
	fileName := ellipsis(file[0], 10)
	fileName = "[" + fileName + ":" + strconv.Itoa(entry.Caller.Line) + "]"
	time := timeColor(fmt.Sprintf("[%s]", entry.Time.Format("01-02-2006 3:04PM")))
	caller := callerColor(fmt.Sprintf("%-15s", fileName))
	b.WriteString(fmt.Sprintf("%s %s %s : %s\n",
		time,
		caller,
		level,
		entry.Message))
	return b.Bytes(), nil
}

func compileColorScheme(s *ColorScheme) *compiledColorScheme {
	return &compiledColorScheme{
		InfoLevelColor:  getCompiledColor(s.InfoLevelStyle, defaultColorScheme.InfoLevelStyle),
		WarnLevelColor:  getCompiledColor(s.WarnLevelStyle, defaultColorScheme.WarnLevelStyle),
		ErrorLevelColor: getCompiledColor(s.ErrorLevelStyle, defaultColorScheme.ErrorLevelStyle),
		FatalLevelColor: getCompiledColor(s.FatalLevelStyle, defaultColorScheme.FatalLevelStyle),
		PanicLevelColor: getCompiledColor(s.PanicLevelStyle, defaultColorScheme.PanicLevelStyle),
		DebugLevelColor: getCompiledColor(s.DebugLevelStyle, defaultColorScheme.DebugLevelStyle),
		CallerColor:     getCompiledColor(s.CallerStyle, defaultColorScheme.CallerStyle),
		TimestampColor:  getCompiledColor(s.TimestampStyle, defaultColorScheme.TimestampStyle),
	}
}

func getCompiledColor(main string, fallback string) func(string) string {
	var style string
	if main != "" {
		style = main
	} else {
		style = fallback
	}
	return ansi.ColorFunc(style)
}

func ellipsis(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[0:maxLen])
}
