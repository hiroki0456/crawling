package config

import (
	"context"
	"testing"

	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/spannertest"
	"cloud.google.com/go/spanner/spansql"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

// MockSpanner holds mock spanner configuration information
type MockSpanner struct {
	// Dbschemas hold DB Schema DDLs, the order matters
	Dbschemas []string
	// Mutations hold DB mutations, inculding inserts, updates, deletes
	Mutations []*spanner.Mutation
}

func SetUpTestDatabase(t *testing.T, cfg *MockSpanner) *spanner.Client {
	t.Helper()
	ctx := context.Background()
	srv, err := spannertest.NewServer("localhost:0")
	if err != nil {
		t.Fatalf("spannertest server setup failed, %+v", err)
	}
	srv.SetLogger(t.Logf)

	for _, schema := range cfg.Dbschemas {
		dbs, err := spansql.ParseDDL("testfile", schema)
		if err != nil {
			t.Fatalf("spannertest failed to parse DDL, %+v", err)
		}
		err = srv.UpdateDDL(dbs)
		if err != nil {
			t.Fatalf("spannertest failed to apply DDL, %+v", err)
		}
	}

	conn, err := grpc.DialContext(ctx, srv.Addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed dialing test spanner, %+v", err)
	}
	client, err := spanner.NewClient(ctx, "projects/test-project/instances/test-instance/databases/test-database", option.WithGRPCConn(conn))
	if err != nil {
		t.Fatalf("failed setting up client, %+v", err)
	}

	_, err = client.Apply(ctx, cfg.Mutations)
	if err != nil {
		t.Fatalf("failed to write to database, %+v", err)
	}
	t.Cleanup(func() {
		client.Close()
		srv.Close()
	})

	return client
}

func InsertStruct(t *testing.T, tableName string, input interface{}) *spanner.Mutation {
	t.Helper()
	result, err := spanner.InsertStruct(tableName, input)
	if err != nil {
		t.Fatalf("insert failed, %+v", err)
	}
	return result
}
