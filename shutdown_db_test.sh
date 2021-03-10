#! /bin/bash

CONTAINER_NAME="miauth_postgres_db_test" 
CONTAINER_ID=$(docker inspect --format="{{.Id}}" ${CONTAINER_NAME} 2> /dev/null)	

docker container stop $CONTAINER_NAME