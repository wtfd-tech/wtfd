package smtp

import (
	"fmt"
	"strconv"
	"net/smtp"
)

/**
 * Dispatch a mail via send config
 */
func DispatchMail(recipient string, replyTo string, subject string, content string, fields map[string]string) error {
	auth := smtp.PlainAuth("", Config.User+"@"+Config.Host, Config.Password, Config.Host)

	recipients := []string{recipient}

	formatContent := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nReply-To: %s\r\n\r\n%s\r\n",
		recipient, Config.User+"@"+Config.Host,  subject, replyTo, content)

	err := smtp.SendMail(Config.Host+":"+strconv.Itoa(Config.Port), auth,
		Config.User+"@"+Config.Host, recipients, []byte(formatContent))

	return err
}
