package main

import (
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/config"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"github.com/zhayt/clean-arch-tmp-forum/internal/transport/http"
	"github.com/zhayt/clean-arch-tmp-forum/internal/transport/http/handler"
	"github.com/zhayt/clean-arch-tmp-forum/internal/transport/http/middleware"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const _path = "./config/config.json"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.NewConfig(_path)
	if err != nil {
		return err
	}

	l := logger.NewLogger()

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		return err
	}

	serv := service.NewService(repo, l)

	controller := handler.NewHandler(serv, l)
	mid := middleware.NewMiddleware(serv)

	httpServer := http.NewServer(cfg, l, controller, mid)

	l.Info.Printf("Start server on http://%v", net.JoinHostPort("localhost", cfg.Server.Port))
	httpServer.Start()

	osSignalCh := make(chan os.Signal, 1)
	signal.Notify(osSignalCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-osSignalCh:
		l.Info.Printf("signal accepted %s", s.String())
	case err := <-httpServer.Notify:
		l.Info.Printf("server closing %w", err)
	}

	if err = httpServer.Shutdown(); err != nil {
		return fmt.Errorf("error while shuttdown server: %w", err)
	}

	return nil
}
