FROM golang:latest AS builder
WORKDIR /app
ADD . .
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o svc ./cmd/svc/main.go

FROM centos:7
ENV container docker
RUN (cd /lib/systemd/system/sysinit.target.wants/; for i in *; do [ $i == \
systemd-tmpfiles-setup.service ] || rm -f $i; done); \
rm -f /lib/systemd/system/multi-user.target.wants/*;\
rm -f /etc/systemd/system/*.wants/*;\
rm -f /lib/systemd/system/local-fs.target.wants/*; \
rm -f /lib/systemd/system/sockets.target.wants/*udev*; \
rm -f /lib/systemd/system/sockets.target.wants/*initctl*; \
rm -f /lib/systemd/system/basic.target.wants/*;\
rm -f /lib/systemd/system/anaconda.target.wants/*;
COPY --from=builder /app/svc /usr/local/bin/
COPY --from=builder /app/deployments/svc-systemd/restart.service /etc/systemd/system/
RUN systemctl enable restart.service
STOPSIGNAL SIGRTMIN+3
ENV GRPC_ADDRESS=gw
ENV GRPC_PORT=50051
VOLUME [ "/sys/fs/cgroup" ]
CMD [ "/sbin/init" ]

