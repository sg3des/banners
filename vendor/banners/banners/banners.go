package banners

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//AppendCallback is type of callback function for append new banner to storage
type AppendCallback func(*Banner)

//Banner structure contains basic fields describe banners in rotation service
type Banner struct {
	//URL to banner
	URL string

	//Count of prepaid shows
	Count int64

	//Available categories
	Categories []string

	//Lock flag 1 or 0
	Lock int32
}

//LoadBanners read csv file and load banners from it to rotation service
func LoadBanners(filename string, separator string, add AppendCallback) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed open file '%s' with banners list, reason: %s", filename, err)
	}
	defer f.Close()

	// banners.Lock()
	// defer banners.Unlock()

	// if banners.categories == nil {
	// 	banners.categories = make(map[string][]*banners.Banner)
	// }

	s := bufio.NewScanner(f)
	for i := 0; s.Scan(); i++ {
		csvLine := s.Text()

		//skip empty and commented lines
		if len(csvLine) == 0 || csvLine[0] == '#' {
			continue
		}

		b, err := parseBannerCSV(csvLine, separator)
		if err != nil {
			log.Println("error while parse csv file %s:%d, ", err, i)
			continue
		}

		// log.Println(b)
		add(b)

		// appendBanner(b)
	}

	// if len(banners.categories) == 0 {
	// 	return fmt.Errorf("in file '%s' banners not founded")
	// }

	return nil
}

func parseBannerCSV(csvLine string, separator string) (*Banner, error) {
	ss := strings.Split(csvLine, separator)

	if len(ss) < 3 {
		return nil, errors.New("not enough fields")
	}

	count, err := strconv.Atoi(ss[1])
	if err != nil {
		return nil, fmt.Errorf("incorrect count value, %s", err)
	}

	b := &Banner{
		URL:        ss[0],
		Count:      int64(count),
		Categories: ss[2:],
	}

	return b, nil
}

func (b *Banner) ContainCategory(cat string) bool {
	for _, s := range b.Categories {
		if s == cat {
			return true
		}
	}

	return false
}
