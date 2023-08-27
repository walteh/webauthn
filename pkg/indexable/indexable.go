package indexable

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

type TableResolver interface {
	// SetTableName(string, string) bool
	ResolveTable(string) (string, bool)
	SetTable(string, string) (string, bool)

	IsInitialized() bool
}

type Indexable interface {
	Identifiable
	PrimaryIndex() *DynamoDBIndex
	SecondaryIndexes() map[string]*DynamoDBIndex
	ResolvableTableName() string
	IsWorthy(string, string) bool
	ToProtobuf() proto.Message
	Combine([]Indexable) error
}

type DynamoDBIndex struct {
	HashKey    string
	HashValue  types.AttributeValue
	RangeKey   string
	RangeValue types.AttributeValue
	Name       string
}

func IndexableMarshalIndexes(i Indexable, attrs map[string]types.AttributeValue) {
	primary := i.PrimaryIndex()

	attrs[primary.HashKey] = primary.HashValue

	if primary.RangeValue != nil && primary.RangeKey != "" {
		attrs[primary.RangeKey] = primary.RangeValue
	}

	for _, index := range i.SecondaryIndexes() {
		if index.HashValue != nil && index.HashKey != "" {
			attrs[index.HashKey] = index.HashValue
			if index.RangeValue != nil && index.RangeKey != "" {
				attrs[index.RangeKey] = index.RangeValue
			}
		}
	}
}

func indexableFindZlmLabel(fields ...reflect.StructField) string {
	label := "zlm."
	for i, f := range fields {
		label += strings.Split(f.Tag.Get("dynamodbav"), ",")[0]
		if i < len(fields)-1 {
			label += "."
		}
	}
	return label
}

func IndexableFindZlm(attrs map[string]types.AttributeValue, fields ...reflect.StructField) float64 {
	label := indexableFindZlmLabel(fields...)
	if v, ok := attrs[label]; ok {
		if n := v.(*types.AttributeValueMemberN); n != nil {
			res, _ := strconv.ParseFloat(n.Value, 64)
			return res
		}
	}
	return 0
}

func IndexableFindZlmArray(attrs map[string]types.AttributeValue, fields ...reflect.StructField) []float64 {
	label := indexableFindZlmLabel(fields...)
	res := []float64{}
	if v, ok := attrs[label]; ok {
		if n := v.(*types.AttributeValueMemberL); n != nil {
			for _, v := range n.Value {
				if n := v.(*types.AttributeValueMemberN); n != nil {
					p, _ := strconv.ParseFloat(n.Value, 64)
					res = append(res, p)
				}
			}
		}
	}
	return res
}

func IndexableApplyIndexes(i Indexable, attrs *types.Update) {
	primary := i.PrimaryIndex()
	if attrs.Key == nil {
		attrs.Key = make(map[string]types.AttributeValue)
	}
	attrs.Key[primary.HashKey] = primary.HashValue
	if primary.RangeValue != nil && primary.RangeKey != "" {
		attrs.Key[primary.RangeKey] = primary.RangeValue
	}
}

/*///////////////////////////////////////////////////////////////////
 ///                     	DynamoDBIndex 						 ///
///////////////////////////////////////////////////////////////////*/

func DynamoDBScalarTypeToAttributeValueMember(s types.ScalarAttributeType) reflect.Type {
	switch s {
	case types.ScalarAttributeTypeS:
		return reflect.TypeOf(&types.AttributeValueMemberS{})
	case types.ScalarAttributeTypeN:
		return reflect.TypeOf(&types.AttributeValueMemberN{})
	case types.ScalarAttributeTypeB:
		return reflect.TypeOf(&types.AttributeValueMemberB{})
	default:
		panic(fmt.Sprintf("unknown type %s", s))
	}
}

