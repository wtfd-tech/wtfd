package wtfd

import (
	"fmt"
	"os"
	"testing"
)

var (
	u         User      = User{Name: "testuser", Hash: []byte("a")}
	fakeu     User      = User{Name: "faketestuser", Hash: []byte("a")} // No Challenges because new users dont have any
	testchall Challenge = Challenge{Name: "testchall", Flag: "testflag"}
)

// TestMain tests all the good ol orm stuff
func TestMain(m *testing.M) {
	os.Remove("./state.db")
	err := ormStart("./testdblog")
	if err != nil {
		fmt.Printf("Database Creation failed: %v", err)
		os.Exit(1)
	}
	challs = append(challs, testchall)
	os.Exit(m.Run())

}
func TestUserCreation(t *testing.T) {

	err := ormNewUser(u)
	if err != nil {
		t.Errorf("Database UserCreate failed: %v", err)

	}

}
func TestUserExists(t *testing.T) {

	hopefullyTrue, err := ormUserExists(u)
	if err != nil {
		t.Errorf("Database UserExists failed: %v", err)
	}
	if !hopefullyTrue {
		t.Errorf("Database UserExists failed: %v", fmt.Sprintf("Real User %v does not exist", u))
	}

	hopefullyFalse, err := ormUserExists(fakeu)
	if err != nil {
		t.Errorf("Database UserExists failed: %v", err)
	}
	if hopefullyFalse {
		t.Errorf("Database UserExists failed: %v", fmt.Sprintf("Fake User %v does exist", fakeu))
	}
}

func TestUserDeletion(t *testing.T) {

	deleteTestUser := User{Name: "deletetestuser"}
_:
	ormNewUser(deleteTestUser)
	err := ormDeleteUser(deleteTestUser)
	if err != nil {
		t.Errorf("Database UserDelete failed: %v", err)
	}

	hopefullyFalse, err := ormUserExists(deleteTestUser)
	if err != nil {
		t.Errorf("Database UserDelete failed: %v", err)
	}
	if hopefullyFalse {
		t.Errorf("Database UserDelete failed: %v", fmt.Sprintf("Deleted User %v does exist", deleteTestUser))
	}
}

func TestPullingUserFromDatabase(t *testing.T) {

	uFromDB, err := ormLoadUser(u.Name)
	if err != nil {
		t.Errorf("Database Pull failed: %v", err)
	}
	if !(uFromDB.Name == u.Name) {
		t.Errorf("Database Pull failed: %v", fmt.Sprintf("User %v != %v", u, uFromDB))
	}

}

func TestDeletingUserFromDatabase(t *testing.T) {

	err := ormDeleteUser(u)
	if err != nil {
		t.Errorf("Database Delete failed: %v", err)
	}
	hopefullyFalse, err := ormUserExists(u)
	if err != nil {
		t.Errorf("Database Delete failed: %v", err)
	}
	if hopefullyFalse {
		t.Errorf("Database UserDelete failed: %v", fmt.Sprintf("Deleted User %v does exist", u))
	}

}

func TestUserSolvesChallenge(t *testing.T) {

	solvetestuser := User{Name: "asdf"}
	_ = ormNewUser(solvetestuser)
	err := ormSolvedChallenge(solvetestuser, testchall)
	if err != nil {
		t.Errorf("Database Challenge Solve failed: %v", err)
	}
	uFromDB, err := ormLoadUser(solvetestuser.Name)
	fmt.Printf("%+v", uFromDB)
	if err != nil {
		t.Errorf("Database Challenge Solve failed: %v", err)
	}
	if !(uFromDB.Completed[0].Name == testchall.Name) {

		t.Errorf("Database Challenge Solve failed: %v",
			fmt.Sprintf("uFromDB.Completed[0].Name == testchall.Name, but uFromDB.Completed[0].Name = %v",
				uFromDB.Completed[0].Name))

	}
}
