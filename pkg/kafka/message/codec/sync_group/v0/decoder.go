package v0

import (
	"fmt"

	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {
	msg := &v0.Request{}

	var err error
	if msg.GroupID, err = reader.String(); err != nil {
		fmt.Println("HERE GROUPID!")
		return nil, err
	}

	if msg.GenerationID, err = reader.Int32(); err != nil {
		fmt.Println("HERE GENERATIONID!")
		return nil, err
	}

	if msg.MemberID, err = reader.String(); err != nil {
		fmt.Println("HERE MEMBERID!")
		return nil, err
	}

	fmt.Printf("Request Before Fail: %v\n", *msg)

	var assignmentsLength int32
	if assignmentsLength, err = reader.ArrayLength(); err != nil {
		fmt.Println("HERE ASSIGNMENT LENGTH!")
		return nil, err
	}

	for i := int32(0); i < assignmentsLength; i++ {
		groupAssignment := v0.GroupAssignments{}

		if groupAssignment.MemberID, err = reader.String(); err != nil {
			fmt.Println("HERE ASSIGNMENT ID!")
			return nil, err
		}

		if groupAssignment.Assignment, err = reader.Bytes(); err != nil {
			fmt.Println("HERE ASSIGNMENT!")
			return nil, err
		}

		msg.Assignments = append(msg.Assignments, groupAssignment)
	}

	return msg, nil
}
