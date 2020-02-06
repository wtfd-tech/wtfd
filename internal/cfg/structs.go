package cfg

import (
	"html/template"
	"time"
)

type BugreportConfig struct {
	ServiceDeskAddress           string  `json:"servicedeskaddress"`
	ServiceDeskRateLimitInterval float64 `json:"servicedeskratelimitinterval"` // See bugreport.go
	ServiceDeskRateLimitReports  int     `json:"servicedeskratelimitreports"`  // See bugreport.go

}

type EmailConfig struct {
	RestrictEmailDomains                 []string      `json:"restrict_email_domains"`
	RequireEmailVerification             bool          `json:"require_email_verification"`
	EmailVerificationTokenLifetimeString string        `json:"email_verification_token_lifetime"`
	EmailVerificationTokenLifetime       time.Duration `json:"-"`
	SMTPRelayString                      string        `json:"smtprelaymailwithport"`
	SMTPRelayPasswd                      string        `json:"smtprelaymailpassword"`
}


type DesignConfig struct {
	Icon        string        `json:"icon"`
	CoinIcon    string        `json:"coinicon"`
	Favicon     string        `json:"favicon"`
	UpperLeft   template.HTML `json:"upperleft"`
	Header      string        `json:"header"`
	SocialMedia template.HTML `json:"social"`
}

// Config stores settings
type Config struct {
	Port             int64           `json:"port"`
	StartDate        string          `json:"startdate"`
	Key              string          `json:"key"`
	ChallengeInfoDir string          `json:"challinfodir"`
	SSHHost          string          `json:"sshhost"`
	BugreportConfig  BugreportConfig `json:"bugreport"`
	EmailConfig      EmailConfig     `json:"email"`
	DesignConfig     DesignConfig    `json:"design"`
}
