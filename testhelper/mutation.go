package testhelper

import (
	"testing"

	"cloud.google.com/go/spanner"
)

// InsertStruct generates the Spanner Mutation with error handling wrapped with
// *testing.T.
func InsertStruct(t *testing.T, tableName string, input interface{}) *spanner.Mutation {
	t.Helper()
	result, err := spanner.InsertStruct(tableName, input)
	if err != nil {
		t.Fatalf("insert failed, %+v", err)
	}
	return result
}
