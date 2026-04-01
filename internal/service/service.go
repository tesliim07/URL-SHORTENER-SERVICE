package service

import (
	"url-shortener-service/config"
	"url-shortener-service/internal/cache"
	"url-shortener-service/internal/repository"

	"fmt"
	"crypto/rand"
	"encoding/base64"
)

type Service struct {
	repo repository.UrlShortenerServiceRepository
	cache cache.URLCache
}

func NewService (repo repository.UrlShortenerServiceRepository, cache cache.URLCache, cfg *config.Config) *Service{
	return &Service{
		repo: repo,
		cache: cache,
	}
}

func (service *Service) ShortenURL(originalURL string) (string, error) {
	code, err := GenerateUniqueCode()
	if err != nil{
		return "", fmt.Errorf("failed to generate code: %w", err)
	}
	// save the code and original URL to the database
	err = service.repo.SaveURL(code, originalURL)
	if err != nil {
		return "", fmt.Errorf("(service) failed to save url: %w", err)
	}
	// cache the code and original URL in Redis
	err = service.cache.SetURL(code, originalURL)
	if err != nil {
		return "", fmt.Errorf("(service) failed to cache url: %w", err)
	}

	//return the full short URL e.g. http://localhost:8080/abc123
	shortURL := fmt.Sprintf("http://localhost:8080/%s", code)
	return shortURL, err
}

func (service *Service) GetOriginalURL(code string) (string, error) {
	originalURL, err := service.cache.GetURL(code)
	if err != nil {
		return "", fmt.Errorf("(service) cache error : %w", err)
	}
	if originalURL != ""{
		return originalURL, nil
	}
	originalURL, err = service.repo.GetOriginalURL(code)
	if err != nil {
		return "", fmt.Errorf("(service) failed to get url: %w", err)
	}
	err = service.cache.SetURL(code, originalURL)
	if err != nil {
		return "", fmt.Errorf("(service) failed to cache url: %w", err)
	}
	return originalURL, nil
}

func GenerateUniqueCode() (string, error){
	//creates a slice of 6 random bytes
	bytes := make([]byte, 6)

	//fill it with random data
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	//convert to a URL safe string and take first 8 characters
	return base64.URLEncoding.EncodeToString(bytes)[:8], nil
}