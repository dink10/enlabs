FROM golang:1.13.0-alpine3.10 as builder

WORKDIR payment
ENV GO111MODULE=on CGO_ENABLED=0
RUN apk add --no-cache git openssh-client

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags '-w -s' -o ./bin/migrate ./tools/migrations

FROM alpine:3.7

COPY --from=builder /go/payment/bin/migrate /migrate

RUN chmod +x /migrate

CMD ["/migrate", "migrate"]