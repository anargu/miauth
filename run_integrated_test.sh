#! /bin/bash

CONTAINER_NAME="miauth_postgres_db_test" 
CONTAINER_ID=$(docker inspect --format="{{.Id}}" ${CONTAINER_NAME} 2> /dev/null)	
if [[ "${CONTAINER_ID}" ]]; then
	echo "container $CONTAINER_NAME already running" 
	docker container start $CONTAINER_NAME 
else 
	echo "could not found container $CONTAINER_NAME..." 
	docker-compose -f test-compose.yml up --d 
fi
DEBUG=true go test
echo "stopping container..." 
docker container stop $CONTAINER_NAME
