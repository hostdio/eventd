#!/bin/bash

curl -XPOST http://localhost:8080 -d \
    '{"type":"hello", "id": "id", "version":"1", "timestamp": "2017-12-24T20:32:00Z", "source": "self"}'