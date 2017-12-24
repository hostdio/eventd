#!/bin/bash

set -e

_sql_dialect=$1


cd databases/postgres/migrations && \
goose $_sql_dialect "user=eventd dbname=event_store sslmode=disable" up && \
cd - # should store current path and revert to that in the exit trigger
