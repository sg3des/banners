package storage

import "banners/banners"

type Storage interface {
	AppendBanner(*banners.Banner)
	LookupBanner(categories []string) (*banners.Banner, bool)
	GetCount() int
}

func checkList(list []*banners.Banner) bool {
	for _, b := range list {
		if b.Count > 0 {
			return true
		}
	}

	return false
}
