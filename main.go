package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tsironi93/WebServer/internal/database"
)

type apiConf struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type parameters struct {
	Body string `json:"body"`
}

type returnVals struct {
	Cleaned_body string `json:"cleaned_body"`
}

func (cfg *apiConf) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConf) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	adminMetrics := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	w.Write([]byte(adminMetrics))
}

func (cfg *apiConf) handlerResetHits(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}

	cfg.db.DeleteAllUsers(r.Context())
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func respondJSON(resp *returnVals, w http.ResponseWriter, returnStatus int) {
	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		respondWithError(w, 500, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(returnStatus)
	w.Write(dat)
}

func hanlderJSONResponce(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "")
		return
	}

	const stars = "****"
	b := strings.ToLower(params.Body)
	badWordList := []string{"kerfuffle", "sharbert", "fornax"}

	for _, w := range badWordList {
		for {
			idx := strings.Index(b, w)
			if idx == -1 {
				break
			}

			params.Body = params.Body[:idx] + stars + params.Body[idx+len(w):]
			b = b[:idx] + stars + b[idx+len(w):]
		}
	}

	respondJSON(&returnVals{params.Body}, w, 200)
}

func (cfg *apiConf) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUser struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	create := createUser{}
	if err := decoder.Decode(&create); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if !strings.Contains(create.Email, "@") || !strings.Contains(create.Email, ".") {
		respondWithError(w, 401, "Wrong email format")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), create.Email)
	if err != nil {
		log.Println("Error creating user:", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	type userData struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	resp := &userData{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		respondWithError(w, 500, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println(err)
		return
	}

	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	cfg := &apiConf{
		db:       dbQueries,
		platform: platform,
	}
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetHits)
	mux.HandleFunc("POST /api/validate_chirp", hanlderJSONResponce)
	mux.HandleFunc("POST /api/users", cfg.HandlerCreateUser)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
