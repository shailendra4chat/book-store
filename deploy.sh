#!/bin/bash

# source the automatically created .env file
set -o allexport; source .env; set +o allexport

# start all services
docker-compose up -d