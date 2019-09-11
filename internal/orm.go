package wtfd

import (
	"os"
	"fmt"
	"errors"
    _ "github.com/mattn/go-sqlite3"
    "github.com/go-xorm/xorm"
	"xorm.io/core"
)

const (
	_ORMUserName            = "User"			// Table name for db
	_ORMChallengeByUserName = "ChallengeByUser" // Table name for db
)

var (
	engine *xorm.Engine
	ErrORMGeneric = errors.New("Database error (check log)")
)

////////////////////////////////////////////////////////////////////////////////
// ORM definitions /////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
type _ORMUser struct {
	Name      string `xorm:"unique"`
	Hash      []byte
	Mail	  string
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
	ormCreateTestDB()
}

func ormCreateTestDB() {
	u := &_ORMUser {
		Name: "TestUser",
		Hash: []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		Mail: "test@example.com",
	}

	cu := &_ORMChallengesByUser {
		UserName: u.Name,
		ChallengeName: "TestChallenge",
	}

	engine.Insert(u)
	engine.Insert(cu)
}

func _ORMGenericError(desc string) error {
	return errors.New(fmt.Sprintf("ORM Error %s", desc))
}

func ormNewUser(user User) error {
	count, err := engine.Count(_ORMUser{Name: user.Name,})
	if err != nil {
		return _ORMGenericError(err.Error())
	}

	if count == 0 {
		// TODO: Insert user
	} else {
		return ErrUserExisting
	}

	return nil
}
