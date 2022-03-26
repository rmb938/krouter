package logical_broker

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/twmb/franz-go/pkg/kmsg"
)

const (
	GroupCoordinatorRedisKeyFmt          = "{group-%s}-coordinator"
	GroupGenerationRedisKeyFmt           = "{group-%s}-generation"
	GroupTopicOffsetRedisKeyFmt          = "{group-%s}-offset-topic"
	GroupTopicPartitionOffsetRedisKeyFmt = GroupTopicOffsetRedisKeyFmt + "-%s-partition-%d"
)

func (c *Controller) FindGroupCoordinator(consumerGroup string) (*kmsg.FindCoordinatorResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	request := kmsg.NewPtrFindCoordinatorRequest()
	request.CoordinatorType = 0
	request.CoordinatorKey = consumerGroup

	response, err := c.franzKafkaClient.Request(ctx, request)
	if err != nil {
		return nil, err
	}

	coordinatorResponse := response.(*kmsg.FindCoordinatorResponse)

	if coordinatorResponse.ErrorCode == int16(errors.None) {
		// Consumers don't actively refresh this
		// so if this expires we should return errors.NotCoordinator and clients will try and FindGroupCoordinator again
		err := c.redisClient.Client.Set(ctx, fmt.Sprintf(GroupCoordinatorRedisKeyFmt, consumerGroup), coordinatorResponse.NodeID, 1*time.Hour).Err()
		if err != nil {
			return nil, err
		}
	}

	return coordinatorResponse, nil
}

func (c *Controller) groupCoordinator(ctx context.Context, consumerGroup string) (int, error) {
	brokerID, err := c.redisClient.Client.Get(ctx, fmt.Sprintf(GroupCoordinatorRedisKeyFmt, consumerGroup)).Int()
	if err != nil {
		if err == redis.Nil {
			return -1, nil
		}
		return -1, err
	}

	return brokerID, nil
}

func (c *Controller) JoinGroup(request *kmsg.JoinGroupRequest) (*kmsg.JoinGroupResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	coordinatorId, err := c.groupCoordinator(ctx, request.Group)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		response := kmsg.NewPtrJoinGroupResponse()
		response.ErrorCode = int16(errors.NotCoordinator)

		return response, nil
	}

	resp, err := c.franzKafkaClient.Broker(coordinatorId).RetriableRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	joinGroupResponse := resp.(*kmsg.JoinGroupResponse)

	if joinGroupResponse.ErrorCode == int16(errors.None) {
		redisGroupGenerationKey := fmt.Sprintf(GroupGenerationRedisKeyFmt, request.Group)
		err = c.redisClient.Client.Watch(ctx, func(tx *redis.Tx) error {
			// TODO: make exp configurable
			return tx.Set(ctx, redisGroupGenerationKey, joinGroupResponse.Generation, 7*24*time.Hour).Err()
		}, redisGroupGenerationKey)
	}

	return joinGroupResponse, err
}

func (c *Controller) SyncGroup(request *kmsg.SyncGroupRequest) (*kmsg.SyncGroupResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	coordinatorId, err := c.groupCoordinator(ctx, request.Group)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		response := kmsg.NewPtrSyncGroupResponse()
		response.ErrorCode = int16(errors.NotCoordinator)

		return response, nil
	}

	resp, err := c.franzKafkaClient.Broker(coordinatorId).RetriableRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*kmsg.SyncGroupResponse), nil
}

func (c *Controller) LeaveGroup(request *kmsg.LeaveGroupRequest) (*kmsg.LeaveGroupResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	coordinatorId, err := c.groupCoordinator(ctx, request.Group)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		response := kmsg.NewPtrLeaveGroupResponse()
		response.ErrorCode = int16(errors.NotCoordinator)

		return response, nil
	}

	resp, err := c.franzKafkaClient.Broker(coordinatorId).RetriableRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*kmsg.LeaveGroupResponse), nil
}

func (c *Controller) HeartBeat(request *kmsg.HeartbeatRequest) (*kmsg.HeartbeatResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	coordinatorId, err := c.groupCoordinator(ctx, request.Group)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		response := kmsg.NewPtrHeartbeatResponse()
		response.ErrorCode = int16(errors.NotCoordinator)

		return response, nil
	}

	resp, err := c.franzKafkaClient.Broker(coordinatorId).RetriableRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	heartBeatResponse := resp.(*kmsg.HeartbeatResponse)

	if heartBeatResponse.ErrorCode == int16(errors.None) {
		redisGroupGenerationKey := fmt.Sprintf(GroupGenerationRedisKeyFmt, request.Group)
		err = c.redisClient.Client.Watch(ctx, func(tx *redis.Tx) error {
			generationId, err := tx.Get(ctx, redisGroupGenerationKey).Int64()
			if err != nil && err != redis.Nil {
				return err
			}

			if err == redis.Nil || generationId != int64(request.Generation) {
				heartBeatResponse = kmsg.NewPtrHeartbeatResponse()
				heartBeatResponse.ErrorCode = int16(errors.IllegalGeneration)
				return nil
			}

			return tx.Expire(ctx, redisGroupGenerationKey, 7*24*time.Hour).Err()
		})
	}

	return heartBeatResponse, err
}

