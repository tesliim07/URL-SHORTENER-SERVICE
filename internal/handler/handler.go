package handler

import (
	"url-shortener-service/internal/service"
	"github.com/rs/zerolog/log"

	"net/http"
	"encoding/json"
	"strings"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler {service : service}
}

// ShortenURL handles POST /shorten
// @Summary      Shorten a URL
// @Description  Shortens a given URL and returns the shortened version
// @Tags         urls
// @Param        url body string true "URL to shorten"
// @Produce      json
// @Success      201 {object} map[string]string "short_url"
// @Router       /shorten [post]
func (handler *Handler) ShortenURL (write http.ResponseWriter, request *http.Request){
	var payload struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil || strings.TrimSpace(payload.URL) == "" {
		log.Error().Err(err).Msg("invalid request body")
		http.Error(write, "invalid request body", http.StatusBadRequest)
		return
	}



	//call the service to shorten the URL
	shortURL, err := handler.service.ShortenURL(payload.URL)
	if err != nil {
		log.Error().Err(err).Msg("failed to shorten url")
		http.Error(write, "failed to shorten url", http.StatusInternalServerError)
		return
	}

	// return the shortened URL as JSON
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusCreated)
	json.NewEncoder(write).Encode(map[string]string{
		"short_url": shortURL,
	})
}

// Redirect handles GetOriginalURL/{code}
// @Summary      Redirect to original URL
// @Description  Redirects to the original URL based on the provided code
// @Tags         urls
// @Param        code path string true "Short code"
// @Produce      json
// @Success      302
// @Router       /{code} [get]
func (handler *Handler) Redirect(write http.ResponseWriter, request *http.Request) {
	code := strings.TrimPrefix(request.URL.Path, "/")
	if code == "" {
		log.Error().Msg("code is required")
		http.Error(write, "code is required", http.StatusBadRequest)
		return
	}
	originalURL, err := handler.service.GetOriginalURL(code)
	if err != nil {
		log.Error().Err(err).Msg("url not found")
		http.Error(write, "url not found", http.StatusNotFound)
		return
	}
	http.Redirect(write, request, originalURL, http.StatusFound)
}