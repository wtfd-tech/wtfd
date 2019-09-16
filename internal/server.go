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
	"sort"
	"strconv"

	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultPort = int64(80)
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key                 = []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	store               = sessions.NewFilesystemStore("", key) // generates filesystem store at os.tempdir
	errUserExisting     = errors.New("User with this name exists")
	errWrongPassword    = errors.New("Wrong Password")
	errUserNotExisting  = errors.New("User with this name does not exist")
	sshHost             = "localhost:2222"
	challs              = Challenges{}
	challcats           = ChallengeCategory{}
	mainpagetemplate    = template.New("")
	leaderboardtemplate = template.New("")
)

// Challenges Array of challenges but in nice with funcitons
type Challenges []Challenge

// ChallengeCategories Array of challengeCategories
type ChallengeCategories []ChallengeCategory

// JSONFile Challenge JSON File
type JSONFile struct {
	Categories []ChallengeCategoryJSON `json:"categories"`
	Challenges []ChallengeJSON         `json:"challenges"`
}

// ChallengeCategory as a go struct
type ChallengeCategory struct {
	Title      string `json:"title"`
	Challenges Challenges
}

// ChallengeCategoryJSON ChallengeCategory as JSON
type ChallengeCategoryJSON struct {
	Title      string   `json:"title"`
	Challenges []string `json:"challs"`
}

// Challenge is a challenge obv
type Challenge struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	Flag        string `json:"flag"`
	Points      int    `json:"points"`
	URI         string `json:"uri"`
	DepCount    int
	MinRow		int
	Row         int
	Solution    string `json:"solution"`
	DepIDs      []string
	Deps        []Challenge
	HasURI      bool // This emerges from URI != ""
}

// ChallengeJSON is Challenge as JSON
type ChallengeJSON struct {
	Name        string   `json:"name"`
	Description string   `json:"desc"`
	Solution    string   `json:"solution"`
	Flag        string   `json:"flag"`
	Points      int      `json:"points"`
	URI         string   `json:"uri"`
	Deps        []string `json:"deps"`
	HasURI      bool     // This emerges from URI != ""
}

// User was ist das wohl
type User struct {
	Name        string
	Hash        []byte
	DisplayName string
	Completed   []Challenge
}

type mainPageData struct {
	PageTitle              string
	Challenges             []Challenge
	SelectedChallengeID    string
	HasSelectedChallengeID bool
	User                   User
	IsUser                 bool
	Points                 int
}

// FillChallengeURI Fill host into each challenge's URI field and set HasURI
func (c Challenges) FillChallengeURI(host string) {
	for i := range c {
		if c[i].URI != "" {
			c[i].HasURI = true
			c[i].URI = fmt.Sprintf(c[i].URI, host)
		} else {
			c[i].HasURI = false
		}
	}
}

// Find finds a challenge from a string
func (c Challenges) Find(id string) (Challenge, error) {
	for _, v := range c {
		if v.Name == id {
			return v, nil
		}
	}
	return Challenge{}, fmt.Errorf("No challenge with this id")
}

// AllDepsCompleted checks if User u has completed all Dependent challenges of c
func (c Challenge) AllDepsCompleted(u User) bool {
	for _, ch := range c.Deps {
		a := false
		for _, uch := range u.Completed {
			if uch.Name == ch.Name {
				a = true
			}
		}
		if a == false {
			return false
		}
	}
	return true
}

// Contains looks if a username is in the datenbank
func Contains(username string) bool {
	_, err := ormLoadUser(username)
	return err == nil
}

// HasSolvedChallenge returns true if u has solved chall
func (u User) HasSolvedChallenge(chall Challenge) bool {
	for _, c := range u.Completed {
		if c.Name == chall.Name {
			return true
		}
	}
	return false
}

// CalculatePoints calculates Points of u
func (u User) CalculatePoints() int {
	points := 0

	for _, c := range u.Completed {
		points += c.Points
	}

	return points
}

// Get gets username based on username
func Get(username string) (User, error) {
	user, err := ormLoadUser(username)
	if err != nil {
		fmt.Printf("Get Error: username: %v, user: %v, err: %v\n", username, user, err)
		return User{}, err
	}
	return user, err

}

func resolveDeps(a []string) []Challenge {
	toReturn := []Challenge{}
	for _, b := range a {
		for _, c := range challs {
			if c.Name == b {
				toReturn = append(toReturn, c)
			}
		}
	}
	return toReturn

}

