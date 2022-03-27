package franz

import (
	"context"

	"github.com/twmb/franz-go/pkg/kbin"
	"github.com/twmb/franz-go/pkg/kmsg"
)

// This file contains custom Franz Messages that are for older versions
//  Franz always auto upgrades, so we need these to support old versions
//  when the newer versions cause issues

type SyncGroupRequestGroupAssignment struct {
	// MemberID is the member this assignment is for.
	MemberID string

	// MemberAssignment is the assignment for this member. This is typically
	// of type GroupMemberAssignment.
	MemberAssignment []byte

	// UnknownTags are tags Kafka sent that we do not know the purpose of.
	UnknownTags kmsg.Tags // v4+

}

// Default sets any default fields. Calling this allows for future compatibility
// if new fields are added to SyncGroupRequestGroupAssignment.
func (v *SyncGroupRequestGroupAssignment) Default() {
}

// NewSyncGroupRequestGroupAssignment returns a default SyncGroupRequestGroupAssignment
// This is a shortcut for creating a struct and calling Default yourself.
func NewSyncGroupRequestGroupAssignment() SyncGroupRequestGroupAssignment {
	var v SyncGroupRequestGroupAssignment
	v.Default()
	return v
}

// SyncGroupRequest is issued by all group members after they receive a a
// response for JoinGroup. The group leader is responsible for sending member
// assignments with the request; all other members do not.
//
// Once the leader sends the group assignment, all members will be replied to.
type SyncGroupRequest struct {
	// Version is the version of this message used with a Kafka broker.
	Version int16

	// Group is the group ID this sync group is for.
	Group string

	// Generation is the group generation this sync is for.
	Generation int32

	// MemberID is the member ID this member is.
	MemberID string

	// InstanceID is the instance ID of this member in the group (KIP-345).
	InstanceID *string // v3+

	// GroupAssignment, sent only from the group leader, is the topic partition
	// assignment it has decided on for all members.
	GroupAssignment []SyncGroupRequestGroupAssignment

	// UnknownTags are tags Kafka sent that we do not know the purpose of.
	UnknownTags kmsg.Tags // v4+

}

func (*SyncGroupRequest) Key() int16                   { return 14 }
func (*SyncGroupRequest) MaxVersion() int16            { return 4 }
func (v *SyncGroupRequest) SetVersion(version int16)   { v.Version = version }
func (v *SyncGroupRequest) GetVersion() int16          { return v.Version }
func (v *SyncGroupRequest) IsFlexible() bool           { return v.Version >= 4 }
func (v *SyncGroupRequest) IsGroupCoordinatorRequest() {}
func (v *SyncGroupRequest) ResponseKind() kmsg.Response {
	return &kmsg.SyncGroupResponse{Version: v.Version}
}

// RequestWith is requests v on r and returns the response or an error.
// For sharded requests, the response may be merged and still return an error.
// It is better to rely on client.RequestSharded than to rely on proper merging behavior.
func (v *SyncGroupRequest) RequestWith(ctx context.Context, r kmsg.Requestor) (*kmsg.SyncGroupResponse, error) {
	kresp, err := r.Request(ctx, v)
	resp, _ := kresp.(*kmsg.SyncGroupResponse)
	return resp, err
}

