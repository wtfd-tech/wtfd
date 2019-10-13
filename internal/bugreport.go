package wtfd

import (
	"errors"
	"fmt"
	"time"
	"strconv"
	"net/smtp"
)


var (
	/* Runtime Parameter */        /* Defaults */
	BRServiceDeskDomain          = "example.com" // Server recieving service desk mails
	BRServiceDeskUser            = "noreply"     // The user recieving the mails
	BRServiceDeskPort            = 25            // server to server smtp port
	BRServiceDeskEnabled         = false         // Is service desk support enabled
	BRRateLimitInterval  float64 = 180           // 3 Minutes
	BRRateLimitReports           = 2             // 2 Reports during interval before beeing rate limited

	userAccess map[string]access = make(map[string]access)
)

type access struct {
	lastBlock  time.Time // Currently unused
	lastAccess []time.Time
}


/**
 * Check if user is rate limited
 */
func BRIsUserRateLimited(u *User) bool {
	record, ok := userAccess[u.Name]
	if !ok {
		return false
	}

	/* Ok if no critical ammount of records */
	if len(record.lastAccess) < BRRateLimitReports {
		return false
	}

	/* Check if earliest record is in interval, then block */
	if time.Since(record.lastAccess[0]).Seconds() < BRRateLimitInterval {
		return true
	}

	return false
}

/**
 * Register a user access
 */
func registerUserAccess(u *User) {
	record, ok := userAccess[u.Name]

	if !ok {
		/* New record */
		record = access{
			lastBlock:  time.Time{},
			lastAccess: []time.Time{time.Now()},
		}
	} else if len(record.lastAccess) < BRRateLimitReports {
		/* No critical ammount of records */
		record.lastAccess = append(record.lastAccess, time.Now())
	} else if !BRIsUserRateLimited(u) {
		/* Cycle access */
		record.lastAccess = record.lastAccess[1:]
		record.lastAccess = append(record.lastAccess, time.Now())
	}
	userAccess[u.Name] = record
}

/**
 * Send bugreport
 */
func BRDispatchBugreport(u *User, subject string, content string) error {
	if !BRServiceDeskEnabled {
		return errors.New("Service Desk is disabled")
	}

	recipient := BRServiceDeskUser + "@" + BRServiceDeskDomain
	recipients := []string{recipient}
	formatContent := fmt.Sprintf("From: %s\nSubject: %s\n\n%s", u.Name, subject, content)


	err := smtp.SendMail(BRServiceDeskDomain+":"+strconv.Itoa(BRServiceDeskPort),
		nil, u.Name, recipients, []byte(formatContent))
	if err == nil {
		registerUserAccess(u)
	}
	return err
}
