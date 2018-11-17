FROM golang:1.11.2-alpine3.8

RUN apk update; apk upgrade

RUN apk add git

RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir $GOPATH/src/webapp
WORKDIR $GOPATH/src/webapp
ADD . $GOPATH/src/webapp/
RUN dep ensure
RUN go build

CMD ["go", "/webapp/webapp"]

