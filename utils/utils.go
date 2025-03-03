package utils

import (
	"fmt"
	"sync"

	"github.com/hex4coder/go-url-shortener/models"
)

type Shortener struct {
	Domain    string
	MaxLength int
}

type ShortenerI interface {
	ReadFile(string) (error, []*models.DataLink)
	MappingToDataLinks(*DataExcel) (error, []*models.DataLink)
	GenerateShortLink([]*models.DataLink) (error, []*models.ShortLink)
}

func NewShortener(domain string, maxlength int) *Shortener {
	return &Shortener{
		Domain:    domain,
		MaxLength: maxlength,
	}
}

func (s *Shortener) MappingToDataLinks(dataexcel *DataExcel) (error, []*models.DataLink) {
	links := []*models.DataLink{}

	if len(dataexcel.Rows) < 2 {
		return fmt.Errorf("tidak ada data dalam file"), nil
	}

	for i, row := range dataexcel.Rows {
		if i < 1 {
			continue
		}

		item := new(models.DataLink)
		item.Lesson = row[2]
		item.Teacher = row[3]
		item.ClassInfo = row[4]
		item.LongUrl = row[5]
		item.Token = row[6]

		if len(item.LongUrl) > 0 {
			links = append(links, item)
		}
	}

	return nil, links
}

func (s *Shortener) ReadFile(filename string) (error, []*models.DataLink) {

	ec := make(chan error)
	dc := make(chan *DataExcel)
	dlc := make(chan *models.Links)

	go func() {
		err, d := ReadExcelFile(filename)

		if err != nil {
			ec <- err
			return
		}

		dc <- d
	}()

	go func(ec chan error, dlc chan *models.Links) {

		select {
		case data := <-dc:
			// run pre processing data
			go func() {
				err, links := s.MappingToDataLinks(data)

				if err != nil {
					ec <- err
					return
				}

				dlc <- &models.Links{DataLinks: links}
			}()
		}

	}(ec, dlc)

	select {

	case err := <-ec:
		return err, nil

	case links := <-dlc:
		return nil, links.DataLinks

	}
}

func (s *Shortener) GenerateShortLink(datalinks []*models.DataLink) (error, []*models.ShortLink) {

	var wg sync.WaitGroup

	type GeneratedShortLinks struct {
		mu         sync.Mutex
		shortlinks []*models.ShortLink
	}

	ds := new(GeneratedShortLinks)
	dc := make(chan *models.ShortLink)
	done := make(chan any)

	go func(listdata []*models.DataLink) {
		for {
			select {
			case <-done:
				return

			case shortlink := <-dc:
				func(newSL *models.ShortLink) {
					ds.mu.Lock()
					ds.shortlinks = append(ds.shortlinks, newSL)

					if len(ds.shortlinks) >= len(listdata) {
						done <- true
					}

					defer ds.mu.Unlock()
				}(shortlink)

			}
		}

	}(datalinks)

	for i, dl := range datalinks {
		wg.Add(1)
		go func(wg *sync.WaitGroup, cc chan *models.ShortLink, index int) {
			defer wg.Done()

			// create new short link
			randomString := GenerateRandomString(s.MaxLength)
			shortlink := fmt.Sprintf("%s/%s", s.Domain, randomString)

			// create new data
			sl := new(models.ShortLink)
			sl.ID = index + 1
			sl.UniqueCode = randomString
			sl.Data = dl
			sl.ShortUrl = shortlink

			// send data to channel
			cc <- sl
		}(&wg, dc, i)
	}

	wg.Wait()

	return nil, ds.shortlinks
}
