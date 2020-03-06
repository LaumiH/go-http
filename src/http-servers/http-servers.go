package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
)

var usernames []string

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("super-secret-key")
	//key =
	//store = sessions.NewCookieStore(key)
	store = sessions.NewFilesystemStore("", key)
)

// http.Handler -> Funktion nimmt http.ResponseWriter und einen http.Request als Argument
//dynamische Requests

//der http.Request erlaubt
//GET via r.URL.Query().Get("keyword")
//POST via Parametern im request body r.Body

func hello(w http.ResponseWriter, req *http.Request) {
	//addCookie(w, "username", "laumi")

	session, _ := store.Get(req, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// fill in the http response
	fmt.Fprintf(w, "hello\n")
	//io.WriteString(w, "hello!")
}

func hello_cookie(w http.ResponseWriter, req *http.Request) {
	cookie := getCookie(w, req, "username")
	if cookie == nil || Contains(usernames, cookie.Value) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	fmt.Fprintf(w, "hello\n")
}

// liest http request header und schreibt sie in die Antwort
func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: /v\n", name, h)
		}
	}
}

func addCookie(w http.ResponseWriter, name string, value string) {
	expiration := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
}

func getCookie(w http.ResponseWriter, req *http.Request, name string) *http.Cookie {
	cookie, _ := req.Cookie(name)
	return cookie
	//fmt.Fprint(w, cookie)
}

//session management

func login(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "cookie-name")

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(req, w)
}

func login_cookie(w http.ResponseWriter, req *http.Request) {
	addCookie(w, "username", "laumi")
	usernames = append(usernames, "laumi")
}

func logout_cookie(w http.ResponseWriter, req *http.Request) {
	addCookie(w, "username", "")
	usernames = usernames[:0]
}

func logout(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(req, w)
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func main() {

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3001"
	}

	//setzt den default router -> /hello und nimmt die funktion von oben
	http.HandleFunc("/", hello)
	http.HandleFunc("/headers", headers)

	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	//port und handler
	//nil sagt dass es den default router von eben nutzen soll
	http.ListenAndServe(":"+PORT, nil)
}

//curl -s http://localhost:8080
//curl -s -I http://localhost:8080/login
//curl -s --cookie "cookie-name=MTQ4NzE5Mz..." http://localhost:8080/secret

//http POST
//curl --data "cart=bananas" localhost:3001
