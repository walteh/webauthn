package structure

import (
	"github.com/walteh/webauthn/pkg/indexable"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var _ indexable.Indexable = (*Ceremony)(nil)

type Ceremony struct {
	PK string `dynamo:"pk"`
}

func (c *Ceremony) ResolvableTableName() string {
	return "challenge"
}

func (c *Ceremony) IsWorthy(primary, secondary string) bool {
	return primary == "challenge"
}

func (c *Ceremony) PrimaryIndex() *indexable.DynamoDBIndex {
	return &indexable.DynamoDBIndex{
		HashKey:  "pk",
		RangeKey: "sk",
	}
}

func (c *Ceremony) SecondaryIndexes() map[string]*indexable.DynamoDBIndex {
	return map[string]*indexable.DynamoDBIndex{}
}

func (c *Ceremony) ID() string {
	return ""
}

func (c *Ceremony) Combine([]indexable.Indexable) error {
	return nil
}

func (c *Ceremony) ToProtobuf() protoreflect.ProtoMessage {
	return nil
}
func NewChallengeQueryable(id string) *Ceremony {
	return &Ceremony{
		PK: id,
	}
}
