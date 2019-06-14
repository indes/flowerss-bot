FROM golang:1.12-alpine as builder
#ENV CGO_ENABLED=0
COPY . /rss
RUN apk add git make gcc libc-dev && \
    cd /rss && make build
# Image starts here
FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /rss/flowerss-bot /bin/
VOLUME /root/.flowerss-bot
WORKDIR /
ENTRYPOINT ["/bin/flowerss-bot"]
