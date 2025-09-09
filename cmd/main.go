package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/magomzr/magomzr-api/handlers"
	"github.com/magomzr/magomzr-api/models"
)

var (
	dynamoClient *dynamodb.Client
)

const (
	postsPath      = "/posts"
	contentType    = "Content-Type"
	encodingError  = "error encoding response"
	appJson        = "application/json"
	invalidReqBody = "invalid request body"
)

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

	w.Header().Set(contentType, appJson)
	if err := json.NewEncoder(w).Encode(cards); err != nil {
		http.Error(w, encodingError, http.StatusInternalServerError)
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

	w.Header().Set(contentType, appJson)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, encodingError, http.StatusInternalServerError)
		return
	}
}

func GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := handlers.GetTags(r.Context(), dynamoClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting tags: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, appJson)
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		http.Error(w, encodingError, http.StatusInternalServerError)
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

	w.Header().Set(contentType, appJson)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, encodingError, http.StatusInternalServerError)
		return
	}
}

func getDrafts(w http.ResponseWriter, r *http.Request) {
	drafts, err := handlers.GetDrafts(r.Context(), dynamoClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting drafts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, appJson)
	if err := json.NewEncoder(w).Encode(drafts); err != nil {
		http.Error(w, encodingError, http.StatusInternalServerError)
		return
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, invalidReqBody, http.StatusBadRequest)
		return
	}

	ok, err := handlers.SavePost(r.Context(), dynamoClient, &post, true)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating post: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, appJson)
	if err := json.NewEncoder(w).Encode(ok); err != nil {
		http.Error(w, "error encoding response: %w", http.StatusInternalServerError)
	}
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "post ID is required in URL", http.StatusBadRequest)
		return
	}

	var post models.Post

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, invalidReqBody, http.StatusBadRequest)
		return
	}

	post.ID = id

	ok, err := handlers.SavePost(r.Context(), dynamoClient, &post, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("error updating post: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, appJson)
	if err := json.NewEncoder(w).Encode(map[string]any{"id": id, "ok": ok}); err != nil {
		http.Error(w, "error encoding response: %w", http.StatusInternalServerError)
	}
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	var reqBody models.ReqBody

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, invalidReqBody, http.StatusBadRequest)
		return
	}
	secretKey := reqBody.SecretKey

	token, err := handlers.GenerateKey(secretKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("error generating token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, appJson)
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, encodingError, http.StatusInternalServerError)
		return
	}
}

func buildRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get(postsPath, getAllPosts)
	r.Get(postsPath+"/{id}", getPostById)
	r.Get("/tags", GetTags)
	r.Get("/tags/{tag}", getPostsByTag)
	r.Post("/token", generateToken)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/drafts", getDrafts)
		r.Post(postsPath, createPost)
		r.Put(postsPath+"/{id}", updatePost)
	})
	return r
}

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	router := buildRouter()

	// Create request body from Lambda request
	var body io.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	}

	// Convert Lambda request → HTTP request WITH body
	httpReq, err := http.NewRequest(req.RequestContext.HTTP.Method, req.RawPath, body)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 500, Body: "Error creando request"}, nil
	}

	// Add context
	httpReq = httpReq.WithContext(ctx)

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
