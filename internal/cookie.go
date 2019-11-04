package wtfd

import (
	"net/http"
	"github.com/wtfd-tech/wtfd/internal/db"
)

func getLoginEmail(r *http.Request) (string, bool) {
	session, _ := store.Get(r, "auth")
	email, ok := session.Values["email"].(string)
	return email, ok && email != ""
}

func getUser(r *http.Request) (db.User, bool) {
	email, loggedIn := getLoginEmail(r)
	if !loggedIn {
		return db.User{}, false
	}
	loadeduser, err := wtfdDB.Get(email)
	return loadeduser, err == nil
}

func loginUser(r *http.Request, w http.ResponseWriter, email string) error {
	session, _ := store.Get(r, "auth")
	session.Values["email"] = email
	return session.Save(r, w)
}

func logoutUser(r *http.Request, w http.ResponseWriter) error {
	session, _ := store.Get(r, "auth")
	session.Values["email"] = ""
	return session.Save(r, w)
}
