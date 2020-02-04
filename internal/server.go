//go:generate packr2
package wtfd

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wtfd-tech/wtfd/internal/cfg"
	"github.com/wtfd-tech/wtfd/internal/db"
	"github.com/wtfd-tech/wtfd/internal/smtp"
	"github.com/wtfd-tech/wtfd/internal/types"

	"github.com/gobuffalo/packr/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	config              cfg.Config
	store               sessions.Store
	wtfdDB              db.DB
	errUserExisting     = errors.New("user with this name exists")
	errWrongPassword    = errors.New("wrong Password")
	errUserNotExisting  = errors.New("user with this name does not exist")
	challs              = types.Challenges{}
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
	User          *db.User
	Config        cfg.Config
	IsUser        bool
	Points        int
	Leaderboard   bool
	AllUsers      []db.User
	GeneratedName string
	Style         template.HTMLAttr
	RowNums       []gridinfo
	ColNums       []gridinfo
}
type leaderboardPageData struct {
	PageTitle     string
	User          *db.User
	Config        cfg.Config
	IsUser        bool
	Points        int
	Leaderboard   bool
	AllUsers      []db.User
	GeneratedName string
	Style         template.HTMLAttr
	RowNums       []gridinfo
	ColNums       []gridinfo
}
type mainPageData struct {
	PageTitle              string
	Challenges             []*types.Challenge
	Leaderboard            bool
	SelectedChallengeID    string
	HasSelectedChallengeID bool
	ADC                    func(*db.User, *types.Challenge) bool
	GeneratedName          string
	Config                 cfg.Config
	User                   *db.User
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
	u, err := wtfdDB.LoadUser(vars["user"])
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error: %v", err)
	}
	userToReturn := db.User{Name: u.Name, DisplayName: u.DisplayName, Points: u.Points, Admin: u.Admin}
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
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Nice Try, %s", user.DisplayName)
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
		u := db.User{Name: r.FormValue("name"), DisplayName: r.FormValue("displayname"), Points: dumb, Admin: isAdmin}
		//                fmt.Printf("a: %#v",u)

		err = wtfdDB.UpdateUser(u)
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
	allUsers, err := wtfdDB.AllUsersSortedByPoints()
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
		Config:        config,
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
	allUsers, err := wtfdDB.AllUsersSortedByPoints()
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error: %v", err)
	}
	data := leaderboardPageData{
		PageTitle:     "foss-ag O-Phasen CTF",
		Config:        config,
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
		Config:                 config,
		Challenges:             challs,
		GeneratedName:          genu,
		ADC:                    AllDepsCompleted,
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
		if _, ok := getUser(r); ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Already logged in")
		} else {
			email := r.Form.Get("username")
			err := wtfdDB.Login(email, r.Form.Get("password"))
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

			if err = wtfdDB.SolvedChallenge(user, completedChallenge); err != nil {
				_ = fmt.Errorf("ORM Error: %s", err)
			}

			user.CalculatePoints()

			if err = wtfdDB.UpdateUser(user); err != nil {
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
		if _, ok := getUser(r); ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Server Error: Already logged in")
		} else {
			// username here means e-mail address
			if !validateEmailAddress(r.Form.Get("username")) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(w, "Server Error: The entered e-mail address is invalid")
			} else {
				// Check if registration is restricted to certain email domains
				if len(config.RestrictEmailDomains) != 0 {
					valid := false
					for _, domain := range config.RestrictEmailDomains {
						if strings.Split(r.Form.Get("username"), "@")[1] == domain {
							valid = true
						}
					}

					if !valid {
						w.WriteHeader(http.StatusBadRequest)
						_, _ = fmt.Fprintf(w, "Server Error: The entered e-mail address is not allowed")
						return
					}
				}

				u, err := db.NewUserStruct(wtfdDB, r.Form.Get("username"), r.Form.Get("password"), r.Form.Get("displayname"))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprintf(w, "Server Error: %v", err)
				} else {
					_ = wtfdDB.NewUser(u)
					login(w, r)
					_ = updateScoreboard()

				}

			}
		}

	}

}

func changePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Invalid Request")
	} else {
		if err := r.ParseForm(); err != nil {
			_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		// Check if user is logged in and get it
		if u, ok := getUser(r); ok {
			// Check if old password matches the entered one
			if bcrypt.CompareHashAndPassword(u.Hash, []byte(r.Form.Get("oldpassword"))) != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(w, "The old password entered is incorrect")
				fmt.Println("Old password wrong")
				return
			}

			// Check if both new passwords are the same
			if r.Form.Get("newpassword") != r.Form.Get("repeatnewpassword") {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(w, "The entered new password are not the same")
				fmt.Println("New passwords wrong")
				return
			}

			// Hash the entered password...
			hash, err := bcrypt.GenerateFromPassword([]byte(r.Form.Get("newpassword")), 14)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(w, "Server Error: %v", err)
				return
			}

			// ...and update it for the current user
			u.Hash = hash

			if wtfdDB.UpdateUser(u) != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(w, "Server Error: %v", err)
				return
			}

			fmt.Println("Done changing")
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "You have to be logged in to change your password")
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	_ = logoutUser(r, w)
	http.Redirect(w, r, "/", http.StatusFound)

}

func reportBug(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Invalid Request")
		return
	}

	/* Check user login */
	user, ok := getUser(r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Server Error: %v", "Not logged in")
		return
	}

	/* Check if user is rate limited */
	if BRIsUserRateLimited(&user) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = fmt.Fprint(w, "Too many requsets")
		return
	}

	/* Read and check form */
	if err = r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Server Error: %v", "Not logged in")
		return
	}
	subject := r.FormValue("subject")
	content := r.FormValue("content")
	if subject == "" || content == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invaild Request")
		return
	}
	/* Prevent field injection (assuming no injection in user.Name is possible) */
	if strings.ContainsRune(subject, '\n') {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invaild Request")
		return
	}

	/* Try to dispatch bugreport */
	if err = BRDispatchBugreport(&user, subject, content); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Server Error: %v", err)
		return
	}

	fmt.Fprint(w, "OK")
}

func requestVerify(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Invalid Request")
		return
	}

	/* Check user login */
	user, ok := getUser(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "Not logged in")
		return
	}

	/* Check user verification */
	if user.VerifiedInfo.IsVerified {
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprintln(w, "Already verified")
		return
	}

	/* Check rate limit */
	//TODO

	token := generateRandomString(32)
	content := fmt.Sprintf("Click here to verify your WTFd account: "+
		"http://%s/verify?token=%s\r\n\r\n"+
		"If you don't know about this, you can ignore this Mail", r.Host, token)

	/* Send mail */
	err = smtp.DispatchMail(user.Name, "WTFd Verification", content, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("SMTP error: %s\n", err.Error())
		return
	}

	/* Setup user info */
	user.VerifiedInfo.VerifyToken = token
	user.VerifiedInfo.VerifyDeadline = time.Now().Add(config.EmailVerificationTokenLifetime)
	if err = wtfdDB.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ORM error: %s\n", err.Error())
		return
	}

}