func (me *DynamoDBIndex) MatchesKeySchema(scheme []types.KeySchemaElement, attrs []types.AttributeDefinition) (bool, string) {

	attrvalue := func(s string) types.AttributeDefinition {
		for _, attr := range attrs {
			if *attr.AttributeName == s {
				return attr
			}
		}
		return types.AttributeDefinition{}
	}

	var hash types.KeySchemaElement
	var rangee types.KeySchemaElement

	for _, key := range scheme {
		if key.KeyType == types.KeyTypeHash {
			hash = key
		} else if key.KeyType == types.KeyTypeRange {
			rangee = key
		}
	}

	attrv := attrvalue(*hash.AttributeName)

	if me.HashKey != *hash.AttributeName {
		return false, fmt.Sprintf("hash key does not match %s != %s", me.HashKey, *hash.AttributeName)
	}

	if DynamoDBScalarTypeToAttributeValueMember(attrv.AttributeType) != reflect.TypeOf(me.HashValue) {
		return false, fmt.Sprintf("hash key type does not match (tf:%s) %s != (go:%s) %s", *attrv.AttributeName, attrv.AttributeType, me.HashKey, reflect.TypeOf(me.HashValue))
	}

	if len(scheme) > 1 {
		attrv := attrvalue(*rangee.AttributeName)
		if me.RangeKey != *rangee.AttributeName {
			return false, fmt.Sprintf("range key does not match %s != %s", me.RangeKey, *rangee.AttributeName)
		}
		if DynamoDBScalarTypeToAttributeValueMember(attrv.AttributeType) != reflect.TypeOf(me.RangeValue) {
			return false, fmt.Sprintf("range key type does not match (tf:%s) %s != (go:%s) %s", *attrv.AttributeName, attrv.AttributeType, me.RangeKey, reflect.TypeOf(me.RangeValue))
		}
	}
	return true, ""
}

func (me *DynamoDBIndex) MatchesGSI(gsi types.GlobalSecondaryIndex, attrs []types.AttributeDefinition) (bool, string) {

	if gsi.IndexName == nil {
		return false, "index name is nil - " + me.Name
	}
	if me.Name != *gsi.IndexName {
		return false, fmt.Sprintf("index name does not match %s != %s", me.Name, *gsi.IndexName)
	}
	return me.MatchesKeySchema(gsi.KeySchema, attrs)
}

func (me *DynamoDBIndex) Query(idx Indexable) *dynamodb.QueryInput {

	var idxname *string
	values := make(map[string]types.AttributeValue)
	names := make(map[string]string)

	if me.Name != "" {
		idxname = aws.String(me.Name)
	}

	values[":hk"] = me.HashValue
	names["#hk"] = me.HashKey
	condition := "#hk = :hk"

	if me.RangeValue != nil {
		values[":rk"] = me.RangeValue
		names["#rk"] = me.RangeKey
		condition = condition + " AND #rk = :rk"
	}

	// tb, ok := resolver.ResolveTable(idx.ResolvableTableName())
	// if !ok {
	// 	tb = idx.ResolvableTableName()
	// }

	return &dynamodb.QueryInput{
		TableName:                 aws.String(idx.ResolvableTableName()),
		KeyConditionExpression:    aws.String(condition),
		ExpressionAttributeValues: values,
		ExpressionAttributeNames:  names,
		IndexName:                 idxname,
	}

}

/*///////////////////////////////////////////////////////////////////
 ///                     	HELPERS 							 ///
///////////////////////////////////////////////////////////////////*/

func (me *DynamoDBIndex) UpdateKey() map[string]types.AttributeValue {
	// only use primary index for updating
	if me.Name != "" {
		return map[string]types.AttributeValue{}
	}

	if me.RangeValue == nil {
		return map[string]types.AttributeValue{
			me.HashKey: me.HashValue,
		}
	}

	return map[string]types.AttributeValue{
		me.HashKey:  me.HashValue,
		me.RangeKey: me.RangeValue,
	}
}

func IndexableAttributeName(I Indexable, f func(s string) bool) string {
	st := reflect.TypeOf(I)
	field, ok := st.FieldByNameFunc(f)
	if !ok {
		return ""
	}
	tag := field.Tag.Get("dynamodbav")
	if tag != "" {
		return tag
	}
	return ""
}

func IndexablePut(I Indexable, exists bool) *types.Put {

	itm, err := attributevalue.MarshalMap(I)
	if err != nil {
		panic(err)
	}

	condition := "attribute_not_exists"

	if exists {
		condition = "attribute_exists"
	}

	q := I.PrimaryIndex().Query(I)

	replacer := strings.NewReplacer("#rk = :rk", condition+"(#rk)", "#hk = :hk", condition+"(#hk)")

	return &types.Put{
		TableName:                q.TableName,
		Item:                     itm,
		ConditionExpression:      aws.String(replacer.Replace(*q.KeyConditionExpression)),
		ExpressionAttributeNames: q.ExpressionAttributeNames,
		// ExpressionAttributeValues:           q.ExpressionAttributeValues,
		ReturnValuesOnConditionCheckFailure: types.ReturnValuesOnConditionCheckFailureAllOld,
	}

}

