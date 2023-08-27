package indexable

import (
	"context"
	"math"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

type DynamoDBAPI interface {
	UpdateItem(context.Context, *dynamodb.UpdateItemInput, ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	TransactGetItems(context.Context, *dynamodb.TransactGetItemsInput, ...func(*dynamodb.Options)) (*dynamodb.TransactGetItemsOutput, error)
	TransactWriteItems(context.Context, *dynamodb.TransactWriteItemsInput, ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
	Scan(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	Table() string
	BatchWriteItem(context.Context, *dynamodb.BatchWriteItemInput, ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
}

type DynamoDBAPIProvisioner interface {
	SetTable(string)
	Url() *url.URL
	DynamoDBAPI
	DescribeLimits(context.Context, *dynamodb.DescribeLimitsInput, ...func(*dynamodb.Options)) (*dynamodb.DescribeLimitsOutput, error)
	CreateTable(context.Context, *dynamodb.CreateTableInput, ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DeleteTable(context.Context, *dynamodb.DeleteTableInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error)
}

var _ DynamoDBAPI = (*defaultDynamoDBAPI)(nil)
var _ DynamoDBAPIProvisioner = (*defaultDynamoDBAPI)(nil)

type defaultDynamoDBAPI struct {
	*dynamodb.Client

	local_endpoint *url.URL
	table          string
}

func (me *defaultDynamoDBAPI) Table() string {
	return me.table
}

func (me *defaultDynamoDBAPI) SetTable(in string) {
	me.table = in
}

func (me *defaultDynamoDBAPI) Url() *url.URL {
	return me.local_endpoint
}

func NewDynamoDBAPI(c *dynamodb.Client, table string) DynamoDBAPI {
	return &defaultDynamoDBAPI{
		Client: c,
		table:  table,
	}
}

func NewDynamoDBAPIProvisioner(c *dynamodb.Client, endpoint *url.URL) DynamoDBAPIProvisioner {
	return &defaultDynamoDBAPI{
		Client:         c,
		table:          "",
		local_endpoint: endpoint,
	}
}

func (me *defaultDynamoDBAPI) TransactWriteItems(ctx context.Context, input *dynamodb.TransactWriteItemsInput, opts ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	return me.Client.TransactWriteItems(ctx, input, opts...)
}

func arrayofptrornil[I interface{}](i *I) []I {
	if i == nil {
		return nil
	}
	return []I{*i}
}

func TransactWriteItemsExpectingConditionalCheckToFail(ctx context.Context, api DynamoDBAPI, input DynamoDBTransaction, opts ...func(*dynamodb.Options)) (res *dynamodb.TransactWriteItemsOutput, err error) {
	log := zerolog.Ctx(ctx).With().Logger()

	r := input.AsTransactWriteItems()

	defer func() {
		if res != nil {
			input.ApplyConsumedCapacity(res.ConsumedCapacity...)
		}
	}()

	if len(r) == 1 {
		var err error
		if z := r[0].Put; z != nil {

			rv := types.ReturnValueNone
			if z.ReturnValuesOnConditionCheckFailure == types.ReturnValuesOnConditionCheckFailureAllOld {
				rv = types.ReturnValueAllOld
			}

			if pi, err := api.PutItem(ctx, &dynamodb.PutItemInput{
				TableName:                   z.TableName,
				Item:                        z.Item,
				ReturnConsumedCapacity:      types.ReturnConsumedCapacityTotal,
				ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
				ConditionExpression:         z.ConditionExpression,
				ExpressionAttributeNames:    z.ExpressionAttributeNames,
				ExpressionAttributeValues:   z.ExpressionAttributeValues,
				ReturnValues:                rv,
			}, opts...); err == nil {
				if pi.ConsumedCapacity != nil {
					pi.ConsumedCapacity.WriteCapacityUnits = pi.ConsumedCapacity.CapacityUnits
				}
				res = &dynamodb.TransactWriteItemsOutput{
					ConsumedCapacity: arrayofptrornil(pi.ConsumedCapacity),
					ItemCollectionMetrics: map[string][]types.ItemCollectionMetrics{
						*z.TableName: arrayofptrornil(pi.ItemCollectionMetrics),
					},
					ResultMetadata: pi.ResultMetadata,
				}
			}

		} else if z := r[0].Update; z != nil {

			if pi, err := api.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				TableName:                   z.TableName,
				Key:                         z.Key,
				UpdateExpression:            z.UpdateExpression,
				ExpressionAttributeNames:    z.ExpressionAttributeNames,
				ExpressionAttributeValues:   z.ExpressionAttributeValues,
				ConditionExpression:         z.ConditionExpression,
				ReturnConsumedCapacity:      types.ReturnConsumedCapacityTotal,
				ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
				ReturnValues:                types.ReturnValue(z.ReturnValuesOnConditionCheckFailure),
			}, opts...); err == nil {
				if pi.ConsumedCapacity != nil {
					pi.ConsumedCapacity.WriteCapacityUnits = pi.ConsumedCapacity.CapacityUnits
				}
				res = &dynamodb.TransactWriteItemsOutput{
					ConsumedCapacity: arrayofptrornil(pi.ConsumedCapacity),
					ItemCollectionMetrics: map[string][]types.ItemCollectionMetrics{
						*z.TableName: arrayofptrornil(pi.ItemCollectionMetrics),
					}, ResultMetadata: pi.ResultMetadata,
				}
			}
		} else if z := r[0].Delete; z != nil {

			if pi, err := api.DeleteItem(ctx, &dynamodb.DeleteItemInput{
				TableName:                   z.TableName,
				Key:                         z.Key,
				ConditionExpression:         z.ConditionExpression,
				ExpressionAttributeNames:    z.ExpressionAttributeNames,
				ReturnConsumedCapacity:      types.ReturnConsumedCapacityTotal,
				ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
				ReturnValues:                types.ReturnValue(z.ReturnValuesOnConditionCheckFailure),
			}, opts...); err == nil {
				if pi.ConsumedCapacity != nil {
					pi.ConsumedCapacity.WriteCapacityUnits = pi.ConsumedCapacity.CapacityUnits
				}
				res = &dynamodb.TransactWriteItemsOutput{
					ConsumedCapacity: arrayofptrornil(pi.ConsumedCapacity),
					ItemCollectionMetrics: map[string][]types.ItemCollectionMetrics{
						*z.TableName: arrayofptrornil(pi.ItemCollectionMetrics),
					},
					ResultMetadata: pi.ResultMetadata,
				}
			}
		}

		if err != nil || res != nil {
			return res, err
		}

	}

	res, err = api.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems:               r,
		ReturnConsumedCapacity:      types.ReturnConsumedCapacityTotal,
		ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
	}, opts...)
	if err != nil {
		// var ok bool
		// var f *types.TransactionCanceledException
		// probs := make([]types.CancellationReason, 0)
		// if f, ok = errd.As[*types.TransactionCanceledException](err); ok && f != nil {
		// 	for _, c := range f.CancellationReasons {
		// 		if *c.Code == "None" {
		// 			log.Trace().Interface("reason", c).Msg("no reason - ignoring")
		// 			continue
		// 		} else if *c.Code == "ConditionalCheckFailed" {
		// 			log.Trace().Interface("reason", c).Msg("conditional check failed - ignoring")
		// 			continue
		// 		}
		// 		// panic("unexpected cancellation reason: " + *c.Code)
		// 		probs = append(probs, c)
		// 	}
		// }
		// if ok && len(probs) == 0 {

		// 	trace.SpanFromContext(ctx).SetStatus(codes.Ok, "")
		// 	if res == nil {
		// 		return &dynamodb.TransactWriteItemsOutput{}, nil
		// 	} else {
		// 		return res, nil
		// 	}
		// }
		log.Error().Err(err).Interface("transaction", input).Msg("transaction failed")
		return nil, err
	}

	// if res.ConsumedCapacity != nil {
	// 	log.Info().Interface("consumed", res.ConsumedCapacity).Msg("consumed capacity")
	// }

	return res, nil
}

type DynamoDBTransaction interface {
	Preconditions() []DynamoDBPrecondition
	AsTransactWriteItems() []types.TransactWriteItem
	SetTable(string)
	GetTable() string
	ResolveTables(TR TableResolver) bool
	ID() string
	SetID(string)
	ApplyConsumedCapacity(...types.ConsumedCapacity)
	ConsumedCapacity() *ConsumedCapacity
}

var _ Identifiable = (*DynamoDBUpdate)(nil)
var _ DynamoDBTransaction = (*DynamoDBUpdate)(nil)

type DynamoDBUpdate struct {
	*types.Update
	id            string
	preconditions []DynamoDBPrecondition
	cc            ConsumedCapacity
}

func (me *DynamoDBUpdate) ConsumedCapacity() *ConsumedCapacity {
	return &me.cc
}

func (me *DynamoDBUpdate) Preconditions() []DynamoDBPrecondition {
	return me.preconditions
}

func (me *DynamoDBUpdate) ApplyConsumedCapacity(cc ...types.ConsumedCapacity) {
	me.cc.ApplyArray(cc...)
}

var _ DynamoDBTransaction = (*DynamoDBGenericTransaction)(nil)

type DynamoDBGenericTransaction struct {
	operations    []types.TransactWriteItem
	id            string
	preconditions []DynamoDBPrecondition
	cc            ConsumedCapacity
}

func (me *DynamoDBGenericTransaction) ApplyConsumedCapacity(cc ...types.ConsumedCapacity) {
	me.cc.ApplyArray(cc...)
}

func (me *DynamoDBGenericTransaction) Id() string {
	return me.id
}

func (me *DynamoDBGenericTransaction) Preconditions() []DynamoDBPrecondition {
	return me.preconditions
}

func (me *DynamoDBGenericTransaction) AsTransactWriteItems() []types.TransactWriteItem {
	return me.operations
}

func (me *DynamoDBGenericTransaction) ConsumedCapacity() *ConsumedCapacity {
	return &me.cc
}

func (me *DynamoDBGenericTransaction) SetTable(in string) {
	panic("table should be set before trying to set on me")
}

func (me *DynamoDBGenericTransaction) GetTable() string {
	panic("get: table should be set before trying to get from me")
}

func (me *DynamoDBUpdate) ResolveTables(TR TableResolver) bool {
	if str, ok := TR.ResolveTable(*me.Update.TableName); ok {
		me.Update.TableName = &str
		return true
	}
	return false
}

func (me *DynamoDBGenericTransaction) ResolveTables(TR TableResolver) (ok bool) {
	for _, op := range me.operations {
		if op.Update != nil {
			if str, ok := TR.ResolveTable(*op.Update.TableName); ok {
				op.Update.TableName = &str
			} else {
				return false
			}
		}
		if op.Put != nil {
			if str, ok := TR.ResolveTable(*op.Put.TableName); ok {
				op.Put.TableName = &str
			} else {
				return false
			}
		}
		if op.Delete != nil {
			if str, ok := TR.ResolveTable(*op.Delete.TableName); ok {
				op.Delete.TableName = &str
			} else {
				return false
			}
		}
		if op.ConditionCheck != nil {
			if str, ok := TR.ResolveTable(*op.ConditionCheck.TableName); ok {
				op.ConditionCheck.TableName = &str
			} else {
				return false
			}
		}
	}
	return true
}

func (me *DynamoDBGenericTransaction) ID() string {
	return me.id
}

func (me *DynamoDBGenericTransaction) SetID(in string) {
	me.id = in
}

func CombineDynamoDBTransactions(transactions ...[]DynamoDBTransaction) []DynamoDBTransaction {
	transact := &DynamoDBGenericTransaction{
		operations:    make([]types.TransactWriteItem, 0),
		id:            "",
		preconditions: []DynamoDBPrecondition{},
	}
	for _, t := range transactions {
		for _, t2 := range t {
			items := t2.AsTransactWriteItems()
			transact.operations = append(transact.operations, items...)
			transact.preconditions = append(transact.preconditions, t2.Preconditions()...)
		}
	}

	if len(transact.operations) >= 25 {
		arr := make([]DynamoDBTransaction, 0)
		for i := 0; i < len(transact.operations); i += 25 {
			arr = append(arr, &DynamoDBGenericTransaction{
				operations:    transact.operations[i:int(math.Min(float64(i+25), float64(len(transact.operations))))],
				id:            transact.id,
				preconditions: transact.preconditions,
			})
		}
		return arr
	}

	return []DynamoDBTransaction{transact}
}

func (me *DynamoDBUpdate) ID() string {
	return me.id
}

func (me *DynamoDBUpdate) SetID(in string) {
	me.id = in
}

func (me *DynamoDBUpdate) GetTable() string {
	return *me.TableName
}

func (me *DynamoDBUpdate) SetTable(in string) {
	me.TableName = aws.String(in)
}

func (me *DynamoDBUpdate) AsUpdateItemInput() *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 me.TableName,
		Key:                       me.Key,
		ExpressionAttributeNames:  me.ExpressionAttributeNames,
		ExpressionAttributeValues: me.ExpressionAttributeValues,
		UpdateExpression:          me.UpdateExpression,
		ConditionExpression:       me.ConditionExpression,
		ReturnConsumedCapacity:    types.ReturnConsumedCapacityTotal,
		// ReturnValues:              ,
	}
}