func countDeps(chall Challenge) int {
	max := 1
	if len(chall.Deps) == 0 {
		return 0

	}
	for _, a := range chall.Deps {
		depcount := countDeps(a)
		if depcount+1 > max {
			max = depcount + 1
		}
	}
	//return len(chall.DepIDs) + max
	return max

}

func countAllDeps() {
	for i := range challs {
		challs[i].DepCount = countDeps(challs[i])
	}
}
func reverseResolveAllDepIDs() {
	for i := range challs {
		for j := range challs {
			if i != j {
				for _, d := range challs[j].Deps {
					if d.Name == challs[i].Name {
						//						fmt.Printf("%s hat %s als revers dep\n", challs[i].Name, challs[j].Name)
						challs[i].DepIDs = append(challs[i].DepIDs, challs[j].Name)
						break
					}
				}
			}
		}
	}
}

func calculateMinRowNum(chall Challenge) int {
	if chall.MinRow != -1 {return chall.MinRow}
	chall.MinRow = 0
	for i := range chall.Deps {
		val := calculateMinRowNum(chall.Deps[i])
		if val >= chall.MinRow {
			chall.MinRow = val+1
		}
	}
	return chall.MinRow
}

func calculateRowNums() {
	maxcol := 0
	for i := range challs {
		if challs[i].DepCount > maxcol {maxcol = challs[i].DepCount}
		calculateMinRowNum(challs[i])
	}

	challmatrix := make([][]Challenge, maxcol+1)
	for i := range challmatrix {
		challmatrix[i] = []Challenge{}
	}

	for i := range challs {
		challmatrix[challs[i].DepCount] = append(challmatrix[challs[i].DepCount], challs[i])
	}
	for col := range challmatrix {
		sort.Slice(challmatrix[col], func(i, j int) bool {
			return challmatrix[col][i].MinRow < challmatrix[col][j].MinRow
		})
		row := 0
		for e := range challmatrix[col] {
			if challmatrix[col][e].MinRow > row {
				row = challmatrix[col][e].MinRow
			}
			challmatrix[col][e].Row = row
			row++
		}
	}
}

func resolveChalls(challcat []ChallengeJSON) {
	i := 0
	idsInChalls := []string{}
	for len(challcat) != 0 {
		//          fmt.Printf("challs: %v, challcat: %v\n",challs,challcat)
		this := challcat[i]
		if bContainsAllOfA(this.Deps, idsInChalls) {
			idsInChalls = append(idsInChalls, this.Name)
			challs = append(challs, Challenge{Name: this.Name, Description: this.Description, Flag: this.Flag, URI: this.URI, Points: this.Points, Deps: resolveDeps(this.Deps), Solution: this.Solution, MinRow: -1, Row: -1})
			challcat[i] = challcat[len(challcat)-1]
			challcat = challcat[:len(challcat)-1]
			i = 0
		} else {
			i++
		}

	}
	countAllDeps()
	reverseResolveAllDepIDs()
	calculateRowNums()
}

// Login checks if password is right for username and returns the User object of it
func Login(username, passwd string) (User, error) {
	user, err := Get(username)
	if err != nil {
		return User{}, err
	}
	if pwdRight := user.ComparePassword(passwd); !pwdRight {
		return User{}, errWrongPassword
	}
	fmt.Printf("User login: %s\n", username)
	return user, nil

}

// ComparePassword checks if the password is valid
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Hash, []byte(password)) == nil
}

// NewUser creates a new user object
func NewUser(name, password string) (User, error) {
	if Contains(name) {
		return User{}, errUserExisting
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}

	fmt.Printf("New User added: %s\n", name)
	return User{Name: name, Hash: hash}, nil

}

func mainpage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hasChall := vars["chall"] != ""
	session, _ := store.Get(r, "auth")
	val := session.Values["User"]
	user := &User{}
	newuser, ok := val.(*User)
	if ok {
		user = newuser
	}
	data := mainPageData{
		PageTitle:              "foss-ag O-Phasen CTF",
		Challenges:             challs,
		HasSelectedChallengeID: hasChall,
		SelectedChallengeID:    vars["chall"],
		User:                   *user,
		IsUser:                 ok,
	}
	err := mainpagetemplate.Execute(w, data)
	if err != nil {
		fmt.Printf("Template error: %v\n", err)

	}

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid Request")

	} else {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		session, _ := store.Get(r, "auth")
		if _, ok := session.Values["User"].(*User); ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Already logged in")
		} else {
			u, err := Login(r.Form.Get("username"), r.Form.Get("password"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Server Error: %v", err)
			} else {
				session.Values["User"] = u
				session.Save(r, w)
				fmt.Fprintf(w, "success")

			}

		}

	}

}

