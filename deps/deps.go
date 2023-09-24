package deps

import (
	"github.com/timrourke/maschera/m/v2/app"
	"github.com/timrourke/maschera/m/v2/config"
	"github.com/timrourke/maschera/m/v2/log"
	"go.uber.org/zap"
)

type deps struct {
	cfg    config.Config
	logger log.Logger
}

func (c deps) Config() config.Config {
	if c.cfg != nil {
		return c.cfg
	}

	c.cfg = config.NewConfig()

	return c.cfg
}

func (c deps) Logger() log.Logger {
	if c.logger != nil {
		return c.logger
	}

	var logFunc func(options ...zap.Option) (*zap.Logger, error)
	switch c.Config().AppEnv() {
	case config.AppEnvDevelopment:
		logFunc = zap.NewDevelopment
		break
	case config.AppEnvProduction:
		logFunc = zap.NewProduction
		break
	case config.AppEnvTest:
		logFunc = func(_ ...zap.Option) (*zap.Logger, error) { return zap.NewNop(), nil }
		break
	}

	zapLogger, err := logFunc()
	if err != nil {
		panic(err)
	}

	c.logger = log.NewLogger(zapLogger)

	return c.logger
}

func BuildApp() app.App {
	c := &deps{}

	return app.NewApp(c.Logger())
}
