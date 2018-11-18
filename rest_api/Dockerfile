FROM golang:1.11.2-alpine3.8

RUN apk update; apk upgrade

RUN apk add git

RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir $GOPATH/src/rest_api
WORKDIR $GOPATH/src/rest_api
ADD . $GOPATH/src/rest_api/
RUN dep ensure
RUN go install


CMD ["rest_api"]

