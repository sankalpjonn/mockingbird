all:
	gofmt -e -s -w ${GOPATH}/src/github.com/sankalpjonn/mockingbird/..
	go vet -v github.com/sankalpjonn/mockingbird/...
	go get -v github.com/sankalpjonn/mockingbird
	go install github.com/sankalpjonn/mockingbird
	exit 0
