package wtfd

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-" //64 chars
)

// Initialize utilities
func utilInit() {
	rand.Seed(time.Now().UnixNano())
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

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letterBytes[int(rand.Int() & 0b00111111)]
	}
	return string(b)
}

func generateUserName() (string, error) {
	var name string
	for _, s := range coolNames {
		if exists, err := wtfdDB.DisplayNameExists(s); !exists {
			if err != nil {
				return "", err

			}
			name = s
			break
		}
	}
	for name == "" {
		name = strconv.FormatInt(rand.Int63(), 10)
		if exists, err := wtfdDB.DisplayNameExists(name); !exists {
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

func validateEmailAddress(emailAddress string) bool {
	if strings.Count(emailAddress, "@") != 1 {
		return false
	}

	// Check if the e-mail address contains an @ symbol
	parts := strings.Split(emailAddress, "@")

	// Check if there are any character before and after the @ symbol
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}

	// Check if there is only one dot after the @ symbol
	if strings.Count(parts[1], ".") != 1 {
		return false
	}

	// Check if there is a dot followed by two or more chars and preceded by at least one char after the @ symbol
	partsAfterAt := strings.Split(parts[1], ".")

	if len(partsAfterAt[0]) < 1 || len(partsAfterAt[1]) < 2 {
		return false
	}

	// If all checks pass, the e-mail address is probably valid
	return true
}
