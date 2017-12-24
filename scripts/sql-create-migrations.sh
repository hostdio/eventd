#!/bin/bash

set -e

cd databases/mysql/migrations && goose create $1 sql && cd - > /dev/null
