-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- This is a total reset of the system. You will lose all your data
DROP TABLE event_store;

CREATE TABLE event_store (
  namespace   VARCHAR(2064),
  type        VARCHAR(2064),
  id          VARCHAR(2064),
  version     VARCHAR(2064),
  source      VARCHAR(2064),
  produced    TIMESTAMP,
  data        JSON,
  metadata    JSON,

  PRIMARY KEY(id)
);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.


DROP TABLE event_store;

CREATE TABLE event_store (
  event_id            VARCHAR(2000) NOT NULL,
  event_type          VARCHAR(2000) NOT NULL,
  event_version       VARCHAR(2000) NOT NULL,
  event_timestamp     TIMESTAMP     NOT NULL,
  event_payload       JSON,
  event_source        VARCHAR(2000) NOT NULL,
  received_timestamp  TIMESTAMP     NOT NULL,
  stored_timestamp    TIMESTAMP     NOT NULL
);
