package wtfd

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3" // needed for xorm
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	"xorm.io/core"
)

const (
	_ORMUserName            = "User"            // Table name for db
	_ORMChallengeByUserName = "ChallengeByUser" // Table name for db
)

var (
	engine *xorm.Engine
)

////////////////////////////////////////////////////////////////////////////////
// ORM definitions /////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
type _ORMUser struct {
	Name        string    `xorm:"unique"`
	DisplayName string    `xorm:"unique"`
	Created     time.Time `xorm:"created" json:"time"`
	Admin       bool
	Hash        []byte
	Points      int
	Verified    int
	VerifyToken string    `xorm:"varchar(32)"`
	VerifyDeadline time.Time
}

type _ORMChallengesByUser struct {
	UserName      string    `json:"username"` // Foregin keys don't exist
	Created       time.Time `xorm:"created" json:"time"`
	ChallengeName string    `json:"name"`
}

func (u _ORMUser) TableName() string {
	return _ORMUserName
}

func (c _ORMChallengesByUser) TableName() string {
	return _ORMChallengeByUserName
}

func ormSync() {
	_ = engine.Sync(_ORMUser{})
	_ = engine.Sync(_ORMChallengesByUser{})
}

////////////////////////////////////////////////////////////////////////////////

// Login checks if password is right for username and returns the User object of it
func Login(username, passwd string) error {
	user, err := Get(username)
	if err != nil {
		return err
	}
	if pwdRight := user.ComparePassword(passwd); !pwdRight {
		return errWrongPassword
	}
	fmt.Printf("User login: %s\n", username)
	return nil

}

// NewUser creates a new user object
func NewUser(name, password, displayname string) (User, error) {
	if Contains(name, displayname) {
		return User{}, errUserExisting
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}

	fmt.Printf("New User added: %s\n", name)
	isAdmin := ormGetUserCount() == 0
	if isAdmin {
		fmt.Printf("New User %s is an Admin\n", name)

	}
	return User{Name: name, Hash: hash, DisplayName: displayname, Admin: isAdmin, VerifiedInfo: VerifyInfo{IsVerified: false}}, nil

}

// Contains looks if a username is in the datenbank
func Contains(username, displayname string) bool {
	count, _ := ormUserExists(User{Name: username, DisplayName: displayname})
	return count
}

// Get gets username based on username
func Get(username string) (User, error) {
	user, err := ormLoadUser(username)
	if err != nil {
		fmt.Printf("Get Error: username: %v, user: %v, err: %v\n", username, user, err)
		return User{}, err
	}
	return user, err

}

func ormStart(logFile string) error {
	var err error
	engine, err = xorm.NewEngine("sqlite3", "./state.db")

	if err != nil {
		return err
	}

	if logFile != "" {
		f, err := os.Create(logFile)
		if err != nil {
			return err
		}
		engine.SetLogger(xorm.NewSimpleLogger(f))
	}

	engine.SetMapper(core.SameMapper{})

	ormSync()
	return nil
}

//noinspection GoUnusedFunction
func _ORMGenericError(desc string) error {
	return fmt.Errorf("ORM Error %s", desc)
}

////////////////////////////////////////////////////////////////////////////////
// DB Operations ///////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// Create new User DB record from a User struct
func ormNewUser(user User) error {
	exists, err := ormUserExists(user)
	if err != nil {
		return err
	}

	if exists {
		return errUserExisting
	}

	_, err = engine.Insert(_ORMUser{
		Name:        user.Name,
		Hash:        user.Hash,
		Admin:       user.Admin,
		DisplayName: user.DisplayName,
		Points:      0,
	})

	return err
}

// ormGetSolveCount returns the number of solves for the Challenge chall
func ormGetSolvesWithTime(u string) []_ORMChallengesByUser {

	var a []_ORMChallengesByUser
	if err := engine.Where("UserName = ?", u).Find(&a); err != nil {
		fmt.Printf("ORM Error: %v\n", err)
		return []_ORMChallengesByUser{}
	}
	return a

}

// ormGetSolveCount returns the number of solves for the Challenge chall
func ormGetUserCount() int64 {

	count, err := engine.Count(_ORMUser{})
	if err != nil {
		return 0
	}
	return count

}

// ormGetSolveCount returns the number of solves for the Challenge chall
func ormGetSolveCount(chall Challenge) int64 {

	count, err := engine.Where("ChallengeName = ?", chall.Name).Count(_ORMChallengesByUser{})
	if err != nil {
		return 0
	}
	return count

}

