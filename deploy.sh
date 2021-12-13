#!/bin/bash

# Export all the env variables from the .env file for all the services
set -a
. ./.env
set +a

# start all services
docker-compose up -d