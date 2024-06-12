package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	m "github.com/mahl/gotext/models"
	re "github.com/mahl/gotext/resources"
)

func parseLang(language string) (string, bool) {
	strings := map[string]string{
		"es": "es",
		"en": "en",
	}

	string, ok := strings[language]
	if !ok {
		return "", false
	}

	return string, true
}

type Params struct {
	N    *uint64
	Lang string
}

func ParseParams(r *http.Request) (*Params, error) {
	langStr := r.URL.Query().Get("lang")
	lang, ok := parseLang(langStr)
	if !ok {
		return nil, errors.New("invalid lang")
	}

	nStr := r.URL.Query().Get("n")
	n, err := strconv.ParseUint(nStr, 10, 64)
	if err != nil {
		return &Params{N: nil, Lang: lang}, nil
	}

	return &Params{N: &n, Lang: lang}, nil
}

func GetTextHandler(w http.ResponseWriter, r *http.Request) {
	params, err := ParseParams(r)
	if err != nil {
		WriteError(w, err.Error())
		return
	}

	words := re.ReadWordFileN(params.N, params.Lang)
	response := &m.Text{Words: words}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
