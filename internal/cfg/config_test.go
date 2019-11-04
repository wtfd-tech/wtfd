package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestJsonConfigGeneration(t *testing.T) {
	os.Remove("config.json") // Cleanup
	generatedConfig, err := getConfigJson()
	if err != nil {
		t.Errorf("Config Generation failed with error: %v", err)
	}
	cjson, err := ioutil.ReadFile("config.json")
	if err != nil {
		t.Errorf("Config File Reading failed with error: %v", err)
	}
	configFromJson := &Config{}
	err = json.Unmarshal(cjson, configFromJson)
	if err != nil {
		t.Errorf("Config File Reading failed with error: %v", err)
	}
	if generatedConfig.Port != configFromJson.Port ||
		generatedConfig.RequireEmailVerification != configFromJson.RequireEmailVerification ||
		generatedConfig.SocialMedia != configFromJson.SocialMedia ||
		generatedConfig.Icon != configFromJson.Icon ||
		generatedConfig.FirstLine != configFromJson.FirstLine ||
		generatedConfig.SecondLine != configFromJson.SecondLine ||
		generatedConfig.Key != configFromJson.Key ||
		generatedConfig.ChallengeInfoDir != configFromJson.ChallengeInfoDir ||
		generatedConfig.SSHHost != configFromJson.SSHHost ||
		generatedConfig.ServiceDeskAddress != configFromJson.ServiceDeskAddress ||
		generatedConfig.SMTPRelayString != configFromJson.SMTPRelayString ||
		generatedConfig.SMTPRelayPasswd != configFromJson.SMTPRelayPasswd ||
		generatedConfig.ServiceDeskRateLimitInterval != configFromJson.ServiceDeskRateLimitInterval ||
		generatedConfig.ServiceDeskRateLimitReports != configFromJson.ServiceDeskRateLimitReports ||
		generatedConfig.EmailVerificationTokenLifetimeString != configFromJson.EmailVerificationTokenLifetimeString ||
		generatedConfig.EmailVerificationTokenLifetime != configFromJson.EmailVerificationTokenLifetime {

		t.Errorf("Config File is not the same as generated Config\ngeneratedConfig: %v\nconfigFromJson: %v", generatedConfig, configFromJson)
	}
	os.Remove("config.json") // Cleanup

}

func TestJsonConfigReading(t *testing.T) {
        configstring := []byte(`{
	"Port": 8080,
	"social": "\u003ca class=\"link sociallink\" href=\"https://github.com/wtfd-tech/wtfd\"\u003e\u003cspan class=\"mdi mdi-github-circle\"\u003e\u003c/span\u003e GitHub\u003c/a\u003e",
	"icon": "icon.svg",
	"firstline": "WTFd",
	"secondline": "CTF",
	"Key": "ED+vKjuFycJk9WQ3jc4GRyeOSGXUOloONxlD9qw8USk=",
	"challinfodir": "../challenges/info/",
	"sshhost": "ctf.wtfd.tech",
	"servicedeskaddress": "-",
	"smtprelaymailwithport": "mail@example.com:25",
	"smtprelaymailpassword": "passwd",
	"ServiceDeskRateLimitInterval": 180,
	"ServiceDeskRateLimitReports": 2,
	"restrict_email_domains": null,
	"require_email_verification": false,
	"email_verification_token_lifetime": "168h"
}`)

	os.Remove("config.json") // Cleanup
        err := ioutil.WriteFile("config.json", configstring, os.FileMode(0600))
	if err != nil {
		t.Errorf("config file writing failed with error: %v", err)
	}
        generatedConfig, err := getConfigJson()
	if err != nil {
		t.Errorf("getConfigJson failed with error: %v", err)
	}
        configFromJson := Config{}
        err = json.Unmarshal(configstring, &configFromJson)
	if err != nil {
		t.Errorf("Config Json Genereation failed with error: %v", err)
	}
	if generatedConfig.Port != configFromJson.Port ||
		generatedConfig.RequireEmailVerification != configFromJson.RequireEmailVerification ||
		generatedConfig.SocialMedia != configFromJson.SocialMedia ||
		generatedConfig.Icon != configFromJson.Icon ||
		generatedConfig.FirstLine != configFromJson.FirstLine ||
		generatedConfig.SecondLine != configFromJson.SecondLine ||
		generatedConfig.Key != configFromJson.Key ||
		generatedConfig.ChallengeInfoDir != configFromJson.ChallengeInfoDir ||
		generatedConfig.SSHHost != configFromJson.SSHHost ||
		generatedConfig.ServiceDeskAddress != configFromJson.ServiceDeskAddress ||
		generatedConfig.SMTPRelayString != configFromJson.SMTPRelayString ||
		generatedConfig.SMTPRelayPasswd != configFromJson.SMTPRelayPasswd ||
		generatedConfig.ServiceDeskRateLimitInterval != configFromJson.ServiceDeskRateLimitInterval ||
		generatedConfig.ServiceDeskRateLimitReports != configFromJson.ServiceDeskRateLimitReports ||
		generatedConfig.EmailVerificationTokenLifetimeString != configFromJson.EmailVerificationTokenLifetimeString ||
		generatedConfig.EmailVerificationTokenLifetime != configFromJson.EmailVerificationTokenLifetime {

		t.Errorf("Config String is not the same as generated Config\ngeneratedConfig: %v\nconfigFromJson: %v", generatedConfig, configFromJson)
	}
	os.Remove("config.json") // Cleanup

}

func TestEnvConfig(t *testing.T) {

}
