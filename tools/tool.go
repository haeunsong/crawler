package tools

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func CheckCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: ", res.StatusCode)
	}
}

func EucKrReaderToUtf8Reader(body io.Reader) io.Reader {
	rInUTF8 := transform.NewReader(body, korean.EUCKR.NewDecoder())
	decBytes, _ := ioutil.ReadAll(rInUTF8)
	decrypted := string(decBytes)
	return strings.NewReader(decrypted)
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
