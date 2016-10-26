package email

import "net/smtp"
import "menud/components/attendees"
import "menud/database/connpool"
import "fmt"
import "strings"
import "menud/config"

func Send(att attendees.Attendee) (err error) {
	toList := make([]string, 1, 1)
	toList[0] = att.Email()

	evt, err := connpool.GetEvent(att.EventId())
	if err != nil {
		return
	}

	organiser, err := connpool.GetUser(evt.UserID())
	if err != nil {
		return
	}

	headers := make([]string, 0, 5)
	headers = append(headers, fmt.Sprintf("Subject: Menu choices for %s", evt.Name()))
	headers = append(headers, fmt.Sprintf("Sender: %s <%s>", organiser.Name(), organiser.Email()))

	body := fmt.Sprintf("Hi %s,\n\nYou've been invited to %s at %s on %s\n\nPlease choose your menu selections at: %s\n\nThanks!\n\n%s\n", att.Name(), evt.Name(), evt.Location(), evt.Date().Format("Mon 2 Jan at 15:04"), att.GetLoginURL(), organiser.Name())

	msg := fmt.Sprintf("%s\r\n\r\n%s", strings.Join(headers, "\r\n"), body)

	err = sendMail(config.MailServer(), config.MailSender(), toList, []byte(msg))
	return
}

func sendMail(addr string, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello("menud"); err != nil {
		return err
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
