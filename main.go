package main

import (
    "fmt"
    "net/http"
    "html/template"
)
//GOOS=linux GOARCH=arm GOARM=6 go build
//scp /Users/bazilio/Works/goserver/test pi@192.168.1.41:GoServer

type page struct {
  Title string
  Msg string
}
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love  %s!", r.URL.Path[1:])
}

func index(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-type","text/html");
  t,_:=template.ParseFiles("index.html");
  t.Execute(w,&page{Title:"Just Page",Msg: "Just Message"});
}

func main() {
    http.HandleFunc("/", index)
    http.ListenAndServe(":8080", nil)
}
