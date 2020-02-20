package cfg

import (
	"html/template"
	"time"
)

type BugreportConfig struct {
	ServiceDeskAddress           string  `yaml:"address"`
	ServiceDeskRateLimitInterval float64 `yaml:"rateLimitInterval"` // See bugreport.go
	ServiceDeskRateLimitReports  int     `yaml:"rateLimitReports"`  // See bugreport.go

}

type EmailConfig struct {
	RestrictEmailDomains                 []string      `yaml:"allowedDomains"`
	RequireEmailVerification             bool          `yaml:"verification"`
	EmailVerificationTokenLifetimeString string        `yaml:"verificationTokenLifetime"`
	EmailVerificationTokenLifetime       time.Duration `yaml:"-"`
	SMTPRelayString                      string        `yaml:"smtpAddressWithPort"`
	SMTPRelayPasswd                      string        `yaml:"smtpPassword"`
}

type FooterLink struct {
	Name string `yaml:"name"`
	Icon string `yaml:"icon"`
	Url  string `yaml:"url"`
}

type DesignConfig struct {
	Icon        string        `yaml:"logo"`
	CoinIcon    string        `yaml:"coinicon"`
	Favicon     string        `yaml:"favicon"`
	UpperLeft   template.HTML `yaml:"logoText"`
	Header      string        `yaml:"title"`
	FooterLinks []FooterLink  `yaml:"links"`
	Slogan      template.HTML `yaml:"slogan"`
}

// Config stores settings
type Config struct {
	Port             int64           `yaml:"port"`
	StartDate        string          `yaml:"startDate"`
	Key              string          `yaml:"cookieKey"`
	ChallengeInfoDir string          `yaml:"challDir"`
	ChallHost        string          `yaml:"challHost"`
	BugreportConfig  BugreportConfig `yaml:"bugreport"`
	EmailConfig      EmailConfig     `yaml:"email"`
	DesignConfig     DesignConfig    `yaml:"design"`
}
