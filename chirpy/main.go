package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync/atomic"

	"github.com/ckm54/go-projects/chirpy/internal/auth"
	"github.com/ckm54/go-projects/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
	platform       string
}

type ErrorRes struct {
	Error string `json:"error"`
}

var marshalError = "{\"error\": \"Error marshaling json\"}"

func main() {
	mux := *http.NewServeMux()
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	const assetsDir = "./assets"
	const port = "8080"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s", err)
	}

	dbQueries := database.New(db)
	apiCfg := apiConfig{
		database: dbQueries,
		platform: platform,
	}

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handleRegister)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handleGetChirp)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	rootFs := http.FileServer(http.Dir("."))

	rootWithMetrics := apiCfg.middlewareMetricsInc(rootFs)
	fs := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(assetsDir)))

	mux.Handle("/app", http.StripPrefix("/app", rootWithMetrics))
	mux.Handle("/app/", http.StripPrefix("/app", rootWithMetrics))

	mux.Handle("/assets/", http.StripPrefix("/assets", fs))

	server := &http.Server{
		Handler: &mux,
		Addr:    ":" + port,
	}

	server.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	hits := cfg.fileserverHits.Load()
	html := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)
	w.Write([]byte(html))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("Forbidden"))
		return
	}

	if err := cfg.database.DeleteUsers(context.Background()); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"message\": \"users deleted\"}"))
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	chirpInfo := database.CreateChirpParams{}
	if err := decoder.Decode(&chirpInfo); err != nil {
		w.Header().Set("Content-Type", "application/json")
		response := ErrorRes{
			Error: "Bad Request",
		}
		data, err := json.Marshal(&response)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(marshalError))
			return
		}
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	cleanedBody, err := validateChirp(chirpInfo.Body)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}

	chirp, err := cfg.database.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: chirpInfo.UserID,
	})

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(&chirp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.database.GetChirps(r.Context())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(&chirps)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("{\"error\": \"Bad Request\"}"))
	}

	chirp, err := cfg.database.GetChirp(r.Context(), id)
	if err != nil {
		// Check if the error is "no rows found"
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Chirp not found"}`))
			return
		}

		// Other DB error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}

	data, err := json.Marshal(&chirp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Comtent-Type", "application/json")
	w.Write(data)

}

func validateChirp(s string) (string, error) {
	maxLen := 140

	if len(s) > maxLen {
		return "", fmt.Errorf("chirp is too long. got %d characters max is %d characters", len(s), maxLen)
	}

	replacements := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}

	return replaceCaseInsensitive(s, replacements), nil

}

func (cfg *apiConfig) handleRegister(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	userInfo := parameters{}
	if err := decoder.Decode(&userInfo); err != nil {
		w.Header().Set("Content-Type", "application/json")
		response := ErrorRes{
			Error: "Bad Request",
		}
		data, err := json.Marshal(&response)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(marshalError))
			return
		}
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	hashedPass, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := cfg.database.CreateUser(context.Background(), database.CreateUserParams{
		Email:          userInfo.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	res := database.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	userData, err := json.Marshal(&res)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(userData)
}

func replaceCaseInsensitive(str string, replacements map[string]string) string {
	result := str
	for old, new := range replacements {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(old))
		result = re.ReplaceAllString(result, new)
	}
	return result
}
