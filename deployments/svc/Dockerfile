FROM golang:latest AS builder
WORKDIR /app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o restarter ./cmd/restarter/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o svc ./cmd/svc/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/restarter .
COPY --from=builder /app/svc .
CMD ["/bin/sh","-c","./restarter"]
