package v1

import (
	"github.com/go-logr/logr"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/models"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_acls/v1"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("describe-acls-v1-handler")
	request := message.(*v1.Request)

	response := &v1.Response{}

	authorizer := broker.GetController().GetAuthorizer()
	acls := authorizer.GetAcls(models.ACLResourceType(request.ResourceTypeFilter), request.ResourceNameFilter,
		models.ACLPatternType(request.PatternTypeFilter), request.PrincipalFilter,
		models.ACLOperation(request.Operation),
		models.ACLPermission(request.PermissionType))

	resources := make(map[uint64]*v1.ResponseResource)
	for _, acl := range acls {
		resource := &v1.ResponseResource{
			ResourceType: int8(acl.ResourceType),
			ResourceName: acl.ResourceName,
			PatternType:  int8(acl.PatternType),
		}

		hashInt, err := hashstructure.Hash(resource, hashstructure.FormatV2, nil)
		if err != nil {
			log.Error(err, "hashing acl for map key")
		}

		if _, ok := resources[hashInt]; !ok {
			resources[hashInt] = resource
		} else {
			resource = resources[hashInt]
		}

		resource.ACLs = append(resource.ACLs, v1.ResponseACL{
			Principal:      acl.Principal,
			Host:           "*",
			Operation:      int8(acl.Operation),
			PermissionType: int8(acl.Permission),
		})
	}

	for _, resource := range resources {
		response.Resources = append(response.Resources, *resource)
	}

	return response, nil
}
