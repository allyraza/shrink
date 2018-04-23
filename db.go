package api

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// Repo represents a database wrapper
type Repo struct {
	Dsn string
	db  *sql.DB
}

// NewRepo creates a new repo
func NewRepo(dsn string) *Repo {
	r := &Repo{
		Dsn: dsn,
	}

	r.connect()

	return r
}

func (r *Repo) connect() {
	db, err := sql.Open("postgres", r.Dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	r.db = db
}

// Insert inserts a url to the database
// returns inserted row's ID if ok otherwise error
func (r *Repo) Insert(url string) (int, error) {
	var id int
	err := r.db.QueryRow(`INSERT INTO urls (name, url) VALUES ($1, $2) RETURNING id`, "", url).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Update @todo
func (r *Repo) Update(id int, name string) bool {
	err := r.db.QueryRow(`UPDATE urls SET name = $1 WHERE id = $2`, name, id)
	return err == nil
}

// Find returns a url for given key
func (r *Repo) Find(key string) (string, error) {
	var url string
	err := r.db.QueryRow(`SELECT url FROM urls WHERE key = $1`, key).Scan(&url)
	if err != nil {
		log.Printf("Url not found for key (%s)", key)
		return "", err
	}

	return url, nil
}
