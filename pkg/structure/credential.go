package structure

import (
	"reflect"

	"github.com/walteh/webauthn/pkg/indexable"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var _ indexable.Indexable = (*Credential)(nil)

type Credential struct {
	Id           string `dynamo:"pk,hash"`
	SessionCount int    `dynamo:"sk,range"`
}

/*///////////////////////////////////////////////////////////////////
 ///                     STRUCT FIELDS 							 ///
///////////////////////////////////////////////////////////////////*/

func CredentialSessionCountStructField() reflect.StructField {
	tmp := Credential{}
	_ = tmp.SessionCount
	f, _ := reflect.TypeOf(tmp).FieldByName("SessionCount")
	return f
}

func (c *Credential) ResolvableTableName() string {
	return "credential"
}

func (c *Credential) IsWorthy(primary, secondary string) bool {
	return primary == "ceremony"
}

func (c *Credential) PrimaryIndex() *indexable.DynamoDBIndex {
	return &indexable.DynamoDBIndex{
		HashKey:  "pk",
		RangeKey: "sk",
	}
}

func (c *Credential) SecondaryIndexes() map[string]*indexable.DynamoDBIndex {
	return map[string]*indexable.DynamoDBIndex{}
}

func (c *Credential) ID() string {
	return ""
}

func (c *Credential) Combine([]indexable.Indexable) error {
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
