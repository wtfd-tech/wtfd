package smtp

type SMTPConfig struct {
	Port     int     // server to server smtp port
	User     string  // user used for sending mails
	Password string  // password for user
	Host     string  // host where the send stmp server runns at
	Enabled  bool    // Is smtp service enabled
}

var Config SMTPConfig
