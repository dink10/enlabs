FROM golang:1.13.0-alpine3.10 as builder

WORKDIR payment
ENV GO111MODULE=on CGO_ENABLED=1
RUN apk add --no-cache git openssh-client

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags '-w -s' -o ./bin/app ./cmd/api

FROM alpine:3.7

COPY --from=builder /go/payment/bin/app /app

RUN chmod +x /app

EXPOSE 3000

CMD ["/app"]
