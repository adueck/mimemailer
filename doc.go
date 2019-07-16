// Copyright 2019 Adam Dueck. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package mimemailer provides a way to easily create send multi-part
// MIME email messages as specified by RFC 2045 and RFC 2046.
//
// This package takes an email with a regular HTML body, converts it
// into a multipart MIME message as per RFC 2045 and RFC 2046, and
// sends it via SMTP.
//
// Given an email like this:
//
// 	Email{
// 		ToAddress: "test@example.com",
// 		ToName: "Test Recipient",
// 		Subject: "Test Email",
// 		Date: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
//		HTML: `<!doctype html>
//	<html xmlns=http://www.w3.org/1999/xhtml style=background:#f3f3f3>
//	<body>
//	<p>Hello 🌍. This email wíll be formatted as a MIME message as per RFC 2045 and RFC 2046 📧</p>
//	</body>
//	</html>`,
//	}
//
//
// It will generate and send a multi-part MIME message over SMTP like this:
//
// 	Subject: Test Email
// 	From: "Sender Name" <sender@example.com>
// 	To: "Test Recipient" <test@example.com>
// 	Date: Tue, 10 Nov 2009 23:00:00 +0000
// 	MIME-Version: 1.0
// 	Content-Type: multipart/alternative; boundary=boundary42
//
// 	--boundary42
// 	Content-Type: text/plain; charset=utf-8
// 	Content-Transfer-Encoding: quoted-printable
//
//
// 	Hello =F0=9F=8C=8D. This email w=C3=ADll be formatted as a MIME message as =
// 	per RFC 2045 and RFC 2046 =F0=9F=93=A7
//
// 	--boundary42
// 	Content-Type: text/html; charset=utf-8
// 	Content-Transfer-Encoding: quoted-printable
//
//
// 	<!doctype html>
// 	<html xmlns=3Dhttp://www.w3.org/1999/xhtml style=3Dbackground:#f3f3f3>
// 	<body>
// 	<p>Hello =F0=9F=8C=8D. This email w=C3=ADll be formatted as a MIME message =
// 	as per RFC 2045 and RFC 2046 =F0=9F=93=A7</p>
// 	</body>
// 	</html>
//
//	--boundary42--
//
// The text version is generated automatically, and the message is converted
// into quoted printable format with CRLF line endings.
//
// After connecting to a SMTP server, you can send multiple messages and re-use
// the same connection.
//
// Example Usage:
//
//	import (
//		"log"
//		"github.com/adueck/mimemailer"
//	)
//
// 	func main() {
//		// Step 1 - Create a new mailer instance with config for SMTP
// 		m, err := mimemailer.NewMailer(mimemailer.Config{
//			Host: 			"smtp.example.com",
//			Port: 			"576",
//			Username: 		"myusername",
//			Password: 		"mysecretpassword",
//			SenderName:		"My Name",
//			SenderAddress:		"email@example.com",
//		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
//		// Step 2 - Connect to SMTP Server
//		err = m.Connect()
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		// Step 3 - Send message(s)
//		err = m.SendEmail(mimemailer.Email{
//			ToAddress: 	"recipient@example.com",
//			ToName:		"Bob Smith",
//			Subject:	"Example Mail",
//			HTML:		"<html><p>Hello Bob</p></html>",
//			Date:		time.Now(),
//		})
//		if err != nil {
//			log.Print(err)
//		}
//
//		// ...
//		// You can send more messages on the persistent connection
//		// When you are done, disconnect from the SMTP server
//
//		// Step 4 - Disconnect from SMTP server
//		err = m.Disconnect()
//		if err != nil {
//			log.Fatal(err)
//		}
// 	}
package mimemailer
