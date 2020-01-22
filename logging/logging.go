package logging

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func Log() *logrus.Entry {
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

	nameParts := strings.Split(name, ".")
	prefix := strings.Join(nameParts[:len(nameParts)-1], ".")

	return l.WithFields(logrus.Fields{
		"prefix": prefix,
	})
}
