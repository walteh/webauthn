package structure_test

import (
	"testing"

	"git.nugg.xyz/go-sdk/tfparse"
	"git.nugg.xyz/go-sdk/x"
	"git.nugg.xyz/webauthn/pkg/structure"
)

func TestIndexesAgainstTerraform(t *testing.T) {

	tables, err := tfparse.ParseDynamoDBTables("../../terraform/dynamo.tf", nil)
	if err != nil {
		t.Fatal(err)
	}

	p, err := x.NewTableResolverFromMap(tables, &structure.Tables{})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		indexable x.Indexable
		table     string
	}{
		{
			name:      "Credential",
			indexable: structure.NewCredentialQueryable(""),
			table:     p.Credential,
		},
		{
			name:      "Ceremony",
			indexable: structure.NewCeremonyQueryable(),
			table:     p.Ceremony,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			table := tables[tt.table]

			ok, reason := tt.indexable.PrimaryIndex().MatchesKeySchema(table.KeySchema, table.AttributeDefinitions)
			if !ok {
				t.Errorf("%s:primary index does not match key schema: %s", tt.table, reason)
			}

			// if len(table.GlobalSecondaryIndexes) != len(tt.indexable.SecondaryIndexes()) {
			// 	t.Errorf("%s: amount of secondary indexes wrong... got %d, want %d", tt.table, len(table.GlobalSecondaryIndexes), len(tt.indexable.SecondaryIndexes()))
			// }

			// sec := tt.indexable.SecondaryIndexes()

			// for _, index := range table.GlobalSecondaryIndexes {

			// 	ok, reason := sec[*index.IndexName].MatchesGSI(index, table.AttributeDefinitions)
			// 	if !ok {
			// 		t.Errorf("%s:secondary index does not match key schema: %s", tt.table, reason)
			// 	}

			// }

		})
	}

}
