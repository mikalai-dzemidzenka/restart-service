FROM golang:latest AS builder
WORKDIR /app
ADD . .
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o main ./cmd/gw/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["/bin/sh","-c","./main"]
