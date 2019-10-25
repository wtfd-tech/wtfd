package wtfd

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"sort"
)

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
	Author      string `json:"author"`
	DepIDs      []string
	Deps        []*Challenge
	HasURI      bool // This emerges from URI != ""
}

// ChallengeJSON is Challenge as JSON
type ChallengeJSON struct {
	Name        string   `json:"name"`
	Description string   `json:"desc"`
	Solution    string   `json:"solution"`
	Author      string   `json:"author"`
	Flag        string   `json:"flag"`
	Points      int      `json:"points"`
	URI         string   `json:"uri"`
	Deps        []string `json:"deps"`
	HasURI      bool     // This emerges from URI != ""
}

// Config stores settings loaded from config.json
type Config struct {
	Port             int64  `json:port`
	Key              string `json:key`
	ChallengeInfoDir string `json:"challinfodir"`
	SSHHost          string `json:"sshhost"`
}

// User, was ist das wohl
type User struct {
  Name        string `json:"name"`
	Hash        []byte
	DisplayName string`json:"displayname"`
	Completed   []*Challenge
	Admin       bool `json:"admin"`
	Points      int `json:"points"`
}

type gridinfo struct {
	Index int
	Pos   int
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

// ComparePassword checks if the password is valid
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Hash, []byte(password)) == nil
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
				if len(cols[i][x].DepIDs) == len(cols[i][y].DepIDs) {
					return stringCompareLess(cols[i][x].Name, cols[i][y].Name)
				} else {
					// Sort as less (higher) if it has more dependecies
					return len(cols[i][x].DepIDs) > len(cols[i][y].DepIDs)
				}
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
func resolveChalls(jsons []*ChallengeJSON) {
	i := 0
	var idsInChalls []string
	for len(jsons) != 0 {
		//          fmt.Printf("challs: %v, jsons: %v\n",challs,jsons)
		this := jsons[i]
		if bContainsAllOfA(this.Deps, idsInChalls) {
			idsInChalls = append(idsInChalls, this.Name)
			challs = append(challs, &Challenge{Name: this.Name, Description: this.Description, Flag: this.Flag, URI: this.URI, Points: this.Points, Deps: resolveDeps(this.Deps), Solution: this.Solution, MinRow: -1, Row: -1, Author: this.Author})
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