func submitFlag(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid Request")

	} else {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		session, _ := store.Get(r, "auth")

		user, ok := session.Values["User"].(*User)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Server Error: %v", "Not logged in")
			return
		}
		completedChallenge, err := challs.Find(r.Form.Get("challenge"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Server Error: %v", err)
			return
		}
		if r.Form.Get("flag") == completedChallenge.Flag {
			user.Completed = append(user.Completed, completedChallenge)
			if err = ormSolvedChallenge(*user, completedChallenge); err != nil {
				fmt.Errorf("ORM Error: %s", err)
			}
			fmt.Fprintf(w, "correct")

		} else {
			fmt.Fprintf(w, "not correct")
		}
		session.Save(r, w)

	}

}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid Request")

	} else {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		session, _ := store.Get(r, "auth")
		if _, ok := session.Values["User"].(*User); ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Already logged in")
		} else {

			if len(r.Form.Get("username")) < 5 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Username must be at least 5 characters")

			} else {
				u, err := NewUser(r.Form.Get("username"), r.Form.Get("password"))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Server Error: %v", err)
				} else {
					session.Values["User"] = u
					ormNewUser(u)
					session.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)

				}

			}
		}

	}

}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth")
	session.Values["User"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)

}

func solutionview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chall, err := challs.Find(vars["chall"])
	if err != nil {
		fmt.Fprintf(w, "ServerError: Challenge with is %s not found", vars["chall"])
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session, _ := store.Get(r, "auth")
	u, ok := session.Values["User"].(*User)
	if !ok {
		fmt.Fprintf(w, "ServerError: not logged in")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !u.HasSolvedChallenge(chall) {
		fmt.Fprintf(w, "did you just try to pull a sneaky on me?")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	md := markdown.ToHTML([]byte(chall.Solution), nil, nil)
	fmt.Fprintf(w, "%s", md)

}

func detailview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chall, err := challs.Find(vars["chall"])
	if err != nil {
		fmt.Fprintf(w, "ServerError: Challenge with is %s not found", vars["chall"])
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	md := markdown.ToHTML([]byte(chall.Description), nil, nil)
	fmt.Fprintf(w, "%s", md)

}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/static/favicon.ico")
}

// Server is the main server func, start it with
//  log.Fatal(wtfd.Server())
func Server() error {
	gob.Register(&User{})

	// Loading challs file
	challsFile, err := os.Open("challs.json")
	if err != nil {
		return err
	}
	defer challsFile.Close()
	var challsStructure JSONFile
	challsFileBytes, _ := ioutil.ReadAll(challsFile)
	if err := json.Unmarshal(challsFileBytes, &challsStructure); err != nil {
		return err
	}
	resolveChalls(challsStructure.Challenges)

	// Load database
	err = ormStart("./dblog")
	if err != nil {
		return err
	}

	// Fill in sshHost
	challs.FillChallengeURI(sshHost)

	// Parse Templates
	mainpagetemplate, err = template.ParseFiles("html/index.html", "html/footer.html", "html/header.html")
	if err != nil {
		return err
	}
	// Http sturf
	r := mux.NewRouter()
	r.HandleFunc("/", mainpage)
	r.HandleFunc("/favicon.ico", favicon)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/register", register)
	r.HandleFunc("/submitflag", submitFlag)
	r.HandleFunc("/{chall}", mainpage)
	r.HandleFunc("/detailview/{chall}", detailview)
	r.HandleFunc("/solutionview/{chall}", solutionview)
	// static
	r.PathPrefix("/static").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("html/static"))))

	Port := defaultPort
	if portenv := os.Getenv("WTFD_PORT"); portenv != "" {
		Port, _ = strconv.ParseInt(portenv, 10, 64)
	}
	fmt.Printf("WTFD Server Starting at port %d\n", Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", Port), r)
}

func bContainsA(a string, b []string) bool {
	for _, c := range b {
		if a == c {
			return true
		}

	}
	return false

}

func bContainsAllOfA(a, b []string) bool {
	for _, c := range a {
		if !bContainsA(c, b) {
			return false
		}
	}
	return true
}
