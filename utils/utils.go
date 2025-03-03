package utils

import (
	"fmt"

	"github.com/hex4coder/go-url-shortener/models"
)

type Shortener struct {
	Domain    string
	MaxLength int
}

func NewShortener(domain string, maxlength int) *Shortener {
	return &Shortener{
		Domain:    domain,
		MaxLength: maxlength,
	}
}

func (s *Shortener) ReadFile(filename string) (error, []*models.DataLink) {

	ec := make(chan error)
	dc := make(chan *DataExcel)

	go func() {
		err, d := ReadExcelFile(filename)

		if err != nil {
			ec <- err
			return
		}

		dc <- d
	}()

	select {
	case err := <-ec:
		return err, nil
	case data := <-dc:
		// processing data
		fmt.Printf("data from excel file: \r\n%v\r\n", data)
	}

	fmt.Println(filename)
	return fmt.Errorf("failed to read filename : %s", filename), nil
}
