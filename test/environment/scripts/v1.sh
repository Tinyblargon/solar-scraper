#!/bin/bash
set -e

influx v1 dbrp create \
  --bucket-id $(influx bucket list -n $DOCKER_INFLUXDB_INIT_BUCKET | awk 'NR>1 {print $1}') \
  --db $V1_DB_NAME \
  --rp $V1_RP_NAME \
  --default \
  --org $DOCKER_INFLUXDB_INIT_ORG

influx v1 auth create \
  --username ${V1_AUTH_USERNAME} \
  --password ${V1_AUTH_PASSWORD} \
  --read-bucket $(influx bucket list -n $DOCKER_INFLUXDB_INIT_BUCKET | awk 'NR>1 {print $1}') \
  --write-bucket $(influx bucket list -n $DOCKER_INFLUXDB_INIT_BUCKET | awk 'NR>1 {print $1}') \
  --org ${DOCKER_INFLUXDB_INIT_ORG}