// Update existing user record (user.Name) with other values from user
// Solved challenges WON'T be updated (refer to ormChallengeSolved)
func ormUpdateUser(user User) error {
	var exists bool
	var err error
	var u _ORMUser

	if exists, err = ormUserExists(user); err != nil {
		return err
	}

	if !exists {
		return errUserNotExisting
	}

	verified := 0
	if user.VerifiedInfo.IsVerified {
		verified = 1
	}

	u = _ORMUser{
		Name:           user.Name,
		Hash:           user.Hash,
		DisplayName:    user.DisplayName,
		Points:         user.Points,
		Admin:          user.Admin,
		Verified:       verified,
		VerifyToken:    user.VerifiedInfo.VerifyToken,
		VerifyDeadline: user.VerifiedInfo.VerifyDeadline,
	}
	// fmt.Printf("a: %#v", u)

	if _, err = engine.Where("Name = ?", user.Name).Update(&u); err != nil {
		return err
	}

	return nil
}

// remove user from db (matches user.Name)
func ormDeleteUser(user User) error {
	var u _ORMUser
	var exists bool
	var err error

	if exists, err = ormUserExists(user); err != nil {
		return err
	}

	if !exists {
		return errUserNotExisting
	}

	if _, err = engine.Where("Name = ?", user.Name).Get(&u); err != nil {
		return err
	}

	if _, err = engine.Where("Name = ?", user.Name).Delete(&u); err != nil {
		return err
	}

	return nil
}

// check if user exists in db
func ormUserExists(user User) (bool, error) {
	count, err := engine.Where("Name = ?", user.Name).Count(_ORMUser{})
	//fmt.Printf("ormUserExists: user: %v, count: %v, err: %v\n", user, count, err)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}
	if count == 1 {
		return true, nil
	}
	return false, errors.New("DB User-table is in an invalid state")

}

func ormDisplayNameExists(name string) (bool, error) {
	count, err := engine.Where("DisplayName = ?", name).Count(_ORMUser{})
	//fmt.Printf("ormUserExists: user: %v, count: %v, err: %v\n", user, count, err)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}
	if count == 1 {
		return true, nil
	}
	return false, errors.New("DB User-table is in an invalid state")

}

// get Challenges{} solved by user (db-state)
func ormChallengesSolved(user User) ([]*Challenge, error) {
	var challenges []*Challenge

	_ = engine.Where("UserName = ?", user.Name).Iterate(_ORMChallengesByUser{}, func(i int, bean interface{}) error {
		relation := bean.(*_ORMChallengesByUser)

		for _, c := range challs {
			if c.Name == relation.ChallengeName {
				challenges = append(challenges, c)
			}
		}

		return nil
	})

	return challenges, nil
}

// Write solved state (user solved chall) in db
func ormSolvedChallenge(user User, chall *Challenge) error {
	var exists bool
	var err error
	var relation _ORMChallengesByUser
	var count int64

	if exists, err = ormUserExists(user); err != nil {
		return err
	}

	if !exists {
		return errUserNotExisting
	}

	count, err = engine.Where("UserName = ?", user.Name).And("ChallengeName = ?", chall.Name).Count(_ORMChallengesByUser{})
	if err != nil {
		return nil
	}

	if count == 1 {
		return errors.New("user already solved this challenge")
	} else if count > 1 {
		return errors.New("DB-State Error: User solved this challenge multiple times")
	}

	relation = _ORMChallengesByUser{
		ChallengeName: chall.Name,
		UserName:      user.Name,
	}

	_, err = engine.Insert(relation)

	return err
}

func ormAllUsersSortedByPoints() ([]_ORMUser, error) {
	var a []_ORMUser
	err := engine.Desc("Points").Find(&a)
	if err != nil {
		return a, err

	}
	return a, nil

}

// load a single user from db (search by name)
// The remaining fields of u will be filled by this function
func ormLoadUser(name string) (User, error) {
	var user _ORMUser
	var u User

	exists, err := ormUserExists(User{Name: name})
	// fmt.Printf("ormLoadUser: name: %v, exists: %v, err: %v\n",name, exists,err)
	if err != nil {
		return User{}, err
	}

	if !exists {
		return User{}, errUserNotExisting
	}

	if _, err := engine.Where("Name = ?", name).Get(&user); err != nil {
		// fmt.Printf("User %s seems to not exist", name)
		return User{}, err
	}

	verified := false
	if user.Verified == 1 {
		verified = true
	}

	u = User{
		Name:        user.Name,
		Hash:        user.Hash,
		DisplayName: user.DisplayName,
		Admin:       user.Admin,
		Points:      user.Points,
		VerifiedInfo: VerifyInfo {
			IsVerified: verified,
			VerifyToken: user.VerifyToken,
			VerifyDeadline: user.VerifyDeadline,
		},
	}

	if u.Completed, err = ormChallengesSolved(u); err != nil {
		return u, err
	}

	return u, nil
}

// get a username by a verify token.
// NOTE: race condition possible, when user changes properties of
// his own user object while verifying the token. Current architecture does
// not allow to solve this without expense
func ormUserByToken(token string) (User, error) {
	var user _ORMUser

	if _, err := engine.Where("VerifyToken = ?", token).Get(&user); err != nil {
		return User{}, err
	}

	// Load all data into real User struct
	return ormLoadUser(user.Name)
}

////////////////////////////////////////////////////////////////////////////////
