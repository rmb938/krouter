package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/rmb938/krouter/pkg/router"
	"github.com/rmb938/krouter/pkg/storage/etcd"
)

func main() {

	brokerId := flag.Int("broker-id", 0, "sets the broker id, must be unique")

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

	if *brokerId < 1 {
		setupLog.Error(nil, "broker id must be set to a non-zero positive number")
		os.Exit(1)
	}

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

	database, err := etcd.NewETCDStorage("localhost:2379")
	if err != nil {
		setupLog.Error(err, "error creating etcd storage")
		os.Exit(1)
	}
	defer database.Close()

	r := &router.Router{
		Log: logr.WithName("router"),
	}

	err = r.StartServer(listenerAddr, advertiseListenerAddr, database, int32(*brokerId))
	if err != nil {
		setupLog.Error(err, "error starting router")
		os.Exit(1)
	}
}
