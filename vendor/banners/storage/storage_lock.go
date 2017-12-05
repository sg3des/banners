package storage

import (
	"banners/banners"
	"math/rand"
	"sync"
	"sync/atomic"
)

//LockStorage is global structure contains all banners in rotation
type LockStorage struct {
	sync.RWMutex
	categories map[string][]*banners.Banner
	count      int64
}

//NewLockStorage initialize new instance of storage with lock mechanic
func NewLockStorage() *LockStorage {
	bs := new(LockStorage)
	bs.categories = make(map[string][]*banners.Banner)

	return bs
}

//AppendBanner to categories
func (bs *LockStorage) AppendBanner(b *banners.Banner) {
	bs.Lock()
	for _, category := range b.Categories {
		bs.categories[category] = append(bs.categories[category], b)
	}
	bs.Unlock()
	atomic.AddInt64(&bs.count, 1)
}

//removeBanner from all categories
func (bs *LockStorage) removeBanner(b *banners.Banner) {
	bs.Lock()
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
	bs.Unlock()
}

//LookupBanner lookup banner by category, return nil,false if banner with count > 0 not found
func (bs *LockStorage) LookupBanner(categories []string) (b *banners.Banner, ok bool) {
	category, ok := bs.lookupCategory(categories)
	if !ok {
		return nil, false
	}

	bs.RLock()
	list, ok := bs.categories[category]
	bs.RUnlock()

	if !ok {
		return nil, false
	}

	i := rand.Intn(len(list))

	b = list[i]
	atomic.AddInt64(&b.Count, -1)
	if atomic.LoadInt64(&b.Count) < 1 {
		bs.removeBanner(b)
	}
	ok = true

	return
}

//GetCount returned total count of banners in rotation
func (bs *LockStorage) GetCount() int {
	return int(atomic.LoadInt64(&bs.count))
}

//LookupCategory lookup some category from request category list and avaliable in rotation service
func (bs *LockStorage) lookupCategory(categories []string) (cat string, ok bool) {
	if len(categories) == 0 {
		return bs.lookupRandomCategory()
	}

	var list []*banners.Banner

	bs.RLock()
	for _, cat = range categories {
		list, ok = bs.categories[cat]
		if !ok {
			continue
		}

		if ok = checkList(list); ok {
			break
		}
	}
	bs.RUnlock()

	return
}

//LookupRandomCategory lookup random category from available in rotation service
func (bs *LockStorage) lookupRandomCategory() (cat string, ok bool) {
	var list []*banners.Banner

	bs.RLock()
	for cat, list = range bs.categories {
		if ok = checkList(list); ok {
			break
		}
	}
	bs.RUnlock()

	return
}
