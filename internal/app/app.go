package app

import (
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/config"
	"github.com/zhayt/clean-arch-tmp-forum/internal/controller"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"github.com/zhayt/clean-arch-tmp-forum/internal/storage/sqlite"
	"github.com/zhayt/clean-arch-tmp-forum/pkg/httpserver"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const migrationPath = "migrations/20230308_init.up.sql"

func Run(config *config.Config) error {
	storage, err := sqlite.NewStorage(config.Database.Driver, config.Database.DSN)
	if err != nil {
		return fmt.Errorf("can't get storage: %w", err)
	}

	err = storage.Init(migrationPath)
	if err != nil {
		return err
	}

	services := service.NewUserService(storage)
	controllers := controller.NewController(services)

	httpServer := httpserver.New(controllers.InitRoute(), httpserver.Port(config.Server.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	return nil
}
