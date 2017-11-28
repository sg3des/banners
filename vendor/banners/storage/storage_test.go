package storage

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

var bannersfile = "testdata/banners.csv"

var bannerscsv = []string{
	`https://octodex.github.com/images/dinotocat.png;2;flight;trains`,
	`https://octodex.github.com/images/dojocat.jpg;2;flight`,
	`https://octodex.github.com/images/orderedlistocat.png;2;flight`,
	`https://octodex.github.com/images/gracehoppertocat.jpg;1;show;britain;benny hill;sketches;tv`,
	`https://octodex.github.com/images/jetpacktocat.png;1;games;minecraft;blocks;sandbox`,
	`https://octodex.github.com/images/steroidtocat.png;1;onlycategory`,
}

func init() {
	log.SetFlags(log.Lshortfile)

	ioutil.WriteFile(bannersfile, []byte(strings.Join(bannerscsv, "\n")), 0644)
}

func TestLoadBanner(t *testing.T) {
	err := LoadBanners(bannersfile, ";")
	if err != nil {
		t.Fatal(err)
	}

	if banners.count != len(bannerscsv) {
		t.Error("count of banners not equal")
	}
}

func TestLookupBanner(t *testing.T) {
	_, ok := LookupBanner([]string{"flight"})
	if !ok {
		t.Error("banner not found by existing category flight")
	}

	_, ok = LookupBanner([]string{"flight", "show"})
	if !ok {
		t.Error("banner not found by existing categories flight and show")
	}

	_, ok = LookupBanner([]string{"onlycategory"})
	if !ok {
		t.Error("banner not found by existing category onlycategory")
	}

	_, ok = LookupBanner([]string{"onlycategory"})
	if ok {
		t.Error("banner should not be founded")
	}

	_, ok = LookupBanner(nil)
	if !ok {
		t.Error("banner not found by any category")
	}
}
