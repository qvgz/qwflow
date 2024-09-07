package mail

import (
	"crypto/tls"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type Smtp struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Address string `json:"address"`
	PassWd  string `json:"passWd"`
	Name    string `json:"name"`
}

type Mail struct {
	To         []string // 接收邮箱
	Subject    string   // 主题
	Body       string   // 内容
	AttachPath string   // 附件路径
}

// 发送邮件
func (mail *Mail) Send(s *Smtp) error {
	// 收件人不能为空
	if len(mail.To) == 0 {
		return fmt.Errorf("%#v can not empty", mail.To)
	}

	m := gomail.NewMessage()

	if s.Name == "" {
		s.Name = s.Address
	}

	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(s.Address, s.Name)},
		"To":      mail.To,
		"Subject": {mail.Subject},
	})

	m.SetBody("text/html", mail.Body)

	if mail.AttachPath != "" {
		m.Attach(mail.AttachPath)
	}

	e := gomail.NewDialer(s.Host, s.Port, s.Address, s.PassWd)
	e.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := e.DialAndSend(m); err != nil {
		// 失败暂停 1s 重发
		time.Sleep(1 * time.Second)
		if err := e.DialAndSend(m); err != nil {
			return err
		}
	}
	return nil
}

// 单独发送邮件给每一个人
func (mail *Mail) SendAlone(s *Smtp) []error {
	var errList []error
	tmpMail := Mail{
		To:         make([]string, 1),
		Subject:    mail.Subject,
		Body:       mail.Body,
		AttachPath: mail.AttachPath,
	}
	for _, to := range mail.To {
		tmpMail.To[0] = to
		if err := tmpMail.Send(s); err != nil {
			errList = append(errList, err)
		}
		// 间隔 0.5 秒
		time.Sleep(500 * time.Millisecond)
	}
	return errList
}
