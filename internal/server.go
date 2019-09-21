package wtfd

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
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
	errUserExisting     = errors.New("user with this name exists")
	errWrongPassword    = errors.New("wrong Password")
	errUserNotExisting  = errors.New("user with this name does not exist")
	sshHost             = "localhost:2222"
	challs              = Challenges{}
	mainpagetemplate    = template.New("")
	leaderboardtemplate = template.New("")
	coolNames           = [...]string{
		"Anstruther's Dark Prophecy",
		"The Unicorn Invasion of Dundee",
		"Angus McFife",
		"Quest for the Hammer of Glory",
		"Magic Dragon",
		"Silent Tears of Frozen Princess",
		"Amulet of Justice",
		"Hail to Crail",
		"Beneath Cowdenbeath",
		"The Epic Rage of Furious Thunder",
		"Infernus Ad Astra",
		"Rise of the Chaos Wizards",
		"Legend of the Astral Hammer",
		"Goblin King of the Darkstorm Galaxy",
		"The Hollywood Hootsman",
		"Victorious Eagle Warfare",
		"Questlords of Inverness, Ride to the Galactic Fortress!",
		"Universe on Fire",
		"Heroes (of Dundee)",
		"Apocalypse 1992",
		"The Siege of Dunkeld (In Hoots We Trust)",
		"Masters of the Galaxy",
		"The Land of Unicorns",
		"Power of the Laser Dragon Fire",
		"Legendary Enchanted Jetpack",
		"Gloryhammer",
		"Hootsforce",
		"Battle for Eternity",
		"The Fires of Ancient Cosmic Destiny",
		"Dundaxian Overture",
		"The Battle of Cowdenbeath",
		"Return of the Astral Demigod of Unst",
		"The Knife of Evil",
		"Transmission",
	}
)

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
func (c Challenges) Find(id string) (*Challenge, error) {
	for _, v := range c {
		if v.Name == id {
			return v, nil
		}
	}
	return &Challenge{}, fmt.Errorf("no challenge with this id")
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
func Contains(username, displayname string) bool {
	count, _ := ormUserExists(User{Name: username, DisplayName: displayname})
	return count
}

// HasSolvedChallenge returns true if u has solved chall
func (u User) HasSolvedChallenge(chall *Challenge) bool {
	for _, c := range u.Completed {
		if c.Name == chall.Name {
			return true
		}
	}
	return false
}

// CalculatePoints calculates Points and updates user.Points
func (u *User) CalculatePoints() {
	points := 0

	for _, c := range u.Completed {
		points += c.Points
	}

	u.Points = points
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

func resolveDeps(a []string) []*Challenge {
	var toReturn []*Challenge
	for _, b := range a {
		for _, c := range challs {
			if c.Name == b {
				toReturn = append(toReturn, c)
			}
		}
	}
	return toReturn

}

func countDeps(chall *Challenge) int {
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

func calculateMinRowNum(chall *Challenge) int {
	if chall.MinRow != -1 {
		return chall.MinRow
	}
	chall.MinRow = len(chall.Deps)-1
        if chall.MinRow < 0 {
          chall.MinRow = 0
        }
	for _, d := range chall.Deps {
		val := calculateMinRowNum(d)
		if val > chall.MinRow {
			chall.MinRow = val + 1
		}
	}
	return chall.MinRow
}

func calculateRowNums() {
	maxcol := 0
	for i := range challs {
		if challs[i].DepCount > maxcol {
			maxcol = challs[i].DepCount
		}
		calculateMinRowNum(challs[i])
	}

	challmatrix := make([][]*Challenge, maxcol+1)
	for i := range challmatrix {
		challmatrix[i] = []*Challenge{}
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
	var idsInChalls []string
	for len(challcat) != 0 {
		//          fmt.Printf("challs: %v, challcat: %v\n",challs,challcat)
		this := challcat[i]
		if bContainsAllOfA(this.Deps, idsInChalls) {
			idsInChalls = append(idsInChalls, this.Name)
			challs = append(challs, &Challenge{Name: this.Name, Description: this.Description, Flag: this.Flag, URI: this.URI, Points: this.Points, Deps: resolveDeps(this.Deps), Solution: this.Solution, MinRow: -1, Row: -1})
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
func Login(username, passwd string) error {
	user, err := Get(username)
	if err != nil {
		return err
	}
	if pwdRight := user.ComparePassword(passwd); !pwdRight {
		return errWrongPassword
	}
	fmt.Printf("User login: %s\n", username)
	return nil

}

// ComparePassword checks if the password is valid
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Hash, []byte(password)) == nil
}

// NewUser creates a new user object
func NewUser(name, password, displayname string) (User, error) {
	if Contains(name, displayname) {
		return User{}, errUserExisting
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}

	fmt.Printf("New User added: %s\n", name)
	return User{Name: name, Hash: hash, DisplayName: displayname}, nil

}

func generateUserName() (string, error) {

	var name string
	for _, s := range coolNames {
		if exists, err := ormDisplayNameExists(s); !exists {
			if err != nil {
				return "", err

			}
			name = s
			break
		}
	}
	for name == "" {
		name = strconv.FormatInt(rand.Int63(), 10)
		if exists, err := ormDisplayNameExists(name); !exists {
			if err != nil {
				return "", err

			}
			name = ""

		}

	}
	return name, nil

}

func leaderboardpage(w http.ResponseWriter, r *http.Request) {
	userobj, ok := getUser(r)
	user := &userobj
	genu := ""
	var err error
	if !ok {
		genu, err = generateUserName()
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error: %v", err)
		}
	}
	allUsers, err := ormAllUsersSortedByPoints()
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error: %v", err)
	}
	data := leaderboardPageData{
		PageTitle:     "foss-ag O-Phasen CTF",
		GeneratedName: genu,
		Leaderboard:   true,
		AllUsers:      allUsers,
		User:          user,
		IsUser:        ok,
	}
	err = leaderboardtemplate.Execute(w, data)
	if err != nil {
		fmt.Printf("Template error: %v\n", err)

	}

}
func mainpage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hasChall := vars["chall"] != ""
	userobj, ok := getUser(r)
	user := &userobj
	genu := ""
	var err error
	if !ok {
		genu, err = generateUserName()
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error: %v", err)
		}

	}
	data := mainPageData{
		PageTitle:              "foss-ag O-Phasen CTF",
		Challenges:             challs,
		GeneratedName:          genu,
		HasSelectedChallengeID: hasChall,
		SelectedChallengeID:    vars["chall"],
		User:                   user,
		IsUser:                 ok,
	}
	err = mainpagetemplate.Execute(w, data)
	if err != nil {
		fmt.Printf("Template error: %v\n", err)

	}

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Invalid Request")

	} else {
		if err := r.ParseForm(); err != nil {
			_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		if _, ok := getLoginEmail(r); ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Already logged in")
		} else {
			email := r.Form.Get("username")
			err := Login(email, r.Form.Get("password"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(w, "Server Error: %v", err)
			} else if err := loginUser(r, w, email); err != nil {
				_, _ = fmt.Fprintf(w, "success")
			}

		}

	}

}

func submitFlag(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Invalid Request")

	} else {
		if err := r.ParseForm(); err != nil {
			_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		user, ok := getUser(r)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "Server Error: %v", "Not logged in")
			return
		}
		completedChallenge, err := challs.Find(r.Form.Get("challenge"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "Server Error: %v", err)
			return
		}
		if r.Form.Get("flag") == completedChallenge.Flag {
			user.Completed = append(user.Completed, completedChallenge)

			if err = ormSolvedChallenge(user, completedChallenge); err != nil {
				_ = fmt.Errorf("ORM Error: %s", err)
			}

			user.CalculatePoints()

			if err = ormUpdateUser(user); err != nil {
				_ = fmt.Errorf("ORM Error: %s", err)
			}

			_, _ = fmt.Fprintf(w, "correct")

		} else {
			_, _ = fmt.Fprintf(w, "not correct")
		}
		if err != nil {
			log.Print(err)
		}
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Invalid Request")

	} else {
		if err := r.ParseForm(); err != nil {
			_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		if _, ok := getLoginEmail(r); ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Already logged in")
		} else {

			if len(r.Form.Get("username")) < 5 {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(w, "Username must be at least 5 characters")

			} else {
				u, err := NewUser(r.Form.Get("username"), r.Form.Get("password"), r.Form.Get("displayname"))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprintf(w, "Server Error: %v", err)
				} else {
					_ = ormNewUser(u)
					http.Redirect(w, r, "/", http.StatusFound)

				}

			}
		}

	}

}

