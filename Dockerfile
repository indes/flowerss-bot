FROM golang:1.13-alpine as builder
#ENV CGO_ENABLED=0
COPY . /rssbot
RUN apk add git make gcc libc-dev && \
    cd /rssbot && make build

# Image starts here
FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /rssbot/rssbot /bin/
VOLUME /root/.rssbot
WORKDIR /root/.rssbot
ENTRYPOINT ["/bin/rssbot"]

