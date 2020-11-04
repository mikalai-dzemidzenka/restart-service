#!/bin/bash

function gen() {
	protoc \
	-I/usr/local/include \
	-I$GOPATH/src \
	-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I. \
	--grpc-gateway_out . \
       	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
       	api/proto/svc.proto
	protoc \
	-I/usr/local/include \
	-I$GOPATH/src \
	-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
       	-I. \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
       	--go-grpc_opt=paths=source_relative \
	api/proto/svc.proto
	mv api/proto/*.go api/proto/pb
}

if [[ $# -lt 1 ]] ; then
  exit 0
else
  MODE=$1
  shift
fi

if [ "$MODE" == "gen" ]; then
  echo ">>> Generating proto files..."
  echo
else
  exit 1
fi

if [ "${MODE}" == "gen" ]; then
  gen
else
  exit 1
fi

