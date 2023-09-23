package config

import (
	"github.com/timrourke/maschera/m/v2/app"
	"github.com/timrourke/maschera/m/v2/log"
	"go.uber.org/zap"
)

type Config interface {
	App() app.App
}

type config struct {
	logger log.Logger
}

func (c config) Logger() log.Logger {
	if c.logger != nil {
		return c.logger
	}

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	c.logger = log.NewLogger(zapLogger)

	return c.logger
}

func BuildApp() app.App {
	c := &config{}

	return app.NewApp(c.Logger())
}
