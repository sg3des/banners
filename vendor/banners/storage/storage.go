package storage

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

//Banner structure contains basic fields describe banners in rotation service
type Banner struct {
	//URL to banner
	URL string

	//Count of prepaid shows
	Count int

	//Available categories
	Categories []string
}

//banners is global storage contains all banners in rotation
var banners struct {
	sync.Mutex
	categories map[string][]*Banner
	count      int
}

//LoadBanners read csv file and load banners from it to rotation service
func LoadBanners(filename string, separator string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed open file '%s' with banners list, reason: %s", filename, err)
	}
	defer f.Close()

	banners.Lock()
	defer banners.Unlock()

	if banners.categories == nil {
		banners.categories = make(map[string][]*Banner)
	}

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

		appendBanner(b)
	}

	if len(banners.categories) == 0 {
		return fmt.Errorf("in file '%s' banners not founded")
	}

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
		Count:      count,
		Categories: ss[2:],
	}

	return b, nil
}

func appendBanner(b *Banner) {
	for _, category := range b.Categories {
		banners.categories[category] = append(banners.categories[category], b)
	}
	banners.count++
}

//LookupBanner lookup banner by category, return nil,false if banner with count > 0 not found
func LookupBanner(categories []string) (b *Banner, ok bool) {
	category, ok := LookupCategory(categories)
	if !ok {
		return nil, false
	}

	banners.Lock()
	list, ok := banners.categories[category]
	banners.Unlock()

	if !ok {
		return nil, false
	}

	i := rand.Intn(len(list))

	b = list[i]
	b.Count--
	if b.Count == 0 {
		removeBanner(category, list, i)
	}
	ok = true

	return
}

func removeBanner(cat string, list []*Banner, i int) {
	banners.Lock()
	list[i] = list[len(list)-1]
	list = list[:len(list)-1]
	if len(list) == 0 {
		delete(banners.categories, cat)
	} else {
		banners.categories[cat] = list
	}
	banners.Unlock()
}

//GetCount returned total count of banners in rotation
func GetCount() int {
	return banners.count
}

//LookupCategory lookup some category from request category list and avaliable in rotation service
func LookupCategory(categories []string) (cat string, ok bool) {
	if len(categories) == 0 {
		return LookupRandomCategory()
	}

	var list []*Banner

	banners.Lock()
	for _, cat = range categories {
		list, ok = banners.categories[cat]
		if !ok {
			continue
		}

		if ok = checkList(list); ok {
			break
		}
	}
	banners.Unlock()

	return
}

//LookupRandomCategory lookup random category from available in rotation service
func LookupRandomCategory() (cat string, ok bool) {
	var list []*Banner

	banners.Lock()
	for cat, list = range banners.categories {
		if ok = checkList(list); ok {
			break
		}
	}
	banners.Unlock()

	return
}

func checkList(list []*Banner) bool {
	for _, b := range list {
		if b.Count > 0 {
			return true
		}
	}

	return false
}
