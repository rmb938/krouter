package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/rmb938/krouter/pkg/router"
)

func main() {
	clusterID := flag.String("cluster-id", "some-cluster-id", "The Kafka cluster ID")
	listener := flag.String("listener", "localhost:29092", "The address to listener on")
	advertiseListener := flag.String("advertise-listener", "localhost:29092", "The address to advertise to clients")

	flag.Parse()

	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableStacktrace = true
	zapConfig.DisableCaller = true
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	zapLog, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("error creating logger: %v", err))
	}
	defer zapLog.Sync()

	logr := zapr.NewLogger(zapLog)
	setupLog := logr.WithName("setup")

	listenerAddr, err := net.ResolveTCPAddr("tcp", *listener)
	if err != nil {
		setupLog.Error(err, "error parsing listener address")
		os.Exit(1)
	}

	advertiseListenerAddr, err := net.ResolveTCPAddr("tcp", *advertiseListener)
	if err != nil {
		setupLog.Error(err, "error parsing advertise listener address")
		os.Exit(1)
	}

	r := &router.Router{
		Log: logr.WithName("router"),
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = r.ListenAndServe(listenerAddr, advertiseListenerAddr, *clusterID)
		if err != nil {
			setupLog.Error(err, "error starting router")
			sigChan <- nil
		}
	}()
	<-sigChan
	if err := r.Shutdown(); err != nil {
		setupLog.Error(err, "Error shutting down router")
	}
}
