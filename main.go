package main

import (
  "log"
  "net/http"
  "net/url"
  "io/ioutil"
)

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
  data, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    w.WriteHeader(400)
    w.Write([]byte(err.Error()))
    return
  }
  w.Write(data)

}

func main()  {
  log.Print("Proxy Go")
  http.HandleFunc("/", rootHandler)
  http.ListenAndServe(":8080", nil)
}
