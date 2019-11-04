package smtp

type config struct {
	// Port
	Port int
	// User
	User string
	// Password
	Password string
	// Host
	Host string
	// Enabled
	Enabled bool
}

// Config is the config for the SMTP Server
var Config config
