package storage

import (
	"banners/banners"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

var (
	bannersfile = "testdata/banners.csv"
	bannerscsv  = []string{
		`https://octodex.github.com/images/dinotocat.png;2;flight;trains`,
		`https://octodex.github.com/images/dojocat.jpg;2;flight`,
		`https://octodex.github.com/images/orderedlistocat.png;2;flight`,
		`https://octodex.github.com/images/gracehoppertocat.jpg;1;show;britain;benny hill;sketches;tv`,
		`https://octodex.github.com/images/jetpacktocat.png;1;games;minecraft;blocks;sandbox`,
		`https://octodex.github.com/images/steroidtocat.png;1;onlycategory`,
	}

	testLock *LockStorage
	testChan *ChanStorage
)

func init() {
	log.SetFlags(log.Lshortfile)
	ioutil.WriteFile(bannersfile, []byte(strings.Join(bannerscsv, "\n")), 0644)

	testLock = NewLockStorage()
	testChan = NewChanStorage()
}

//
//LOCK STORAGE
//

func TestLockLoadBanners(t *testing.T) {
	err := banners.LoadBanners(bannersfile, ";", testLock.AppendBanner)
	if err != nil {
		t.Fatal(err)
	}

	if testLock.GetCount() != len(bannerscsv) {
		t.Error("count of banners not equal")
	}
}

func TestLockBanner(t *testing.T) {
	_, ok := testLock.LookupBanner([]string{"flight"})
	if !ok {
		t.Error("banner not found by existing category flight")
	}

	_, ok = testLock.LookupBanner([]string{"flight", "show"})
	if !ok {
		t.Error("banner not found by existing categories flight and show")
	}

	_, ok = testLock.LookupBanner([]string{"onlycategory"})
	if !ok {
		t.Error("banner not found by existing category onlycategory")
	}

	_, ok = testLock.LookupBanner([]string{"onlycategory"})
	if ok {
		t.Error("banner should not be founded")
	}

	_, ok = testLock.LookupBanner(nil)
	if !ok {
		t.Error("banner not found by any category")
	}
}

//
//CHAN STORAGE
//

func TestChanLoadBanners(t *testing.T) {
	err := banners.LoadBanners(bannersfile, ";", testChan.AppendBanner)
	if err != nil {
		t.Fatal(err)
	}

	if testChan.GetCount() != len(bannerscsv) {
		t.Error("count of banners not equal")
	}
}

func TestChanBanner(t *testing.T) {
	_, ok := testChan.LookupBanner([]string{"flight"})
	if !ok {
		t.Error("banner not found by existing category flight")
	}

	_, ok = testChan.LookupBanner([]string{"flight", "show"})
	if !ok {
		t.Error("banner not found by existing categories flight and show")
	}

	_, ok = testChan.LookupBanner([]string{"onlycategory"})
	if !ok {
		t.Error("banner not found by existing category onlycategory")
	}

	_, ok = testChan.LookupBanner([]string{"onlycategory"})
	if ok {
		t.Error("banner should not be founded")
	}

	_, ok = testChan.LookupBanner(nil)
	if !ok {
		t.Error("banner not found by any category")
	}
}