func (c *Controller) base64Topic(topic string) string {
	return base64.StdEncoding.EncodeToString([]byte(topic))
}

func (c *Controller) OffsetFetch(group, topic string, partition int32) (int64, error) {
	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer redisContextCancel()

	redisGroupOffsetKey := fmt.Sprintf(GroupTopicPartitionOffsetRedisKeyFmt, group, c.base64Topic(topic), partition)

	var offset int64
	err := c.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		var err error
		offset, err = tx.Get(redisContext, redisGroupOffsetKey).Int64()

		if err == redis.Nil {
			return nil
		}

		return err
	}, redisGroupOffsetKey)

	return offset, err
}

func (c *Controller) OffsetFetchAllTopics(group string) (map[string]map[int32]int64, error) {
	offsets := make(map[string]map[int32]int64)

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer redisContextCancel()

	redisGroupOffsetKeyPrefix := fmt.Sprintf(GroupTopicOffsetRedisKeyFmt+"-", group)

	err := c.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		scan := tx.Scan(redisContext, 0, fmt.Sprintf("%s*", redisGroupOffsetKeyPrefix), 0).Iterator()

		for scan.Next(redisContext) {
			key := scan.Val()
			offset, err := tx.Get(redisContext, key).Int64()

			if err == redis.Nil {
				continue
			}

			if err != nil {
				return err
			}

			topicParts := strings.Split(strings.TrimPrefix(key, redisGroupOffsetKeyPrefix), "-partition-")

			topicByte, err := base64.StdEncoding.DecodeString(topicParts[0])
			if err != nil {
				return fmt.Errorf("err decoding base64 topic name: %s: %w", topicParts[0], err)
			}
			topic := string(topicByte)

			partition, err := strconv.Atoi(topicParts[1])
			if err != nil {
				return err
			}

			if offsets[topic] == nil {
				offsets[topic] = make(map[int32]int64)
			}
			offsets[topic][int32(partition)] = offset
		}

		return nil
	})

	return offsets, err
}

func (c *Controller) OffsetCommit(group, topic string, groupGenerationId, partition int32, offset int64) (errors.KafkaError, error) {
	// TODO: to expire these, every 5 minutes do a `SCAN MATCH group-offset-*` and see if a generation exists, if it doesn't delete it

	kafkaError := errors.None

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer redisContextCancel()

	// TODO: if resetting offsets do we need to check if the group is empty?
	//  kafka java cli tools do it themselves but do we still need to check?

	redisGroupGenerationKey := fmt.Sprintf(GroupGenerationRedisKeyFmt, group)
	redisGroupOffsetKey := fmt.Sprintf(GroupTopicPartitionOffsetRedisKeyFmt, group, c.base64Topic(topic), partition)
	err := c.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		// generation will be -1 when we are resetting offsets
		if groupGenerationId != -1 {
			generationId, err := tx.Get(redisContext, redisGroupGenerationKey).Int64()
			if err != nil && err != redis.Nil {
				return err
			}

			if err == redis.Nil || generationId != int64(groupGenerationId) {
				kafkaError = errors.IllegalGeneration
				return nil
			}
		}

		_, err := tx.TxPipelined(redisContext, func(pipeliner redis.Pipeliner) error {
			return tx.Set(redisContext, redisGroupOffsetKey, offset, 0).Err()
		})

		return err
	}, redisGroupGenerationKey, redisGroupOffsetKey)

	return kafkaError, err
}

func (c *Controller) DescribeGroup(group string) (*kmsg.DescribeGroupsResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	coordinatorId, err := c.groupCoordinator(ctx, group)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		response := kmsg.NewPtrDescribeGroupsResponse()

		responseGroup := kmsg.NewDescribeGroupsResponseGroup()
		responseGroup.ErrorCode = int16(errors.NotCoordinator)
		responseGroup.Group = group

		response.Groups = append(response.Groups, responseGroup)

		return response, nil
	}

	request := kmsg.NewPtrDescribeGroupsRequest()
	request.Groups = append(request.Groups, group)

	resp, err := c.franzKafkaClient.Broker(coordinatorId).RetriableRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*kmsg.DescribeGroupsResponse), nil
}
