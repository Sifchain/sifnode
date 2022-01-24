package client

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Sifchain/sifnode/tools/siflisten/events"
	_ "github.com/lib/pq"
)

type client struct {
	dbc *sql.DB
}

func (c client) GetCursor(ctx context.Context, name string) (int64, error) {
	row := c.dbc.QueryRowContext(ctx, "select position from cursors where name = $1", name)

	var position int64
	err := row.Scan(&position)
	if err != nil {
		return 0, err
	}

	return position, nil
}

func (c client) SetCursor(ctx context.Context, name string, position int64) error {
	_, err := c.dbc.ExecContext(ctx, "insert into cursors (name, position) values ($1, $2) on conflict (name) do update set position = $3", name, position, position)
	if err != nil {
		return err
	}

	return nil
}

func (c client) CreateEvent(ctx context.Context, ev *events.Event) error {
	attrs, err := json.Marshal(ev.Attributes)
	if err != nil {
		return err
	}

	_, err = c.dbc.ExecContext(ctx, "insert into events (type, height, metadata) values ($1, $2, $3)", ev.EventType, ev.Height, attrs)
	if err != nil {
		return err
	}

	return err
}

func (c client) GetEvent(ctx context.Context, cursorPosition int64) (*events.Event, error) {
	panic("implement me")
}

func (c client) GetEvents(ctx context.Context, cursorPosition int64) ([]*events.Event, error) {
	rows, err := c.dbc.QueryContext(ctx, "select id, type, height, metadata from events where id > $1", cursorPosition)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var evs []*events.Event
	for rows.Next() {
		var ev events.Event
		err := rows.Scan(&ev.ID, &ev.EventType, &ev.Height, &ev.Metadata)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(ev.Metadata, &ev.Attributes)
		if err != nil {
			return nil, err
		}

		evs = append(evs, &ev)
	}

	return evs, nil
}

func (c client) ConsumeForever(ctx context.Context, cursorName string, consume events.Consumer) {
	position, err := c.GetCursor(ctx, cursorName)
	if err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(time.Second):
			evs, err := c.GetEvents(ctx, position)

			for _, ev := range evs {
				err := consume(ctx, ev)
				if err != nil {
					break
				}

				position = ev.ID
			}

			err = c.SetCursor(ctx, cursorName, position)
			if err != nil {
				return
			}
		}
	}
}

func NewClient(dbc *sql.DB) events.Client {
	return &client{
		dbc: dbc,
	}
}
