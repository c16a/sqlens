package connectors

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
)

type SqlServerConnectionOpts struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

type SqlServerConnector struct {
	conn *sql.Conn
}

func (p *SqlServerConnector) Connect(ctx context.Context) error {
	return p.conn.PingContext(ctx)
}

func NewMssqlConnector(ctx context.Context, opts *SqlServerConnectionOpts) (*SqlServerConnector, error) {
	url := fmt.Sprintf("sqlserver://%s:%s@%s:%d/database=%s", opts.User, opts.Password, opts.Host, opts.Port, opts.DatabaseName)
	db, err := sql.Open("sqlserver", url)
	if err != nil {
		return nil, err
	}
	conn, err := db.Conn(ctx)
	return &SqlServerConnector{conn: conn}, nil
}

func (p *SqlServerConnector) Query(ctx context.Context, query string) (*QueryResult, error) {
	rows, err := p.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
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

func (p *SqlServerConnector) Close(ctx context.Context) error {
	return p.conn.Close()
}

func (p *SqlServerConnector) Execute(ctx context.Context, query string) (int64, error) {
	tag, err := p.conn.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected()
}