func verify(w http.ResponseWriter, r *http.Request) {
	var err error
	var user db.User

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Invalid Request")
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Invalid Request")
		return
	}

	/* load user */
	user, err = wtfdDB.UserByToken(token)
	if err != nil {
		// Any error just ends in an invalid token
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "Invalid Token")
		// print real error
		log.Printf("ORM Error: %s\n", err.Error())
		return
	}

	/* check verify deadline */
	if user.VerifiedInfo.VerifyDeadline.Before(time.Now()) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "Invalid Token")
		return
	}

	user.VerifiedInfo.IsVerified = true
	user.VerifiedInfo.VerifyDeadline = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local) // invalidate

	/* write user back to db */
	if err = wtfdDB.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ORM error: %s", err.Error())
		return
	}

	_, _ = fmt.Fprintf(w, "Succesfully verified \"%s\"\n", user.Name)
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
	_, _ = fmt.Fprintf(w, "%s<br><p>Solves: %d</p>", chall.Description, wtfdDB.GetSolveCount(*chall))

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
	utilInit()

	gob.Register(&db.User{})

	var err error
	config, err = cfg.GetConfig()
	if err != nil {
		return err
	}

	// setup servicedesk vars
	if config.ServiceDeskAddress == "-" {
		BRServiceDeskEnabled = false
	} else {
		BRServiceDeskAddress = config.ServiceDeskAddress
		smtp.Config.Password = config.SMTPRelayPasswd

		// Parse relay mail string
		split := strings.Split(config.SMTPRelayString, ":")

		if len(split) < 2 {
			return errors.New("Invalid smtprelaymailwithport format")
		}
		if smtp.Config.Port, err = strconv.Atoi(split[1]); err != nil {
			return err
		}
		split = strings.Split(split[0], "@")
		if len(split) < 2 {
			return errors.New("Invalid smtprelaymailwithport format")
		}
		smtp.Config.User = split[0]
		smtp.Config.Host = split[1]

		smtp.Config.Enabled = true
		BRServiceDeskEnabled = true
	}
	BRRateLimitReports = config.ServiceDeskRateLimitReports
	BRRateLimitInterval = config.ServiceDeskRateLimitInterval
	if BRServiceDeskEnabled {
		fmt.Printf("ServiceDesk running at %s (Send via %s@%s:%d)  (Max %dR/%.02fs)\n",
			BRServiceDeskAddress, smtp.Config.User, smtp.Config.Host, smtp.Config.Port,
			BRRateLimitReports, BRRateLimitInterval)
	} else {
		fmt.Println("ServiceDesk disabled")
	}

	// Email verification
	config.EmailVerificationTokenLifetime, err = time.ParseDuration(config.EmailVerificationTokenLifetimeString)
	if err != nil {
		log.Printf("Could not parse email_verification_lifetime: %s", err.Error())
		return err
	}

	key, err := base64.StdEncoding.DecodeString(config.Key)
	if err != nil {
		return err
	}

	store = sessions.NewFilesystemStore("", key) // generates filesystem store at os.tempdir

	//Load challs from dirs
	var challsStructure []*types.ChallengeJSON

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

			jsonStruct types.ChallengeJSON

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
	wtfdDB, err = db.StartDB(&challs)
	if err != nil {
		return err
	}

	// Fill in SSHHost
	challs.FillChallengeURI(config.SSHHost)
	// Packr

	box := packr.New("Box", "../html")
	maintemplatetext, err := box.FindString("index.html")
	if err != nil {
		return err
	}
	headertemplatetext, err := box.FindString("header.html")
	if err != nil {
		return err
	}
	footertemplatetext, err := box.FindString("footer.html")
	if err != nil {
		return err
	}
	admintemplatetext, err := box.FindString("admin.html")
	if err != nil {
		return err
	}
	leaderboardtemplatetext, err := box.FindString("leaderboard.html")
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
	r.HandleFunc("/changepassword", changePassword)
	r.HandleFunc("/submitflag", submitFlag)
	r.HandleFunc("/ws", leaderboardWS)
	r.HandleFunc("/reportbug", reportBug)
	r.HandleFunc("/request_verify", requestVerify)
	r.HandleFunc("/verify", verify)
	r.HandleFunc("/reportbug", reportBug)
	r.HandleFunc("/detailview/{chall}", detailview)
	r.HandleFunc("/solutionview/{chall}", solutionview)
	r.HandleFunc("/getUserData/{user}", getUserData)
	r.HandleFunc("/uriview/{chall}", uriview)
	r.HandleFunc("/authorview/{chall}", authorview)
	// static
	r.PathPrefix("/static").Handler(
		http.FileServer(box))
	r.HandleFunc("/dist/"+config.Icon, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, config.Icon)
	})
	r.HandleFunc("/dist/coinicon.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, config.CoinIcon)
	})
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, config.Favicon)
	})

	Port := config.Port
	if portenv := os.Getenv("WTFD_PORT"); portenv != "" {
		Port, _ = strconv.ParseInt(portenv, 10, 64)
	}
	fmt.Printf("WTFD Server Starting at port %d\n", Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", Port), r)
}
