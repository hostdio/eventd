#!/bin/bash

../bin/cloud_sql_proxy \
  -instances=hostd-eventd:europe-west1:master-instance-2=tcp:3306
