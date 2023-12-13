package logging

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"gopkg.in/gookit/color.v1"
)

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

const (
	OutputFormatText = "text"
	OutputFormatJSON = "json"
)

const DefaultLogLevel = LevelInfo

var lineBreakRE = regexp.MustCompile(`\r?\n`)

type CliLogger struct {
	level        LogLevel
	writer       io.Writer
	outputFormat string
}

// New creates a new Interface instance.
func New(level LogLevel) *CliLogger {
	return &CliLogger{
		level:        level,
		writer:       os.Stdout,
		outputFormat: OutputFormatText,
	}
}

// New creates a new Interface instance.
func NewWithWriter(writer io.Writer) *CliLogger {
	return &CliLogger{
		level:        DefaultLogLevel,
		writer:       writer,
		outputFormat: OutputFormatText,
	}
}

// NewDefaultLogger creates a new Interface instance.
func NewDefaultLogger() *CliLogger {
	return &CliLogger{
		level:        DefaultLogLevel,
		writer:       os.Stdout,
		outputFormat: OutputFormatText,
	}
}

func (l *CliLogger) printLog(level LogLevel, msg, icon string) {
	if level > l.level {
		return
	}

	if l.outputFormat == OutputFormatJSON {
		_, _ = fmt.Fprintf(l.writer, "{\"status\":\"%s\",\"message\":\"%s\"}\n", icon, msg)
		return
	}

	_, _ = fmt.Fprintf(l.writer, "[%s] %s\n", icon, lineBreakRE.ReplaceAllLiteralString(msg, " "))
}

func (l *CliLogger) Success(msg string) {
	icon := color.Success.Render("✔")

	if l.outputFormat == OutputFormatJSON {
		icon = "success"
	}

	l.printLog(LevelInfo, msg, icon)
}

func (l *CliLogger) Info(msg string) {
	icon := color.Info.Render("ℹ")

	if l.outputFormat == OutputFormatJSON {
		icon = "info"
	}

	l.printLog(LevelInfo, msg, icon)
}

func (l *CliLogger) Warn(msg string) {
	icon := color.Warn.Render("⚠")

	if l.outputFormat == OutputFormatJSON {
		icon = "warn"
	}

	l.printLog(LevelWarn, msg, icon)
}

func (l *CliLogger) Error(msg string) {
	icon := color.Danger.Render("✖")

	if l.outputFormat == OutputFormatJSON {
		icon = "error"
	}

	l.printLog(LevelError, msg, icon)
}

func (l *CliLogger) Debug(msg string) {
	icon := color.Info.Render("✎")

	if l.outputFormat == OutputFormatJSON {
		icon = "debug"
	}

	l.printLog(LevelDebug, msg, icon)
}

func (l *CliLogger) SetDebugLevel() {
	l.level = LevelDebug
}

func (l *CliLogger) SetOutputFormat(of string) {
	l.outputFormat = of
}
