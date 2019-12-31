package notifier

import (
	"net/smtp"
	"time"
)

const (
	mailSource         string = ""
	mailSourcePassword string = ""
)

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) serverName() string {
	return s.host + ":" + s.port
}

var (
	mailRecipients []string   = []string{}
	server         smtpServer = smtpServer{host: "smtp.gmail.com", port: "587"}
)

func sendMail(recipient, subject, msg string) error {
	from := mailSource
	password := mailSourcePassword
	to := []string{recipient}
	message := []byte("To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		msg + "\r\n")
	auth := smtp.PlainAuth("", from, password, server.host)
	return smtp.SendMail(server.serverName(), auth, from, to, message)
}

// SendMail sends a mail to configured recipients
func SendMail(subject, msg string) error {
	for _, recipient := range mailRecipients {
		err := sendMail(recipient, subject, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// NotifyCalderaState sends a mail to configured recipients about the state of the device
func NotifyCalderaState(active, manual bool) {
	subject := "Calefacción "
	msg := "Se ha encendico la calefacción de manera "
	if active {
		subject += "Encendida "
	} else {
		subject += "Apagada "
		msg = "Se ha apagado la calefacción de manera "
	}
	subject += time.Now().Format("15:04:05")
	if manual {
		msg += "manual"
	} else {
		msg += "automática"
	}
	SendMail(subject, msg)
}
