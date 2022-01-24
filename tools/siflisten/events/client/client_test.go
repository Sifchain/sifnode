package client

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Sifchain/sifnode/tools/siflisten/events"
	"github.com/Sifchain/sifnode/tools/siflisten/events/schema"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func SetupTest(t *testing.T) (context.Context, *sql.DB, events.Client) {
	ctx := context.Background()
	db, err := sql.Open("postgres", "postgres://localhost?sslmode=disable")
	require.NoError(t, err)

	err = schema.Create(ctx, db)
	require.NoError(t, err)

	client := NewClient(db)

	return ctx, db, client
}

func TestCreateEvent(t *testing.T) {
	ctx, dbc, eventsClient := SetupTest(t)
	defer dbc.Close()

	attrs := []events.Attribute{{Key: "a", Value: "b"}}
	err := eventsClient.CreateEvent(ctx, &events.Event{
		EventType:  "test",
		Height:     1,
		Attributes: attrs,
	})
	require.NoError(t, err)

	evs, err := eventsClient.GetEvents(ctx, 0)
	require.NoError(t, err)
	require.Len(t, evs, 1)
	require.Equal(t, int64(1), evs[0].ID)
	require.Equal(t, "test", evs[0].EventType)
	require.Equal(t, int32(1), evs[0].Height)
	require.Equal(t, attrs, evs[0].Attributes)
}
