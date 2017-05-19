package main

import (
    "fmt"
    //"html"
    "log"
    "net/http"
	"io/ioutil"
	"strings"
)

type viewHandler struct{}

var (
	head string = `<!DOCTYPE html>
<html>
<head>
<title>Page Title</title>
</head>
<body>

<h1>This is a Heading</h1>`
	tail string =	`</body></html>`
)

func main() {
	http.HandleFunc("/List", List)
	http.Handle("/", http.FileServer(http.Dir("./Vids/")))
	//http.Handle("/", new(viewHandler))
	
    //http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    //    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    //})


    log.Fatal(http.ListenAndServe(":8080", nil))

}

func List(w http.ResponseWriter, r *http.Request){
	fileText:=""
	files, err := ioutil.ReadDir("./Vids")
	if(err != nil){
		fileText = `<p>No Files Found</p>`
	}
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name()), "mp4"){
			fileText += ("<video width=\"400\" controls><source src=\""+file.Name()+"\" type=\"video/mp4\">Your browser does not support HTML5 video.</video>")
		}else{
			fileText += (`<p>` + file.Name() + `</p>`)
		}
		
	}
	fmt.Fprintf(w, head + fileText + tail)
}