package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	URL_LENGTH = 4
)

var welcomeTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
	<title>Shortnr</title>
	<style>
	body {
		margin: 60px;
		padding: 0;
	}

	.welcome {
		margin: 0;
		padding: 0;
		font-family: "Helvetica Neue", sans-serif;
		font-weight: lighter;
		font-size: 45px;
		color: #aaa;
	}
	</style>
</head>
<body>
	<h1 class="welcome">Welcome</h1>
</body>
</html>`

var (
	repo *Repo
)

func init() {
	repo = NewRepo(os.Getenv("PG_URL"))
}

func GenerateRandomStr(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// WelcomeHandler handles / route
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(welcomeTemplate))
}

// ShortHandler shortens a given url (param url)
func ShortHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	id, err := repo.Insert(url)
	if err != nil {
		log.Printf("Failed to create url: %v", err)
		http.Error(w, "Failed to create url", http.StatusInternalServerError)
		return
	}

	key, _ := GenerateRandomStr(URL_LENGTH)
	queue := make(chan bool)
	go func() {
		repo.Update(id, key)
		queue <- true
	}()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"status": 200, "status_text": "ok", "short_url": %s}`, key)
	<-queue
}

// RedirectHandler redirects to url based on hash value
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("redirecting..."))
}

func verifyMethod(next http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		next(w, r)
	}
}

func verifyParam(next http.HandlerFunc, keys ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		for _, key := range keys {
			if params.Get(key) == "" {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
		}

		next(w, r)
	}
}

// New creates a mux and registers routes
func New() *http.ServeMux {
	r := http.NewServeMux()
	shortHandler := verifyMethod(verifyParam(ShortHandler, "url"), http.MethodPost)
	r.HandleFunc("/short", shortHandler)
	r.HandleFunc("/", verifyMethod(WelcomeHandler, http.MethodGet))
	r.HandleFunc("/r", verifyMethod(RedirectHandler, http.MethodGet))
	return r
}
