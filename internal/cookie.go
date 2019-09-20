package wtfd

import (
	"net/http"
)

func getLoginEmail(r *http.Request) (string, bool) {
	session, _ := store.Get(r, "auth")
	email, ok := session.Values["email"].(string)
	return email, ok && email != ""
}

func getUser(r *http.Request) (User, bool) {
	email, loggedIn := getLoginEmail(r)
	if !loggedIn {
		return User{}, false
	}
	loadeduser, err := Get(email)
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