func NewTableResolver[I TableResolver](tabs []string, TR I) (I, error) {

	for _, t := range tabs {
		if t == "" {
			continue
		}
		if _, ok := TR.SetTable(string(t), t); !ok {
			return TR, errors.New("failed to set table name")
		}
	}

	if !TR.IsInitialized() {
		return TR, errors.New("missing table name")
	}

	return TR, nil
}

func NewTableResolverFromMap[I TableResolver, O any](tabs map[string]O, TR I) (I, error) {
	keys := []string{}
	for k := range tabs {
		keys = append(keys, k)
	}
	return NewTableResolver(keys, TR)
}

func MustNewTableResolver[I TableResolver](tabs []string, TR I) I {
	tab, err := NewTableResolver(tabs, TR)
	if err != nil {
		panic(err)
	}
	return tab
}

/*///////////////////////////////////////////////////////////////////
 ///                     SINGLE TABLE RESOLVER 					 ///
///////////////////////////////////////////////////////////////////*/

var _ TableResolver = (*SingleTableResolver)(nil)

type SingleTableResolver struct {
	table string
}

func NewSingleTableResolver(name string) *SingleTableResolver {
	return &SingleTableResolver{name}
}

func (me *SingleTableResolver) IsInitialized() bool {
	return me.table != ""
}

func (me *SingleTableResolver) ResolveTable(name string) (string, bool) {
	return me.SetTable(name, "")
}

func NilTableResolver(me TableResolver, name *string) (string, bool) {
	if name == nil {
		// returning true so that the caller can continue
		return "invalid", true
	}
	return me.ResolveTable(*name)
}

func (me *SingleTableResolver) SetTable(name, newval string) (string, bool) {
	if newval != "" {
		me.table = newval
	}

	if me.table == "" {
		return "", false
	}

	return me.table, true
}

/*///////////////////////////////////////////////////////////////////
 ///                     	UPDATES 							 ///
///////////////////////////////////////////////////////////////////*/

type DynamoDBSet struct {
	Field     []reflect.StructField
	Operation string
}

type DynamoDBCondition struct {
	Field      []reflect.StructField
	Expression string
	OR         *DynamoDBCondition
	AND        *DynamoDBCondition
}

type DynamoDBRemove struct {
	Field []reflect.StructField
}

func IndexableResolveTablesForTransactWriteItem(I Indexable, TR TableResolver, updates *types.TransactWriteItem) *types.TransactWriteItem {
	if updates.Update != nil {
		table, ok := NilTableResolver(TR, updates.Update.TableName)
		if !ok {
			table = *updates.Update.TableName
		}
		updates.Update.TableName = aws.String(table)
	} else if updates.Put != nil {
		table, ok := NilTableResolver(TR, updates.Put.TableName)
		if !ok {
			table = *updates.Put.TableName
		}
		updates.Put.TableName = aws.String(table)
	} else if updates.Delete != nil {
		table, ok := NilTableResolver(TR, updates.Delete.TableName)
		if !ok {
			table = *updates.Delete.TableName
		}
		updates.Delete.TableName = aws.String(table)
	} else if updates.ConditionCheck != nil {
		table, ok := NilTableResolver(TR, updates.ConditionCheck.TableName)
		if !ok {
			table = *updates.ConditionCheck.TableName
		}
		updates.ConditionCheck.TableName = aws.String(table)
	}

	return updates
}

func IndexableResolveTablesForTransactWriteItemsInput(I Indexable, TR TableResolver, updates *dynamodb.TransactWriteItemsInput) *dynamodb.TransactWriteItemsInput {
	for _, u := range updates.TransactItems {
		IndexableResolveTablesForTransactWriteItem(I, TR, &u)
	}
	return updates
}

