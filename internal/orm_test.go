package wtfd

import (
	"fmt"
	"os"
	"testing"
)

var (
	u     User = User{Name: "testuser", Hash: []byte("a")}
	fakeu User = User{Name: "faketestuser", Hash: []byte("a")} // No Challenges because new users dont have any
)

// TestMain tests all the good ol orm stuff
func TestMain(m *testing.M) {
	err := ormStart("./testdblog")
	if err != nil {
		fmt.Printf("Database Creation failed: %v", err)
		os.Exit(1)
	}
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
