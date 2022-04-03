package v1

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/models"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/create_acls/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("create-acls-v1-handler")
	request := message.(*v1.Request)

	response := &v1.Response{}

	for _, creation := range request.Creations {
		result := v1.ResponseResult{}

		err := broker.GetController().GetAuthorizer().CreateAcl(
			models.ACLOperation(creation.Operation),
			models.ACLResourceType(creation.ResourceType),
			models.ACLPatternType(creation.ResourcePatternType),
			creation.ResourceName,
			creation.Principal,
			models.ACLPermission(creation.PermissionType),
		)
		if err != nil {
			log.Error(err, "Error creating ACL")
			result.ErrCode = errors.UnknownServerError
			result.ErrMessage = kmsg.StringPtr(fmt.Errorf("error creating acl: %w", err).Error())
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}
