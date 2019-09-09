package wtfd

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key                = []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	store              = sessions.NewCookieStore(key)
	ErrUserExisting    = errors.New("User with this name exists")
	ErrWrongPassword   = errors.New("Wrong Password")
	ErrUserNotExisting = errors.New("User with this name does not exist")
	users              = Users{}
	challs             = Challenges{}
)

type Users []User
type Challenges []Challenge

type Challenge struct {
	Title       string `json:"title"`
	Description string `json:"desc"`
	Flag        string `json:"flag"`
	Points      int    `json:"points"`
}

type User struct {
	Name      string
	Hash      []byte
	Completed []*Challenge
}

type MainPageData struct {
	PageTitle  string
	Challenges []Challenge
	User       User
}

func (u *Users) Contains(username string) bool {
	_, err := u.Get(username)
	return err != nil
}

func (u *Users) Get(username string) (User, error) {
	for _, user := range *u {
		if user.Name == username {
			return user, nil
		}
	}
	return User{}, ErrUserNotExisting

}

func (u *Users) Login(username, passwd string) (User, error) {
	user, err := u.Get(username)
	if err != nil {
		return User{}, err
	}
	if pwdRight := user.ComparePassword(passwd); !pwdRight {
		return User{}, ErrWrongPassword
	}
	return user, nil

}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Hash, []byte(password)) == nil
}

func (u *User) New(name, password string) (User, error) {
	if users.Contains(name) {
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
	user, _ := session.Values["User"].(*User)

	data := MainPageData{
		PageTitle:  "foss-ag O-Phasen CTF",
		Challenges: challs,
		User:       *user,
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid Request")

	} else {
		session, _ := store.Get(r, "auth")
		if session.Values["User"] != "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Already logged in")
		} else {
			u, err := users.Login(r.Form["username"], r.Form["password"])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Server Error: %v", err)
			} else {
				session.Values["User"] = u
				session.Save(r, w)
				http.Redirect(w, r, mainpage)

			}

		}

	}

}

func logout(w http.ResponseWriter, r *http.Request) {

}

func Server() error {
	gob.Register(&User{})

        // Loading challs file
	challsFile, err := os.Open("challs.json")
	if err != nil {
		return err
	}
	defer challsFile.Close()
	challsFileBytes, _ := ioutil.ReadAll(challsFile)
	if err := json.Unmarshal(challsFileBytes, &challs); err != nil {
		return err
	}

        // Loading template files


	r := mux.NewRouter()
	r.HandleFunc("/", mainpage)

	return http.ListenAndServe(":80", r)
}
