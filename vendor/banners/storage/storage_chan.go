package storage

import (
	"banners/banners"
	"math/rand"
	"sync/atomic"
)

type ChanStorage struct {
	categories map[string][]*banners.Banner
	count      int64

	chAppend    chan (*banners.Banner)
	chLookupIn  chan ([]string)
	chLookupOut chan (*banners.Banner)
}

func NewChanStorage() *ChanStorage {
	bs := new(ChanStorage)
	bs.categories = make(map[string][]*banners.Banner)

	bs.chAppend = make(chan (*banners.Banner))
	bs.chLookupIn = make(chan ([]string))
	bs.chLookupOut = make(chan (*banners.Banner))

	go bs.listener()

	return bs
}

func (bs *ChanStorage) listener() {
	for {
		select {
		case b := <-bs.chAppend:
			bs.appendBanner(b)
		case categoriess := <-bs.chLookupIn:
			b := bs.lookupBanner(categoriess)
			bs.chLookupOut <- b
		}
	}
}

func (bs *ChanStorage) appendBanner(b *banners.Banner) {
	for _, category := range b.Categories {
		bs.categories[category] = append(bs.categories[category], b)
	}
	bs.count++
}

//AppendBanner to categories
func (bs *ChanStorage) AppendBanner(b *banners.Banner) {
	bs.chAppend <- b
}

//removeBanner from all categories
func (bs *ChanStorage) removeBanner(b *banners.Banner) {
	for _, cat := range b.Categories {
		list := bs.categories[cat]

		for i, ban := range list {
			if ban == b {
				list[i] = list[len(list)-1]
				list = list[:len(list)-1]

				if len(list) == 0 {
					delete(bs.categories, cat)
				} else {
					bs.categories[cat] = list
				}

				break
			}
		}
	}
}

func (bs *ChanStorage) lookupBanner(categories []string) (b *banners.Banner) {
	category, ok := bs.lookupCategory(categories)
	if !ok {
		return nil
	}

	list, ok := bs.categories[category]
	if !ok {
		return nil
	}

	i := rand.Intn(len(list))

	b = list[i]
	b.Count--
	if b.Count == 0 {
		bs.removeBanner(b)
	}

	return
}

//LookupBanner lookup banner by category, return nil,false if banner with count > 0 not found
func (bs *ChanStorage) LookupBanner(categories []string) (b *banners.Banner, ok bool) {

	bs.chLookupIn <- categories
	b = <-bs.chLookupOut
	if b != nil {
		ok = true
	}

	return
}

//GetCount returned total count of banners in rotation
func (bs *ChanStorage) GetCount() int {
	return int(atomic.LoadInt64(&bs.count))
}

//LookupCategory lookup some category from request category list and avaliable in rotation service
func (bs *ChanStorage) lookupCategory(categories []string) (cat string, ok bool) {
	if len(categories) == 0 {
		return bs.lookupRandomCategory()
	}

	var list []*banners.Banner

	for _, cat = range categories {
		list, ok = bs.categories[cat]
		if !ok {
			continue
		}

		if ok = checkList(list); ok {
			break
		}
	}

	return
}

//LookupRandomCategory lookup random category from available in rotation service
func (bs *ChanStorage) lookupRandomCategory() (cat string, ok bool) {
	var list []*banners.Banner

	for cat, list = range bs.categories {
		if ok = checkList(list); ok {
			break
		}
	}

	return
}
