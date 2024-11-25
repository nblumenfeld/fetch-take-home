FROM golang:1.23.3-alpine

RUN apk add build-base

WORKDIR /app

ADD . /app
RUN go build -o /fetch-take-home

EXPOSE 8080

CMD [ "/fetch-take-home" ]
