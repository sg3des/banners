package main

import (
	"banners/banners"
	"banners/storage"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/sg3des/argum"
)

var bs storage.Storage

//Args is global struct contains current configuration
var args struct {
	Config string `toml:"-" argum:"-c,--config" help:"path to configuration file, toml format`
	Debug  bool   `toml:"debug" argum:"-d,--debug" help:"enable debug mode"`

	HTTPAddr     string `toml:"http-addr" help:"listening address"`
	CSVFile      string `toml:"csv-file" help:"path to file with banners, CSV format"`
	CSVSeparator string `toml:"csv-separator" help:"field seperator" default:";"`
	StorageType  string `toml:"storage-type" argum:"lock|chan|slice" help:"select mechanism type of storage"`
}

func init() {
	argum.Version = "0.2.a171205"
	argum.MustParse(&args)

	if args.Debug {
		log.SetFlags(log.Lshortfile)
	}

	if err := LoadConfig(args.Config); err != nil {
		log.Fatalln(err)
	}

	//overwrite values config values with arguments
	argum.MustParse(&args)
}

func main() {
	debugOutput("Initialize storage:", args.StorageType)
	switch args.StorageType {
	case "chan":
		bs = storage.NewChanStorage()
	case "lock":
		bs = storage.NewLockStorage()
	case "slice":
		bs = storage.NewSliceStorage()
	default:
		log.Fatalln("unexpected storage type")
	}

	debugOutput("Load banners from:", args.CSVFile)
	if err := banners.LoadBanners(args.CSVFile, args.CSVSeparator, bs.AppendBanner); err != nil {
		log.Fatalln(err)
	}
	debugOutput("Loaded banners:", bs.GetCount())

	debugOutput("Enable web server on:", args.HTTPAddr)
	if err := enableServer(args.HTTPAddr); err != nil {
		log.Fatalln(err)
	}
}

func debugOutput(ss ...interface{}) {
	if args.Debug {
		for _, s := range ss {
			fmt.Print(s, " ")
		}
		fmt.Println()
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

	b, ok := bs.LookupBanner(req.Categories)
	if !ok {
		err = fmt.Errorf("banner not found")
		log.Println(err)
		http.Error(w, err.Error(), 402)
		return
	}

	debugOutput("request banner by", r.RemoteAddr, b.URL, b.Count)

	fmt.Fprintf(w, "<!DOCTYPE html><img src='%s' width=200>", b.URL)
}

func reload(w http.ResponseWriter, r *http.Request) {
	banners.LoadBanners(args.CSVFile, args.CSVSeparator, bs.AppendBanner)
}
