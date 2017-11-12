package main

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

// ParseResult Link result after parsed the HTML into DOM
type ParseResult struct {
	script []string
	link   []string
	img    []string
	dom    *html.Node
}

func processLink(n *html.Node, r *ParseResult, host *string) {
	var key string
	var attr *html.Attribute
	var kind string = n.Data

	if n.Data == "img" || n.Data == "script" {
		key = "src"

	} else if n.Data == "a" || n.Data == "link" {
		key = "href"
	} else {
		return
	}
	for i := range n.Attr {
		if n.Attr[i].Key == key {
			attr = &n.Attr[i]
			break
		}
	}
	if attr != nil {
		if kind == "a" {
			if (!strings.HasPrefix(attr.Val, "//")) && strings.HasPrefix(attr.Val, "/") {
				attr.Val = "http://" + *host + attr.Val
			}
		} else if kind == "link" || kind == "img" || kind == "script" {
			if !(strings.HasPrefix(attr.Val, "http://") || strings.HasPrefix(attr.Val, "https://") || strings.HasPrefix(attr.Val, "//")) {
				sep := ""
				if !strings.HasPrefix(attr.Val, "/") {
					sep = "/"
				}
				attr.Val = "/" + *host + sep + attr.Val
			}
		}
	}
}
func domHandle(n *html.Node, r *ParseResult, host *string) {

	if n.Type == html.ElementNode {
		switch n.Data {
		case "img", "a", "script", "link":
			processLink(n, r, host)
			break
		}

	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		domHandle(c, r, host)
	}
}

// ParseHTML Parse the HTML and return the result
func ParseHTML(r io.ReadCloser, host string) *ParseResult {
	var result ParseResult
	dom, err := html.Parse(r)
	if err != nil {
		log.Printf("Can't parse HTML: %s", err.Error())
		return nil
	}
	result.dom = dom
	domHandle(dom, &result, &host)
	return &result
}
