package types

import (
	"fmt"
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
