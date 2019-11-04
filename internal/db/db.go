package db

import (
	"fmt"
	"github.com/wtfd-tech/wtfd/internal/types"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	challs *types.Challenges
)

type DB interface {
	NewUser(u User) error
	GetSolvesWithTime(u string) []_ORMChallengesByUser
	Get(username string) (User, error)
	Contains(username, displayname string) bool
	Login(username, password string) error
	GetUserCount() int64
	GetSolveCount(chall types.Challenge) int64
	UpdateUser(u User) error
	DeleteUser(u User) error
	SolvedChallenge(u User, c *types.Challenge) error
	UserExists(u User) (bool, error)
	DisplayNameExists(n string) (bool, error)
	ChallengesSolved(u User) ([]*types.Challenge, error)
	LoadUser(username string) (User, error)
	UserByToken(token string) (User, error)
	AllUsersSortedByPoints() ([]User, error)
}

// VerifyInfo saves if a users e-mail is verified
type VerifyInfo struct {
	IsVerified     bool
	VerifyToken    string
	VerifyDeadline time.Time
}

// User was ist das wohl
type User struct {
	Name         string `json:"name"`
	Hash         []byte
	DisplayName  string `json:"displayname"`
	Completed    []*types.Challenge
	Admin        bool `json:"admin"`
	Points       int  `json:"points"`
	VerifiedInfo VerifyInfo
	Created      time.Time
}

// StartDB starts the db and returns a DB interface
func StartDB(c *types.Challenges) (DB, error) {
	challs = c
	x := xormdb{}
	if err := x.Start(); err != nil {
		return xormdb{}, err
	}
	return x, nil
}

// HasSolvedChallenge returns true if u has solved chall
func (u User) HasSolvedChallenge(chall *types.Challenge) bool {
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

// ComparePassword checks if the password is valid
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Hash, []byte(password)) == nil
}

// NewUserStruct creates a new user object
func NewUserStruct(d DB, name, password, displayname string) (User, error) {
	if d.Contains(name, displayname) {
		return User{}, errUserExisting
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}

	fmt.Printf("New User added: %s\n", name)
	isAdmin := d.GetUserCount() == 0
	if isAdmin {
		fmt.Printf("New User %s is an Admin\n", name)

	}
	return User{Name: name, Hash: hash, DisplayName: displayname, Admin: isAdmin, VerifiedInfo: VerifyInfo{IsVerified: false}}, nil
}
