package issues

import (
	"context"
	"testing"

	"github.com/ClickHouse/clickhouse-go/v2"
	clickhouse_tests "github.com/ClickHouse/clickhouse-go/v2/tests"
	"github.com/stretchr/testify/require"
)

func Test1113(t *testing.T) {
	t.Skip("Object JSON type is deprecated. Test is kept for a historical reference.")

	var (
		conn, err = clickhouse_tests.GetConnection("issues", clickhouse.Settings{
			"max_execution_time":             60,
			"allow_experimental_object_type": true,
		}, nil, &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		})
	)
	ctx := context.Background()
	require.NoError(t, err)
	const ddl = "CREATE TABLE test_1113 (col_1 JSON, col_2 JSON) Engine MergeTree() ORDER BY tuple()"
	require.NoError(t, conn.Exec(ctx, ddl))
	defer func() {
		conn.Exec(ctx, "DROP TABLE IF EXISTS test_1113")
	}()

	batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO test_1113")
	require.NoError(t, err)

	v1 := map[string]struct {
		Str string
	}{"a": {Str: "value"}}
	v2 := map[string]any{}

	require.NoError(t, batch.Append(v1, v2))
	require.NoError(t, batch.Send())
}
