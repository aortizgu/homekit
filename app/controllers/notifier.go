package controllers

import "net/smtp"

const (
	mailSource         string = "homekit.ortiz.gutierrez.adrian@gmail.com"
	mailSourcePassword string = "homekit.adrian"
)

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) serverName() string {
	return s.host + ":" + s.port
}

var (
	mailRecipients []string   = []string{"ortiz.gutierrez.adrian@gmail.com", "lauranton1592@gmail.com"}
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
	err := smtp.SendMail(server.serverName(), auth, from, to, message)
	if err != nil {
		return err
	}
	return nil
}

func SendMail(subject, msg string) error {
	for _, recipient := range mailRecipients {
		err := sendMail(recipient, subject, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func NotifyCalderaState(active, manual bool) {
	subject := "Calefacción "
	msg := "Se ha encendico la calefacción de manera "
	if active {
		subject += "Encendida"
	} else {
		subject += "Apagada"
	}
	if manual {
		msg += "manual"
	} else {
		msg += "automática"
	}
	SendMail(subject, msg)
}
