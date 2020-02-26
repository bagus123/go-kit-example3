
# go-kit-example3

this sample use this library
- [go-kit](https://github.com/go-kit/kit "go-kit")
- [go-chi](https://github.com/go-chi/chi "go-chi")

## Installation

```shell
# 1. clone repository
git clone https://github.com/bagus123/go-kit-example3.git

# 2. downloads all dependencies and build binary
go build -o bin/app cmd/main.go

# run from binary
./bin/app 

# or run from source
go run cmd/main.go 

# note
# clean cache go
go clean --modcache

# remove unused module
go mod tidy
```

## Test

```shell

# create todo
curl --location --request POST 'http://localhost:8000/todos' \
--header 'Content-Type: application/json' \
--data-raw '{
	"id":"1",
	"username":"john doe",
	"text":"running",
	"completed": false
}'

# get all todo
curl --location --request GET 'http://localhost:8000/todos?username=tubagus' \
--header 'Content-Type: application/json'



```