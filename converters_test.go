package mimemailer

import (
	"testing"
	"bytes"
)

var (
	qpInput string
	qpExpected string
)

func init() {
	qpInput = "J'interdis aux marchands de vanter trop leur marchandises. Car ils se font vite pédagogues et t'enseignent comme but ce qui n'est par essence qu'un moyen, et te trompant ainsi sur la route à suivre les voilà bientôt qui te dégradent, car si leur musique est vulgaire ils te fabriquent pour te la vendre une âme vulgaire."
	// Result of the quoted printable function will be given in CRLF
	qpExpected = string(makeAllCRLF([]byte(`J'interdis aux marchands de vanter trop leur marchandises. Car ils se font =
vite p=C3=A9dagogues et t'enseignent comme but ce qui n'est par essence qu'=
un moyen, et te trompant ainsi sur la route =C3=A0 suivre les voil=C3=A0 bi=
ent=C3=B4t qui te d=C3=A9gradent, car si leur musique est vulgaire ils te f=
abriquent pour te la vendre une =C3=A2me vulgaire.`)))
}

func TestConvertToQuotedPrintable(t *testing.T) {
	result, error := convertToQuotedPrintable(qpInput)
	if error != nil {
		t.Error("Error running convertToQuotedPrintable")
		return
	}
	if result != qpExpected {
		t.Error("Error converting to quoted printable")
	}
}

func TestMakeAllCRLF(t *testing.T) {
	mixed := []byte{78, 23, 13, 10, 94, 10, 23, 13}
	expected := []byte{78, 23, 13, 10, 94, 13, 10, 23, 13, 10}
	result := makeAllCRLF(mixed)
	if bytes.Equal(result, expected) == false {
		t.Error("Error converting to all CRLF")
	}
}