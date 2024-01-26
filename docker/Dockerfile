FROM golang:1.21.5-alpine3.19

RUN apk update && apk add --no-cache build-base

EXPOSE 8080

RUN mkdir /opt/helloworld
WORKDIR /opt/helloworld
RUN mkdir ./sqlite

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o helloworld cmd/main.go
RUN mv helloworld /usr/local/bin/

ENTRYPOINT ["helloworld"]