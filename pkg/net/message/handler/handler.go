package handler

import (
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/net/message"
)

type MessageHandler interface {
	Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error
}
