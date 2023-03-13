package utils

import "github.com/sirupsen/logrus"

func LogError(error error) {
	logrus.Errorf(error.Error())
}
