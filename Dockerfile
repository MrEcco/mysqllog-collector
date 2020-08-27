FROM golang:1.12.5-alpine3.9 as builder

# Prepare buildement
RUN apk add --update --no-cache git \
   && mkdir -p /output
RUN go get github.com/hpcloud/tail

# Build
COPY . /go/src/mysqllog-collector
WORKDIR /go/src/mysqllog-collector
RUN go build -i -o /output/mysqllog-collector .

FROM busybox:1.30.1-musl

COPY --from=builder /output/* /
ENTRYPOINT ["/mysqllog-collector"]
