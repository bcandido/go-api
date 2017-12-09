package app

import (
	"github.com/op/go-logging"
	"os"
)

func GetSimpleLog(module string) *logging.Logger {
	log := logging.MustGetLogger(module)
	format := logging.MustStringFormatter(`%{message}`)

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)
	logging.SetFormatter(format)
	log.SetBackend(backend1Leveled)

	return log
}