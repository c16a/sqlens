package connectors

import "context"

type DatabaseConnector interface {
	Query(ctx context.Context, query string) (*QueryResult, error)
	Close(ctx context.Context) error
	Execute(ctx context.Context, query string) (int64, error)
	Connect(ctx context.Context) error
}

type QueryResult struct {
	Columns []string
	Rows    []map[string]any
}
