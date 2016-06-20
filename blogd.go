package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
	"flag"
)

var mtime_map = map[int]string{}
var mtime_int = []int{}
var blog_html string

func blogMain() {
	for {
		headers := ""

		for k := range mtime_map {
			delete(mtime_map, k)
		}
		mtime_int = mtime_int[:0]

		header, err := ioutil.ReadFile("./header.html")
		if err != nil {
			log.Fatal("No header file?")
		}

		pretext, err := ioutil.ReadFile("./pretext.html")
		if err != nil {
			log.Fatal("No pretext file?")
		}

		posttext, err := ioutil.ReadFile("./posttext.html")
		if err != nil {
			log.Fatal("No posttext file?")
		}

		files, err := ioutil.ReadDir("./posts/")
		if err != nil {
			log.Fatal("No posts directory?")
		}

		for _, f := range files {
			nm, err := os.Stat("./posts/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			nixtime := int(nm.ModTime().Unix())
			mtime_map[nixtime] = "./posts/" + f.Name()
			mtime_int = append(mtime_int, nixtime)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(mtime_int)))
		headers = string(header)
		for k := 0; k < len(mtime_int); k++ {
			postsdat, err := ioutil.ReadFile(mtime_map[mtime_int[k]])
			postsdat = blackfriday.MarkdownCommon(postsdat)
			if err != nil {
				log.Fatal(err)
			}
			headers = headers + string(pretext) + string(postsdat) + string(posttext)
		}
		footer, err := ioutil.ReadFile("./footer.html")
		if err != nil {
			log.Fatal("No footer file?")
		}
		blog_html = headers + string(footer)
		time.Sleep(10 * time.Second)
	}
}

func blogHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, blog_html)
}

func main() {
	var ListenIP string
	flag.StringVar(&ListenIP, "listen", "127.0.0.1:8080", "The IP to listen on.") 
	flag.Parse()

	assets, err := filepath.Abs("assets")

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assets))))

	go blogMain()
	http.HandleFunc("/", blogHandler)
	log.Fatal(http.ListenAndServe(ListenIP, nil))
}
