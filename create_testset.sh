#!/bin/sh

curl -X POST localhost:8888/api/user/create -d "name=test&description=test&password=test"
curl -X POST localhost:8888/api/user/create -d "name=test1&description=test&password=test"

curl -X POST localhost:8888/api/task/create -d "title=task&description=test&key=bla&points=10&storyline=s1"
curl -X POST localhost:8888/api/task/create -d "title=task1&description=test&key=bla&points=10&storyline=s1"
curl -X POST localhost:8888/api/task/create -d "title=task2&description=test&key=bla&points=10&storyline=s2"
curl -X POST localhost:8888/api/task/create -d "title=task3&description=test&key=bla&points=10&storyline=s2"
