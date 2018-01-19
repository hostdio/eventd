package postgres

import (
	"context"
	"database/sql"

	"github.com/hostdio/eventd/eventkit"

	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func New(connStr string) (*Client, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Client{
		db: db,
	}, nil
}

type Client struct {
	db *sql.DB
}

func (c Client) Close() error {
	return c.db.Close()
}

var (
	insertQuery = `
	INSERT INTO event_store(
		event_id,
		event_type,
		event_version,
		event_timestamp,
		event_payload,
		event_source,
		received_timestamp,
		stored_timestamp)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		NOW(),
		NOW()
	);`
)

func (c Client) Store(ctx context.Context, event eventkit.Event) error {
	stmt, prepErr := c.db.PrepareContext(ctx, insertQuery)
	if prepErr != nil {
		return errors.Wrap(prepErr, "postgres store: Could not prepare query")
	}
	defer stmt.Close()
	_, execErr := stmt.ExecContext(ctx,
		event.ID,
		event.Type,
		event.Version,
		event.Produced,
		event.Data.JSON(),
		event.Source)

	if execErr != nil {
		return execErr
	}
	return nil
}

var (
	scanQuery = `
	SELECT
		event_id,
		event_type,
		event_version,
		event_timestamp,
		event_payload,
		event_source,
		received_timestamp,
		stored_timestamp
	FROM
		event_store
	WHERE
		stored_timestamp >= $1
    LIMIT $2;
	`
)

type persistedEvent struct {
	ID                string
	Type              string
	Version           string
	Timestamp         time.Time
	Payload           string
	Source            string
	StoredTimestamp   time.Time
	ReceivedTimestamp time.Time
}

func (c Client) Scan(ctx context.Context, from time.Time, limit int) ([]eventkit.Event, error) {
	stmt, prepErr := c.db.PrepareContext(ctx, scanQuery)
	if prepErr != nil {
		panic(prepErr)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, from, limit)
	if err != nil {
		return nil, err
	}
	events := []eventkit.Event{}
	for rows.Next() {
		var pevent persistedEvent
		if err := rows.Scan(
			&pevent.ID,
			&pevent.Type,
			&pevent.Version,
			&pevent.Timestamp,
			&pevent.Payload,
			&pevent.Source,
			&pevent.StoredTimestamp,
			&pevent.ReceivedTimestamp,
		); err != nil {
			return nil, err
		}

		event := eventkit.Event{
			Namespace: "not set",
			Type:      pevent.Type,
			ID:        pevent.ID,
			Version:   pevent.Version,
			Source:    pevent.Source,
			Produced:  pevent.Timestamp,
			// Data:      pevent.Payload,
			// Metadata: map[string]interface{}{"exception": "not implemented"},
		}
		events = append(events, event)

	}

	return events, nil
}
