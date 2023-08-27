package mail

import (
	"testing"

	"github.com/SergeyPanov/bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewEmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test email SimpleBank"
	content := "test content"

	to := []string{"panovsy@gmail.com"}
	attachedFiles := []string{"../Makefile"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachedFiles)
	require.NoError(t, err)
}
