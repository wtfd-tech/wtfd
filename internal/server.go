package wtfd

import (
	"encoding/gob"
	"encoding/json"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	defaultPort = int64(8080)
)

var (
	config Config
	store  sessions.Store

	errUserExisting     = errors.New("user with this name exists")
	errWrongPassword    = errors.New("wrong Password")
	errUserNotExisting  = errors.New("user with this name does not exist")
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
	maxcol = 0
	maxrow = 0
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

func calculateRowNums() {
	cols := make(map[int][]*Challenge)

	for _, chall := range challs {
		col := chall.DepCount
		cols[col] = append(cols[col], chall)
		if col > maxcol {
			maxcol = col
		}
	}

	fmt.Println("col\t[         <name>]\tmin\trow")
	for i := 0; i <= maxcol; i++ {
		if _, ok := cols[i]; !ok {
			continue
		} //Skip empty columns

		for _, chall := range cols[i] {
			chall.MinRow = 0
			for _, dep := range chall.Deps {
				if dep.Row > chall.MinRow {
					chall.MinRow = dep.Row
				}
			}
		}

		sort.Slice(cols[i], func(x, y int) bool {
			if cols[i][x].MinRow == cols[i][y].MinRow {
				return stringCompareLess(cols[i][x].Name, cols[i][y].Name)
			} else {
				return cols[i][x].MinRow < cols[i][y].MinRow
			}
		})

		row := 0
		for j := 0; j < len(cols[i]); j++ {
			if row < cols[i][j].MinRow {
				row = cols[i][j].MinRow
			}
			cols[i][j].Row = row
			if row > maxrow {
				maxrow = row
			}
			row++
			fmt.Printf("%1d\t[%15s]\t%3d %3d\n", i, cols[i][j].Name, cols[i][j].MinRow, cols[i][j].Row)
		}
	}
}

// https://stackoverflow.com/a/35099450
func stringCompareLess(si, sj string) bool {
	var siLower = strings.ToLower(si)
	var sjLower = strings.ToLower(sj)
	if siLower == sjLower {
		return si < sj
	}
	return siLower < sjLower
}

func resolveChalls(jsons []*ChallengeJSON) {
	i := 0
	var idsInChalls []string
	for len(jsons) != 0 {
		//          fmt.Printf("challs: %v, jsons: %v\n",challs,jsons)
		this := jsons[i]
		if bContainsAllOfA(this.Deps, idsInChalls) {
			idsInChalls = append(idsInChalls, this.Name)
			challs = append(challs, &Challenge{Name: this.Name, Description: this.Description, Flag: this.Flag, URI: this.URI, Points: this.Points, Deps: resolveDeps(this.Deps), Solution: this.Solution, MinRow: -1, Row: -1})
			jsons[i] = jsons[len(jsons)-1]
			jsons = jsons[:len(jsons)-1]
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
		RowNums:       make([]gridinfo, 0),
		ColNums:       make([]gridinfo, 0),
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
	rnums := make([]gridinfo, maxrow+1)
	for i := 0; i <= maxrow; i++ {
		rnums[i] = gridinfo{
			Index: i,
			Pos:   i + 1,
		}
	}
	cnums := make([]gridinfo, maxcol+1)
	for i := 0; i <= maxcol; i++ {
		cnums[i] = gridinfo{
			Index: i,
			Pos:   i + 1,
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
		RowNums:                rnums,
		ColNums:                cnums,
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
                                        login(w,r)

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

	_, _ = fmt.Fprintf(w, "%s", chall.Solution)

}

func detailview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chall, err := challs.Find(vars["chall"])
	if err != nil {
		_, _ = fmt.Fprintf(w, "ServerError: Challenge with is %s not found", vars["chall"])
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = fmt.Fprintf(w, "%s", chall.Description)

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

	var key []byte

	//Test if config file exists
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		// Generate a new key
		key = securecookie.GenerateRandomKey(32)

		//Write default config to disk
		config = Config{
			Key:              base64.StdEncoding.EncodeToString(key),
			Port:             defaultPort,
			ChallengeInfoDir: "../challenges/info/",
			SSHHost:          "ctf.foss-ag.de",
		}
		configBytes, _ := json.MarshalIndent(config, "", "\t")
		_ = ioutil.WriteFile("config.json", configBytes, os.FileMode(0600))
	} else {
		//Load config file
		var (
			configBytes []byte
			err         error
		)

		if configBytes, err = ioutil.ReadFile("config.json"); err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(configBytes, &config); err != nil {
			return err
		}

		// Decode key
		key, err = base64.StdEncoding.DecodeString(config.Key)
		if err != nil {
			log.Fatal("Could not decode config.json:Key")
		}
	}

	store = sessions.NewFilesystemStore("", key) // generates filesystem store at os.tempdir

	//Load challs from dirs
	var challsStructure []*ChallengeJSON

	files, err := ioutil.ReadDir(config.ChallengeInfoDir)
	if err != nil {
		return err
	}

	for _, current := range files {
		var (
			challDir     = config.ChallengeInfoDir + "/" + current.Name()
			jsonName     = challDir + "/meta.json"
			readmeName   = challDir + "/README.md"
			solutionName = challDir + "/SOLUTION.md"

			jsonBytes     []byte
			readmeBytes   []byte
			solutionBytes []byte

			jsonStruct ChallengeJSON

			err error
		)

		// Check if meta.json exists and load it into a ChallengeJSON struct
		if !current.IsDir() {
			continue
		}
		if jsonBytes, err = ioutil.ReadFile(jsonName); err != nil {
			log.Println(err)
			continue
		}
		if json.Unmarshal(jsonBytes, &jsonStruct) != nil {
			log.Println(err)
			continue
		}

		// Set name from folder name
		jsonStruct.Name = current.Name()

		// Load and compile markdown files
		if readmeBytes, err = ioutil.ReadFile(readmeName); err == nil {
			jsonStruct.Description = string(markdown.ToHTML(readmeBytes, nil, nil))
		} else {
			jsonStruct.Description = "<i>Description unavailable</i>"
		}

		if solutionBytes, err = ioutil.ReadFile(solutionName); err == nil {
			jsonStruct.Solution = string(markdown.ToHTML(solutionBytes, nil, nil))
		} else {
			jsonStruct.Description = "<i>Solution unavailable</i>"
		}

		challsStructure = append(challsStructure, &jsonStruct)
	}

	fixDeps(challsStructure)
	resolveChalls(challsStructure)

	// Load database
	err = ormStart("./dblog")
	if err != nil {
		return err
	}

	// Fill in SSHHost
	challs.FillChallengeURI(config.SSHHost)
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

	Port := config.Port
	if portenv := os.Getenv("WTFD_PORT"); portenv != "" {
		Port, _ = strconv.ParseInt(portenv, 10, 64)
	}
	fmt.Printf("WTFD Server Starting at port %d\n", Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", Port), r)
}

func fixDeps(jsons []*ChallengeJSON) {
	challsByName := make(map[string]*ChallengeJSON)
	for _, chall := range jsons {
		challsByName[chall.Name] = chall
	}
	for _, chall := range jsons {
		keepDep := make(map[string]bool)

		//Inititalize maps
		for _, dep := range chall.Deps {
			keepDep[dep] = true
		}

		//Kick out redundant challenges
		for _, dep := range chall.Deps {
			for _, depdep := range challsByName[dep].Deps {
				if _, ok := keepDep[depdep]; ok {
					keepDep[depdep] = false
				}
			}
		}

		//Rebould dependency array
		var newdeps []string
		for name, keep := range keepDep {
			if keep {
				newdeps = append(newdeps, name)
			}
		}

		//Write to struct
		chall.Deps = newdeps
	}
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
