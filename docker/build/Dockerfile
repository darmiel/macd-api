FROM golang:1.16 AS builder
LABEL maintainer="darmiel <hi@d2a.io>"
WORKDIR /usr/src/app

# install golang dependencies
# https://github.com/montanaflynn/golang-docker-cache
COPY go.mod go.sum ./
RUN ls -la
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o macd .

FROM alpine AS git
RUN apk update && apk add --update git
RUN git clone https://github.com/vishnubob/wait-for-it
RUN mv wait-for-it/wait-for-it.sh /

FROM debian
COPY --from=builder /usr/src/app/macd .
COPY --from=git /wait-for-it.sh .
RUN chmod +x /wait-for-it.sh

# ca
RUN apt-get update
RUN apt-get install -y ca-certificates

ENTRYPOINT ["/macd"]