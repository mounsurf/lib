package notify

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

type Mail struct {
	Host     string
	Port     int
	Username string
	Email    string
	Password string
	Prefix   string
}

func (m *Mail) Send(email string, subject string, body string) error {
	subject = m.Prefix + subject
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Email)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	d := gomail.NewDialer(m.Host, m.Port, m.Username, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(msg)
}

/*
demo:

func main() {
	m := Mail{
		Host:     "smtp.163.com",
		Port:     25,
		Username: "test",
		Email:    "test@163.com",
		Password: "test",
	}
	err := m.Send("test@163.com", "subjectTest", "bodyTest")
	if err != nil {
		panic(err)
	}
}
*/
