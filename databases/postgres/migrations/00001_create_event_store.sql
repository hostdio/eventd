-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE event_store (
  event_id            VARCHAR(2000) NOT NULL,
  event_type          VARCHAR(2000) NOT NULL,
  event_version       VARCHAR(2000) NOT NULL,
  event_timestamp     TIMESTAMP     NOT NULL,
  event_payload       JSON,
  event_source        VARCHAR(2000) NOT NULL,
  received_timestamp  TIMESTAMP     NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE event_store;
