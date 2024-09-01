#!/usr/bin/env sh

docker build -t social -f ./Dockerfile .

docker tag social dockerhub.com/voocel/social:latest
docker push dockerhub.com/voocel/social:latest