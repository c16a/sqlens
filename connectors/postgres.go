package connectors

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type PgConnectionOpts struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

type PostgresConnector struct {
	conn *pgx.Conn
}

func (p *PostgresConnector) Connect(ctx context.Context) error {
	return p.conn.Ping(ctx)
}

func NewPostgresConnector(opts *PgConnectionOpts) (*PostgresConnector, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", opts.User, opts.Password, opts.Host, opts.Port, opts.DatabaseName)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &PostgresConnector{conn: conn}, nil
}

func (p *PostgresConnector) Query(ctx context.Context, query string) (*QueryResult, error) {
	rows, err := p.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column descriptions
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}

	var results []map[string]any

	// Iterate over rows and map to a slice of maps
	for rows.Next() {
		values := make([]any, len(columns))
		valuePointers := make([]any, len(columns))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]any)
		for i, colName := range columns {
			rowMap[colName] = values[i]
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &QueryResult{
		Columns: columns,
		Rows:    results,
	}, nil
}

func (p *PostgresConnector) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}

func (p *PostgresConnector) Execute(ctx context.Context, query string) (int64, error) {
	tag, err := p.conn.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
