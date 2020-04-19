package types

import (
	"fmt"
)

// Challenges Array of challenges but in nice with funcitons
type Challenges []*Challenge

// Challenge is a challenge obv
type Challenge struct {
	Name        string `json:"name"`
	Title       string
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

// ChallengeYAML is Challenge as YAML
type ChallengeYAML struct {
	Name        string   `yaml:"name"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"desc"`
	Solution    string   `yaml:"solution"`
	Author      string   `yaml:"author"`
	Flag        string   `yaml:"flag"`
	Points      int      `yaml:"points"`
	URI         string   `yaml:"uri"`
	Deps        []string `yaml:"deps"`
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
