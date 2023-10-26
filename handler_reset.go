package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiconfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("hits reset to: %d", cfg.fileServerHits)))
}
