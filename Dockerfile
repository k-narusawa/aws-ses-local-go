FROM golang:1.22.3-alpine3.18 AS builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

ENV GO111MODULE=on
ENV CGO_ENABLED=1

WORKDIR /app

COPY . .

RUN make build

FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 8080

COPY --from=builder /app/engine /app/
COPY --from=builder /app/views /app/views

CMD /app/engine
