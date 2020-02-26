package todo

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	httptransport "github.com/go-kit/kit/transport/http"
)

// ErrMissingParam is thrown when an http request is missing a URL Parameter
var ErrMissingParam = errors.New("Missing parameter")

// MakeHTTPHandler creates http transport layer for the Todo service
func MakeHTTPHandler(endpoints TodoEndpoints) http.Handler {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(chiMiddleware.DefaultCompress)

	todoRouter := chi.NewRouter()

	todoRouter.Get("/", httptransport.NewServer(
		endpoints.GetAllForUserEndPoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Get("/{id}", httptransport.NewServer(
		endpoints.GetByIDEndpoint,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Post("/", httptransport.NewServer(
		endpoints.AddEndpoint,
		decodeAddRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Put("/{id}", httptransport.NewServer(
		endpoints.UpdateEndpoint,
		decodeUpdateRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Delete("/{id}", httptransport.NewServer(
		endpoints.DeleteEndpoint,
		decodeDeleteRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Mount("/todos", todoRouter)

	return r
}

func decodeGetRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	username := r.URL.Query().Get("username")
	log.Printf("username : %s", username)
	return GetAllForUserRequest{username}, err
}

func decodeGetByIDRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrMissingParam
	}
	return GetByIDRequest{id}, err
}

func decodeAddRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var todo Todo
	err = render.Decode(r, &todo)
	if err != nil {
		return nil, err
	}
	return AddRequest{todo}, err
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrMissingParam
	}
	var todo Todo
	err = render.Decode(r, &todo)
	if err != nil {
		return nil, err
	}
	return UpdateRequest{id, todo}, err
}

func decodeDeleteRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrMissingParam
	}
	return DeleteRequest{id}, err
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(error); ok && err != nil {
		encodeError(ctx, err, w)
		return json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInconsistentIDs, ErrMissingParam:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
