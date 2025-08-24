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
	"github.com/magomzr/magomzr-api/models"
	"github.com/magomzr/magomzr-api/pkg"
)

var dynamoClient *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("error loading AWS config: %v", err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
}

// To avoid returning all posts data, this method just maps
// the relevant fields using the Card struct.
func getAllPosts(w http.ResponseWriter, r *http.Request) {
	cards, err := handlers.GetAllPosts(r.Context(), dynamoClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting cards: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cards); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func getPostById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := handlers.GetPostById(r.Context(), dynamoClient, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting post: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := handlers.GetTags(r.Context(), dynamoClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting tags: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
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

func getDrafts(w http.ResponseWriter, r *http.Request) {
	drafts, err := handlers.GetDrafts(r.Context(), dynamoClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting drafts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(drafts); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Post creado")
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	var reqBody models.ReqBody

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	secretKey := reqBody.SecretKey

	token, err := pkg.GenerateKey(secretKey)
	if err != nil {
		http.Error(w, "key generation failed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func buildRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/posts", getAllPosts)
	r.Get("/posts/{id}", getPostById)
	r.Get("/tags", GetTags)
	r.Get("/tags/{tag}", getPostsByTag)
	r.Post("/token", generateToken)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/drafts", getDrafts)
		r.Post("/posts", createPost)
	})
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
