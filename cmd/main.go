package main

import (
    "net/http"

    "github.com/bagus123/go-kit-example3"
)

func main() {

    service := todo.NewInmemTodoService()

    endpoints := todo.MakeTodoEndpoints(service)

    err := http.ListenAndServe(":8000", todo.MakeHTTPHandler(endpoints))
    if err != nil {
        panic(err)
    }
}