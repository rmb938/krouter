package v0

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_configs/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("describe-configs-v0-handler")
	request := message.(*v0.Request)

	response := &v0.Response{}

	response.ThrottleDuration = time.Duration(0)

	for _, resource := range request.Resources {
		result := v0.DescribeConfigResultResponse{
			ResourceType: resource.ResourceType,
			ResourceName: resource.ResourceName,
		}

		if resource.ResourceType != int8(2) {
			result.ErrCode = errors.InvalidRequest
			result.ErrMessage = func(s string) *string { return &s }("Only support describing topic resources")
			response.Results = append(response.Results, result)
			continue
		}

		_, topic := broker.GetTopic(result.ResourceName)

		if topic == nil {
			result.ErrCode = errors.UnknownTopicOrPartition
			result.ErrMessage = func(s string) *string { return &s }("Topic does not exist")
			response.Results = append(response.Results, result)
			continue
		}

		if resource.ConfigurationKeys != nil {
			for _, configName := range resource.ConfigurationKeys {
				if configValue, ok := topic.Config[configName]; ok {
					configResult := v0.DescribeConfigResultConfigResponse{
						Name:      configName,
						Value:     configValue,
						ReadOnly:  false, // TODO: this
						Default:   false, // TODO: this
						Sensitive: false, // TODO: this
					}

					result.Configs = append(result.Configs, configResult)
				}
			}
		} else {
			for key, value := range topic.Config {
				configResult := v0.DescribeConfigResultConfigResponse{
					Name:      key,
					Value:     value,
					ReadOnly:  false, // TODO: this
					Default:   false, // TODO: this
					Sensitive: false, // TODO: this
				}

				result.Configs = append(result.Configs, configResult)
			}
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}
