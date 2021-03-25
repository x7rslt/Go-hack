package zip_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/alexmullins/zip"
)

func TestCreatzip(t *testing.T) {
	contents := []byte("hello world")
	fzip, err := os.Create(`./test.zip`)
	if err != nil {
		log.Fatal(err)
	}
	zipw := zip.NewWriter(fzip)
	defer zipw.Close()
	w, err := zipw.Encrypt(`test.txt`, `golang`) //compress test.txt with password "golang"
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(w, bytes.NewReader(contents))
	if err != nil {
		log.Fatal(err)
	}
	zipw.Flush()
}

func TestUnzip(t *testing.T) {
	zipfile := `./test.zip`

	zipr, err := zip.OpenReader(zipfile)
	if err != nil {
		log.Fatal(err)
	}
	for _, z := range zipr.File {
		z.SetPassword("golang")
		rr, err := z.Open()
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(os.Stdout, rr)
		if err != nil {
			log.Fatal(err)
		}
		rr.Close()
	}

}
