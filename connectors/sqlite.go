package connectors

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type SqliteConnectionOpts struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

type SqliteConnector struct {
	conn *sql.Conn
}

func NewSqliteConnector(ctx context.Context, opts *SqliteConnectionOpts) (*SqliteConnector, error) {
	url := fmt.Sprintf("%s:%s@%s:%d/%s", opts.User, opts.Password, opts.Host, opts.Port, opts.DatabaseName)
	db, err := sql.Open("sqlite", url)
	if err != nil {
		return nil, err
	}
	conn, err := db.Conn(ctx)
	return &SqliteConnector{conn: conn}, nil
}

func (p *SqliteConnector) Connect(ctx context.Context) error {
	return p.conn.PingContext(ctx)
}

func (p *SqliteConnector) Query(ctx context.Context, query string) (*QueryResult, error) {
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

func (p *SqliteConnector) Close(ctx context.Context) error {
	return p.conn.Close()
}

func (p *SqliteConnector) Execute(ctx context.Context, query string) (int64, error) {
	tag, err := p.conn.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected()
}
