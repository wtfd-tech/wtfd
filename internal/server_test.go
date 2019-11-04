package wtfd

import (
	"encoding/json"
	"fmt"
	"testing"
	"github.com/wtfd-tech/wtfd/internal/types"
	"github.com/wtfd-tech/wtfd/internal/db"
)

func TestChallenges(t *testing.T) {
	jsonstring := `[
    {
      "name": "chall-a",
      "desc": "aioöhsogöaerghaeörkglkfjgaöoerilgeoörgijk",
      "flag": "FOSS{b}",
      "solution": "# You nice person, \nyou just gotta put the b to the a, obviously",
      "points": 8,
      "uri": "ssh://asdf@%s"
    },
    {
      "name": "chall-b",
      "deps": ["chall-a"],
      "desc": "aassddfaassddf",
      "flag": "FOSS{a}",
      "points": 1
    },
    {
      "name": "chall-c",
      "desc": "aassddfaassddf",
      "flag": "FOSS{a}",
      "deps": ["chall-b"],
      "points": 1,
      "uri": "ssh://asdf@%s"
    },
    {
      "name": "chall-d",
      "desc": "aassddfaassddf",
      "flag": "FOSS{a}",
      "deps": ["chall-a"],
      "points": 1,
      "uri": "ssh://asdf@%s"
    },
    {
      "name": "chall-f",
      "deps": ["chall-b","chall-c", "chall-d"],
      "desc": "aassddfaassddf",
      "flag": "FOSS{a}",
      "points": 1
    },
    {
      "name": "chall-h",
      "deps": ["chall-n","chall-c"],
      "desc": "aassddfaassddf",
      "flag": "FOSS{a}",
      "points": 1
    },
    {
      "name": "chall-n",
      "deps": ["chall-a"],
      "desc": "aassddfaassddf",
      "uri": "ssh://asdf@%s",
      "flag": "FOSS{a}",
      "points": 1
    },
    {
      "name": "chall-e",
      "deps": ["chall-b"],
      "desc": "aassddfaassddf",
      "flag": "FOSS{a}",
      "points": 1
    }
  ]`
	var challsStructure []*types.ChallengeJSON
	if err := json.Unmarshal([]byte(jsonstring), &challsStructure); err != nil {
		t.Errorf("JSON Unmarshal Error: %v", err)
	}
	resolveChalls(challsStructure)
	sshHost := "localhost:2222"
	challs.FillChallengeURI(sshHost)
	challn, err := challs.Find("chall-n")
	if err != nil {
		t.Errorf("Chall not found err: %v", err)
	}
	if challn.DepCount != 1 {
		t.Errorf("ChallDepCount is wrong: %v", fmt.Errorf("%v has %v deps instead of 2", challn, challn.DepCount))
	}
	if challn.DepIDs[0] != "chall-h" {
		t.Errorf("ChallDepID is wrong: %v", fmt.Errorf("%v has %v as reverse dep instead of chall-h", challn, challn.DepIDs[0]))
	}
	if !challn.HasURI {
		t.Errorf("ChallHasUri is wrong: %v", fmt.Errorf("challn.HasURI is false instead of true"))
	}
	if challn.URI != fmt.Sprintf("ssh://asdf@%s", sshHost) {
		t.Errorf("ChallUri is wrong: %v", fmt.Errorf("challn.URI is %v instead of 'ssh://asdf@%s'", challn.URI, sshHost))

	}
	t.Run("TestAllDepsCompleted", func(t *testing.T) {

		un := db.User{Name: "a", Completed: []*types.Challenge{{Name: "chall-b"}}}
		uy := db.User{Name: "a", Completed: []*types.Challenge{{Name: "chall-a"}}}
		if AllDepsCompleted(&un, challn) {
			t.Errorf("AllDepsCompleted is wrong: %v", fmt.Errorf("user hasn't completed dependent but AllDepsCompleted thinks it has"))
		}
		if !AllDepsCompleted(&uy, challn) {
			t.Errorf("AllDepsCompleted is wrong: %v", fmt.Errorf("user has completed dependent but AllDepsCompleted thinks it hasn't"))
		}

	})

}
