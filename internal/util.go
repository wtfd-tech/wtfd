package wtfd

import (
	"strings"
	"strconv"
	"math/rand"
)

// https://stackoverflow.com/a/35099450
func stringCompareLess(si, sj string) bool {
	var siLower = strings.ToLower(si)
	var sjLower = strings.ToLower(sj)
	if siLower == sjLower {
		return si < sj
	}
	return siLower < sjLower
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
