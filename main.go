package main

import (
	"log"
	"net/http"
	"os"

	"github.com/batursertan/Chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiconfig struct {
	fileServerHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	godotenv.Load(".env")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET enviroment variable is not set")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiconfig{
		fileServerHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/", fsHandler)
	router.Handle("/app", fsHandler)
	router.Handle("/assets/logo.png", http.FileServer(http.Dir("/assets")))

	apirouter := chi.NewRouter()
	apirouter.Get("/healthz", handlerReadiness)
	apirouter.Get("/reset", apiCfg.handlerReset)

	apirouter.Post("/polka/webhooks", apiCfg.handlerWebhook)

	apirouter.Post("/login", apiCfg.handlerLogin)

	apirouter.Post("/refresh", apiCfg.handlerRefresh)
	apirouter.Post("/revoke", apiCfg.handlerRevoke)

	apirouter.Post("/users", apiCfg.handlerUsersCreate)
	apirouter.Put("/users", apiCfg.handlerUsersUpdate)

	apirouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apirouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	apirouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	apirouter.Delete("/chirps/{chirpID}", apiCfg.handlerChirpsDelete)
	router.Mount("/api", apirouter)

	adminrouter := chi.NewRouter()
	adminrouter.Get("/metrics", apiCfg.handlerMetrics)
	router.Mount("/admin", adminrouter)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
