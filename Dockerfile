FROM golang:1.11 as builder

## Create a directory and Add Code
RUN mkdir -p /go/src/github.com/orvice/monitor-server
WORKDIR /go/src/github.com/orvice/monitor-server
ADD .  /go/src/github.com/orvice/monitor-server

RUN go get
RUN CGO_ENABLED=0 go build


FROM orvice/go-runtime

COPY --from=builder /go/src/github.com/orvice/monitor-server/monitor-server .

ENTRYPOINT [ "./monitor-server" ]