package handler

import (
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/net/message"
)

type MessageHandler interface {
	Handle(client *client.Client, message message.Message, correlationId int32) error
}
