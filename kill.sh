#!/bin/bash
docker kill united

docker container prune 

docker system prune -a