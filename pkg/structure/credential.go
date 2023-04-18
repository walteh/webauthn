package structure

import (
	"git.nugg.xyz/go-sdk/x"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var _ x.Indexable = (*Credential)(nil)

type Credential struct {
	Id string `dynamo:"pk,hash"`
}

func (c *Credential) ResolvableTableName() string {
	return "ceremony"
}

func (c *Credential) IsWorthy(primary, secondary string) bool {
	return primary == "ceremony"
}

func (c *Credential) PrimaryIndex() *x.DynamoDBIndex {
	return &x.DynamoDBIndex{
		HashKey:  "pk",
		RangeKey: "sk",
	}
}

func (c *Credential) SecondaryIndexes() map[string]*x.DynamoDBIndex {
	return map[string]*x.DynamoDBIndex{}
}

func (c *Credential) ID() string {
	return ""
}

func (c *Credential) Combine([]x.Indexable) error {
	return nil
}

func (c *Credential) ToProtobuf() protoreflect.ProtoMessage {
	return nil
}

func NewCredentialQueryable(id string) *Credential {
	return &Credential{
		Id: id,
	}
}
