package repository

import (
	"database/sql"
	"url-shortener-service/config"
	"fmt"
	// _ "github.com/jackc/pgx/v5/stdlib"  // ← replace lib/pq
	_ "github.com/lib/pq"
)

type UrlShortenerServiceRepository interface {
	SaveURL(code string, originalURL string) error
	GetOriginalURL(code string) (string, error)
}

type Repo struct {
	db *sql.DB
}

// NewRepo creates a new Repo and connects to PostgreSQL
func NewRepo(cfg *config.Config) (*Repo, error) {
	// connectionString := fmt.Sprintf(
	// 	"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	cfg.DBHost, 
	// 	cfg.DBPort, 
	// 	cfg.DBUser, 
	// 	cfg.DBPassword, 
	// 	cfg.DBName, 
	// )
	connectionString := fmt.Sprintf(
    "postgres://%s:%s@%s:%s/%s?sslmode=disable",
    cfg.DBUser,
    cfg.DBPassword,
    cfg.DBHost,
    cfg.DBPort,
    cfg.DBName,
)
	fmt.Println("Connection string:", connectionString)

	// db, err := sql.Open("pgx", connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &Repo{db : db}, err
}

func (repo *Repo) SaveURL(code string, originalURL string) error {
	_, err := repo.db.Exec(
		"INSERT INTO urls (code, original_url) VALUES ($1, $2)",
		code,
		originalURL,
	)
	if err != nil{
		return fmt.Errorf("failed to save url: %w", err)
	}
	return err
}

func (repo *Repo) GetOriginalURL(code string) (string, error){
	var originalURL string
	err := repo.db.QueryRow(
		"SELECT original_url FROM urls WHERE code = $1",
		code,
	).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("code not found: %s", code)
	}
	if err != nil{
		return "", fmt.Errorf("failed to get original url: %w", err)
	}
	return originalURL, err
}