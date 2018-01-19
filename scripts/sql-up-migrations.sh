#!/bin/bash

set -e

goose -dir plugins/postgres/migrations postgres "user=eventd dbname=event_store sslmode=disable" up
