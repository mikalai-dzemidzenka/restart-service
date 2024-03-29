FROM golang:latest AS builder
WORKDIR /app
ADD . .
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o svc ./cmd/svc/main.go

FROM alpine:latest
LABEL maintainer="danesh@riseup.net" \
    org.label-schema.name="alpine" \
    org.label-schema.vendor="Dockage" \
    org.label-schema.description="Docker image uses openRC as a process supervision on Alpine Linux" \
    org.label-schema.version="3.11-openrc" \
    org.label-schema.license="MIT"
RUN echo '@edge http://dl-cdn.alpinelinux.org/alpine/edge/main' >> /etc/apk/repositories \
    && echo '@edgecommunity http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories \
    && echo '@testing http://dl-cdn.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories \
    && apk add --no-cache openrc su-exec ca-certificates wait4ports \
    # Disable getty's
    && sed -i 's/^\(tty\d\:\:\)/#\1/g' /etc/inittab \
    && sed -i \
        # Change subsystem type to "docker"
        -e 's/#rc_sys=".*"/rc_sys="docker"/g' \
        # Allow all variables through
        -e 's/#rc_env_allow=".*"/rc_env_allow="\*"/g' \
        # Start crashed services
        -e 's/#rc_crashed_stop=.*/rc_crashed_stop=NO/g' \
        -e 's/#rc_crashed_start=.*/rc_crashed_start=YES/g' \
        # Define extra dependencies for services
        -e 's/#rc_provide=".*"/rc_provide="loopback net"/g' \
        /etc/rc.conf \
    # Remove unnecessary services
    && rm -f /etc/init.d/hwdrivers \
            /etc/init.d/hwclock \
            /etc/init.d/hwdrivers \
            /etc/init.d/modules \
            /etc/init.d/modules-load \
            /etc/init.d/modloop \
    # Can't do cgroups
    && sed -i 's/\tcgroup_add_service/\t#cgroup_add_service/g' /lib/rc/sh/openrc-run.sh \
    && sed -i 's/VSERVER/DOCKER/Ig' /lib/rc/sh/init.sh
COPY --from=builder /app/svc /bin/svc
COPY --from=builder /app/deployments/svc/restart /etc/init.d/
RUN chmod +x /etc/init.d/restart
RUN openrc -n
RUN rc-update add restart default

CMD ["/sbin/init"]
