package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/skodnik/go-contenttype/contenttype"
)

type ApiConfig struct {
	Platform       string
	SignSecret     string
	Queries        *database.Queries
	FileserverHits atomic.Int32
}

// getters
func (cfg *ApiConfig) hits() string {
	return fmt.Sprintf("Hits: %d\n", cfg.FileserverHits.Load())
}

// helpers
func logError(r, f string, err error) {
	msg := fmt.Sprintf(
		"[ERROR] %s: function %s failed, %s",
		r, f, err.Error())
	log.Println(msg)
}

func logInfo(r, s string) {
	msg := fmt.Sprintf("[INFO] %s: %s", r, s)
	log.Println(msg)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("Content-Type", contenttype.ApplicationJSON)
	w.WriteHeader(code)
	data := []byte(fmt.Sprintf("{\"error\": \"%s\"}", msg))
	w.Write(data)
}

func respondWithUnexpeted(
	w http.ResponseWriter, r,
	f string, err error,
) {
	logError(r, f, err)
	respondWithError(
		w, http.StatusInternalServerError,
		errorMesages[http.StatusInternalServerError],
	)
}

func respondWithInfoError(
	w http.ResponseWriter, pattern string, code int, extra ...string,
) {
	msg := errorMesages[code]
	if len(extra) > 0 {
		msg = msg + ": " + strings.Join(extra, " ")
	}
	logInfo(pattern, fmt.Sprintf("%d %s", code, msg))
	respondWithError(w, code, msg)
}

func respondWithJSON(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Add("Content-Type", contenttype.ApplicationJSON)
	w.WriteHeader(code)
	w.Write(payload)
}
