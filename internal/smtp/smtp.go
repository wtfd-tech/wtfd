package smtp

import (
	"fmt"
	"strconv"
	"net/smtp"
)

/**
 * Dispatch a mail via send config
 */
func DispatchMail(recipient string, subject string, content string, fields map[string]string) error {

	auth := smtp.PlainAuth("", Config.User+"@"+Config.Host, Config.Password, Config.Host)

	recipients := []string{recipient}

	extraFields := ""
	for field, value := range(fields) {
		extraFields += fmt.Sprintf("%s: %s\r\n", field, value)
	}

	formatContent := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n%s\r\n",
		recipient, Config.User+"@"+Config.Host,  subject, extraFields, content)

	err := smtp.SendMail(Config.Host+":"+strconv.Itoa(Config.Port), auth,
		Config.User+"@"+Config.Host, recipients, []byte(formatContent))

	return err
}
