package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tsironi93/WebServer/docs"
	"github.com/tsironi93/WebServer/internal/database"
)

type apiConf struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	JWTSecret      string
}

func loadEnvAndConnect() apiConf {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("Secret must be set")
	}

	db, err := sql.Open("postgres", "postgres://postgres:@localhost:5432/chirpy?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to DB:", err)
	}
	log.Println("Successfully connected to DB!")

	dbQueries := database.New(db)
	return apiConf{
		db:        dbQueries,
		platform:  platform,
		JWTSecret: secret,
	}
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	cfg := loadEnvAndConnect()
	mux := http.NewServeMux()

	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetHits)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", cfg.HandlerCreateUser)
	mux.HandleFunc("POST /api/login", cfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.HandlerTokenRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.HandlerTokenRevoke)

	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.HandlerSingleChirp)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
