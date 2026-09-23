package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/popey17/urlshortener/internal/shortener"
	"github.com/popey17/urlshortener/internal/store"
)

type Handler struct {
	Store   *store.Store
	BaseUrl string
}

type request struct {
	URL string `json:"url"`
}

type responseStatus struct {
	Error string `json:"error"`
}

type createResponse struct {
	Code       string `json:"code"`
	ShortURL   string `json:"short_url"`
	URL        string `json:"url"`
	ClickCount int    `json:"click_count"`
}

func writeJson(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(responseStatus{Error: msg})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}

func (h *Handler) CreateUrl(w http.ResponseWriter, r *http.Request) {
	var req request
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// http.Error(w, "invalid JSON body", http.StatusBadRequest)
		writeJson(w, 400, err.Error())
		return
	}

	req.URL = strings.TrimSpace(req.URL)

	if err := shortener.ValidateURL(req.URL); err != nil {
		writeJson(w, 400, err.Error())
		return
	}

	code, err := shortener.GenerateCode()
	if err != nil {
		writeJson(w, 500, err.Error())
		return
	}

	addedUrl, err := h.Store.Create(code, req.URL)
	if err != nil {
		writeJson(w, 500, err.Error())
		return
	}

	short_url := h.BaseUrl + "/" + addedUrl.Code

	response := createResponse{
		Code:       addedUrl.Code,
		ShortURL:   short_url,
		URL:        addedUrl.Original,
		ClickCount: addedUrl.ClickCount,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	url, err := h.Store.GetByCode(code)
	if err != nil {
		writeJson(w, http.StatusNotFound, err.Error())
		return
	}
	err = h.Store.IncrementClicks(code)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, url.Original, http.StatusFound)
}