func (v *SyncGroupRequest) AppendTo(dst []byte) []byte {
	version := v.Version
	_ = version
	isFlexible := version >= 4
	_ = isFlexible
	{
		v := v.Group
		if isFlexible {
			dst = kbin.AppendCompactString(dst, v)
		} else {
			dst = kbin.AppendString(dst, v)
		}
	}
	{
		v := v.Generation
		dst = kbin.AppendInt32(dst, v)
	}
	{
		v := v.MemberID
		if isFlexible {
			dst = kbin.AppendCompactString(dst, v)
		} else {
			dst = kbin.AppendString(dst, v)
		}
	}
	if version >= 3 {
		v := v.InstanceID
		if isFlexible {
			dst = kbin.AppendCompactNullableString(dst, v)
		} else {
			dst = kbin.AppendNullableString(dst, v)
		}
	}
	{
		v := v.GroupAssignment
		if isFlexible {
			dst = kbin.AppendCompactArrayLen(dst, len(v))
		} else {
			dst = kbin.AppendArrayLen(dst, len(v))
		}
		for i := range v {
			v := &v[i]
			{
				v := v.MemberID
				if isFlexible {
					dst = kbin.AppendCompactString(dst, v)
				} else {
					dst = kbin.AppendString(dst, v)
				}
			}
			{
				v := v.MemberAssignment
				if isFlexible {
					dst = kbin.AppendCompactBytes(dst, v)
				} else {
					dst = kbin.AppendBytes(dst, v)
				}
			}
			if isFlexible {
				dst = kbin.AppendUvarint(dst, 0+uint32(v.UnknownTags.Len()))
				dst = v.UnknownTags.AppendEach(dst)
			}
		}
	}
	if isFlexible {
		dst = kbin.AppendUvarint(dst, 0+uint32(v.UnknownTags.Len()))
		dst = v.UnknownTags.AppendEach(dst)
	}
	return dst
}

func (v *SyncGroupRequest) ReadFrom(src []byte) error {
	v.Default()
	b := kbin.Reader{Src: src}
	version := v.Version
	_ = version
	isFlexible := version >= 4
	_ = isFlexible
	s := v
	{
		var v string
		if isFlexible {
			v = b.CompactString()
		} else {
			v = b.String()
		}
		s.Group = v
	}
	{
		v := b.Int32()
		s.Generation = v
	}
	{
		var v string
		if isFlexible {
			v = b.CompactString()
		} else {
			v = b.String()
		}
		s.MemberID = v
	}
	if version >= 3 {
		var v *string
		if isFlexible {
			v = b.CompactNullableString()
		} else {
			v = b.NullableString()
		}
		s.InstanceID = v
	}
	{
		v := s.GroupAssignment
		a := v
		var l int32
		if isFlexible {
			l = b.CompactArrayLen()
		} else {
			l = b.ArrayLen()
		}
		if !b.Ok() {
			return b.Complete()
		}
		if l > 0 {
			a = make([]SyncGroupRequestGroupAssignment, l)
		}
		for i := int32(0); i < l; i++ {
			v := &a[i]
			v.Default()
			s := v
			{
				var v string
				if isFlexible {
					v = b.CompactString()
				} else {
					v = b.String()
				}
				s.MemberID = v
			}
			{
				var v []byte
				if isFlexible {
					v = b.CompactBytes()
				} else {
					v = b.Bytes()
				}
				s.MemberAssignment = v
			}
			if isFlexible {
				s.UnknownTags = internalReadTags(&b)
			}
		}
		v = a
		s.GroupAssignment = v
	}
	if isFlexible {
		s.UnknownTags = internalReadTags(&b)
	}
	return b.Complete()
}

// NewPtrSyncGroupRequest returns a pointer to a default SyncGroupRequest
// This is a shortcut for creating a new(struct) and calling Default yourself.
func NewPtrSyncGroupRequest() *SyncGroupRequest {
	var v SyncGroupRequest
	v.Default()
	return &v
}

// Default sets any default fields. Calling this allows for future compatibility
// if new fields are added to SyncGroupRequest.
func (v *SyncGroupRequest) Default() {
}

// NewSyncGroupRequest returns a default SyncGroupRequest
// This is a shortcut for creating a struct and calling Default yourself.
func NewSyncGroupRequest() SyncGroupRequest {
	var v SyncGroupRequest
	v.Default()
	return v
}

// internalReadTags reads tags in a reader and returns the tags from a
// duplicated inner kbin.Reader.
func internalReadTags(b *kbin.Reader) kmsg.Tags {
	var t kmsg.Tags
	for num := b.Uvarint(); num > 0; num-- {
		key, size := b.Uvarint(), b.Uvarint()
		t.Set(key, b.Span(int(size)))
	}
	return t
}
