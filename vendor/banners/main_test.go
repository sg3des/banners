package main

import (
	"banners/banners"
	"banners/storage"
	"net/http"
	"testing"
	"time"
)

var testaddr = "127.0.0.1:8090"

func init() {
	bs = storage.NewTreeStorage()

	banners.LoadBanners("storage/testdata/banners.csv", ";", bs.AppendBanner)
}

func TestEnableServer(t *testing.T) {
	go func(t *testing.T) {
		err := enableServer(testaddr)
		if err != nil {
			t.Fatal(err)
		}
	}(t)
}

func TestIndex(t *testing.T) {
	time.Sleep(time.Second) // crunch ^^

	requests := []string{
		"?category[]=auto&category[]=trains",
		"?category[]=onlycategory",
		"",
	}

	for _, req := range requests {
		resp, err := http.Get("http://" + testaddr + "/" + req)
		if err != nil {
			t.Error(err)
		}

		if resp.ContentLength == 0 {
			t.Errorf("response by request %s is empty", req)
		}
	}

}

func TestIndexInvalid(t *testing.T) {
	time.Sleep(time.Second) // crunch ^^

	requests := []string{
		"?undefined=empty",
	}

	for _, req := range requests {
		resp, err := http.Get("http://" + testaddr + "/" + req)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode == 200 {
			t.Errorf("should be error by request %s", req)
		}
	}

}

//
//benchmarks

func benchmarkStorage(b *testing.B, categories []string) {
	banners.LoadBanners("storage/testdata/banners.csv", ";", bs.AppendBanner)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		bs.LookupBanner(categories)
	}
}

func BenchmarkLockStorage(b *testing.B) {
	b.StopTimer()
	bs = storage.NewLockStorage()
	categories := []string{"flight", "onlycategory"}

	benchmarkStorage(b, categories)
}

func BenchmarkLockStorageAny(b *testing.B) {
	b.StopTimer()
	bs = storage.NewLockStorage()

	benchmarkStorage(b, nil)
}

func BenchmarkTreeStorage(b *testing.B) {
	b.StopTimer()
	bs = storage.NewTreeStorage()
	categories := []string{"flight", "onlycategory"}

	benchmarkStorage(b, categories)
}

func BenchmarkTreeStorageAny(b *testing.B) {
	b.StopTimer()
	bs = storage.NewTreeStorage()

	benchmarkStorage(b, nil)
}

func BenchmarkSliceStorage(b *testing.B) {
	b.StopTimer()
	bs = storage.NewSliceStorage()
	categories := []string{"flight", "onlycategory"}

	benchmarkStorage(b, categories)
}

func BenchmarkSliceStorageAny(b *testing.B) {
	b.StopTimer()
	bs = storage.NewSliceStorage()

	benchmarkStorage(b, nil)
}
