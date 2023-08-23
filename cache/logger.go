package cache

import (
	"context"
	"fmt"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type EmptyLogger struct {
}

func (el EmptyLogger) Warnf(format string, args ...interface{}) {
	warnPrefix := "[Warn] "
	fmt.Printf(warnPrefix+format, args)
	return
}

func (el EmptyLogger) Errorf(format string, args ...interface{}) {
	errorPrefix := "[Error] "
	fmt.Printf(errorPrefix+format, args)
	return
}

func (el EmptyLogger) Infof(format string, args ...interface{}) {
	infoPrefix := "[Info] "
	fmt.Printf(infoPrefix+format, args)
	return
}

type LoggerBuilder func(context.Context) Logger
