package controllers

import (
	"code.cloudfoundry.org/lager"
	"github.com/go-logr/logr"
)

type lagrLogger struct {
	logger logr.Logger
}

func NewLagrLogger(logger logr.Logger) lager.Logger {
	return &lagrLogger{
		logger: logger,
	}
}

func (l *lagrLogger) Session(task string, data ...lager.Data) lager.Logger {
	return NewLagrLogger(l.logger.WithName(task))
}

func (l *lagrLogger) Debug(action string, data ...lager.Data) {
	l.logger.Info("DEBUG "+action, lagerDataToLogrValues(data...)...)
}

func (l *lagrLogger) Info(action string, data ...lager.Data) {
	l.logger.Info(action, lagerDataToLogrValues(data...)...)
}

func (l *lagrLogger) Error(action string, err error, data ...lager.Data) {
	l.logger.Error(err, action, lagerDataToLogrValues(data...)...)
}

func (l *lagrLogger) Fatal(action string, err error, data ...lager.Data) {
	l.logger.Error(err, "FATAL "+action, lagerDataToLogrValues(data...)...)
	panic(err)
}

func (l *lagrLogger) WithData(data lager.Data) lager.Logger {
	return NewLagrLogger(l.logger.WithValues(lagerDataToLogrValues(data)...))
}

func (l *lagrLogger) RegisterSink(lager.Sink) {}

func (l *lagrLogger) SessionName() string {
	return ""
}

func lagerDataToLogrValues(dataCollections ...lager.Data) []interface{} {
	values := []interface{}{}

	for _, data := range dataCollections {
		for k, v := range data {
			values = append(values, k)
			values = append(values, v)
		}
	}

	return values
}
