package mimemailer

import (
	"bytes"
	"mime/quotedprintable"
)

// converts text to quoted printable format re: RFC 2045
func convertToQuotedPrintable(text string) (string, error) {
	var outputBuffer bytes.Buffer
	w := quotedprintable.NewWriter(&outputBuffer)
	_, err := w.Write([]byte(text))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}
	converted := outputBuffer.String()
	return converted, nil
}

// Takes a byte array with mixed line endings and outputs one with all CRLF endigns
func makeAllCRLF(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	// Now that every line ending is a \n
	// replace every \n with \r
	d = bytes.Replace(d, []byte{10}, []byte{13, 10}, -1)
	return d
}