tag = latest

code:
	gofmt -e -s -w ${GOPATH}/src/github.com/sankalpjonn/mockingbird/..
	go vet -v github.com/sankalpjonn/mockingbird/...
	go get -v github.com/sankalpjonn/mockingbird
	go install github.com/sankalpjonn/mockingbird
	exit 0

image:
	docker build -t mockingbird:$(tag) .

run: image
	docker run -p 8000:8000 -d mockingbird:$(tag)
