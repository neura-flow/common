package mail

import (
	"context"
	"fmt"
	"testing"

	"github.com/neura-flow/common/log"
)

func TestSMTPSender_Send(t *testing.T) {
	sender := createSender()
	mail := &Mail{
		From:    "hello@163.com",
		Subject: "test",
		To:      []string{"hello2@163.com"},
		Text:    []byte("this.is.the.test.mail"),
	}
	if err := sender.Send(mail); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("email send success \n")
}

func createSender() *SMTPSender {
	cfg := &SMTPSendConfig{
		Host:      "smtp.163.com",
		Port:      465,
		User:      "",
		Pwd:       "",
		Timeout:   30,
		TLSEnable: true,
	}
	reader, err := NewSMTPSender(context.Background(), log.DefaultLogger(), cfg)
	if err != nil {
		panic(err)
	}
	return reader
}
