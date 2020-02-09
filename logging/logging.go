package logging

import (
	"runtime"
	"strings"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	prefixedFormatter "github.com/x-cray/logrus-prefixed-formatter"
)

type Fields logrus.Fields

func Stdout(stdout string) Fields {
	return Fields{"stdout": stdout}
}

func StdoutFields(stdout string, fields Fields) Fields {
	mergo.Merge(
		&fields,
		Fields{"stdout": stdout},
	)
	return fields
}

type Logger struct {
	*logrus.Logger
}

type Options struct {
	Level logrus.Level
}

var defaultOptions = Options{
	Level: logrus.InfoLevel,
}

var logger = &Logger{logrus.New()}

func GetLevel(lvl string) logrus.Level {
	switch lvl {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warning":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	default:
		return logrus.InfoLevel
	}
}

func Configure(options ...Options) {
	var _options Options
	if len(options) > 0 {
		_options = options[0]
	}
	l := logrus.New()
	l.SetFormatter(&prefixedFormatter.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	l.SetLevel(_options.Level)
	logger = &Logger{l}
}

func Log() *prefixed {

	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	parts := strings.Split(nameFull, "/")
	name := strings.Split(nameFull, "/")[len(parts)-1]

	prefix := name
	nameParts := strings.Split(name, ".")
	if nameParts[len(nameParts)-1] == "func1" {
		prefix = strings.Join(nameParts[:len(nameParts)-1], ".")
	}

	return &prefixed{logger.WithFields(logrus.Fields{
		"prefix": prefix,
	})}
}

type prefixed struct {
	prefix *logrus.Entry
}

func (p *prefixed) Trace(fields Fields, args ...interface{}) *prefixed {
	p.prefix.WithFields(logrus.Fields(fields)).Trace(args...)
	return p
}

func (p *prefixed) Debug(fields Fields, args ...interface{}) *prefixed {
	p.prefix.WithFields(logrus.Fields(fields)).Debug(args...)
	return p
}

func (p *prefixed) Info(fields Fields, args ...interface{}) *prefixed {
	p.prefix.WithFields(logrus.Fields(fields)).Info(args...)
	return p
}

func (p *prefixed) Warn(fields Fields, args ...interface{}) *prefixed {
	p.prefix.WithFields(logrus.Fields(fields)).Warn(args...)
	return p
}

func (p *prefixed) WarnIf(fields Fields, condition bool, args ...interface{}) *prefixed {
	if condition {
		p.prefix.WithFields(logrus.Fields(fields)).Warn(args...)
	}
	return p
}

func (p *prefixed) Stderr(stderr string, err error, args ...interface{}) *prefixed {
	if strings.TrimSpace(stderr) != "" {
		entry := p.prefix.WithFields(logrus.Fields{"stderr": stderr})
		entry.Warn(args...)
	}
	if err != nil {
		entry := p.prefix.WithFields(logrus.Fields{"err": err})
		entry.Error(args...)
	}
	return p
}

func (p *prefixed) Error(err error, args ...interface{}) *prefixed {
	if err != nil {
		p.prefix.WithFields(logrus.Fields{"err": err}).Error(args...)
	}
	return p
}

func (p *prefixed) ErrorFields(err error, fields Fields, args ...interface{}) *prefixed {
	if err != nil {
		mergo.Merge(
			&fields,
			Fields{"err": err},
		)
		p.prefix.WithFields(logrus.Fields(fields)).Error(args...)
	}
	return p
}
