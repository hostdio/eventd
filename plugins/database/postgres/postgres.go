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
		namespace
		type
		id
		version
		source
		produced
		data
		metadata)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	);`
)

func (c Client) Store(ctx context.Context, event eventkit.Event) error {
	stmt, prepErr := c.db.PrepareContext(ctx, insertQuery)
	if prepErr != nil {
		return errors.Wrap(prepErr, "postgres store: Could not prepare query")
	}
	defer stmt.Close()
	_, execErr := stmt.ExecContext(ctx,
		event.Namespace,
		event.Type,
		event.ID,
		event.Version,
		event.Source,
		event.Produced,
		event.Data.JSON(),
		event.Metadata.JSON())

	if execErr != nil {
		return execErr
	}
	return nil
}

var (
	scanQuery = `
	SELECT
		namespace,
		type,
		id,
		version,
		source,
		produced,
		data,
		metadata
	FROM
		event_store
	WHERE
		produced >= $1
    LIMIT $2;
	`
)

type persistedEvent struct {
	Namespace string
	Type      string
	ID        string
	Version   string
	Source    string
	Produced  time.Time
	Data      []byte
	Metadata  []byte
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
			&pevent.Namespace,
			&pevent.Type,
			&pevent.ID,
			&pevent.Version,
			&pevent.Source,
			&pevent.Produced,
			&pevent.Data,
			&pevent.Metadata,
		); err != nil {
			return nil, err
		}

		event := eventkit.Event{
			Namespace: "not set",
			Type:      pevent.Type,
			ID:        pevent.ID,
			Version:   pevent.Version,
			Source:    pevent.Source,
			Produced:  pevent.Produced,
			// Data:      pevent.Payload,
			// Metadata: map[string]interface{}{"exception": "not implemented"},
		}
		events = append(events, event)

	}

	return events, nil
}
