package main

import (
	"os"
	"os/signal"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"

	"github.com/wcharczuk/katwillmarry.com/pkg/config"
	"github.com/wcharczuk/katwillmarry.com/pkg/controller"
)

func main() {
	var cfg config.Config
	if err := configutil.Read(&cfg); !configutil.IsIgnored(err) {
		logger.FatalExit(err)
	}

	log := logger.NewFromConfig(&cfg.Logger)

	/*
		// uncomment when we need oauth ...
		auth, err := oauth.NewFromConfig(&cfg.OAuth)
		if err != nil {
			logger.FatalExit(err)
		}
	*/

	app := web.NewFromConfig(&cfg.Web)
	app.WithLogger(log)
	app.Register(&controller.Index{Log: log})

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		log.SyncFatalExit(app.Start())
	}()

	go func() {
		<-quit
		if err := app.Shutdown(); err != nil {
			log.SyncFatal(err)
		}
		close(done)
	}()
	<-done
}
