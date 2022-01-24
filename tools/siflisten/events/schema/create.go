package schema

import (
	"context"
	"database/sql"
)

func Create(ctx context.Context, dbc *sql.DB) error {
	_, err := dbc.ExecContext(ctx, `create temporary table events (
    	id serial primary key,
    	type varchar,
    	height int,
    	metadata bytea
	)`)
	if err != nil {
		return err
	}

	_, err = dbc.ExecContext(ctx, `create temporary table cursors (
    	name varchar primary key,
    	position bigint
	)`)

	return err
}
