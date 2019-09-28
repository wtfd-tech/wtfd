package wtfd

// Challenges Array of challenges but in nice with funcitons
type Challenges []*Challenge

// Challenge is a challenge obv
type Challenge struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	Flag        string `json:"flag"`
	Points      int    `json:"points"`
	URI         string `json:"uri"`
	DepCount    int
	MinRow      int
	Row         int
	Solution    string `json:"solution"`
	DepIDs      []string
	Deps        []*Challenge
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

// Config stores settings loaded from config.json
type Config struct {
	ChallengeInfoDir	string	`json:"challinfodir"`
}

// User, was ist das wohl
type User struct {
	Name        string
	Hash        []byte
	DisplayName string
	Completed   []*Challenge
	Points      int
}

type leaderboardPageData struct {
	PageTitle     string
	User          *User
	IsUser        bool
	Points        int
	Leaderboard   bool
	AllUsers      []_ORMUser
	GeneratedName string
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
}
