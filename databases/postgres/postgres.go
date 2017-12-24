package postgres

import (
	"database/sql"
	"context"
	"github.com/hostdio/eventd/api"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"encoding/json"
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
	if prepErr!= nil {
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
	if err := json.Unmarshal([]byte(payload),&p); err != nil {
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