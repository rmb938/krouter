package apiServer

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
)

type APIServer struct {
	log    logr.Logger
	server *http.Server
	broker *logical_broker.LogicalBroker
}

func NewAPIServer(log logr.Logger, addr string, broker *logical_broker.LogicalBroker) *APIServer {
	s := &APIServer{log: log, broker: broker}

	s.server = &http.Server{Addr: addr}

	return s
}

func (s *APIServer) ListenAndServe() error {
	s.log.Info("Starting API Server")
	return s.server.ListenAndServe()
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
