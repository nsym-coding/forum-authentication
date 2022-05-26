#!/bin/bash
docker image build -t team-forum .

 docker container run -p 8080:8080 --detach --name united team-forum

 docker image prune -a