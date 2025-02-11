#!/bin/bash

## Build image
docker build -t forum .

## Run container
docker run -p 9000:9000 --name forum_container forum