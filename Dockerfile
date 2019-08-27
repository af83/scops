# Usage :
# docker build -t scops .

FROM golang:1.12 AS builder

WORKDIR /go/src/scops
COPY . .

ENV GO111MODULE=on
RUN go get -d -v ./...
RUN go install -v ./...

FROM debian:latest

RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /go/bin/scops ./
RUN chmod +x ./scops
CMD ["./scops"]
