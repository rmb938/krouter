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

type MultiFlag []string

func (m *MultiFlag) String() string {
	if len(*m) == 0 {
		return ""
	}
	return fmt.Sprint(*m)
}

func (m *MultiFlag) Set(val string) error {
	*m = append(*m, val)
	return nil
}

func main() {
	var clusterID string
	var listener string
	var advertiseListener string
	var redisAddressesFlag MultiFlag

	flag.StringVar(&clusterID, "cluster-id", "some-cluster-id", "The Kafka cluster ID")
	flag.StringVar(&listener, "listener", "localhost:29092", "The address to listener on")
	flag.StringVar(&advertiseListener, "advertise-listener", "localhost:29092", "The address to advertise to clients")
	flag.Var(&redisAddressesFlag, "redis-addresses", "List of redis addresses (default \"localhost:6379\")")

	flag.Parse()

	redisAddresses := []string(redisAddressesFlag)

	if len(redisAddresses) == 0 {
		redisAddresses = append(redisAddresses, "localhost:6379")
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableStacktrace = true
	zapConfig.DisableCaller = true
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	zapLog, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("error creating logger: %v", err))
	}
	defer zapLog.Sync()

	logr := zapr.NewLogger(zapLog)
	setupLog := logr.WithName("setup")

	listenerAddr, err := net.ResolveTCPAddr("tcp", listener)
	if err != nil {
		setupLog.Error(err, "error parsing listener address")
		os.Exit(1)
	}

	advertiseListenerAddr, err := net.ResolveTCPAddr("tcp", advertiseListener)
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
		err = r.ListenAndServe(listenerAddr, advertiseListenerAddr, clusterID, redisAddresses)
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
