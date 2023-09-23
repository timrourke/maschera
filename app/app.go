package app

import "github.com/timrourke/maschera/m/v2/log"

type App interface {
	Run() error
}

type app struct {
	logger log.Logger
}

func (a app) Run() error {
	a.logger.Info("Running app")

	return nil
}

func NewApp(logger log.Logger) App {
	return &app{logger: logger}
}
