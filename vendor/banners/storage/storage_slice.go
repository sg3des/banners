package storage

import (
	"banners/banners"
	"math/rand"
	"sync/atomic"
)

//SliceStorage contains all banners in slice
type SliceStorage struct {
	list []*banners.Banner
}

//NewSliceStorage initialize new instance of slice storage
func NewSliceStorage() *SliceStorage {
	return new(SliceStorage)
}

//AppendBanner add banner to slice
func (bs *SliceStorage) AppendBanner(b *banners.Banner) {
	bs.list = append(bs.list, b)
}

//LookupBanner lookup banner by categories
func (bs *SliceStorage) LookupBanner(categories []string) (b *banners.Banner, ok bool) {

	//x is random position from list
	x := rand.Intn(len(bs.list))

	//lookup banner from x to end
	b, ok = bs.lookupBanner(bs.list[x:], categories)
	if ok {
		atomic.AddInt64(&b.Count, -1)
		return b, ok
	}

	//lookup banner from start to x
	b, ok = bs.lookupBanner(bs.list[:x], categories)
	if ok {
		atomic.AddInt64(&b.Count, -1)
		return b, ok
	}

	return nil, false
}

func (bs *SliceStorage) lookupBanner(list []*banners.Banner, categories []string) (*banners.Banner, bool) {

	for _, b := range list {
		if atomic.LoadInt64(&b.Count) < 1 && atomic.LoadInt32(&b.Lock) == 1 {
			continue
		}

		atomic.StoreInt32(&b.Lock, 1)
		ok := bs.suitBanner(b, categories)
		atomic.StoreInt32(&b.Lock, 0)

		if ok {
			return b, true
		}
	}

	return nil, false
}

func (bs *SliceStorage) suitBanner(b *banners.Banner, categories []string) bool {
	if len(categories) == 0 || len(b.Categories) == 0 {
		return true
	}

	//x is random position from slice of categories
	x := rand.Intn(len(categories))

	//lookup banner with category from x to end
	for _, cat := range categories[x:] {
		if b.ContainCategory(cat) {
			return true
		}
	}

	//lookup banner with category from start to x
	for _, cat := range categories[:x] {
		if b.ContainCategory(cat) {
			return true
		}
	}

	return false
}

func (bs *SliceStorage) GetCount() int {
	return len(bs.list)
}
