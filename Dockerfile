FROM golang:1.9

RUN mkdir -p /go/src/github.com/sankalpjonn/mockingbird
ADD . /go/src/github.com/sankalpjonn/mockingbird
ADD run.sh run.sh
RUN go vet -v github.com/sankalpjonn/mockingbird/...
RUN go get -v github.com/sankalpjonn/mockingbird
RUN go install github.com/sankalpjonn/mockingbird

RUN apt-get update
RUN apt-get install -y redis-server

CMD ["sh", "run.sh"]
