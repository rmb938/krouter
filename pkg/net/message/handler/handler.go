package handler

import (
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/net/message"
)

type MessageHandler interface {
	Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error)
}
