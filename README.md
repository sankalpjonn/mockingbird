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

For creating a mock, you can also use mbcli

```sh
go get -u github.com/sankalpjonn/mockingbird/cmd/mbcli/...

mbcli -create -H "mockheader1=mockvalue1" -H "mockheader2=mockvalue2" -status 200 -body="awesome!!!" -ttl=120
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

You can even run the docker image. If you do this, you do not have to install redis. To run the docker image run

```
$ make run
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

## Record and Playback

Mockingbird can create stub mappings from requests it has received by proxying the actual API end point. This allows you to “record” stub mappings from interaction with existing APIs.

To get started with the record feature, install the mocking bird client

```sh
go get -u github.com/sankalpjonn/mockingbird/cmd/mbcli/...
```

Then run
```sh
mbcli -record
```

You will be greeted with this and a shell will open
************************************
WELCOME TO THE MOCKING BIRD RECORDER

Please enter `help` for assistance
************************************
Now set the domain that you want to record by running

```sh
>>> domain http://api.jsonbin.io
```
Now start recording all calls to this domain by runninng
```sh
>>> start
started recording ...

Please use localhost:8080 to access http://api.jsonbin.io.

use the `stop` command to stop recording
```

Now make a request to the target API through Mockingbird proxy
```sh
curl -XGET 'http://localhost:8080/b/5b2be0376c6ba17a60dbddec'
```
No run stop
```sh
>>> stop
stopped recording
```
You should see that a file has been created. Something like `http_api.jsonbin.io_GET__b_5b2be0376c6ba17a60dbddec` under the folder `mockingbird_recordings`. This indicates that a stub has been recorded and requesting the url again will serve the recorded result. You can verify this by switching off wifi before making the next request.

```sh
curl -XGET 'http://localhost:8080/b/5b2be0376c6ba17a60dbddec'

{
    "msg": "this is my first mockingbird recording"
}

```