func (me *DynamoDBUpdate) AsTransactWriteItem() *types.TransactWriteItem {
	return &types.TransactWriteItem{
		Update: me.Update,
	}
}

func (me *DynamoDBUpdate) AsTransactWriteItems() []types.TransactWriteItem {
	return []types.TransactWriteItem{*me.AsTransactWriteItem()}
}

func (me *DynamoDBUpdate) AsTransactWriteItemInput() *dynamodb.TransactWriteItemsInput {
	return &dynamodb.TransactWriteItemsInput{
		TransactItems:          me.AsTransactWriteItems(),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func ChunkedTransactWriteItems(items []types.TransactWriteItem, size int) [][]types.TransactWriteItem {
	return chunkBy(items, size)
}

// func (c *Client) FullScan(ctx context.Context, table *string) (*dynamodb.ScanOutput, error) {

// }

func FullScan(ctx context.Context, api DynamoDBAPI, input *dynamodb.ScanInput) (out *dynamodb.ScanOutput, err error) {

	out = &dynamodb.ScanOutput{
		Items: make([]map[string]types.AttributeValue, 0),
	}

	paginator := dynamodb.NewScanPaginator(api, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		out.Items = append(out.Items, page.Items...)
	}

	return out, err
}

func FullQuery(ctx context.Context, api DynamoDBAPI, input *dynamodb.QueryInput) (out *dynamodb.QueryOutput, err error) {

	out = &dynamodb.QueryOutput{
		Items: make([]map[string]types.AttributeValue, 0),
	}

	paginator := dynamodb.NewQueryPaginator(api, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		out.Items = append(out.Items, page.Items...)
	}

	return out, err
}

type DynamoDBAVMap[I attributevalue.Unmarshaler] struct {
	b map[string]bool
	m map[string]I
}
type DynamoDBAVMapIDBuilder[I attributevalue.Unmarshaler] func(val I) string

func NewDynamoDBAVMap[I attributevalue.Unmarshaler](out []map[string]types.AttributeValue, builder DynamoDBAVMapIDBuilder[I]) (*DynamoDBAVMap[I], error) {
	r := &DynamoDBAVMap[I]{
		b: make(map[string]bool),
		m: make(map[string]I),
	}

	for _, item := range out {
		i := new(I)
		if err := attributevalue.UnmarshalMap(item, i); err != nil {
			return nil, err
		}

		zid := builder(*i)

		r.b[zid] = true
		r.m[zid] = *i
	}

	return r, nil
}

func (me *DynamoDBAVMap[I]) Get(id string) I {
	return me.m[id]
}

func (me *DynamoDBAVMap[I]) Has(id string) bool {
	return me.b[id]
}
