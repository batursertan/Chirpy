package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type apiconfig struct {
	fileServerHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiconfig{
		fileServerHits: 0,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/", fsHandler)
	router.Handle("/app", fsHandler)
	router.Handle("/assets/logo.png", http.FileServer(http.Dir("/assets")))

	apirouter := chi.NewRouter()
	apirouter.Get("/healthz", handlerReadiness)
	apirouter.Get("/reset", apiCfg.handlerReset)
	apirouter.Post("/validate_chirp", handleChirpsValidate)
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

func handleChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		w.WriteHeader(500)
		return
	}
	const maxChirpLenght = 140
	if len(params.Body) > maxChirpLenght {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleanedBody(params.Body, badWords)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})

}

func getCleanedBody(body string, badWords map[string]struct{}) string {

	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")

	return cleaned
}
