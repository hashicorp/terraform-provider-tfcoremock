#!/bin/bash

protoc --proto_path=. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative tfproto6.proto
