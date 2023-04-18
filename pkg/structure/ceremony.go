package structure

import (
	"git.nugg.xyz/go-sdk/x"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var _ x.Indexable = (*Ceremony)(nil)

type Ceremony struct {
}

func (c *Ceremony) ResolvableTableName() string {
	return "ceremony"
}

func (c *Ceremony) IsWorthy(primary, secondary string) bool {
	return primary == "ceremony"
}

func (c *Ceremony) PrimaryIndex() *x.DynamoDBIndex {
	return &x.DynamoDBIndex{
		HashKey:  "pk",
		RangeKey: "sk",
	}
}

func (c *Ceremony) SecondaryIndexes() map[string]*x.DynamoDBIndex {
	return map[string]*x.DynamoDBIndex{}
}

func (c *Ceremony) ID() string {
	return ""
}

func (c *Ceremony) Combine([]x.Indexable) error {
	return nil
}

func (c *Ceremony) ToProtobuf() protoreflect.ProtoMessage {
	return nil
}
func NewCeremonyQueryable() *Ceremony {
	return &Ceremony{}
}
