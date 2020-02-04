package cfg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"time"

	"github.com/gorilla/securecookie"
)

const (
	defaultPort                 = int64(8080)
	bRRateLimitReports          = 2   // 2 Reports during interval before beeing rate limited
	bRRateLimitInterval float64 = 180 // 3 Minutes
)

// Config stores settings
type Config struct {
	Port                                 int64         `json:"port"`
	SocialMedia                          template.HTML `json:"social"`
	Icon                                 string        `json:"icon"`
	CoinIcon                             string        `json:"coinicon"`
	Favicon                              string        `json:"favicon"`
	UpperLeft                            template.HTML `json:"upperleft"`
	Key                                  string        `json:"key"`
	Header                               string        `json:"header"`
	ChallengeInfoDir                     string        `json:"challinfodir"`
	SSHHost                              string        `json:"sshhost"`
	ServiceDeskAddress                   string        `json:"servicedeskaddress"`
	SMTPRelayString                      string        `json:"smtprelaymailwithport"`
	SMTPRelayPasswd                      string        `json:"smtprelaymailpassword"`
	ServiceDeskRateLimitInterval         float64       `json:"servicedeskratelimitinterval"` // See bugreport.go
	ServiceDeskRateLimitReports          int           `json:"servicedeskratelimitreports"`  // See bugreport.go
	RestrictEmailDomains                 []string      `json:"restrict_email_domains"`
	RequireEmailVerification             bool          `json:"require_email_verification"`
	EmailVerificationTokenLifetimeString string        `json:"email_verification_token_lifetime"`
	EmailVerificationTokenLifetime       time.Duration `json:"-"`
}

// GetConfig returns a config struct generated either from a config.json or (TODO) from Environment
func GetConfig() (Config, error) {
	_, useEnv := os.LookupEnv("WTFD_USE_ENV_CONFIG")

	if useEnv {
		return getConfigEnv()
	}
	return getConfigJSON()
}

func getConfigJSON() (Config, error) {
	config := Config{}

	var key []byte

	//Test if config file exists
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		// Generate a new key
		key = securecookie.GenerateRandomKey(32)

		//Write default config to disk
		config = Config{
			Key:                                  base64.StdEncoding.EncodeToString(key),
			Port:                                 defaultPort,
			ChallengeInfoDir:                     "../challenges/info/",
			Header:                               "WTFd CTF",
			ServiceDeskAddress:                   "-", // service desk disabled
			SMTPRelayString:                      "mail@example.com:25",
			SMTPRelayPasswd:                      "passwd",
			ServiceDeskRateLimitReports:          bRRateLimitReports,
			ServiceDeskRateLimitInterval:         bRRateLimitInterval,
			SSHHost:                              "ctf.wtfd.tech",
			RestrictEmailDomains:                 nil,
			RequireEmailVerification:             false,
			SocialMedia:                          `<a class="link sociallink" href="https://github.com/wtfd-tech/wtfd"><span class="mdi mdi-github-circle"></span> GitHub</a>`,
			CoinIcon:                             "coinicon.svg",
			Favicon:                              "favicon.svg",
			Icon:                                 "icon.svg",
			UpperLeft:                            "// WTFd<br>//CTF",
			EmailVerificationTokenLifetimeString: "168h", // One week
		}
		configBytes, _ := json.MarshalIndent(config, "", "\t")
		err = ioutil.WriteFile("config.json", configBytes, os.FileMode(0600))
		if err != nil {
			return config, err
		}
	} else {
		//Load config file
		var (
			configBytes []byte
			err         error
		)

		if configBytes, err = ioutil.ReadFile("config.json"); err != nil {
			return config, err
		}
		if err := json.Unmarshal(configBytes, &config); err != nil {
			return config, err
		}
	}
	fmt.Fprintf(os.Stderr, "a: %v", config)
	return config, nil

}

func getConfigEnv() (Config, error) {
	return Config{}, fmt.Errorf("Environment Config not supported yet")
}
