package main

import (
	"fmt"
	"net/http"
	"regexp"
	"unicode/utf8"

	"github.com/skodnik/go-contenttype/contenttype"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("Content-Type", contenttype.ApplicationJSON)
	w.WriteHeader(code)
	data := []byte(fmt.Sprintf("{\"error\": \"%s\"}", msg))
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Add("Content-Type", contenttype.ApplicationJSON)
	w.WriteHeader(code)
	w.Write(payload)
}

func (v *ValidationRes) checkChirpLength() {
	if utf8.RuneCountInString(v.Body) > charLimit {
		v.Error = tooLong
		v.Status = http.StatusBadRequest
		return
	}
	v.Valid = true
	v.Error = ""
	v.Status = http.StatusOK
}

func (v *ValidationRes) cleanBody() {
	if v.Valid {
		body := v.Body
		for _, werd := range sanitationWords {
			re := regexp.MustCompile(`(?i)` + werd)
			body = re.ReplaceAllString(body, "****")
		}
		v.ClanedBody = body
	}
}
