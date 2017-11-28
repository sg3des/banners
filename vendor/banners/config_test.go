package main

import (
	"io/ioutil"
	"log"
	"testing"
)

var testconfigfile = "testdata/banners.conf"
var testconfigdata = `
http-addr = ":8080"
csv-file = "./testdata/banners.csv"
csv-separator = ";"
`

func init() {
	log.SetFlags(log.Lshortfile)

	ioutil.WriteFile(testconfigfile, []byte(testconfigdata), 0644)
}

func TestLoadConfig(t *testing.T) {
	err := LoadConfig("testdata/banners.conf")
	if err != nil {
		t.Error(err)
	}

	err = LoadConfig("")
	if err != nil {
		t.Error(err)
	}
}
