FROM golang:1.9 as builder

## Create a directory and Add Code
RUN mkdir -p /go/src/github.com/orvice/monitor-server
WORKDIR /go/src/github.com/orvice/monitor-server
ADD .  /go/src/github.com/orvice/monitor-server

# Download and install any required third party dependencies into the container.
RUN go-wrapper download
# RUN go-wrapper install
RUN CGO_ENABLED=0 go build


FROM alpine

COPY --from=builder /go/src/github.com/orvice/monitor-server/monitor-server .

RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
# Change TimeZone
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai
# Clean APK cache
RUN rm -rf /var/cache/apk/*

ENTRYPOINT [ "./monitor-server" ]