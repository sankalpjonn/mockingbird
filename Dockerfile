FROM golang:1.9

RUN mkdir -p /go/src/github.com/sankalpjonn/mockingbird

ADD bird /go/src/github.com/sankalpjonn/mockingbird/bird
ADD validator /go/src/github.com/sankalpjonn/mockingbird/validator
ADD main.go /go/src/github.com/sankalpjonn/mockingbird/main.go
ADD cmd /go/src/github.com/sankalpjonn/mockingbird/cmd
ADD run.sh run.sh

RUN go vet -v github.com/sankalpjonn/mockingbird/...
RUN go get -v github.com/sankalpjonn/mockingbird
RUN go install github.com/sankalpjonn/mockingbird
RUN go install github.com/sankalpjonn/mockingbird/cmd/mbird-cli
RUN apt-get update
RUN apt-get install -y redis-server

CMD ["sh", "run.sh"]