func logout(w http.ResponseWriter, r *http.Request) {
	_ = logoutUser(r, w)
	http.Redirect(w, r, "/", http.StatusFound)

}

func solutionview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chall, err := challs.Find(vars["chall"])
	if err != nil {
		_, _ = fmt.Fprintf(w, "ServerError: Challenge with is %s not found", vars["chall"])
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, ok := getUser(r)
	if !ok {
		_, _ = fmt.Fprintf(w, "ServerError: not logged in")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !u.HasSolvedChallenge(chall) {
		_, _ = fmt.Fprintf(w, "did you just try to pull a sneaky on me?")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	md := markdown.ToHTML([]byte(chall.Solution), nil, nil)
	_, _ = fmt.Fprintf(w, "%s", md)

}

func detailview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chall, err := challs.Find(vars["chall"])
	if err != nil {
		_, _ = fmt.Fprintf(w, "ServerError: Challenge with is %s not found", vars["chall"])
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	md := markdown.ToHTML([]byte(chall.Description), nil, nil)
	_, _ = fmt.Fprintf(w, "%s", md)

}

func uriview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chall, err := challs.Find(vars["chall"])
	if err != nil {
		_, _ = fmt.Fprintf(w, "ServerError: Challenge with is %s not found", vars["chall"])
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = fmt.Fprint(w, chall.URI)
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
	leaderboardtemplate, err = template.ParseFiles("html/leaderboard.html", "html/footer.html", "html/header.html")
	if err != nil {
		return err
	}
	// Http sturf
	r := mux.NewRouter()
	r.HandleFunc("/", mainpage)
	r.HandleFunc("/leaderboard", leaderboardpage)
	r.HandleFunc("/favicon.ico", favicon)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/register", register)
	r.HandleFunc("/submitflag", submitFlag)
	r.HandleFunc("/{chall}", mainpage)
	r.HandleFunc("/detailview/{chall}", detailview)
	r.HandleFunc("/solutionview/{chall}", solutionview)
	r.HandleFunc("/uriview/{chall}", uriview)
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
