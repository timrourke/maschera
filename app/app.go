package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-multierror"
	"github.com/segmentio/kafka-go"
	"github.com/timrourke/maschera/m/v2/log"
	"github.com/timrourke/maschera/m/v2/pii"
	"go.uber.org/zap"
)

type App interface {
	Run(ctx context.Context) error
	Shutdown() error
}

type app struct {
	kafkaPIIReader *kafka.Reader
	logger         log.Logger
	piiMasker      pii.Masker
	shutdownFns    []func() error
}

func (a *app) Run(ctx context.Context) error {
	a.logger.Info("Running app")

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		err := a.piiMasker.Mask(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			a.logger.Error("Error masking PII", zap.Error(err))
		}
	}()

	<-ch

	a.logger.Info("Gracefully shutting down maschera...")

	err := a.Shutdown()

	return err
}

func (a *app) Shutdown() error {
	var result *multierror.Error

	for _, fn := range a.shutdownFns {
		if err := fn(); err != nil {
			e := multierror.Append(result, err)
			if e != nil {
				panic(e)
			}
		}
	}

	return result.ErrorOrNil()
}

func NewApp(logger log.Logger, piiMasker pii.Masker) App {
	return &app{
		logger:    logger,
		piiMasker: piiMasker,
		shutdownFns: []func() error{
			piiMasker.Shutdown,
		},
	}
}