func ApplyConditions(res *types.Update, cond []DynamoDBCondition) {

	if res.ExpressionAttributeNames == nil {
		res.ExpressionAttributeNames = map[string]string{}
	}

	cache := map[string]string{}

	for k, v := range res.ExpressionAttributeNames {
		cache[v] = k
	}

	cacheorarg := func(arg string) string {
		if _, ok := cache[arg]; ok && cache[arg] != "" {
			return cache[arg]
		} else {
			cache[arg] = fmt.Sprintf("#arg_%d", len(cache))
		}
		return cache[arg]
	}

	if len(cond) == 0 {
		return
	}

	var recurse func(int, *DynamoDBCondition) string
	recurse = func(idx int, c *DynamoDBCondition) string {
		expr := ""
		for y, f := range c.Field {
			need := strings.Split(f.Tag.Get("dynamodbav"), ",")[0]
			arg := cacheorarg(need)
			res.ExpressionAttributeNames[arg] = need
			expr += arg
			if y != len(c.Field)-1 {
				expr += "."
				continue
			}
		}

		expr = strings.Replace(c.Expression, "%s", expr, -1)

		if c.OR != nil {
			return "(" + expr + " OR " + recurse(idx, c.OR) + ")"
		}

		if c.AND != nil {
			return "(" + expr + " AND " + recurse(idx, c.AND) + ")"
		}

		return expr
	}

	update := ""

	for i, c := range cond {
		if i != 0 {
			update += " AND "
		}
		update += recurse(i, &c)
	}

	res.ConditionExpression = aws.String(update)
}

func ApplyValues(res *types.Update, values map[string]types.AttributeValue) {

	if len(values) == 0 {
		return
	}

	if res.ExpressionAttributeValues == nil {
		res.ExpressionAttributeValues = map[string]types.AttributeValue{}
	}

	for k, v := range values {
		res.ExpressionAttributeValues[k] = v
	}

}

func ApplySet(res *types.Update, updates []DynamoDBSet) {

	if len(updates) == 0 {
		return
	}

	if res.ExpressionAttributeNames == nil {
		res.ExpressionAttributeNames = map[string]string{}
	}

	cache := map[string]string{}

	for k, v := range res.ExpressionAttributeNames {
		cache[v] = k
	}

	cacheorarg := func(arg string) string {
		if _, ok := cache[arg]; ok && cache[arg] != "" {
			return cache[arg]
		} else {
			cache[arg] = fmt.Sprintf("#arg_%d", len(cache))
		}
		return cache[arg]
	}

	// build update expression
	update := "SET "
	for i, u := range updates {

		expr := ""
		for y, f := range u.Field {
			need := strings.Split(f.Tag.Get("dynamodbav"), ",")[0]
			arg := cacheorarg(need)
			res.ExpressionAttributeNames[arg] = need
			expr += arg
			if y != len(u.Field)-1 {
				expr += "."
				continue
			}
		}

		update += expr + " " + strings.Replace(u.Operation, "%s", expr, -1)

		if i != len(updates)-1 && len(updates) > 1 {
			update += ", "
		}

	}

	res.UpdateExpression = aws.String(update)

}

func ApplyRemove(res *types.Update, updates []DynamoDBRemove) {

	if len(updates) == 0 {
		return
	}

	// expressionKeys := make(map[string]string)

	// // build update expression
	// update := "REMOVE "
	// for i, u := range updates {

	// 	expr := ""
	// 	for y, f := range u.Field {
	// 		arg := fmt.Sprintf("#arg_%d_%d", i, y)
	// 		raw := f.Tag.Get("dynamodbav")
	// 		expressionKeys[arg] = strings.Split(raw, ",")[0]
	// 		expr += arg
	// 		if y != len(u.Field)-1 {
	// 			expr += "."
	// 			continue
	// 		}
	// 	}

	// 	update += expr + ", "
	// }

	// res.UpdateExpression = aws.String(strings.TrimSuffix(update, ", "))

	if res.ExpressionAttributeNames == nil {
		res.ExpressionAttributeNames = map[string]string{}
	}

	cache := map[string]string{}

	for k, v := range res.ExpressionAttributeNames {
		cache[v] = k
	}

	cacheorarg := func(arg string) string {
		if _, ok := cache[arg]; ok && cache[arg] != "" {
			return cache[arg]
		} else {
			cache[arg] = fmt.Sprintf("#arg_%d", len(cache))
		}
		return cache[arg]
	}

	// build update expression
	update := "REMOVE "
	for i, u := range updates {

		expr := ""
		for y, f := range u.Field {
			need := strings.Split(f.Tag.Get("dynamodbav"), ",")[0]
			arg := cacheorarg(need)
			res.ExpressionAttributeNames[arg] = need
			expr += arg
			if y != len(u.Field)-1 {
				expr += "."
				continue
			}
		}

		update += expr + " "

		if i != len(updates)-1 && len(updates) > 1 {
			update += ", "
		}

	}

	res.UpdateExpression = aws.String(update)

}

