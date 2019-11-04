package wtfd

import (
	"fmt"
	"sort"
        "github.com/wtfd-tech/wtfd/internal/types"
        "github.com/wtfd-tech/wtfd/internal/db"

)


type gridinfo struct {
	Index int
	Pos   int
}


func resolveDeps(a []string) []*types.Challenge {
	var toReturn []*types.Challenge
	for _, b := range a {
		for _, c := range challs {
			if c.Name == b {
				toReturn = append(toReturn, c)
			}
		}
	}
	return toReturn

}

func countDeps(chall *types.Challenge) int {
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
	cols := make(map[int][]*types.Challenge)

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
				}
				// Sort as less (higher) if it has more dependecies
				return len(cols[i][x].DepIDs) > len(cols[i][y].DepIDs)
			}
			return cols[i][x].MinRow < cols[i][y].MinRow
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
func resolveChalls(jsons []*types.ChallengeJSON) {
	i := 0
	var idsInChalls []string
	for len(jsons) != 0 {
		//          fmt.Printf("challs: %v, jsons: %v\n",challs,jsons)
		this := jsons[i]
		if bContainsAllOfA(this.Deps, idsInChalls) {
			idsInChalls = append(idsInChalls, this.Name)
			challs = append(challs, &types.Challenge{Name: this.Name, Description: this.Description, Flag: this.Flag, URI: this.URI, Points: this.Points, Deps: resolveDeps(this.Deps), Solution: this.Solution, MinRow: -1, Row: -1, Author: this.Author})
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

func fixDeps(jsons []*types.ChallengeJSON) {
	challsByName := make(map[string]*types.ChallengeJSON)
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

// AllDepsCompleted checks if User u has completed all Dependent challenges of c
func AllDepsCompleted(u *db.User, c *types.Challenge) bool {
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

