package structure

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"git.nugg.xyz/go-sdk/env"
	"git.nugg.xyz/go-sdk/x"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

var _ x.TableResolver = (*Tables)(nil)

type Tables struct {
	Credential string
	Ceremony   string
}

func ParseUnknown(ctx context.Context, name string, tr x.TableResolver, keys map[string]types.AttributeValue, img map[string]types.AttributeValue) (x.Indexable, error) {
	var primary string
	var secondary string

	log := zerolog.Ctx(ctx).With().Str("table", name).Logger()

	if pk, ok := keys["pk"].(*types.AttributeValueMemberS); ok {
		primary = pk.Value
	} else {
		return nil, fmt.Errorf("pk not found")
	}

	if sk, ok := keys["sk"].(*types.AttributeValueMemberS); ok {
		secondary = sk.Value
	} else if sk, ok := keys["sk"].(*types.AttributeValueMemberN); ok {
		secondary = sk.Value
	} else {
		secondary = ""
	}

	log.Trace().Str("primary", primary).Str("secondary", secondary).Msg("key parse success")

	for _, v := range []x.Indexable{
		&Credential{},
		&Ceremony{},
	} {
		specific, ok := tr.ResolveTable(v.ResolvableTableName())
		if ok && specific == name {
			log.Trace().Str("table", name).Str("type", reflect.TypeOf(v).String()).Msg("table match")
			if v.IsWorthy(primary, secondary) {
				if err := attributevalue.UnmarshalMap(img, &v); err != nil {
					return nil, err
				} else {
					log.Debug().Any("before", img).Any("after", v).Msg("table worthy")
					return v, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("unable to parse unknown table %s", name)
}

func (me *Tables) ResolveTable(name string) (str string, ok bool) {
	return me.SetTable(name, "")
}
func (me *Tables) SetTable(name string, newval string) (str string, ok bool) {

	name = strings.ToLower(name)
	name = strings.Replace(name, "-", "", -1)
	name = strings.Replace(name, "_", "", -1)

	resolve := func(str string) (string, string) {
		if newval == "" {
			return str, str
		}
		return newval, newval
	}

	ok = true

	if strings.Contains(name, (&Credential{}).ResolvableTableName()) {
		str, me.Credential = resolve(me.Credential)
	} else if strings.Contains(name, (&Ceremony{}).ResolvableTableName()) {
		str, me.Ceremony = resolve(me.Ceremony)
	} else {
		ok = false
	}
	return
}

func (me *Tables) IsInitialized() bool {
	for _, v := range me.AsArray() {
		if v == "" {
			return false
		}
	}
	return true
}

func (me *Tables) AsArray() []string {
	v := reflect.ValueOf(me).Elem()
	arr := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		arr[i] = v.Field(i).String()
	}
	return arr
}

func MustGetTables(fmtr string, key string) *Tables {
	c := func(str string) string {
		return strings.Replace(fmtr, key, str, -1)
	}
	return &Tables{
		Ceremony:   (env.MustGet(c("CEREMONY"))),
		Credential: (env.MustGet(c("CREDENTIAL"))),
	}
}

func DummyTables() *Tables {
	return &Tables{
		Credential: ("credentials"),
		Ceremony:   ("ceremonies"),
	}
}

func DummyTablesArray() []string {
	return DummyTables().AsArray()
}
