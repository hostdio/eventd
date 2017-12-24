#!/bin/bash

set -e

createdb event_store;
createuser eventd;
psql event_store -c "GRANT ALL PRIVILEGES ON DATABASE event_store TO eventd;";
