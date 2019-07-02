package mimemailer

import (
	"bytes"
	"github.com/kylelemons/godebug/diff"
	"net/mail"
	"testing"
	"time"
)

var (
	emailHTML     string
	expectedEmail []byte
)

func init() {
	emailHTML = `<!doctype html>
<html xmlns=http://www.w3.org/1999/xhtml style=background:#f3f3f3>
<body>
<p>Hello 🌍. This email wíll be formatted as a MIME message as per RFC 2045 and RFC 2046 📧</p>
</body>
</html>`
	expectedEmail = makeAllCRLF([]byte(`Subject: Test Email
From: "Sender Name" <sender@example.com>
To: "Test Recipient" <test@example.com>
Date: Tue, 10 Nov 2009 23:00:00 +0000
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary=boundary42

--boundary42
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable


Hello =F0=9F=8C=8D. This email w=C3=ADll be formatted as a MIME message as =
per RFC 2045 and RFC 2046 =F0=9F=93=A7

--boundary42
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable


<!doctype html>
<html xmlns=3Dhttp://www.w3.org/1999/xhtml style=3Dbackground:#f3f3f3>
<body>
<p>Hello =F0=9F=8C=8D. This email w=C3=ADll be formatted as a MIME message =
as per RFC 2045 and RFC 2046 =F0=9F=93=A7</p>
</body>
</html>

--boundary42--
`))
}

func TestNewMailer(t *testing.T) {
	_, err := NewMailer(Config{})
	if err != nil {
		t.Error("Error creating new mailer")
	}
}

func TestIsConnected(t *testing.T) {
	m, err := NewMailer(Config{})
	if err != nil {
		t.Error("Error creating new mailer")
	}
	c := m.IsConnected()
	if c != false {
		t.Error("Error on IsConnected")
	}
}

func TestMakeEmail(t *testing.T) {
	e := Email{
		ToAddress: "test@example.com",
		ToName:    "Test Recipient",
		Subject:   "Test Email",
		HTML:      emailHTML,
		Date:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	a := mail.Address{Name: "Sender Name", Address: "sender@example.com"}
	made, err := e.make(a)
	if err != nil {
		t.Error("Error on TestMakeEmail")
	}
	if bytes.Equal(made, expectedEmail) == false {
		t.Error("Error on TestMakeEmail - Email doesn't match expected output")
		t.Log(diff.Diff(string(made), string(expectedEmail)))
	}

}
