package models

type ACLOperation int8

const (
	ACLOperationAny ACLOperation = iota + 1
	ACLOperationAll
	ACLOperationRead
	ACLOperationWrite
	ACLOperationCreate
	ACLOperationDelete
	ACLOperationAlter
	ACLOperationDescribe
	ACLOperationClusterAction
	ACLOperationDescribeConfigs
	ACLOperationAlterConfigs
	ACLOperationIdempotentWrite
)

type ACLResourceType int8

const (
	ACLResourceTypeAny ACLResourceType = iota + 1
	ACLResourceTypeTopic
	ACLResourceTypeGroup
	ACLResourceTypeCluster
	ACLResourceTypeTransactionalID
	ACLResourceTypeDelegationToken
)

type ACLPatternType int8

const (
	ACLPatternTypeAny ACLPatternType = iota + 1
	ACLPatternTypeMatch
	ACLPatternTypeLiteral
	ACLPatternTypePrefixed
)

type ACLPermission int8

const (
	ACLPermissionAny = iota + 1
	ACLPermissionDeny
	ACLPermissionAllow
)

type ACL struct {
	Operation    ACLOperation
	ResourceType ACLResourceType
	PatternType  ACLPatternType
	ResourceName string
	Principal    string
	Permission   ACLPermission `hash:"ignore"` // ignore this field while hashing so hash doesn't change
}
