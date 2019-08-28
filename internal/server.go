package wtfd

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key             = []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	store           = sessions.NewCookieStore(key)
	ErrUserExisting = errors.New("User with this name exists")
	ErrUserNotExisting = errors.New("User with this name does not exist")
)

type Users []User

type Challenge struct {
	Title       string
	Description string
	Flag        string
	Points      int
}

type User struct {
	Name      string
	Hash      []byte
	Completed []Challenge
}

type MainPageData struct {
	PageTitle  string
	Challenges []Challenge
	User       User
}

func (u *Users) Contains(username string) bool {
	for _, user := range *u {
		if user.Name == username {
			return true
		}
	}
	return false
}

func (u *Users) Get(username string) 


func (u *User) New(name, password string) (User, error) {
	if Users.Contains(name) {
		return User{}, ErrUserExisting
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}

	return User{Name: name, Hash: hash}, nil

}

func mainpage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth")
        user := Users.Get(session.Values["username"])

        data := MainPageData{
          PageTitle: "foss-ag O-Phasen CTF",
          Challenges: Challenges,
          User: User(user),
        }

}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth")


}

func logout(w http.ResponseWriter, r *http.Request) {

}

func Server() error {
	r := mux.NewRouter()
	r.HandleFunc("/", mainpage)

	return http.ListenAndServe(":80", r)
}
