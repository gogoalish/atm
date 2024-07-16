package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogoalish/atm/config"
	"github.com/gogoalish/atm/internal/app"
	"github.com/gogoalish/atm/internal/logger"
	"github.com/gogoalish/atm/internal/server"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal("error init config", err)
	}

	l, err := logger.New()
	if err != nil {
		log.Fatal("error init logger", err)
	}
	defer l.Sync()

	accountCntrl := app.NewAccountsController()

	router := server.NewRouter(accountCntrl, l)
	httpServer := server.New(cfg, router)
	l.Info(fmt.Sprintf("server is listening on: http://%s:%s", cfg.Host, cfg.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("main - signal:" + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Sprint("main - httpServer.Notify: ", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Sprint("main - httpServer.Shutdown: ", err))
	}

}
