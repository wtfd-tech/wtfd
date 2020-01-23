package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3" // needed for xorm
	"github.com/wtfd-tech/wtfd/internal/types"
	"xorm.io/core"
)

const (
	_ORMUserName            = "User"            // Table name for db
	_ORMChallengeByUserName = "ChallengeByUser" // Table name for db
)

var (
	errUserExisting    = errors.New("user with this name exists")
	errWrongPassword   = errors.New("wrong Password")
	errUserNotExisting = errors.New("user with this name does not exist")
)

type xormdb struct {
	engine *xorm.Engine
}

////////////////////////////////////////////////////////////////////////////////
// ORM definitions /////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
type _ORMUser struct {
	Name           string    `xorm:"unique"`
	DisplayName    string    `xorm:"unique"`
	Created        time.Time `xorm:"created" json:"time"`
	Admin          int
	Hash           []byte
	Points         int
	Verified       int
	VerifyToken    string `xorm:"varchar(32)"`
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

func (d xormdb) ormSync() {
	_ = d.engine.Sync(_ORMUser{})
	_ = d.engine.Sync(_ORMChallengesByUser{})
}

////////////////////////////////////////////////////////////////////////////////

func (d *xormdb) Start() error {
	var err error
	d.engine, err = xorm.NewEngine("sqlite3", "./state.db")

	if err != nil {
		return err
	}

	d.engine.SetMapper(core.SameMapper{})

	d.ormSync()
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
func (d xormdb) NewUser(user User) error {
	exists, err := d.UserExists(user)
	if err != nil {
		return err
	}

	if exists {
		return errUserExisting
	}

	admin := 1
	if user.Admin {
		admin = 2
	}

	_, err = d.engine.Insert(_ORMUser{
		Name:        user.Name,
		Hash:        user.Hash,
		Admin:       admin,
		DisplayName: user.DisplayName,
		Points:      0,
	})

	return err
}

// ormGetSolveCount returns the number of solves for the Challenge chall
func (d xormdb) GetSolvesWithTime(u string) []_ORMChallengesByUser {

	var a []_ORMChallengesByUser
	if err := d.engine.Where("UserName = ?", u).Find(&a); err != nil {
		fmt.Printf("ORM Error: %v\n", err)
		return []_ORMChallengesByUser{}
	}
	return a

}

// ormGetSolveCount returns the number of solves for the Challenge chall
func (d xormdb) GetUserCount() int64 {

	count, err := d.engine.Count(_ORMUser{})
	if err != nil {
		return 0
	}
	return count

}

// ormGetSolveCount returns the number of solves for the Challenge chall
func (d xormdb) GetSolveCount(chall types.Challenge) int64 {

	count, err := d.engine.Where("ChallengeName = ?", chall.Name).Count(_ORMChallengesByUser{})
	if err != nil {
		return 0
	}
	return count

}

// Update existing user record (user.Name) with other values from user
// Solved challenges WON'T be updated (refer to ormChallengeSolved)
func (d xormdb) UpdateUser(user User) error {
	var exists bool
	var err error
	var u _ORMUser

	if exists, err = d.UserExists(user); err != nil {
		return err
	}

	if !exists {
		return errUserNotExisting
	}

	verified := 1
	if user.VerifiedInfo.IsVerified {
		verified = 2
	}

	admin := 1
	if user.Admin {
		admin = 2
	}

	u = _ORMUser{
		Name:           user.Name,
		Hash:           user.Hash,
		DisplayName:    user.DisplayName,
		Points:         user.Points,
		Admin:          admin,
		Verified:       verified,
		VerifyToken:    user.VerifiedInfo.VerifyToken,
		VerifyDeadline: user.VerifiedInfo.VerifyDeadline,
	}
	// fmt.Printf("a: %#v", u)

	if _, err = d.engine.Where("Name = ?", user.Name).Update(&u); err != nil {
		return err
	}

	return nil
}

// remove user from db (matches user.Name)
func (d xormdb) DeleteUser(user User) error {
	var u _ORMUser
	var exists bool
	var err error

	if exists, err = d.UserExists(user); err != nil {
		return err
	}

	if !exists {
		return errUserNotExisting
	}

	if _, err = d.engine.Where("Name = ?", user.Name).Get(&u); err != nil {
		return err
	}

	if _, err = d.engine.Where("Name = ?", user.Name).Delete(&u); err != nil {
		return err
	}

	return nil
}

// check if user exists in db
func (d xormdb) UserExists(user User) (bool, error) {
	count, err := d.engine.Where("Name = ?", user.Name).Count(_ORMUser{})
	//fmt.Printf("UserExists: user: %v, count: %v, err: %v\n", user, count, err)
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

func (d xormdb) DisplayNameExists(name string) (bool, error) {
	count, err := d.engine.Where("DisplayName = ?", name).Count(_ORMUser{})
	//fmt.Printf("UserExists: user: %v, count: %v, err: %v\n", user, count, err)
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
func (d xormdb) ChallengesSolved(user User) ([]*types.Challenge, error) {
	var challenges []*types.Challenge

	_ = d.engine.Where("UserName = ?", user.Name).Iterate(_ORMChallengesByUser{}, func(i int, bean interface{}) error {
		relation := bean.(*_ORMChallengesByUser)

		for _, c := range *challs {
			if c.Name == relation.ChallengeName {
				challenges = append(challenges, c)
			}
		}

		return nil
	})

	return challenges, nil
}
func (d xormdb) AllUsersSortedByPoints() ([]User, error) {
	var users []User

	_ = d.engine.Desc("Points").Iterate(_ORMUser{}, func(i int, bean interface{}) error {
		u := bean.(*_ORMUser)

		if u.Name != "" {
			user, _ := d.LoadUser(u.Name)
			//fmt.Printf("%s, %#v\n", user.Name, user.Created.String())

			users = append(users, user)
		}

		return nil
	})

	return users, nil

}

// Write solved state (user solved chall) in db
func (d xormdb) SolvedChallenge(user User, chall *types.Challenge) error {
	var exists bool
	var err error
	var relation _ORMChallengesByUser
	var count int64

	if exists, err = d.UserExists(user); err != nil {
		return err
	}

	if !exists {
		return errUserNotExisting
	}

	count, err = d.engine.Where("UserName = ?", user.Name).And("ChallengeName = ?", chall.Name).Count(_ORMChallengesByUser{})
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

	_, err = d.engine.Insert(relation)

	return err
}

// load a single user from db (search by name)
// The remaining fields of u will be filled by this function
func (d xormdb) LoadUser(name string) (User, error) {
	var user _ORMUser
	var u User

	exists, err := d.UserExists(User{Name: name})
	// fmt.Printf("ormLoadUser: name: %v, exists: %v, err: %v\n",name, exists,err)
	if err != nil {
		return User{}, err
	}

	if !exists {
		return User{}, errUserNotExisting
	}

	if _, err := d.engine.Where("Name = ?", name).Get(&user); err != nil {
		// fmt.Printf("User %s seems to not exist", name)
		return User{}, err
	}

	verified := false
	if user.Verified == 2 {
		verified = true
	}
	admin := false
	if user.Admin == 2 {
		admin = true
	}

	u = User{
		Name:        user.Name,
		Hash:        user.Hash,
		DisplayName: user.DisplayName,
		Admin:       admin,
		Points:      user.Points,
		Created:     user.Created,
		VerifiedInfo: VerifyInfo{
			IsVerified:     verified,
			VerifyToken:    user.VerifyToken,
			VerifyDeadline: user.VerifyDeadline,
		},
	}

	if u.Completed, err = d.ChallengesSolved(u); err != nil {
		return u, err
	}

	return u, nil
}

// ormUserByToken gets a username by a verify token.
// NOTE: race condition possible, when user changes properties of
// his own user object while verifying the token. Current architecture does
// not allow to solve this without expense
func (d xormdb) UserByToken(token string) (User, error) {
	var user _ORMUser

	if _, err := d.engine.Where("VerifyToken = ?", token).Get(&user); err != nil {
		return User{}, err
	}

	// Load all data into real User struct
	return d.LoadUser(user.Name)
}

// Contains looks if a username is in the datenbank
func (d xormdb) Contains(username, displayname string) bool {
	count, _ := d.UserExists(User{Name: username, DisplayName: displayname})
	return count
}

// Get gets username based on username
func (d xormdb) Get(username string) (User, error) {
	user, err := d.LoadUser(username)
	if err != nil {
		fmt.Printf("Get Error: username: %v, user: %v, err: %v\n", username, user, err)
		return User{}, err
	}
	return user, err

}

// Login checks if password is right for username and returns the User object of it
func (d xormdb) Login(username, passwd string) error {
	user, err := d.Get(username)
	if err != nil {
		return err
	}
	if pwdRight := user.ComparePassword(passwd); !pwdRight {
		return errWrongPassword
	}
	fmt.Printf("User login: %s\n", username)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
