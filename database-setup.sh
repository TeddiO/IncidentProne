#!/bin/bash

# Run this once on your local system to get it pre-populated with data.

docker exec -i incidentprone-postgres-1 psql -U postgres -c 'CREATE DATABASE incidentprone;'
docker exec -i incidentprone-postgres-1 psql -U postgres -d incidentprone < backup-plain-postgres.sql