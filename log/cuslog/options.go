package cuslog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

const FmtEmptySeparate = ""

type Level uint8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var LevelNameMapping = map[Level]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warn",
	ErrorLevel: "error",
	PanicLevel: "panic",
	FatalLevel: "fatal",
}

var errUnmarshalNilLevel = errors.New("can not unmarshal nil level")

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO", "": // make the zero value useful
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "panic", "PANIC":
		*l = PanicLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	default:
		return false
	}
	return true
}

func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) || !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unregocongnized level:%q", text)
	}
	return nil
}

type options struct {
	output        io.Writer
	level         Level
	stdLevel      Level
	formatter     Formatter
	disableCaller bool
}

type Option func(*options)

func initOptions(opts ...Option) (o *options) {
	o = &options{}
	for _, opt := range opts {
		opt(o)
	}
	if o.output == nil {
		o.output = os.Stderr
	}
	if o.formatter == nil {
		o.formatter = &TextFormatter{}
	}
	return
}

func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

func WithLevel(level Level) Option {
	return func(o *options) {
		o.level = level
	}
}

func WithStdLevel(level Level) Option {
	return func(o *options) {
		o.stdLevel = level
	}
}

func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.formatter = formatter

	}

}
func WithDisableCaller(caller bool) Option {
	return func(o *options) {
		o.disableCaller = caller

	}
}
