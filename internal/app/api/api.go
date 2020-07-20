package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/dink10/enlabs/internal/app/api/provider"
	"github.com/dink10/enlabs/internal/pkg/config"
	"github.com/dink10/enlabs/internal/pkg/database"
	"github.com/dink10/enlabs/internal/pkg/logger"
	"github.com/dink10/enlabs/internal/pkg/payments"
	"github.com/dink10/enlabs/internal/pkg/payments/storage"
	"github.com/dink10/enlabs/internal/pkg/router"
	"github.com/dink10/enlabs/internal/pkg/server"
)

// Run runs application.
func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg Config
	err := config.LoadConfig(&cfg)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	err = logger.Init(&cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	db, err := database.Connect(ctx, &cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer database.Close(db)

	paymentStorage := storage.NewPaymentStorage(db)
	paymentService, err := payments.NewService(paymentStorage)
	if err != nil {
		return fmt.Errorf("failed to init service: %v", err)
	}
	paymentProvider := provider.NewPaymentProvider(paymentService)

	r := router.NewDefaultRouter(cfg.Server.LogRequests)
	r.AddSubRouter("/v1", router.Routes{
		"/payments": paymentProvider.Router(),
	})

	go cancelOnSignal(cancel)

	httpServer := server.New(&cfg.Server, r.Handler())

	return httpServer.Run(ctx)
}

func cancelOnSignal(cancelFunc context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	sig := <-signals
	logrus.Infof("got signal %s, canceling app context", sig)

	cancelFunc()
}
