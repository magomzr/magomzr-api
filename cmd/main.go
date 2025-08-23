package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/magomzr/magomzr-api/handlers"
)

var dynamoClient *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("error loading AWS config: %v", err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := handlers.GetAllPosts(r.Context(), dynamoClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("error obteniendo posts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func getPostsByTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	posts, err := handlers.GetPostsByTag(r.Context(), dynamoClient, tag)
	if err != nil {
		http.Error(w, fmt.Sprintf("error obteniendo posts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
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
	r.Get("/tags/{tag}", getPostsByTag)
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
