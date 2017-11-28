package main

import (
	"banners/storage"
	"net/http"
	"testing"
	"time"
)

var testaddr = "127.0.0.1:8090"

func init() {
	storage.LoadBanners("storage/testdata/banners.csv", ";")
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
