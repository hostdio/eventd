#!/bin/bash

set -e

_project_id=$1

./scripts/gcloud-project.sh $_project_id

cd terraform && \
terraform apply \
  -var google_cloud_project_id=$_project_id \
  -var google_cloud_credentials=$HOME/.config/gcloud/application_default_credentials.json && \
cd -
