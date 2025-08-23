package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"
)

func getPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Lista de posts desde GitHub Actions")
}

func createPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Post creado")
}

func getPostByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Fprintf(w, "Detalle del post %s", id)
}

func buildRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/posts", getPosts)
	r.Post("/posts", createPost)
	r.Get("/posts/{id}", getPostByID)
	return r
}

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	router := buildRouter()

	// Convert Lambda request → HTTP request
	httpReq, err := http.NewRequest(req.RequestContext.HTTP.Method, req.RawPath, nil)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 500, Body: "Error creando request"}, nil
	}

	// Copiar headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Usar recorder para capturar la respuesta
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httpReq)

	// Convertir HTTP response → Lambda response
	resp := events.LambdaFunctionURLResponse{
		StatusCode: rec.Code,
		Headers:    map[string]string{},
		Body:       rec.Body.String(),
	}

	for k, v := range rec.Header() {
		resp.Headers[k] = v[0]
	}

	return resp, nil
}

func main() {
	lambda.Start(handler)
}
