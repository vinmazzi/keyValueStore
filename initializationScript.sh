#!/bin/bash

curl -X PUT --data "This is my second Message." localhost:8080/v1/key/second
curl -X PUT --data "This is my third Message." localhost:8080/v1/key/third
curl -X PUT --data "This is my forth Message." localhost:8080/v1/key/forth
curl -X PUT --data "This is my fith Message." localhost:8080/v1/key/fith
curl -X PUT --data "This is my sixth Message." localhost:8080/v1/key/sixth
curl -X PUT --data "This is my seventh Message." localhost:8080/v1/key/seventh
curl -X PUT --data "This is my eighth Message." localhost:8080/v1/key/eighth
