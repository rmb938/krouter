package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/zapr"
	"github.com/rmb938/krouter/apiServer"
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
	var brokerID int64
	var listener string
	var advertiseListener string
	var apiListener string
	var redisAddressesFlag MultiFlag

	flag.StringVar(&clusterID, "cluster-id", "some-cluster-id", "The Kafka cluster ID")
	flag.Int64Var(&brokerID, "broker-id", -1, "The Broker ID")
	flag.StringVar(&listener, "listener", "127.0.0.1:29092", "The address to listener on")
	flag.StringVar(&advertiseListener, "advertise-listener", "localhost:29092", "The address to advertise to clients")
	flag.StringVar(&apiListener, "api-listener", "127.0.0.1:6060", "The api address to listener on")
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

	if brokerID < 1 {
		setupLog.Error(fmt.Errorf("broker ID must be at least %d", 1), "error parsing brokerID")
		os.Exit(1)
	}

	if brokerID > math.MaxInt32 {
		setupLog.Error(fmt.Errorf("broker ID must be at most %d", math.MaxInt32), "error parsing brokerID")
		os.Exit(1)
	}

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

	kafkaRouter, err := router.NewRouter(logr.WithName("router"), listenerAddr, advertiseListenerAddr, clusterID, int32(brokerID), redisAddresses)
	if err != nil {
		setupLog.Error(err, "error creating router")
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	api := apiServer.NewAPIServer(logr.WithName("api-server"), apiListener, kafkaRouter.Broker)

	go func() {
		if err := api.ListenAndServe(); err != http.ErrServerClosed {
			setupLog.Error(err, "error running api listener")
			sigChan <- nil
		}
	}()

	go func() {
		err = kafkaRouter.ListenAndServe()
		if err != nil {
			setupLog.Error(err, "error starting router")
			sigChan <- nil
		}
	}()
	<-sigChan
	if err := api.Shutdown(context.TODO()); err != nil {
		setupLog.Error(err, "Error shutting down api server")
	}
	if err := kafkaRouter.Shutdown(); err != nil {
		setupLog.Error(err, "Error shutting down router")
	}
}
