package mail

import (
	"testing"

	"github.com/Rexkizzy22/micro-bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmai(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test mail"
	content := `
		<h1>A test mail</h1>
		<p> This is a test mail from <a href="http://github.com/RexKizzy22" target="_blank">Kizito</a></p>
	`
	to := []string{"kizitoiriogbe@gmail.com"}
	attachFiles := []string{"../NOTE.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}