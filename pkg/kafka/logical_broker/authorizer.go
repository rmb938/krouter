package logical_broker

import (
	"context"
	"sync"

	"github.com/go-logr/logr"
	"github.com/puzpuzpuz/xsync"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/internal_topics_pb"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/models"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

type Authorizer struct {
	log logr.Logger

	kafkaClient *kgo.Client

	aclsSyncOnce sync.Once
	syncedChan   chan struct{}

	acls *xsync.MapOf[*models.ACL]
}

func NewAuthorizer(log logr.Logger, kafkaClient *kgo.Client) (*Authorizer, error) {
	authorizer := &Authorizer{
		log:         log.WithName("authorizer"),
		kafkaClient: kafkaClient,
		syncedChan:  make(chan struct{}),
		acls:        xsync.NewMapOf[*models.ACL](),
	}

	return authorizer, nil
}

func (c *Authorizer) CreateTestACLs() error {
	acls := []struct {
		Key   *internal_topics_pb.ACLMessageKey
		Value *internal_topics_pb.ACLMessageValue
	}{
		{
			Key: &internal_topics_pb.ACLMessageKey{
				Operation:    internal_topics_pb.ACLMessageKey_OPERATION_WRITE,
				ResourceType: internal_topics_pb.ACLMessageKey_RESOURCE_TYPE_TOPIC,
				PatternType:  internal_topics_pb.ACLMessageKey_PATTERN_TYPE_LITERAL,
				ResourceName: "topic-1",
				Principal:    "User:user1",
			},
			Value: &internal_topics_pb.ACLMessageValue{
				Permission: internal_topics_pb.ACLMessageValue_PERMISSION_ALLOW,
			},
		},
		{
			Key: &internal_topics_pb.ACLMessageKey{
				Operation:    internal_topics_pb.ACLMessageKey_OPERATION_READ,
				ResourceType: internal_topics_pb.ACLMessageKey_RESOURCE_TYPE_TOPIC,
				PatternType:  internal_topics_pb.ACLMessageKey_PATTERN_TYPE_LITERAL,
				ResourceName: "topic-1",
				Principal:    "User:user2",
			},
			Value: &internal_topics_pb.ACLMessageValue{
				Permission: internal_topics_pb.ACLMessageValue_PERMISSION_DENY,
			},
		},
		{
			Key: &internal_topics_pb.ACLMessageKey{
				Operation:    internal_topics_pb.ACLMessageKey_OPERATION_READ,
				ResourceType: internal_topics_pb.ACLMessageKey_RESOURCE_TYPE_GROUP,
				PatternType:  internal_topics_pb.ACLMessageKey_PATTERN_TYPE_PREFIXED,
				ResourceName: "my-group1",
				Principal:    "User:user2",
			},
			Value: &internal_topics_pb.ACLMessageValue{
				Permission: internal_topics_pb.ACLMessageValue_PERMISSION_DENY,
			},
		},
	}

	for _, acl := range acls {
		aclMessageKeyBytes, err := proto.Marshal(acl.Key)
		if err != nil {
			return err
		}

		aclMessageValueBytes, err := proto.Marshal(acl.Value)
		if err != nil {
			return err
		}

		record := kgo.KeySliceRecord(aclMessageKeyBytes, aclMessageValueBytes)
		record.Topic = InternalTopicAcls
		produceResp := c.kafkaClient.ProduceSync(context.TODO(), record)
		if produceResp.FirstErr() != nil {
			return produceResp.FirstErr()
		}
	}

	return nil
}

func (c *Authorizer) WaitSynced() {
	<-c.syncedChan
}

func (c *Authorizer) GetAcls(resourceType models.ACLResourceType, resourceName *string, patternType models.ACLPatternType, principal *string, operation models.ACLOperation, permission models.ACLPermission) []*models.ACL {
	var acls []*models.ACL

	c.acls.Range(func(_ string, acl *models.ACL) bool {
		if resourceType != models.ACLResourceTypeAny && resourceType != acl.ResourceType {
			return true
		}

		if resourceName != nil && *resourceName != acl.ResourceName {
			return true
		}

		if patternType != models.ACLPatternTypeAny && patternType != acl.PatternType {
			return true
		}

		if principal != nil && *principal != acl.Principal {
			return true
		}

		if operation != models.ACLOperationAny && operation != acl.Operation {
			return true
		}

		if permission != models.ACLPermissionAny && permission != acl.Permission {
			return true
		}

		acls = append(acls, acl)
		return true
	})

	return acls
}

func (c *Authorizer) CreateAcl(operation models.ACLOperation, resourceType models.ACLResourceType, patternType models.ACLPatternType, resourceName string, principal string, permission models.ACLPermission) error {

	aclKey := &internal_topics_pb.ACLMessageKey{
		Operation:    internal_topics_pb.ACLMessageKey_Operation(operation),
		ResourceType: internal_topics_pb.ACLMessageKey_ResourceType(resourceType),
		PatternType:  internal_topics_pb.ACLMessageKey_PatternType(patternType),
		ResourceName: resourceName,
		Principal:    principal,
	}

	aclValue := &internal_topics_pb.ACLMessageValue{
		Permission: internal_topics_pb.ACLMessageValue_Permission(permission),
	}

	aclMessageKeyBytes, err := proto.Marshal(aclKey)
	if err != nil {
		return err
	}

	aclMessageValueBytes, err := proto.Marshal(aclValue)
	if err != nil {
		return err
	}

	record := kgo.KeySliceRecord(aclMessageKeyBytes, aclMessageValueBytes)
	record.Topic = InternalTopicAcls
	produceResp := c.kafkaClient.ProduceSync(context.TODO(), record)
	if produceResp.FirstErr() != nil {
		return produceResp.FirstErr()
	}

	return nil
}
