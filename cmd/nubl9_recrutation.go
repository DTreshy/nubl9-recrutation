package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/DTreshy/nubl9-recrutation/internal/config"
	"github.com/DTreshy/nubl9-recrutation/internal/flags"
	"github.com/DTreshy/nubl9-recrutation/internal/logger"
	"github.com/DTreshy/nubl9-recrutation/server"
)

func main() {
	defer os.Exit(1)

	f := flags.Parse()

	cfg, err := config.New(f.ConfigPath)
	if err != nil {
		fmt.Println(err)
		runtime.Goexit()
	}

	log, err := logger.New(cfg.Log)
	if err != nil {
		log.Sugar().Error(err.Error())
		runtime.Goexit()
	}

	httpServer := server.New(log)

	if err := httpServer.Routes(); err != nil {
		log.Sugar().Error(err.Error())
	}

	httpCloseC := make(chan struct{})

	go func() {
		if err := httpServer.Run(cfg.Net.HTTPBind); err != nil {
			log.Sugar().Error(err.Error())
		}

		close(httpCloseC)
	}()

	defer func() {
		if httpStopErr := httpServer.Shutdown(); httpStopErr != nil {
			log.Sugar().Error(httpStopErr.Error())
		}

		log.Sugar().Info("Stopped HTTP server")
	}()

	interruptC := make(chan os.Signal, 1)

	signal.Notify(interruptC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(interruptC)
		close(interruptC)
	}()

	select {
	case <-interruptC:
	case <-httpCloseC:
	}
}
