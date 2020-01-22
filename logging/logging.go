package logging

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func Log() *logrus.Entry {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	parts := strings.Split(nameFull, "/")
	name := strings.Split(nameFull, "/")[len(parts)-1]

	return l.WithFields(logrus.Fields{
		"caller": name,
	})
}
