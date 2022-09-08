FROM golang:1.17.8-alpine3.15 AS builder

WORKDIR /app/src

COPY go.* .
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ../bin/server ./services/mock_dealer/*.go

FROM alpine:3.15

COPY --from=builder /app/bin/server /usr/local/bin/server

CMD ["server"]
