package cfg

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"gopkg.in/yaml.v2"
)

const (
	defaultPort                 = int64(8080)
	bRRateLimitReports          = 2   // 2 Reports during interval before beeing rate limited
	bRRateLimitInterval float64 = 180 // 3 Minutes
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
	SocialMedia template.HTML `json:"social"`
	Icon        string        `json:"icon"`
	CoinIcon    string        `json:"coinicon"`
	Favicon     string        `json:"favicon"`
	UpperLeft   template.HTML `json:"upperleft"`
	Header      string        `json:"header"`
}

// Config stores settings
type Config struct {
	Port             int64           `json:"port"`
	StartDate        string    `json:"startdate"`
	Key              string          `json:"key"`
	ChallengeInfoDir string          `json:"challinfodir"`
	SSHHost          string          `json:"sshhost"`
	BugreportConfig  BugreportConfig `json:"bugreport"`
	EmailConfig      EmailConfig     `json:"email"`
	DesignConfig     DesignConfig    `json:"design"`
}

// GetConfig returns a config struct generated either from a config.json or (TODO) from Environment
func GetConfig() (Config, error) {
	_, useEnv := os.LookupEnv("WTFD_USE_ENV_CONFIG")

	if useEnv {
		return getConfigEnv()
	}
	return getConfigYAML()
}

func getConfigYAML() (Config, error) {
	config := Config{}

	var key []byte

	//Test if config file exists
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		// Generate a new key
		key = securecookie.GenerateRandomKey(32)

		//Write default config to disk
		config = Config{
			Key:              base64.StdEncoding.EncodeToString(key),
			SSHHost:          "ctf.wtfd.tech",
			Port:             defaultPort,
                        StartDate: time.Now().Format(time.RubyDate),
			ChallengeInfoDir: "../challenges/info/",
			BugreportConfig: BugreportConfig{
				ServiceDeskAddress:           "-", // service desk disabled
				ServiceDeskRateLimitReports:  bRRateLimitReports,
				ServiceDeskRateLimitInterval: bRRateLimitInterval,
			},
			EmailConfig: EmailConfig{
				SMTPRelayString:                      "mail@example.com:25",
				SMTPRelayPasswd:                      "passwd",
				EmailVerificationTokenLifetimeString: "168h", // One week
				RestrictEmailDomains:                 nil,
				RequireEmailVerification:             false,
			},
			DesignConfig: DesignConfig{
				Header:      "WTFd CTF",
				SocialMedia: `<a class="link sociallink" href="https://github.com/wtfd-tech/wtfd"><span class="mdi mdi-github-circle"></span> GitHub</a>`,
				CoinIcon:    "coinicon.svg",
				Favicon:     "favicon.svg",
				Icon:        "icon.svg",
				UpperLeft:   "// WTFd<br>//CTF",
			},
		}
		configBytes, _ := yaml.Marshal(config)
		err = ioutil.WriteFile("config.yaml", configBytes, os.FileMode(0600))
		if err != nil {
			return config, err
		}
	} else {
		//Load config file
		var (
			configBytes []byte
			err         error
		)

		if configBytes, err = ioutil.ReadFile("config.yaml"); err != nil {
			return config, err
		}
		if err := yaml.Unmarshal(configBytes, &config); err != nil {
			return config, err
		}
	}
        if _, err := time.Parse(time.RubyDate, config.StartDate); err != nil {
			return config, err
        }
	return config, nil

}

func getConfigEnv() (Config, error) {
	return Config{}, fmt.Errorf("Environment Config not supported yet")
}
