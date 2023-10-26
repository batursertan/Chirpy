package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiconfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("hits reset to: %d", cfg.fileServerHits)))
}