func IndexableSet(ctx context.Context, idx Indexable, lm LastModifier, to attributevalue.Marshaler, rs ...reflect.StructField) (txs []DynamoDBTransaction) {

	upd := &types.Update{}

	IndexableApplyIndexes(idx, upd)

	g := &DynamoDBUpdate{
		Update:        upd,
		preconditions: []DynamoDBPrecondition{},
	}

	txs = append(txs, g)

	toset := []DynamoDBSet{
		{
			Field:     rs,
			Operation: "= :A",
		}}

	tocond := []DynamoDBCondition{}

	toval := map[string]types.AttributeValue{}

	toremove := []DynamoDBRemove{}

	if lm != nil {
		toset = append(toset, DynamoDBSet{
			Field: []reflect.StructField{
				{
					Tag: reflect.StructTag(`dynamodbav:"ttl"`),
				},
			},
			Operation: "= :TTL",
		})
		if lm.IsPermanent() {
			toval[":TTL"] = NULL(true)
		} else {
			toval[":TTL"] = N(time.Now().Add(time.Hour).Unix())
		}

		lmfieldName := fmt.Sprintf(`zlm.%s`, strings.Split(rs[0].Tag.Get("dynamodbav"), ",")[0])

		lmfield := []reflect.StructField{{
			Tag: reflect.StructTag(fmt.Sprintf(`dynamodbav:"%s"`, lmfieldName))}}

		toset = append(toset, DynamoDBSet{
			Field:     lmfield,
			Operation: "= :LM",
		})

		tocond = append(tocond, DynamoDBCondition{
			Field:      lmfield,
			Expression: "attribute_not_exists(%s)",
			OR: &DynamoDBCondition{
				Field:      lmfield,
				Expression: "%s < :LM",
			},
		})

		g.preconditions = append(g.preconditions, MakeLastModifierPrecondition(idx, lmfieldName, lm))

		toval[":LM"] = F(lm.At())

	}

	attr, err := attributevalue.Marshal(to)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Interface("attribtue_marshaler", to).Msg("failed to marshal attribute")
		return
	}

	toval[":A"] = attr

	ApplySet(upd, toset)

	ApplyConditions(upd, tocond)

	ApplyValues(upd, toval)

	ApplyRemove(upd, toremove)

	upd.TableName = aws.String(idx.ResolvableTableName())

	return
}

func IndexableID(x Indexable) (string, error) {

	ranger := ""
	hasher := ""

	switch r := x.PrimaryIndex().RangeValue.(type) {
	case *types.AttributeValueMemberS:
		ranger = r.Value
	case *types.AttributeValueMemberN:
		ranger = r.Value
	case *types.AttributeValueMemberB:
		ranger = base64.StdEncoding.EncodeToString(r.Value)
	default:
		ranger = ""
	}

	switch h := x.PrimaryIndex().HashValue.(type) {
	case *types.AttributeValueMemberS:
		hasher = h.Value
	case *types.AttributeValueMemberN:
		hasher = h.Value
	case *types.AttributeValueMemberB:
		hasher = base64.StdEncoding.EncodeToString(h.Value)
	default:
		return "", fmt.Errorf("unknown hash type %T", h)
	}

	if ranger == "" {
		return hasher, nil
	}

	return hasher + ":" + ranger, nil
}

