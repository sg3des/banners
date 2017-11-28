package main

import (
	"fmt"
	"log"
	"net/http"

	"banners/storage"

	"github.com/gorilla/schema"
	"github.com/sg3des/argum"
)

var args struct {
	Config string `argum:"--config" help:"path to configuration file, toml format"`
	Debug  bool   `argum:"--debug" help:"enable debug mode"`
}

func init() {
	argum.Version = "0.1.a171128"
	argum.MustParse(&args)

	if args.Debug {
		log.SetFlags(log.Lshortfile)
	}
}

func main() {
	if args.Debug {
		fmt.Println("Load configuragion")
	}
	if err := LoadConfig(args.Config); err != nil {
		log.Fatalln(err)
	}

	if args.Debug {
		fmt.Println("Load banners from:", Config.CSVFile)
	}
	if err := storage.LoadBanners(Config.CSVFile, Config.CSVSeparator); err != nil {
		log.Fatalln(err)
	}
	if args.Debug {
		fmt.Println("Loaded", storage.GetCount())
	}

	if args.Debug {
		fmt.Println("Enable web server on:", Config.HTTPAddr)
	}
	if err := enableServer(Config.HTTPAddr); err != nil {
		log.Fatalln(err)
	}
}

func enableServer(addr string) error {
	http.HandleFunc("/", index)
	http.HandleFunc("/reload", reload)

	return http.ListenAndServe(addr, nil)
}

var indexDecoder = schema.NewDecoder()

type Request struct {
	Categories []string `schema:"category[]"`
}

func index(w http.ResponseWriter, r *http.Request) {
	var req Request

	err := indexDecoder.Decode(&req, r.URL.Query())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	b, ok := storage.LookupBanner(req.Categories)
	if !ok {
		err = fmt.Errorf("banner not found")
		log.Println(err)
		http.Error(w, err.Error(), 402)
		return
	}

	if args.Debug {
		log.Println("request banner by", r.RemoteAddr, b.URL, b.Count)
	}

	fmt.Fprintf(w, "<!DOCTYPE html><img src='%s' width=200>", b.URL)
}

func reload(w http.ResponseWriter, r *http.Request) {
	storage.LoadBanners(Config.CSVFile, Config.CSVSeparator)
}
