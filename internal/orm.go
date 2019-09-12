package wtfd

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"xorm.io/core"
)

const (
	_ORMUserName            = "User"            // Table name for db
	_ORMChallengeByUserName = "ChallengeByUser" // Table name for db
)

var (
	engine        *xorm.Engine
	ErrORMGeneric = errors.New("Database error (check log)")
)

////////////////////////////////////////////////////////////////////////////////
// ORM definitions /////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
type _ORMUser struct {
	Name string `xorm:"unique"`
	Hash []byte
}

type _ORMChallengesByUser struct {
	UserName      string // Foregin keys don't exist
	ChallengeName string
}

func (u _ORMUser) TableName() string {
	return _ORMUserName
}

func (c _ORMChallengesByUser) TableName() string {
	return _ORMChallengeByUserName
}

func ormSync() {
	engine.Sync(_ORMUser{})
	engine.Sync(_ORMChallengesByUser{})
}

////////////////////////////////////////////////////////////////////////////////

func ormStart(logFile string) {
	var err error
	engine, err = xorm.NewEngine("sqlite3", "./state.db")

	if err != nil {
		panic(fmt.Sprintf("Could not start xorm engine: %s\n", err.Error()))
	}

	if logFile != "" {
		f, err := os.Create(logFile)
		if err != nil {
			fmt.Errorf("Could not create DB Logfile: %s\n", err.Error())
		} else {
			engine.SetLogger(xorm.NewSimpleLogger(f))
		}
	}

	engine.SetMapper(core.SameMapper{})

	ormSync()
}

func _ORMGenericError(desc string) error {
	return errors.New(fmt.Sprintf("ORM Error %s", desc))
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
		return ErrUserExisting
	}

	_, err = engine.Insert(_ORMUser{
		Name: user.Name,
		Hash: user.Hash,
	})

	return err
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
		return ErrUserNotExisting
	}

	u = _ORMUser{
		Name: user.Name,
		Hash: user.Hash,
	}

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
		return ErrUserNotExisting
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
	var count int64
	var err error

	if count, err = engine.Count(_ORMUser{Name: user.Name}); err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}
	if count == 1 {
		return true, nil
	} else {
		return false, errors.New("DB User-table is in an invalid state")
	}

	return false, nil
}

// get Challenges{} solved by user (db-state)
func ormChallengesSolved(user User) ([]Challenge, error) {
	challenges := make([]Challenge, 0)

	engine.Where("UserName = ?", user.Name).Iterate(_ORMChallengesByUser{}, func(i int, bean interface{}) error {
		relation := bean.(*_ORMChallengesByUser)

		for i, _ := range challs {
			if challs[i].Name == relation.ChallengeName {
				challenges = append(challenges, challs[i])
			}
		}

		return nil
	})

	return challenges, nil
}

// Write solved state (user solved chall) in db
func ormSolvedChallenge(user User, chall Challenge) error {
	var exists bool
	var err error
	var relation _ORMChallengesByUser
	var count int64

	if exists, err = ormUserExists(user); err != nil {
		return err
	}

	if !exists {
		return ErrUserNotExisting
	}

	count, err = engine.Where("UserName = ?", user.Name).And("ChallengeName = ?", chall.Name).Count(_ORMChallengesByUser{})
	if err != nil {
		return nil
	}

	if count == 1 {
		return errors.New("User already solved this challenge")
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

// load a single user from db (search by name)
// The remaining fields of u will be filled by this function
func ormLoadUser(name string) (User, error) {
	var exists bool
	var err error
	var user _ORMUser
	var u User

	if exists, err = ormUserExists(User{Name: name}); err != nil {
		return User{}, err
	}

	if !exists {
		return User{}, ErrUserNotExisting
	}

	if _, err = engine.Where("Name = ?", name).Get(&user); err != nil {
		return User{}, err
	}

	u = User{
		Name: user.Name,
		Hash: user.Hash,
	}

	if u.Completed, err = ormChallengesSolved(u); err != nil {
		return u, err
	}

	return u, nil
}

////////////////////////////////////////////////////////////////////////////////
