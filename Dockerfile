FROM golang:1.7.1-alpine

RUN apk update && apk upgrade && apk add git

RUN mkdir -p /go/src/github.com/byuoitav
ADD . /go/src/github.com/byuoitav/pulse-eight-neo-microservice

WORKDIR /go/src/github.com/byuoitav/pulse-eight-neo-microservice
RUN go get -d -v
RUN go install -v

CMD ["/go/bin/pulse-eight-neo-microservice"]

EXPOSE 8011
