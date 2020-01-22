package logging

import (
	"runtime"
	"strings"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type Fields logrus.Fields

type logger struct {
	prefixEntry *logrus.Entry
}

func Log() *logger {
	l := logrus.New()
	l.SetFormatter(&prefixed.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	parts := strings.Split(nameFull, "/")
	name := strings.Split(nameFull, "/")[len(parts)-1]

	prefix := name
	nameParts := strings.Split(name, ".")
	if nameParts[len(nameParts)-1] == "func1" {
		prefix = strings.Join(nameParts[:len(nameParts)-1], ".")
	}

	return &logger{l.WithFields(logrus.Fields{
		"prefix": prefix,
	})}
}

func (l *logger) Stdout(stdout string) *logger {
	l.prefixEntry.WithFields(logrus.Fields{"stdout": stdout}).Info()
	return l
}

func (l *logger) StdoutFields(stdout string, fields Fields) *logger {
	mergo.Merge(
		&fields,
		Fields{"stdout": stdout},
	)
	l.prefixEntry.WithFields(logrus.Fields(fields)).Info()
	return l
}

func (l *logger) Info(fields Fields) *logger {
	l.prefixEntry.WithFields(logrus.Fields(fields)).Info()
	return l
}

func (l *logger) Warn(fields Fields) *logger {
	l.prefixEntry.WithFields(logrus.Fields(fields)).Warn()
	return l
}

func (l *logger) Stderr(stderr string, err error) *logger {
	if strings.TrimSpace(stderr) != "" {
		entry := l.prefixEntry.WithFields(logrus.Fields{"stderr": stderr})
		entry.Warn()
	}
	if err != nil {
		entry := l.prefixEntry.WithFields(logrus.Fields{"err": err})
		entry.Error()
	}
	return l
}

func (l *logger) Error(err error) *logger {
	if err != nil {
		l.prefixEntry.WithFields(logrus.Fields{"err": err}).Error()
	}
	return l
}
