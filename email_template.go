package mimemailer

const emailTemplate = `Subject: {{ .Subject }}
From: {{ .From }}
To: {{ .To }}
Date: {{ .Date }}{{ if .ListUnsubscribe }}
List-Unsubscribe: {{ .ListUnsubscribe }}{{ end }}
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary=boundary42

--boundary42
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

{{/* This is a Quoted Printable encoded text version of the message */}}
{{ .TextQP }}

--boundary42
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

{{/* This is a Quoted Printable encoded html version of the message */}}
{{ .HTMLQP }}

--boundary42--
`
