package indexable

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/exp/slices"
)

type DynamoDBPrecondition func(ctx context.Context, api DynamoDBAPI, tr TableResolver, cc DynamoDBTransaction) (bool, error)

// Path: x/precondition.go
// Compare this snippet from dynamo/indexable.go:
// package dynamo

func MakePrecondition[I Indexable](idx I, check func(self map[string]types.AttributeValue) bool) DynamoDBPrecondition {
	return func(ctx context.Context, api DynamoDBAPI, tr TableResolver, cc DynamoDBTransaction) (bool, error) {
		tbl, ok := tr.ResolveTable(idx.ResolvableTableName())
		if !ok {
			return false, fmt.Errorf("table not found: %s", idx.ResolvableTableName())
		}
		getter, err := api.GetItem(ctx, &dynamodb.GetItemInput{
			Key:                    idx.PrimaryIndex().UpdateKey(),
			TableName:              aws.String(tbl),
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		})
		if getter != nil && getter.ConsumedCapacity != nil {
			cc.ApplyConsumedCapacity(*getter.ConsumedCapacity)
		}
		if err != nil {
			return false, err
		}

		return check(getter.Item), nil
	}

}

func MakeDoesNotExistPrecondition[I Indexable](idx I, field string) DynamoDBPrecondition {
	return MakePrecondition(idx, func(self map[string]types.AttributeValue) bool {
		return self == nil || self[field] == nil
	})
}

func MakeLastModifierPrecondition[I Indexable](idx I, field string, lm LastModifier) DynamoDBPrecondition {
	return MakePrecondition(idx, func(self map[string]types.AttributeValue) bool {
		if self == nil || self[field] == nil {
			return true
		} else {
			switch a := self[field].(type) {
			case *types.AttributeValueMemberN:
				zz, err := strconv.ParseFloat(a.Value, 64)
				if err != nil {
					return true
				}
				if zz < lm.At() {
					return true
				}
			case *types.AttributeValueMemberL:
				var out []float64
				err := attributevalue.UnmarshalList(a.Value, &out)
				if err != nil {
					return true
				}
				if !slices.Contains(out, lm.At()) {
					return true
				}
			}
		}

		return false
	})
}

func CheckPreconditions(ctx context.Context, api DynamoDBAPI, txn DynamoDBTransaction, tr TableResolver) (bool, error) {
	for _, cond := range txn.Preconditions() {
		res, err := cond(ctx, api, tr, txn)

		if err != nil {
			return true, err
		}
		if !res {
			return false, nil
		}
	}
	return true, nil
}
