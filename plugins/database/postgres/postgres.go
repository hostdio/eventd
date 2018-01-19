package postgres

import (
	"context"
	"database/sql"

	"github.com/hostdio/eventd/api"

	"encoding/json"
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
		$7,
		NOW()
	);`
)

func (c Client) Store(ctx context.Context, event api.PublishedEvent) error {
	stmt, prepErr := c.db.PrepareContext(ctx, insertQuery)
	if prepErr != nil {
		return errors.Wrap(prepErr, "postgres store: Could not prepare query")
	}
	defer stmt.Close()
	_, execErr := stmt.ExecContext(ctx,
		event.ID,
		event.Type,
		event.Version,
		event.Timestamp,
		payloadToJSON(event.Payload),
		event.Source,
		event.ReceivedTimestamp)

	if execErr != nil {
		return execErr
	}
	return nil
}

func payloadToJSON(payload string) []byte {
	if payload == "" {
		return []byte("{}")
	}
	var p interface{}
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		panic(err)
	}
	switch p.(type) {
	case map[string]interface{}:
		return []byte(payload)
	default:
		payloadJSON := map[string]interface{}{
			"payload": payload,
		}
		byt, marshalErr := json.Marshal(payloadJSON)
		if marshalErr != nil {
			panic(marshalErr)
		}
		return byt
	}
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

func (c Client) Scan(ctx context.Context, from time.Time, limit int) ([]api.PersistedEvent, error) {
	stmt, prepErr := c.db.PrepareContext(ctx, scanQuery)
	if prepErr != nil {
		panic(prepErr)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, from, limit)
	if err != nil {
		return nil, err
	}
	events := []api.PersistedEvent{}
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
		event := api.PersistedEvent{
			PublishedEvent: &api.PublishedEvent{
				PublishEvent: &api.PublishEvent{
					BaseEvent: &api.BaseEvent{
						ID:        pevent.ID,
						Type:      pevent.Type,
						Version:   pevent.Version,
						Timestamp: pevent.Timestamp,
						Payload:   pevent.Payload,
						Source:    pevent.Source,
					},
				},
				ReceivedTimestamp: pevent.ReceivedTimestamp,
			},
			StoredTimestamp: pevent.StoredTimestamp,
		}

		events = append(events, event)

	}

	return events, nil
}
