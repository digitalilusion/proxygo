package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var result *ParseResult
var injections []*html.Node

func rootHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Request: %s", r.URL.Path)
	uri := r.URL.Path[1:]
	url, err := url.ParseRequestURI("http://" + uri)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	// Make the request
	resp, err := http.Get(url.String())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	contentType := resp.Header.Get("Content-type")
	if strings.Contains(contentType, "/html") {
		result = ParseHTML(resp.Body, r.URL.Path[1:])
		if result == nil {
			w.WriteHeader(400)
			w.Write([]byte("Can't parse HTML"))
			return
		}
		w.Header().Set("Content-type", contentType)
		if injections != nil && len(injections) > 0 {
			var fInjection func(*html.Node)
			fInjection = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "body" {
					for _, k := range injections {
						k.Parent = nil
						k.NextSibling = nil
						n.AppendChild(k)
						//log.Print(k)
					}
					return
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					fInjection(c)
				}
			}
			fInjection(result.dom)
		}
		html.Render(w, result.dom)
	} else {
		w.Header().Set("Content-type", contentType)
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
		resp.Body.Close()
	}

}

func init() {
	InitCache()
	// Check for injections
	_, err := os.Stat("injects")
	if err != nil {
		return
	}
	files, _ := ioutil.ReadDir("injects")
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".html") {
			continue
		}
		f, err := os.Open("injects/" + f.Name())
		if err != nil {
			log.Printf("Can't open inject file: %s - %s", f.Name(), err.Error())
			continue
		}
		els, err := html.ParseFragment(f, nil)
		if err != nil {
			log.Printf("Can't open inject file: %s - %s", f.Name(), err.Error())
			continue
		}
		for _, inj := range els {

			var findBody func(*html.Node)
			findBody = func(n *html.Node) {
				if n.Type == html.ElementNode && (n.Data == "head" || n.Data == "body") {
					for cc := n.FirstChild; cc != nil && strings.TrimSpace(cc.Data) != ""; cc = cc.NextSibling {
						cc.Parent = nil
						cc.NextSibling = nil
						injections = append(injections, cc)
					}
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					findBody(c)
				}
			}
			findBody(inj)

		}

	}

}

func main() {
	log.Print("Proxy Go")
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8080", nil)
}
