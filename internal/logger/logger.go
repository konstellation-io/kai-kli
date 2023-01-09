package logger

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"gopkg.in/gookit/color.v1"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks -mock_names=Interface=MockLogger

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

const DefaultLogLevel = LevelInfo

var lineBreakRE = regexp.MustCompile(`\r?\n`)

type Interface interface {
	Success(msg string)
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type CliLogger struct {
	level  LogLevel
	writer io.Writer
}

// New creates a new Interface instance.
func New(level LogLevel) *CliLogger {
	return &CliLogger{
		level:  level,
		writer: os.Stdout,
	}
}

// New creates a new Interface instance.
func NewWithWriter(writer io.Writer) *CliLogger {
	return &CliLogger{
		level:  DefaultLogLevel,
		writer: writer,
	}
}

// NewDefaultLogger creates a new Interface instance.
func NewDefaultLogger() *CliLogger {
	return &CliLogger{
		level:  DefaultLogLevel,
		writer: os.Stdout,
	}
}

func (l *CliLogger) printLog(level LogLevel, msg, icon string) {
	if level > l.level {
		return
	}

	_, _ = fmt.Fprintf(l.writer, "[%s] %s\n", icon, lineBreakRE.ReplaceAllLiteralString(msg, " "))
}

func (l *CliLogger) Success(msg string) {
	icon := color.Success.Render("✔")
	l.printLog(LevelInfo, msg, icon)
}

func (l *CliLogger) Info(msg string) {
	icon := color.Info.Render("ℹ")
	l.printLog(LevelInfo, msg, icon)
}

func (l *CliLogger) Warn(msg string) {
	icon := color.Warn.Render("⚠")
	l.printLog(LevelWarn, msg, icon)
}

func (l *CliLogger) Error(msg string) {
	icon := color.Danger.Render("✖")
	l.printLog(LevelError, msg, icon)
}

func (l *CliLogger) Debug(msg string) {
	icon := color.Info.Render("✎")
	l.printLog(LevelDebug, msg, icon)
}
