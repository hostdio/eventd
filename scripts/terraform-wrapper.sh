#!/bin/bash

set -e

_cmd=$1
_project_id=$2

cd terraform

terraform $_cmd \
  -var google_cloud_project_id=$_project_id

cd - > /dev/null
