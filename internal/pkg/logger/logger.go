package logger

import "github.com/sirupsen/logrus"

func InitLogger() {

	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetLevel(logrus.InfoLevel)
}
