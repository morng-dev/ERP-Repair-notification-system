package mail

import (
	"context"
	"fmt"

	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	mail "github.com/wneessen/go-mail"
)

type SmtpMailer struct {
	Host     string
	Username string
	Password string
	From     string
	AppUrl   string
	Port     int
}

func NewSmtpMailer(host, username, password, from, appUrl string, port int) services.Mailer {
	return &SmtpMailer{
		Host:     host,
		Username: username,
		Password: password,
		From:     from,
		AppUrl:   appUrl,
		Port:     port,
	}
}

func (m *SmtpMailer) SendResetPassword(ctx context.Context, to string, resetToken string) error {
	resetUrl := fmt.Sprintf("%s/reset-password?token=%s", m.AppUrl, resetToken)

	msg := mail.NewMsg()

	if err := msg.From(m.From); err != nil {
		return err
	}

	if err := msg.To(to); err != nil {
		return err
	}
	msg.Subject("รีเซตรหัสผ่านของคุณ")
	msg.SetBodyString(mail.TypeTextHTML, fmt.Sprintf(`
		<h2>รีเซ็ตรหัสผ่าน</h2>

		<p>เราได้รับคำขอรีเซ็ตรหัสผ่านสำหรับบัญชีของคุณ</p>

		<p><a href="%s">กดที่นี่เพื่อตั้งรหัสผ่านใหม่</a></p>

		<p>ลิงก์นี้จะหมดอายุใน 15 นาที</p>

		<p>หากคุณไม่ได้เป็นผู้ร้องขอ กรุณาเพิกเฉยต่ออีเมลนี้</p>
	`, resetUrl))

	client, err := mail.NewClient(
		m.Host,
		mail.WithPort(m.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(m.Username),
		mail.WithPassword(m.Password),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return err
	}
	return client.DialAndSendWithContext(ctx, msg)
}
