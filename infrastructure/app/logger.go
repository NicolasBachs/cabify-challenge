package app

import (
	"github.com/sirupsen/logrus"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/config"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/domain/logger"
)

type logrusLogger struct {
	env    string
	logger *logrus.Logger
}

var Logger logger.Logger

func NewAppLogger() logger.Logger {
	logger := logrus.New()

	if config.IsProductionEnv() {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}

	return &logrusLogger{
		logger: logger,
		env:    config.AppConfig.App.Environment,
	}
}

func (log *logrusLogger) Debug(service string, msg string, args ...interface{}) {
	log.logger.WithFields(logrus.Fields{
		"service": service,
		"env":     log.env,
	}).Debugf(msg, args...)
}

func (log *logrusLogger) Info(service string, msg string, args ...interface{}) {
	log.logger.WithFields(logrus.Fields{
		"service": service,
		"env":     log.env,
	}).Infof(msg, args...)
}

func (log *logrusLogger) Warn(service string, msg string, args ...interface{}) {
	log.logger.WithFields(logrus.Fields{
		"service": service,
		"env":     log.env,
	}).Warnf(msg, args...)
}

func (log *logrusLogger) Error(service string, msg string, args ...interface{}) {
	log.logger.WithFields(logrus.Fields{
		"service": service,
		"env":     log.env,
	}).Errorf(msg, args...)
}
