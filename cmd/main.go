package main

import (
	"flag"
	"os"
	"os/signal"
	"receipt-processor-challenge/config"
	. "receipt-processor-challenge/internal"
	"receipt-processor-challenge/internal/handler"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	configFile = flag.String("config", "./config.yaml", "config file path")
	log        *logrus.Logger
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	flag.Parse()
	cfg, err := config.Setup(*configFile)
	if err != nil {
		panic(err)
	}

	log = logrus.New()
	log.SetLevel(logrus.DebugLevel)
	h := handler.NewHandler(log)
	router := NewProcessorService(h)
	server := fasthttp.Server{
		ReadBufferSize: 10240,
		Handler:        router.Handler,
		Logger:         log,
	}
	log.Infof("server started on port %s", cfg.App.Port)
	log.Fatal(server.ListenAndServe(":" + cfg.App.Port))
	<-sigChan
}
