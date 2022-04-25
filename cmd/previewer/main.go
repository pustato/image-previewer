package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/pustato/image-previewer/internal/app"
	"github.com/pustato/image-previewer/internal/cache"
	"github.com/pustato/image-previewer/internal/client"
	"github.com/pustato/image-previewer/internal/logger"
	"github.com/pustato/image-previewer/internal/resizer"
	"github.com/pustato/image-previewer/internal/server"
)

var (
	port      = flag.String("port", "8000", "service port")
	cacheDir  = flag.String("cacheDir", "/tmp/cache", "directory to store cache")
	cacheSize = flag.String("cacheSize", "10M", "directory to store cache")
	logLevel  = flag.String("logLevel", "debug", "logging level")
)

func main() {
	resultCode := 0
	defer func() {
		os.Exit(resultCode)
	}()

	flag.Parse()
	if port == nil || *port == "" {
		flag.PrintDefaults()
		return
	}

	logg, err := logger.New(*logLevel)
	if err != nil {
		fmt.Println("create logger: " + err.Error())
		resultCode = 1
		return
	}

	cacheSizeBytes, err := bytefmt.ToBytes(*cacheSize)
	if err != nil {
		logg.Error("invalid cache size: " + err.Error())
		resultCode = 1
		return
	}

	clientInstance := client.New(time.Second)
	resizerInstance := resizer.New()

	appInstance := app.New(clientInstance, resizerInstance)
	cachedApp, err := cache.New(appInstance, cacheSizeBytes, *cacheDir)
	if err != nil {
		logg.Error("create cached app: " + err.Error())
		resultCode = 1
		return
	}

	srv := server.New(net.JoinHostPort("0.0.0.0", *port), cachedApp, logg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		logg.Info("stopping server...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := srv.Stop(ctx); err != nil {
			resultCode = 1
			logg.Error("stop server:" + err.Error())
		}
	}()

	logg.Info("starting server on " + *port)
	if err := srv.Start(); err != nil {
		resultCode = 1
		logg.Error("start server: " + err.Error())
	}
}
