package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	cfg.FileserverHits.Store(0)
	cfg.DbPtr.ResetUsers(req.Context())
	cfg.DbPtr.ResetChirps(req.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0; users and chirps databases returned to initial state."))
}