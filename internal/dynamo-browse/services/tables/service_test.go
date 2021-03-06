package tables_test

import (
	"context"
	"testing"

	"github.com/lmika/audax/internal/dynamo-browse/providers/dynamo"
	"github.com/lmika/audax/internal/dynamo-browse/services/tables"
	"github.com/lmika/audax/test/testdynamo"
	"github.com/stretchr/testify/assert"
)

func TestService_Describe(t *testing.T) {
	tableName := "service-test-data"

	client := testdynamo.SetupTestTable(t, testData)
	provider := dynamo.NewProvider(client)

	t.Run("return details of the table", func(t *testing.T) {
		ctx := context.Background()

		service := tables.NewService(provider)
		ti, err := service.Describe(ctx, tableName)
		assert.NoError(t, err)

		// Hash first, then range, then columns in alphabetic order
		assert.Equal(t, ti.Name, tableName)
		assert.Equal(t, "pk", ti.Keys.PartitionKey, "pk")
		assert.Equal(t, "sk", ti.Keys.SortKey, "sk")
		assert.Equal(t, []string{"pk", "sk"}, ti.DefinedAttributes)
	})
}

func TestService_Scan(t *testing.T) {
	tableName := "service-test-data"

	client := testdynamo.SetupTestTable(t, testData)
	provider := dynamo.NewProvider(client)

	t.Run("return all columns and fields in sorted order", func(t *testing.T) {
		ctx := context.Background()

		service := tables.NewService(provider)
		ti, err := service.Describe(ctx, tableName)
		assert.NoError(t, err)

		rs, err := service.Scan(ctx, ti)
		assert.NoError(t, err)

		// Hash first, then range, then columns in alphabetic order
		assert.Equal(t, rs.TableInfo, ti)
		assert.Equal(t, rs.Columns(), []string{"pk", "sk", "alpha", "beta", "gamma"})
	})
}

var testData = []testdynamo.TestData{
	{
		TableName: "service-test-data",
		Data: []map[string]interface{}{
			{
				"pk":    "abc",
				"sk":    "222",
				"alpha": "This is another some value",
				"beta":  1231,
			},
			{
				"pk":    "abc",
				"sk":    "111",
				"alpha": "This is some value",
			},
			{
				"pk":    "bbb",
				"sk":    "131",
				"beta":  2468,
				"gamma": "foobar",
			},
		},
	},
}
