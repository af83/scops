# Usage :
# docker build -t scops .

FROM golang:1.12 AS builder

WORKDIR /go/src/scops
COPY . .

ENV GO111MODULE=on
RUN go get -d -v ./...
RUN go install -v ./...

FROM debian:latest

WORKDIR /app
COPY --from=builder /go/bin/scops ./
RUN chmod +x ./scops
CMD ["./scops"]
