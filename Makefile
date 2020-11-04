.PHONY: proto

proto:
	@protoc \
	-I/usr/local/include \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I. \
	--grpc-gateway_out api/proto/pb \
       	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
       	api/proto/svc.proto \
	&& \
	protoc \
	-I/usr/local/include \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
       	-I. \
	--go_out=api/proto/pb \
	--go_opt=paths=source_relative \
	--go-grpc_out=api/proto/pb \
       	--go-grpc_opt=paths=source_relative \
	api/proto/svc.proto
