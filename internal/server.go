package wtfd

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
	admintemplate       = template.New("")
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

type adminPageData struct {
	PageTitle     string
	User          *User
	IsUser        bool
	Points        int
	Leaderboard   bool
	AllUsers      []_ORMUser
	GeneratedName string
	Style         template.HTMLAttr
	RowNums       []gridinfo
	ColNums       []gridinfo
}
type leaderboardPageData struct {
	PageTitle     string
	User          *User
	IsUser        bool
	Points        int
	Leaderboard   bool
	AllUsers      []_ORMUser
	GeneratedName string
	Style         template.HTMLAttr
	RowNums       []gridinfo
	ColNums       []gridinfo
}
type mainPageData struct {
	PageTitle              string
	Challenges             []*Challenge
	Leaderboard            bool
	SelectedChallengeID    string
	HasSelectedChallengeID bool
	GeneratedName          string
	User                   *User
	IsUser                 bool
	Points                 int
	RowNums                []gridinfo
	ColNums                []gridinfo
}

func getUserData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userobj, _ := getUser(r)
	user := &userobj
	if user.Admin == false {
		fmt.Fprintf(w, "Nice Try, %s", user.DisplayName)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	u, err := ormLoadUser(vars["user"])
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error: %v", err)
	}
	userToReturn := User{Name: u.Name, DisplayName: u.DisplayName, Points: u.Points, Admin: u.Admin}
	jsonToReturn, err := json.Marshal(&userToReturn)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error: %v", err)
		return
	}
	w.Write(jsonToReturn)
}
func adminpage(w http.ResponseWriter, r *http.Request) {
	userobj, ok := getUser(r)
	user := &userobj
	if user.Admin == false {
		fmt.Fprintf(w, "Nice Try, %s", user.DisplayName)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		fmt.Printf("a: %#v", r.FormValue("admin"))
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error: %v", err)
			return
		}
		dumb, err := strconv.Atoi(r.FormValue("points"))
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error: %v", err)
			return
		}
		isAdmin := r.FormValue("admin") == "on"
		u := User{Name: r.FormValue("name"), DisplayName: r.FormValue("displayname"), Points: dumb, Admin: isAdmin}
		//                fmt.Printf("a: %#v",u)

		err = ormUpdateUser(u)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error: %v", err)
			return
		}

		r.Method = "GET"
		adminpage(w, r)
		return
	}
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
	data := adminPageData{
		PageTitle:     "foss-ag O-Phasen CTF",
		GeneratedName: genu,
		Leaderboard:   false,
		AllUsers:      allUsers,
		User:          user,
		IsUser:        ok,
		RowNums:       make([]gridinfo, 0),
		ColNums:       make([]gridinfo, 0),
	}
	err = admintemplate.Execute(w, data)
	if err != nil {
		fmt.Printf("Template error: %v\n", err)

	}

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
			_ = updateScoreboard()

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
					login(w, r)
					_ = updateScoreboard()

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
	_, _ = fmt.Fprintf(w, "%s<br><p>Solves: %d</p>", chall.Description, ormGetSolveCount(*chall))

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

func authorview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chall, err := challs.Find(vars["chall"])
	if err != nil {
		_, _ = fmt.Fprintf(w, "ServerError: Challenge with is %s not found", vars["chall"])
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = fmt.Fprint(w, chall.Author)
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
	// Packr

	box := packr.New("Box", "./html")
	maintemplatetext, err := box.FindString("html/index.html")
	if err != nil {
		return err
	}
	headertemplatetext, err := box.FindString("html/header.html")
	if err != nil {
		return err
	}
	footertemplatetext, err := box.FindString("html/footer.html")
	if err != nil {
		return err
	}
	admintemplatetext, err := box.FindString("html/admin.html")
	if err != nil {
		return err
	}
	leaderboardtemplatetext, err := box.FindString("html/leaderboard.html")
	if err != nil {
		return err
	}

	// Parse Templates
	admintemplate, err = template.New("admin").Parse(admintemplatetext)
	if err != nil {
		return err
	}
	_, err = admintemplate.Parse(headertemplatetext)
	if err != nil {
		return err
	}
	_, err = admintemplate.Parse(footertemplatetext)
	if err != nil {
		return err
	}
	mainpagetemplate, err = template.New("main").Parse(maintemplatetext)
	if err != nil {
		return err
	}
	_, err = mainpagetemplate.Parse(headertemplatetext)
	if err != nil {
		return err
	}
	_, err = mainpagetemplate.Parse(footertemplatetext)
	if err != nil {
		return err
	}
	leaderboardtemplate, err = template.New("leader").Parse(leaderboardtemplatetext)
	_, err = leaderboardtemplate.Parse(headertemplatetext)
	if err != nil {
		return err
	}
	_, err = leaderboardtemplate.Parse(footertemplatetext)
	if err != nil {
		return err
	}
	go leaderboardMessageServer(serverChan)
	// Http sturf
	r := mux.NewRouter()
	r.HandleFunc("/", mainpage)
	r.HandleFunc("/leaderboard", leaderboardpage)
	r.HandleFunc("/admin", adminpage)
	r.HandleFunc("/favicon.ico", favicon)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/register", register)
	r.HandleFunc("/submitflag", submitFlag)
	r.HandleFunc("/ws", leaderboardWS)
	r.HandleFunc("/{chall}", mainpage)
	r.HandleFunc("/detailview/{chall}", detailview)
	r.HandleFunc("/solutionview/{chall}", solutionview)
	r.HandleFunc("/getUserData/{user}", getUserData)
	r.HandleFunc("/uriview/{chall}", uriview)
	r.HandleFunc("/authorview/{chall}", authorview)
	// static
	r.PathPrefix("/static").Handler(
		http.FileServer(box))

	Port := config.Port
	if portenv := os.Getenv("WTFD_PORT"); portenv != "" {
		Port, _ = strconv.ParseInt(portenv, 10, 64)
	}
	fmt.Printf("WTFD Server Starting at port %d\n", Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", Port), r)
}
