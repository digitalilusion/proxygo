package main

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

var result *ParseResult

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
	/*
		dom, err := html.Parse(resp.Body)
		resp.Body.Close()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
	*/
	result = ParseHTML(resp.Body, r.URL.Path)
	if result == nil {
		w.WriteHeader(400)
		w.Write([]byte("Can't parse HTML"))
		return
	}

	/*
	  data, err := ioutil.ReadAll(resp.Body)
	  resp.Body.Close()
	  if err != nil {
	    w.WriteHeader(400)
	    w.Write([]byte(err.Error()))
	    return
	  }
	*/

	/*
	  xmlroot, err := xmlpath.ParseHTML(resp.Body)
	  if err != nil {
	    w.WriteHeader(400)
	    w.Write([]byte(err.Error()))
	    return
	  }
	  path := xmlpath.MustCompile(`//img/@src`)
	  iter := path.Iter(xmlroot)
	  for iter.Next() {

	        log.Println(iter.Node().String())
	  }
	*/
	w.Header().Set("Content-type", "text/html")
	html.Render(w, result.dom)

}

func init() {
	InitCache()
}

func main() {
	log.Print("Proxy Go")
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8080", nil)
}
