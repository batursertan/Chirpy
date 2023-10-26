package main

import (
	"log"
	"net/http"

	"github.com/batursertan/Chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

type apiconfig struct {
	fileServerHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiconfig{
		fileServerHits: 0,
		DB:             db,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/", fsHandler)
	router.Handle("/app", fsHandler)
	router.Handle("/assets/logo.png", http.FileServer(http.Dir("/assets")))

	apirouter := chi.NewRouter()
	apirouter.Get("/healthz", handlerReadiness)
	apirouter.Get("/reset", apiCfg.handlerReset)
	apirouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apirouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)
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
