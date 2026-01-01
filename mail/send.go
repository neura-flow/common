package mail

import (
	"context"
	"crypto/tls"

	"github.com/neura-flow/common/log"
	"gopkg.in/gomail.v2"
)

const (
	from    = "From"
	to      = "To"
	cc      = "Cc"
	subject = "Subject"
)

type SMTPSendConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	TLSEnable bool   `json:"tlsEnable"`
	Timeout   int    `json:"timeout"`
	User      string `json:"user"`
	Pwd       string `json:"pwd"`
}

type Mail struct {
	Subject string
	From    string
	To      []string
	Cc      []string
	Text    []byte
	HTML    []byte
}

type SMTPSender struct {
	ctx    context.Context
	logger log.Logger
	cfg    *SMTPSendConfig

	dialer *gomail.Dialer
}

func NewSMTPSender(ctx context.Context, logger log.Logger, cfg *SMTPSendConfig) (*SMTPSender, error) {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pwd)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &SMTPSender{
		ctx:    ctx,
		logger: logger,
		cfg:    cfg,
		dialer: dialer,
	}, nil
}

func (s *SMTPSender) Send(m *Mail) error {
	msg := gomail.NewMessage()
	msg.SetHeaders(map[string][]string{
		from:    {msg.FormatAddress(m.From, m.From)},
		to:      m.To,
		cc:      m.Cc,
		subject: {m.Subject},
	})
	if len(m.HTML) > 0 {
		msg.SetBody("text/html", string(m.HTML))
	} else {
		msg.SetBody("text/plain", string(m.Text))
	}
	return s.dialer.DialAndSend(msg)
}