func IndexableIncrement(ctx context.Context, idx Indexable, lm LastModifier, amount *types.AttributeValueMemberN, rs ...reflect.StructField) (txs []DynamoDBTransaction) {

	upd := &types.Update{}
	IndexableApplyIndexes(idx, upd)

	g := &DynamoDBUpdate{
		Update:        upd,
		preconditions: []DynamoDBPrecondition{},
	}

	txs = append(txs, g)

	lmfieldName := fmt.Sprintf(`zlm.%s`, strings.Split(rs[0].Tag.Get("dynamodbav"), ",")[0])

	lmfield := []reflect.StructField{{
		Tag: reflect.StructTag(fmt.Sprintf(`dynamodbav:"%s"`, lmfieldName))}}

	ApplySet(upd, []DynamoDBSet{
		{
			Field:     rs,
			Operation: "= if_not_exists(%s, :ZERO) + :A",
		},
		{
			Field:     lmfield,
			Operation: "= list_append(if_not_exists(%s, :LMSE), :LMS)",
		},
		{
			Field: []reflect.StructField{
				{
					Tag: reflect.StructTag(`dynamodbav:"ttl"`),
				},
			},
			Operation: "= :TTL",
		},
	})

	ApplyConditions(upd, []DynamoDBCondition{
		{
			Field:      lmfield,
			Expression: "attribute_not_exists(%s)",
			OR: &DynamoDBCondition{
				Field:      lmfield,
				Expression: "NOT contains(%s, :LM)",
			}}})

	g.preconditions = append(g.preconditions, MakeLastModifierPrecondition(idx, lmfieldName, lm))

	toval := map[string]types.AttributeValue{
		":A":    amount,
		":LMS":  L(F(lm.At())),
		":LM":   F(lm.At()),
		":LMSE": L(),
		":ZERO": N(0),
	}

	if lm.IsPermanent() {
		toval[":TTL"] = NULL(true)
	} else {
		toval[":TTL"] = N(time.Now().Add(time.Hour).Unix())
	}

	ApplyValues(upd, toval)

	upd.TableName = aws.String(idx.ResolvableTableName())

	upd.ReturnValuesOnConditionCheckFailure = types.ReturnValuesOnConditionCheckFailureAllOld

	return
}

func IndexableUpsert(ctx context.Context, idx Indexable, lm LastModifier) (txs []DynamoDBTransaction) {

	upd := &types.Update{}
	IndexableApplyIndexes(idx, upd)

	root := &DynamoDBUpdate{
		Update:        upd,
		preconditions: []DynamoDBPrecondition{},
	}

	txs = append(txs, root)

	rs := []reflect.StructField{{Tag: reflect.StructTag("dynamodbav:\"index\"")}}

	lmfieldName := fmt.Sprintf(`zlm.%s`, strings.Split(rs[0].Tag.Get("dynamodbav"), ",")[0])

	lmfield := []reflect.StructField{{
		Tag: reflect.StructTag(fmt.Sprintf(`dynamodbav:"%s"`, lmfieldName))}}

	ApplySet(upd, []DynamoDBSet{
		{
			Field:     lmfield,
			Operation: "= if_not_exists(%s, :LM)",
		},
	})

	ApplyConditions(upd, []DynamoDBCondition{
		{
			Field:      lmfield,
			Expression: "attribute_not_exists(%s)",
			OR: &DynamoDBCondition{
				Field:      lmfield,
				Expression: "%s < :LM",
			},
		}})

	root.preconditions = append(root.preconditions, MakeLastModifierPrecondition(idx, lmfieldName, lm))

	ApplyValues(upd, map[string]types.AttributeValue{
		":LM": F(lm.At()),
	})

	upd.TableName = aws.String(idx.ResolvableTableName())

	return
}

func IndexableSetConstant(ctx context.Context, idx Indexable, to attributevalue.Marshaler, rs reflect.StructField) (txs []DynamoDBTransaction) {

	upd := &types.Update{}

	IndexableApplyIndexes(idx, upd)

	g := &DynamoDBUpdate{
		Update:        upd,
		preconditions: []DynamoDBPrecondition{},
	}

	txs = append(txs, g)

	toset := []DynamoDBSet{
		{
			Field:     []reflect.StructField{rs},
			Operation: "= :A",
		}}

	tocond := []DynamoDBCondition{}

	toval := map[string]types.AttributeValue{}

	toremove := []DynamoDBRemove{}

	field := strings.Split(rs.Tag.Get("dynamodbav"), ",")[0]

	g.preconditions = append(g.preconditions, MakeDoesNotExistPrecondition(idx, field))

	attr, err := attributevalue.Marshal(to)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Interface("attribtue_marshaler", to).Msg("failed to marshal attribute")
		return
	}

	toval[":A"] = attr

	ApplySet(upd, toset)

	ApplyConditions(upd, tocond)

	ApplyValues(upd, toval)

	ApplyRemove(upd, toremove)

	upd.TableName = aws.String(idx.ResolvableTableName())

	return
}
