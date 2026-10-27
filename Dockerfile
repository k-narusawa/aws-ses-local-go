FROM public.ecr.aws/docker/library/golang:1.27.3-alpine AS builder

RUN apk update && apk upgrade && \
    apk --update add make bash

ENV CGO_ENABLED=0

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

CMD /app/engine
