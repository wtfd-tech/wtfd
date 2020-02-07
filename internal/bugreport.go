package wtfd

import (
	"errors"
	"github.com/wtfd-tech/wtfd/internal/db"
	"github.com/wtfd-tech/wtfd/internal/smtp"
	"time"
)

var (
	// BRServiceDeskAddress is the address where the Server is recieving service desk mails
	BRServiceDeskAddress = "mail@example.com"
	// BRServiceDeskEnabled sets that service desk support is enabled
	BRServiceDeskEnabled = false
	// BRRateLimitInterval sets the rate limiting interval, defaults to 3m
	BRRateLimitInterval float64 = 180
	// BRRateLimitReports sets the Reports during interval needed before beeing rate limited, defaults to 2
	BRRateLimitReports = 2
	// BRSMTPPort sets the server to server smtp port
	BRSMTPPort = 25
	// BRSMTPUser sets the server to server smtp user
	BRSMTPUser = "sender"
	// BRSMTPPassword sets the server to server smtp password
	BRSMTPPassword = "passwd"
	// BRSMTPHost sets the server to server smtp host
	BRSMTPHost = "example.com"

	userAccess map[string]access = make(map[string]access)
)

type access struct {
	lastBlock  time.Time // Currently unused
	lastAccess []time.Time
}

// BRIsUserRateLimited Checks if user is rate limited
func BRIsUserRateLimited(u *db.User) bool {
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
func registerUserAccess(u *db.User) {
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

// BRDispatchBugreport sends a bugreport
func BRDispatchBugreport(u *db.User, subject string, content string) error {
	if !BRServiceDeskEnabled {
		return errors.New("Service Desk is disabled")
	}
	if !smtp.Config.Enabled {
		return errors.New("SMTP is disabled")
	}

	fields := make(map[string]string)
	fields["ReplyTo"] = u.Name

	err := smtp.DispatchMail(BRServiceDeskAddress, subject, content, fields)
	if err == nil {
		registerUserAccess(u)
	}
	return err
}
