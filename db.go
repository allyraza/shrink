package api

import (
	"database/sql"
	"os"

	"log"
)

// Repo represents a database wrapper
type Repo struct {
	Dsn string
	db  *sql.DB
}

// NewRepo creates a new repo
func NewRepo() *Repo {
	r := &Repo{
		Dsn: os.Getenv("SHORTNR_DATABASE"),
	}

	r.connect()

	return r
}

func (r *Repo) connect() {
	db, err := sql.Open("postgres", r.Dsn)
	if err != nil {
		log.Fatal("Could not connect to the database")
	}

	r.db = db
}

// Insert inserts a url to the database
// returns inserted row's ID if ok otherwise error
func (r *Repo) Insert(url string) (int, error) {
	var id int
	err := r.db.QueryRow(`INSERT INTO urls (url) VALUES ($1) RETURING id`).Scan(&id)
	if err != nil {
		log.Printf("Url (%v) could not inserted to the database.", url)
		return -1, err
	}

	return id, nil
}

// Update @todo
func (r *Repo) Update(id int, key string, url string) {

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
