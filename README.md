# mockingbird

`mockingbird` makes mocking of API responses easy and provides a high throughput (>10000 RPS).

## Create Mock response

#### Request

```sh
$ curl -X POST \
  http://localhost:8000/egg \
  -d '
    {
  	"headers": {
  		"Content-Type": "application/json",
  		"mockheader1": "mockvalue1",
          "mockheader2": "mockvalue2"
  	},
  	"status_code": 200,
  	"body": "awesome!!!",
      "ttl": 120
  }'
```

### Response
```
{
  "egg_id": "5b70b541-7be3-4def-87c8-028cdc4e1141"
}
```

## Get mock response

#### Request
```sh
$ curl http://localhost:8000/egg/5b70b541-7be3-4def-87c8-028cdc4e1141 -v
```
#### Response
```
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8000 (#0)
> GET /egg/5b70b541-7be3-4def-87c8-028cdc4e1141 HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/7.47.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Mockheader1: mockvalue1
< Mockheader2: mockvalue2
< Date: Sat, 09 Jun 2018 12:15:23 GMT
< Content-Length: 10
<
* Connection #0 to host localhost left intact
awesome!!!
```


## Installation

__Minimum Go version:__ Go 1.9
__Minimum Redis version:__ Redis 3.0.7

Use [`go get`](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies) to install and update:

```sh
$ go get -u github.com/sankalpjonn/mockingbird/...
```

Redis needs to be installed and running locally. Refer to [`this`](https://redis.io/topics/quickstart) to install redis,
or run
```sh
$ sudo apt install redis-server
$ sudo service redis-server start
```

## Usage

From the command line, `mockingbird` can be run by providing the host. Default is 0.0.0.0:8000

```sh
$ mockingbird [options]
```

Available options:

```
  -host          run the server on this host (ip:port)
```
